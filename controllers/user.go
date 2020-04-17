package controllers

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/app"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserController struct {
	userSvr echoapp.UserService
	echoapp.BaseController
}

func NewUserController() *UserController {
	return &UserController{
		userSvr: app.MustUserService(),
	}
}

func (sCtl *UserController) AddUserScore(ctx echo.Context) error {
	type Params struct {
		UserId int `json:"user_id"`
		Score  int `json:"score"`
	}
	params := &Params{}
	if err := ctx.Bind(params); err != nil {
		return sCtl.Fail(ctx, echoapp.Error_ArgumentError, err.Error(), err)
	}
	user, err := sCtl.userSvr.GetUserById(ctx, params.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return sCtl.Fail(ctx, echoapp.Error_NotFound, "用户不存在", err)
		} else {
			return sCtl.Fail(ctx, echoapp.Error_DBError, "系统异常", err)
		}
	}
	sCtl.userSvr.AddScore(ctx, user, params.Score)
	return sCtl.Success(ctx, nil)
}
