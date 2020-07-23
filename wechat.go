package echoapp

import (
	"context"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
)

const WxAuthCallBack = "wxAuthCallBack"

type WechatService interface {
	GetAuthCodeUrl(comId uint) (url string, err error)
	GetUserInfo(ctx context.Context, comId uint, code string) (*mpoauth2.UserInfo, error)
	UnifiedOrder(order *Order, openId string) (interface{}, error)
}
