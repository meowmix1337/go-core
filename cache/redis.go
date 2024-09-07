package cache

import (
	"context"
	"time"

	redis "github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) (*redisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Err(err).Msg("error connecting to redis")
		return nil, err
	}

	return &redisCache{
		client: client,
	}, nil
}

// Get retrieves data given a key
func (rc *redisCache) Get(ctx context.Context, key string) (interface{}, error) {
	result, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, CacheMissErr
	} else if err != nil {
		return nil, err
	}

	return result, nil
}

// Set adds the value for a given key
func (rc *redisCache) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	return rc.client.Set(ctx, key, value, time.Second*time.Duration(ttl)).Err()
}

// Delete removes the item from the cache given the key
func (rc *redisCache) Delete(ctx context.Context, key string) error {
	_, err := rc.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// Purge clear all items in the cache
func (rc *redisCache) Purge(ctx context.Context) {
	err := rc.client.FlushDB(ctx).Err()
	if err != nil {
		log.Err(err).Msg("failed to flush redis DB")
	}
}

// Size returns the number of elements in the cache
func (rc *redisCache) Size(ctx context.Context) uint64 {
	size, err := rc.client.DBSize(ctx).Result()
	if err != nil {
		log.Err(err).Msg("failed to return size of redis cache")
		return 0
	}
	return uint64(size)
}
