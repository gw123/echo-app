package echoapp

import (
	"time"
)

type GoodsBrief struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ComId       int       `json:"com_id"`
	Price       float32   `json:"price"`
	RealPrice   float32   `json:"real_price"`
	Num         int       `json:"num"`
	SaleNum     int       `json:"sale_num"`
	Status      string    `json:"status"`
	ExpressType string    `json:"express_type"`
	Express     int       `json:"express"`
	Tags        string    `json:"tags"`
	Name        string    `json:"name"`
	SmallCover  string    `json:"small_cover"`
	Covers      string    `json:"covers"`
	Desc        string    `json:"desc"`
}

func (*GoodsBrief) TableName() string {
	return "goods"
}

type Goods struct {
	GoodsBrief
	Body      string `json:"body"`
	Infos     string `json:"infos"`
	GoodsType string `json:"goods_type"`
}

type BannerBrief struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
	Title     string     `json:"title"`
	Cover     string     `json:"cover"`
	Visit     int        `json:"visit"`
	EndAt     time.Time  `json:"end_at"`
	Type      string     `json:"type"`
	GoodsId   int        `json:"goods_id"`
}

func (*BannerBrief) TableName() string {
	return "activies"
}

type Banner struct {
	Body string `json:"body"`
}

type GoodsService interface {
	GetIndexBanner(comId int) ([]*BannerBrief, error)
	GetActivityById(id int) (*Banner, error)
	AddActivityPv(goodsId int) error
	AddGoodsPv(goodsId int) error
	GetGoodsById(goodsId int) (*Goods, error)
	GetGoodsList(comId, lastId, limit int) ([]*GoodsBrief, error)
	GetRecommendGoodsList(comId, lastId, limit int) ([]*GoodsBrief, error)
	Save(goods *Goods) error
	GetGoodsByCode(code string) (*Goods, error)
	UpdateCachedGoods(goods *Goods) (err error)
}
