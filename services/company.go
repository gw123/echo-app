package services

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v7"
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	//redis 相关的key
	RedisCompanyKey = "Company:%d"
)

func FormatCompanyRedisKey(comId uint) string {
	return fmt.Sprintf(RedisCompanyKey, comId)
}

type CompanyService struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewCompanyService(db *gorm.DB, redis *redis.Client) *CompanyService {
	return &CompanyService{
		db:    db,
		redis: redis,
	}
}

func (c CompanyService) GetCompanyById(comId uint) (*echoapp.Company, error) {
	company := &echoapp.Company{}
	if err := c.db.Where("id = ?", comId).First(company).Error; err != nil {
		return nil, errors.Wrap(err, "db query")
	}
	company.SmsChannels = make(map[string]*echoapp.SmsChannel)
	company.WxTemplateTypes = make(map[string]*echoapp.TemplateType)
	// sms channel
	var channels []echoapp.SmsChannel
	c.db.Where("com_id = ?", comId).Find(&channels)
	for _, sch := range channels {
		company.SmsChannels[sch.Type] = &sch
	}
	// wxTplMessage
	var wxTpls []*echoapp.TemplateType
	c.db.Where("com_id = ?", comId).Find(&wxTpls)
	for _, tpl := range wxTpls {
		company.WxTemplateTypes[tpl.Type] = tpl
	}

	return company, nil
}

func (c CompanyService) GetCompanyList(offsetId uint, limit int) ([]*echoapp.Company, error) {
	list := []*echoapp.Company{}
	if err := c.db.Where("id > ?", offsetId).
		Order("id asc").Limit(limit).
		Find(&list).Error; err != nil {
		return nil, errors.Wrap(err, "db err")
	}
	return list, nil
}

func (c CompanyService) GetCachedCompanyById(comId uint) (*echoapp.Company, error) {
	user := &echoapp.Company{}
	data, err := c.redis.Get(FormatCompanyRedisKey(comId)).Result()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(data), user); err != nil {
		return nil, err
	}
	return user, nil
}

func (c CompanyService) UpdateCachedCompany(company *echoapp.Company) (err error) {
	data, err := json.Marshal(company)
	if err != nil {
		return errors.Wrap(err, "redis set")
	}
	if err = c.redis.Set(FormatCompanyRedisKey(company.Id), data, 0).Err(); err != nil {
		return errors.Wrap(err, "redis set")
	}
	return err
}

func (c *CompanyService) GetQuickNav(comId uint) ([]*echoapp.Nav, error) {
	var navs []*echoapp.Nav
	if err := c.db.Debug().Where("com_id = ?", comId).Find(&navs).Error; err != nil {
		return nil, errors.Wrap(err, "query nav")
	}
	return navs, nil
}
