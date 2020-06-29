package services

import (
	"encoding/json"
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func NewActivityService(db *gorm.DB, redis *redis.Client) *ActivityService {
	return &ActivityService{
		db:    db,
		redis: redis,
	}
}

type ActivityService struct {
	db    *gorm.DB
	redis *redis.Client
}

func (aSvr ActivityService) GetActivityList(comId uint, lastId uint, limit int) ([]*echoapp.Activity, error) {
	if limit == 0 || limit > 12 {
		limit = 6
	}
	var activityList []*echoapp.Activity
	if err := aSvr.db.Where("com_id = ? and status ='online'", comId).
		Where("id < ?", lastId).
		Limit(limit).Find(&activityList).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return activityList, nil
}

func (aSvr ActivityService) GetActivityDetail(id uint) (*echoapp.Activity, error) {
	var activity echoapp.Activity
	if err := aSvr.db.Where("id = ? and status ='online'", id).First(&activity).Error; err != nil {
		return nil, errors.Wrap(err, "getActivityDetail")
	}
	return &activity, nil
}

func (aSvr ActivityService) GetBannerList(comId uint, position string, limit int) ([]*echoapp.BannerBrief, error) {
	if position == "" {
		position = "index"
	}
	if limit == 0 || limit > 12 {
		limit = 6
	}
	var banners []*echoapp.BannerBrief

	if err := aSvr.db.Where("com_id = ? and status ='online'", comId).
		Where("position = ?", position).
		Limit(limit).Find(&banners).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}

	return banners, nil
}

func (aSvr *ActivityService) GetCachedBannerList(comId uint, position string) ([]*echoapp.BannerBrief, error) {
	bannerList, err := func() ([]*echoapp.BannerBrief, error) {
		var bannerList []*echoapp.BannerBrief
		data, err := aSvr.redis.Get(echoapp.FormatBannerListRedisKey(comId, position)).Result()
		if err != nil {
			return bannerList, errors.Wrap(err, "GetCachedBannerList->redis.Get")
		}
		if data == "" {
			return bannerList, nil
		}
		err = json.Unmarshal([]byte(data), &bannerList)
		if err != nil {
			return bannerList, errors.Wrap(err, "GetCachedBannerList->json.Unmarshal")
		}
		return bannerList, nil
	}()

	if err != nil {
		return bannerList, errors.Wrap(err, "GetCachedBannerList")
	}

	if len(bannerList) == 0 {
		bannerList, err = aSvr.GetBannerList(comId, position, 6)
		if err != nil {
			return bannerList, errors.Wrap(err, "GetCachedBannerList->GetBannerList")
		}
	}

	return bannerList, nil
}

func (aSvr *ActivityService) UpdateCachedBannerList(comId uint, position string) error {
	bannerList, err := aSvr.GetBannerList(comId, position, 6)
	if err != nil {
		return err
	}
	data, err := json.Marshal(bannerList)
	err = aSvr.redis.Set(echoapp.FormatBannerListRedisKey(comId, position), data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (aSvr ActivityService) GetIndexBanner(comId uint) ([]*echoapp.BannerBrief, error) {
	panic("implement me")
}

func (aSvr ActivityService) AddActivityPv(goodsId uint) error {
	panic("implement me")
}

func (aSvr ActivityService) GetNotifyDetail(id int) (*echoapp.Notify, error) {
	notify := &echoapp.Notify{}
	if err := aSvr.db.Where("id = ?", id).
		First(notify).Error; err != nil {
		return nil, errors.Wrap(err, "db not find notify")
	}
	return notify, nil
}

func (aSvr ActivityService) GetNotifyList(comId uint, lastId, limit int) ([]*echoapp.Notify, error) {
	list := []*echoapp.Notify{}
	if err := aSvr.db.Where("com_id = ? and type='notify' ", comId).Order("id desc").Find(&list).Error; err != nil {
		return nil, errors.Wrap(err, "db not find notify list")
	}
	return list, nil
}
