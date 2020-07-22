package controllers

import (
	"errors"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
	"strconv"
)

type ActivityController struct {
	actSvr echoapp.ActivityService
	echoapp.BaseController
}

func NewActivityController(comSvr echoapp.CompanyService, actSvr echoapp.ActivityService) *ActivityController {
	return &ActivityController{
		actSvr: actSvr,
	}
}

func (sCtl *ActivityController) GetActivityList(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	activityList, err := sCtl.actSvr.GetActivityList(comId, lastId, limit)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, activityList)
}

func (sCtl *ActivityController) GetActivityDetail(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	activity, err := sCtl.actSvr.GetActivityDetail(uint(id))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, activity)
}

func (sCtl *ActivityController) GetCouponsByGoodsId(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	goodsId, err := strconv.Atoi(ctx.QueryParam("goods_id"))
	if goodsId == 0 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	coupons, err := sCtl.actSvr.GetCouponsByGoodsId(comId, uint(goodsId))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, coupons)
}

func (sCtl *ActivityController) GetCouponsByPosition(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	position := ctx.QueryParam("position")
	if position == "" {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", errors.New("position is null"))
	}

	coupons, err := sCtl.actSvr.GetCouponsByPosition(comId, position)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, coupons)
}

func (sCtl *ActivityController) GetUserCoupons(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	lastId, _ := echoapp_util.GetCtxListParams(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	coupons, err := sCtl.actSvr.GetUserCoupons(comId, uint(userId), lastId)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, coupons)
}

func (sCtl *ActivityController) GetCouponsByOrder(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	order := &echoapp.Order{}
	if err := ctx.Bind(order); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	order.UserId = uint(userId)
	glog.Infof("len : %d", len(order.GoodsList))
	glog.Infof("len : %+v", order)
	userCoupons, coupons, err := sCtl.actSvr.GetUserCouponsByOrder(comId, order)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}

	response := map[string]interface{}{
		"user":      userCoupons,
		"available": coupons,
	}
	return sCtl.Success(ctx, response)
}

type CreateUserCouponParams struct {
	CouponId uint `json:"coupon_id"`
	Position uint `json:"position"` //领取位置
	From     uint `json:"from"`     //从什么平台领取
}

func (sCtl *ActivityController) CreateUserCoupon(ctx echo.Context) error {
	comId := echoapp_util.GetCtxComId(ctx)
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	params := &CreateUserCouponParams{}
	if err := ctx.Bind(params); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}
	if err := sCtl.actSvr.CreateUserCoupon(comId, uint(userId), params.CouponId); err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	return sCtl.Success(ctx, nil)
}
