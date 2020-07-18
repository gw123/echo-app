package echoapp

import (
	"fmt"
	"time"
)

/***
 type 优惠券类型
	   register 注册就送
       daily    每日领取一张, 鼓励老用户常来
       daily-3  每日使用3张, sho
       once     限制每人只能领取一次
       normal   只要领取就能使用

 range 可以领取和使用的商品列表   [1, 2, 3, 4]
*/
const (
	CouponRangeTypeAll   = "all"
	CouponRangeTypeRange = "range"

	CouponTypeRegister = "register"
	CouponTypeDaily    = "daily"
	CouponTypeDaily3   = "daily3"
	CouponTypeOnce     = "once"
	CouponTypeNormal   = "normal"

	RedisGoodsCoupons    = "RedisGoodsCoupons:%d:%d"
	RedisPositionConpons = "RedisPositionCoupons:%d:%s"
	RedisActivityConpons = "RedisActivityCoupons:%d:%d"
	RedisAllGoodsConpons = "RedisAllGoodsCoupons:%d"
	RedisConpon          = "RedisCoupon:%d"
)

func FormatBannerListRedisKey(comId uint, position string) string {
	return fmt.Sprintf(RedisBannerListKey, comId, position)
}

func FormatActivityRedisKey(id uint) string {
	return fmt.Sprintf(RedisActivityKey, id)
}

func FormatGoodsCouponsKey(comId, goodsId uint) string {
	return fmt.Sprintf(RedisGoodsCoupons, comId, goodsId)
}

func FormatPositionCouponsKey(comId uint, position string) string {
	return fmt.Sprintf(RedisPositionConpons, comId, position)
}

func FormatActivityCouponsKey(comId uint, activityId uint) string {
	return fmt.Sprintf(RedisActivityConpons, comId, activityId)
}

func FormatAllGoodsCoupons(comId uint) string {
	return fmt.Sprintf(RedisAllGoodsConpons, comId)
}

func FormatCoupon(couponId uint) string {
	return fmt.Sprintf(RedisConpon, couponId)
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

/***
 type 优惠券类型
	   register 注册就送
       daily    每日领取一张, 鼓励老用户常来
       daily-3  每日使用3张, sho
       once     限制每人只能领取一次
       normal   只要领取就能使用

 range 可以领取和使用的商品列表   [1, 2, 3, 4]
*/
type Coupon struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	Type       string    `json:"type"`
	Cover      string    `json:"cover"`
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	Href       string    `json:"href"`        //跳转地址
	MinConsume uint      `json:"min_consume"` // 最低消费
	Amount     float32   `json:"amount"`      // 优惠券金额
	RangeType  string    `json:"range_type"`  //商品使用范围类型 all ,range , single
	Range      []uint    `json:"range"`       //商品适用范围
	Total      uint      `json:"total"`
	UsedTotal  uint      `json:"used_total"`
	ExpireAt   time.Time `json:"expire_at"` //整体过期时间
	Duration   uint      `json:"duration"`  //领取后可以保留的时间
	StartAt    time.Time `json:"start_at"`  //开始可以领取优惠券时间
	CreatedAt  time.Time `json:"created_at"`
}

func (*Coupon) TableName() string {
	return "coupons"
}

//用户优惠券
type UserCoupon struct {
	Id         uint      `gorm:"primary_key" json:"id"`
	CouponId   uint      `json:"coupon_id"` //优惠券
	UserId     uint      `json:"user_id"`
	DeletedAt  time.Time //删除优惠券
	CreatedAt  time.Time //领取时间
	StartAt    time.Time //可以开始使用时间
	UsedAt     time.Time //使用时间
	ExpireAt   time.Time //优惠券过期时间
	BaseCoupon Coupon    `json:"base_coupon"`
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
	//获取当前商品可以领取的优惠券列表
	GetCouponsByGoodsId(comId uint, goodsId uint) ([]*Coupon, error)
	//获取当前位置可以领取的优惠券, 首页,支付页面,购物车
	GetCouponsByPosition(ComId uint, position string) ([]*Coupon, error)
	//获取活动可以领取优惠券
	GetCouponsByActivity(ComId uint, activityId uint) ([]*Coupon, error)
	//获取当前订单可以使用的优惠券
	GetCouponsByOrder(ComId uint, order Order) ([]*Coupon, error)
	//获取用户领取的优惠券列表
	GetCouponsByUser(ComId uint, userId, lastId uint) ([]*Coupon, error)
}
