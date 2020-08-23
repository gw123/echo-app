package echoapp

import (
	"fmt"
	"time"
)

func FormatActivityRedisKey(id uint) string {
	return fmt.Sprintf(RedisActivityKey, id)
}

type Activity struct {
	ComId uint `json:"com_id"`
	BannerBrief
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Title     string     `json:"title"`
	Visit     int        `json:"visit"`
	EndAt     time.Time  `json:"end_at"`
	Type      string     `json:"type"`
	GoodsId   int        `json:"goods_id"`
	Body      string     `json:"body"`
}

func (a *Activity) TableName() string {
	return "activies"
}

//活动优惠券
type ActivityCoupon struct {
	Id         uint `gorm:"primary_key" json:"id"`
	ActivityId uint `json:"aid"`
	CouponId   uint `json:"cid"`
	DeletedAt  time.Time
}

func (*ActivityCoupon) TableName() string {
	return "activity_coupons"
}

func (*UserCoupon) TableName() string {
	return "user_coupons"
}

type ActivityService interface {
	//GetActivityById(id int) (*Banner, error)
	AddActivityPv(goodsId uint) error
	GetActivityList(comId uint, lastId uint, limit int) ([]*Activity, error)
	GetActivityDetail(id uint) (*Activity, error)

	//coupon
	GetCachedCouponById(couponId uint) (*Coupon, error)
	GetCachedCouponsByIds(couponIds []uint) ([]*Coupon, error)

	//获取当前商品可以领取的优惠券列表
	GetCouponsByGoodsId(comId uint, goodsId uint) ([]*Coupon, error)
	//获取当前位置可以领取的优惠券, 首页,支付页面,购物车
	GetCouponsByPosition(ComId uint, position string) ([]*Coupon, error)
	//获取活动可以领取优惠券
	GetCouponsByActivity(ComId uint, activityId uint) ([]*Coupon, error)
	//获取当前订单可以使用的优惠券
	GetUserCouponsByOrder(ComId uint, order *Order) ([]*Coupon, []*Coupon, error)
	//获取用户领取的优惠券列表
	GetUserCoupons(ComId uint, userId, lastId uint) ([]*UserCoupon, error)
	//创建用户优惠券
	CreateUserCoupon(comId uint, userId uint, couponId uint) error
	UpdateCachedCouponsByComId(comId uint, lastId uint) ([]*Coupon, error)
	GetUserCouponById(comId, userId, couponId uint) (*Coupon, error)
	//获取商品页面关联商品的一个活动
	GetGoodsActivity(comID uint, goodsID uint) (*Activity, error)
}
