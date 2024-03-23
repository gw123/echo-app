package echoapp

import (
	"context"
	"net/http"
	"time"

	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/labstack/echo"
	"github.com/silenceper/wechat/v2/officialaccount/js"
	"github.com/silenceper/wechat/v2/officialaccount/server"
)

const WxAuthCallBack = "wxAuthCallBack"

type WxPreOrderResponse struct {
	AppID     string `json:"appId"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
	Timestamp string `json:"timestamp"`
}

/***
微信模板消息类型
*/
type TemplateType struct {
	ID                 uint              `json:"id"`
	Type               string            `json:"type"`
	ComID              uint              `json:"com_id"`
	TemplateID         string            `json:"template_id"`
	KeywordColorMap    map[string]string `json:"color_map" gorm:"-"` //模板消息关键字颜色
	KeywordColorMapStr string            `gorm:"color_map" json:"-"`
	CreatedAt          time.Time         `json:"created_at"`
	DeletedAt          *time.Time        `json:"deleted_at"`
}

func (w *TemplateType) GetKeywordColor(key string) string {
	if w == nil || w.KeywordColorMap == nil {
		return ""
	}
	return w.KeywordColorMap[key]
}

type TemplateDataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

type TemplateDataItemMap map[string]*TemplateDataItem

type BaseTemplateMessage struct {
	ComID  uint
	Openid string
	First  string
	Remark string
}

func (b *BaseTemplateMessage) GetOpenid() string {
	return b.Openid
}

func (b *BaseTemplateMessage) GetComID() uint {
	return b.ComID
}

// TemplateMessage
type TemplateMessage interface {
	GetItems() TemplateDataItemMap
	GetMsgType() string
	GetOpenid() string
	GetComID() uint
	GetUrl() string
	GetMiniAppID() string // 当返回空跳转到公众号网页 ，如果配置跳转到小程序
}

type WechatService interface {
	GetOfficialServer(ctx echo.Context, comID uint) (*server.Server, error)
	GetAuthCodeUrl(comId uint, uri string) (url string, err error)
	GetUserInfo(ctx context.Context, comId uint, code string) (*mpoauth2.UserInfo, error)
	UnifiedOrder(ctx echo.Context, order *Order, openId string) (*WxPreOrderResponse, error)
	QueryOrder(order *Order) (string, error)
	PayCallback(r *http.Request) (*wechat.NotifyRequest, error)
	GetJsConfig(comID uint, url string) (*js.Config, error)
	Refund(order *Order, openId string) (*wechat.RefundResponse, error)
	RefundCallback(r *http.Request) (*wechat.RefundNotify, error)
	QueryRefund(order *Order) (string, error)
	SendTplMessage(ctx context.Context, message TemplateMessage) (int64, error)
}
