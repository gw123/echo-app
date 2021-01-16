package sms_tpls

import "github.com/gw123/gworker"

type SmsTplAPiIpl struct {
	pusher gworker.Producer
}

func NewSmsTplApi(pusher gworker.Producer) *SmsTplAPiIpl {
	return &SmsTplAPiIpl{
		pusher: pusher,
	}
}
