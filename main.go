package main

import (
	"fmt"

	cacheServer "github.com/sajeelwaien/consistent-hashing/cacheserver"
	"github.com/spaolacci/murmur3"
)

func main() {
	hashFunction := murmur3.New64WithSeed(100)
	server := cacheServer.NewCacheServer(hashFunction, 3, 2, 5)

	server.AddRecord("User1", "Sajeel")
	server.AddRecord("User2", "Khawaja")
	key, val, err := server.GetRecord("User3")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Found", key, val)
}
