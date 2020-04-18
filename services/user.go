package services

import (
	"sync"

	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"
	"github.com/gw123/echo-app/models"
	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserService struct {
	db *gorm.DB
	mu sync.Mutex
}

func NewUserService(dbs *gorm.DB) *UserService {
	return &UserService{
		db: dbs,
	}
}

func (uSvr UserService) AddScore(ctx echo.Context, user *echoapp.User, amount int) error {
	user.Score += amount
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,增加积分: %d", user.Id, amount)
	return uSvr.Save(ctx, user)
}

func (uSvr UserService) SubScore(ctx echo.Context, user *echoapp.User, amount int) error {
	user.Score -= amount
	echoapp_util.ExtractEntry(ctx).Infof("UserId: %d ,消耗积分: %d", user.Id, amount)
	return uSvr.Save(ctx, user)
}

func (uSvr UserService) Login(ctx echo.Context, param *echoapp.LoginParam) (*echoapp.User, error) {
	user := &echoapp.User{}
	res := uSvr.db.Where("phone=? AND pwd=? ", param.Mobile, param.Password).Find(user)
	//panic("implement me")
	if err := res.Error; err != nil && res.RecordNotFound() {
		return nil, err
	}
	echoapp_util.ExtractEntry(ctx).Infof("mobile:%s,pwd:%s", param.Mobile, param.Password)
	return user, nil

}

func (uSvr UserService) GetUserById(ctx echo.Context, userId int) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := uSvr.db.Where("id = ?", userId).First(user).Error; err != nil {
		return nil, err
	}
	echoapp_util.ExtractEntry(ctx).Infof("userid:%d", userId)
	return user, nil
}

func (uSvr UserService) Save(ctx echo.Context, user *echoapp.User) error {
	return uSvr.db.Save(user).Error
}

func (uSvr UserService) CreateUser(c echo.Context, user *echoapp.RegisterUser) error {
	if err := uSvr.db.Create(user).Error; err != nil && uSvr.db.NewRecord(user) {
		return errors.Wrap(err, "user create fail")
	}
	return nil
}

func (t *UserService) RegisterUser(c echo.Context, param *echoapp.RegisterUser) error {

	err := t.db.Table("users").Where("phone=?", param.Phone)
	if err.Error != nil && err.RecordNotFound() {
		return errors.Wrap(err.Error, "Record has Found")
	} else if err.RecordNotFound() {
		return errors.New("Record has Found")
	}
	echoapp_util.ExtractEntry(c).Infof("mobile:%s,pwd:%s", param.Phone, param.Pwd)
	return t.CreateUser(c, param)
}

func (t *UserService) Addroles(c echo.Context, param *models.Role) error {

	res := t.db.Table("roles").Where("name=?", param.Name)
	if res.Error != nil && res.RecordNotFound() {
		return errors.Wrap(res.Error, "Record has Found")
	} else if res.RecordNotFound() {
		return errors.New("Record has Found")
	}
	err := t.db.Create(param)
	if err.Error != nil && t.db.NewRecord(param) {
		return errors.Wrap(err.Error, "role  create failed")
	} else if t.db.NewRecord(param) {
		return errors.New("not NewRecord")
	}
	echoapp_util.ExtractEntry(c).Infof("role name:%s", param.Name)
	return nil
}
func (t *UserService) AddPermission(c echo.Context, param *models.Permission) error {

	res := t.db.Table("permissions").Where("name=?", param.Name)
	if res.Error != nil && res.RecordNotFound() {
		return errors.Wrap(res.Error, "Record has Found")
	} else if res.RecordNotFound() {
		return errors.New("Record has Found")
	}
	err := t.db.Create(param)
	if err.Error != nil && t.db.NewRecord(param) {
		return errors.Wrap(err.Error, "role  create failed")
	} else if t.db.NewRecord(param) {
		return errors.New("not NewRecord")
	}
	echoapp_util.ExtractEntry(c).Infof("permission name:%s", param.Name)
	return nil
}

/*
func (t *UserService) RoleHasPermission(c echo.Context,param *models.RoleandermissionParam) (*models.Role_Has_Permission,error) {
	res:=t.db.Table("roles").Where("name=?",param.Role)
	if res.Error != nil && !res.RecordNotFound() {
		return errors.Wrap(res.Error,"Record has Found")
	} else if !res.RecordNotFound() {
		return errors.New("Record has Found")
	}
	err := t.db.Create(param)
	if err.Error != nil && t.db.NewRecord(param) {
		return errors.Wrap(err.Error, "role  create failed")
	}else if t.db.NewRecord(param){
		return  errors.New("not NewRecord")
	}
	return nil
}
*/
