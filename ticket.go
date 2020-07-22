package echoapp

import (
	"fmt"
	"time"
)

type Ticket struct {
	ID         int64     `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"-"`
	Code       string    `json:"code"`
	GoodsId    int       `json:"-"`
	OrderNo    string    `json:"-"`
	OrderId    uint      `json:"-"`
	Mobile     string    `json:"mobile"`
	Name       string    `json:"name"`
	Number     int       `json:"number"`
	Status     string    `json:"status"  gorm:"status" `
	UsedNumber int       `json:"used_number"`
	Username   string    `json:"username"`
	UsedAt     time.Time `json:"used_at"`
	Cover      string    `json:"cover"`
	ComId      int       `json:"com_id"`
	UserId     int       `json:"user_id"`
	Rand       int64     `json:"-"`
}

func (c *Ticket) AfterFind() error {
	if c.Rand != 0 {
		c.Code = fmt.Sprintf("%d%d", c.Rand, c.Rand+1234+c.ID)
	}
	return nil
}

type CodeTicket struct {
	BayAt       time.Time `json:"bay_at"`
	ComId       int       `json:"com_id"`
	GoodsCover  string    `json:"goods_cover"`
	GoodsId     uint      `json:"goods_id"`
	GoodsName   string    `json:"goods_name"`
	OrderNo     string    `json:"order_no"`
	OrderStatus string    `json:"order_status"`
	Username    string    `json:"username"`
	UserId      int       `json:"user_id"`
	Tickets     []*Ticket `json:"tickets"`
	XcxCover    string    `json:"xcx_cover"`
}