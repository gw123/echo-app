package echoapp

import (
	"github.com/streadway/amqp"
)

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

type RabbitMqPool interface {
	MQ(name string) (*amqp.Connection, error)
}
