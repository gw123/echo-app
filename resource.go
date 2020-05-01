package echoapp

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

/*
type Category struct {
	gorm.Model
	Title string `json:"title"`
}



type TagGroup struct {
	gorm.Model
	TagId   uint `json:"tag_id"`
	GroupId uint `json:"group_id"`
}

func (TagGroup) TableName() string {
	return "tag_group"
}

type Group struct {
	gorm.Model
	Title      string    `gorm:"column:title" json:"title"`
	UserId     uint      `gorm:"column:user_id"json:"user_id"`
	Covers     []string  `gorm:"-" json:"covers"`
	CoversStr  string    `gorm:"column:covers" json:"covers_str"`
	Desc       string    `gorm:"column:desc" json:"desc"`
	CategoryId uint      `gorm:"column:category_id" json:"category_id"`
	TagIds     []uint    `gorm:"-" json:"tags"`
	Display    string    `json:"display"`
	Chapters   []Chapter `json:"chapters"`
}

func (g *Group) BeforeSave() (err error) {
	var str []byte
	if str, err = json.Marshal(g.Covers); err != nil {
		return err
	}
	g.CoversStr = string(str)
	return
}

type Chapter struct {
	gorm.Model
	Title     string     `json:"title"`
	GroupId   uint       `json:"group_id"`
	ParentId  uint       `json:"parent_id"`
	Resources []Resource `json:"resources"`
}
*/
type Tag struct {
	gorm.Model
	Title string `json:"title"`
}
type Resource struct {
	gorm.Model
	//Title  string `json:"title"`
	UserId uint `json:"user_id"`
	// GroupId   uint       `json:"group_id"`
	// ChapterId uint       `json:"chapter_id"`
	//	Covers    string     `json:"covers"`
	Type string `json:"type"`
	//	Article   *Article   `json:"article" gorm:"ForeignKey:rid"`
	//	Testpaper *Testpaper `json:"testpaper" gorm:"ForeignKey:rid"`
	Path      string `json:"path" gorm:"unique"`
	Price     string `json:"price"`
	Status    string `json:"status"`
	Download  string `json:"download"`
	Name      string `json:"name"`
	Privilege string `json:"privilege"`
	TagId     int    `json:"tag_id"`
}

/*
type Article struct {
	gorm.Model
	Rid     int    `json:"rid"`
	Content string `json:"content"`
}

type Testpaper struct {
	gorm.Model
	Rid     int    `json:"rid"`
	Content string `json:"content"`
}*/
type ResFilePathParam struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type ResourceService interface {
	//保存上传的资源到数据库
	SaveResource(resource *Resource) error
	//通过资源ID查找资源
	GetResourceById(c echo.Context, id uint) (*Resource, error)
	//通过tagId查找资源
	GetResourcesByTagId(c echo.Context, tagId uint, from int, limit int) ([]*Resource, error)
	//用户上传的资源
	GetSelfResources(c echo.Context, userId uint, from int, limit int) ([]*Resource, error)
	//用户购买的资源
	GetUserPaymentResources(c echo.Context, userId uint, from int, limit int) ([]*Resource, error)
	//MD5文件path
	GetMd5String(path string) string
	//Md5文件内容
	Md5SumFile(file string) (string, error)
	//更改资源
	ModifyResource(resource *Resource) error
	//删除资源
	DeleteResource(resource *Resource) error
	//通过PATH查找资源
	GetResourceByPath(path string) (*Resource, error)
	//本地上传文件到服务端
	UploadFile(c echo.Context, formname, uploadpath string) (string, error)
	//查看资源文件 ，每页有 limit 条数据
	GetResourceList(c echo.Context, from, limit int) ([]ResFilePathParam, error)
}
