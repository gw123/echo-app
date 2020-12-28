package echoapp

import (
	"fmt"
	"github.com/gw123/glog"
	"time"
)

const (
	IdHashSalt = 1234567

	TicketStatusNormal  = "normal"
	TicketStatusUsed    = "used"
	TicketstatusRefund  = "refund"
	TicketstatusOverdue = "overdue"
	TicketstatusInvalid = "invalid"
)

type Ticket struct {
	ID         int64      `gorm:"primary_key" json:"id"`
	Cover      string     `json:"cover" gorm:"-"`
	ComId      uint       `json:"com_id"`
	UserId     uint       `json:"user_id"`
	StaffID    uint       `json:"staff_id"`
	Rand       int64      `json:"-"`
	Code       string     `json:"code" gorm:"-"` //前端显示用 这个参数通过rand 计算出来不需要入库
	GoodsId    uint       `json:"-"`
	OrderNo    string     `json:"-"`
	OrderId    uint       `json:"-"`
	Mobile     string     `json:"mobile"`
	Name       string     `json:"name" gorm:"-"`
	Source     string     `json:"name" gorm:"source"`
	Number     uint       `json:"number"`
	Status     string     `json:"status"  gorm:"status" `
	UsedNumber uint       `json:"used_number"`
	Username   string     `json:"username"`
	AddressID  uint       `gorm:"address_id" json:"address_id"`
	UsedAt     *time.Time `json:"used_at"`
	OverdueAt  *time.Time `json:"overdue_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"-"`
}

func (t *Ticket) AfterFind() error {
	if t.Rand != 0 {
		t.Code = fmt.Sprintf("%d%d", t.Rand, t.Rand+IdHashSalt+t.ID)
	}
	return nil
}

func (t *Ticket) UpdateStatus() {
	if t.OverdueAt != nil {
		if t.OverdueAt.Sub(time.Now()) < 0 {
			t.Status = TicketstatusOverdue
		}
	}
}

func (t *Ticket) IsValid() error {
	switch t.Status {
	case TicketstatusRefund:
		return ErrRefund
	case TicketstatusInvalid:
		return ErrTicketInvaild
	case TicketstatusOverdue:
		return ErrTicketOverdue
	case TicketStatusUsed:
		return ErrTicketUsed
	case TicketStatusNormal:
	}

	if t.UsedAt != nil {
		glog.Errorf("ticketId:%d,状态和used_at冲突", t.ID)
		return ErrTicketUsed
	}

	if t.OverdueAt != nil {
		if t.OverdueAt.Sub(time.Now()) < 0 {
			return ErrTicketOverdue
		}
	}
	return nil
}

type CodeTicket struct {
	BayAt       time.Time `json:"bay_at"`
	ComId       uint      `json:"com_id"`
	GoodsCover  string    `json:"goods_cover"`
	GoodsId     uint      `json:"goods_id"`
	GoodsName   string    `json:"goods_name"`
	OrderNo     string    `json:"order_no"`
	OrderStatus string    `json:"order_status"`
	Username    string    `json:"username"`
	UserId      uint      `json:"user_id"`
	Tickets     []*Ticket `json:"tickets"`
	XcxCover    string    `json:"xcx_cover"`
}

type TicketService interface {
	GetTicketByCode(code string) (*Ticket, error)
	PreCreateTicket(order *Order, source string, address *Address, goods *CartGoodsItem) *Ticket
	GetTicketsByOrder(order *Order) ([]*Ticket, error)
	//验票的相关逻辑,因为和事务有关系交给调用着去更新
	CheckTicket(ticket *Ticket, num uint, staffID uint) error
}
