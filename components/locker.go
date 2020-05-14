package echoapp_util

import (
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
	"time"
)

//获取锁 包装了重试的逻辑
func GetLock(client redis.Client, key string, timeout time.Duration) (*redislock.Lock, error) {
	max := timeout / (200 * time.Millisecond)
	locker := redislock.New(client)
	lock, err := locker.Obtain(key, timeout, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(time.Millisecond*200), int(max)),
		Metadata:      "my data",
		Context:       nil,
	})
	return lock, err
}

//释放锁
func RelaseLock(lock *redislock.Lock) {
	lock.Release()
}

//获取锁后执行回调, 把需要锁住的数据做更新
func GetLockThenCall(client redis.Client, key string, timeout time.Duration, fu func() error) error {
	max := timeout / (200 * time.Millisecond)
	locker := redislock.New(client)
	lock, err := locker.Obtain(key, timeout, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(time.Millisecond*100), int(max)),
		Metadata:      "my data",
		Context:       nil,
	})
	if err != nil {
		return err
	}
	defer lock.Release()
	return fu()
}
