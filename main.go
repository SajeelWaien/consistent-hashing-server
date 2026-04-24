package main

import (
	"fmt"

	"github.com/SajeelWaien/bloom-filter-go/bloomfilter"
	"github.com/google/uuid"
)

func main() {
	values := make([]string, 0)
	for i := 0; i < 100; i++ {
		values = append(values, uuid.New().String())
	}

	hashFunc := bloomfilter.InitHashFunc(3)
	bf := bloomfilter.NewBloomFilter(50)

	for i := 0; i < len(values); i++ {
		bf.Add(values[i], &hashFunc)
	}

	bf.Print("test")
	fmt.Println(bf.Contains(uuid.NewString(), &hashFunc))
	fmt.Println(values[5], bf.Contains(values[5], &hashFunc))
}
