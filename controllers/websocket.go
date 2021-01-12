package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
)

type WsController struct {
	echoapp.BaseController
	wsSvr  echoapp.WsService
	usrSvr echoapp.UserService
}

func NewWsController(usrSvr echoapp.UserService, wsSvr echoapp.WsService) *WsController {
	temp := new(WsController)
	temp.wsSvr = wsSvr
	temp.usrSvr = usrSvr
	return temp
}

func (c *WsController) Index(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "index", nil)
}

func (c *WsController) CreateWsClient(ctx echo.Context) error {
	//token := ctx.QueryParam("token")
	clientID := ctx.QueryParam("ClientID")
	echoapp_util.ExtractEntry(ctx).Infof("创建新的客户端")
	userId, err := echoapp_util.GetCtxtUserId(ctx)
	if err != nil {
		return c.Fail(ctx, echoapp.CodeArgument, err.Error(), err)
	}

	echoapp_util.ExtractEntry(ctx).Infof("user_id:%d 登录成功", userId)
	if err := c.wsSvr.AddWsClient(ctx, uint(userId), clientID); err != nil {
		return c.Fail(ctx, echoapp.CodeArgument, "创建ws客户端失败", err)
	}

	return c.Success(ctx, nil)
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
