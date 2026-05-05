package bloomfilter

import (
	"fmt"
	"hash"

	"github.com/spaolacci/murmur3"
)

type BloomFilter struct {
	filter   []uint8
	size     int
	hashFunc []hash.Hash32
}

type OptionFunc func(*BloomFilter)

func (b *BloomFilter) Add(key string) {
	for i := range len(b.hashFunc) {
		index := hashKey(key, len(b.filter), b.hashFunc[i])
		aIndex := index / 8
		bIndex := index % 8
		b.filter[aIndex] = b.filter[aIndex] | (1 << uint8(bIndex))
		// fmt.Printf("ADD: index: %d, aIndex: %d, bIndex: %d, bytes: %08b\n", index, aIndex, bIndex, b.filter[aIndex])
	}
}

func (b *BloomFilter) Contains(key string) bool {
	for i := range len(b.hashFunc) {
		index := hashKey(key, len(b.filter), b.hashFunc[i])
		aIndex := index / 8
		bIndex := index % 8
		// fmt.Printf("CONTAINS: index: %d, aIndex: %d, bIndex: %d, bytes: %08b\n", index, aIndex, bIndex, b.filter[aIndex])
		exists := b.filter[aIndex] & (1 << uint8(bIndex))
		if exists == 0 {
			return false
		}
	}
	return true
}

func (b *BloomFilter) Print(key string) {
	// fmt.Printf("%08b ", b.filter)
	for _, b := range b.filter {
		fmt.Printf("%08b ", b)
	}
	fmt.Println()
}

func NewBloomFilter(size int, opts ...OptionFunc) *BloomFilter {
	filter := &BloomFilter{
		filter:   make([]uint8, size),
		size:     size,
		hashFunc: nil,
	}

	for _, opt := range opts {
		opt(filter)
	}

	return filter
}

func InitHashFunc(size int) OptionFunc {
	hashFunc := make([]hash.Hash32, 0)
	for i := 0; i < size; i++ {
		// seed := uint32(rand.Int31())
		// hashFunc = append(hashFunc, murmur3.New32WithSeed(seed))
		hashFunc = append(hashFunc, murmur3.New32WithSeed(uint32(i)))
	}
	return func(bf *BloomFilter) {
		bf.hashFunc = hashFunc
	}
}

func hashKey(key string, bfSize int, hasher hash.Hash32) uint32 {
	hasher.Reset()
	hasher.Write([]byte(key))
	result := hasher.Sum32() % uint32(bfSize*8)
	return result
}
