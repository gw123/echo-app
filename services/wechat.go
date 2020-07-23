package services

import (
	"context"
	"fmt"
	mpoauth2 "github.com/chanxuehong/wechat/mp/oauth2"
	"github.com/chanxuehong/wechat/oauth2"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/glog"
	"github.com/iGoogle-ink/gopay"
	"github.com/iGoogle-ink/gopay/wechat"
	"github.com/iGoogle-ink/gotil"
	"github.com/pkg/errors"
)

type WechatService struct {
	authRedirectUrl string
	comSvr          echoapp.CompanyService
}

func NewWechatService(comSvr echoapp.CompanyService, authUrl string) *WechatService {
	return &WechatService{
		authRedirectUrl: authUrl,
		comSvr:          comSvr,
	}
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

/***

 */
func (we *WechatService) UnifiedOrder(order *echoapp.Order, openId string) (interface{}, error) {
	com, err := we.comSvr.GetCachedCompanyById(order.ComId)
	if err != nil {
		return nil, errors.Wrapf(err, "GetEndPoint 获取com失败：%d", order.ComId)
	}

	var appID string
	if order.ClientType == "official" {
		appID = com.WxOfficialAppId
	} else {
		appID = com.WxMiniAppId
	}

	client := wechat.NewClient(appID, com.WxPaymentMchId, com.WxPaymentKey, false)
	client.SetCountry(wechat.China)

	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", gotil.GetRandomString(32))
	bm.Set("body", "小程序测试支付")
	bm.Set("out_trade_no", order.OrderNo)
	bm.Set("total_fee", 1)
	bm.Set("spbill_create_ip", order.ClientIP)
	bm.Set("notify_url", "http://www.gopay.ink")
	bm.Set("trade_type", wechat.TradeType_Mini)
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
	sign := wechat.GetParamSign(appID, com.WxPaymentMchId, com.WxPaymentKey, bm)
	// sign, _ := wechat.GetSanBoxParamSign("wxdaa2ab9ef87b5497", mchId, apiKey, body)
	bm.Set("sign", sign)
	resp, err := client.UnifiedOrder(bm)
	if err != nil {
		return nil, errors.Wrap(err, "UnifiedOrder")
	}
	ok, err := wechat.VerifySign(com.WxPaymentKey, wechat.SignType_MD5, resp)
	if err != nil {
		return nil, errors.Wrap(err, "UnifiedOrder")
	}
	if !ok {
		return nil, errors.New("下单失败")
	}
	return resp, nil
}
