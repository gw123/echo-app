package activity

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/glog"
	"github.com/pkg/errors"
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

	glog.Infof("activity begin")
	if err := a.Dao.CreateUserActivity(userActivity); err != nil {
		return err
	}

	glog.Infof("activity success")
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

// 领取Vip奖品的活动驱动器 这里使用购买vip订单支付完成的事件驱动
type AwardVipActivityDriverListener struct {
	AwardVipActivityDriver
	userSvr echoapp.UserService
}

func NewAwardVipActivityDriverListener(dao Dao, userSvr echoapp.UserService) *AwardVipActivityDriverListener {
	return &AwardVipActivityDriverListener{
		AwardVipActivityDriver: AwardVipActivityDriver{AwardActivityDriver{Dao: dao}},
		userSvr:                userSvr,
	}
}

func (a AwardVipActivityDriverListener) Name() string {
	return "award-vip"
}

func (a AwardVipActivityDriverListener) OnEvent(event interface{}) error {
	order, ok := event.(*echoapp.Order)
	if !ok {
		return errors.Errorf("need order type", order.OrderNo)
	}

	if len(order.GoodsList) == 0 {
		return errors.Errorf("goods list len is 0", order.OrderNo)
	}

	//
	if order.GoodsType == echoapp.GoodsTypeVip {
		glog.Info("微信支付成功回调: 创建vip会员")
		user, err := a.userSvr.GetUserById(int64(order.UserId))
		if err != nil {
			return errors.Wrap(err, "get user at set user vip")
		}

		if err := a.userSvr.SetVipLevel(user, 1); err != nil {
			return errors.Wrap(err, "SetVipLevel")
		}

		act, err := a.Dao.GetGoodsActivity(order.GoodsList[0].GoodsId)
		if err != nil {
			glog.Errorf("GetGoodsActivity %s", err.Error())
			return errors.Wrap(err, "GetGoodsActivity ")
		}

		err = a.Begin(act, user)
		if err != nil {
			glog.Errorf("activity begin err %s", err.Error())
		}
	}
	return nil
}
