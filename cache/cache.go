package cache

type Cache interface {
	// Get retrieves data given a key
	Get(key string) (interface{}, error)

	// Set adds the value for a given key
	Set(key string, value interface{}, ttl int) error

	// Delete removes the item from the cache given the key
	Delete(key string) error

	// Purge clear all items in the cache
	Purge()

	// Size returns the number of elements in the cache
	Size() int64
}
