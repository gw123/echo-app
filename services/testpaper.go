package services

import (
	echoapp "github.com/gw123/echo-app"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type TestpaperService struct {
	db *gorm.DB
	//redis *redis.Client

}

func NewTestpaperService(db *gorm.DB) *TestpaperService {
	return &TestpaperService{
		db: db,
	}
}
func (t *TestpaperService) SaveTestpaper(testpaper *echoapp.Testpaper) error {
	t.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.Testpaper{})
	return t.db.Create(testpaper).Error
}
func (t *TestpaperService) GetTestpaperById(id int64) (*echoapp.Testpaper, error) {
	var test = &echoapp.Testpaper{}
	if err := t.db.Table("testpapers1").Where("id=?", id).Find(test).Error; err != nil {
		return nil, errors.Wrap(err, "GetTestPaperById")
	}
	return test, nil
}

func (t *TestpaperService) SaveUserTestAnswer(answer *echoapp.UserAnswer) error {
	t.db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&echoapp.UserAnswer{})
	return t.db.Create(answer).Error
}
