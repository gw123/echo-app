package services

import (
	"encoding/json"
	"fmt"

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
}

func NewGoodsService(db *gorm.DB, redis *redis.Client) *GoodsService {
	return &GoodsService{
		db:    db,
		redis: redis,
	}
}

func (u *GoodsService) GetGoodsByName(name string) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	res := u.db.Where("name=?", name).Find(goods)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsByName")
	}
	return goods, nil
}

func (u *GoodsService) SaveTag(tag *echoapp.GoodsTag) error {
	u.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.GoodsTag{})
	tmptag := &echoapp.GoodsTag{}
	if u.db.Table("tags").Where("name=?", tag.Name).Find(tmptag).RecordNotFound() {
		return u.db.Create(tag).Error
	}
	return nil
}

func (u *GoodsService) GetIndexBanner(comId int) ([]*echoapp.BannerBrief, error) {
	banners := []*echoapp.BannerBrief{}
	if err := u.db.Where("com_id = ? and status ='online'", comId).Find(&banners).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return banners, nil
}

func (u *GoodsService) GetActivityById(id int) (*echoapp.Banner, error) {
	banner := echoapp.Banner{}
	if err := u.db.Where("id = ?", id).Find(&banner).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return &banner, nil
}

func (u *GoodsService) AddActivityPv(activityId int) error {
	banner := echoapp.Banner{}
	if err := u.db.Model(&banner).Where("id = ?", activityId).
		Update("set visit=visit+1").Error; err != nil {
		return errors.Wrap(err, "getIndexBanner")
	}
	return nil
}

func (u *GoodsService) AddGoodsPv(goodsId int) error {
	goods := echoapp.Goods{}
	if err := u.db.Model(&goods).Where("id = ?", goodsId).
		Update("set visit=visit+1").Error; err != nil {
		return errors.Wrap(err, "getIndexBanner")
	}
	return nil
}

func (u *GoodsService) GetGoodsByCode(code string) (*echoapp.Goods, error) {
	panic("implement me")
}

func (u *GoodsService) GetGoodsList(comId, lastId uint, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := u.db.Where("com_id = ? and id > ?", comId, lastId).
		Where("status = 'publish'").
		Order("id desc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (u *GoodsService) GetTagGoodsList(comID uint, tagID int, lastID uint, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := u.db.Where("com_id = ? and id > ?", comID, lastID).
		Where("status = 'publish'").
		Where("match (tags) against (? in boolean mode)", tagID).
		Order("id desc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (u *GoodsService) GetRecommendGoodsList(comId, lastId uint, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := u.db.Where("com_id = ? and id > ?", comId, lastId).
		Where("recommend_lvl > 0 and status = 'publish'").
		Order("recommend_lvl desc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (u *GoodsService) GetCachedGoodsById(goodsId uint) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	data, err := u.redis.Get(FormatGoodsRedisKey(goodsId)).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), goods); err != nil {
		return nil, err
	}
	return goods, nil
}

func (u *GoodsService) UpdateCachedGoods(goods *echoapp.Goods) (err error) {
	//r := time.Duration(rand.Int63n(180))
	data, err := json.Marshal(goods)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	//fmt.Println(string(data))
	err = u.redis.Set(FormatGoodsRedisKey(goods.ID), data, 0).
		Err()
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	return err
}

//基础方法自动更新cache
func (u *GoodsService) Save(goods *echoapp.Goods) error {
	if err := u.db.Save(goods).Error; err != nil {
		return errors.Wrap(err, "db save goods")
	}
	return u.UpdateCachedGoods(goods)
}

func (rsv *GoodsService) GetGoodsByTagId(tagId uint, from, limit int) ([]*echoapp.Goods, error) {
	var goodslist []*echoapp.Goods

	res := rsv.db.Offset(from*limit).Limit(limit).Where("tag_id=?", tagId).Find(goodslist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetGoodsByTagId")
	}

	return goodslist, nil
}

func (rsv *GoodsService) GetUserPaymentGoods(userId uint, from int, limit int) ([]*echoapp.Goods, error) {
	var goodslist []*echoapp.Goods
	res := rsv.db.Where("user_id=?", userId).Offset(from * limit).Limit(limit).Find(goodslist)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetUserPaymentGoods")
	}
	return goodslist, nil
}

func (rsv *GoodsService) DeleteGoods(goods *echoapp.Goods) error {
	return rsv.db.Delete(goods).Error
}

func (rsv *GoodsService) GetTagByName(name string) (*echoapp.GoodsTag, error) {
	goods := &echoapp.GoodsTag{}
	res := rsv.db.Where("name=?", name).Find(goods)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GoodsService->GetTagsByName")
	}
	return goods, nil
}

func (u *GoodsService) GetGoodsById(goodsId int) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	if err := u.db.Where(" id = ?", goodsId).First(goods).Error; err != nil {
		return nil, err
	}
	return goods, nil
}

func (u *GoodsService) GetCartGoodsList(comID uint, userID uint) (*echoapp.Cart, error) {
	cart := &echoapp.Cart{}
	if err := u.db.Where("com_id = ? and  user_id = ?", comID, userID).
		First(cart).Error; err != nil {
		return nil, err
	}
	return cart, nil
}

func (u *GoodsService) DelCartGoods(comID uint, userID uint, goodsID uint, skuID uint) error {
	item := &echoapp.CartGoodsItem{GoodsId: goodsID, SkuID: skuID}
	return u.UpdateCartGoods(comID, userID, item)
}

func (u *GoodsService) AddCartGoods(comID uint, userID uint, goodsItem *echoapp.CartGoodsItem) error {
	return u.UpdateCartGoods(comID, userID, goodsItem)
}

func (u *GoodsService) UpdateCartGoods(comID uint, userID uint, goodsItem *echoapp.CartGoodsItem) error {
	cart, err := u.GetCartGoodsList(comID, userID)
	if err != nil {
		return errors.Wrap(err, "GetCartGoodsList")
	}
	isNew := true
	for index, item := range cart.Content {
		if item.GoodsId == goodsItem.GoodsId && item.SkuID == goodsItem.SkuID {
			isNew = false
			if goodsItem.Num == 0 {
				cart.Content = append(cart.Content[0:index], cart.Content[index:]...)
				break
			}
			item.Num = goodsItem.Num
			break
		}
	}

	if isNew {
		if goodsItem.Num == 0 {
			return errors.New("删除商品不存在购物车")
		}
		cart.Content = append(cart.Content, goodsItem)
	}

	if err := u.db.Save(cart).Error; err != nil {
		return errors.Wrap(err, "SaveCart")
	}
	return nil
}
