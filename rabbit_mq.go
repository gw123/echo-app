package echoapp

import (
	"github.com/streadway/amqp"
)

type RabbitMqPool interface {
	MQ(name string) (*amqp.Connection, error)
}
