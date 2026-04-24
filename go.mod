module github.com/sajeelwaien/consistent-hashing

go 1.26.2

replace github.com/sajeelwaien/consistent-hashing/bloomfilter => ./BloomFilter

require (
	github.com/google/uuid v1.6.0
	github.com/sajeelwaien/consistent-hashing/bloomfilter v1.0.0
)

require github.com/spaolacci/murmur3 v1.1.0 // indirect
