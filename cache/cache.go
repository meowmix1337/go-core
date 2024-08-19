package cache

import "context"

type Cache interface {
	// Get retrieves data given a key
	Get(ctx context.Context, key string) (interface{}, error)

	// Set adds the value for a given key
	Set(ctx context.Context, key string, value interface{}, ttl int) error

	// Delete removes the item from the cache given the key
	Delete(ctx context.Context, key string) error

	// Purge clear all items in the cache
	Purge(ctx context.Context)

	// Size returns the number of elements in the cache
	Size(ctx context.Context) uint64
}
