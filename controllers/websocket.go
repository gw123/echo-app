package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/services"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
)

type WsController struct {
	echoapp.BaseController
	wsSvr  echoapp.WsService
	usrSvr echoapp.UserService
}

func NewWsController(usrSvr echoapp.UserService) *WsController {
	temp := new(WsController)
	temp.wsSvr = services.NewWsService()
	temp.usrSvr = usrSvr
	return temp
}

func (c *WsController) Index(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "index", nil)
}

func (c *WsController) CreateWsClient(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	echoapp_util.ExtractEntry(ctx).Infof("CreateWsClient: %s", token)
	user, err := c.usrSvr.GetUserByToken(token)
	if err != nil {
		return c.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}
	echoapp_util.ExtractEntry(ctx).Infof("user_id:%d 登录成功", user.Id)
	return c.wsSvr.AddWsClient(ctx, token)
}

func (c *WsController) SendCmd(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	cmdParams := &echoapp.WsEventCmd{}
	if err := ctx.Bind(cmdParams); err != nil {
		return c.Fail(ctx, echoapp.CodeArgument, "参数错误", err)
	}

	cmdParams.EventType = echoapp.WsEventTypeCmd
	cmdParams.CreatedAt = time.Now().Unix()
	cmdParams.RequestId = fmt.Sprintf("%d", rand.Int())
	return c.wsSvr.SendWsClientEvent(cmdParams, token)
}
