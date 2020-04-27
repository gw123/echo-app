package controllers

import (
	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"
<<<<<<< HEAD
	"github.com/gw123/echo-app/app"
	echoapp_util "github.com/gw123/echo-app/util"
=======
>>>>>>> refs/remotes/origin/master
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserController struct {
	userSvr echoapp.UserService
	echoapp.BaseController
}

func NewUserController(usrSvr echoapp.UserService) *UserController {
	return &UserController{
		userSvr: usrSvr,
	}
}

func (sCtl *UserController) AddUserScore(ctx echo.Context) error {
	params := &echoapp.UserScoreParam{}
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
func (uctr *UserController) SubUserScore(c echo.Context) error {
	param := &echoapp.UserScoreParam{}
	if err := c.Bind(param); err != nil {
		return uctr.Fail(c, echoapp.Error_ArgumentError, err.Error(), err)
	}
	user, err := uctr.userSvr.GetUserById(c, param.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uctr.Fail(c, echoapp.Error_ArgumentError, "uesr not found", err)
		} else {
			return uctr.Fail(c, echoapp.Error_ArgumentError, "system error", err)
		}
	}
	uctr.userSvr.SubScore(c, user, param.Score)
	return uctr.Success(c, nil)
}
func (t *UserController) Adduser(c echo.Context) error {
	postUser := &echoapp.RegisterUser{}
	if err := c.Bind(postUser); err != nil {
		return t.Fail(c, echoapp.Error_ArgumentError, "输入错误", err)
	}

	err := t.userSvr.RegisterUser(c, postUser)
	if err != nil {
		return t.Fail(c, echoapp.Error_ArgumentError, "RegisterUser", err)
	} else {
		return t.Success(c, nil)
	}
}

func (t *UserController) Login(c echo.Context) error {
	queryparam := &echoapp.LoginParam{}
	if err := c.Bind(queryparam); err != nil {
		return t.Fail(c, echoapp.Error_ArgumentError, err.Error(), errors.Wrap(err, "输入有误"))
	}
	res, err := t.userSvr.Login(c, queryparam)
	if err != nil {
		return t.Fail(c, echoapp.Error_NotFound, "Controller Login", err)
	}
	return t.Success(c, res)

}

func (u *UserController) Addroles(c echo.Context) error {

	roles := &echoapp.Role{}
	err := c.Bind(roles)
	if err != nil {
		return u.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "Bind"))
	}
	err = u.userSvr.Addroles(c, roles)
	if err != nil {
		return u.Fail(c, echoapp.Error_ArgumentError, "角色添加失败", errors.Wrap(err, "Addroles"))
	}
	return u.Success(c, nil)
}

func (u *UserController) Addpermissions(c echo.Context) error {
	permissions := &echoapp.Permission{}
	err := c.Bind(permissions)
	if err != nil {
		return u.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "Bind"))
	}
	err = u.userSvr.AddPermission(c, permissions)

	if err != nil {
		return u.Fail(c, echoapp.Error_ArgumentError, "权限添加失败", errors.Wrap(err, "AddPermission"))
	}
	return u.Success(c, nil)
}

func (t *UserController) RoleHasPermission(c echo.Context) error {

	Param := &echoapp.RoleandPermissionParam{}
	if err := c.Bind(Param); err != nil {
		return t.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "输入有误"))
	}
	newRHP, err := t.userSvr.RoleHasPermission(c, Param)
	echoapp_util.ExtractEntry(c).Infof("role:%s,permission:%s", Param.Role, Param.Permission)
	if err != nil {
		return t.Fail(c, echoapp.Error_ArgumentError, "", errors.Wrap(err, "RoleHasPermission"))
	}
	return t.Success(c, newRHP)
}
