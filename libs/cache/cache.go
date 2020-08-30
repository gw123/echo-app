package cache

import (
	"github.com/go-redis/redis/v7"
	"time"
)

type SimpleCache struct {
	client *redis.Client
}

func NewSmipleCache(client *redis.Client) *SimpleCache {
	return &SimpleCache{client: client}
}

func (c *SimpleCache) Get(key string) interface{} {
	return c.client.Get(key).Val()
}

func (c *SimpleCache) Set(key string, val interface{}, timeout time.Duration) error {
	return c.client.Set(key, val, 0).Err()
}

func (c *SimpleCache) IsExist(key string) bool {
	return c.client.Exists(key).Val() == 1
}

func (c *SimpleCache) Delete(key string) error {
	return c.client.Del(key).Err()
}
