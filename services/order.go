package services

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
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
	wechat  echoapp.WechatService
	tkSvr   echoapp.TicketService
}

func NewOrderService(db *gorm.DB,
	redis *redis.Client,
	goodsSvr echoapp.GoodsService,
	actSvr echoapp.ActivityService,
	wechat echoapp.WechatService,
	tkSvr echoapp.TicketService,
) *OrderService {
	help := &OrderService{
		db:      db,
		redis:   redis,
		goodSvr: goodsSvr,
		actSvr:  actSvr,
		wechat:  wechat,
		tkSvr:   tkSvr,
	}
	return help
}

func (oSvr *OrderService) GetGoodsById(id uint) (*echoapp.GoodsBrief, error) {
	goods := &echoapp.GoodsBrief{}
	if err := oSvr.db.Where("id = ?", id).First(goods).Error; err != nil {
		return nil, err
	}
	return goods, nil
}

//生成随机订单 增加
func (oSvr *OrderService) makeTicketNo(userId uint) string {
	date := time.Now().Format("200601021504")
	r := rand.Int31n(89999) + 10000
	idHash := userId % 1137
	str := fmt.Sprintf("%s%d%d", date, idHash, r)
	return str
}

func (oSvr *OrderService) GetTicketByCode(code string) (*echoapp.CodeTicket, error) {
	ticket, err := oSvr.DeTicketCode(code)
	if err != nil {
		return nil, errors.Wrap(err, "DeTicketCode")
	}
	order, err := oSvr.GetOrderById(ticket.OrderId)
	if err != nil {
		return nil, errors.Wrap(err, "GetOrderById")
	}
	goods, err := oSvr.GetGoodsById(ticket.GoodsId)
	if err != nil {
		return nil, errors.Wrap(err, "DeTicketCode")
	}
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

func (oSvr *OrderService) UniPreOrder(order *echoapp.Order, user *echoapp.User) (*echoapp.UnifiedOrderResp, error) {
	///glog.DefaultLogger().Infof("user: %+v,openid:%s", user, user.Openid)
	resp, err := oSvr.wechat.UnifiedOrder(order, user.Openid)
	if err != nil {
		return nil, errors.Wrap(err, "下单失败")
	}

	return &echoapp.UnifiedOrderResp{
		WxPreOrderResponse: *resp,
		OrderNo:            order.OrderNo,
	}, nil
}

func (oSvr *OrderService) CreateOrderTickets(order *echoapp.Order, db *gorm.DB) error {
	order.Tickets = make([]*echoapp.Ticket, 0)
	for _, cartGoods := range order.GoodsList {
		goods, err := oSvr.goodSvr.GetGoodsById(cartGoods.GoodsId)
		if err != nil {
			return err
		}
		//glog.Infof("ID: %d ,goodsType %s", goods.ID, goods.GoodsType)
		if goods.GoodsType == echoapp.GoodsTypeTicket || goods.GoodsType == echoapp.GoodsTypeRoom {
			ticket := oSvr.tkSvr.PreCreateTicket(order, order.Source, order.Address, cartGoods)
			if err := db.Save(ticket).Error; err != nil {
				return err
			}
			order.Tickets = append(order.Tickets, ticket)
			glog.Infof("createTicket ID:%d, ticketName:%s ,OrderNo:%s ,Mobile:%s", ticket.ID, ticket.Name, ticket.OrderNo, ticket.Mobile)
		}
	}
	return nil
}

func (oSvr *OrderService) RefundOrderTickets(order *echoapp.Order, tx *gorm.DB) ([]*echoapp.Ticket, error) {
	tickets, err := oSvr.tkSvr.GetTicketsByOrder(order)
	if err != nil {
		return nil, errors.Wrap(err, "GetTicketsByOrder")
	}

	for _, ticket := range tickets {
		ticket.Status = echoapp.TicketstatusRefund
		if err := tx.Model(ticket).UpdateColumn("status", ticket.Status).Error; err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "create ticket")
		}
	}
	return tickets, nil
}

