package app_components

import (
	"testing"

	"github.com/davecgh/go-spew/spew"

	echoapp "github.com/gw123/echo-app"
)

func TestGetShopDb(t *testing.T) {
	db, err := GetShopDb()
	if err != nil {
		t.Error(err)
	}
	goods := &echoapp.Goods{}
	if err := db.First(goods).Error; err != nil {
		t.Error(err)
	}
	spew.Dump(goods)
}

func TestGetUserDb(t *testing.T) {
	db, err := GetUserDb()
	if err != nil {
		t.Error(err)
	}
	user := &echoapp.User{}
	if err := db.First(user).Error; err != nil {
		t.Error(err)
	}
	spew.Dump(user)
}
