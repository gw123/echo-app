package services

import (
	"context"
	"fmt"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/oauth2"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/glog"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/iGoogle-ink/gotil"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	wx "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/js"
	"github.com/silenceper/wechat/v2/officialaccount/server"
	"net/http"
)

type WechatService struct {
	authRedirectUrl string
	comSvr          echoapp.CompanyService
	wx              *wx.Wechat
	jsHost          string
	reids           *redis.Client
}

func NewWechatService(comSvr echoapp.CompanyService, authUrl string, jsHost string, redis *redis.Client) *WechatService {
	wx := wx.NewWechat()
	return &WechatService{
		comSvr:          comSvr,
		authRedirectUrl: authUrl,
		jsHost:          jsHost,
		wx:              wx,
		reids:           redis,
	}
}

func (we *WechatService) GetJsConfig(comID uint, url string) (*js.Config, error) {
	cfg, err := we.GetComOfficialCfg(comID)
	if err != nil {
		return nil, err
	}
	wxOfficial := we.wx.GetOfficialAccount(cfg)
	jsConfig := wxOfficial.GetJs()

	glog.Infof("GetJsConfig Url: %s", url)
	config, err := jsConfig.GetConfig(url)
	if err != nil {
		return nil, errors.Wrap(err, "GetJSConfig")
	}
	return config, nil
}

func (we *WechatService) GetOfficialServer(ctx echo.Context, comID uint) (*server.Server, error) {
	cfg, err := we.GetComOfficialCfg(comID)
	if err != nil {
		return nil, err
	}
	wxOfficial := we.wx.GetOfficialAccount(cfg)
	server := wxOfficial.GetServer(ctx.Request(), ctx.Response())
	////设置接收消息的处理方法
	//server.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
	//	//TODO
	//	//回复消息：演示回复用户发送的消息
	//	text := message.NewText(msg.Content)
	//	return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	//})
	//if err = server.Serve(); err != nil {
	//	return nil, err
	//}
	//server.Send()
	return server, nil
}

func (we *WechatService) GetComOfficialCfg(comID uint) (*offConfig.Config, error) {
	com, err := we.comSvr.GetCompanyById(comID)
	if err != nil {
		return nil, errors.Wrap(err, "GetCompanyById")
	}
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:          com.WxOfficialAppId,
		AppSecret:      com.WxOfficialSecret,
		Token:          com.WxToken,
		EncodingAESKey: com.WxOfficialAesKey,
		Cache:          memory,
	}
	return cfg, nil
}

func (we *WechatService) GetAuthCodeUrl(comId uint) (url string, err error) {
	com, err := we.comSvr.GetCachedCompanyById(comId)
	if err != nil {
		return "", errors.Wrap(err, "WxLogin 获取com失败")
	}
	url = fmt.Sprintf("%s/%d/wxAuthCallBack", we.authRedirectUrl, comId)
	url = mpoauth2.AuthCodeURL(com.WxOfficialAppId, url, "snsapi_userinfo", "")
	return
}

func (we *WechatService) GetEndPoint(comId uint) (*mpoauth2.Endpoint, error) {
	com, err := we.comSvr.GetCachedCompanyById(comId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetEndPoint 获取com失败：%d", comId)
	}
	oauth2Endpoint := mpoauth2.NewEndpoint(com.WxOfficialAppId, com.WxOfficialSecret)
	return oauth2Endpoint, nil
}

func (we *WechatService) GetUserInfo(ctx context.Context, comId uint, code string) (*mpoauth2.UserInfo, error) {
	endPoint, err := we.GetEndPoint(comId)
	if err != nil {
		return nil, errors.Wrap(err, "GetEndPoint")
	}

	oauth2Client := oauth2.Client{
		Endpoint: endPoint,
	}
	token, err := oauth2Client.ExchangeToken(code)
	if err != nil {
		return nil, errors.Wrap(err, "ExchangeToken")
	}
	glog.Infof("token: %+v\r\n", token)

	userinfo, err := mpoauth2.GetUserInfo(token.AccessToken, token.OpenId, "", nil)
	if err != nil {
		//echoapp_util.ExtractEntry(ctx).WithError(err).Error("code为空")
		return nil, errors.Wrap(err, "ExchangeToken")
	}
	return userinfo, nil
}

