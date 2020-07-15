package controllers

import (
	"strconv"

	"github.com/gw123/glog"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserController struct {
	userSvr echoapp.UserService
	smsSvr  echoapp.SmsService
	goodSvr echoapp.GoodsService
	echoapp.BaseController
}

func NewUserController(usrSvr echoapp.UserService, goodsSvr echoapp.GoodsService) *UserController {
	return &UserController{
		userSvr: usrSvr,
		goodSvr: goodsSvr,
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
	addr, err := sCtl.userSvr.GetUserAddrById(int64(addrParam.ID))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	if err := sCtl.userSvr.DelUserAddress(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, echoapp.ErrDb.Error(), err)
	}
	return sCtl.Success(ctx, nil)
}

// func (sCtl *UserController) GetUserCollectionList(ctx echo.Context) error {
// 	lastId, limitint := echoapp_util.GetCtxListParams(ctx)
// 	// limit := ctx.QueryParam("limit")
// 	// limitint, _ := strconv.Atoi(limit)
// 	userId, err := echoapp_util.GetCtxtUserId(ctx)
// 	if err != nil {
// 		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
// 	}
// 	addressList, err := sCtl.userSvr.GetUserCollectionList(userId, lastId, limitint)
// 	if err != nil {
// 		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
// 	}
// 	return sCtl.Success(ctx, addressList)
// }
type CollectParams struct {
	TargetId uint   `json:"target_id"`
	Type     string `json:"type"`
}

func (sCtl *UserController) IsCollect(ctx echo.Context) error {
	params := &CollectParams{}
	if err := ctx.Bind(params); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}

	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	glog.Infof("是否收藏 UserID:%d,type:%s,targetId:%d", userId, params.Type, params.TargetId)
	res, err := sCtl.userSvr.IsCollect(userId, params.TargetId, params.Type)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	return sCtl.Success(ctx, res)
}
func (sCtl *UserController) GetUserCollectionList(ctx echo.Context) error {
	targetType := ctx.QueryParam("targetType")
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	lastCursor, _ := strconv.Atoi(ctx.QueryParam("last_id"))
	//lastId := uint64(lastCursor)
	collecttionList, err := sCtl.userSvr.GetCachedUserCollectionTypeSet(userId, targetType)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	type GoodsInfo struct {
		Price      float32 `json:"price"`
		Name       string  `json:"name"`
		SmallCover string  `json:"small_cover"`
		GoodsType  string  `json:"goods_type" `
	}
	var goodslist []*GoodsInfo
	var limit = 20
	var temp = lastCursor + limit
	if temp > len(collecttionList) {
		temp = len(collecttionList)
	}
	if targetType == "goods" {
		for i := lastCursor; i < temp; i++ {
			tempGoods := &GoodsInfo{}
			targetId, _ := strconv.Atoi(collecttionList[i])
			goods, err := sCtl.goodSvr.GetGoodsById(uint(targetId))
			//goods, err := sCtl.goodSvr.GetCachedGoodsById(uint(targetId))
			if err != nil {
				return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
			}
			tempGoods.Name = goods.Name
			tempGoods.Price = goods.Price
			tempGoods.GoodsType = goods.GoodsType
			tempGoods.SmallCover = goods.SmallCover
			goodslist = append(goodslist, tempGoods)
		}
	}
	collectionMap := make(map[string]interface{})
	collectionMap["userId"] = userId
	collectionMap["Type"] = targetType
	collectionMap["target"] = goodslist
	return sCtl.Success(ctx, collectionMap)
}
func (sCtl *UserController) AddUserCollection(ctx echo.Context) error {
	addr := &echoapp.Collection{}
	if err := ctx.Bind(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	addr.UserID = userId
	if addr.Type == "goods" {
		res, err := sCtl.goodSvr.GetGoodsById(addr.TargetId)
		if err != nil {
			return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
		}
		if res.Status != "publish" {
			return sCtl.Fail(ctx, echoapp.CodeArgument, "商品已下架", err)
		}
	}
	if err := sCtl.userSvr.CreateUserCollection(addr); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)

	}
	return sCtl.Success(ctx, addr)
}

type DelCollectionParams struct {
	Type     string `json:"type"`
	TargetId uint   `json:"target_id"`
}

func (sCtl *UserController) DelUserCollection(ctx echo.Context) error {
	params := &DelCollectionParams{}
	if err := ctx.Bind(params); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	if err := sCtl.userSvr.DelUserCollection(userId, params.Type, params.TargetId); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, nil)
}

func (sCtl *UserController) AddUserHistory(ctx echo.Context) error {
	his := &echoapp.History{}
	if err := ctx.Bind(his); err != nil {
		glog.Info("bind")
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		glog.Info("GetCtxtUserId")
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	his.UserID = userId
	if err := sCtl.userSvr.CreateUserHistory(his); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, his)
}
func (sCtl *UserController) GetUserHistoryList(ctx echo.Context) error {
	lastId, limitint := echoapp_util.GetCtxListParams(ctx)
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, echoapp.ErrArgument.Error(), err)
	}
	hisList, err := sCtl.userSvr.GetUserHistoryList(userId, lastId, limitint)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	type GoodsInfo struct {
		//BrowsTime string
		Price      float32 `json:"price"`
		Name       string  `json:"name"`
		SmallCover string  `json:"small_cover"`
		GoodsType  string  `json:"goods_type" `
	}
	hisResMap := make(map[string][]*GoodsInfo)
	var goodslist []*GoodsInfo
	hisListLen := len(hisList)
	browseTime := hisList[0].CreatedAt.Format("2006-01-02")
	goods, err := sCtl.goodSvr.GetGoodsById(hisList[0].TargetId)
	//goods, err := sCtl.goodSvr.GetCachedGoodsById(hisList[0].TargetId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	tempGoods := &GoodsInfo{}
	tempGoods.Name = goods.Name
	tempGoods.Price = goods.Price
	tempGoods.GoodsType = goods.GoodsType
	tempGoods.SmallCover = goods.SmallCover
	goodslist = append(goodslist, tempGoods)
	for i := 1; i < hisListLen; i++ {
		tempGoods := &GoodsInfo{}
		//if hisList[i].Type=="goods"{
		goods, err := sCtl.goodSvr.GetGoodsById(hisList[0].TargetId)
		//goods, err := sCtl.goodSvr.GetCachedGoodsById(hisList[i].TargetId)
		if err != nil {
			glog.Info("GetCachedGoodsById")
			continue
		}
		tempGoods.Name = goods.Name
		tempGoods.Price = goods.Price
		tempGoods.GoodsType = goods.GoodsType
		tempGoods.SmallCover = goods.SmallCover
		//}
		curTime := hisList[i].CreatedAt.Format("2006-01-02")
		if curTime == browseTime {
			goodslist = append(goodslist, tempGoods)
		} else {
			hisResMap[browseTime] = goodslist
			browseTime = curTime
			goodslist = nil
			goodslist = append(goodslist, tempGoods)
		}
		hisResMap[browseTime] = goodslist
	}
	return sCtl.Success(ctx, hisResMap)
}
