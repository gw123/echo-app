package echoapp

//
type BaseMqMsg struct {
	MsgType   string `json:"msg_type"`
	CreatedAt string `json:"created_at"`
	ExprAt    string `json:"expr_at"`
	TryTimes  int    `json:"try_times"`
	MessageId string `json:"message_id"`
	TraceId   string `json:"trace_id"`
	Sender    string `json:"sender"`
	//Payload   string
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
