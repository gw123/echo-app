package controllers

import (
	"strconv"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type OrderController struct {
	orderSvr echoapp.OrderService
	userSvr  echoapp.UserService
	echoapp.BaseController
}

func NewOrderController(orderSvr echoapp.OrderService, userSvr echoapp.UserService) *OrderController {
	return &OrderController{
		orderSvr: orderSvr,
		userSvr:  userSvr,
	}
}

func (oCtl *OrderController) GetOrdereById(c echo.Context) error {
	param := &Param{}
	if err := c.Bind(param); err != nil {
		return oCtl.Fail(c, echoapp.CodeArgument, "", errors.Wrap(err, "Bind"))
	}
	res, err := oCtl.orderSvr.GetOrderById(param.ID)
	if err != nil {
		return oCtl.Fail(c, echoapp.CodeDBError, "", errors.Wrap(err, "GetResourceById"))
	}
	return oCtl.Success(c, res)
}

type Param struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	From   int  `json:"from"`
	Limit  int  `json:"limit"`
	TagID  uint `json:"tag_id"`
	Status uint `json:"status"`
}

func (oCtl *OrderController) GetUserPaymentOrder(c echo.Context) error {
	params := &Param{}
	if err := c.Bind(params); err != nil {
		return oCtl.Fail(c, echoapp.CodeArgument, "", errors.Wrap(err, "Bind"))
	}
	res, err := oCtl.orderSvr.GetUserPaymentOrder(c, params.UserID, params.From, params.Limit)
	if err != nil {
		return oCtl.Fail(c, echoapp.CodeDBError, "", errors.Wrap(err, "GetUserPaymentOrder"))
	}
	return oCtl.Success(c, res)
}

func (oCtl *OrderController) GetOrderList(c echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(c)
	if err != nil {
		return oCtl.Fail(c, echoapp.CodeArgument, "授权失败", err)
	}
	last, limit := echoapp_util.GetCtxListParams(c)
	status := c.QueryParam("status")
	filelist, err := oCtl.orderSvr.GetUserOrderList(c, uint(userID), status, last, limit)
	if err != nil {
		return oCtl.Fail(c, echoapp.CodeArgument, "OrderCtrl->GetOrderList", err)
	}
	echoapp_util.ExtractEntry(c).Infof("status: %s from:%s,limit:%s", status, last, limit)
	return oCtl.Success(c, filelist)
}

func (oCtl *OrderController) GetTicketByCode(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	company, err := echoapp_util.GetCtxCompany(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	codeTicket, err := oCtl.orderSvr.GetTicketByCode(code)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	codeTicket.XcxCover = company.XcxCover
	return oCtl.Success(ctx, codeTicket)
}

func (oCtl *OrderController) GetOrderStatistics(ctx echo.Context) error {
	return nil
}

func (oCtl *OrderController) PreOrder(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		glog.Info("preOrder argument err")
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		glog.Error("preOrder getCtxUser err")
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	oCtl.setOrderInfoFromCtx(ctx, params)

	addr, err := oCtl.userSvr.GetUserAddrById(int64(params.AddressId))
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, "无效的收货地址", err)
	}
	params.Address = addr

	if err := oCtl.orderSvr.PreCheckOrder(params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ExpressStatus = echoapp.OrderStatusToShip
	glog.Info("begin uniPreOrder")
	resp, err := oCtl.orderSvr.UniPreOrder(params, user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return oCtl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
		} else {
			return oCtl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
		}
	}
	return oCtl.Success(ctx, resp)
}

