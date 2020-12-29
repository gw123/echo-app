package components

import (
	"github.com/go-redis/redis/v7"
	"github.com/pkg/errors"
	"sync"
)

type RedisPoolService struct {
	redisMap       map[string]*redis.Client
	redisOptionMap map[string]*redis.Options
	mu             sync.Mutex
}

func NewRedisPool(options map[string]*redis.Options) *RedisPoolService {
	return &RedisPoolService{
		redisOptionMap: options,
		redisMap:       map[string]*redis.Client{},
	}
}

func (dSvr RedisPoolService) Redis(redisName string) (*redis.Client, error) {
	client, ok := dSvr.redisMap[redisName]
	if !ok || client == nil {
		redisOption, ok := dSvr.redisOptionMap[redisName]
		if !ok {
			return nil, errors.New("notfound RedisName:" + redisName)
		}
		client = redis.NewClient(redisOption)
		if err := client.Ping().Err(); err != nil {
			return nil, errors.Wrap(err, "gorm.open")
		}
		//防止多线程并发操作
		dSvr.mu.Lock()
		defer dSvr.mu.Unlock()
		dSvr.redisMap[redisName] = client
	}
	return client, nil
}

