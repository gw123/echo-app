package app_components

import (
	"sync"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/services"
)

var smsSvr echoapp.SmsService
var smsSvrOnce sync.Once

func GetSmsService() (echoapp.SmsService, error) {
	var err2 error
	smsSvrOnce.Do(func() {
		db, err := GetShopDb()
		if err != nil {
			err2 = err
			return
		}
		redis, err := GetRedis()
		if err != nil {
			err2 = err
			return
		}
		smsSvr = services.NewSmsService(db, redis)
	})

	return smsSvr, err2
}
