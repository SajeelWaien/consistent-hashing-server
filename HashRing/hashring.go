package hashring

import (
	"errors"
	"fmt"
	"hash"
	"slices"
	"sort"
	"sync"

	"github.com/sajeelwaien/consistent-hashing/node"
	cacheNode "github.com/sajeelwaien/consistent-hashing/node"
	"github.com/spaolacci/murmur3"
)

var (
	HashError           = errors.New("Error hashing the key")
	ExistError          = errors.New("Node already exists in the hash ring")
	NotExistError       = errors.New("Node does not exist in the hash ring")
	ErrNoNodesAvailable = errors.New("No nodes available for current key")
)

type HashRingConfig struct {
	NumberOfReplicas int8
	VirtualNodeCount int8
	LoggingEnabled   bool
	HashFunction     hash.Hash64
}

type HashRingConfigFunc func(*HashRingConfig)

func WithHashFunction(hashFunc hash.Hash64) HashRingConfigFunc {
	return func(config *HashRingConfig) {
		config.HashFunction = hashFunc
	}
}

func WithReplicationFactor(replicas int8) HashRingConfigFunc {
	return func(config *HashRingConfig) {
		config.NumberOfReplicas = replicas
	}
}

func WithVirtualNodeCount(virtualNodeCount int8) HashRingConfigFunc {
	return func(config *HashRingConfig) {
		config.VirtualNodeCount = virtualNodeCount
	}
}

func WithLoggingEnabled(enabled bool) HashRingConfigFunc {
	return func(config *HashRingConfig) {
		config.LoggingEnabled = enabled
	}
}

type HashRing struct {
	vNodes       sync.Map
	hostNodes    sync.Map
	sortedHashes []uint64
	mutex        sync.RWMutex
	config       HashRingConfig
}

func NewHashRing(configOptions ...HashRingConfigFunc) *HashRing {
	config := &HashRingConfig{
		NumberOfReplicas: 2,
		VirtualNodeCount: 4,
		LoggingEnabled:   false,
		HashFunction:     murmur3.New64WithSeed(100),
	}

	for _, opt := range configOptions {
		opt(config)
	}

	hashRing := &HashRing{
		config:       *config,
		sortedHashes: make([]uint64, 0),
		vNodes:       sync.Map{},
		hostNodes:    sync.Map{},
	}

	return hashRing
}

func (ring *HashRing) hash(key string) (uint64, error) {
	hash := ring.config.HashFunction
	hash.Reset()
	_, err := hash.Write([]byte(key))
	if err != nil {
		return 0, err
	}
	return hash.Sum64(), nil
}

func (ring *HashRing) searchKey(h uint64) int {
	idx := sort.Search(len(ring.sortedHashes), func(i int) bool {
		return ring.sortedHashes[i] >= h
	})

	if idx == len(ring.sortedHashes) {
		idx = 0
	}

	return idx
}

func (ring *HashRing) AddNode(node cacheNode.ICacheNode) error {
	// Implementation for adding a node to the hash ring
	ring.mutex.Lock()
	defer ring.mutex.Unlock()

	nodeId := node.GetID()
	_, exists := ring.hostNodes.Load(nodeId)
	if exists {
		return ExistError
	}

	for i := 0; i < int(ring.config.VirtualNodeCount); i++ {
		hash, err := ring.hash(fmt.Sprintf("%s-%d", nodeId, i))
		if err != nil {
			return err
		}
		ring.sortedHashes = append(ring.sortedHashes, hash)
		ring.vNodes.Store(hash, node)
	}

	ring.hostNodes.Store(nodeId, true)
	slices.Sort(ring.sortedHashes)
	return nil
}

func (ring *HashRing) RemoveNode(cacheNode node.ICacheNode) error {
	nodeId := cacheNode.GetID()
	hash, err := ring.hash(nodeId)
	if err != nil {
		return HashError
	}
	_, exists := ring.hostNodes.Load(hash)
	if !exists {
		return NotExistError
	}
	ring.hostNodes.Delete(hash)
	// remove all virtual nodes
	newKeys := make([]uint64, 0, len(ring.sortedHashes))
	for _, h := range ring.sortedHashes {
		val, ok := ring.vNodes.Load(h)
		if ok && val.(node.ICacheNode).GetID() == nodeId {
			ring.vNodes.Delete(h)
			continue
		}
		newKeys = append(newKeys, h)
	}
	ring.sortedHashes = newKeys
	return nil
}

func (ring *HashRing) GetPrimaryNode(key string) (node.ICacheNode, error) {
	hash, err := ring.hash(key)
	if err != nil {
		return nil, HashError
	}

	ring.mutex.RLock()
	defer ring.mutex.RUnlock()

	if len(ring.sortedHashes) == 0 {
		return nil, ErrNoNodesAvailable
	}

	idx := ring.searchKey(hash)
	vHash := ring.sortedHashes[idx]
	cacheNode, exists := ring.vNodes.Load(vHash)
	if !exists {
		return nil, ErrNoNodesAvailable
	}
	return cacheNode.(node.ICacheNode), nil

}

func (ring *HashRing) GetNodesForKey(key string) ([]node.ICacheNode, error) {
	ring.mutex.RLock()
	defer ring.mutex.RUnlock()

	if len(ring.sortedHashes) == 0 {
		return nil, ErrNoNodesAvailable
	}

	hash, err := ring.hash(key)
	if err != nil {
		return nil, HashError
	}

	seenNodes := make(map[string]bool)
	nodeList := make([]node.ICacheNode, 0, ring.config.NumberOfReplicas)

	start := ring.searchKey(hash)
	i := start

	for len(nodeList) < int(ring.config.NumberOfReplicas) {
		vHash := ring.sortedHashes[i%len(ring.sortedHashes)]
		curr, exists := ring.vNodes.Load(vHash)

		if !exists {
			i++
			continue
		}

		var n node.ICacheNode = curr.(node.ICacheNode)
		identifier := n.GetID()

		if _, exists := seenNodes[identifier]; !exists {
			seenNodes[identifier] = true
			nodeList = append(nodeList, n)
		}

		i++
		if i-start > len(ring.sortedHashes) {
			break
		}
	}

	if len(nodeList) == 0 {
		return nil, ErrNoNodesAvailable
	}

	return nodeList, nil

}
