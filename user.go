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
	Id       int64  `json:"id"`
	ComId    int    `json:"com_id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Score    int    `json:"score"`
	Token    string `json:"-"`
	JwsToken string `gorm:"-" json:"jws_token"`
}

type Company struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type Role struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Label string `json:"label"`
}

type UserRole struct {
	RoleId    int    `json:"role_id"`
	ModelType string `json:"model_type"`
	ModelId   uint64 `json:"model_id"`
}

type UserService interface {
	AddScore(ctx echo.Context, user *User, amount int) error
	SubScore(ctx echo.Context, user *User, amount int) error
	Login(ctx echo.Context, param *LoginParam) (*User, error)
	Register(ctx echo.Context, param RegisterParam) (*User, error)
	GetUserById(userId int64) (*User, error)
	Save(user *User) error
	GetUserByToken(token string) (*User, error)
}
