package echoapp_util

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	"time"
)

//先从缓存读取数据,如果不存在调用 call方法获取后在存放到数据库中
func GetCache(client *redis.Client, key string, out interface{}, call func() (interface{}, error)) (string, error) {
	val := client.Get(key).Val()
	//是否需要重新刷新缓存
	isOK := false
	if val != "" {
		if out != nil {
			err := json.Unmarshal([]byte(val), out)
			if err != nil {
				isOK = false
			} else {
				isOK = true
			}
		} else {
			isOK = true
		}
	}

	if val == "" || !isOK {
		newVal, err := call()
		if err != nil {
			return "", err
		}
		tmp, err := json.Marshal(newVal)
		if err != nil {
			return "", err
		}
		if err := client.Set(key, tmp, time.Hour).Err(); err != nil {
			return "", err
		}
		if err := json.Unmarshal([]byte(tmp), out); err != nil {
			return "", err
		}
		return string(tmp), err
	}

	return val, nil
}
