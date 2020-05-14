package services

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/pkg/errors"
	"sync"
)

type DbPoolService struct {
	dbMap       map[string]*gorm.DB
	dbOptionMap map[string]echoapp.DBOption
	mu          sync.Mutex
}

func NewDbPool(options map[string]echoapp.DBOption) *DbPoolService {
	return &DbPoolService{
		dbOptionMap: options,
		dbMap:       map[string]*gorm.DB{},
	}
}

func (dSvr DbPoolService) Db(dbName string) (*gorm.DB, error) {
	client, ok := dSvr.dbMap[dbName]
	if !ok || client == nil {
		dbOption, ok := dSvr.dbOptionMap[dbName]
		if !ok {
			return nil, errors.New("notfound DbName:" + dbName)
		}
		var err error
		client, err = gorm.Open(dbOption.Driver, dbOption.DSN)
		if err != nil {
			return nil, errors.Wrap(err, "gorm.open")
		}
		//防止多线程并发操作
		dSvr.mu.Lock()
		defer dSvr.mu.Unlock()
		dSvr.dbMap[dbName] = client
	}
	return client, nil
}
