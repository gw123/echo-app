package services

import (
	"encoding/json"
	"fmt"
	"github.com/gw123/glog"
	"github.com/olivere/elastic/v7"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/components"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	//redis 相关的key
	RedisGoodsKey = "Goods:%d"
)

func FormatGoodsRedisKey(goodsId uint) string {
	return fmt.Sprintf(RedisGoodsKey, goodsId)
}

type GoodsService struct {
	db    *gorm.DB
	redis *redis.Client
	jws   *components.JwsHelper
	es    *elastic.Client
}

func NewGoodsService(db *gorm.DB, redis *redis.Client, es *elastic.Client) *GoodsService {
	return &GoodsService{
		db:    db,
		redis: redis,
		es:    es,
	}
}

func (gSvr *GoodsService) GetGoodsByName(name string) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	res := gSvr.db.Where("name=?", name).Find(goods)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsByName")
	}
	return goods, nil
}

func (gSvr *GoodsService) SaveTag(tag *echoapp.GoodsTag) error {
	gSvr.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.GoodsTag{})
	tmptag := &echoapp.GoodsTag{}
	if gSvr.db.Table("tags").Where("name=?", tag.Name).Find(tmptag).RecordNotFound() {
		return gSvr.db.Create(tag).Error
	}
	return nil
}

func (gSvr *GoodsService) GetGoodsTags(comID uint) ([]*echoapp.GoodsTag, error) {
	var goodsTags []*echoapp.GoodsTag
	if err := gSvr.db.Where("com_id=?", comID).Find(&goodsTags).Error; err != nil {
		return nil, errors.Wrap(err, "query tags")
	}
	return goodsTags, nil
}

func (gSvr *GoodsService) GetIndexBanner(comId int) ([]*echoapp.BannerBrief, error) {
	banners := []*echoapp.BannerBrief{}
	if err := gSvr.db.Where("com_id = ? and status ='online'", comId).Find(&banners).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return banners, nil
}

func (gSvr *GoodsService) GetActivityById(id int) (*echoapp.Banner, error) {
	banner := echoapp.Banner{}
	if err := gSvr.db.Where("id = ?", id).Find(&banner).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return &banner, nil
}

func (gSvr *GoodsService) AddActivityPv(activityId int) error {
	banner := echoapp.Banner{}
	if err := gSvr.db.Model(&banner).Where("id = ?", activityId).
		Update("set visit=visit+1").Error; err != nil {
		return errors.Wrap(err, "getIndexBanner")
	}
	return nil
}

func (gSvr *GoodsService) AddGoodsPv(goodsId int) error {
	glog.Infof("AddGoodsPv %d", goodsId)
	goods := echoapp.Goods{}
	if err := gSvr.db.Debug().Model(&goods).Where("id = ?", goodsId).
		Update("pv", gorm.Expr("pv+1")).Error; err != nil {
		return errors.Wrap(err, "AddGoodsPv")
	}
	return nil
}

func (gSvr *GoodsService) GetGoodsByCode(code string) (*echoapp.Goods, error) {
	panic("implement me")
}

func (gSvr *GoodsService) GetGoodsList(comId, lastId uint, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := gSvr.db.Where("com_id = ? and id > ?", comId, lastId).
		Where("status = 'publish'").
		Order("id desc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (gSvr *GoodsService) GetGoodsListByKeyword(comId uint, keyword string, lastId uint, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := gSvr.db.Where("com_id = ? and id > ?", comId, lastId).
		Where("status = 'publish'").
		Where("name like ?", "%"+keyword+"%").
		Order("id desc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (gSvr *GoodsService) GetTagGoodsList(comID uint, tagID int, lastID uint, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := gSvr.db.Where("com_id = ? and id > ?", comID, lastID).
		Where("status = 'publish'").
		Where("match (tags) against (? in boolean mode)", tagID).
		Order("id desc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (gSvr *GoodsService) GetRecommendGoodsList(comId, lastId uint, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := gSvr.db.Where("com_id = ? and id > ?", comId, lastId).
		Where("recommend_lvl > 0 and status = 'publish'").
		Order("recommend_lvl desc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (gSvr *GoodsService) GetGoodsById(goodsId uint) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	if err := gSvr.db.Where(" id = ?", goodsId).First(goods).Error; err != nil {
		return nil, err
	}
	return goods, nil
}

func (gSvr *GoodsService) GetCachedGoodsById(goodsId uint) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	data, err := gSvr.redis.Get(FormatGoodsRedisKey(goodsId)).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), goods); err != nil {
		return nil, err
	}
	return goods, nil
}

