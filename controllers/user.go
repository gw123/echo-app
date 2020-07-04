package controllers

import (
	"strconv"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserController struct {
	userSvr echoapp.UserService
	smsSvr  echoapp.SmsService
	echoapp.BaseController
}

func NewUserController(usrSvr echoapp.UserService) *UserController {
	return &UserController{
		userSvr: usrSvr,
	}
}

func (sCtl *UserController) AddUserScore(ctx echo.Context) error {
	type Params struct {
		UserId int64 `json:"user_id"`
		Score  int   `json:"score"`
	}
	params := &Params{}
	if err := ctx.Bind(params); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	user, err := sCtl.userSvr.GetUserById(params.UserId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return sCtl.Fail(ctx, echoapp.CodeNotFound, "用户不存在", err)
		} else {
			return sCtl.Fail(ctx, echoapp.CodeDBError, "系统异常", err)
		}
	}
	sCtl.userSvr.AddScore(ctx, user, params.Score)
	return sCtl.Success(ctx, nil)
}

func (sCtl *UserController) Login(ctx echo.Context) error {
	param := &echoapp.LoginParam{}
	if err := ctx.Bind(param); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "登录失败", err)
	}
	//comId, err := strconv.Atoi(ctx.QueryParam("com_id"))
	//if err != nil {
	//	return sCtl.Fail(ctx, echoapp.CodeArgument, "参数校验失败", err)
	//}
	//
	//param.ComId = comId
	user, err := sCtl.userSvr.Login(ctx, param)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeInnerError, "登录失败", err)
	}
	return sCtl.Success(ctx, user)
}

func (sCtl *UserController) Register(ctx echo.Context) error {
	return nil
}

func (sCtl *UserController) GetUserRoles(ctx echo.Context) error {
	return nil
}

func (sCtl *UserController) Logout(ctx echo.Context) error {
	return nil
}

func (sCtl *UserController) SendVerifyCodeSms(ctx echo.Context) error {
	//com, err := util.GetCtxCompany(ctx)
	//if err != nil {
	//	return sCtl.Fail(ctx, echoapp.CodeArgument, "", err)
	//}
	//
	//options := &echoapp.SendMessageOptions{
	//	Token:         "",
	//	ComId:         0,
	//	PhoneNumbers:  nil,
	//	SignName:      "",
	//	TemplateCode:  "",
	//	TemplateParam: "",
	//}
	//sCtl.smsSvr.SendMessage(options)
	return nil
}

func (sCtl *UserController) GetVerifyPic(ctx echo.Context) error {
	return nil
}

func (sCtl *UserController) GetUserInfo(ctx echo.Context) error {
	//echoapp_util.ExtractEntry(ctx).Info("getUserInfo")
	user, err := util.GetCtxtUser(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeNotFound, "未发现用户", err)
	}
	return sCtl.Success(ctx, user)
}

func (sCtl *UserController) CheckHasRoles(ctx echo.Context) error {
	return nil
}

func (sCtl *UserController) Jscode2session(ctx echo.Context) error {
	comId := util.GetCtxComId(ctx)
	code := ctx.QueryParam("code")
	sCtl.userSvr.Jscode2session(comId, code)
	return nil
}

func (sCtl *UserController) GetUserAddressList(ctx echo.Context) error {
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	addressList, err := sCtl.userSvr.GetUserAddressList(userId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	return sCtl.Success(ctx, addressList)
}
func (sCtl *UserController) GetUserDefaultAddress(ctx echo.Context) error {
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	address, err := sCtl.userSvr.GetCachedUserDefaultAddrById(userId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeCacheError, err.Error(), err)
	}
	return sCtl.Success(ctx, address)
}

func (sCtl *UserController) CreateUserAddress(ctx echo.Context) error {
	addr := &echoapp.Address{}
	if err := ctx.Bind(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	addr.UserID = userId
	if err := sCtl.userSvr.CreateUserAddress(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, addr)
}

func (sCtl *UserController) UpdateUserAddress(ctx echo.Context) error {
	addrParam := &echoapp.Address{}
	if err := ctx.Bind(addrParam); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	addrParam.UserID = userId
	//addrId, _ := echoapp_util.GetCtxtAddrId(ctx)
	//addrParam.AddrId = addrId

	// addr, err := sCtl.userSvr.GetUserAddrById(addrParam.AddrId)
	// if err != nil {
	// 	return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	// }
	// addr.Username = addrParam.Username
	// addr.Mobile = addrParam.Mobile
	// addr.Address = addrParam.Address
	// addr.Checked = addrParam.Checked
	if err := sCtl.userSvr.UpdateUserAddress(addrParam); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, echoapp.ErrDb.Error(), err)
	}
	return sCtl.Success(ctx, addrParam)
}

func (sCtl *UserController) DelUserAddress(ctx echo.Context) error {
	addrParam := &echoapp.Address{}
	if err := ctx.Bind(addrParam); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	var addr *echoapp.Address
	addr, err := sCtl.userSvr.GetUserAddrById(addrParam.AddrId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	if err := sCtl.userSvr.DelUserAddress(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, echoapp.ErrDb.Error(), err)
	}
	return sCtl.Success(ctx, nil)
}

func (sCtl *UserController) GetCartGoodsList(context echo.Context) error {
	return nil
}

func (sCtl *UserController) AddCartGoods(context echo.Context) error {
	return nil
}

func (sCtl *UserController) DelCartGoods(context echo.Context) error {
	return nil
}

func (sCtl *UserController) UpdateCartGoods(context echo.Context) error {
	return nil
}

func (sCtl *UserController) GetUserCollectionList(ctx echo.Context) error {
	lastId, limitint := echoapp_util.GetCtxListParams(ctx)
	// limit := ctx.QueryParam("limit")
	// limitint, _ := strconv.Atoi(limit)
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	addressList, err := sCtl.userSvr.GetUserCollectionList(userId, lastId, limitint)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	return sCtl.Success(ctx, addressList)
}
func (sCtl *UserController) GetCacheUserCollectionById(ctx echo.Context) error {
	targetId := ctx.QueryParam("targetId")
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	res, err := sCtl.userSvr.GetCachedUserCollectionById(userId, targetId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	return sCtl.Success(ctx, res)
}
func (sCtl *UserController) GetUserCacheCollectionList(ctx echo.Context) error {
	// limit := ctx.QueryParam("limit")
	// limitint, _ := strconv.Atoi(limit)
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	addressList, err := sCtl.userSvr.GetCachedUserCollectionListById(userId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeCacheError, err.Error(), err)
	}
	return sCtl.Success(ctx, addressList)
}
func (sCtl *UserController) CreateUserCollection(ctx echo.Context) error {
	addr := &echoapp.Collection{}
	if err := ctx.Bind(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	addr.UserID = userId
	if err := sCtl.userSvr.CreateUserCollection(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, addr)
}

func (sCtl *UserController) DelUserCollection(ctx echo.Context) error {
	targetId := ctx.QueryParam("targetId")
	targetIdInt, _ := strconv.ParseInt(targetId, 10, 64)
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	res, err := sCtl.userSvr.GetUserCollectionById(targetIdInt, userId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	if err := sCtl.userSvr.DelUserCollection(res); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, nil)
}
