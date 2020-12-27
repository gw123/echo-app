package services

import (
	"math"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type Video struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewVideoService(db *gorm.DB, redis *redis.Client) *Video {
	return &Video{
		db:    db,
		redis: redis,
	}
}

func (aSvr Video) GetVideoList(comId uint, lastId uint, limit int) ([]*echoapp.Video, error) {
	if limit == 0 || limit > 12 {
		limit = 6
	}
	if lastId == 0 {
		lastId = math.MaxUint32
	}
	var VideoList []*echoapp.Video
	if err := aSvr.db.Debug().Where("com_id = ? and status = ?", comId, echoapp.ViodeStatusOnline).
		Where("id < ?", lastId).
		Limit(limit).Find(&VideoList).Error; err != nil {
		return nil, errors.Wrap(err, "getIndexBanner")
	}
	return VideoList, nil
}

func (aSvr Video) GetVideoDetail(id uint) (*echoapp.Video, error) {
	var Video echoapp.Video
	if err := aSvr.db.Where("id = ? and status = ?", id, echoapp.ViodeStatusOnline).First(&Video).Error; err != nil {
		return nil, errors.Wrap(err, "getVideoDetail")
	}
	return &Video, nil
}
