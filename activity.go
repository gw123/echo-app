package echoapp

import "time"

type Activity struct {
	Id int `json:"id"`
}

func (a *Activity) TableName() string {
	return "activies"
}

type Notify struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notify *Notify) TableName() string {
	return "activies"
}

type ActivityService interface {
	GetNotifyList(comId int, lastId, limit int) ([]*Notify, error)
	GetNotifyDetail(id int) (*Notify, error)
}
