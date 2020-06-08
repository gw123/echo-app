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
	"github.com/medivhzhan/weapp"
	"github.com/pkg/errors"
	"time"
)

const (
	//redis 相关的key
	RedisUserKey          = "User:%d"
	RedisUserXCXOpenidKey = "UserXCXOpenid:%d"
	RedisSmsLoginCodeKey  = "SmsLoginCode"

	//登录方式
	LoginMethodPassword = "password"
	LoginMethodSms      = "sms"
)

func FormatUserRedisKey(userId int64) string {
	return fmt.Sprintf(RedisUserKey, userId)
}

func FormatOpenidRedisKey(userId int64) string {
	return fmt.Sprintf(RedisUserXCXOpenidKey, userId)
}

type UserService struct {
	db    *gorm.DB
	redis *redis.Client
	jws   *components.JwsHelper
}

func (u *UserService) GetUserByToken(token string) (*echoapp.User, error) {
	panic("implement me")
}

func NewUserService(db *gorm.DB, redis *redis.Client, jws *components.JwsHelper) *UserService {
	return &UserService{
		db:    db,
		redis: redis,
		jws:   jws,
	}
}

func (u *UserService) UpdateJwsToken(user *echoapp.User) (err error) {
	t := u.redis.TTL(FormatUserRedisKey(user.Id)).Val()
	if t < time.Hour {
		user.JwsToken, err = u.jws.CreateToken(user.Id, "")
		if err != nil {
			return errors.Wrap(err, "createToken")
		}
	}
	if err := u.UpdateCachedUser(user); err != nil {
		return errors.Wrap(err, "updateCachedUser")
	}
	return nil
}

func (u *UserService) GetUserByOpenId(comId int, openId string) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := u.db.Where("com_id = ? and openid = ? ", comId, openId).First(user).Error; err != nil {
		return nil, errors.Wrap(err, "db error")
	}
	if user.IsStaff {
		u.db.
			Table("model_has_roles as r").
			Select("roles.id ,roles.name,roles.label").
			Joins("inner join roles on roles.id = r.role_id").
			Where("r.model_id = ?", user.Id).
			Find(&user.Roles)
	}
	return user, nil
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
		if err := u.db.Debug().
			Select("id").
			Where("com_id = ? and mobile =? and password = ?", param.ComId, param.Username, sign).
			First(user).Error; err != nil {
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
	newUser, err := u.GetUserById(user.Id)
	if err != nil {
		return nil, errors.Wrap(err, "CreateToken")
	}
	newUser.JwsToken = token
	if err := u.UpdateCachedUser(newUser); err != nil {
		echoapp_util.ExtractEntry(ctx).Warnf("redis set user err :%s", err.Error())
	}
	return newUser, nil
}

func (u *UserService) GetUserList(comId, currentMaxId, limit int) ([]*echoapp.User, error) {
	userList := []*echoapp.User{}
	if err := u.db.Where("com_id = ? and id > ?", comId, currentMaxId).
		Order("id asc").Limit(limit).
		Find(&userList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return userList, nil
}

func (u *UserService) Register(ctx echo.Context, param echoapp.RegisterParam) (*echoapp.User, error) {
	panic("implement me")
}

func (u *UserService) GetCachedUserById(userId int64) (*echoapp.User, error) {
	user := &echoapp.User{}
	data, err := u.redis.Get(FormatUserRedisKey(userId)).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) UpdateCachedUser(user *echoapp.User) (err error) {
	//r := time.Duration(rand.Int63n(180))
	data, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	//fmt.Println(string(data))
	err = u.redis.Set(FormatUserRedisKey(user.Id), data, 0).
		Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return err
}

func (u *UserService) GetUserById(userId int64) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := u.db.Where(" id = ?", userId).First(user).Error; err != nil {
		return nil, err
	}
	if user.IsStaff {
		u.db.
			Table("model_has_roles as r").
			Select("roles.id ,roles.name,roles.label").
			Joins("inner join roles on roles.id = r.role_id").
			Where("r.model_id = ?", user.Id).
			Find(&user.Roles)
	}
	return user, nil
}

//基础方法自动更新cache
func (u *UserService) Save(user *echoapp.User) error {
	if err := u.db.Save(user).Error; err != nil {
		return errors.Wrap(err, "db save user")
	}
	return u.UpdateCachedUser(user)
}

//自动注册微信用户
func (u *UserService) AutoRegisterWxUser(user *echoapp.User) (err error) {
	user.JwsToken, err = u.jws.CreateToken(user.Id, "")
	if err != nil {
		return errors.Wrap(err, "createToken")
	}

	if err := u.Save(user); err != nil {
		return errors.Wrap(err, "save user")
	}
	//更新缓存
	return nil
}

//解析当前用户，如果用户未注册自动注册
func (u *UserService) Jscode2session(comId int, code string) (*echoapp.User, error) {
	company := &echoapp.Company{}
	_, err := echoapp_util.GetCache(
		u.redis,
		FormatCompanyRedisKey(comId),
		company,
		func() (interface{}, error) {
			return nil, echoapp.ErrNotFoundCache
		})
	if err != nil {
		return nil, errors.Wrap(err, "未找到 "+FormatCompanyRedisKey(comId))
	}
	res, err := weapp.Login(company.WxMiniAppId, company.WxMinSecret, code)
	if err != nil {
		return nil, errors.Wrap(err, "微信登录失败")
	}
	fmt.Printf("返回结果: %#v", res)
	user, err := u.GetUserByOpenId(comId, res.OpenID)
	if err == gorm.ErrRecordNotFound {
		user = &echoapp.User{
			Nickname: "未设置用户名",
			ComId:    comId,
			Openid:   res.OpenID,
		}
		err := u.AutoRegisterWxUser(user)
		if err != nil {
			return nil, errors.Wrap(err, "微信用户保存失败")
		}
	} else if err != nil {
		return nil, errors.Wrapf(err, "查找失败请重试")
	}

	return user, nil
}
