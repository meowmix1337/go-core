package cache

import (
	"errors"
	"sync"
	"time"
)

var (
	CacheMissErr = errors.New("cache miss")
)

// InMemoryCache is a very dumb and simple cache
type InMemoryCache struct {
	cache map[string]cacheItem
	mu    sync.Mutex
}

type cacheItem struct {
	value      interface{}
	expiration int64
}

func (i *cacheItem) isExpired() bool {
	return time.Now().After(time.Unix(i.expiration, 0))
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		cache: make(map[string]cacheItem),
	}
}

func (c *InMemoryCache) Get(key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.cache[key]
	if !found || item.isExpired() {
		return nil, CacheMissErr
	}

	return item.value, nil
}

func (c *InMemoryCache) Set(key string, value interface{}, ttl int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheItem := cacheItem{
		value:      value,
		expiration: time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}

	c.cache[key] = cacheItem

	return nil
}

func (c *InMemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
	return nil
}

func (c *InMemoryCache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]cacheItem)
}

func (c *InMemoryCache) Size() int64 {
	var size int64
	for _, item := range c.cache {
		if !item.isExpired() {
			size++
		}
	}
	return size
}
