package app_components

import (
	"testing"
	"time"
)

func TestGetRedis(t *testing.T) {
	redis, err := GetRedis()
	if err != nil {
		t.Error(err)
	}
	redis.Set("test", "123", time.Minute)

	t.Log(redis.Get("test").Val())
}
