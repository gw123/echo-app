package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	tasksapp "github.com/gw123/echo-app"
	"github.com/gw123/glog"
	"github.com/pkg/errors"
)

const RedisSmsCode = "Mobile:%d:%s"
const RegionId = "cn-hangzhou"

type SmsService struct {
	redis  *redis.Client
	comSvr tasksapp.CompanyService
	db     *gorm.DB
}

func NewSmsService(db *gorm.DB, redis *redis.Client) *SmsService {
	return &SmsService{
		db:    db,
		redis: redis,
	}
}

func (sSvr *SmsService) GetChannel(comId uint, smsTpl string) (*echoapp.SmsChannel, error) {
	var channel echoapp.SmsChannel
	if err := sSvr.db.Where("com_id = ? and type = ?", comId, smsTpl).Find(&channel).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func (sSvr *SmsService) SendVerifyCodeSms(comId uint, phone string, code string) error {
	options := &tasksapp.SendMessageOptions{
		Type:          "loginCode",
		ComId:         comId,
		PhoneNumbers:  []string{phone},
		TemplateParam: fmt.Sprintf(`{"code":"%s"}`, code),
	}

	rkey := fmt.Sprintf(RedisSmsCode, comId, phone)
	duration := sSvr.redis.TTL(rkey).Val()
	if duration > time.Minute {
		return errors.New(fmt.Sprintf("请在%d秒后尝试", (duration-time.Minute)/time.Second))
	}
	if err := sSvr.SendMessage(options); err != nil {
		return errors.Wrap(err, "短信发生失败")
	}
	if err := sSvr.redis.Set(rkey, code, 2*time.Minute).Err(); err != nil {
		return errors.Wrap(err, "cache set")
	}
	return nil
}

func (sSvr *SmsService) CheckVerifyCode(comId uint, phone string, code string) bool {
	rCode := sSvr.redis.Get(fmt.Sprintf(RedisSmsCode, comId, phone)).Val()
	if code == "" && rCode != code {
		return false
	}
	return true
}

func (sSvr SmsService) SendMessage(opt *tasksapp.SendMessageOptions) error {
	smsChannel, err := sSvr.GetChannel(opt.ComId, opt.Type)
	if err != nil {
		return errors.Wrap(err, "SendMessage->GetChannel")
	}

	client, err := dysmsapi.NewClientWithAccessKey(RegionId, smsChannel.Key, smsChannel.Secret)
	if err != nil {
		glog.Errorf("newClient %s", err)
		return errors.Wrap(err, "SendMessage->NewClientWithAccessKey")
	}

	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = strings.Join(opt.PhoneNumbers, ",")
	request.SignName = smsChannel.SignName
	request.TemplateCode = smsChannel.TemplateCode
	request.TemplateParam = opt.TemplateParam
	response, err := client.SendSms(request)
	if err != nil {
		glog.Errorf("response %s", err)
		return errors.Wrap(err, "SendMessage->SendSms")
	}
	if response.Message != "OK" {
		return errors.New("请求失败:" + response.RequestId + "," + response.Message)
	}
	glog.Warnf("Message :%s ,code :%s", response.Message, response.Code)
	return nil
}
