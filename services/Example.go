package services

import "time"

type Example struct {
}

const TimeFormat = "2006-01-02 15:04:05"

func (e Example) GetTime() string {
	return time.Now().Format(TimeFormat)
}
