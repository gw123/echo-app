package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/silenceper/wechat/v2/officialaccount/message"

	"github.com/silenceper/wechat/v2/officialaccount"

	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/oauth2"
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
)

type WechatService struct {
	authRedirectUrl  string
	comSvr           echoapp.CompanyService
	wx               *wx.Wechat
	jsHost           string
	reids            *redis.Client
	officialAccounts map[uint]*officialaccount.OfficialAccount
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
	if err != nil || com.WxOfficialAppId == "" {
		return "", errors.Wrapf(err, "WxLogin 获取com.WxOfficialAppId失败: %d", comId)
	}
	url = fmt.Sprintf("%s/%d/wxAuthCallBack", we.authRedirectUrl, comId)
	url = mpoauth2.AuthCodeURL(com.WxOfficialAppId, url, "snsapi_userinfo", "")
	return
}

func (we *WechatService) GetEndPoint(comId uint) (*mpoauth2.Endpoint, error) {
	com, err := we.comSvr.GetCachedCompanyById(comId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetEndPoint 获取com.WxOfficialAppId失败：%d", comId)
	}
	oauth2Endpoint := mpoauth2.NewEndpoint(com.WxOfficialAppId, com.WxOfficialSecret)
	return oauth2Endpoint, nil
}

func (we *WechatService) GetClient(comId uint) (*mpoauth2.Endpoint, error) {
	com, err := we.comSvr.GetCachedCompanyById(comId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetEndPoint 获取com.WxOfficialAppId失败：%d", comId)
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

func (we *WechatService) GetAccessToken(comID uint) (string, error) {
	officialaccount, err := we.GetOfficialByComID(comID)
	if err != nil {
		return "", err
	}
	return officialaccount.GetAccessToken()
}

func (we *WechatService) GetOfficialByComID(comID uint) (*officialaccount.OfficialAccount, error) {
	officialAccount, ok := we.officialAccounts[comID]
	if !ok {
		cfg, err := we.GetComOfficialCfg(comID)
		if err != nil {
			return nil, err
		}
		officialAccount = we.wx.GetOfficialAccount(cfg)
	}
	return officialAccount, nil
}

// SendTplMessage 发送模板消息跳转到公众号网页
func (we *WechatService) SendTplMessageRaw(ctx context.Context, comID uint, openid, templateID, URL string, items echoapp.TemplateDataItemMap) (int64, error) {
	account, err := we.GetOfficialByComID(comID)
	if err != nil {
		return 0, err
	}

	data := make(map[string]*message.TemplateDataItem)
	for key, item := range items {
		data[key] = &message.TemplateDataItem{
			Value: item.Value,
			Color: item.Color,
		}
	}

	return account.GetTemplate().Send(&message.TemplateMessage{
		ToUser:     openid,
		TemplateID: templateID,
		URL:        URL,
		Color:      "",
		Data:       data,
	})
}

// SendTplMessage 发送模板消息跳转到公众号网页
func (we *WechatService) SendTplMessage(ctx context.Context, msg echoapp.TemplateMessage) (int64, error) {
	com, err := we.comSvr.GetCompanyById(msg.GetComID())
	if err != nil {
		return 0, errors.Wrap(err, "GetCompanyById")
	}

	wxTemplateType := com.GetTemplateType(msg.GetMsgType())
	if wxTemplateType == nil {
		return 0, errors.New("unset template type : " + msg.GetMsgType())
	}

	data := make(map[string]*message.TemplateDataItem)
	for key, val := range msg.GetItems() {
		data[key] = &message.TemplateDataItem{
			Value: val.Value,
		}

		if color := wxTemplateType.GetKeywordColor(key); color != "" {
			data[key].Color = color
		}
	}

	account, err := we.GetOfficialByComID(msg.GetComID())
	if err != nil {
		return 0, err
	}

	return account.GetTemplate().Send(&message.TemplateMessage{
		TemplateID: wxTemplateType.TemplateID,
		ToUser:     msg.GetOpenid(),
		URL:        msg.GetUrl(),
		Data:       data,
		MiniProgram: struct {
			AppID    string `json:"appid"`
			PagePath string `json:"pagepath"`
		}{
			AppID:    msg.GetMiniAppID(),
			PagePath: msg.GetUrl(),
		},
	})
}

func (we *WechatService) UnifiedOrderByPHP(order *echoapp.Order, openId string) (*echoapp.WxPreOrderResponse, error) {
	params := make(map[string]string)
	params["openid"] = openId
	params["comID"] = strconv.Itoa(int(order.ComId))
	params["order_no"] = order.OrderNo
	params["total"] = strconv.Itoa(int(order.RealTotal))
	url := "http://shop.laravelschool.xyt/jspay?com_id=14&openid=" + openId + "&order_no=" + order.OrderNo + "&total=1"
	req, err := we.getRequest("GET", url, params)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	type WxPreOrderResponse struct {
		Code int                         `json:"code"`
		Data *echoapp.WxPreOrderResponse `json:"data"`
	}
	res := WxPreOrderResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return res.Data, err
}

func (we *WechatService) getRequest(method string, u string, data map[string]string) (*http.Request, error) {
	//post要提交的数据
	dataUrlVal := url.Values{}
	for key, val := range data {
		dataUrlVal.Add(key, val)
	}
	req, err := http.NewRequest(method, u, strings.NewReader(dataUrlVal.Encode()))
	if err != nil {
		return nil, err
	}
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

/****/
func (we *WechatService) UnifiedOrder(order *echoapp.Order, openId string) (*echoapp.WxPreOrderResponse, error) {
	com, err := we.comSvr.GetCachedCompanyById(order.ComId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetEndPoint 获取com失败：%d", order.ComId)
	}

	var appID string
	if order.ClientType == echoapp.ClientWxOfficial {
		appID = com.WxOfficialAppId
	} else {
		appID = com.WxMiniAppId
	}

	client := wechat.NewClient(appID, com.WxPaymentMchId, com.WxPaymentKey, true)
	client.SetCountry(wechat.China)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gotil.GetRandomString(32))
	bm.Set("body", order.GoodsList[0].Name)
	bm.Set("out_trade_no", order.OrderNo)
	bm.Set("total_fee", int(order.RealTotal))
	bm.Set("spbill_create_ip", order.ClientIP)
	bm.Set("notify_url", "http://m.xytschool.com/wx_order")
	//
	bm.Set("trade_type", wechat.TradeType_JsApi)
	bm.Set("sign_type", wechat.SignType_MD5)
	bm.Set("openid", openId)

	resp, err := client.UnifiedOrder(bm)
	if err != nil {
		return nil, errors.Wrap(err, "UnifiedOrder")
	}

	ok, err := wechat.VerifySign(com.WxPaymentKey, wechat.SignType_MD5, resp)
	if err != nil {
		return nil, errors.Wrap(err, "UnifiedOrder")
	}
	if !ok {
		glog.Infof("weixin unified order 签名验证失败")
		return nil, errors.New(resp.ReturnCode + ":" + resp.ReturnMsg)
	}

	// 微信统一下单后，获取微信小程序支付、APP支付、微信内H5支付所需要的 paySign
	//(这里改sdk 不太友好 需要自己对返回结果再次签名算出paySign，前面接口返回的sign 是请求的签名和paySign不是一回事)
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	packages := "prepay_id=" + resp.PrepayId // 此处的 wxRsp.PrepayId ,统一下单成功后得到
	paySign := wechat.GetMiniPaySign(com.WxOfficialAppId, resp.NonceStr, packages, wechat.SignType_MD5, timeStamp, com.WxPaymentKey)
	preResp := &echoapp.WxPreOrderResponse{
		AppID:     resp.Appid,
		NonceStr:  resp.NonceStr,
		Package:   packages,
		SignType:  wechat.SignType_MD5,
		PaySign:   paySign,
		Timestamp: timeStamp,
	}
	return preResp, nil
}

/***/
func (we *WechatService) QueryOrder(order *echoapp.Order) (string, error) {
	com, err := we.comSvr.GetCachedCompanyById(order.ComId)
	if err != nil {
		return "", errors.Wrapf(err, "GetEndPoint 获取com失败：%d", order.ComId)
	}

	var appID string
	if order.ClientType == echoapp.ClientWxOfficial {
		appID = com.WxOfficialAppId
	} else {
		appID = com.WxMiniAppId
	}
	//glog.Infof("clientType: %s , appId: %s ", order.ClientType, appID)

	client := wechat.NewClient(appID, com.WxPaymentMchId, com.WxPaymentKey, true)
	client.SetCountry(wechat.China)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gotil.GetRandomString(32))
	bm.Set("transaction_id", order.TransactionId)
	bm.Set("out_trade_no", order.OrderNo)
	//todo 查询订单接口
	resp, _, err := client.QueryOrder(bm)
	if err != nil {
		return echoapp.OrderStatusUnpay, errors.Wrap(err, "queryOrder")
	}

	glog.Info(resp.ReturnCode + " -- " + resp.TradeState + " -- " + resp.TotalFee + " -- " + strconv.Itoa(int(order.RealTotal)))
	if resp.ResultCode == "SUCCESS" {
		if resp.TradeState == "SUCCESS" {
			if resp.TotalFee == strconv.Itoa(int(order.RealTotal)) {
				return echoapp.OrderPayStatusPaid, nil
			} else {
				glog.Errorf("查询成功但系统订单金额和微信订单金额不一致:orderNo %s,wxRes %s ,local %s", order.OrderNo, resp.TotalFee, strconv.Itoa(int(order.RealTotal)))
				return echoapp.OrderPayStatusUnpay, errors.Errorf("查询成功但系统订单金额和微信订单金额不一致:orderNo %s,wxRes %s ,local %s", order.OrderNo, resp.TotalFee, strconv.Itoa(int(order.RealTotal)))
			}
		} else if resp.TradeState == "REFUND" {
			return echoapp.OrderPayStatusRefund, nil
		} else {
			return echoapp.OrderPayStatusUnpay, nil
		}
	}
	return echoapp.OrderPayStatusUnpay, errors.New(resp.ReturnMsg)
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
