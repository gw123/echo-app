package jobs

import "time"

type SendWxTplMsg struct {
	ComId        uint
	UserId       uint
	Openid       string
	Source       string
	SourceDetail string
	Msg          string
}

func (s *SendWxTplMsg) GetName() string {
	return "send-wx-tpl-msg"
}

func (s *SendWxTplMsg) RetryCount() int {
	return 2
}

func (s *SendWxTplMsg) Delay() time.Duration {
	return 0
}
