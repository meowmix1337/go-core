# Cache
This is a interface to allow me to use a cache. LRU, simple, Redis or memcache

## Add the dependency
```
go get github.com/meowmix1337/go-core
```

## Purpose
1. Can connect to different caches (memcache, redis, in-memory/lru) and swap between each one with ease

## Usage
LRU Cache
```go
import "github.com/meowmix1337/go-core/cache"

// set a default capacity
cache := cache.NewLRUCache(5000)

// no need to handle err since LRU won't return an error
cache.Set(ctx, "key", "value", 60) // insert a cached item that has a time to live of 60 seconds

value, err := cache.Get(ctx, "key")
if errors.Is(err, CacheMissErr) {
    // handle cache miss
}
// use value

// if you need to purge the whole cache
cache.Purge()
```

Simple Cache
> [!CAUTION]
> This simple cache is really dumb and unbounded. Seriously, don't use this. Just use LRU if you need in-memory
```go
import "github.com/meowmix1337/go-core/cache"

// set a default capacity
cache := cache.NewInMemoryCache()

// no need to handle err since simple cache won't return an error
cache.Set(ctx, "key", "value", 60) // insert a cached item that has a time to live of 60 seconds

value, err := cache.Get(ctx, "key")
if errors.Is(err, CacheMissErr) {
    // handle cache miss
}
// use value

// if you need to purge the whole cache
cache.Purge()
```