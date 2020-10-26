package echoapp

import (
	"context"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/labstack/echo"
	"github.com/silenceper/wechat/v2/officialaccount/js"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"net/http"
)

const WxAuthCallBack = "wxAuthCallBack"

type WechatService interface {
	GetOfficialServer(ctx echo.Context, comID uint) (*server.Server, error)
	GetAuthCodeUrl(comId uint) (url string, err error)
	GetUserInfo(ctx context.Context, comId uint, code string) (*mpoauth2.UserInfo, error)
	UnifiedOrder(order *Order, openId string) (*wechat.UnifiedOrderResponse, error)
	QueryOrder(order *Order) (string, error)
	PayCallback(r *http.Request) (*wechat.NotifyRequest, error)
	GetJsConfig(comID uint, url string) (*js.Config, error)
	Refund(order *Order, openId string) (*wechat.RefundResponse, error)
	RefundCallback(r *http.Request) (*wechat.RefundNotify, error)
	QueryRefund(order *Order) (string, error)
}
