package echoapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gw123/echo-app/observer"

	"github.com/gw123/glog"
)

const (
	// 商品关联活动的状态
	GoodsActivityStatusOffline = 0
	GoodsActivityStatusOnline  = 1
	GoodsActivityStatusOverDue = 2

	//  用户参加活动的当前状态
	UserActivityStatusSuccess = 1
	UserActivityStatusIng     = 2
	UserActivityStatusOverdue = 3
	UserActivityStatusFail    = 4

	// 奖品变化的类型 添加或者消耗
	AwardHistoryTypeActivity = "activity"
	AwardHistoryTypeVIP      = "vip"
	AwardHistoryTypeUse      = "used"

	// 活动驱动类型
	ActivityDriverTypeActivity = "score"
	ActivityDriverTypeVIP      = "vip"
	ActivityDriverTypeAward    = "award"
)

func FormatActivityRedisKey(id uint) string {
	return fmt.Sprintf(RedisActivityKey, id)
}

// 活动中的奖品
type SubReward struct {
	GoodsId   uint   `json:"goods_id"`
	GoodsName string `json:"goods_name"`
	Num       uint   `json:"num"`
}

// 活动中的优惠券
type SubCoupon struct {
	CouponId   uint   `json:"coupon_id"`
	CouponName string `json:"coupon_name"`
	Num        uint   `json:"num"`
}

type Activity struct {
	ComId uint `json:"com_id"`
	BannerBrief
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  *time.Time   `sql:"index"`
	Title      string       `json:"title"`
	Visit      int          `json:"visit"`
	EndAt      time.Time    `json:"end_at"`
	Type       string       `json:"type"`
	Body       string       `json:"body"`
	Driver     string       `json:"driver"`
	RewardsStr string       `json:"-" gorm:"column:rewards"`
	Rewards    []*SubReward `json:"rewards" gorm:"-"`
	CouponsStr string       `json:"-" gorm:"column:coupons"`
	Coupons    []*SubCoupon `json:"coupon" gorm:"-"`
	Score      string       `json:"score" `
}

func (a *Activity) TableName() string {
	return "activies"
}

func (a *Activity) AfterFind() error {
	if a.RewardsStr == "" {
		a.Rewards = make([]*SubReward, 0)
	} else {
		if err := json.Unmarshal([]byte(a.RewardsStr), &a.Rewards); err != nil {
			glog.Errorf("Cart AfterFind %s", err)
		}
	}

	if a.CouponsStr == "" {
		a.Coupons = make([]*SubCoupon, 0)
	} else {
		if err := json.Unmarshal([]byte(a.CouponsStr), &a.Coupons); err != nil {
			glog.Errorf("Cart AfterFind %s", err)
		}
	}
	return nil
}

func (a *Activity) BeforeSave() error {
	data, err := json.Marshal(a.Rewards)
	if err != nil {
		glog.Errorf("Cart BeforeSave %s", err)
	}
	a.RewardsStr = string(data)

	data, err = json.Marshal(a.Coupons)
	if err != nil {
		glog.Errorf("Cart BeforeSave %s", err)
	}
	a.CouponsStr = string(data)
	return nil
}