func (gSvr *GoodsService) UpdateCachedGoods(goods *echoapp.Goods) (err error) {
	//r := time.Duration(rand.Int63n(180))
	data, err := json.Marshal(goods)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	//fmt.Println(string(data))
	err = gSvr.redis.Set(FormatGoodsRedisKey(goods.ID), data, 0).
		Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return err
}

//基础方法自动更新cache
func (gSvr *GoodsService) Save(goods *echoapp.Goods) error {
	if err := gSvr.db.Save(goods).Error; err != nil {
		return errors.Wrap(err, "db save goods")
	}
	return gSvr.UpdateCachedGoods(goods)
}

func (gSvr *GoodsService) GetGoodsByTagId(tagId uint, from, limit int) ([]*echoapp.Goods, error) {
	var goodslist []*echoapp.Goods

	res := gSvr.db.Offset(from*limit).Limit(limit).Where("tag_id=?", tagId).Find(goodslist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsByTagId")
	}

	return goodslist, nil
}

func (gSvr *GoodsService) GetUserPaymentGoods(userId uint, from int, limit int) ([]*echoapp.Goods, error) {
	var goodslist []*echoapp.Goods
	res := gSvr.db.Where("user_id=?", userId).Offset(from * limit).Limit(limit).Find(goodslist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetUserPaymentGoods")
	}
	return goodslist, nil
}

func (gSvr *GoodsService) DeleteGoods(goods *echoapp.Goods) error {
	return gSvr.db.Delete(goods).Error
}

func (gSvr *GoodsService) GetTagByName(name string) (*echoapp.GoodsTag, error) {
	goods := &echoapp.GoodsTag{}
	res := gSvr.db.Where("name=?", name).Find(goods)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetTagsByName")
	}
	return goods, nil
}

func (gSvr *GoodsService) GetCartGoodsList(comID uint, userID uint) (*echoapp.Cart, error) {
	cart := &echoapp.Cart{}
	if err := gSvr.db.Debug().Where("com_id = ? and  user_id = ?", comID, userID).
		First(cart).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			cart.Content = []*echoapp.CartGoodsItem{}
			cart = &echoapp.Cart{
				ComId:  comID,
				UserID: userID,
			}
			return cart, nil
		}
		return nil, err
	}
	if cart.Content == nil {
		cart.Content = []*echoapp.CartGoodsItem{}
	}
	return cart, nil
}

func (gSvr *GoodsService) DelCartGoods(comID uint, userID uint, goodsID uint, skuID uint) error {
	item := &echoapp.CartGoodsItem{GoodsId: goodsID, SkuID: skuID}
	return gSvr.updateCartGoods(comID, userID, item, "del")
}

func (gSvr *GoodsService) AddCartGoods(comID uint, userID uint, goodsItem *echoapp.CartGoodsItem) error {
	return gSvr.updateCartGoods(comID, userID, goodsItem, "add")
}

