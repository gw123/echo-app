package sms_tpls

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/gworker"
)

type SmsTplAPiIpl struct {
	pusher gworker.Producer
	smsSvr echoapp.SmsService
}

func NewSmsTplApi(pusher gworker.Producer, smsSvr echoapp.SmsService) *SmsTplAPiIpl {
	return &SmsTplAPiIpl{
		pusher: pusher,
		smsSvr: smsSvr,
	}
}
