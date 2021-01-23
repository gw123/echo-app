package app_components

import (
	"sync"

	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
)

var shopDb *gorm.DB
var shopDbOnce sync.Once

func GetShopDb() (*gorm.DB, error) {
	client, err := echoapp.GetApolloClient()

	if err != nil {
		return nil, errors.Wrap(err, "getShopDb")
	}

	shopDbOnce.Do(func() {
		dsn := client.GetStringValue("database.shop.dsn", "")
		shopDb, err = gorm.Open("mysql", dsn)
	})

	return shopDb, err
}

var userDb *gorm.DB
var userDbOnce sync.Once

func GetUserDb() (*gorm.DB, error) {
	client, err := echoapp.GetApolloClient()

	if err != nil {
		return nil, errors.Wrap(err, "getUserDb")
	}

	userDbOnce.Do(func() {
		dsn := client.GetStringValue("database.user.dsn", "")
		userDb, err = gorm.Open("mysql", dsn)
	})
	return userDb, err
}
