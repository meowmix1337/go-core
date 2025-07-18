# cache

This package provides a flexible caching solution for Go applications, offering both in-memory and Redis-backed implementations. It defines a common `Cache` interface, allowing you to easily swap caching strategies based on your application's needs.

## Features

- **`Cache` Interface**: A unified interface for caching operations (`Get`, `Set`, `Delete`, `Purge`, `Size`).
- **`InMemoryCache`**: A simple, thread-safe in-memory cache suitable for small-scale caching or testing.
- **`LRUCache`**: A thread-safe in-memory cache that implements the Least Recently Used (LRU) eviction policy, ideal for managing a fixed-size cache.
- **`RedisCache`**: A cache implementation backed by Redis, providing persistent and distributed caching capabilities.
- **TTL (Time-To-Live)**: All cache implementations support setting a time-to-live for cached items, ensuring data freshness.
- **Error Handling**: Provides a `CacheMissErr` for when a requested key is not found in the cache.

## Installation

To use this package, you need to have Go installed.

```bash
go get github.com/meowmix1337/go-core/cache
```

If you plan to use `RedisCache`, you also need to install the Redis Go client:

```bash
go get github.com/redis/go-redis/v9
```

## Usage

### `Cache` Interface

All cache implementations adhere to the `Cache` interface:

```go
type Cache interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Set(ctx context.Context, key string, value interface{}, ttl int) error
	Delete(ctx context.Context, key string) error
	Purge(ctx context.Context)
	Size(ctx context.Context) uint64
}
```

### `InMemoryCache`

The `InMemoryCache` is a basic, unbounded cache.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/meowmix1337/go-core/cache"
)

func main() {
	inMemoryCache := cache.NewInMemoryCache()

	// Set a value with a TTL of 5 seconds
	err := inMemoryCache.Set(context.Background(), "myKey", "myValue", 5)
	if err != nil {
		log.Fatalf("Failed to set cache: %v", err)
	}
	fmt.Println("Set 'myKey' to 'myValue' with 5s TTL.")

	// Get the value
	val, err := inMemoryCache.Get(context.Background(), "myKey")
	if err != nil {
		log.Fatalf("Failed to get cache: %v", err)
	}
	fmt.Printf("Got 'myKey': %v\n", val)

	// Wait for TTL to expire
	time.Sleep(6 * time.Second)

	// Try to get again (should be a cache miss)
	val, err = inMemoryCache.Get(context.Background(), "myKey")
	if err == cache.CacheMissErr {
		fmt.Println("'myKey' is now a cache miss as expected.")
	} else if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	} else {
		fmt.Printf("Unexpectedly got 'myKey': %v (should have expired)\n", val)
	}

	// Delete a key
	inMemoryCache.Set(context.Background(), "anotherKey", "anotherValue", 60)
	fmt.Printf("Cache size before delete: %d\n", inMemoryCache.Size(context.Background()))
	err = inMemoryCache.Delete(context.Background(), "anotherKey")
	if err != nil {
		log.Fatalf("Failed to delete cache: %v", err)
	}
	fmt.Printf("Cache size after delete: %d\n", inMemoryCache.Size(context.Background()))

	// Purge the cache
	inMemoryCache.Set(context.Background(), "key1", "val1", 60)
	inMemoryCache.Set(context.Background(), "key2", "val2", 60)
	fmt.Printf("Cache size before purge: %d\n", inMemoryCache.Size(context.Background()))
	inMemoryCache.Purge(context.Background())
	fmt.Printf("Cache size after purge: %d\n", inMemoryCache.Size(context.Background()))
}
```

### `LRUCache`

The `LRUCache` evicts the least recently used items when its capacity is reached.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/meowmix1337/go-core/cache"
)

func main() {
	// Create an LRU cache with a capacity of 3
	lruCache := cache.NewLRUCache(3)

	lruCache.Set(context.Background(), "keyA", "valueA", 60) // A
	lruCache.Set(context.Background(), "keyB", "valueB", 60) // B, A
	lruCache.Set(context.Background(), "keyC", "valueC", 60) // C, B, A

	fmt.Printf("LRU Cache size: %d\n", lruCache.Size(context.Background()))

	// Access keyA, making it most recently used
	_, err := lruCache.Get(context.Background(), "keyA") // A, C, B
	if err != nil {
		log.Fatalf("Failed to get keyA: %v", err)
	}
	fmt.Println("Accessed keyA.")

	// Add a new item, which should evict keyB (least recently used)
	lruCache.Set(context.Background(), "keyD", "valueD", 60) // D, A, C

	fmt.Printf("LRU Cache size after adding keyD: %d\n", lruCache.Size(context.Background()))

	// Try to get keyB (should be a cache miss)
	_, err = lruCache.Get(context.Background(), "keyB")
	if err == cache.CacheMissErr {
		fmt.Println("keyB was evicted as expected.")
	} else if err != nil {
		log.Fatalf("Unexpected error getting keyB: %v", err)
	} else {
		fmt.Printf("Unexpectedly got keyB: %v\n", err)
	}
}
```

### `RedisCache`

The `RedisCache` connects to a Redis server for distributed caching.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/meowmix1337/go-core/cache"
)

func main() {
	// Replace with your Redis server address, password, and DB number
	redisAddr := os.Getenv("REDIS_ADDR")     // e.g., "localhost:6379"
	redisPassword := os.Getenv("REDIS_PASS") // e.g., "" (no password)
	redisDB := 0                             // Redis database number

	if redisAddr == "" {
		log.Fatal("REDIS_ADDR environment variable not set.")
	}

	redisCache, err := cache.NewRedisCache(redisAddr, redisPassword, redisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis.")

	// Set a value with a TTL of 10 seconds
	err = redisCache.Set(context.Background(), "redisKey", "redisValue", 10)
	if err != nil {
		log.Fatalf("Failed to set Redis cache: %v", err)
	}
	fmt.Println("Set 'redisKey' to 'redisValue' with 10s TTL.")

	// Get the value
	val, err := redisCache.Get(context.Background(), "redisKey")
	if err != nil {
		log.Fatalf("Failed to get Redis cache: %v", err)
	}
	fmt.Printf("Got 'redisKey': %v\n", val)

	// Delete a key
	err = redisCache.Delete(context.Background(), "redisKey")
	if err != nil {
		log.Fatalf("Failed to delete Redis cache: %v", err)
	}
	fmt.Println("Deleted 'redisKey'.")

	// Try to get again (should be a cache miss)
	_, err = redisCache.Get(context.Background(), "redisKey")
	if err == cache.CacheMissErr {
		fmt.Println("'redisKey' is now a cache miss as expected.")
	} else if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}

	// Purge the entire Redis database (use with caution in production!)
	// redisCache.Purge(context.Background())
	// fmt.Println("Redis cache purged.")
}
```

## Contributing

Feel free to open issues or pull requests if you have suggestions or improvements.
