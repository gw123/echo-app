package services

import (
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

type ResourceService struct {
	db    *gorm.DB
	redis *redis.Client
	opt   *echoapp.ResourceServerOption
}

func NewResourceService(db *gorm.DB, redis *redis.Client, opt *echoapp.ResourceServerOption) echoapp.ResourceService {
	fmt.Println(*opt)
	return &ResourceService{
		db:    db,
		redis: redis,
		opt:   opt,
	}
}

func (r ResourceService) GetUploadToken(comId int) (string, error) {
	putPolicy := storage.PutPolicy{
		Scope:            r.opt.BucketName,
		CallbackURL:      r.opt.CallbackUrl,
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)","type":"$(x:type)"}`,
		ReturnURL:        "",
		ReturnBody:       `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","name":"$(x:name)","type":"$(x:type)"}`,
		CallbackBodyType: "application/json",
	}
	mac := qbox.NewMac(r.opt.AccessKey, r.opt.SecretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken, nil
}

func (r ResourceService) SaveResource(resource echoapp.Resource) error {
	panic("implement me")
}

func (r ResourceService) GetResourceById(id int) (*echoapp.Resource, error) {
	panic("implement me")
}

func (r ResourceService) GetFileById(id int) (*echoapp.File, error) {
	panic("implement me")
}
func (r ResourceService) GetResourcesByTagID(tagID int, lastId, limit int) ([]*echoapp.Resource, error) {
	panic("implement me")
}

func (r ResourceService) GetSelfResources(userId int, lastId, limit int) ([]*echoapp.Resource, error) {
	panic("implement me")
}

func (r ResourceService) GetUserPaymentResources(userId int, from, limit int) ([]*echoapp.Resource, error) {
	panic("implement me")
}
