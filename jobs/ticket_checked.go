package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

type TicketChecked struct {
	echoapp.Ticket
}

func (s *TicketChecked) Trace() []string {
	return nil
}

func (s *TicketChecked) GetName() string {
	return "ticket-checked"
}

func (s *TicketChecked) RetryCount() int {
	return 5
}

func (s *TicketChecked) Delay() time.Duration {
	return 0
}