func (gSvr *GoodsService) updateCartGoods(comID uint, userID uint, goodsItem *echoapp.CartGoodsItem, action string) error {
	if action != "del" {
		if err := gSvr.IsValidCartGoods(goodsItem); err != nil {
			return errors.Wrap(err, "商品校验失败")
		}
	}

	cart, err := gSvr.GetCartGoodsList(comID, userID)
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return errors.Wrap(err, "GetCartGoodsList")
	}

	if err != nil && gorm.IsRecordNotFoundError(err) {
		cart = &echoapp.Cart{
			ComId:  comID,
			UserID: userID,
		}
	}

	isNew := true
	for index, item := range cart.Content {
		if item.GoodsId == goodsItem.GoodsId && item.SkuID == goodsItem.SkuID {
			isNew = false
			if goodsItem.Num == 0 {
				cart.Content = append(cart.Content[0:index], cart.Content[index+1:]...)
				break
			}
			if action == "add" {
				item.Num += 1
			} else {
				item.Num = goodsItem.Num
			}
			break
		}
	}

	if isNew {
		if goodsItem.Num == 0 {
			return errors.New("删除商品不存在购物车")
		}
		cart.Content = append(cart.Content, goodsItem)
	}

	if err := gSvr.db.Save(cart).Error; err != nil {
		return errors.Wrap(err, "SaveCart")
	}
	return nil
}

func (gSvr *GoodsService) UpdateCartGoods(comID uint, userID uint, goodsItem *echoapp.CartGoodsItem) error {
	return gSvr.updateCartGoods(comID, userID, goodsItem, "update")
}

func (gSvr *GoodsService) IsValidCartGoods(item *echoapp.CartGoodsItem) error {
	//todo sku商品
	if item.SkuID == 0 {
		goods, err := gSvr.GetGoodsById(item.GoodsId)
		if err != nil {
			return errors.Wrapf(err, "获取商品失败 %d", item.GoodsId)
		}
		if item.Name != goods.Name {
			return errors.New(item.Name + "商品名称发生变动")
		}
		if item.Price != goods.Price {
			return errors.New(goods.Name + "原价对比失败")
		}
		if item.RealPrice != goods.RealPrice {
			return errors.New(goods.Name + "商品价格发生变动")
		}
		if uint(goods.Num) < item.Num {
			return errors.New(goods.Name + "商品数量不足")
		}
	} else {
		if len(item.LabelCombine) == 0 {
			return errors.New("请选择商品属性后下单")
		}

		goods, err := gSvr.GetSkuById(item.GoodsId, item.LabelCombine)

		if err != nil {
			return errors.Wrap(err, "获取商品失败")
		}

		if item.SkuName != goods.SkuName {
			return errors.New(item.Name + "商品名称发生变动")
		}

		if item.RealPrice != goods.RealPrice {
			glog.Errorf("cart price:%f, realPrice:%f", item.RealPrice, goods.RealPrice)
			return errors.New(item.SkuName + "商品价格发生变动")
		}

		if uint(goods.Num) < item.Num {
			return errors.New(goods.SkuName + "商品数量不足")
		}
	}
	return nil
}

func (gSvr *GoodsService) IsValidCartGoodsList(itemList []*echoapp.CartGoodsItem) error {
	for _, item := range itemList {
		if err := gSvr.IsValidCartGoods(item); err != nil {
			return err
		}
	}
	return nil
}

func (gSvr *GoodsService) ClearCart(comID uint, userID uint) error {
	cart, err := gSvr.GetCartGoodsList(comID, userID)
	if gorm.IsRecordNotFoundError(err) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "clearCart")
	}
	if err := gSvr.db.Debug().Table("carts").Delete(cart).Error; err != nil {
		return errors.Wrap(err, "delete")
	}
	return nil
}

func (gSvr *GoodsService) GetSkuById(goodsId uint, labelCombine map[string]string) (*echoapp.Sku, error) {
	goods, err := gSvr.GetGoodsById(goodsId)
	if err != nil {
		return nil, err
	}

	if goods.Status != echoapp.GoodsStatusPublish {
		return nil, errors.New("商品已经下架")
	}

	for _, sku := range goods.Skus {
		if sku.IsLabelCombine(labelCombine) {
			return sku, nil
		}
	}
	return nil, errors.New("not found")
}
