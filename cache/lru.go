package cache

import (
	"container/list"
	"context"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

const (
	MinimumCapacity uint64 = 10
	DefaultCapacity uint64 = 500
)

type lruCache struct {
	mu        sync.RWMutex
	capacity  uint64
	cache     map[string]*list.Element
	cacheList *list.List // doubly linked list
}

func NewLRUCache(capacity uint64) *lruCache {
	if capacity < MinimumCapacity {
		log.Warn().Msg(fmt.Sprintf("minimum capacity is %v, but got %v. Defaulting to %v", MinimumCapacity, capacity, DefaultCapacity))
		capacity = DefaultCapacity
	}
	cache := &lruCache{
		capacity:  capacity,
		cache:     make(map[string]*list.Element),
		cacheList: list.New(),
	}

	return cache
}

func (c *lruCache) Get(ctx context.Context, key string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	element, found := c.cache[key]

	if !found || element == nil {
		return nil, CacheMissErr
	}

	cacheItem := element.Value.(*cacheItem)

	// delete expired
	if cacheItem.isExpired() {
		// unlock the read lock and acquire write lock
		c.mu.RUnlock()
		c.mu.Lock()
		c.removeElement(element)
		// now we need to unlock the write and re-acquire the read lock since we'll defer unlock it
		c.mu.Unlock()
		c.mu.RLock()
		return nil, CacheMissErr
	}

	// if exist, move it to the front
	// must not be expired
	// if already exists, move to front of list
	c.cacheList.MoveToFront(element)
	return cacheItem.value, nil

}

func (c *lruCache) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// if already exists, move to front of list
	if element, moved := c.moveToFront(key); moved {
		// update the expiration since we've access the existing element
		element.Value = newCacheItem(key, value, int64(ttl))
		return nil
	}

	// create new item
	newCacheItem := newCacheItem(key, value, int64(ttl))
	newElement := c.cacheList.PushFront(newCacheItem)
	c.cache[key] = newElement

	// if new size is larger than capacity, evict last element
	if c.Size(ctx) > c.capacity {
		c.evict()
	}

	return nil
}

func (c *lruCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	element := c.cache[key]
	if element != nil {
		c.cacheList.Remove(element)
		delete(c.cache, element.Value.(*cacheItem).key)
	}

	return nil
}

func (c *lruCache) Purge(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[string]*list.Element)
	c.cacheList = list.New()
}

func (c *lruCache) Size(ctx context.Context) uint64 {
	return uint64(c.cacheList.Len())
}

func (c *lruCache) moveToFront(key string) (*list.Element, bool) {
	// if already exists, move to front of list
	if element, found := c.cache[key]; found {
		c.cacheList.MoveToFront(element)
		return element, true
	}
	return nil, false
}

func (c *lruCache) evict() {
	// get the last element in the list, remove it from the list and the map
	lastElement := c.cacheList.Back()
	if lastElement != nil {
		c.removeElement(lastElement)
		return
	}
	log.Debug().Msg("last element was nil, nothing to evict")
}

func (c *lruCache) removeElement(element *list.Element) {
	delete(c.cache, element.Value.(*cacheItem).key)
	c.cacheList.Remove(element)
}
