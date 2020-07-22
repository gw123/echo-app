package echoapp

import (
	"time"

	"github.com/labstack/echo"
)

type RegisterParam struct {
	ComId    int    `json:"com_id"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type LoginParam struct {
	ComId int `json:"com_id"`
	//auth type sms|password ;认证方式  短信验证码|账号密码
	Method   string `json:"method"`
	Username string `json:"username"`
	SmsCode  string `json:"sms_code"`
	Password string `json:"password"`
}

type User struct {
	Id         int64      `json:"id"`
	ComId      int        `json:"com_id"`
	Name       string     `json:"name"`
	Nickname   string     `json:"nickname"`
	Avatar     string     `json:"avatar"`
	Sex        string     `json:"sex"`
	City       string     `json:"city"`
	Email      string     `json:"email"`
	Mobile     string     `json:"mobile"`
	Score      int        `json:"score"`
	Openid     string     `gorm:"xcx_openid" json:"-"`
	Unionid    string     `gorm:"unionid" json:"-"`
	IsStaff    bool       `json:"is_staff"`
	IsVip      string     `json:"is_vip"`
	VipLevel   string     `json:"vip_level"`
	JwsToken   string     `gorm:"-" json:"jws_token"`
	SessionKey string     `gorm:"session_key" json:"-"`
	Roles      []*Role    `json:"roles" gorm:"many2many:model_has_roles;ForeignKey:model_id;AssociationForeignKey:role_id"`
	Address    []*Address `json:"address" gorm :"ForeignKey:UserID" `
}
type Address struct {
	//gorm.Model
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	//DeletedAt  *time.Time `sql:"index"`
	UserID     int64  `json:"user_id"`
	Username   string `json:"username" gorm:"username"` //收件人
	Mobile     string `json:"mobile"`
	Address    string `json:"address" gorm:"address"`
	Checked    bool   `json:"checked" gorm:"checked"`
	CityId     int64  `json:"city_id"`
	DistrictId int64  `json:"district_id"`
	ProvinceId int64  `json:"province_id"`
	Code       string `json:"code"`
}

func (*Address) TableName() string {
	return "user_address"
}

type Collection struct {
	//gorm.Model
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	TargetId  uint       `json:"target_id"`
	Type      string     `json:"type"`
	UserID    int64      `json:"user_id"`
}

func (*Collection) TableName() string {
	return "user_collection"
}

type Role struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type UserRole struct {
	RoleId int `json:"role_id"`
	//ModelType string `json:"model_type"`
	ModelId uint64 `json:"model_id"`
}

func (*UserRole) TableName() string {
	return "model_has_roles"
}

type History struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Type      string     `json:"type"`
	TargetId  uint       `json:"target_id"`
	UserID    int64      `json:"user_id"`
	ComId     uint       `json:"com_id"`
}

func (*History) TableName() string {
	return "user_history"
}

type UserService interface {
	AddScore(ctx echo.Context, user *User, amount int) error
	SubScore(ctx echo.Context, user *User, amount int) error
	Login(ctx echo.Context, param *LoginParam) (*User, error)
	Register(ctx echo.Context, param *RegisterParam) (*User, error)
	GetUserById(userId int64) (*User, error)
	GetCachedUserById(userId int64) (*User, error)
	GetUserList(comId, currentMaxId, limit int) ([]*User, error)
	GetUserByOpenId(comId uint, openId string) (*User, error)
	Save(user *User) error
	GetUserByToken(token string) (*User, error)
	UpdateJwsToken(user *User) error
	UpdateCachedUser(user *User) (err error)
	Jscode2session(comId uint, code string) (*User, error)
	AutoRegisterWxUser(user *User) (err error)
	//Jscode2session(comId int, code string) (*User, error)
	GetUserAddressList(userId int64) ([]*Address, error)
	CreateUserAddress(address *Address) error
	UpdateUserAddress(address *Address) error
	DelUserAddress(address *Address) error
	GetUserAddrById(addrId int64) (*Address, error)
	GetCachedUserDefaultAddrById(userId int64) (*Address, error)

	//GetCachedUserCollectionListById(userId int64) ([]*Collection, error)
	//GetUserCollectionList(userId int64, lastId uint, limit int) ([]*Collection, error)
	CreateUserCollection(collection *Collection) error
	DelUserCollection(userId int64, collectType string, targetId uint) error
	//GetUserCollectionById(userId int64, targetType string, targetId uint) (*Collection, error)
	IsCollect(userId int64, targetId uint, targetType string) (bool, error)
	GetCachedUserCollectionTypeSet(userId int64, targetType string) ([]string, error)
	//History
	UpdateCacheUserHistory(history *History) (err error)
	GetUserHistoryList(userId int64, lastId uint, limit int) ([]*History, error)
	GetCacheUserHistoryList(len uint) ([]string, error)
	CreateUserHistory(history *History) error
	GetCacheUserHistoryHotZset(comId uint, targetYype string) ([]string, error)
}
