package services

import (
	"github.com/gw123/echo-app/models"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type IntegralServices struct {
	db *gorm.DB
}

func NewIntegralServices() *IntegralServices {
	newIntegral := new(IntegralServices)
	return newIntegral
}

type Op struct {
	Phone          string `json:phone`
	RechargeAmount int    `json:"recharge_amout"`
	CostAmount     int    `json:"cost_amount"`
}

func (t *IntegralServices) Integral(param *Op) error {
	user := &models.User{}
	integral := &models.Integral{}
	err := t.db.Where("phone=?", param.Phone).Find(user)
	if err.Error != nil && !err.RecordNotFound() {
		return errors.Wrap(err.Error, "没有此用户")
	}
	err = t.db.Where("user_id=?", user.ID).Find(integral)
	reinterral := integral.CurrentIntegal + param.RechargeAmount

	integral = &models.Integral{
		UserID:         user.ID,
		CurrentIntegal: reinterral,
	}
	res := t.db.Save(integral)
	if res.Error != nil && t.db.NewRecord(integral) {
		return errors.Wrap(err.Error, "user integral save fail")
	}
	return nil
}
