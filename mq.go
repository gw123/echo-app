package echoapp

//
type BaseMqMsg struct {
	ComId string `json:"com_id"`
}

type MsgQueue interface {
	PushUpdateCahceJob(job UpdateCacheJob) error
	PushSmsJob(job SendMessageJob) error
}

type MqPool interface {
	AddMq(name string, mq MsgQueue)
	RemoveMq(name string)
	MQ(name string) (MsgQueue, error)
}
