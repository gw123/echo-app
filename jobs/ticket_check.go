package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

type TicketCheck struct {
	echoapp.Ticket
}

func (s *TicketCheck) GetName() string {
	return "ticket-check"
}

func (s *TicketCheck) RetryCount() int {
	return 5
}

func (s *TicketCheck) Delay() time.Duration {
	return 0
}
