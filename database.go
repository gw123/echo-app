package echoapp

import "github.com/jinzhu/gorm"

type DbPool interface {
	Db(dbname string) (*gorm.DB, error)
}
