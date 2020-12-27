package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

type TicketCreated struct {
	echoapp.Ticket
}

func (s *TicketCreated) GetName() string {
	return "ticket-created"
}

func (s *TicketCreated) RetryCount() int {
	return 5
}

func (s *TicketCreated) Delay() time.Duration {
	return 0
}
