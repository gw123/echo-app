package echoapp

import (
	"encoding/json"
	"time"

	"github.com/gw123/glog"
	"github.com/jinzhu/gorm"
)

const (
	GoodsTypeGoods   = "goods"
	GoodsTypeTicket  = "ticket"
	GoodsTypeRoom    = "room"
	GoodsTypeCombine = "combine"
	//订单中多种商品混合可能出现
	GoodsTypeMix = "mix"

	GoodsStatusPublish = "publish"
	GoodsStatusOffline = "offline"
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
	Pv          int       `json:"pv"`
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
	//下面這樣的会影响 goods结果导致goods 为空
	GoodsType string `json:"goods_type" gorm:"goods_type" `
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

//子商品
type Sku struct {
	ID           uint              `json:"id"`
	SkuName      string            `json:"sku_name"`
	Price        float32           `json:"price"`
	RealPrice    float32           `json:"real_price"`
	Cover        string            `json:"cover"`
	Num          uint              `json:"num"`
	LabelCombine map[string]string `json:"label_combine"` //sku组合 通过组合确定sku j
}

//判断组合的属性和当前商品是否匹配
func (sku *Sku) IsLabelCombine(labels map[string]string) bool {
	var count, num = 0, len(sku.LabelCombine)
	for key, val := range labels {
		if skuVal, ok := sku.LabelCombine[key]; ok && val == skuVal {
			count++
			continue
		} else {
			return false
		}
	}
	if num != count {
		return false
	}
	return true
}

//子商品同一个属性列表, 商品展示选择属性时候使用
type GoodsLabel struct {
	Label   string   `json:"label"`
	Type    string   `json:"type"`
	Options []string `json:"options"`
}

type GoodsLabels []GoodsLabel

type Goods struct {
	GoodsBrief
	SkusStr      string        `gorm:"column:skus" json:"-"`
	SkuLabelsStr string        `json:"-" gorm:"column:sku_labels"`
	Skus         []*Sku        `gorm:"column:-" json:"skus"`
	SkuLabels    []*GoodsLabel `json:"sku_labels" gorm:"column:-"`
	Body         string        `json:"body"`
	Infos        string        `json:"infos"`
	Service      string        `json:"service"`
}

func (goods *Goods) AfterFind() error {
	if goods.SkuLabelsStr != "" {
		if err := json.Unmarshal([]byte(goods.SkuLabelsStr), &goods.SkuLabels)
			err != nil {
			glog.Error(goods.SkuLabelsStr + "skuiLabel str---" + err.Error())
			return err
		}
	}
	if goods.SkusStr != "" {
		if err := json.Unmarshal([]byte(goods.SkusStr), &goods.Skus); err != nil {
			glog.Error(goods.SkusStr + "---" + err.Error())
			return err
		}
	}
	return goods.GoodsBrief.AfterFind()
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

// Cart
type Cart struct {
	gorm.Model
	ComId      uint             `json:"com_id"`
	UserID     uint             `json:"user_id"`
	Status     string           `json:"status"`
	Content    []*CartGoodsItem `json:"content" gorm:"-"`
	ContentStr string           `json:"-" gorm:"column:content"`
}

func (Cart) TableName() string {
	return "carts"
}

func (c *Cart) AfterFind() error {
	err := json.Unmarshal([]byte(c.ContentStr), &c.Content)
	if err != nil {
		glog.Errorf("Cart AfterFind %s", err)
	}
	return nil
}

func (c *Cart) BeforeSave() error {
	data, err := json.Marshal(c.Content)
	if err != nil {
		glog.Errorf("Cart BeforeSave %s", err)
	}
	c.ContentStr = string(data)
	return nil
}

type CartGoodsItem struct {
	GoodsId      uint              `json:"goods_id"`
	Name         string            `json:"name"`
	SkuName      string            `json:"sku_name"`
	SkuID        uint              `json:"sku_id"`
	Num          uint              `json:"num"`
	Price        float32           `json:"price"`
	RealPrice    float32           `json:"real_price"`
	Cover        string            `json:"cover"`
	LabelCombine map[string]string `json:"label_combine"`
}

type GoodsService interface {
	AddGoodsPv(goodsId int) error
	GetGoodsById(goodsId uint) (*Goods, error)
	GetGoodsByName(name string) (*Goods, error)
	GetGoodsList(comId, lastId uint, limit int) ([]*GoodsBrief, error)
	GetTagByName(name string) (*GoodsTag, error)
	SaveTag(tag *GoodsTag) error
	GetRecommendGoodsList(comId, lastId uint, limit int) ([]*GoodsBrief, error)
	GetGoodsListByKeyword(comId uint, keyword string, lastId uint, limit int) ([]*GoodsBrief, error)
	Save(goods *Goods) error
	GetGoodsByCode(code string) (*Goods, error)
	UpdateCachedGoods(goods *Goods) (err error)
	GetTagGoodsList(comID uint, tagID int, lastID uint, limit int) ([]*GoodsBrief, error)
	GetCachedGoodsById(goodsId uint) (*Goods, error)
	//购物车
	GetCartGoodsList(comID uint, userID uint) (*Cart, error)
	DelCartGoods(comID uint, userID uint, goodsID uint, skuID uint) error
	AddCartGoods(comID uint, userID uint, goods *CartGoodsItem) error
	UpdateCartGoods(comID uint, userID uint, goods *CartGoodsItem) error
	ClearCart(id uint, u uint) error
	IsValidCartGoods(item *CartGoodsItem) error
	IsValidCartGoodsList(itemList []*CartGoodsItem) error
	GetGoodsTags(id uint) ([]*GoodsTag, error)
}
