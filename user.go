<<<<<<< HEAD
package echoapp

import "github.com/labstack/echo"

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
	Id         int64   `json:"id"`
	ComId      int     `json:"com_id"`
	Name       string  `json:"name"`
	Nickname   string  `json:"nickname"`
	Avatar     string  `json:"avatar"`
	Sex        string  `json:"sex"`
	City       string  `json:"city"`
	Email      string  `json:"email"`
	Mobile     string  `json:"mobile"`
	Score      int     `json:"score"`
	Openid     string  `json:"xcx_openid"`
	Unionid    string  `json:"unionid"`
	IsStaff    bool    `json:"is_staff"`
	IsVip      string  `json:"is_vip"`
	VipLevel   string  `json:"vip_level"`
	JwsToken   string  `gorm:"-" json:"jws_token"`
	SessionKey string  `gorm:"-" json:"session_key"`
	Roles      []*Role `json:"roles" gorm:"many2many:model_has_roles;ForeignKey:model_id;AssociationForeignKey:role_id"`
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

func (*UserRole) Table() string {
	return "model_has_roles"
}

type UserService interface {
	AddScore(ctx echo.Context, user *User, amount int) error
	SubScore(ctx echo.Context, user *User, amount int) error
	Login(ctx echo.Context, param *LoginParam) (*User, error)
	Register(ctx echo.Context, param *RegisterParam) (*User, error)
	GetUserById(userId int64) (*User, error)
	GetUserList(comId, currentMaxId, limit int) ([]*User, error)
	GetUserByOpenId(comId int, openId string) (*User, error)
	Save(user *User) error
	GetUserByToken(token string) (*User, error)
	UpdateJwsToken(user *User) error
	UpdateCachedUser(user *User) (err error)
	GetCachedUserById(userId int64) (*User, error)
	Jscode2session(comId int, code string) (*User, error)
	AutoRegisterWxUser(user *User) (err error)
	//Jscode2session(comId int, code string) (*User, error)

}
=======
package echoapp

import "github.com/labstack/echo"

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
	Id         int64   `json:"id"`
	ComId      int     `json:"com_id"`
	Name       string  `json:"name"`
	Nickname   string  `json:"nickname"`
	Avatar     string  `json:"avatar"`
	Sex        string  `json:"sex"`
	City       string  `json:"city"`
	Email      string  `json:"email"`
	Mobile     string  `json:"mobile"`
	Score      int     `json:"score"`
	Openid     string  `grom:"xcx_openid" json:"-"`
	Unionid    string  `gorm:"unionid" json:"-"`
	IsStaff    bool    `json:"is_staff"`
	IsVip      string  `json:"is_vip"`
	VipLevel   string  `json:"vip_level"`
	JwsToken   string  `gorm:"-" json:"jws_token"`
	SessionKey string  `gorm:"session_key" json:"-"`
	Roles      []*Role `json:"roles" gorm:"many2many:model_has_roles;ForeignKey:model_id;AssociationForeignKey:role_id"`
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

type UserService interface {
	AddScore(ctx echo.Context, user *User, amount int) error
	SubScore(ctx echo.Context, user *User, amount int) error
	Login(ctx echo.Context, param *LoginParam) (*User, error)
	Register(ctx echo.Context, param RegisterParam) (*User, error)
	GetUserById(userId int64) (*User, error)
	GetUserList(comId, currentMaxId, limit int) ([]*User, error)
	GetUserByOpenId(comId int, openId string) (*User, error)
	Save(user *User) error
	GetUserByToken(token string) (*User, error)
	UpdateJwsToken(user *User) error
	UpdateCachedUser(user *User) (err error)
	GetCachedUserById(userId int64) (*User, error)
	Jscode2session(comId int, code string) (*User, error)
	AutoRegisterWxUser(user *User) (err error)
	//Jscode2session(comId int, code string) (*User, error)
}
>>>>>>> develop
