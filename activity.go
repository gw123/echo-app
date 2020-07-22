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
	BannerBrief
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Title     string     `json:"title"`
	Visit     int        `json:"visit"`
	EndAt     time.Time  `json:"end_at"`
	Type      string     `json:"type"`
	GoodsId   int        `json:"goods_id"`
	ComId     uint       `json:"com_id"`
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

type Coupon struct {
	Id         uint   `gorm:"primary_key" json:"id"`
	Type       string `json:"type"`
	Aid        string `json:"aid"` //关联活动id
	Cover      string `json:"cover"`
	Name       string
	Href       string `json:"href"`
	CreatedAt  time.Time
	MinConsume uint   `json:"min_consume"` // 最低消费
	Range      []uint `json:"range"`       //商品适用范围
	Total      uint
	UsedTotal  uint
	ExpireAt   uint
	StartAt    uint
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
	GetBannerList(comId uint, position string, limit int) ([]*BannerBrief, error)
	UpdateCachedBannerList(comId uint, position string) error
	GetCachedBannerList(comId uint, position string) ([]*BannerBrief, error)
	GetIndexBanner(comId uint) ([]*BannerBrief, error)
	//GetActivityById(id int) (*Banner, error)
	AddActivityPv(goodsId uint) error
	//
	GetActivityList(comId uint, lastId uint, limit int) ([]*Activity, error)
	GetActivityDetail(id uint) (*Activity, error)
}
