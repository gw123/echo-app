package app_components

import (
	"sync"

	"github.com/spf13/viper"

	"github.com/jinzhu/gorm"
)

var shopDb *gorm.DB
var shopDbOnce sync.Once

func GetShopDb() (*gorm.DB, error) {
	var err error
	shopDbOnce.Do(func() {
		dsn := viper.GetString("database.shop.dsn")
		shopDb, err = gorm.Open("mysql", dsn)
	})

	return shopDb, err
}

var userDb *gorm.DB
var userDbOnce sync.Once

func GetUserDb() (*gorm.DB, error) {
	var err error

	userDbOnce.Do(func() {
		dsn := viper.GetString("database.user.dsn")
		userDb, err = gorm.Open("mysql", dsn)
	})
	return userDb, err
}
