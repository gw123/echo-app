package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/components"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
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
func (t *UserService) Register(ctx echo.Context, param *echoapp.RegisterParam) (*echoapp.User, error) {

	err := t.db.Table("users").Where("phone=?", param.Mobile)
	if err.Error != nil && err.RecordNotFound() {
		return nil, errors.Wrap(err.Error, "Record has Found")
	}
	echoapp_util.ExtractEntry(ctx).Infof("mobile:%s,pwd:%s", param.Mobile, param.Password)
	return nil, t.Create(param)
}

func (uSvr UserService) Create(user *echoapp.RegisterParam) error {
	if err := uSvr.db.Create(user).Error; err != nil && uSvr.db.NewRecord(user) {
		return errors.Wrap(err, "user create fail")
	}
	return nil
}

func (u *UserService) Save(user *echoapp.User) error {
	return u.db.Save(user).Error
}

func (u *UserService) GetUserByToken(token string) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := u.db.Where("api_token = ?", token).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) GetUserById(userId int64) (*echoapp.User, error) {
	//if err := u.redis.Get(fmt.Sprintf("%s%d", RedisUserKey, userId)).Result(); err != nil {
	//	return nil, err
	//}
	user := &echoapp.User{}
	if err := u.db.Where("id = ?", userId).First(user).Error; err != nil {
		return nil, err
	}
	//echoapp_util.ExtractEntry(ctx).Infof("userid:%d", userId)
	return user, nil
}
func (uSvr UserService) AddScore(ctx echo.Context, user *echoapp.User, amount int) error {
	user.Score += amount
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,增加积分: %d", user.Id, amount)
	return uSvr.Save(user)
}
func (uSvr UserService) SubScore(ctx echo.Context, user *echoapp.User, amount int) error {
	user.Score -= amount
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,消耗积分: %d", user.Id, amount)
	return uSvr.Save(user)
}
func (t *UserService) Addroles(c echo.Context, param *echoapp.Role) error {

	res := t.db.Table("roles").Where("name=?", param.Name)
	if res.Error != nil && res.RecordNotFound() {
		return errors.Wrap(res.Error, "Record has Found")
	}
	err := t.db.Create(param)
	if err.Error != nil && t.db.NewRecord(param) {
		return errors.Wrap(err.Error, "role create failed")
	} else if t.db.NewRecord(param) {
		return errors.New("not NewRecord")
	}
	echoapp_util.ExtractEntry(c).Infof("create role name:%s", param.Name)
	return nil
}

func (t *UserService) AddPermission(c echo.Context, param *echoapp.Permission) error {

	res := t.db.Table("permissions").Where("name=?", param.Name)
	if res.Error != nil && res.RecordNotFound() {
		return errors.Wrap(res.Error, "Record has Found")
	}
	err := t.db.Create(param)
	if err.Error != nil && t.db.NewRecord(param) {
		return errors.Wrap(err.Error, "permiseeion create failed")
	} else if t.db.NewRecord(param) {
		return errors.New("not NewRecord")
	}
	echoapp_util.ExtractEntry(c).Infof("create permission name:%s", param.Name)
	return nil
}

func (t *UserService) RoleHasPermission(c echo.Context, param *echoapp.RoleandPermissionParam) (*echoapp.RoleHasPermission, error) {
	role := &echoapp.Role{}
	permission := &echoapp.Permission{}
	res := t.db.Where("name=?", param.Role).Find(role)
	if res.Error != nil && res.RecordNotFound() {
		return nil, errors.Wrap(res.Error, "Role record has Found")
	} else if res.RecordNotFound() {
		return nil, errors.New("Role Record has Found")
	}
	res = t.db.Where("name=?", param.Permission)
	if res.Error != nil && res.RecordNotFound() {
		return nil, errors.Wrap(res.Error, "Permission record has Found")
	} else if res.RecordNotFound() {
		return nil, errors.New("Permission Record has Found")
	}
	rolehaspermission := &echoapp.RoleHasPermission{
		RoleId:       role.Id,
		PermissionId: permission.Id,
	}
	err := t.db.Create(rolehaspermission)
	if err.Error != nil && t.db.NewRecord(param) {
		return nil, errors.Wrap(err.Error, "rolehasper create failed")
	} else if t.db.NewRecord(param) {
		return nil, errors.New("not NewRecord")
	}
	return rolehaspermission, nil
}
