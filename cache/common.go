package cache

import "time"

type cacheItem struct {
	key        string
	value      interface{}
	expiration int64
}

func newCacheItem(key string, value interface{}, ttl int64) *cacheItem {
	return &cacheItem{
		key:        key,
		value:      value,
		expiration: time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
	}
}

func (i *cacheItem) isExpired() bool {
	return time.Now().After(time.Unix(i.expiration, 0))
}
