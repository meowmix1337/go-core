package cache

import (
	"context"
	"sync"
)

// InMemoryCache is a very dumb and simple cache
// defer to LRU if possible unless you want pain
type InMemoryCache struct {
	cache map[string]*cacheItem
	mu    sync.Mutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]*cacheItem),
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.cache[key]
	if !found || item.isExpired() {
		return nil, CacheMissErr
	}

	return item.value, nil
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheItem := newCacheItem(key, value, int64(ttl))

	c.cache[key] = cacheItem

	return nil
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
	return nil
}

func (c *InMemoryCache) Purge(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*cacheItem)
}

func (c *InMemoryCache) Size(ctx context.Context) uint64 {
	var size uint64
	for _, item := range c.cache {
		if !item.isExpired() {
			size++
		}
	}
	return size
}
