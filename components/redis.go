package components

import (
	"github.com/go-redis/redis/v7"
)

func NewRedisClient(opt *redis.Options) (*redis.Client, error) {
	redisClient := redis.NewClient(opt)
	if _, err := redisClient.Ping().Result(); err != nil {
		return nil, err
	}
	return redisClient, nil
}

func NewRedisClusterClient(opt *redis.ClusterOptions) (*redis.ClusterClient, error) {
	redisClient := redis.NewClusterClient(opt)
	if _, err := redisClient.Ping().Result(); err != nil {
		return nil, err
	}
	return redisClient, nil
}