/***
开通会员
*/
func (oCtl *OrderController) OpenVip(ctx echo.Context) error {
	params := &echoapp.Order{}
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	if user.VipLevel > 0 {
		// 如果已经是会员不需要重复办理
		return oCtl.Fail(ctx, echoapp.CodeArgument, "您已经是会员", err)
	}

	if err := ctx.Bind(params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.GoodsType = echoapp.GoodsTypeVip
	oCtl.setOrderInfoFromCtx(ctx, params)
	if err := oCtl.orderSvr.PreCheckOrder(params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	resp, err := oCtl.orderSvr.UniPreOrder(params, user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return oCtl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
		} else {
			return oCtl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
		}
	}
	return oCtl.Success(ctx, resp)
}

func (oCtl *OrderController) setOrderInfoFromCtx(ctx echo.Context, order *echoapp.Order) {
	order.ComId = echoapp_util.GetCtxComId(ctx)
	user, _ := echoapp_util.GetCtxtUser(ctx)
	order.UserId = uint(user.Id)

	if order.InviterId == uint(user.Id) {
		// 自己邀请自己不算
		order.InviterId = 0
	}

	clientType := echoapp_util.GetClientTypeByUA(ctx.Request().UserAgent())
	order.Source = clientType
	switch clientType {
	case echoapp.ClientWxOfficial:
		order.Source = "公众号"
	case echoapp.ClientWxMiniApp:
		order.Source = "小程序"
	}
	order.ClientType = clientType
	order.ClientIP = ctx.RealIP()
}

/***
  查询订单的支付结果
*/
func (oCtl *OrderController) QueryOrder(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ComId = echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(user.Id)
	}

	order, err := oCtl.orderSvr.GetUserOrderDetail(ctx, uint(user.Id), params.OrderNo)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "获取订单失败"+err.Error(), err)
	}

	// 办理会员订单没有收货地址
	if order.GoodsType != echoapp.GoodsTypeVip {
		addr := &echoapp.Address{}
		addr, err = oCtl.userSvr.GetUserAddrById(int64(order.AddressId))
		if err != nil {
			return oCtl.Fail(ctx, echoapp.CodeArgument, "错误的收货地址", err)
		}
		order.Address = addr
	}

	resp, err := oCtl.orderSvr.QueryOrderAndUpdate(order, echoapp.OrderStatusPaid)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return oCtl.Success(ctx, resp)
}

//func (orderCtrl *OrderController) PlaceSelfOrder(ctx echo.Context) error {
//	params := &echoapp.Order{}
//	if err := ctx.Bind(params); err != nil {
//		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
//	}
//	user, err := echoapp_util.GetCtxtUser(ctx)
//	if err != nil {
//		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
//	} else {
//		params.UserId = uint(user.Id)
//	}
//	_, err = orderCtrl.orderSvr.UniPreOrder(params, user)
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return orderCtrl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
//		} else {
//			return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常", err)
//		}
//	}
//	return orderCtrl.Success(ctx, nil)
//}

func (oCtl *OrderController) CancelOrder(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ComId = echoapp_util.GetCtxComId(ctx)
	if userId, err := echoapp_util.GetCtxtUserId(ctx); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(userId)
	}

	err := oCtl.orderSvr.CancelOrder(params)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "系统异常", err)
	}
	return oCtl.Success(ctx, nil)
}

func (oCtl *OrderController) Refund(ctx echo.Context) error {
	params := &echoapp.Order{}
	if err := ctx.Bind(params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	params.ComId = echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	order, err := oCtl.orderSvr.GetUserOrderDetail(ctx, uint(user.Id), params.OrderNo)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, err.Error(), err)
	}
	if err := oCtl.orderSvr.Refund(order, user); err != nil {
		return err
	}
	return nil
}

func (oCtl *OrderController) GetOrderDetail(ctx echo.Context) error {
	orderNo := ctx.QueryParam("order_no")
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	order, err := oCtl.orderSvr.GetUserOrderDetail(ctx, uint(user.Id), orderNo)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "获取订单失败", err)
	}
	return oCtl.Success(ctx, order)
}

type CheckTicketRequest struct {
	Code string `json:"code"`
	Num  uint   `json:"num"`
}

func (oCtl *OrderController) CheckTicket(ctx echo.Context) error {
	request := CheckTicketRequest{}
	if err := ctx.Bind(&request); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, "授权失败", err)
	}

	if err := oCtl.orderSvr.CheckTicket(request.Code, request.Num, uint(user.Id)); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, "门票校验失败", err)
	}
	return oCtl.Success(ctx, nil)
}

