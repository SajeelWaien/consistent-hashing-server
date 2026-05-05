package cacheserver

import (
	"fmt"
	"hash"
	"math/rand"

	"github.com/sajeelwaien/consistent-hashing/bloomfilter"
	"github.com/sajeelwaien/consistent-hashing/hashring"
	"github.com/sajeelwaien/consistent-hashing/node"
)

type CacheServer struct {
	nodes       []*node.Node
	hashRing    *hashring.HashRing
	bloomFilter *bloomfilter.BloomFilter
}

func NewCacheServer(hashFunction hash.Hash64, nodeCount int, replicationFactor int8, virtualNodeCount int8) *CacheServer {
	if nodeCount <= 0 {
		return nil
	}

	bloomFilter := bloomfilter.NewBloomFilter(500, bloomfilter.InitHashFunc(3, hashFunction))

	hashRing := hashring.NewHashRing(
		hashring.WithHashFunction(hashFunction),
		hashring.WithReplicationFactor(replicationFactor),
		hashring.WithVirtualNodeCount(virtualNodeCount),
		hashring.WithLoggingEnabled(true),
	)

	nodes := make([]*node.Node, nodeCount)
	for i := 0; i < nodeCount; i++ {
		identifier := fmt.Sprintf("%c_%d_node_%d", 'a'+rand.Intn(26), rand.Intn(1000), i)

		n := node.NewNode(identifier)
		nodes[i] = n
		err := hashRing.AddNode(n)
		if err != nil {
			return nil
		}
	}
	return &CacheServer{
		nodes:       nodes,
		hashRing:    hashRing,
		bloomFilter: bloomFilter,
	}
}

func (cs *CacheServer) AddRecord(key string, value string) error {
	nodes, err := cs.hashRing.GetNodesForKey(key)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		err := node.Set(key, value)
		if err != nil {
			return err
		}
	}

	cs.bloomFilter.Add(key)
	return nil

}

func (cs *CacheServer) GetRecord(key string) (string, any, error) {
	if !cs.bloomFilter.Contains(key) {
		return "", nil, fmt.Errorf("Key '%s' not found in cache (Bloom Filter)", key)
	}

	node, err := cs.hashRing.GetPrimaryNode(key)
	if err != nil {
		return "", nil, err
	}

	val, error := node.Get(key)
	if error != nil {
		return "", nil, error
	}

	// Implementation for getting a record from the cache
	return key, val, nil
}
