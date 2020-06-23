package echoapp

import (
	"encoding/json"
	"time"
)

type Testpaper struct {
	ID          int64 `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time  `sql:"index"`
	Rid         int         `json:"rid"`
	Content     string      `json:"content" gorm:"size:1024"`
	QuestionArr []*Question `json:"questions" gorm:"-"`
	QuestionStr string      `json:"-" gorm:"column:questions;size:1024"`
}

func (*Testpaper) TableName() string {
	return "testpapers1"
}

type Question struct {
	ID             int64    `gorm:"primary_key"`
	Content        string   `json:"content" gorm:"size:256;not null"`
	Type           string   `josn:"type"`
	Options        []string `json:"options" gorm:"-"`
	OptionsStr     string   `json:"-" gorm:"column:options;size:1024 "`
	StandardAnswer string   `json:"standard_answer" gorm:"standard_answer"`
}

type UserAnswer struct {
	ID          int64 `gorm:"primary_key"`
	CreatedAt   time.Time
	UserId      int64
	TestpaperId int64
	//	QA          map[int]string `json:"user_answers" gorm:"-"`
	QAStr string `gorm:"column:user_answers;size:1024" json:"-"`
}

func (t *Testpaper) BeforeCreate() (err error) {
	qstr, err := json.Marshal(t.QuestionArr)
	if err != nil {
		return err
	}
	t.QuestionStr = string(qstr)
	return
}
func (t *Testpaper) AfterFind() error {
	err := json.Unmarshal([]byte(t.QuestionStr), &t.QuestionArr)
	if err != nil {
		return err
	}
	return nil
}

// func (t *UserAnswer) AfterFind() error {
// 	err := json.Unmarshal([]byte(t.QAStr), &t.QA)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (t *UserAnswer) BeforeCreate() (err error) {
// 	qstr, err := json.Marshal(t.QA)
// 	if err != nil {
// 		return err
// 	}
// 	t.OptionsStr = string(qstr)
// 	return
// }

type TestpaperService interface {
	SaveTestpaper(*Testpaper) error
	GetTestpaperById(id int64) (*Testpaper, error)
	SaveUserTestAnswer(answer *UserAnswer) error
}
