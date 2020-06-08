package services

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sync"
	"time"

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
	rsv.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Resource{})
	return rsv.db.Create(resource).Error
}
func (rsv *ResourceService) GetResourceById(c echo.Context, id int64) (*echoapp.Resource, error) {
	resource := &echoapp.Resource{}
	res := rsv.db.Where("ID=?", id).First(resource)
	if res.Error != nil {
		return nil, res.Error
	}
	return resource, nil
}
func (rsv *ResourceService) GetResourcesByTagId(c echo.Context, tagId int64, from int, limit int) ([]*echoapp.Resource, error) {
	var resources []*echoapp.Resource
	res := rsv.db.Where("tag_id=?", tagId).Limit(limit).Offset(from * limit).Find(&resources)
	if res.Error != nil {
		return nil, res.Error
	}
	return resources, nil
}
func (rsv *ResourceService) GetUserPaymentResources(c echo.Context, userId int64, from int, limit int) ([]*echoapp.Resource, error) {
	var resource []*echoapp.Resource
	res := rsv.db.Where("user_id=? AND status=?", userId, "paid").Limit(limit).Offset(from * limit).Find(&resource)
	if res.Error != nil {
		return nil, res.Error
	}
	//echoapp_util.ExtractEntry(c).Info("UserID:%d,from:%d,limit:%d", userId, from, limit)
	return resource, nil
}

func (rsv *ResourceService) GetSelfResources(c echo.Context, userId int64, from int, limit int) ([]*echoapp.Resource, error) {
	var reslist []*echoapp.Resource
	status := "user_upload"
	res := rsv.db.Table("resources").Where("user_id=? AND status=?", userId, status).Limit(limit).Offset(from * limit).Find(&reslist)
	if res.Error != nil {
		return nil, res.Error
	}
	return reslist, nil
}

func (rsv *ResourceService) GetResourceByName(name string) (*echoapp.Resource, error) {
	resource := &echoapp.Resource{}
	res := rsv.db.Where("name=?", name).Find(resource)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "ResourceServicec->GetResourceByName")
	}
	return resource, nil

}
func (rsv *ResourceService) ModifyResource(resource *echoapp.Resource) error {
	return rsv.db.Save(resource).Error
}
func (rsv *ResourceService) DeleteResource(resource *echoapp.Resource) error {
	return rsv.db.Delete(resource).Error
}

func (rsv *ResourceService) GetResourceList(c echo.Context, from, limit int) ([]*echoapp.GetResourceOptions, error) {
	//var options []*echoapp.GetResourceOptions
	options := []*echoapp.GetResourceOptions{}
	res := rsv.db.Table("resources").Offset(limit * from).Limit(limit).Find(&options)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "rsv GetResourceList")
	}
	return options, nil
}
func (rsv *ResourceService) GetResourceByMd5Path(c echo.Context, path string) (*echoapp.Resource, error) {
	resource := &echoapp.Resource{}
	res := rsv.db.Where("path=?", path).Find(resource)
	if res.Error != nil {
		return nil, errors.Wrap(res.Error, "Servicec->GetResourceByMd5")
	}
	return resource, nil
}
func (rsv *ResourceService) UploadFile(c echo.Context, formname, uploadpath string, maxfilesize int64) (map[string]string, error) {
	//r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)
	file, err := c.FormFile(formname)
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	filetype := echoapp_util.GetFileType(file.Filename)
	fullPath := uploadpath + "/" + filetype + "/" + file.Filename
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	size, err := io.Copy(dst, src)
	if err != nil {
		return nil, err
	}
	if size > maxfilesize {
		return nil, errors.New("File size over limit")
	}
	res := map[string]string{
		"filename":   file.Filename,
		"uploadpath": fullPath,
	}
	return res, nil
}
func (rsv *ResourceService) DownloadFile(durl, localpath string) (string, error) {
	uri, err := url.ParseRequestURI(durl)
	if err != nil {
		return "", errors.Wrap(err, "url err")
	}
	filename := path.Base(uri.Path)

	req, err := http.NewRequest("GET", durl, nil)
	if err != nil {
		return "", errors.Wrap(err, "DoRequest->http.NewRequest")
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "DoRequest->Do")
	}
	http.DefaultClient.Timeout = time.Second * 60 //超时设置
	file, err := os.Create(localpath + filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", err
	}
	if resp.Body == nil {
		return "", errors.New("Body is Null")
	}
	defer resp.Body.Close()
	return filename, nil
}
