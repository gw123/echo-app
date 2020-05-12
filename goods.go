package echoapp

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Goods struct {
	//gorm.Model
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	UserId     uint       `json:"user_id"`
	ComId      int        `json:"com_id"`
	TagStr     string     `josn:"tags"`
	Tags       string     `josn:"-" gorm:"-"`
	Name       string     `json:"name"`
	Price      float32    `json:"price"`
	Body       string     `json:"body"`
	RealPrice  float32    `json:"real_price"`
	GoodType   string     `json:"good_type"`
	Status     string     `json:"status"`
	SmallCover string     `gorm:"type:varchar(2048)" json:"small_cover"`
	Covers     string     `json:"covers"`
	Pages      int        `json:"pages"`
}

type Tags struct {
	gorm.Model
	GoodsId uint `goods_id`
	//TagId   uint   `tag_id`
	Name  string `json:"name"`
	ComId string `json:"com_id"`
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
	SmallCover string `json:"small_cover"`
	Cover      string `json:"cover"`
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

	//查看资源文件 ，每页有 limit 条数据
	GetGoodsList(c echo.Context, from, limit int) ([]*GetGoodsOptions, error)
}