func (oSvr *OrderService) QueryOrderAndUpdate(order *echoapp.Order, shouldStatus string) (*echoapp.Order, error) {
	orderBeferStatuss := order.Status
	if order.Status == echoapp.OrderStatusPaid || order.Status == echoapp.OrderStatusRefund {
		//订单已经是最终的状态
		return order, nil
	}

	var currentStatus string
	var err error
	switch order.ClientType {
	case echoapp.ClientWxOfficial, echoapp.ClientWxMiniApp:
		if shouldStatus == echoapp.OrderPayStatusPaid {
			currentStatus, err = oSvr.wechat.QueryOrder(order)
		} else if shouldStatus == echoapp.OrderStatusRefund {
			currentStatus, err = oSvr.wechat.QueryRefund(order)
		}
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("暂时不支持该支付方式")
	}

	if currentStatus == echoapp.OrderPayStatusUnpay {
		return order, nil
	}

	//从未支付到支付成功
	if orderBeferStatuss == echoapp.OrderStatusUnpay && currentStatus == echoapp.OrderStatusPaid {
		order.Status = echoapp.OrderPayStatusPaid
		tx := oSvr.db.Begin()
		//创建票
		err := oSvr.CreateOrderTickets(order, tx)
		if err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "create ticket")
		}

		if err := tx.Model(order).UpdateColumn("status", order.Status).Error; err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "change order currentStatus")
		}
		tx.Commit()
		//发送通知
	} else if currentStatus == echoapp.OrderPayStatusRefund {
		order.Status = echoapp.OrderStatusRefund
		tx := oSvr.db.Begin()
		_, err := oSvr.RefundOrderTickets(order, tx)
		if err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "refund ticket")
		}
		if err := tx.Model(order).UpdateColumn("status", order.Status).Error; err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, "change order currentStatus")
		}
		tx.Commit()

	}
	return order, nil
}

func (oSvr *OrderService) GetUserAddress(addrId uint) (*echoapp.Address, error) {
	panic(addrId)
}

func (oSvr *OrderService) PreCheckOrder(order *echoapp.Order) error {
	//todo 校验ip , 校验客户端类型
	glog.JsonLogger().WithField("coupons", order.Coupons).Info("订单优惠券")
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
		totalGoodsAmount += goods.RealPrice * float32(goods.Num)
	}

	if len(order.Coupons) > 0 {
		glog.JsonLogger().Info("校验优惠券")
		for _, coupon := range order.Coupons {
			realCoupon, err := oSvr.actSvr.GetUserCouponById(order.ComId, order.UserId, coupon.Id)
			if err != nil {
				glog.JsonLogger().Error("下单失败,优惠券校验失败")
				return errors.Wrap(err, "下单失败,优惠券校验失败")
			}
			if realCoupon.Amount != coupon.Amount {
				glog.JsonLogger().Error("下单失败,优惠券金额校验错误")
				return errors.New("下单失败,优惠券金额校验错误")
			}
			if realCoupon.IsExpire() {
				glog.JsonLogger().Error("下单失败,优惠券已经过期")
				return errors.New("下单失败,优惠券已经过期")
			}
			if realCoupon.RangeType == echoapp.CouponRangeTypeAll {
				//todo 检查商品是否满足条件
			}
			totalCouponAmount += realCoupon.Amount
		}
	}

	realTotal := totalGoodsAmount - totalCouponAmount
	if realTotal != order.RealTotal {
		glog.JsonLogger().Info("订单金额校验失败 %f , %f", realTotal, order.RealTotal)
		return errors.New("订单金额校验失败")
	}

	tx := oSvr.db.Begin()
	var userCoupon echoapp.UserCoupon
	if len(order.Coupons) > 0 {
		glog.JsonLogger().Info("开始核销优惠券")
		for _, coupon := range order.Coupons {
			glog.JsonLogger().WithField("coupon", coupon).Info("核销优惠券")
			if err := tx.Model(&userCoupon).
				Where("id = ?", coupon.Id).
				Update("used_at", time.Now()).Error; err != nil {
				tx.Rollback()
				return errors.Wrap(err, "优惠券核销失败")
			}
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

func (oSvr *OrderService) GetUserPaymentOrder(ctx echo.Context, userId uint, from, limit int) ([]*echoapp.Order, error) {
	var orderlist []*echoapp.Order
	res := oSvr.db.Where("user_id=?", userId).Offset(from * limit).Limit(limit).Find(orderlist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "OrderService->GetUserPaymentOrder")
	}
	//echoapp_util.ExtractEntry(c).Info("UserID:%d,from:%d,limit:%d", userId, from, limit)
	return orderlist, nil
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

func (oSvr *OrderService) GetUserOrderDetial(ctx echo.Context, userId uint, orderNo string) (*echoapp.Order, error) {
	var (
		order echoapp.Order
	)

	if err := oSvr.db.Table("orders").
		Where("user_id = ? and order_no = ?", userId, orderNo).
		First(&order).Error; err != nil {
		return nil, errors.Wrap(err, "db query")
	}

	if order.PayStatus == echoapp.OrderPayStatusPaid && order.GoodsType == echoapp.GoodsTypeTicket {
		ticketList, err := oSvr.tkSvr.GetTicketsByOrder(&order)
		if err != nil {
			return nil, errors.Wrap(err, "db query")
		}
		order.Tickets = ticketList
	}

	return &order, nil
}

func (oSvr *OrderService) CancelOrder(o *echoapp.Order) error {
	if err := oSvr.db.Model(o).Where("user_id = ? and order_no = ?", o.UserId, o.OrderNo).
		Update("status", echoapp.OrderStatusCancel).Error; err != nil {
		return errors.Wrap(err, "cancelOrder")
	}
	return nil
}

func (oSvr *OrderService) CheckTicket(code string, num uint, staffID uint) error {
	ticket, err := oSvr.tkSvr.GetTicketByCode(code)
	if err != nil {
		return errors.Wrap(err, "GetTicketByCode")
	}
	if err := ticket.IsValid(); err != nil {
		return err
	}

	order, err := oSvr.GetOrderById(ticket.OrderId)
	if err != nil {
		return err
	}
	err = oSvr.tkSvr.CheckTicket(ticket, num, staffID)
	if err != nil {
		return errors.Wrap(err, "CheckTicket")
	}
	order.ExpressStatus = echoapp.OrderStatusSigned

	tx := oSvr.db.Begin()
	if err := tx.Save(ticket).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update ticket")
	}

	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "update order")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.Wrap(err, "tx err")
	}
	return nil
}

