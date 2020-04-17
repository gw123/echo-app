package services

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"sync"
)

type UserService struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (uSvr UserService) AddScore(ctx echo.Context, user *echoapp.User, amount int) error {
	user.Score += amount
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,增加积分: %d", user.Id, amount)
	return uSvr.Save(ctx, user)
}

func (uSvr UserService) SubScore(ctx echo.Context, user *echoapp.User, amount int) error {
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,消耗积分: %d", user.Id, amount)
	return uSvr.Save(ctx, user)
}

func (uSvr UserService) Login(ctx echo.Context, param echoapp.LoginParam) (*echoapp.User, error) {
	panic("implement me")
}

func (uSvr UserService) GetUserById(ctx echo.Context, userId int) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := uSvr.db.Where("id = ?", userId).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (uSvr UserService) Save(ctx echo.Context, user *echoapp.User) error {
	return uSvr.db.Save(user).Error
}
