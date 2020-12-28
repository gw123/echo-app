package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/gw123/glog"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type OrderController struct {
	orderSvc echoapp.OrderService
	userSvr  echoapp.UserService
	echoapp.BaseController
}

func NewOrderController(orderSvr echoapp.OrderService, userSvr echoapp.UserService) *OrderController {
	return &OrderController{
		orderSvc: orderSvr,
		userSvr:  userSvr,
	}
}

func (orderCtrl *OrderController) GetOrdereById(c echo.Context) error {
	param := &Param{}
	if err := c.Bind(param); err != nil {
		return orderCtrl.Fail(c, echoapp.CodeArgument, "", errors.Wrap(err, "Bind"))
	}
	res, err := orderCtrl.orderSvc.GetOrderById(param.ID)
	if err != nil {
		return orderCtrl.Fail(c, echoapp.CodeDBError, "", errors.Wrap(err, "GetResourceById"))
	}
	return orderCtrl.Success(c, res)
}

type Param struct {
	ID     uint `json:"id"`
	UserID uint `json:"user_id"`
	From   int  `json:"from"`
	Limit  int  `json:"limit"`
	TagID  uint `json:"tag_id"`
	Status uint `json:"status"`
}

func (orderCtrl *OrderController) GetUserPaymentOrder(c echo.Context) error {
	params := &Param{}
	if err := c.Bind(params); err != nil {
		return orderCtrl.Fail(c, echoapp.CodeArgument, "", errors.Wrap(err, "Bind"))
	}
	res, err := orderCtrl.orderSvc.GetUserPaymentOrder(c, params.UserID, params.From, params.Limit)
	if err != nil {
		return orderCtrl.Fail(c, echoapp.CodeDBError, "", errors.Wrap(err, "GetUserPaymentOrder"))
	}
	return orderCtrl.Success(c, res)
}

func (orderCtrl *OrderController) GetOrderList(c echo.Context) error {
	userID, err := echoapp_util.GetCtxtUserId(c)
	if err != nil {
		return orderCtrl.Fail(c, echoapp.CodeArgument, "授权失败", err)
	}
	last, limit := echoapp_util.GetCtxListParams(c)
	status := c.QueryParam("status")
	filelist, err := orderCtrl.orderSvc.GetUserOrderList(c, uint(userID), status, last, limit)
	if err != nil {
		return orderCtrl.Fail(c, echoapp.CodeArgument, "OrderCtrl->GetOrderList", err)
	}
	echoapp_util.ExtractEntry(c).Infof("status: %s from:%s,limit:%s", status, last, limit)
	return orderCtrl.Success(c, filelist)
}

func (orderCtrl *OrderController) GetTicketByCode(ctx echo.Context) error {
	code := ctx.QueryParam("code")
	company, err := echoapp_util.GetCtxCompany(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	codeTicket, err := orderCtrl.orderSvc.GetTicketByCode(code)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeDBError, err.Error(), err)
	}
	codeTicket.XcxCover = company.XcxCover
	return orderCtrl.Success(ctx, codeTicket)
}

func (orderCtrl *OrderController) GetOrderStatistics(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) PreOrder(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ComId = echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(user.Id)
	}

	addr, err := orderCtrl.userSvr.GetUserAddrById(int64(params.AddressId))
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, "无效的收货地址", err)
	}
	params.Address = addr

	if params.InviterId == uint(user.Id) {
		// 自己邀请自己不算
		params.InviterId = 0
	}

	if params.InviterId != 0 {
		glog.Infof("inviterId :%d ", params.InviterId)
	}

	clientType := echoapp_util.GetClientTypeByUA(ctx.Request().UserAgent())
	switch clientType {
	case echoapp.ClientWxOfficial:
		params.Source = "公众号"
	case echoapp.ClientWxMiniApp:
		params.Source = "小程序"
	}

	params.ClientType = clientType
	params.ClientIP = ctx.RealIP()
	if err := orderCtrl.orderSvc.PreCheckOrder(params); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	resp, err := orderCtrl.orderSvc.UniPreOrder(params, user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return orderCtrl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
		} else {
			return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
		}
	}
	return orderCtrl.Success(ctx, resp)
}

