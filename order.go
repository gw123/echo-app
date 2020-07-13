package echoapp

import (
	"encoding/json"
	"fmt"
	"github.com/gw123/glog"
	"time"

	"github.com/labstack/echo"
)

const (
	OrderStatusUnpay     = "unpay"
	OrderStatusPaid      = "paid"
	OrderStatusRefund    = "refund"
	OrderStatusShipping  = "shipping"
	OrderStatusSigned    = "signed"
	OrderStatusCancel    = "cancel"
	OrderStatusCommented = "commented"
)

type Order struct {
	ID            uint             `json:"id" gorm:"primary_key"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	PaidAt        time.Time        `json:"paid_at"`
	Source        string           `json:"source"`
	PayMethod     string           `json:"pay_method"`
	ComId         uint             `json:"com_id"`
	ShopId        uint             `json:"shop_id"`
	OrderNo       string           `json:"order_no"`
	UserId        uint             `json:"user_id"`
	SellerId      uint             `json:"seller_id"`
	InviterId     uint             `json:"inviter_id"`
	PayStatus     string           `json:"pay_status" gorm:"column:status"`
	Status        string           `json:"status" gorm:"-"`
	ExpressStatus string           `json:"express_status"`
	AddressId     uint             `json:"address_id"`
	Total         float32          `json:"total"`
	RealTotal     float32          `json:"real_total"`
	GoodsList     []*CartGoodsItem `json:"goodsList" gorm:"-"`
	GoodsListStr  string           `json:"-" gorm:"column:goodslist"`
	GoodsType     string           `json:"goods_type"`
	TransactionId string           `json:"transaction_id"`
	Note          string           `json:"note"`
	Info          string           `json:"info"`
	//Score         string           `score` //积分
}

func (o *Order) BeforeSave() error {
	data, err := json.Marshal(o.GoodsList)
	if err != nil {
		return err
	}
	o.GoodsListStr = string(data)
	if o.PaidAt.IsZero() {
		o.PaidAt = time.Now()
	}
	return nil
}

func (o *Order) AfterFind() error {
	switch o.PayStatus {
	case OrderStatusUnpay:
		fallthrough
	case OrderStatusRefund:
		o.Status = o.PayStatus
	case OrderStatusPaid:
		if o.ExpressStatus == OrderStatusShipping {
			o.Status = OrderStatusShipping
		} else if o.ExpressStatus == OrderStatusSigned {
			o.Status = OrderStatusSigned
		} else if o.ExpressStatus == OrderStatusCommented {
			o.Status = OrderStatusCommented
		}
	default:
		glog.Warn("unknow pay_status")
	}

	o.GoodsList = make([]*CartGoodsItem, 0)
	//	glog.Infof("goodsListStr %s", o.GoodsListStr)
	err := json.Unmarshal([]byte(o.GoodsListStr), &o.GoodsList)
	return err
}

type GetOrderOptions struct {
	PayMethod     string    `json:"pay_method"`
	ShopId        int       `json:"shop_id"`
	OrderNo       string    `json:"order_no"`
	Status        string    `json:"status"`
	Total         float32   `json:"total"`
	GoodsList     string    `json:"goods_list"`
	GoodsType     string    `json:"goods_type"`
	TransactionId string    `json:"transaction_id"`
	Note          string    `json:"note"`
	Info          string    `json:"info"`
	CreatedAt     time.Time `json:"created_at"`
	PaidAt        time.Time `json:"paid_at"`
	//Score         string    `score`
}

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

type OrderService interface {
	GetTicketByCode(code string) (*CodeTicket, error)
	//保存上传的资源到数据库
	PlaceOrder(order *Order) error
	//通过资源ID查找资源
	GetOrderById(id uint) (*Order, error)
	GetOrderByOrderNo(orderNo string) (*Order, error)

	ModifyOrder(order *Order) error

	GetUserPaymentOrder(c echo.Context, userId uint, from, limit int) ([]*Order, error)
	//查看资源文件 ，每页有 limit 条数据
	GetOrderList(c echo.Context, from, limit int) ([]*GetOrderOptions, error)
	GetUserOrderList(c echo.Context, userId uint, status string, lastId uint, limit int) ([]*Order, error)
	CancelOrder(o *Order) error
}
