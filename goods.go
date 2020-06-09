package echoapp

import (
	"time"
<<<<<<< HEAD

	"github.com/labstack/echo"
)

type Goods struct {
	//gorm.Model
	ID         int64 `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	UserId     int64      `json:"user_id"`
	ComId      int        `json:"com_id"`
	TagStr     string     `josn:"tags"`
	Tags       string     `josn:"-" gorm:"-"`
	Name       string     `json:"name"`
	Price      float32    `json:"price"`
	Body       string     `json:"body"`
	RealPrice  float32    `json:"real_price"`
	GoodType   string     `json:"good_type"`
	Status     string     `json:"status"`
	SmallCover string     ` json:"small_cover"`
	Covers     string     `gorm:"type:varchar(2048)" json:"covers"`
	Pages      int        `json:"pages"`
}

type Tags struct {
	ID        int64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `json:"name"`
	ComId     string `json:"com_id"`
}
type GetGoodsOptions struct {
	Id         uint   `json:"id"`
	ComId      int    `json:"com_id"`
	TagId      uint   `josn:"tag_id"`
	Name       string `json:"name"`
	Price      string `json:"price"`
	Body       string `json:"body"`
	RealPrice  int    `json:"real_price"`
	GoodType   string `json:"good_type"`
	Status     string `json:"status"`
	Cover      string `gorm:"type:varchar(2048)" json:"covers"`
	SmallCover string `json:"small_cover"`
}
type GoodsService interface {
	SaveTags(tags *Tags) error
	//保存上传的资源到数据库
	SaveGoods(goods *Goods) error
	//通过资源ID查找资源
	GetGoodsById(c echo.Context, id uint) (*Goods, error)
	//通过tagId查找资源
	GetGoodsByTagId(c echo.Context, tagId uint, from, limit int) ([]*Goods, error)

	//用户购买的资源
	GetUserPaymentGoods(c echo.Context, userId uint, from int, limit int) ([]*Goods, error)

	ModifyGoods(goods *Goods) error

	DeleteGoods(goods *Goods) error
	GetGoodsByName(name string) (*Goods, error)
	GetTagsByName(name string) (*Tags, error)
	//查看资源文件 ，每页有 limit 条数据
	GetGoodsList(c echo.Context, from, limit int) ([]*GetGoodsOptions, error)
=======
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
>>>>>>> develop
}
