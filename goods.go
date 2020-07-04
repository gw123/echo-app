package echoapp

import (
	"encoding/json"
	"time"

	"github.com/gw123/glog"
)

const (
	GoodsTypeGoods   = "goods"
	GoodsTypeTicket  = "ticket"
	GoodsTypeRoom    = "room"
	GoodsTypeCombine = "combine"
)

type GoodsBrief struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	UserID      uint      `json:"user_id"`
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
	CoversStr   string    `json:"-" gorm:"column:covers"`
	Covers      []string  `gorm:"-" json:"covers"`
	Desc        string    `json:"desc"`
	GoodsType   string    `json:"goods_type" gorm:"goods_type" `
}

func (g *GoodsBrief) AfterFind() error {
	err := json.Unmarshal([]byte(g.CoversStr), &g.Covers)
	if err != nil {
		glog.Errorf("Goods AfterFind")
	}
	return nil
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

type GoodsTag struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	ComId     uint       `json:"com_id"`
	Name      string     `json:"name"`
	Icon      string     `json:"icon"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
}

// type BannerBrief struct {
// 	ID        uint       `gorm:"primary_key" json:"id"`
// 	CreatedAt time.Time  `json:"created_at"`
// 	UpdatedAt time.Time  `json:"updated_at"`
// 	DeletedAt *time.Time `sql:"index"`
// }

//type Banner struct {
//	Body string `json:"body"`
//}

type GoodsService interface {
	AddGoodsPv(goodsId int) error
	GetGoodsById(goodsId int) (*Goods, error)
	GetGoodsByName(name string) (*Goods, error)
	GetGoodsList(comId, lastId uint, limit int) ([]*GoodsBrief, error)
	GetTagByName(name string) (*GoodsTag, error)
	SaveTag(tag *GoodsTag) error
	GetRecommendGoodsList(comId, lastId uint, limit int) ([]*GoodsBrief, error)
	Save(goods *Goods) error
	GetGoodsByCode(code string) (*Goods, error)
	UpdateCachedGoods(goods *Goods) (err error)
	GetTagGoodsList(comID uint, tagID int, lastID uint, limit int) ([]*GoodsBrief, error)
}
