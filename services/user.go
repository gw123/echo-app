package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gw123/glog"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/components"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/medivhzhan/weapp"
	"github.com/pkg/errors"
)

const (
	//redis 相关的key
	RedisUserKey            = "User:%d"
	RedisUserXCXOpenidKey   = "UserXCXOpenid:%d"
	RedisSmsLoginCodeKey    = "SmsLoginCode"
	RedisUserXCXAddrKey     = "UserXCXDefaultAddr:%d"
	RedisUserCollectTypeKey = "UserCollectType:%d,%s"
	//RedisUserCollectionKey  = "UserCollection:%d"
	RedisUserHistoryListKey = "UserHistoryList"
	RedisUserHistoryLockKey = "UserHistoryLock"
	RedisUserHistoryHotKey  = "CompanyTypeZset:%d,%s"
	//登录方式
	LoginMethodPassword = "password"
	LoginMethodSms      = "sms"

	RedisUserCodeKey = "UserCode:%d"
)

func FormatUserRedisKey(userId int64) string {
	return fmt.Sprintf(RedisUserKey, userId)
}

func FormatUserCodeRedisKey(rand int32) string {
	return fmt.Sprintf(RedisUserCodeKey, rand)
}

func FormatOpenidRedisKey(userId int64) string {
	return fmt.Sprintf(RedisUserXCXOpenidKey, userId)
}

func FormatUserAddrRedisKey(userId int64) string {
	return fmt.Sprintf(RedisUserXCXAddrKey, userId)
}

func FormatUserCollectionTypeRedisKey(userId int64, collectType string) string {
	return fmt.Sprintf(RedisUserCollectTypeKey, userId, collectType)
}

func FormatUserHistoryHotKey(comId uint, targetType string) string {
	return fmt.Sprintf(RedisUserHistoryHotKey, comId, targetType)
}

type UserService struct {
	db          *gorm.DB
	redis       *redis.Client
	jws         *components.JwsHelper
	hashIdsSalt string
}

func (u *UserService) GetUserByToken(token string) (*echoapp.User, error) {
	panic("implement me")
}

