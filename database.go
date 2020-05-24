package echoapp

import (
	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
)

type DbPool interface {
	Db(dbname string) (*gorm.DB, error)
}

type RedisPool interface {
	Redis(dbname string) (*redis.Client, error)
}
