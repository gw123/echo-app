package controllers

import (
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type OrderController struct {
	orderSvc echoapp.OrderService

	echoapp.BaseController
}

func NewOrderController(orderSvr echoapp.OrderService) *OrderController {
	return &OrderController{
		orderSvc: orderSvr,
	}
}

func (orderCtrl *OrderController) PlaceOrder(ctx echo.Context) error {
	params := &echoapp.Order{}
	if err := ctx.Bind(params); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	err := orderCtrl.orderSvc.PlaceOrder(params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return orderCtrl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
		} else {
			return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常", err)
		}
	}
	return orderCtrl.Success(ctx, nil)
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
	last,limit := echoapp_util.GetCtxListParams(c)
	status := c.QueryParam("status")
	filelist, err := orderCtrl.orderSvc.GetUserOrderList(c, uint(userID),status, last, limit)
	if err != nil {
		return orderCtrl.Fail(c, echoapp.CodeArgument, "OrderCtrl->GetOrderList", err)
	}
	echoapp_util.ExtractEntry(c).Infof("status: %s from:%s,limit:%s", status, last,limit)
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
	if userId, err := echoapp_util.GetCtxtUserId(ctx); err != nil {
		return orderCtrl.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	} else {
		params.UserId = uint(userId)
	}

	err := orderCtrl.orderSvc.PlaceOrder(params)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return orderCtrl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
		} else {
			return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常", err)
		}
	}
	return orderCtrl.Success(ctx, nil)
}

func (orderCtrl *OrderController) CreateOrder(ctx echo.Context) error {

	return nil
}

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
	return nil
}

func (orderCtrl *OrderController) GetOrderDetail(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) CheckTicket(ctx echo.Context) error {

	return nil
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

//
func (orderCtrl *OrderController) GetCartGoodsList(ctx echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) AddCartGoods(context echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) DelCartGoods(context echo.Context) error {
	return nil
}

func (orderCtrl *OrderController) UpdateCartGoods(context echo.Context) error {
	return nil
}