func NewUserService(db *gorm.DB, redis *redis.Client, jws *components.JwsHelper, salt string) *UserService {
	return &UserService{
		db:          db,
		redis:       redis,
		jws:         jws,
		hashIdsSalt: salt,
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

func (u *UserService) GetUserByOpenId(comId uint, openId string) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := u.db.Where("com_id = ? and openid = ? ", comId, openId).First(user).Error; err != nil {
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

func (u *UserService) CheckVerifyCode(comId uint, phone string, code string) bool {
	rCode := u.redis.Get(fmt.Sprintf(RedisSmsCode, comId, phone)).Val()
	if code == "" || len(code) < 4 || rCode != code {
		return false
	}
	return true
}

func (u *UserService) Login(ctx echo.Context, param *echoapp.LoginParam) (*echoapp.User, error) {
	echoapp_util.ExtractEntry(ctx).Warnf("login params :%+v", param)
	user := &echoapp.User{}
	if param.Method == LoginMethodSms {
		ok := u.CheckVerifyCode(param.ComId, param.Username, param.SmsCode)
		if !ok {
			return nil, errors.New("短信验证码不匹配")
		}

		var err error
		user, err = u.GetUserByMobile(param.ComId, param.Username)
		if err == gorm.ErrRecordNotFound {
			user = &echoapp.User{}
			user.Mobile = param.Username
			user.Nickname = param.Username
			user.ComId = param.ComId
			if err := u.db.Create(user).Error; err != nil {
				return nil, errors.Wrap(err, "创建用户失败")
			}
		} else if err != nil {
			return nil, errors.Wrap(err, "GetUserByMobile")
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
	data["client_id"] = ctx.Request().Header.Get("ClientID")
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "Marshal")
	}

	token, err := u.jws.CreateToken(user.Id, string(payload))

	user.JwsToken = token
	if err := u.UpdateCachedUser(user); err != nil {
		echoapp_util.ExtractEntry(ctx).Warnf("redis set user err :%s", err.Error())
	}
	return user, nil
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

func (u *UserService) GetCachedUserDefaultAddrById(userId int64) (*echoapp.Address, error) {
	addr := &echoapp.Address{}
	data, err := u.redis.Get(FormatUserAddrRedisKey(userId)).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), addr); err != nil {
		return nil, err
	}
	return addr, nil
}

func (u *UserService) UpdateCachedUser(user *echoapp.User) (err error) {
	data, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	err = u.redis.Set(FormatUserRedisKey(user.Id), data, 0).
		Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return err
}

func (u *UserService) UpdateCachedUserDefaultAddr(addr *echoapp.Address) (err error) {
	//r := time.Duration(rand.Int63n(180))
	data, err := json.Marshal(addr)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	//fmt.Println(string(data))
	err = u.redis.Set(FormatUserAddrRedisKey(addr.UserID), data, 0).
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

func (u *UserService) GetUserByMobile(comId uint, mobile string) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := u.db.Where("com_id = ? and mobile = ?", comId, mobile).First(user).Error; err != nil {
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
func (u *UserService) AutoRegisterWxUser(newUser *echoapp.User) (user *echoapp.User, err error) {
	data := make(map[string]interface{})
	data["username"] = newUser.Nickname
	data["com_id"] = newUser.ComId
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrap(err, "Marshal")
	}

	user, err = u.GetUserByOpenId(newUser.ComId, newUser.Openid)
	if err == nil {
		//用户存在更新jwstoken
		glog.Infof("用户已经存在 %+v", user)
		user.JwsToken, err = u.jws.CreateToken(user.Id, string(payload))
		if err != nil {
			return nil, errors.Wrap(err, "createToken")
		}
		if err := u.Save(user); err != nil {
			return nil, errors.Wrap(err, "save newUser")
		}
		u.UpdateCachedUser(user)
		return user, nil
	} else if err == gorm.ErrRecordNotFound {
		//用户不存在创建新的用户
		newUser.JwsToken, err = u.jws.CreateToken(newUser.Id, string(payload))
		if err != nil {
			return nil, errors.Wrap(err, "createToken")
		}
		if err := u.Save(newUser); err != nil {
			return nil, errors.Wrap(err, "save newUser")
		}

		u.UpdateCachedUser(newUser)
		return newUser, nil
	}
	//更新缓存
	return nil, errors.Wrap(err, "GetUserByOpenId")
}

//自动注册微信用户
func (u *UserService) ChangeUserJwsToken(user *echoapp.User) (err error) {
	data := make(map[string]interface{})
	data["username"] = user.Nickname
	data["com_id"] = user.ComId
	payload, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "Marshal")
	}

	user.JwsToken, err = u.jws.CreateToken(user.Id, string(payload))
	if err != nil {
		return errors.Wrap(err, "createToken")
	}

	if err := u.Save(user); err != nil {
		return errors.Wrap(err, "save user")
	}
	return nil
}

// 设置用户vip等级
func (u *UserService) SetVipLevel(user *echoapp.User, level int16) (err error) {
	user.VipLevel = level
	if err := u.Save(user); err != nil {
		return errors.Wrap(err, "set vip level")
	}
	glog.Infof("userid : %d ,set vip level %d", user.Id, level)
	return nil
}

//解析当前用户，如果用户未注册自动注册
func (u *UserService) Jscode2session(comId uint, code string) (*echoapp.User, error) {
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
	//fmt.Printf("返回结果: %#v", res)
	user, err := u.GetUserByOpenId(comId, res.OpenID)
	if err == gorm.ErrRecordNotFound {
		user = &echoapp.User{
			Nickname: "未设置用户名",
			ComId:    comId,
			Openid:   res.OpenID,
		}
		_, err := u.AutoRegisterWxUser(user)
		if err != nil {
			return nil, errors.Wrap(err, "微信用户保存失败")
		}
	} else if err != nil {
		return nil, errors.Wrapf(err, "查找失败请重试")
	}

	return user, nil
}

//解析当前用户，如果用户未注册自动注册
func (u *UserService) RegisterWechatUser(comId uint, newUser *echoapp.User) (*echoapp.User, error) {
	//fmt.Printf("返回结果: %#v", res)
	user, err := u.GetUserByOpenId(comId, newUser.Openid)
	if err == gorm.ErrRecordNotFound {
		newUser.ComId = comId
		_, err := u.AutoRegisterWxUser(newUser)
		if err != nil {
			return nil, errors.Wrap(err, "微信用户保存失败")
		}
		user = newUser
	} else if err != nil {
		return nil, errors.Wrapf(err, "查找失败请重试")
	}

	return user, nil
}

func (t *UserService) Register(ctx echo.Context, param *echoapp.RegisterParam) (*echoapp.User, error) {

	err := t.db.Table("users").Where("phone=?", param.Mobile)
	if err.Error != nil && err.RecordNotFound() {
		return nil, errors.Wrap(err.Error, "RecordAwardHistory has Found")
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

func (uSvr *UserService) AddScore(ctx echo.Context, user *echoapp.User, amount int) error {
	user.Score += amount
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,增加积分: %d", user.Id, amount)
	return uSvr.Save(user)
}

func (uSvr *UserService) SubScore(ctx echo.Context, user *echoapp.User, amount int) error {
	user.Score -= amount
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,消耗积分: %d", user.Id, amount)
	return uSvr.Save(user)
}

func (uSvr *UserService) AddScoreByUserId(comID, userID uint, score int, source string, detail string, note string) error {
	user, err := uSvr.GetUserById(int64(userID))
	if err != nil {
		return err
	}
	if user.ComId != comID {
		return errors.New("comID userID not match")
	}

	user.Score += score
	return uSvr.Save(user)
}

func (uSvr *UserService) GetUserAddressList(userId int64) ([]*echoapp.Address, error) {
	var addrList []*echoapp.Address
	if err := uSvr.db.Table("user_address").Where("user_id=?", userId).Order("updated_at DESC").Find(&addrList).Error; err != nil {
		return nil, errors.Wrap(err, "GetUserAddrList")
	}
	return addrList, nil
}

func (uSvr *UserService) CreateUserAddress(address *echoapp.Address) error {
	if address.ProvinceId == 0 {
		return errors.New("请选择省份")
	}
	if address.CityId == 0 {
		return errors.New("请选择城市")
	}
	if address.DistrictId == 0 {
		return errors.New("请选择区县")
	}
	if len(address.Mobile) != 11 {
		return errors.New("请输入11位正确手机号")
	}
	if len(address.Address) <= 5 || len(address.Address) > 100 {
		return errors.New("详细地址应为5-100个字符之间")
	}
	if address.Checked {
		addrList, err := uSvr.GetUserAddressList(address.UserID)
		if err != nil {
			return err
		}
		for _, addr := range addrList {
			if addr.Checked {
				addr.Checked = false
				uSvr.db.Save(addr)
				break
			}

		}
	}
	if err := uSvr.db.Create(address).Error; err != nil {
		return err
	}
	if address.Checked {
		if err := uSvr.UpdateCachedUserDefaultAddr(address); err != nil {
			return errors.Wrap(err, "UpdateCachedUserDefaultAddr")
		}
	}
	return nil
}
func (uSvr *UserService) GetUserAddrById(addrId int64) (*echoapp.Address, error) {
	res := &echoapp.Address{}
	if err := uSvr.db.Where("id=?", addrId).First(res).Error; err != nil {
		return nil, errors.Wrap(err, "GetUsrAddrById")
	}
	return res, nil
}

func (uSvr *UserService) UpdateUserAddress(address *echoapp.Address) error {
	if address.ProvinceId == 0 {
		return errors.New("请选择省份")
	}
	if address.CityId == 0 {
		return errors.New("请选择城市")
	}
	if address.DistrictId == 0 {
		return errors.New("请选择区县")
	}
	if len(address.Mobile) != 11 {
		return errors.New("请输入11位正确手机号")
	}
	if len(address.Address) <= 5 || len(address.Address) > 100 {
		return errors.New("详细地址应为5-100个字符之间")
	}
	addr, err := uSvr.GetUserAddrById(int64(address.ID))
	if err != nil {
		return err
	}
	addr.Address = address.Address
	addr.CityId = address.CityId
	addr.DistrictId = address.DistrictId
	addr.ProvinceId = address.ProvinceId
	addr.Username = address.Username
	addr.Code = address.Code
	addr.Mobile = address.Mobile
	if (addr.Checked == true) && (address.Checked == false) {
		if err := uSvr.redis.Del(FormatUserAddrRedisKey(address.UserID)).Err(); err != nil {
			return err
		}
	}
	addr.Checked = address.Checked
	if address.Checked && addr.Checked == false {
		addrList, err := uSvr.GetUserAddressList(address.UserID)
		if err != nil {
			return err
		}
		for _, addr := range addrList {
			if addr.Checked {
				addr.Checked = false
				uSvr.db.Save(addr)
				break
			}
		}

	}
	if err := uSvr.db.Save(addr).Error; err != nil {
		return err
	}

	if address.Checked {
		if err := uSvr.UpdateCachedUserDefaultAddr(addr); err != nil {
			return errors.Wrap(err, "UpdateCachedUserDefaultAddr")
		}
	}
	return nil
}

func (uSvr *UserService) DelUserAddress(address *echoapp.Address) error {
	if err := uSvr.db.Delete(address).Error; err != nil {
		return err
	}
	if address.Checked {
		if err := uSvr.redis.Del(FormatUserAddrRedisKey(address.UserID)).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (u *UserService) IsCollect(userId int64, targetId uint, targetType string) (bool, error) {
	ok, err := u.redis.SIsMember(FormatUserCollectionTypeRedisKey(userId, targetType), targetId).Result()
	return ok, err
}

func (uSvr *UserService) CreateUserCollection(collection *echoapp.Collection) error {
	if collection.TargetId == 0 {
		return errors.New("商品不存在")
	}
	ok, err := uSvr.IsCollect(collection.UserID, collection.TargetId, collection.Type)
	if err != nil {
		return errors.Wrap(err, "redis Sismember")
	}

	if ok {
		glog.Info("isCollect用户已经收藏")
		return nil
	}
	return uSvr.UpdateCacheUserCollection(collection)
}

func (u *UserService) GetCachedUserCollectionTypeSet(userId int64, targetType string) ([]string, error) {
	dataArr, err := u.redis.SMembers(FormatUserCollectionTypeRedisKey(userId, targetType)).Result()
	if err != nil {
		return nil, err
	}
	return dataArr, nil
}

func (uSvr *UserService) DelUserCollection(userId int64, collectType string, targetId uint) error {
	ok, err := uSvr.IsCollect(userId, targetId, collectType)
	if err != nil {
		return errors.Wrap(err, "redis Sismember")
	}
	if !ok {
		return nil
	}
	if err := uSvr.redis.SRem(
		FormatUserCollectionTypeRedisKey(userId, collectType),
		targetId).
		Err(); err != nil {
		return err
	}
	return nil
}

func (u *UserService) UpdateCacheUserCollection(collection *echoapp.Collection) (err error) {
	len, err := u.redis.SCard(FormatUserCollectionTypeRedisKey(collection.UserID, collection.Type)).Result()
	if err != nil {
		return errors.Wrap(err, "redis get Hlen")
	}
	if len > 1000 {
		return errors.New("key field beyond the limit")
	}
	err = u.redis.SAdd(FormatUserCollectionTypeRedisKey(collection.UserID, collection.Type), collection.TargetId).Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return errors.Wrap(err, "UpdateCacheUserCollection:sadd")
}

func (uSvr *UserService) CreateUserHistory(history *echoapp.History) error {
	if history.TargetId <= 0 {
		return errors.New("目标不存在")
	}
	len, err := uSvr.redis.LLen(RedisUserHistoryListKey).Result()
	//fmt.Println(len)
	if err != nil {
		if len == 0 {
			return errors.Wrap(err, "redis Llen is nil")
		}
		return errors.Wrap(err, "resdis LLen")
	}
	if len >= 2 {
		if ok := uSvr.redis.SetNX(RedisUserHistoryLockKey, 1, time.Second*5).Val(); ok {
			targetArr, err := uSvr.GetCacheUserHistoryList(uint(len))
			if err != nil {
				uSvr.redis.Del(RedisUserHistoryLockKey)
				return errors.Wrap(err, "CreateUserHistory->GetCacheUserHistoryList")
			}
			for _, val := range targetArr {
				var resHistory = &echoapp.History{}
				if err := json.Unmarshal([]byte(val), resHistory); err != nil {
					glog.DefaultLogger().WithField(RedisUserHistoryListKey, val)
					continue
				}
				if err := uSvr.db.Create(resHistory).Error; err != nil {
					glog.DefaultLogger().WithField(RedisUserHistoryListKey, val)
					continue
				}
				delHisVal, err := uSvr.redis.RPop(RedisUserHistoryListKey).Result()
				if err != nil {
					glog.DefaultLogger().WithField(RedisUserHistoryListKey, "RPop:"+delHisVal)
					continue
				}
			}
			uSvr.redis.Del(RedisUserHistoryLockKey)
		} else {
			glog.DefaultLogger().Warn("CreateUserHistory 获取锁失败")
		}
	}
	if err := uSvr.UpdateCacheUserHistoryHot(history); err != nil {
		glog.DefaultLogger().WithField(FormatUserHistoryHotKey(history.ComId, history.Type), history.TargetId)
	}
	return uSvr.UpdateCacheUserHistory(history)
}

func (u *UserService) GetCacheUserHistoryList(length uint) ([]string, error) {
	if length > 500 {
		length = 500
	}
	dataArr, err := u.redis.LRange(RedisUserHistoryListKey, 0, int64(length)).Result()

	if err != nil {
		return nil, err
	}
	return dataArr, nil
}

func (u *UserService) GetCacheUserHistoryHotZset(comId uint, targetType string) ([]string, error) {
	setLen, err := u.redis.ZCard(FormatUserHistoryHotKey(comId, targetType)).Result()
	//fmt.Println(setLen)
	if err != nil {
		return nil, err
	}
	dataArr, err := u.redis.ZRevRange(FormatUserHistoryHotKey(comId, targetType), 0, setLen-1).Result()
	return dataArr, err
}

func (u *UserService) GetUserHistoryList(userId int64, lastId uint, limit int) ([]*echoapp.History, error) {
	var historyList []*echoapp.History
	if err := u.db.Debug().
		Table("user_history").
		Where("user_id=? AND id>?", userId, lastId).
		Limit(limit).
		Order("created_at asc").
		Find(&historyList).
		Error; err != nil {
		return nil, errors.Wrap(err, "GetUserCollectList")
	}
	return historyList, nil
}

func (u *UserService) UpdateCacheUserHistory(history *echoapp.History) (err error) {
	data, err := json.Marshal(history)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	err = u.redis.LPush(RedisUserHistoryListKey, string(data)).Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return err
}

func (u *UserService) UpdateCacheUserHistoryHot(history *echoapp.History) (err error) {
	member := strconv.Itoa(int(history.TargetId))
	err = u.redis.ZIncrBy(FormatUserHistoryHotKey(history.ComId, history.Type), 1, member).Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return err
}

// GetUserCodeAndUpdate 更新并且获取UserCode 随机数
func (u *UserService) GetUserCodeAndUpdate(user *echoapp.User) (string, error) {
	randNum := rand.Int31n(887654321)
	randNum += 113456789
	for {
		// 循环生成不重复的randNum
		exist, err := u.redis.Exists(FormatUserCodeRedisKey(randNum)).Result()
		if err != nil {
			return "", errors.Wrap(err, "GetUserCodeAndUpdate")
		}
		if exist == 0 {
			break
		}
		randNum = rand.Int31n(387654321)
	}

	if err := u.redis.Set(FormatUserCodeRedisKey(randNum), user.Id, time.Second*30).Err(); err != nil {
		return "", errors.Wrap(err, "GetUserCodeAndUpdate.")
	}

	return strconv.Itoa(int(randNum)), nil
}

// GetUserByUserCode 通过userCode获取User
func (u *UserService) GetUserIdByUserCode(code string) (int64, error) {
	randNum, err := strconv.ParseInt(code, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "GetUserByUserCode ParseInt.")
	}

	val, err := u.redis.Get(FormatUserCodeRedisKey(int32(randNum))).Result()
	if err != nil {
		return 0, errors.Wrap(err, "GetUserCodeAndUpdate.")
	}

	userID, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "GetUserByUserCode ParseInt.")
	}

	return userID, nil
}