/****/
func (we *WechatService) UnifiedOrder(order *echoapp.Order, openId string) (*wechat.UnifiedOrderResponse, error) {
	com, err := we.comSvr.GetCachedCompanyById(order.ComId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetEndPoint 获取com失败：%d", order.ComId)
	}

	var appID string
	if order.ClientType == "wx_official" {
		appID = com.WxOfficialAppId
	} else {
		appID = com.WxMiniAppId
	}
	glog.Infof("clientType: %s , appId: %s ,openId: %s", order.ClientType, appID, openId)

	client := wechat.NewClient(appID, com.WxPaymentMchId, com.WxPaymentKey, false)
	client.SetCountry(wechat.China)
	glog.Infof("orderNo : %s, %s", order.OrderNo, order.ClientIP)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gotil.GetRandomString(32))
	bm.Set("body", "测试支付003")
	bm.Set("out_trade_no", order.OrderNo)
	bm.Set("total_fee", 100)
	bm.Set("spbill_create_ip", order.ClientIP)
	bm.Set("notify_url", "http://www.gopay.ink")
	bm.Set("trade_type", wechat.TradeType_H5)
	bm.Set("device_info", "WEB")
	bm.Set("sign_type", wechat.SignType_MD5)
	bm.Set("openid", openId)

	// 嵌套json格式数据（例如：H5支付的 scene_info 参数）
	h5Info := make(map[string]string)
	h5Info["type"] = "Wap"
	h5Info["wap_url"] = "http://www.gopay.ink"
	h5Info["wap_name"] = "H5测试支付"

	sceneInfo := make(map[string]map[string]string)
	sceneInfo["h5_info"] = h5Info

	bm.Set("scene_info", sceneInfo)

	// 参数 sign ，可单独生成赋值到BodyMap中；也可不传sign参数，client内部会自动获取
	// 如需单独赋值 sign 参数，需通过下面方法，最后获取sign值并在最后赋值此参数
	glog.Infof("com : %+v", com)
	glog.Infof("bm %+v", bm)
	//sign := wechat.GetParamSign(appID, com.WxPaymentMchId, com.WxPaymentKey, bm)
	//todo 目前是测试阶段
	//sign, _ := wechat.GetSanBoxParamSign(appID, com.WxPaymentMchId, com.WxPaymentKey, bm)
	//bm.Set("sign", sign)
	resp, err := client.UnifiedOrder(bm)
	if err != nil {
		return nil, errors.Wrap(err, "UnifiedOrder")
	}
	glog.Infof("Resp: ")
	spew.Dump(resp)
	//todo 线上要补上签名校验
	//ok, err := wechat.VerifySign(com.WxPaymentKey, wechat.SignType_MD5, resp)
	//if err != nil {
	//	return nil, errors.Wrap(err, "UnifiedOrder")
	//}
	//if !ok {
	//	return nil, errors.New(resp.ReturnCode + ":" + resp.ReturnMsg)
	//}
	return resp, nil
}

/***/
func (we *WechatService) QueryOrder(order *echoapp.Order) (string, error) {
	com, err := we.comSvr.GetCachedCompanyById(order.ComId)
	if err != nil {
		return "", errors.Wrapf(err, "GetEndPoint 获取com失败：%d", order.ComId)
	}

	var appID string
	if order.ClientType == "official" {
		appID = com.WxOfficialAppId
	} else {
		appID = com.WxMiniAppId
	}
	glog.Infof("clientType: %s , appId: %s ", order.ClientType, appID)

	client := wechat.NewClient(appID, com.WxPaymentMchId, com.WxPaymentKey, false)
	client.SetCountry(wechat.China)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gotil.GetRandomString(32))
	bm.Set("transaction_id", order.TransactionId)
	bm.Set("out_trade_no", order.OrderNo)
	//todo 查询订单接口
	resp, bmResp, err := client.QueryOrder(bm)

	glog.Infof("%+v", bmResp)
	glog.Infof("%+v", resp)
	if err != nil {
		return echoapp.OrderStatusUnpay, errors.Wrap(err, "queryOrder")
	}

	if resp.ResultCode == "SUCCESS" {
		return echoapp.OrderStatusPaid, nil
	}

	return echoapp.OrderStatusUnpay, nil
}

