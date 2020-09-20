package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

type OrderCreate struct {
	*echoapp.Order
}

func (s *OrderCreate) Trace() []string {
	return nil
}

func (s *OrderCreate) GetName() string {
	return "order-create"
}

func (s *OrderCreate) RetryCount() int {
	return 8
}

func (s *OrderCreate) Delay() time.Duration {
	return 5 * time.Second
}
