package app

import (
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/components"
	"github.com/gw123/echo-app/services"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var App *EchoApp

type EchoApp struct {
	IsHealth        bool
	areaSvc         echoapp.AreaService
	smsSvc          echoapp.SmsService
	UserSvr         echoapp.UserService
	dbPool          echoapp.DbPool
	redisPool       echoapp.RedisPool
	CompanySvr      echoapp.CompanyService
	GoodsSvr        echoapp.GoodsService
	ResourceService echoapp.ResourceService
	CommentSvr      echoapp.CommentService
	OrderSvr        echoapp.OrderService
	ActivitySvr     echoapp.ActivityService
	WsSvr           echoapp.WsService
	TestpaperSvr    echoapp.TestpaperService
	WechatService   echoapp.WechatService
}

func init() {
	App = &EchoApp{
		IsHealth: true,
	}
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
	comSvr := MustGetCompanyService()
	redis := MustGetRedis("")
	smsSvc := services.NewSmsService(comSvr, redis)
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

func MustGetRedLock(dbName string) *redislock.Client {
	redisClient := MustGetRedis(dbName)
	return redislock.New(redisClient)
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
		return nil, errors.Wrap(err, "GetUserSerevice->GetDb")
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
		panic(errors.Wrap(err, "MustuserSer-> GetUserSvr"))
	}
	return userSvr
}

func MustGetWsService() echoapp.WsService {
	App.WsSvr = services.NewWsService()
	return App.WsSvr
}

func GetGoodsService() (echoapp.GoodsService, error) {
	if App.GoodsSvr != nil {
		return App.GoodsSvr, nil
	}
	goodsDb, err := GetDb("shop")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "GetRedis")
	}

	App.GoodsSvr = services.NewGoodsService(goodsDb, redis)
	return App.GoodsSvr, nil
}

func MustGetGoodsService() echoapp.GoodsService {
	goodsSvr, err := GetGoodsService()
	if err != nil {
		panic(errors.Wrap(err, "GetUserSvr"))
	}
	return goodsSvr
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

func GetResourceService() (echoapp.ResourceService, error) {
	if App.ResourceService != nil {
		return App.ResourceService, nil
	}
	shopDb, err := GetDb("resource")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	// redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "GetRedis")
	// }
	//App.ResourceService = services.NewResourceService(shopDb, redis, echoapp.ConfigOpts.ResourceOptions)
	App.ResourceService = services.NewResourceService(shopDb)
	return App.ResourceService, nil
}

func MustGetResourceService() echoapp.ResourceService {
	resource, err := GetResourceService()
	if err != nil {
		panic(errors.Wrap(err, "GetUserSvr"))
	}
	return resource
}
func GetCommentService() (echoapp.CommentService, error) {
	if App.CompanySvr != nil {
		return App.CommentSvr, nil
	}
	commentDb, err := GetDb("goods")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	// redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "GetRedis")
	// }
	App.CommentSvr = services.NewCommentService(commentDb)
	return App.CommentSvr, nil
}

func MustGetCommentService() echoapp.CommentService {
	comment, err := GetCommentService()
	if err != nil {
		panic(errors.Wrap(err, "GetCommentSvr"))
	}
	return comment
}

func GetOrderService() (echoapp.OrderService, error) {
	if App.OrderSvr != nil {
		return App.OrderSvr, nil
	}
	goodsDb, err := GetDb("goods")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "GetRedis")
	}

	goodsSvr := MustGetGoodsService()
	actSvr := MustGetActivityService()
	App.OrderSvr = services.NewOrderService(goodsDb, redis, goodsSvr, actSvr)
	return App.OrderSvr, nil
}

func MustGetOrderService() echoapp.OrderService {
	svr, err := GetOrderService()
	if err != nil {
		panic(errors.Wrap(err, "GetUserSvr"))
	}
	return svr
}

func GetActivityService() (echoapp.ActivityService, error) {
	if App.ActivitySvr != nil {
		return App.ActivitySvr, nil
	}
	shopDb, err := GetDb("shop")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	if err != nil {
		return nil, errors.Wrap(err, "GetRedis")
	}
	lock := MustGetRedLock("")
	App.ActivitySvr = services.NewActivityService(shopDb, redis, lock)
	return App.ActivitySvr, nil
}

func MustGetActivityService() echoapp.ActivityService {
	svr, err := GetActivityService()
	if err != nil {
		panic(errors.Wrap(err, "GetUserSvr"))
	}
	return svr
}

func GetWechatService() (echoapp.WechatService, error) {
	if App.WechatService != nil {
		return App.WechatService, nil
	}
	com := MustGetCompanyService()
	App.WechatService = services.NewWechatService(com, echoapp.ConfigOpts.Wechat.AuthRedirectUrl)
	return App.WechatService, nil
}

func MustGetWechatService() echoapp.WechatService {
	svr, err := GetWechatService()
	if err != nil {
		panic(errors.Wrap(err, "GetUserSvr"))
	}
	return svr
}

func GetTestpaperService() (echoapp.TestpaperService, error) {
	if App.TestpaperSvr != nil {
		return App.TestpaperSvr, nil
	}
	shopDb, err := GetDb("shop")
	if err != nil {
		return nil, errors.Wrap(err, "GetDb")
	}
	// redis, err := components.NewRedisClient(echoapp.ConfigOpts.Redis)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "GetRedis")
	// }

	App.TestpaperSvr = services.NewTestpaperService(shopDb)
	return App.TestpaperSvr, nil
}

func MustGetTestpaperService() echoapp.TestpaperService {
	svr, err := GetTestpaperService()
	if err != nil {
		panic(errors.Wrap(err, "GetTestPapeSvr"))
	}
	return svr
}
