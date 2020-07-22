package echoapp

import (
	"fmt"
	"time"
)

func FormatBannerListRedisKey(comId uint, position string) string {
	return fmt.Sprintf(RedisBannerListKey, comId, position)
}

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

type Notify struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notify *Notify) TableName() string {
	return "activies"
}

type Banner struct {
	BannerBrief
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Title     string     `json:"title"`
	Visit     int        `json:"visit"`
	EndAt     time.Time  `json:"end_at"`
	GoodsId   int        `json:"goods_id"`
	ComId     uint       `json:"com_id"`
}

func (b *Banner) AfterFind() error {
	if b.Type == "goods" {
		b.Href = fmt.Sprintf("/pages/product/product?id=%d&com_id=%d", b.GoodsId, b.ComId)
	} else if b.Type == "activity" {
		b.Href = fmt.Sprintf("/pages/activity/detail?id=%d&com_id=%d", b.Id, b.ComId)
	}
	return nil
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

type BannerBrief struct {
	Id    uint   `gorm:"primary_key" json:"id"`
	Cover string `json:"cover"`
	Href  string `json:"href"`
	Type  string `json:"type"`
	//背景颜色
	Background string `json:"background"`
}

func (*BannerBrief) TableName() string {
	return "activies"
}

type ActivityService interface {
	GetNotifyList(comId uint, lastId, limit int) ([]*Notify, error)
	GetNotifyDetail(id int) (*Notify, error)
	//
	GetBannerList(comId uint, position string, limit int) ([]*BannerBrief, error)
	UpdateCachedBannerList(comId uint, position string) error
	GetCachedBannerList(comId uint, position string) ([]*BannerBrief, error)
	GetIndexBanner(comId uint) ([]*BannerBrief, error)
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
}
