package entity

import "time"

const (
	staHour = 1
	//comId =14
)

type Statistics struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Date      string `json:"date"`
	TargetId  int    `json:"target_id"`
	Total     int64  `json:"total"`
	Type      string `json:"type"`
}

func (*Statistics) TableName() string {
	return "Statistics"
}
