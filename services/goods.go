package services

import (
	"sync"

	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type GoodsService struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewGoodsService(db *gorm.DB) *GoodsService {
	help := &GoodsService{
		db: db,
	}
	return help
}

func (rsv *GoodsService) SaveGoods(goods *echoapp.Goods) error {
	rsv.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Goods{})
	good := &echoapp.Goods{}
	if rsv.db.Table("goods").Where("user_id=? AND name=?", goods.UserId, goods.Name).Find(good).RecordNotFound() {
		return rsv.db.Create(goods).Error
	}
	return nil
}

func (rsv *GoodsService) SaveTags(tags *echoapp.Tags) error {
	rsv.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Tags{})
	tag := &echoapp.Tags{}
	if rsv.db.Table("tags").Where("name=?", tags.Name).Find(tag).RecordNotFound() {
		return rsv.db.Create(tags).Error
	}
	return nil
}

func (rsv *GoodsService) GetGoodsById(c echo.Context, id uint) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	res := rsv.db.Where("ID=?", id).First(goods)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsById")
	}
	//echoapp_util.ExtractEntry(c).Info("goodsID:%d", id)
	return goods, nil
}
func (rsv *GoodsService) GetGoodsByTagId(c echo.Context, tagId uint, from, limit int) ([]*echoapp.Goods, error) {
	var goodslist []*echoapp.Goods

	res := rsv.db.Offset(from*limit).Limit(limit).Where("tag_id=?", tagId).Find(goodslist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsByTagId")
	}
	//echoapp_util.ExtractEntry(c).Info("TagID:%d", tagId)
	return goodslist, nil
}
func (rsv *GoodsService) GetUserPaymentGoods(c echo.Context, userId uint, from int, limit int) ([]*echoapp.Goods, error) {
	var goodslist []*echoapp.Goods
	res := rsv.db.Where("user_id=?", userId).Offset(from * limit).Limit(limit).Find(goodslist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetUserPaymentGoods")
	}
	//echoapp_util.ExtractEntry(c).Info("UserID:%d,from:%d,limit:%d", userId, from, limit)
	return goodslist, nil
}
func (rsv *GoodsService) ModifyGoods(goods *echoapp.Goods) error {
	return rsv.db.Save(goods).Error
}
func (rsv *GoodsService) DeleteGoods(goods *echoapp.Goods) error {
	return rsv.db.Delete(goods).Error
}
func (rsv *GoodsService) GetGoodsList(c echo.Context, from, limit int) ([]*echoapp.GetGoodsOptions, error) {
	var goodsoptions []*echoapp.GetGoodsOptions
	res := rsv.db.Table("goods").Offset(limit * from).Limit(limit).Find(goodsoptions)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsList")
	}
	return goodsoptions, nil
}
func (rsv *GoodsService) GetGoodsByName(name string) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	res := rsv.db.Where("name=?", name).Find(goods)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsByName")
	}
	return goods, nil
}
func (rsv *GoodsService) GetTagsByName(name string) (*echoapp.Tags, error) {
	goods := &echoapp.Tags{}
	res := rsv.db.Where("name=?", name).Find(goods)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetTagsByName")
	}
	return goods, nil
}
