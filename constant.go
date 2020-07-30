package echoapp

import "github.com/pkg/errors"

const (
	TimeFormat   = "2006-01-02 15:04:05"
	DateFormat   = "2006-01-02"
	HostURL      = "host_url"
	RequestRoot  = "static_root"
	AssetVersion = "asset_version"

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
	CodeNoAuth     = 401
	CodeDBError    = 402
	CodeCacheError = 403
	CodeArgument   = 404
	CodeNotAllow   = 405
	CodeEtcdError  = 406
	CodeInnerError = 501
)

var ErrNotFoundCache = errors.New("not found cache item")
var ErrNotFoundDb = errors.New("not found db item")
var ErrDb = errors.New("db exec err")
var ErrNotFoundEtcd = errors.New("not found etcd item")
var ErrArgument = errors.New("argument error")
var ErrNotLogin = errors.New("请登录后尝试")
var ErrNotAuth = errors.New("not auth")
var ErrNotAllow = errors.New("not allow")
var ErrRefund = errors.New("订单已经退款")
var ErrTicketInvaild = errors.New("无效的门票")
var ErrTicketOverdue = errors.New("门票已经过期")
var ErrTicketUsed = errors.New("门票已经被使用")
var ErrTicketNotEnough = errors.New("已购门票数量不足,请核对数量")
