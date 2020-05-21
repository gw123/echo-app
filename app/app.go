package app

import (
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/components"
	"github.com/gw123/echo-app/services"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var App *EchoApp

//私有变量 防止未初始化调用
type EchoApp struct {
	areaSvc    echoapp.AreaService
	smsSvc     echoapp.SmsService
	UserSvr    echoapp.UserService
	dbPool     echoapp.DbPool
	redisPool  echoapp.RedisPool
	CompanySvr echoapp.CompanyService
}

func init() {
	App = &EchoApp{}
}

func GetAreaService() (echoapp.AreaService, error) {
	if App.areaSvc != nil {
		return App.areaSvc, nil
	}
	areaSvc, err := services.NewAreaService(echoapp.ConfigOpts.Asset.AreaRoot)
	if err != nil {
		return nil, errors.Wrap(err, "GetAreaService")
	}
	App.areaSvc = areaSvc
	return areaSvc, nil
}

func MustGetAreaService() echoapp.AreaService {
	areaSvc, err := GetAreaService()
	if err != nil {
		panic(err)
	}
	return areaSvc
}

func GetSmsService() (echoapp.SmsService, error) {
	if App.smsSvc != nil {
		return App.smsSvc, nil
	}
	smsSvc := services.NewSmsService(echoapp.ConfigOpts.SmsOptionTokenMap)
	App.smsSvc = smsSvc
	return smsSvc, nil
}
func MustGetSmsService() echoapp.SmsService {
	areaSvc, err := GetSmsService()
	if err != nil {
		panic(err)
	}
	return areaSvc
}

func GetDb(dbname string) (*gorm.DB, error) {
	if dbname == "" {
		dbname = "default"
	}
	if App.dbPool != nil {
		return App.dbPool.Db(dbname)
	}
	dbPool := services.NewDbPool(echoapp.ConfigOpts.DBMap)
	App.dbPool = dbPool
	return dbPool.Db(dbname)
}

func MustGetDb(dbName string) *gorm.DB {
	db, err := GetDb(dbName)
	if err != nil {
		panic(errors.Wrap(err, "MustGetDb->GetDb"))
	}
	return db
}

func GetRedis(dbname string) (*redis.Client, error) {
	if dbname == "" {
		dbname = "default"
	}
	if App.redisPool != nil {
		return App.redisPool.Redis(dbname)
	}
	redisPool := components.NewRedisPool(echoapp.ConfigOpts.RedisMap)
	App.redisPool = redisPool
	return redisPool.Redis(dbname)
}

func MustGetRedis(dbName string) *redis.Client {
	db, err := GetRedis(dbName)
	if err != nil {
		panic(errors.Wrap(err, "MustGetRedis->GetRedis"))
	}
	return db
}

func GetJwsHelper() (*components.JwsHelper, error) {
	jws, err := components.NewJwsHelper(echoapp.ConfigOpts.Jws)
	if err != nil {
		return nil, errors.Wrap(err, "NewJwsHelper")
	}
	return jws, nil
}

func MustGetJwsHelper() *components.JwsHelper {
	userSvr, err := GetJwsHelper()
	if err != nil {
		panic(errors.Wrap(err, "GetJwsHelper"))
	}
	return userSvr
}

func GetUserService() (echoapp.UserService, error) {
	if App.UserSvr != nil {
		return App.UserSvr, nil
	}
	userDb, err := GetDb("user")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "GetRedis")
	}
	jws, err := GetJwsHelper()
	if err != nil {
		return nil, errors.Wrap(err, "GetJws")
	}
	App.UserSvr = services.NewUserService(userDb, redis, jws)
	return App.UserSvr, nil
}

func MustGetUserService() echoapp.UserService {
	userSvr, err := GetUserService()
	if err != nil {
		panic(errors.Wrap(err, "GetUserSvr"))
	}
	return userSvr
}

func GetCompanyService() (echoapp.CompanyService, error) {
	if App.CompanySvr != nil {
		return App.CompanySvr, nil
	}
	shopDb, err := GetDb("shop")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "GetRedis")
	}
	App.CompanySvr = services.NewCompanyService(shopDb, redis)
	return App.CompanySvr, nil
}

func MustGetCompanyService() echoapp.CompanyService {
	company, err := GetCompanyService()
	if err != nil {
		panic(errors.Wrap(err, "GetUserSvr"))
	}
	return company
}
