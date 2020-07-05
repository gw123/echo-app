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
	Type      string     `json:"type"`
	GoodsId   int        `json:"goods_id"`
}

type BannerBrief struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Cover string `json:"cover"`
	Href  string `json:"href"`
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
