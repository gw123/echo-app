package activity

import (
	echoapp "github.com/gw123/echo-app"
)

// 领取Vip奖品的活动通用
type AwardVipActivityDriver struct {
	AwardActivityDriver
}

func NewAwardVipActivityDriver(dao Dao) *AwardVipActivityDriver {
	return &AwardVipActivityDriver{
		AwardActivityDriver: AwardActivityDriver{Dao: dao},
	}
}

func (a AwardVipActivityDriver) Begin(activity *echoapp.Activity, user *echoapp.User) error {
	userActivity := &echoapp.UserActivity{
		UserID:     uint(user.Id),
		ActivityID: activity.Id,
		Status:     echoapp.UserActivityStatusIng,
		Metadata:   "",
		Log:        "",
	}
	if err := a.Dao.CreateUserActivity(userActivity); err != nil {
		return err
	}

	// vip 是充值后立马成功
	if err := a.Success(userActivity); err != nil {
		return err
	}

	return nil
}

func (a AwardVipActivityDriver) Record(award *echoapp.AwardHistory) error {
	award.Method = echoapp.AwardHistoryTypeVIP
	return a.Dao.RecordAwardHistory(award)
}
