package echoapp

import (
	"context"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/iGoogle-ink/gopay/wechat"
)

const WxAuthCallBack = "wxAuthCallBack"

type WechatService interface {
	GetAuthCodeUrl(comId uint) (url string, err error)
	GetUserInfo(ctx context.Context, comId uint, code string) (*mpoauth2.UserInfo, error)
	UnifiedOrder(order *Order, openId string) (*wechat.UnifiedOrderResponse, error)
	QueryOrder(order *Order) (string, error)
}
