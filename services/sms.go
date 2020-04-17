package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"strings"
	"sync"
)

const RegionId = "cn-hangzhou"

type SmsService struct {
	smsOptionMap map[string]echoapp.SmsOption
	clientMap    map[string]*dysmsapi.Client
	mu           sync.Mutex
}

func NewSmsService(options map[string]echoapp.SmsOption) *SmsService {
	return &SmsService{
		smsOptionMap: options,
		clientMap: map[string]*dysmsapi.Client{},
	}
}

func (mSvr SmsService) SendMessage(ctx echo.Context, opt echoapp.SendMessageOptions) error {
	echoapp_util.ExtractEntry(ctx).Info("发送短信,请求内容:", opt)
	client, ok := mSvr.clientMap[opt.Token]
	if !ok || client == nil {
		smsOption, ok := mSvr.smsOptionMap[opt.Token]
		if !ok {
			return errors.New("notfound sms token")
		}

		var err error
		client, err = dysmsapi.NewClientWithAccessKey(RegionId, smsOption.AccessKey, smsOption.AccessSecret)
		if err != nil {
			return errors.Wrap(err, "SendMessage->NewClientWithAccessKey")
		}
		//防止多线程并发操作
		mSvr.mu.Lock()
		defer mSvr.mu.Unlock()
		mSvr.clientMap[opt.Token] = client
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = strings.Join(opt.PhoneNumbers, ",")
	request.SignName = opt.SignName
	request.TemplateCode = opt.TemplateCode
	request.TemplateParam = opt.TemplateParam
	response, err := client.SendSms(request)
	if err != nil {
		return errors.Wrap(err, "SendMessage->SendSms")
	}
	if response.Message != "OK" {
		return errors.New("请求失败:" + response.RequestId + "," + response.Message)
	}
	return nil
}
