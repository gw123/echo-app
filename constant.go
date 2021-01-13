package echoapp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const (
	ServerIndex       = "http://m.xytschool.com/index-dev/%d#%s"
	TimeFormat        = "2006-01-02 15:04:05"
	TimeHourMinFormat = "2006-01-02 15:04"
	DateFormat        = "2006-01-02"
	HostURL           = "host_url"
	RequestRoot       = "static_root"
	AssetVersion      = "asset_version"

	//缓存key
	//redis 相关的key
	RedisUserKey          = "User:%d"
	RedisUserXCXOpenidKey = "UserXCXOpenid:%d"
	RedisSmsLoginCodeKey  = "SmsLoginCode"
	RedisCompanyKey       = "Company:%d"
	RedisGoodsKey         = "Goods:%d"
	//Goods:comId:position
	RedisBannerListKey = "Goods:%d:%s"

	RedisActivityKey = "Activity:%d"
	RedisArticleKey  = "Article:%d"

	RedisUserDefaultAddrKey = "UserDefaultAddr:%d"
	RedisUserCollectionKey  = "UserCollection:%d"

	ClientWxOfficial = "wx_official"
	ClientWxMiniApp  = "wx_miniapp"
	ClientAliMinApp  = "ali_minapp"
	ClientWap        = "h5"

	SourceMiniApp    = "小程序"
	SourceWxOfficial = "公众号"
	SourceTongcheng  = "同城"
	SourceMeituan    = "美团"
)

const (
	CodeNotFound   = 400
	CodeNoLogin    = 401
	CodeNoAuth     = 402
	CodeDBError    = 403
	CodeCacheError = 404
	CodeArgument   = 405
	CodeNotAllow   = 406
	CodeEtcdError  = 407
	CodeInnerError = 501
)

var ErrNotFoundCache = errors.New("缓冲不存在或者已经过期")
var ErrNotFoundDb = errors.New("未找到资源")
var ErrDb = errors.New("数据库错误")
var ErrNotFoundEtcd = errors.New("不存在的配置")
var ErrArgument = errors.New("参数错误")
var ErrNotLogin = errors.New("用户未登录,请登录后尝试")
var ErrNotAuth = errors.New("未授权")
var ErrNotAllow = errors.New("不允许的操作")
var ErrRefund = errors.New("订单已经退款")
var ErrTicketInvaild = errors.New("无效的门票")
var ErrTicketOverdue = errors.New("门票已经过期")
var ErrTicketUsed = errors.New("门票已经被使用")
var ErrTicketNotEnough = errors.New("已购门票数量不足,请核对数量")
var ErrOrderFormat = errors.New("订单格式不对")

type AppError interface {
	error
	GetOuter() string
	GetInner() error
	GetCode() int
	WithInner(err error) AppError
}

type appError struct {
	outer string
	inner error
	code  int
}

func NewAppError(code int, outer string, inner error) AppError {
	if outer == "" {
		outer = inner.Error()
	}
	return &appError{
		code:  code,
		outer: outer,
		inner: inner,
	}
}

func (a *appError) GetCode() int {
	return a.code
}

func (a *appError) WithInner(err error) AppError {
	a.inner = err
	return a
}

func (a *appError) Error() string {
	return a.inner.Error()
}

func (a *appError) GetInner() error {
	return a.inner
}

func (a *appError) GetOuter() string {
	return a.outer
}

var AppErrNotFoundCache = NewAppError(CodeNotFound, "", ErrNotFoundCache)
var AppErrNotFoundDb = NewAppError(CodeNotFound, "", ErrNotFoundDb)
var AppErrDb = NewAppError(CodeInnerError, "", ErrDb)
var AppErrNotFoundEtcd = NewAppError(CodeNotFound, "", ErrNotFoundEtcd)
var AppErrArgument = NewAppError(CodeArgument, "", ErrArgument)
var AppErrNotLogin = NewAppError(CodeNoLogin, "", ErrNotLogin)
var AppErrNotAuth = NewAppError(CodeNoAuth, "", ErrNotAuth)
var AppErrNotAllow = NewAppError(CodeNotAllow, "", ErrNotAllow)
var AppErrRefund = NewAppError(CodeNotAllow, "", ErrRefund)
var AppErrTicketInvaild = NewAppError(CodeNotAllow, "", ErrTicketInvaild)
var AppErrTicketOverdue = NewAppError(CodeNotAllow, "", ErrTicketOverdue)
var AppErrTicketUsed = NewAppError(CodeNotAllow, "", ErrTicketUsed)
var AppErrTicketNotEnough = NewAppError(CodeNotAllow, "", ErrTicketNotEnough)
var AppErrOrderFromat = NewAppError(CodeArgument, "", ErrOrderFormat)

// 微信模板消息
/***
TplMsgCreateTicket 创建订单的模板消息
*/
type TplMsgOrderPaid struct {
	BaseTemplateMessage
	OrderNO string
	Amount  float32
}

func (t *TplMsgOrderPaid) GetUrl() string {
	return fmt.Sprintf(ServerIndex, t.BaseTemplateMessage.ComID, "/pages/order/order?state=2")
}

func (t *TplMsgOrderPaid) GetMiniAppID() string {
	return ""
}

func (t *TplMsgOrderPaid) GetMsgType() string {
	return "order_paid"
}

func (t *TplMsgOrderPaid) GetItems() TemplateDataItemMap {
	items := TemplateDataItemMap{
		"first":    &TemplateDataItem{Value: t.First},
		"keyword1": &TemplateDataItem{Value: t.OrderNO},
		"keyword2": &TemplateDataItem{Value: fmt.Sprintf("%.2f", t.Amount)},
		"remark":   &TemplateDataItem{Value: t.Remark},
	}
	return items
}

/***
TplMsgCreateTicket 创建订单的模板消息
*/
type TplMsgCreateTicket struct {
	BaseTemplateMessage
	UserName   string
	OrderNO    string
	TicketName string
	Num        uint
	Amount     float32
	CheckCode  string
	CreatedAt  time.Time
}

func (t *TplMsgCreateTicket) GetUrl() string {
	return fmt.Sprintf(ServerIndex, t.BaseTemplateMessage.ComID, "/pages/ticket/index?code="+t.CheckCode)
}

func (t *TplMsgCreateTicket) GetMiniAppID() string {
	return ""
}

func (t *TplMsgCreateTicket) GetMsgType() string {
	return "create_ticket"
}

func (t *TplMsgCreateTicket) GetItems() TemplateDataItemMap {
	items := TemplateDataItemMap{
		"first":    &TemplateDataItem{Value: t.First},
		"keyword1": &TemplateDataItem{Value: t.TicketName},
		"keyword2": &TemplateDataItem{Value: strconv.Itoa(int(t.Num))},
		"keyword3": &TemplateDataItem{Value: fmt.Sprintf("%.2f", t.Amount)},
		"keyword4": &TemplateDataItem{Value: t.CreatedAt.String()},
		"keyword5": &TemplateDataItem{Value: t.OrderNO},
		"remark":   &TemplateDataItem{Value: t.Remark},
	}
	return items
}

type TplMsgCheckTicket struct {
}

type TplMsgRefund struct {
}
