package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

type OrderPaid struct {
	Order *echoapp.Order
}

func (s *OrderPaid) GetName() string {
	return "order-paid"
}

func (s *OrderPaid) RetryCount() int {
	return 2
}

func (s *OrderPaid) Delay() time.Duration {
	return 0
}
