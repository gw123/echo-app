package services

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"sync"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/labstack/echo"
	"github.com/pkg/errors"

	"github.com/jinzhu/gorm"
)

type ResourceService struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewResourceService(db *gorm.DB) *ResourceService {
	help := &ResourceService{
		db: db,
	}
	return help
}

func (rsv *ResourceService) SaveResource(resource *echoapp.Resource) error {
	rsv.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Resource{})
	return rsv.db.Create(resource).Error
}
func (rsv *ResourceService) GetResourceById(c echo.Context, id uint) (*echoapp.Resource, error) {
	resource := &echoapp.Resource{}
	res := rsv.db.Where("ID=?", id).First(resource)
	if res.Error != nil {
		return nil, res.Error
	}
	echoapp_util.ExtractEntry(c).Info("ResourceID:%d", id)
	return resource, nil
}
func (rsv *ResourceService) GetResourcesByTagID(c echo.Context, tagID uint, from int, limit int) ([]*echoapp.Resource, error) {
	tag := &echoapp.Tag{}
	res := rsv.db.Where("ID=?", tagID).Find(tag)
	if res.Error != nil {
		return nil, res.Error
	}
	echoapp_util.ExtractEntry(c).Info("TagID:%d", tagID)
	var resource []*echoapp.Resource
	//[] resource := &echoapp.Resource{}
	res = rsv.db.Where("title=?", tag.Title).Limit(limit).Offset(from * limit).Find(resource)
	if res.Error != nil {
		return nil, res.Error
	}
	return resource, nil
}
func (rsv *ResourceService) GetUserPaymentResources(c echo.Context, userId uint, from int, limit int) ([]*echoapp.Resource, error) {
	var resource []*echoapp.Resource
	res := rsv.db.Where("user_id=?", userId).Limit(limit).Offset(from * limit).Find(resource)
	if res.Error != nil {
		return nil, res.Error
	}
	echoapp_util.ExtractEntry(c).Info("ID:%d", userId)
	return resource, nil
}

func (rsv *ResourceService) GetSelfResources(c echo.Context, userId uint, from int, limit int) ([]*echoapp.Resource, error) {
	return nil, nil
}
func (rsv *ResourceService) GetMd5String(path string) string {
	h := md5.New()
	h.Write([]byte(path))
	return hex.EncodeToString(h.Sum(nil))
}
func (rsv *ResourceService) Md5SumFile(file string) (value [md5.Size]byte, err error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	value = md5.Sum(data)
	return value, nil
}

func (rsv *ResourceService) GetResourceByPath(path string) (*echoapp.Resource, error) {
	resource := &echoapp.Resource{}
	res := rsv.db.Where("path=?", path).Find(resource)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GetResourcePath-path in DB")
	}
	return resource, nil

}
func (rsv *ResourceService) ModifyResource(resource *echoapp.Resource) error {
	return rsv.db.Save(resource).Error
}
func (rsv *ResourceService) DeleteResource(resource *echoapp.Resource) error {
	return rsv.db.Delete(resource).Error
}
