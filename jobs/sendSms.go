package jobs

import "time"

type SendSms struct {
	ComId        uint
	UserId       uint
	Score        int
	Source       string
	SourceDetail string
	Note         string
}

func (s *SendSms) Trace() []string {
	return nil
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
