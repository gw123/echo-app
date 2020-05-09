package echoapp

import (
	"github.com/go-redis/redis"
)

type CacheOptions struct {
	Addr     string `yaml:"addr" mapstructure:"addr"`
	Password string `yaml:"password" mapstructure:"password"`
	PoolSize int    `yaml:"pool_size" mapstructure:"pool_size"`
}


func InitRedis(co *CacheOptions) error {
	redisOptions := &redis.Options{
		Addr:     co.Addr,
		Password: co.Password,
		PoolSize: co.PoolSize,
	}

	RedisClient := redis.NewClient(redisOptions)
	if _, err := RedisClient.Ping().Result(); err != nil {
		return err
	}
	return nil
}



