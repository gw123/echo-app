package controllers

import (
	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
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

func (t *UserController) Register(c echo.Context) error {
	postUser := &echoapp.RegisterParam{}
	if err := c.Bind(postUser); err != nil {
		return t.Fail(c, echoapp.Err_Argument, "输入错误", err)
	}

	user, err := t.userSvr.Register(c, postUser)
	if err != nil {
		return t.Fail(c, echoapp.Err_Argument, "RegisterUser", err)
	}
	return t.Success(c, user)
}
func (sCtl *UserController) Login(ctx echo.Context) error {
	param := &echoapp.LoginParam{}
	if err := ctx.Bind(param); err != nil {
		return sCtl.Fail(ctx, echoapp.Err_Argument, "登录失败", err)
	}
	user, err := sCtl.userSvr.Login(ctx, param)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.Err_InnerError, "登录失败", err)
	}
	return sCtl.Success(ctx, user)
}
func (sCtl *UserController) AddUserScore(ctx echo.Context) error {
	params := &echoapp.UserScoreParam{}
	if err := ctx.Bind(params); err != nil {
		return sCtl.Fail(ctx, echoapp.Err_Argument, err.Error(), err)
	}
	user, err := sCtl.userSvr.GetUserById(params.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return sCtl.Fail(ctx, echoapp.Err_NotFound, "用户不存在", err)
		} else {
			return sCtl.Fail(ctx, echoapp.Err_DBError, "系统异常", err)
		}
	}
	sCtl.userSvr.AddScore(ctx, user, params.Score)
	return sCtl.Success(ctx, nil)
}
func (uctr *UserController) SubUserScore(c echo.Context) error {
	param := &echoapp.UserScoreParam{}
	if err := c.Bind(param); err != nil {
		return uctr.Fail(c, echoapp.Err_Argument, err.Error(), err)
	}
	user, err := uctr.userSvr.GetUserById(param.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return uctr.Fail(c, echoapp.Err_NotFound, "uesr not found", err)
		} else {
			return uctr.Fail(c, echoapp.Err_DBError, "system error", err)
		}
	}
	uctr.userSvr.SubScore(c, user, param.Score)
	return uctr.Success(c, nil)
}
func (u *UserController) AddUserRoles(c echo.Context) error {

	roles := &echoapp.Role{}
	err := c.Bind(roles)
	if err != nil {
		return u.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "Bind"))
	}
	err = u.userSvr.Addroles(c, roles)
	if err != nil {
		return u.Fail(c, echoapp.Err_Argument, "角色添加失败", errors.Wrap(err, "Addroles"))
	}
	return u.Success(c, nil)
}
func (sCtl *UserController) GetUserRoles(ctx echo.Context) error {
	return nil
}
func (sCtl *UserController) CheckHasRoles(ctx echo.Context) error {
	return nil
}
func (u *UserController) AddPermissions(c echo.Context) error {
	permissions := &echoapp.Permission{}
	err := c.Bind(permissions)
	if err != nil {
		return u.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "Bind"))
	}
	err = u.userSvr.AddPermission(c, permissions)

	if err != nil {
		return u.Fail(c, echoapp.Err_Argument, "权限添加失败", errors.Wrap(err, "AddPermission"))
	}
	return u.Success(c, nil)
}
func (t *UserController) RoleHasPermission(c echo.Context) error {

	Param := &echoapp.RoleandPermissionParam{}
	if err := c.Bind(Param); err != nil {
		return t.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "输入有误"))
	}
	newRHP, err := t.userSvr.RoleHasPermission(c, Param)
	echoapp_util.ExtractEntry(c).Infof("role:%s,permission:%s", Param.Role, Param.Permission)
	if err != nil {
		return t.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "RoleHasPermission"))
	}
	return t.Success(c, newRHP)

}

func (sCtl *UserController) GetUserInfo(ctx echo.Context) error {
	echoapp_util.ExtractEntry(ctx).Info("getUserInfo")
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.Err_NotFound, "未发现用户", err)
	}
	return sCtl.Success(ctx, user)
}

func (sCtl *UserController) Logout(ctx echo.Context) error {
	return nil
}

func (sCtl *UserController) SendVerifyCodeSms(ctx echo.Context) error {
	return nil
}

func (sCtl *UserController) GetVerifyPic(ctx echo.Context) error {
	return nil
}