// 商品关联的活动
type GoodsActivity struct {
	ID         uint      `json id:"id"`
	GoodsID    uint      `json:"goods_id"`
	ActivityID uint      `json:"activity_id"`
	Status     uint8     `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

//  记录用户参加的活动，和活动的状态
type UserActivity struct {
	ID         uint       `json id:"id"`
	UserID     uint       `json:"user_id"`
	ActivityID uint       `json:"activity_id"`
	Status     uint8      `json:"status"`
	Metadata   string     `json:"metadata"`
	Log        string     `json:"log"`
	CreatedAt  time.Time  `json:"created_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

// 用户通过参加活动领取的奖品
type UserAward struct {
	ID        uint        `json:"id"`
	ComID     uint        `json:"com_id"`
	UserID    uint        `json:"user_id"`
	GoodsID   uint        `json:"goods_id"`
	Num       uint        `json:"num"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Goods     *GoodsBrief `gorm:"-" json:"goods"`
}

// 用户领取和核销商品的历史记录表
type AwardHistory struct {
	ID             uint        `json:"id"`
	ComID          uint        `json:"com_id"`
	UserID         uint        `json:"user_id"`
	StaffID        uint        `json:"staff_id"`
	GoodsID        uint        `json:"goods_id"`
	UserActivityID uint        `json:"user_activity_id"`
	ActivityID     uint        `json:"activity_id"`
	Method         string      `json:"method"`
	Num            int         `json:"num"`
	CreatedAt      time.Time   `json:"created_at"`
	Goods          *GoodsBrief `gorm:"-" json:"goods"`
}

// 活动优惠券
type ActivityCoupon struct {
	Id         uint `gorm:"primary_key" json:"id"`
	ActivityId uint `json:"aid"`
	CouponId   uint `json:"cid"`
	DeletedAt  time.Time
}

func (*ActivityCoupon) TableName() string {
	return "activity_coupons"
}

type ActivityService interface {
	// GetActivityById(id int) (*Banner, error)
	AddActivityPv(goodsId uint) error
	GetActivityList(comId uint, lastId uint, limit int) ([]*Activity, error)
	GetActivityDetail(id uint) (*Activity, error)

	// coupon
	GetCachedCouponById(couponId uint) (*Coupon, error)
	GetCachedCouponsByIds(couponIds []uint) ([]*Coupon, error)

	// 获取当前位置可以领取的优惠券, 首页,支付页面,购物车
	GetCouponsByPosition(ComId uint, position string) ([]*Coupon, error)
	// 获取当前商品可以领取的优惠券列表
	GetCouponsByGoodsId(comId uint, goodsId uint) ([]*Coupon, error)

	// 获取活动可以领取优惠券
	GetCouponsByActivity(ComId uint, activityId uint) ([]*Coupon, error)
	// 获取当前订单可以使用的优惠券
	GetUserCouponsByOrder(ComId uint, order *Order) ([]*Coupon, []*Coupon, error)
	// 获取用户领取的优惠券列表
	GetUserCoupons(ComId uint, userId, lastId uint) ([]*UserCoupon, error)
	// 创建用户优惠券
	CreateUserCoupon(comId uint, userId uint, couponId uint) error
	UpdateCachedCouponsByComId(comId uint, lastId uint) ([]*Coupon, error)
	GetUserCouponById(comId, userID, couponID uint) (*Coupon, error)
	// 获取商品页面关联商品的一个活动
	GetGoodsActivity(goodsID uint) (*Activity, error)

	// 获取用户获得的奖品列表
	GetUserAwards(userId, lastId, limit uint) ([]*UserAward, error)

	//获取用户的奖品
	GetUserAward(awardID uint) (*UserAward, error)

	CheckUserAward(userID, awardID uint, num int) error

	// 获取用户领取奖品的历史记录
	GetAwardHistoryByUserID(userId, lastId, limit uint) ([]*AwardHistory, error)
	// 获取员工验票记录
	StaffCheckedAwards(comID, staffID, lastID, limit uint) ([]*AwardHistory, error)
}

//
type ActivityDriver interface {
	// 活动开始 生成用户活动记录
	Begin(activity *Activity, user *User) error
	// 活动状态变化（拼团人数变化 邀请好友）
	Change(userActivity *UserActivity) error
	// 活动成功
	//1. 更新userActivity状态
	//2. 发送奖品 调用SendUserAward()
	//3. 添加奖品历史记录表 调用Record() 添加userAward记录
	Success(userActivity *UserActivity) error
	// 活动失败
	Fail(userActivity *UserActivity) error
	// 发送活动奖品
	SendUserAward(userActivity *UserActivity) error
	// 记录活动奖品
	Record(award *AwardHistory) error
}

// 活动开始 生成用户活动记录
// 使用工厂方法创建活动驱动
type ActivityDriverFactory interface {
	GetActivityDriver(activityType string) (ActivityDriver, error)
}

// 装饰器在原来类基础上加上一个订单变化的监听器
type ActivityDriverWithOrderListener interface {
	ActivityDriver
	observer.Observer
}
