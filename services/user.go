package services

import (
	"sync"

	"github.com/pkg/errors"

	echoapp "github.com/gw123/echo-app"

	echoapp_util "github.com/gw123/echo-app/util"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

type UserService struct {
	db *gorm.DB
	mu sync.Mutex
}

func (uSvr UserService) GetUserByToken(token string) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := uSvr.db.Where("api_token = ?", token).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
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
	if err := res.Error; err != nil && res.RecordNotFound() {
		return nil, err
	}
	echoapp_util.ExtractEntry(ctx).Infof("mobile:%s,pwd:%s", param.Mobile, param.Password)
	return user, nil

}

func (uSvr UserService) GetUserById(ctx echo.Context, userID int) (*echoapp.User, error) {
	user := &echoapp.User{}
	if err := uSvr.db.Where("id = ?", userID).First(user).Error; err != nil {
		return nil, err
	}
	echoapp_util.ExtractEntry(ctx).Infof("userid:%d", userID)
	return user, nil
}

func (uSvr UserService) Save(ctx echo.Context, user *echoapp.User) error {
	return uSvr.db.Save(user).Error
}

func (uSvr UserService) Create(c echo.Context, user *echoapp.RegisterUser) error {
	if err := uSvr.db.Create(user).Error; err != nil && uSvr.db.NewRecord(user) {
		return errors.Wrap(err, "user create fail")
	}
	return nil
}

func (t *UserService) RegisterUser(c echo.Context, param *echoapp.RegisterUser) error {

	err := t.db.Table("users").Where("phone=?", param.Phone)
	if err.Error != nil && err.RecordNotFound() {
		return errors.Wrap(err.Error, "Record has Found")
	}
	echoapp_util.ExtractEntry(c).Infof("mobile:%s,pwd:%s", param.Phone, param.Pwd)
	return t.Create(c, param)
}

func (t *UserService) Addroles(c echo.Context, param *echoapp.Role) error {

	res := t.db.Table("roles").Where("name=?", param.Name)
	if res.Error != nil && res.RecordNotFound() {
		return errors.Wrap(res.Error, "Record has Found")
	}
	err := t.db.Create(param)
	if err.Error != nil && t.db.NewRecord(param) {
		return errors.Wrap(err.Error, "role create failed")
	} else if t.db.NewRecord(param) {
		return errors.New("not NewRecord")
	}
	echoapp_util.ExtractEntry(c).Infof("create role name:%s", param.Name)
	return nil
}
func (t *UserService) AddPermission(c echo.Context, param *echoapp.Permission) error {

	res := t.db.Table("permissions").Where("name=?", param.Name)
	if res.Error != nil && res.RecordNotFound() {
		return errors.Wrap(res.Error, "Record has Found")
	}
	err := t.db.Create(param)
	if err.Error != nil && t.db.NewRecord(param) {
		return errors.Wrap(err.Error, "permiseeion create failed")
	} else if t.db.NewRecord(param) {
		return errors.New("not NewRecord")
	}
	echoapp_util.ExtractEntry(c).Infof("create permission name:%s", param.Name)
	return nil
}

func (t *UserService) RoleHasPermission(c echo.Context, param *echoapp.RoleandPermissionParam) (*echoapp.RoleHasPermission, error) {
	role := &echoapp.Role{}
	permission := &echoapp.Permission{}
	res := t.db.Where("name=?", param.Role).Find(role)
	if res.Error != nil && res.RecordNotFound() {
		return nil, errors.Wrap(res.Error, "Role record has Found")
	} else if res.RecordNotFound() {
		return nil, errors.New("Role Record has Found")
	}
	res = t.db.Where("name=?", param.Permission)
	if res.Error != nil && res.RecordNotFound() {
		return nil, errors.Wrap(res.Error, "Permission record has Found")
	} else if res.RecordNotFound() {
		return nil, errors.New("Permission Record has Found")
	}
	rolehaspermission := &echoapp.RoleHasPermission{
		RoleID:       role.ID,
		PermissionID: permission.ID,
	}
	err := t.db.Create(rolehaspermission)
	if err.Error != nil && t.db.NewRecord(param) {
		return nil, errors.Wrap(err.Error, "rolehasper create failed")
	} else if t.db.NewRecord(param) {
		return nil, errors.New("not NewRecord")
	}
	return rolehaspermission, nil
}
