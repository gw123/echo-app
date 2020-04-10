package echoapp

import (
	"github.com/jinzhu/gorm"
)

type DbPool interface {
	GetDbByName(dbname string) *gorm.DB
}