func (oCtl *OrderController) CheckTicketList(ctx echo.Context) error {
	return nil
}

func (oCtl *OrderController) GetTicketList(ctx echo.Context) error {
	return nil
}

func (oCtl *OrderController) GetTicketDetail(ctx echo.Context) error {
	return nil
}

func (oCtl *OrderController) FetchThirdTicket(ctx echo.Context) error {
	return nil
}

func (oCtl *OrderController) CheckTicketByStaff(ctx echo.Context) error {
	return nil
}

func (oCtl *OrderController) CheckTicketBySelf(ctx echo.Context) error {
	return nil
}

func (oCtl *OrderController) WxPayCallback(ctx echo.Context) error {
	err := oCtl.orderSvr.WxPayCallback(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return oCtl.Success(ctx, "SUCCESS")
}

func (oCtl *OrderController) WxRefundCallback(ctx echo.Context) error {
	err := oCtl.orderSvr.WxRefundCallback(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return oCtl.Success(ctx, "SUCCESS")
}

func (oCtl *OrderController) QueryRefund(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ComId = echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(user.Id)
	}

	order, err := oCtl.orderSvr.GetUserOrderDetail(ctx, uint(user.Id), params.OrderNo)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "获取订单失败"+err.Error(), err)
	}

	addr := &echoapp.Address{}
	addr, err = oCtl.userSvr.GetUserAddrById(int64(order.AddressId))
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, "错误的收货地址", err)
	}
	order.Address = addr

	resp, err := oCtl.orderSvr.QueryRefundOrderAndUpdate(order)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return oCtl.Success(ctx, resp)
}

func (oCtl *OrderController) Appointment(ctx echo.Context) error {
	params := &echoapp.Appointment{}

	if err := ctx.Bind(params); err != nil {
		echoapp_util.ExtractEntry(ctx).WithError(err).Error("appointment argument err")
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		echoapp_util.ExtractEntry(ctx).WithError(err).Error("appointment getCtxUser err")
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	comID := echoapp_util.GetCtxComId(ctx)
	params.ComID = comID
	params.Username = user.Name
	params.UserID = uint(user.Id)
	if params.AddressId != 0 {
		// 身份信息来自Address选择
		addr, err := oCtl.userSvr.GetUserAddrById(int64(params.AddressId))
		if err != nil {
			echoapp_util.ExtractEntry(ctx).WithError(err).Error("appointment getUserAddrById err")
			return oCtl.Fail(ctx, echoapp.CodeArgument, "无效的收货地址", err)
		}
		params.IDCard = addr.Code
		params.IDCardType = echoapp.IDCardTypeID
	}
	params.Status = echoapp.AppointmentStatusUnused

	if err := oCtl.orderSvr.Appointment(ctx, params); err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	return oCtl.Success(ctx, nil)
}

func (oCtl *OrderController) GetAppointmentDetail(ctx echo.Context) error {
	idStr := ctx.QueryParam("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	appointment, err := oCtl.orderSvr.GetAppointmentDetail(ctx, id)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "获取预约详情失败", err)
	}

	if int64(appointment.UserID) != user.Id {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "获取预约详情失败", errors.New("没有权限"))
	}

	return oCtl.Success(ctx, appointment)
}

func (oCtl *OrderController) GetAppointmentList(ctx echo.Context) error {
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeArgument, "获取用户信息失败", err)
	}

	status := ctx.QueryParam("status")
	comID := echoapp_util.GetCtxComId(ctx)
	lastID, _ := echoapp_util.GetCtxListParams(ctx)
	list, err := oCtl.orderSvr.GetAppointmentList(ctx, comID, uint(user.Id), lastID, status)
	if err != nil {
		return oCtl.Fail(ctx, echoapp.CodeInnerError, "获取预约详情失败", err)
	}
	return oCtl.Success(ctx, list)
}
