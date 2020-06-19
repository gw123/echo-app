package services

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"strconv"
	"sync"
)

type OrderService struct {
	db    *gorm.DB
	mu    sync.Mutex
	redis *redis.Client
}

func NewOrderService(db *gorm.DB, redis *redis.Client) *OrderService {
	help := &OrderService{
		db:    db,
		redis: redis,
	}
	return help
}

func (oSvr *OrderService) GetTicketByCode(code string) (*echoapp.CodeTicket, error) {
	ticket, err := oSvr.DeTicketCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "DeTicketCode")
	}
	fmt.Printf("%+v \n", ticket)
	order, err := oSvr.GetOrderById(ticket.OrderId)
	if err != nil {
		return nil, errors.Wrap(err, "GetOrderById")
	}
	goods, err := oSvr.GetGoodsById(ticket.GoodsId)
	if err != nil {
		return nil, errors.Wrap(err, "DeTicketCode")
	}
	fmt.Printf("%+v \n", goods)
	ticket.Name = goods.Name
	ticket.Cover = goods.SmallCover

	var tickets []*echoapp.Ticket
	if goods.GoodsType == echoapp.GoodsTypeCombine {
		err := oSvr.db.Debug().
			Select("tickets.id,tickets.rand,tickets.used_number,tickets.used_at,tickets.status,tickets.number,goods.name,goods.small_cover as cover").
			Where("pid = ? and order_id = ?", ticket.GoodsId, ticket.OrderId).
			Joins("left join goods on goods.id = tickets.goods_id").
			Find(&tickets).Error
		if err != nil {
			return nil, errors.Wrap(err, "get tickets")
		}
	} else {
		tickets = append(tickets, ticket)
	}

	codeTicket := &echoapp.CodeTicket{
		BayAt:       ticket.CreatedAt,
		ComId:       ticket.ComId,
		GoodsCover:  goods.SmallCover,
		GoodsId:     goods.ID,
		GoodsName:   goods.Name,
		OrderNo:     ticket.OrderNo,
		OrderStatus: order.Status,
		Username:    ticket.Username,
		UserId:      ticket.UserId,
		Tickets:     tickets,
	}
	return codeTicket, nil
}

func (oSvr *OrderService) GetGoodsById(id int) (*echoapp.GoodsBrief, error) {
	goods := &echoapp.GoodsBrief{}
	if err := oSvr.db.Where("id = ?", id).First(goods).Error; err != nil {
		return nil, err
	}
	return goods, nil
}

func (oSvr *OrderService) DeTicketCode(code string) (*echoapp.Ticket, error) {
	rand, _ := strconv.ParseInt(code[0:8], 10, 64)
	temp, _ := strconv.ParseInt(code[8:], 10, 64)
	if rand == 0 || temp == 0 {
		return nil, errors.New("code is not vaild")
	}
	tId := temp - 1234 - rand
	ticket := &echoapp.Ticket{}
	if err := oSvr.db.Where("id = ?", tId).First(ticket).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	if ticket.Rand != rand {
		return nil, errors.New("校验失败")
	}

	return ticket, nil
}

func (oSvr *OrderService) PlaceOrder(order *echoapp.Order) error {
	oSvr.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Order{})
	return oSvr.db.Create(order).Error
}

func (oSvr *OrderService) GetOrderById(id uint) (*echoapp.Order, error) {
	order := &echoapp.Order{}
	res := oSvr.db.Where("id=?", id).First(order)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetOrderById")
	}
	return order, nil
}

func (oSvr *OrderService) GetOrderByOrderNo(orderNo string) (*echoapp.Order, error) {
	order := &echoapp.Order{}
	res := oSvr.db.Where("order_no=?", orderNo).First(order)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetOrderByOrderNo")
	}
	return order, nil
}

func (oSvr *OrderService) GetUserPaymentOrder(c echo.Context, userId uint, from, limit int) ([]*echoapp.Order, error) {
	var orderlist []*echoapp.Order
	res := oSvr.db.Where("user_id=?", userId).Offset(from * limit).Limit(limit).Find(orderlist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetUserPaymentOrder")
	}
	//echoapp_util.ExtractEntry(c).Info("UserID:%d,from:%d,limit:%d", userId, from, limit)
	return orderlist, nil
}
func (oSvr *OrderService) ModifyOrder(Order *echoapp.Order) error {
	return oSvr.db.Save(Order).Error
}

func (oSvr *OrderService) GetOrderList(c echo.Context, from, limit int) ([]*echoapp.GetOrderOptions, error) {
	var orderoptions []*echoapp.GetOrderOptions

	res := oSvr.db.Table("orders").Offset(limit * from).Limit(limit).Find(orderoptions)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetOrderList")
	}
	return orderoptions, nil
}
