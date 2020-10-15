package activity

import (
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
)

type Dao interface {
	CreateUserActivity(userActivity *echoapp.UserActivity) error
	UpdateUserActivity(userActivity *echoapp.UserActivity) error
	AddUserAward(userId uint, goodsId uint, num uint) error
	AddAwardHistory(awardHistory *echoapp.AwardHistory) error
	GetActivity(ActivityID uint) (*echoapp.Activity, error)
	RecordAwardHistory(awardHistory *echoapp.AwardHistory) error
}

type ActivityDao struct {
	db    *gorm.DB
	cache *redis.Client
}

func NewActivityDao(db *gorm.DB, cache *redis.Client) *ActivityDao {
	return &ActivityDao{db: db, cache: cache}
}

func (a ActivityDao) CreateUserActivity(userActivity *echoapp.UserActivity) error {
	return a.db.Create(userActivity).Error
}

func (a ActivityDao) UpdateUserActivity(userActivity *echoapp.UserActivity) error {
	return a.db.Save(userActivity).Error
}

func (a ActivityDao) AddUserAward(userId uint, goodsId uint, num uint) error {
	userAward := &echoapp.UserAward{}
	if err := a.db.Where("user_id = ? and goods_id = ?", userId, goodsId).First(userAward).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		} else {
			//
			userAward.UserID = userId
			userAward.GoodsID = goodsId
		}
	}

	userAward.Num += num
	return a.db.Save(userAward).Error
}

func (a ActivityDao) AddAwardHistory(awardHistory *echoapp.AwardHistory) error {
	return a.db.Save(awardHistory).Error
}

func (a ActivityDao) GetActivity(activityID uint) (*echoapp.Activity, error) {
	activity := &echoapp.Activity{}
	if err := a.db.Where("id = ?", activityID).First(activity).Error; err != nil {
		return nil, err
	}
	return activity, nil
}

func (a ActivityDao) RecordAwardHistory(awardHistory *echoapp.AwardHistory) error {
	return a.db.Save(awardHistory).Error
}
