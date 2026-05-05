# Consistent Hashing Cache Server

A robust, distributed cache server implementation in Go utilizing consistent hashing, virtual nodes, and replication to ensure high availability and balanced load distribution.

## 🚀 Features

- **Consistent Hashing**: Minimizes reorganization of keys when nodes are added or removed.
- **Virtual Nodes**: Ensures a more uniform distribution of data across physical nodes, preventing hotspots.
- **Data Replication**: Supports replicating data across multiple nodes for fault tolerance and high availability.
- **Bloom Filter Optimization**: Uses a Bloom filter to quickly check if a key exists before performing a full cache lookup, reducing unnecessary operations.
- **Thread-Safe**: Designed with concurrency in mind using `sync.RWMutex` and `sync.Map`.

## 📁 Project Structure

```text
.
├── BloomFilter/    # Space-efficient probabilistic data structure
├── CacheServer/    # Main orchestrator and API entry point
├── HashRing/       # Consistent hashing logic and node management
├── Node/           # Individual storage unit implementation
├── main.go         # Example usage and entry point
└── go.mod          # Project dependencies
```

## 🛠️ Getting Started

### Prerequisites

- Go 1.26 or higher

### Installation

Clone the repository and install dependencies:

```bash
git clone https://github.com/sajeelwaien/consistent-hashing-server.git
cd consistent-hashing-server
go mod download
```

### Usage Example

You can run the example provided in `main.go`:

```go
package main

import (
	"fmt"
	cacheServer "github.com/sajeelwaien/consistent-hashing/cacheserver"
	"github.com/spaolacci/murmur3"
)

func main() {
	// Initialize hash function
	hashFunction := murmur3.New64WithSeed(100)
	
	// Create a new cache server: 3 nodes, replication factor of 2, 5 virtual nodes per node
	server := cacheServer.NewCacheServer(hashFunction, 3, 2, 5)

	// Add records
	server.AddRecord("User1", "Sajeel")
	server.AddRecord("User2", "Khawaja")

	// Retrieve records
	key, val, err := server.GetRecord("User1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Printf("Found Key: %s, Value: %v\n", key, val)
}
```

Run it using:
```bash
go run main.go
```

## 🧩 Components

### Hash Ring
The core component that manages the distribution of keys. It maps nodes and keys to a 64-bit hash space. By using virtual nodes, it achieves better balance even with a small number of physical servers.

### Bloom Filter
A probabilistic data structure used to test whether an element is a member of a set. In this project, it's used as a "gatekeeper" for the cache, providing a fast "no" for keys that definitely don't exist.

### Cache Node
Represents a storage unit that holds the actual data (in-memory map). It implements the `ICacheNode` interface for flexibility.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.
