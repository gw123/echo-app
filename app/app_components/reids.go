package app_components

import (
	"sync"

	"github.com/spf13/viper"

	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client
var redisClientOnce sync.Once

func GetRedis() (*redis.Client, error) {
	var err error

	redisClientOnce.Do(func() {
		redisOption := &redis.Options{
			Addr:     viper.GetString("redis.addr"),
			Password: viper.GetString("redis.password"),
			PoolSize: 15,
			DB:       0,
		}
		redisClient = redis.NewClient(redisOption)
		err = redisClient.Ping().Err()
	})

	return redisClient, err
}
