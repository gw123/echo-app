package services

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type TicketService struct {
	mu sync.Mutex
	db *gorm.DB
}

func NewTicketService(db *gorm.DB) *TicketService {
	return &TicketService{db: db}
}

func (tkSvr *TicketService) DeTicketCode(code string) (*echoapp.Ticket, error) {
	rand, _ := strconv.ParseInt(code[0:8], 10, 64)
	temp, _ := strconv.ParseInt(code[8:], 10, 64)
	if rand == 0 || temp == 0 {
		return nil, errors.New("code is not vaild")
	}
	tId := temp - echoapp.IdHashSalt - rand
	ticket := &echoapp.Ticket{}
	if err := tkSvr.db.Where("id = ?", tId).First(ticket).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	if ticket.Rand != rand {
		return nil, errors.New("校验失败")
	}

	return ticket, nil
}

func (tkSvr *TicketService) GetTicketByCode(code string) (*echoapp.Ticket, error) {
	ticket, err := tkSvr.DeTicketCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "DeTicketCode")
	}
	return ticket, nil
}

// 拼装门票数据,为了使用事物 保存ticket应该放到 saveOrder那部分
func (tkSvr *TicketService) PreCreateTicket(order *echoapp.Order, source string, address *echoapp.Address, goods *echoapp.CartGoodsItem) *echoapp.Ticket {
	r := rand.Int31n(89999999) + 10000000
	ticket := &echoapp.Ticket{
		GoodsId:    goods.GoodsId,
		OrderNo:    order.OrderNo,
		OrderId:    order.ID,
		Name:       goods.Name,
		Number:     goods.Num,
		Status:     echoapp.TicketStatusNormal,
		UsedNumber: 0,
		Username:   address.Username,
		UsedAt:     nil,
		Cover:      goods.Cover,
		ComId:      order.ComId,
		UserId:     order.UserId,
		Rand:       int64(r),
		Source:     source,
		AddressID:  address.ID,
	}
	return ticket
}

func (tkSvr TicketService) GetTicketsByOrder(order *echoapp.Order) ([]*echoapp.Ticket, error) {
	var tickets []*echoapp.Ticket
	if err := tkSvr.db.Where("order_id = ?", order.ID).Find(&tickets).Error; err != nil {
		return nil, errors.Wrap(err, "db query tickets")
	}
	return tickets, nil
}

//验票的相关逻辑,因为和事务有关系交给调用着去更新
func (tkSvr TicketService) CheckTicket(ticket *echoapp.Ticket, num uint, staffID uint) error {
	if err := ticket.IsValid(); err != nil {
		return err
	}
	ticket.StaffID = staffID
	ticket.UsedNumber += num

	if ticket.UsedNumber > ticket.Number {
		return echoapp.ErrTicketNotEnough
	} else if ticket.UsedNumber == ticket.Number {
		ticket.Status = echoapp.TicketStatusUsed
	} else {

	}

	now := time.Now()
	ticket.UsedAt = &now
	return nil
}
