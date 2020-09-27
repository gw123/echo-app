package jobs

import (
	"time"
)

type UserScoreChange struct {
	ComId        uint
	UserId       uint
	Score        int
	Source       string
	SourceDetail string
	Note         string
}

func (s *UserScoreChange) GetName() string {
	return "user-score-change"
}

func (s *UserScoreChange) RetryCount() int {
	return 2
}

func (s *UserScoreChange) Delay() time.Duration {
	return 0
}
