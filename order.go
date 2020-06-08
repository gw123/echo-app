package echoapp

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Order struct {
	gorm.Model
	Source        string  `json:"source"`
	PayMethod     string  `json:"pay_method"`
	ComId         int     `json:"com_id"`
	ShopId        int     `json:"shop_id"`
	OrderNo       string  `json:"order_no"`
	UserId        uint    `json:"user_id"`
	InviterId     int     `json:"inviter_id"`
	Status        string  `json:"status"`
	Total         float32 `json:"total"`
	RealTotal     float32 `json:"real_total"`
	PaidAt        time.Time
	GoodsList     string `json:"goods_list"`
	GoodsType     string `json:"goods_type"`
	TransactionId string `json:"transaction_id"`
	Note          string `json:"note"`
	Info          string `json:"info"`
	Score         string `score` //积分
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
	CreatedAt     time.Time `created_at`
	PaidAt        time.Time `json:"paid_at"`
	Score         string    `score`
}
type OrderService interface {
	//保存上传的资源到数据库
	PlaceOrder(order *Order) error
	//通过资源ID查找资源
	GetOrderById(c echo.Context, id uint) (*Order, error)

	ModifyOrder(order *Order) error

	GetUserPaymentOrder(c echo.Context, userId uint, from, limit int) ([]*Order, error)
	//查看资源文件 ，每页有 limit 条数据
	GetOrderList(c echo.Context, from, limit int) ([]*GetOrderOptions, error)
}
