package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/glog"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type OrderService struct {
	db      *gorm.DB
	mu      sync.Mutex
	redis   *redis.Client
	goodSvr echoapp.GoodsService
	actSvr  echoapp.ActivityService
}

func NewOrderService(db *gorm.DB,
	redis *redis.Client,
	goodsSvr echoapp.GoodsService,
	actSvr echoapp.ActivityService,
) *OrderService {
	help := &OrderService{
		db:      db,
		redis:   redis,
		goodSvr: goodsSvr,
		actSvr:  actSvr,
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
	glog.GetLogger().WithField("coupons", order.Coupons).Info("订单优惠券")
	//oSvr.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Order{})
	order.PayStatus = echoapp.OrderStatusUnpay
	order.OrderNo = fmt.Sprintf("%d%d%d", time.Now().Unix()/3600, order.UserId%9999, rand.Int31n(80000)+10000)
	if err := oSvr.goodSvr.IsValidCartGoodsList(order.GoodsList); err != nil {
		return errors.Wrap(err, "下单失败")
	}

	var (
		totalCouponAmount float32 = 0
		totalGoodsAmount  float32 = 0
	)
	for _, goods := range order.GoodsList {
		totalGoodsAmount += goods.RealPrice
	}
	glog.GetLogger().Info("校验优惠券")
	for _, coupon := range order.Coupons {
		realCoupon, err := oSvr.actSvr.GetUserCouponById(order.ComId, order.UserId, coupon.Id)
		if err != nil {
			glog.GetLogger().Error("下单失败,优惠券校验失败")
			return errors.Wrap(err, "下单失败,优惠券校验失败")
		}
		if realCoupon.Amount != coupon.Amount {
			glog.GetLogger().Error("下单失败,优惠券金额校验错误")
			return errors.New("下单失败,优惠券金额校验错误")
		}
		if realCoupon.IsExpire() {
			glog.GetLogger().Error("下单失败,优惠券已经过期")
			return errors.New("下单失败,优惠券已经过期")
		}
		if realCoupon.RangeType == echoapp.CouponRangeTypeAll {
			//todo 检查商品是否满足条件
		}
		totalCouponAmount += realCoupon.Amount
	}

	if totalGoodsAmount-totalCouponAmount != order.RealTotal {
		glog.GetLogger().Info("订单金额校验失败")
		return errors.New("订单金额校验失败")
	}

	tx := oSvr.db.Begin()
	var userCoupon echoapp.UserCoupon
	glog.GetLogger().Info("开始核销优惠券")
	for _, coupon := range order.Coupons {
		glog.GetLogger().WithField("coupon", coupon).Info("核销优惠券")
		if err := tx.Model(&userCoupon).
			Where("id = ?", coupon.Id).
			Update("used_at", time.Now()).Error; err != nil {
			tx.Rollback()
			return errors.Wrap(err, "优惠券核销失败")
		}
	}
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "订单保存失败")
	}
	tx.Commit()
	return nil
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

func (oSvr *OrderService) GetUserOrderList(c echo.Context, userId uint, status string, lastId uint, limit int) ([]*echoapp.Order, error) {
	var orderoptions []*echoapp.Order
	query := oSvr.db
	if lastId != 0 {
		query = query.Where("id < ? ", lastId)
	}

	switch status {
	case echoapp.OrderStatusUnpay:
		fallthrough
	case echoapp.OrderStatusPaid:
		fallthrough
	case echoapp.OrderStatusRefund:
		fallthrough
	case echoapp.OrderStatusCommented:
		query = query.Where("status = ?", status)
	case echoapp.OrderStatusShipping:
		fallthrough
	case echoapp.OrderStatusSigned:
		query = query.Where("status= ? and express_status=?", echoapp.OrderStatusPaid, status)
	default:
		glog.Warn("test === unknow")
	}

	res := query.Debug().Table("orders").
		Where("user_id = ?", userId).
		Limit(limit).Order("id desc").Find(&orderoptions)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetOrderList")
	}
	return orderoptions, nil
}

func (oSvr *OrderService) CancelOrder(o *echoapp.Order) error {
	if err := oSvr.db.Model(o).Where("user_id = ? and order_no = ?", o.UserId, o.OrderNo).
		Update("status", echoapp.OrderStatusCancel).Error; err != nil {
		return errors.Wrap(err, "cancelOrder")
	}
	return nil
}
func (oSvr *OrderService) StatisticCompanySalesByDate(start, end string) (*echoapp.CompanySalesSatistic, error) {
	salesStatistic := &echoapp.CompanySalesSatistic{}
	if err := oSvr.db.Table("orders").
		Where("created_at>=? and created_at<=?", start, end).
		Select("com_id,sum(total) as company_sales").Group("com_id").
		Find(&salesStatistic).Error; err != nil {
		return nil, errors.Wrap(err, "query")
	}
	salesStatistic.Date = fmt.Sprintf("%s--%s", start, end)
	if err := oSvr.db.Create(salesStatistic).Error; err != nil {
		return nil, errors.Wrap(err, "Create")
	}
	return salesStatistic, nil
}
func (oSvr *OrderService) StatisticComGoodsSalesByDate(start, end string, comId uint) (*echoapp.GoodsSalesSatistic, error) {
	salesStatistic := &echoapp.GoodsSalesSatistic{}
	if err := oSvr.db.Table("orders").Where("com_id=?", comId).
		Where("created_at>=? and created_at<=?", start, end).
		Select("goods_id,sum(total) as goods_sales").Group("goods_id").
		Find(&salesStatistic).Error; err != nil {
		return nil, errors.Wrap(err, "query")
	}
	salesStatistic.Date = fmt.Sprintf("%s--%s", start, end)
	salesStatistic.ComId = comId
	if err := oSvr.db.Create(salesStatistic).Error; err != nil {
		return nil, errors.Wrap(err, "Create")
	}
	return salesStatistic, nil
}
