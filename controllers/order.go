package controllers

// import (
// 	echoapp "github.com/gw123/echo-app"
// 	echoapp_util "github.com/gw123/echo-app/util"
// 	"github.com/jinzhu/gorm"
// 	"github.com/labstack/echo"
// 	"github.com/pkg/errors"
// )

// type OrderController struct {
// 	orderSvc echoapp.OrderService

// 	echoapp.BaseController
// }

// func NewOrderController(orderSvr echoapp.OrderService) *OrderController {
// 	return &OrderController{
// 		orderSvc: orderSvr,
// 	}
// }

// func (orderCtrl *OrderController) PlaceOrder(ctx echo.Context) error {
// 	params := &echoapp.Order{}
// 	if err := ctx.Bind(params); err != nil {
// 		return orderCtrl.Fail(ctx, echoapp.Err, err.Error(), err)
// 	}
// 	err := orderCtrl.orderSvc.PlaceOrder(params)
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return orderCtrl.Fail(ctx, echoapp.CodeArgument, "用户不存在", err)
// 		} else {
// 			return orderCtrl.Fail(ctx, echoapp.CodeInnerError, "系统异常", err)
// 		}
// 	}
// 	return orderCtrl.Success(ctx, nil)
// }
// func (orderCtrl *OrderController) GetOrdereById(c echo.Context) error {
// 	param := &Param{}
// 	if err := c.Bind(param); err != nil {
// 		return orderCtrl.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "Bind"))
// 	}
// 	res, err := orderCtrl.orderSvc.GetOrderById(c, param.ID)
// 	if err != nil {
// 		return orderCtrl.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "GetResourceById"))
// 	}
// 	return orderCtrl.Success(c, res)
// }

// type Param struct {
// 	ID     uint `json:"id"`
// 	UserID uint `json:"user_id"`
// 	From   int  `json:"from"`
// 	Limit  int  `json:"limit"`
// 	TagID  uint `json:"tag_id"`
// }

// func (orderCtrl *OrderController) GetUserPaymentOrder(c echo.Context) error {
// 	params := &Param{}
// 	if err := c.Bind(params); err != nil {
// 		return orderCtrl.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "Bind"))
// 	}
// 	res, err := orderCtrl.orderSvc.GetUserPaymentOrder(c, params.UserID, params.From, params.Limit)
// 	if err != nil {
// 		return orderCtrl.Fail(c, echoapp.Err_Argument, "", errors.Wrap(err, "GetUserPaymentOrder"))
// 	}
// 	return orderCtrl.Success(c, res)
// }

// func (orderCtrl *OrderController) GetOrderList(c echo.Context) error {
// 	params := &Param{}
// 	filelist, err := orderCtrl.orderSvc.GetOrderList(c, params.From, params.Limit)
// 	if err != nil {
// 		return orderCtrl.Fail(c, echoapp.Err_Argument, "OrderCtrl->GetOrderList", err)
// 	}
// 	echoapp_util.ExtractEntry(c).Infof("from:%s,limit:%s", params.From, params.Limit)
// 	return orderCtrl.Success(c, filelist)
// }
