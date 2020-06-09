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

func (u *GoodsService) GetGoodsByName(name string) (*echoapp.Goods, error) {
	panic("implement me")
}

func (u *GoodsService) GetTagByName(name string) (*echoapp.GoodsTag, error) {
	panic("implement me")
}

func (u *GoodsService) SaveTag(tag *echoapp.GoodsTag) error {
	panic("implement me")
}

func NewGoodsService(db *gorm.DB, redis *redis.Client) *GoodsService {
	return &GoodsService{
		db:    db,
		redis: redis,
	}
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
		Update("set visit=visit+1").Error;
		err != nil {
		return errors.Wrap(err, "getIndexBanner")
	}
	return nil
}

func (u *GoodsService) AddGoodsPv(goodsId int) error {
	goods := echoapp.Goods{}
	if err := u.db.Model(&goods).Where("id = ?", goodsId).
		Update("set visit=visit+1").Error;
		err != nil {
		return errors.Wrap(err, "getIndexBanner")
	}
	return nil
}

func (u *GoodsService) GetGoodsByCode(code string) (*echoapp.Goods, error) {
	panic("implement me")
}

func (u *GoodsService) GetGoodsByToken(token string) (*echoapp.Goods, error) {
	panic("implement me")
}

func (u *GoodsService) GetGoodsList(comId, lastId, limit int) ([]*echoapp.GoodsBrief, error) {
	var goodsList []*echoapp.GoodsBrief
	if err := u.db.Where("com_id = ? and id > ?", comId, lastId).
		Where("status = 'publish'").
		Order("id asc").Limit(limit).
		Find(&goodsList).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return goodsList, nil
}

func (u *GoodsService) GetRecommendGoodsList(comId, lastId, limit int) ([]*echoapp.GoodsBrief, error) {
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

func (u *GoodsService) GetGoodsById(goodsId int) (*echoapp.Goods, error) {
	goods := &echoapp.Goods{}
	if err := u.db.Where(" id = ?", goodsId).First(goods).Error; err != nil {
		return nil, err
	}

	return goods, nil
}

//基础方法自动更新cache
func (u *GoodsService) Save(goods *echoapp.Goods) error {
	if err := u.db.Save(goods).Error; err != nil {
		return errors.Wrap(err, "db save goods")
	}
	return u.UpdateCachedGoods(goods)
}
