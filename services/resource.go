package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"

	echoapp "github.com/gw123/echo-app"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
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
	rsv.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Resource{}, &echoapp.Tag{})
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
func (rsv *ResourceService) GetResourcesByTagId(c echo.Context, tagID uint, from int, limit int) ([]*echoapp.Resource, error) {
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
	echoapp_util.ExtractEntry(c).Info("UserID:%d,from:%d,limit:%d", userId, from, limit)
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
func (rsv *ResourceService) Md5SumFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "ReadFile", err
	}
	rawMd5 := md5.Sum(data)
	sign := fmt.Sprintf("%x", rawMd5)
	return sign, nil
}

func (rsv *ResourceService) GetResourceByPath(path string) (*echoapp.Resource, error) {
	resource := &echoapp.Resource{}
	res := rsv.db.Where("name=?", path).Find(resource)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "Servicec->GetResourcePath")
	}

	return resource, nil

}
func (rsv *ResourceService) ModifyResource(resource *echoapp.Resource) error {
	return rsv.db.Save(resource).Error
}
func (rsv *ResourceService) DeleteResource(resource *echoapp.Resource) error {
	return rsv.db.Delete(resource).Error
}

func (rsv *ResourceService) UploadFile(c echo.Context, formname, uploadpath string) (string, error) {
	file, err := c.FormFile(formname)
	if err != nil {
		return "", err
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	var fullPath string
	if rsv.GetFileType(file.Filename) == "pdf" {
		fullPath = uploadpath + "/pdf/" + file.Filename
	} else if rsv.GetFileType(file.Filename) == "pptx" {
		fullPath = uploadpath + "/ppt/" + file.Filename
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	echoapp_util.ExtractEntry(c).Info("uploadpath:%s,fullpath:%s", uploadpath, fullPath)
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	return fullPath, nil
}
func (rsv *ResourceService) GetFileType(filename string) string {
	s := path.Ext(filename)
	return s[1:]
}
func (rsv *ResourceService) GetResourceList(c echo.Context, from, limit int) ([]echoapp.ResFilePathParam, error) {
	var filepath []echoapp.ResFilePathParam
	//res := rsv.db.Table("resources").Find(&filepath)
	res := rsv.db.Table("resources").Offset(limit * from).Limit(limit).Find(&filepath)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "GetPDFList")
	}
	return filepath, nil
}
