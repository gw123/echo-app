package app_components

import (
	"sync"

	"github.com/davecgh/go-spew/spew"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
)

var redisClient *redis.Client
var redisClientOnce sync.Once

func GetRedis() (*redis.Client, error) {
	cfg, err := echoapp.GetApolloClient()
	if err != nil {
		return nil, err
	}

	redisClientOnce.Do(func() {
		redisOption := &redis.Options{
			Addr:     cfg.GetStringValue("redis.addr", ""),
			Password: cfg.GetStringValue("redis.password", ""),
			PoolSize: 15,
			DB:       0,
		}
		spew.Dump(redisOption)
		redisClient = redis.NewClient(redisOption)
		err = redisClient.Ping().Err()
	})

	return redisClient, err
}
