package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/components"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"time"
)

const (
	RedisUserKey         = "User:"
	RedisSmsLoginCodeKey = "SmsLoginCode"
	//登录方式
	LoginMethodPassword = "password"
	LoginMethodSms      = "sms"
)

//下面方式可以用到 context.context 中防止字符在key冲突的问题,
//因为echo的context只能传入字符串，所以这里改成一个个特殊字符串
//type loggerKey struct{}
//type userKey struct{}
//type userIdKey struct{}
//type comKey struct{}
//type userRolesKey struct{}
//
//var ctxLoggerKey = &loggerKey{}
//var ctxUserKey = &userKey{}
//var ctxComKey = &comKey{}
//var ctxUserIddKey = &userIdKey{}
//var ctxUserRolesKey = &userRolesKey{}

type UserService struct {
	db    *gorm.DB
	redis *redis.Client
	jws   *components.JwsHelper
}

func NewUserService(db *gorm.DB, redis *redis.Client, jws *components.JwsHelper) *UserService {
	return &UserService{
		db:    db,
		redis: redis,
		jws:   jws,
	}
}
func (u *UserService) AddScore(ctx echo.Context, user *echoapp.User, amount int) error {
	panic("implement me")
}

func (u *UserService) SubScore(ctx echo.Context, user *echoapp.User, amount int) error {
	panic("implement me")
}

func (u *UserService) Login(ctx echo.Context, param *echoapp.LoginParam) (*echoapp.User, error) {
	echoapp_util.ExtractEntry(ctx).Warnf("login params :%+v", param)
	user := &echoapp.User{}
	if param.Method == LoginMethodSms {
		code := u.redis.HGet(RedisSmsLoginCodeKey, param.Username).Val()
		if code == "" {
			return nil, errors.New("请求过期")
		}
		if code != param.SmsCode {
			return nil, errors.New("code not match")
		}
	} else {
		sign := echoapp_util.Md5(echoapp_util.Md5(param.Password))
		if err := u.db.Debug().Where("com_id = ? and mobile =? and password = ?", param.ComId, param.Username, sign).First(user).Error; err != nil {
			return nil, errors.Wrap(err, "db query")
		}
	}

	data := make(map[string]interface{})
	data["username"] = param.Username
	data["com_id"] = param.ComId
	data["avatar"] = user.Avatar
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "Marshal")
	}

	token, err := u.jws.CreateToken(user.Id, string(payload))
	if err != nil {
		return nil, errors.Wrap(err, "CreateToken")
	}
	user.JwsToken = token
	if err := u.redis.Set(fmt.Sprintf("%s%d", RedisUserKey, user.Id), user, time.Hour*24*15).Err(); err != nil {
		echoapp_util.ExtractEntry(ctx).Warnf("redis set user err :%s", err.Error())
	}
	return user, nil
}

func (u *UserService) Register(ctx echo.Context, param echoapp.RegisterParam) (*echoapp.User, error) {
	panic("implement me")
}

func (u *UserService) GetCachedUserById(userId int64) (*echoapp.User, error) {
	user := &echoapp.User{}
	data, err := u.redis.Get(fmt.Sprintf("%s%d", RedisUserKey, userId)).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetUserById(userId int64) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := u.db.Where("id = ?", userId).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) Save(user *echoapp.User) error {
	panic("implement me")
}

func (u *UserService) GetUserByToken(token string) (*echoapp.User, error) {
	panic("implement me")
}
