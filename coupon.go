package echoapp

import (
	"encoding/json"
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
	//
	CouponStatusNotUse = "notuse"
	CouponStatusUsed   = "used"
	CouponStatusAll    = "all"
	//
	RedisGoodsCoupons      = "RedisGoodsCoupons:%d:%d"
	RedisPositionConpons   = "RedisPositionCoupons:%d:%s"
	RedisActivityConpons   = "RedisActivityCoupons:%d:%d"
	RedisAllGoodsConpons   = "RedisAllGoodsCoupons:%d"
	RedisConpon            = "RedisCoupon:%d"
	RedisMutexCreateCoupon = "RedisMutexCreateCoupon:%d"
)

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

func FormatRedisMutexCreateCoupon(couponId uint) string {
	return fmt.Sprintf(RedisMutexCreateCoupon, couponId)
}

type CouponBase struct {
	Id     uint    `gorm:"primary_key" json:"id"`
	Name   string  `json:"name"`
	Amount float32 `json:"amount"`
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
	Name       string    `json:"name"`
	Desc       string    `json:"desc"`
	Href       string    `json:"href"`                  //跳转地址
	MinConsume uint      `json:"min_consume"`           // 最低消费
	Amount     float32   `json:"amount"`                // 优惠券金额
	RangeType  string    `json:"range_type"`            //商品使用范围类型 all ,range , single
	RangeStr   string    `json:"-" gorm:"column:range"` //商品适用范围
	Range      []uint    `json:"range" gorm:"-"`        //商品适用范围
	Total      uint      `json:"-"`
	UsedTotal  uint      `json:"-"`
	ExpireAt   time.Time `json:"-"`        //整体过期时间
	Duration   uint      `json:"-"`        //领取后可以保留的时间
	StartAt    time.Time `json:"start_at"` //开始可以领取优惠券时间
	CreatedAt  time.Time `json:"created_at"`
}

func (*Coupon) TableName() string {
	return "coupons"
}

func (c *Coupon) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, c)
}

func (c *Coupon) IsExpire() bool {
	return c.ExpireAt.Sub(time.Now()) <= 0
}

// 用户优惠券
type UserCoupon struct {
	Id         uint       `gorm:"primary_key" json:"id"`
	CouponId   uint       `json:"coupon_id"` //优惠券
	UserId     uint       `json:"user_id"`
	DeletedAt  *time.Time //删除优惠券
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at"` //领取时间
	StartAt    time.Time  //可以开始使用时间
	UsedAt     *time.Time //使用时间
	ExpireAt   time.Time  //优惠券过期时间
	BaseCoupon *Coupon    `json:"base_coupon" gorm:"-"`
	ComId      uint
}

func (*UserCoupon) TableName() string {
	return "user_coupons"
}
