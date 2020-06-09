package echoapp

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type Resource struct {
	gorm.Model
	UserId     int64      `json:"user_id"`
	Type       string     `json:"type"`
	Testpaper  *Testpaper `json:"testpaper" gorm:"ForeignKey:rid"`
	Path       string     `json:"path" gorm:"unique"`
	Status     string     `json:"status"`
	Download   string     `json:"download"`
	Name       string     `json:"name"`
	Privilege  string     `json:"privilege"`
	TagId      int64      `json:"tag_id"`
	GoodsId    int64      `json:"goods_id"`
	Covers     string     `gorm:"type:varchar(2048)" json:"covers"`
	SmallCover string     `json:"small_cover"`
	Pages      int        `json:"pages"`
}
type Testpaper struct {
	gorm.Model
	Rid     int    `json:"rid"`
	Content string `json:"content"`
}
type GetResourceOptions struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type File struct {
	gorm.Model
	Rid     int    `json:"rid"`
	Content string `json:"content"`
}

type ResourceService interface {
	//保存上传的资源到数据库
	SaveResource(resource *Resource) error
	//删除资源
	DeleteResource(resource *Resource) error
	//更改资源
	ModifyResource(resource *Resource) error

	//通过资源ID查找资源
	GetResourceById(c echo.Context, id int64) (*Resource, error)
	//通过tagId查找资源
	GetResourcesByTagId(c echo.Context, tagId int64, from int, limit int) ([]*Resource, error)
	//用户上传的资源
	GetSelfResources(c echo.Context, userId int64, from int, limit int) ([]*Resource, error)
	//用户购买的资源
	GetUserPaymentResources(c echo.Context, userId int64, from int, limit int) ([]*Resource, error)
	//通过文件name（加后缀）查找资源
	GetResourceByName(path string) (*Resource, error)
	//通过文件name MD5（加后缀）查找资源
	GetResourceByMd5Path(c echo.Context, file string) (*Resource, error)
	//查看资源文件 ，每页有 limit 条数据
	GetResourceList(c echo.Context, from, limit int) ([]*GetResourceOptions, error)

	//本地上传文件到服务端
	UploadFile(c echo.Context, formname, uploadpath string, maxfilesize int64) (map[string]string, error)
	//资源下载
	DownloadFile(durl, localpath string) (string, error)
	//Md5文件内容
	//Md5SumFile(file string) (string, error)
}