/***
  查询订单的支付结果
*/
func (orderCtrl *OrderController) QueryOrder(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ComId = echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(user.Id)
	}

	order, err := orderCtrl.orderSvc.GetUserOrderDetial(ctx, uint(user.Id), params.OrderNo)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "获取订单失败"+err.Error(), err)
	}

	addr := &echoapp.Address{}
	addr, err = orderCtrl.userSvr.GetUserAddrById(int64(order.AddressId))
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, "错误的收货地址", err)
	}
	order.Address = addr

	resp, err := orderCtrl.orderSvc.QueryOrderAndUpdate(order, echoapp.OrderStatusPaid)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return orderCtrl.Success(ctx, resp)
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
//	_, err = orderCtrl.orderSvc.UniPreOrder(params, user)
//	if err != nil {
//		if err == gorm.ErrRecordNotFound {
//			return orderCtrl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
//		} else {
//			return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常", err)
//		}
//	}
//	return orderCtrl.Success(ctx, nil)
//}

func (orderCtrl *OrderController) CancelOrder(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ComId = echoapp_util.GetCtxComId(ctx)
	if userId, err := echoapp_util.GetCtxtUserId(ctx); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(userId)
	}

	err := orderCtrl.orderSvc.CancelOrder(params)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常", err)
	}
	return orderCtrl.Success(ctx, nil)
}

func (orderCtrl *OrderController) Refund(ctx echo.Context) error {
	params := &echoapp.Order{}
	if err := ctx.Bind(params); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	params.ComId = echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	order, err := orderCtrl.orderSvc.GetUserOrderDetial(ctx, uint(user.Id), params.OrderNo)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, err.Error(), err)
	}
	if err := orderCtrl.orderSvc.Refund(order, user); err != nil {
		return err
	}
	return nil
}

func (orderCtrl *OrderController) GetOrderDetail(ctx echo.Context) error {
	orderNo := ctx.QueryParam("order_no")
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	order, err := orderCtrl.orderSvc.GetUserOrderDetial(ctx, uint(user.Id), orderNo)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "获取订单失败", err)
	}
	return orderCtrl.Success(ctx, order)
}

type CheckTicketRequest struct {
	Code string `json:"code"`
	Num  uint   `json:"num"`
}

func (orderCtrl *OrderController) CheckTicket(ctx echo.Context) error {
	request := CheckTicketRequest{}
	if err := ctx.Bind(&request); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, "授权失败", err)
	}

	if err := orderCtrl.orderSvc.CheckTicket(request.Code, request.Num, uint(user.Id)); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, "门票校验失败", err)
	}
	return orderCtrl.Success(ctx, nil)
}

func (orderCtrl *OrderController) CheckTicketList(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) GetTicketList(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) GetTicketDetail(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) FetchThirdTicket(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) CheckTicketByStaff(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) CheckTicketBySelf(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) WxPayCallback(ctx echo.Context) error {
	err := orderCtrl.orderSvc.WxPayCallback(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return orderCtrl.Success(ctx, "SUCCESS")
}

func (orderCtrl *OrderController) WxRefundCallback(ctx echo.Context) error {
	err := orderCtrl.orderSvc.WxRefundCallback(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return orderCtrl.Success(ctx, "SUCCESS")
}

func (orderCtrl *OrderController) QueryRefund(ctx echo.Context) error {
	params := &echoapp.Order{}

	if err := ctx.Bind(params); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	params.ComId = echoapp_util.GetCtxComId(ctx)
	user, err := echoapp_util.GetCtxtUser(ctx)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(user.Id)
	}

	order, err := orderCtrl.orderSvc.GetUserOrderDetial(ctx, uint(user.Id), params.OrderNo)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "获取订单失败"+err.Error(), err)
	}

	addr := &echoapp.Address{}
	addr, err = orderCtrl.userSvr.GetUserAddrById(int64(order.AddressId))
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, "错误的收货地址", err)
	}
	order.Address = addr

	resp, err := orderCtrl.orderSvc.QueryRefundOrderAndUpdate(order)
	if err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常"+err.Error(), err)
	}
	return orderCtrl.Success(ctx, resp)
}
