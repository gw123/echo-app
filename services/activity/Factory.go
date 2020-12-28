package activity

import (
	"sync"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
)

var once sync.Once
var driverFactory *DriverFactory

// 活动驱动工厂类
type DriverFactory struct {
	db                     *gorm.DB
	cache                  *redis.Client
	activityDao            Dao
	awardVipActivityDriver echoapp.ActivityDriver
	once                   sync.Once
}

// 单例方式使用工厂，暴露出一个获取实例的方法
func GetSingleDriverFactory(db *gorm.DB, cache *redis.Client) *DriverFactory {
	once.Do(func() {
		dao := NewActivityDao(db, cache)
		driverFactory = &DriverFactory{
			db:          db,
			cache:       cache,
			activityDao: dao,
		}
	})
	return driverFactory
}

func (f *DriverFactory) GetAwardVipActivityDriver() *AwardVipActivityDriver {
	f.once.Do(func() {
		f.awardVipActivityDriver = NewAwardVipActivityDriver(f.activityDao)
	})
	return NewAwardVipActivityDriver(f.activityDao)
}

func (f *DriverFactory) GetActivityDriver(driverType string) (echoapp.ActivityDriver, error) {
	var driver echoapp.ActivityDriver
	switch driverType {
	case echoapp.ActivityDriverTypeVIP:
		driver = f.GetAwardVipActivityDriver()
	}

	return driver, nil
}
