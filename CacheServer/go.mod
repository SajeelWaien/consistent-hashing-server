module github.com/sajeelwaien/consistent-hashing/cacheserver

go 1.26.2

replace github.com/sajeelwaien/consistent-hashing/bloomfilter => ../BloomFilter

replace github.com/sajeelwaien/consistent-hashing/node => ../Node

replace github.com/sajeelwaien/consistent-hashing/hashring => ../HashRing

require (
	github.com/sajeelwaien/consistent-hashing/bloomfilter v1.0.0
	github.com/sajeelwaien/consistent-hashing/hashring v1.0.0
	github.com/sajeelwaien/consistent-hashing/node v1.0.0
)

require github.com/spaolacci/murmur3 v1.1.0 // indirect
