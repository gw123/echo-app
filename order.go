package echoapp

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"

	"github.com/gw123/glog"
	"github.com/labstack/echo"
)

const (
	OrderStatusUnpay     = "unpay"
	OrderStatusPaid      = "paid"
	OrderStatusRefund    = "refund"
	OrderStatusShipping  = "shipping"
	OrderStatusToShip    = "unship"
	OrderStatusSigned    = "signed"
	OrderStatusCancel    = "cancel"
	OrderStatusCommented = "commented"

	OrderPayStatusUnpay  = "unpay"
	OrderPayStatusPaid   = "paid"
	OrderPayStatusRefund = "refund"
)

type Order struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PaidAt    time.Time `json:"paid_at"`
	Source    string    `json:"source"`
	PayMethod string    `json:"pay_method"`
	ComId     uint      `json:"com_id"`
	ShopId    uint      `json:"shop_id"`
	OrderNo   string    `json:"order_no"`
	UserId    uint      `json:"user_id"`
	SellerId  uint      `json:"seller_id"`
	InviterId uint      `json:"inviter_id"`
	PayStatus string    `json:"pay_status" gorm:"column:status"`
	// Status 是一个根据pay_status和express_status计算出来的字段
	Status        string           `json:"status" gorm:"-"`
	ExpressStatus string           `json:"express_status"`
	AddressId     uint             `json:"address_id"`
	Total         float32          `json:"total"`
	RealTotal     float32          `json:"real_total"`
	GoodsList     []*CartGoodsItem `json:"goodsList" gorm:"-"`
	GoodsListStr  string           `json:"-" gorm:"column:goodslist"`

	Coupons       []*CouponBase `json:"coupons" gorm:"-"`
	CouponsStr    string        `json:"-" gorm:"column:coupons"`
	GoodsType     string        `json:"goods_type"`
	TransactionId string        `json:"transaction_id"`
	Note          string        `json:"note"`
	Info          string        `json:"info"`
	ClientIP      string        `json:"client_ip"`
	ClientType    string        `json:"client_type"`
	Tickets       []*Ticket     `json:"tickets" gorm:"-"`
	Address       *Address      `json:"address" gorm:"-"`
	//Score         string           `score` //积分
}

func (o *Order) OrderCheck() error {
	if len(o.GoodsList) == 0 {
		return AppErrOrderFromat
	}

	switch o.GoodsType {
	case GoodsTypeVip:
		if len(o.GoodsList) != 1 {
			return errors.Errorf("vip订单商品只能有一个")
		}
	case GoodsTypeTicket:
	case GoodsTypeRoom:
	}
	return nil
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

	if len(o.Coupons) > 0 {
		couponStr, err := json.Marshal(o.Coupons)
		if err != nil {
			return err
		}
		o.CouponsStr = string(couponStr)
	} else {
		o.CouponsStr = "[]"
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
		glog.Warn("unknow pay_status:" + o.PayStatus)
	}

	o.GoodsList = make([]*CartGoodsItem, 0)
	//	glog.Infof("goodsListStr %s", o.GoodsListStr)
	err := json.Unmarshal([]byte(o.GoodsListStr), &o.GoodsList)
	return err
}

type OrderGoods struct {
	ID        int64     `json:"id"`
	ComID     uint      `json:"com_id"`
	OrderID   uint      `json:"order_id"`
	GoodsID   uint      `json:"goods_id"`
	Num       uint      `json:"num"`
	Status    string    `json:"status"`
	RealPrice float32   `json:"real_price"`
	CreatedAt time.Time `json:"created_at"`
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
	CreatedAt     time.Time `json:"created_at"`
	PaidAt        time.Time `json:"paid_at"`
	Info          string    `json:"info"`
	//Score         string    `score`
}

type CompanySalesSatistic struct {
	ID            uint    `gorm:"primary_key"`
	AllSalesTotal float64 `json:"all_sales_total"`
	Date          string  `json:"date"`
	ComId         uint    `json:"com_id"`
	//GoodsSalesTotal int64  `json:"goods_sales"`
	//GoodsId         int    `json:"goods_id"`
	//Status string `json:"status"`
}

type GoodsSalesSatistic struct {
	ID uint `gorm:"primary_key"`
	//AllSalesTotal   int64  `json:"company_sales"`
	Date            string  `json:"date"`
	ComId           uint    `json:"com_id"`
	GoodsSalesTotal float64 `json:"goods_sales" gorm:"column:goods_sales"`
	GoodsId         int     `json:"goods_id"`
	//Status          string `json:"status"`
}

func (*CompanySalesSatistic) TableName() string {
	return "company_sales"
}
func (*GoodsSalesSatistic) TableName() string {
	return "goods_sales "
}

type UnifiedOrderResp struct {
	WxPreOrderResponse
	OrderNo string `json:"order_no"`
}

type OrderService interface {
	GetTicketByCode(code string) (*CodeTicket, error)

	//
	UniPreOrder(order *Order, user *User) (*UnifiedOrderResp, error)
	//预下单校验订单的接口
	PreCheckOrder(order *Order) error

	//查询订单支付状态
	QueryOrderAndUpdate(order *Order, shouldStatus string) (*Order, error)

	GetUserPaymentOrder(c echo.Context, userId uint, from, limit int) ([]*Order, error)
	//查看资源文件 ，每页有 limit 条数据
	GetUserOrderList(c echo.Context, userId uint, status string, lastId uint, limit int) ([]*Order, error)

	//取消订单
	CancelOrder(order *Order) error
	//验票
	CheckTicket(code string, num uint, staffID uint) error
	//通过资源ID查找资源
	GetOrderById(id uint) (*Order, error)
	//
	GetOrderByOrderNo(orderNo string) (*Order, error)
	//
	GetUserOrderDetail(ctx echo.Context, userId uint, orderNo string) (*Order, error)

	//退款
	Refund(order *Order, user *User) error

	WxPayCallback(ctx echo.Context) error

	WxRefundCallback(ctx echo.Context) error

	QueryRefundOrderAndUpdate(order *Order) (*Order, error)

	// StatisticComGoodsSalesByDate(start, end string, comId uint) (*GoodsSalesSatistic, error)

	// StatisticCompanySalesByDate(start, end string) (*CompanySalesSatistic, error)
}
