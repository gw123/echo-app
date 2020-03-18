package echoapp

import (
	"github.com/jinzhu/gorm"
)

type DatabaseOptions struct {
	Options []*SingleOption `yaml:"options"`
}

type SingleOption struct {
	Name   string `yaml:"name"`
	Driver string `yaml:"driver"`
	Dsn    string `yaml:"dsn"`
}

type DB struct {
	*gorm.DB
}

type DatabaseService interface {
	GetDbByName(dbname string) *DB
}
