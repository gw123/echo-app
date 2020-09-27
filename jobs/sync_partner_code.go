package jobs

import (
	"time"

	echoapp "github.com/gw123/echo-app"
)

type SyncPartnerCode struct {
	echoapp.Ticket
}

func (s *SyncPartnerCode) GetName() string {
	return "sync-partner-code"
}

func (s *SyncPartnerCode) RetryCount() int {
	return 5
}

func (s *SyncPartnerCode) Delay() time.Duration {
	return 0
}