func (oSvr *OrderService) WxPayCallback(ctx echo.Context) error {
	resp, err := oSvr.wechat.PayCallback(ctx.Request())
	if err != nil {
		return errors.Wrap(err, "WxPayCallback")
	}
	order := &echoapp.Order{}
	if err := oSvr.db.Where("transaction_id = ?", resp.TransactionId).First(order).Error; err != nil {
		return errors.Wrap(err, "WxPayCallback")
	}
	// 调用query 接口再次确认支付结果是否ok,并且更新订单状态
	_, err = oSvr.QueryOrderAndUpdate(order, echoapp.OrderPayStatusPaid)
	if err != nil {
		echoapp_util.ExtractEntry(ctx).WithError(err).Error("订单状查询失败")
		return errors.Wrap(err, "WxPayCallback")
	}
	return nil
}

func (oSvr *OrderService) Refund(order *echoapp.Order, user *echoapp.User) error {
	oSvr.wechat.Refund(order, user.Openid)
	return nil
}

func (oSvr *OrderService) WxRefundCallback(ctx echo.Context) error {
	resp, err := oSvr.wechat.RefundCallback(ctx.Request())
	if err != nil {
		return errors.Wrap(err, "WxPayCallback")
	}
	order := &echoapp.Order{}
	if err := oSvr.db.Where("transaction_id = ?", resp.TransactionId).First(order).Error; err != nil {
		return errors.Wrap(err, "WxPayCallback")
	}
	// 调用query 接口再次确认支付结果是否ok,并且更新订单状态
	_, err = oSvr.QueryOrderAndUpdate(order, echoapp.OrderStatusRefund)
	if err != nil {
		echoapp_util.ExtractEntry(ctx).WithError(err).Error("订单状查询失败")
		return errors.Wrap(err, "WxPayCallback")
	}
	return nil
}

func (oSvr *OrderService) QueryRefundOrderAndUpdate(order *echoapp.Order) (*echoapp.Order, error) {
	return oSvr.QueryOrderAndUpdate(order, echoapp.OrderPayStatusRefund)
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
