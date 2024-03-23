package activity

import (
	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
)

type Dao interface {
	CreateUserActivity(userActivity *echoapp.UserActivity) error
	UpdateUserActivity(userActivity *echoapp.UserActivity) error
	AddUserAward(comID, userId uint, goodsId uint, num uint) error
	AddAwardHistory(awardHistory *echoapp.AwardHistory) error
	GetActivity(ActivityID uint) (*echoapp.Activity, error)
	RecordAwardHistory(awardHistory *echoapp.AwardHistory) error
	GetGoodsActivity(goodsID uint) (*echoapp.Activity, error)
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

func (a ActivityDao) AddUserAward(comID, userId uint, goodsId uint, num uint) error {
	userAward := &echoapp.UserAward{}
	if err := a.db.Where("user_id = ? and goods_id = ?", userId, goodsId).First(userAward).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		} else {
			//
			userAward.UserID = userId
			userAward.GoodsID = goodsId
			userAward.ComID = comID
		}
	}
	userAward.Num += num
	return a.db.Debug().Save(userAward).Error
}

func (a ActivityDao) AddAwardHistory(awardHistory *echoapp.AwardHistory) error {
	return a.db.Debug().Save(awardHistory).Error
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

//获取商品详情页 某个商品的关联活动
func (aSvr ActivityDao) GetGoodsActivity(goodsId uint) (*echoapp.Activity, error) {
	var goodsActivity echoapp.GoodsActivity
	var activity echoapp.Activity
	var err error
	//优先获取单独给这个商品配置的活动
	err = aSvr.db.Debug().Where("goods_id = ? and status = ? ", goodsId, echoapp.GoodsActivityStatusOnline).
		First(&goodsActivity).Error
	if err != nil {
		return nil, err
	}

	err = aSvr.db.Debug().Where("id = ?", goodsActivity.ActivityID).
		First(&activity).Error
	if err != nil {
		return nil, err
	}
	return &activity, nil
}
