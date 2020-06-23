package services

import (
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

func (aSvr ActivityService) GetNotifyDetail(id int) (*echoapp.Notify, error) {
	notify := &echoapp.Notify{}
	if err := aSvr.db.Where("id = ?", id).
		First(notify).Error; err != nil {
		return nil, errors.Wrap(err, "db not find notify")
	}
	return notify, nil
}

func (aSvr ActivityService) GetNotifyList(comId int, lastId, limit int) ([]*echoapp.Notify, error) {
	list := []*echoapp.Notify{}
	if err := aSvr.db.Where("com_id = ? and type='notify' ", comId).Order("id desc").Find(&list).Error; err != nil {
		return nil, errors.Wrap(err, "db not find notify list")
	}
	return list, nil
}
