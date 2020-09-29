package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

type SendSms struct {
	echoapp.SendMessageOptions
}

func (s *SendSms) GetName() string {
	return "send-sms"
}

func (s *SendSms) RetryCount() int {
	return 2
}

func (s *SendSms) Delay() time.Duration {
	return 0
}
