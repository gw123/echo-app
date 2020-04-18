package echoapp

import (
	"github.com/gw123/echo-app/models"
	"github.com/labstack/echo"
)

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
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"mobile"`
	Pwd   string `json:"password"`
	//Email  string `json:"email"`
	//Avatar string `json:"avatar"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

func (r *RegisterUser) TableName() string {
	return "users"
}

type User struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Phone  string `json:"mobile"`
	Score  int    `json:"score"`
	Role   string `json:"role"`
}

type UserService interface {
	Save(ctx echo.Context, user *User) error
	AddScore(ctx echo.Context, user *User, amount int) error
	SubScore(ctx echo.Context, user *User, amount int) error
	Login(ctx echo.Context, param *LoginParam) (*User, error)
	GetUserById(ctx echo.Context, userId int) (*User, error)
	CreateUser(c echo.Context, user *RegisterUser) error
	Addroles(c echo.Context, param *models.Role) error
	AddPermission(c echo.Context, param *models.Permission) error
	RegisterUser(c echo.Context, param *RegisterUser) error
}