//查询支付结果
func (we *WechatService) PayCallback(r *http.Request) (*wechat.NotifyRequest, error) {
	resp, err := wechat.ParseNotify(r)
	if err != nil {
		return nil, err
	}
	if resp.ReturnCode == "SUCCESS" && resp.ResultCode == "SUCCESS" {
		return resp, nil
	} else {
		return nil, errors.Errorf("Code:%s, %s", resp.ErrCode, resp.ErrCodeDes)
	}
}

func (we *WechatService) Refund(order *echoapp.Order, openId string) (*wechat.RefundResponse, error) {
	com, err := we.comSvr.GetCachedCompanyById(order.ComId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetEndPoint 获取com失败：%d", order.ComId)
	}

	var appID string
	if order.ClientType == "wx_official" {
		appID = com.WxOfficialAppId
	} else {
		appID = com.WxMiniAppId
	}
	glog.Infof("clientType: %s , appId: %s ,openId: %s", order.ClientType, appID, openId)

	client := wechat.NewClient(appID, com.WxPaymentMchId, com.WxPaymentKey, false)
	client.SetCountry(wechat.China)
	glog.Infof("orderNo : %s, %s", order.OrderNo, order.ClientIP)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gotil.GetRandomString(32))
	bm.Set("body", "测试支付003")
	bm.Set("out_trade_no", order.OrderNo)
	bm.Set("total_fee", order.RealTotal)
	bm.Set("spbill_create_ip", order.ClientIP)
	bm.Set("notify_url", "http://www.gopay.ink")
	bm.Set("trade_type", wechat.TradeType_H5)
	bm.Set("device_info", "WEB")
	bm.Set("sign_type", wechat.SignType_MD5)
	bm.Set("openid", openId)
	bm.Set("transaction_id", order.TransactionId)
	bm.Set("out_refund_no", order.OrderNo)
	bm.Set("refund_fee", order.RealTotal)

	glog.Infof("com : %+v", com)
	glog.Infof("bm %+v", bm)
	//sign := wechat.GetParamSign(appID, com.WxPaymentMchId, com.WxPaymentKey, bm)
	//todo 目前是测试阶段
	//sign, _ := wechat.GetSanBoxParamSign(appID, com.WxPaymentMchId, com.WxPaymentKey, bm)
	//bm.Set("sign", sign)
	resp, _, err := client.Refund(bm, "", "", "")
	if err != nil {
		return nil, errors.Wrap(err, "UnifiedOrder")
	}
	//todo 线上要补上签名校验
	ok, err := wechat.VerifySign(com.WxPaymentKey, wechat.SignType_MD5, resp)
	if err != nil {
		return nil, errors.Wrap(err, "UnifiedOrder")
	}
	if !ok {
		return nil, errors.New(resp.ReturnCode + ":" + resp.ReturnMsg)
	}
	return resp, nil
}

func (we *WechatService) RefundCallback(r *http.Request) (*wechat.RefundNotify, error) {
	resp, err := wechat.ParseRefundNotify(r)
	if err != nil {
		return nil, errors.Wrap(err, "ParseRefundNotify")
	}
	return wechat.DecryptRefundNotifyReqInfo(resp.ReqInfo, "")
}

func (we *WechatService) QueryRefund(order *echoapp.Order) (string, error) {
	com, err := we.comSvr.GetCachedCompanyById(order.ComId)
	if err != nil {
		return "", errors.Wrapf(err, "GetEndPoint 获取com失败：%d", order.ComId)
	}

	var appID string
	if order.ClientType == "official" {
		appID = com.WxOfficialAppId
	} else {
		appID = com.WxMiniAppId
	}
	glog.Infof("clientType: %s , appId: %s ", order.ClientType, appID)

	client := wechat.NewClient(appID, com.WxPaymentMchId, com.WxPaymentKey, false)
	client.SetCountry(wechat.China)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gotil.GetRandomString(32))
	bm.Set("transaction_id", order.TransactionId)
	bm.Set("out_trade_no", order.OrderNo)
	//todo 查询订单接口 校验返回参数
	resp, bmResp, err := client.QueryRefund(bm)
	glog.Infof("%+v", bmResp)
	glog.Infof("%+v", resp)

	if err != nil {
		return echoapp.OrderStatusUnpay, errors.Wrap(err, "queryOrder")
	}

	if resp.ResultCode == "SUCCESS" {
		return echoapp.OrderStatusRefund, nil
	}

	return echoapp.OrderStatusUnpay, nil
}
