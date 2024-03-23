package controllers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/labstack/echo"
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
		return sCtl.AppErr(ctx, echoapp.AppErrNotLogin.WithInner(err))
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

//获取商品页面关联商品的一个活动
func (sCtl *ActivityController) GetActivityByGoodsId(ctx echo.Context) error {
	goodsId, err := strconv.Atoi(ctx.QueryParam("goods_id"))
	if goodsId == 0 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	act, err := sCtl.actSvr.GetGoodsActivity(uint(goodsId))

	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	return sCtl.Success(ctx, act)
}

// 获取用户的奖品
func (sCtl *ActivityController) GetUserAwards(ctx echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.AppErr(ctx, echoapp.AppErrNotLogin.WithInner(err))
	}
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	userAwards, err := sCtl.actSvr.GetUserAwards(uint(userID), lastId, uint(limit))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, userAwards)
}

// 奖品领取的历史记录
func (sCtl *ActivityController) GetUserAwardHistories(ctx echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.AppErr(ctx, echoapp.AppErrNotLogin.WithInner(err))
	}
	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	userAwardsHistories, err := sCtl.actSvr.GetAwardHistoryByUserID(uint(userID), lastId, uint(limit))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, userAwardsHistories)
}

// 奖品领取的历史记录
func (sCtl *ActivityController) StaffCheckedAwards(ctx echo.Context) error {
	comID := echoapp_util.GetCtxComId(ctx)
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.AppErr(ctx, echoapp.AppErrNotLogin.WithInner(err))
	}

	lastId, limit := echoapp_util.GetCtxListParams(ctx)
	userAwardsHistories, err := sCtl.actSvr.StaffCheckedAwards(comID, uint(userID), lastId, uint(limit))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeDBError, "系统错误", err)
	}
	return sCtl.Success(ctx, userAwardsHistories)
}

// 获取一个商品的授权码
func (sCtl *ActivityController) GetUserAwardCode(ctx echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return sCtl.AppErr(ctx, echoapp.AppErrNotLogin.WithInner(err))
	}

	userAwardID, err := strconv.Atoi(ctx.QueryParam("id"))
	if err != nil || userAwardID == 0 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "id参数错误", err)
	}

	num, err := strconv.Atoi(ctx.QueryParam("num"))
	if err != nil || num == 0 {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "num参数错误", err)
	}

	award, err := sCtl.actSvr.GetUserAward(uint(userAwardID))
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	echoapp_util.ExtractEntry(ctx).Infof("award %+v", award)
	if uint(userID) != award.UserID {
		return sCtl.Fail(ctx, echoapp.CodeNoAuth, "没有权限", echoapp.AppErrNotAuth)
	}

	code, err := echoapp_util.EncodeInt64(int64(award.ID), echoapp.ConfigOpts.Jws.HashIdsSalt)
	if err != nil {
		return sCtl.Fail(ctx, echoapp.CodeArgument, "code编码失败", err)
	}
	code = fmt.Sprintf("%s-%d-%d", code, num, time.Now().Unix())
	resp := map[string]string{
		"code": code,
	}
	return sCtl.Success(ctx, resp)
}

// 获取一个商品的授权码
func (sCtl *ActivityController) CheckAwardByCode(ctx echo.Context) error {
	echoapp_util.ExtractEntry(ctx).Infof("checkAwardByCode")
	staffID, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		echoapp_util.ExtractEntry(ctx).Error("getCtxUserId")
		return sCtl.AppErr(ctx, echoapp.AppErrNotLogin.WithInner(err))
	}

	// todo 检查角色 是否为老板或是员工
	code := ctx.QueryParam("code")
	if err != nil {
		echoapp_util.ExtractEntry(ctx).Error("query param")
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	tempArr := strings.Split(code, "-")
	if len(tempArr) != 3 {
		echoapp_util.ExtractEntry(ctx).Error("query param len != 3, code", code)
		return sCtl.Fail(ctx, echoapp.CodeArgument, "参数错误2", err)
	}

	code = tempArr[0]
	numStr := tempArr[1]
	timeStr := tempArr[2]

	num, _ := strconv.Atoi(numStr)
	updateTime, _ := strconv.Atoi(timeStr)
	if time.Now().Unix()-int64(updateTime) > 3600 {
		echoapp_util.ExtractEntry(ctx).Error("code overdue")
		return sCtl.Fail(ctx, echoapp.CodeArgument, "核销码已经过期", err)
	}
	echoapp_util.ExtractEntry(ctx).Infof("code %s,num:%d ,time:%d", code, num, updateTime)
	userAwardID, err := echoapp_util.DecodeInt64(code, echoapp.ConfigOpts.Jws.HashIdsSalt)
	if err != nil {
		echoapp_util.ExtractEntry(ctx).Error("code overdue")
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	echoapp_util.ExtractEntry(ctx).Infof("userawardId:%d", userAwardID)
	award, err := sCtl.actSvr.GetUserAward(uint(userAwardID))
	if err != nil {
		echoapp_util.ExtractEntry(ctx).Error("code overdue")
		return sCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	if err := sCtl.actSvr.CheckUserAward(uint(staffID), award.ID, num); err != nil {
		echoapp_util.ExtractEntry(ctx).Error("code overdue")
		return sCtl.Fail(ctx, echoapp.CodeInnerError, err.Error(), err)
	}

	return sCtl.Success(ctx, nil)
}
