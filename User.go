package echoapp

import "github.com/labstack/echo"

type RegisterParam struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type LoginParam struct {
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type RegisterUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Mobile   string `json:"mobile"`
	Password string `json:"password2"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func (r *RegisterUser) TableName() string {
	return "users"
}

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	Score  int    `json:"score"`
}

type UserService interface {
	Save(ctx echo.Context, user *User) error
	AddScore(ctx echo.Context, user *User, amount int) error
	SubScore(ctx echo.Context, user *User, amount int) error
	Login(ctx echo.Context, param LoginParam) (*User, error)
	GetUserById(ctx echo.Context, userId int) (*User, error)
}
