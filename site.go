package echoapp

import (
	"context"
	"fmt"
	"time"
)

const (
	BannerPositionScroll1  = "scroll_1"
	BannerPositionAD       = "ad_1"
	BannerPositionCategory = "category"
)

func FormatBannerListRedisKey(comId uint, position string) string {
	return fmt.Sprintf(RedisBannerListKey, comId, position)
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

type BannerBrief struct {
	Id         uint   `gorm:"primary_key" json:"id"`
	Cover      string `json:"cover"`
	Href       string `json:"href"`
	Type       string `json:"type"`
	TargetId   int    `json:"-"`
	ComId      uint   `json:"-"`
	Background string `json:"background"` //背景颜色
}

func (b *BannerBrief) AfterFind() error {
	if b.Type == "goods" {
		b.Href = fmt.Sprintf("/pages/product/product?id=%d&com_id=%d", b.TargetId, b.ComId)
	} else if b.Type == "activity" {
		b.Href = fmt.Sprintf("/pages/activity/detail?id=%d&com_id=%d", b.TargetId, b.ComId)
	}
	return nil
}

func (*BannerBrief) TableName() string {
	return "banners"
}

type Banner struct {
	BannerBrief
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Title     string     `json:"title"`
	Visit     int        `json:"visit"`
	EndAt     time.Time  `json:"end_at"`
}

// 每个页面的banner图
type PageBanners struct {
	Scroll1    []*BannerBrief `json:"scroll_1"`
	ADs        []*BannerBrief `json:"ads"`
	Categories []*BannerBrief `json:"categories"`
}

type SiteService interface {
	GetNotifyList(comId uint, lastId, limit int) ([]*Notify, error)
	GetNotifyDetail(id int) (*Notify, error)
	//
	GetBannerList(ctx context.Context, comId uint, page, position string, limit int) ([]*BannerBrief, error)
	UpdateCachedBannerList(ctx context.Context, comId uint, position string) error
	GetCachedBannerList(ctx context.Context, comId uint, position string) ([]*BannerBrief, error)
	GetIndexBanner(comId uint) ([]*BannerBrief, error)
	GetIndexPageBanners(ctx context.Context, comId uint) (*PageBanners, error)
	//GetActivityById(id int) (*Banner, error)
}
