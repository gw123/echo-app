package services

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	echoapp "github.com/gw123/echo-app"
	"github.com/pkg/errors"
	"strings"
)

const RegionId = "cn-hangzhou"

type SmsService struct {
	smsOptionMap map[string]echoapp.SmsOption
	clientMap    map[string]*dysmsapi.Client
}

func NewSmsService(options map[string]echoapp.SmsOption) *SmsService {
	return &SmsService{smsOptionMap: options}
}

func (mSvr SmsService) SendMessage(opt echoapp.SendMessageOptions) error {
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
		return errors.New(response.RequestId + "," + response.Message)
	}
	return nil
}
