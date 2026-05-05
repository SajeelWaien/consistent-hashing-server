package bloomfilter

import (
	"fmt"
	"hash"
)

type BloomFilter struct {
	filter   []uint8
	size     int
	hashFunc []hash.Hash64
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

func InitHashFunc(size int, hasher hash.Hash64) OptionFunc {
	hashFunc := make([]hash.Hash64, 0)
	for i := 0; i < size; i++ {
		// seed := uint64(rand.Int31())
		// hashFunc = append(hashFunc, murmur3.New32WithSeed(seed))
		hashFunc = append(hashFunc, hasher)
	}
	return func(bf *BloomFilter) {
		bf.hashFunc = hashFunc
	}
}

func hashKey(key string, bfSize int, hasher hash.Hash64) uint64 {
	hasher.Reset()
	hasher.Write([]byte(key))
	result := hasher.Sum64() % uint64(bfSize*8)
	return result
}
