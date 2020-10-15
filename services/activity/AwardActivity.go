package activity

import (
	echoapp "github.com/gw123/echo-app"
)

// 领取奖品的活动通用
type AwardActivityDriver struct {
	Dao Dao
}

func NewAwardActivityDriver(dao Dao) *AwardActivityDriver {
	return &AwardActivityDriver{Dao: dao}
}

func (a AwardActivityDriver) Begin(activity *echoapp.Activity, user *echoapp.User) error {
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
	return nil
}

func (a AwardActivityDriver) Change(userActivity *echoapp.UserActivity) error {
	if err := a.Dao.UpdateUserActivity(userActivity); err != nil {
		return err
	}
	return nil
}

func (a AwardActivityDriver) Success(userActivity *echoapp.UserActivity) error {
	userActivity.Status = echoapp.UserActivityStatusSuccess
	if err := a.Dao.UpdateUserActivity(userActivity); err != nil {
		return err
	}
	// 发送奖品
	a.SendUserAward(userActivity)

	return nil
}

func (a AwardActivityDriver) Fail(userActivity *echoapp.UserActivity) error {
	userActivity.Status = echoapp.UserActivityStatusFail
	if err := a.Dao.UpdateUserActivity(userActivity); err != nil {
		return err
	}
	return nil
}

func (a AwardActivityDriver) SendUserAward(userActivity *echoapp.UserActivity) error {
	activity, err := a.Dao.GetActivity(userActivity.ID)
	if err != nil {
		return err
	}

	for _, award := range activity.Rewards {
		if err := a.Dao.AddUserAward(userActivity.UserID, award.GoodsId, award.Num); err != nil {
			return err
		}
		awardHistory := &echoapp.AwardHistory{
			UserID:         userActivity.UserID,
			GoodsID:        award.GoodsId,
			UserActivityID: userActivity.ID,
			ActivityID:     userActivity.ActivityID,
			Method:         echoapp.AwardHistoryTypeActivity,
			Num:            award.Num,
		}
		if err := a.Record(awardHistory); err != nil {
			return err
		}
	}
	return nil
}

func (a AwardActivityDriver) Record(award *echoapp.AwardHistory) error {
	return a.Dao.RecordAwardHistory(award)
}
