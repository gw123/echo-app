package services

import (
	"sync"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type OrderService struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewOrderService(db *gorm.DB) *OrderService {
	help := &OrderService{
		db: db,
	}
	return help
}

func (rsv *OrderService) PlaceOrder(order *echoapp.Order) error {
	rsv.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Order{})
	return rsv.db.Create(order).Error
}
func (rsv *OrderService) GetOrderById(c echo.Context, id uint) (*echoapp.Order, error) {
	order := &echoapp.Order{}
	res := rsv.db.Where("id=?", id).First(order)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetOrderById")
	}
	echoapp_util.ExtractEntry(c).Info("OrderID:%d", id)
	return order, nil
}

func (rsv *OrderService) GetUserPaymentOrder(c echo.Context, userId uint, from, limit int) ([]*echoapp.Order, error) {
	var orderlist []*echoapp.Order
	res := rsv.db.Where("user_id=?", userId).Offset(from * limit).Limit(limit).Find(orderlist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetUserPaymentOrder")
	}
	echoapp_util.ExtractEntry(c).Info("UserID:%d,from:%d,limit:%d", userId, from, limit)
	return orderlist, nil
}
func (rsv *OrderService) ModifyOrder(Order *echoapp.Order) error {
	return rsv.db.Save(Order).Error
}

func (rsv *OrderService) GetOrderList(c echo.Context, from, limit int) ([]*echoapp.GetOrderOptions, error) {
	var orderoptions []*echoapp.GetOrderOptions

	res := rsv.db.Table("orders").Offset(limit * from).Limit(limit).Find(orderoptions)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetOrderList")
	}
	return orderoptions, nil
}
