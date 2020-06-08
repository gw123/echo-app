package echoapp

import (
	"time"
)

/***
`type` varchar(198) COLLATE utf8mb4_unicode_ci NOT NULL,
`target_id` int(11) NOT NULL,
`user_id` int(11) NOT NULL,
`parent_id` int(11) NOT NULL DEFAULT '0',
`created_at` datetime NOT NULL,
*/

type Comment struct {
	ID        int32     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at" gorm:"index:created_at"`
	Ip        string    `json:"ip" form:"ip" gorm:"size:16"`
	Mac       string    `json:"mac" form:"mac" gorm:"size:40"`
	ClientId  string    `json:"client_id" form:"client_id" gorm:"size:40"`
	Content   string    `json:"content" form:"content" gorm:"size:4096"`
	Type      string    `json:"type" form:"type" gorm:"size:10;index:type_targetId" `
	TargetId  int32     `json:"target_id" form:"target_id" gorm:"index:type_targetId"`
	UserId    int32     `json:"user_id" form:"user_id" gorm:"index:userId"`
	ParentId  int32     `json:"parent_id" form:"parent_id" gorm:"index:parentId"`
	Image     string    `json:"image"`
	Reply     string    `json:"reply"`
}

func (u Comment) TableName() string {
	return "comments"
}

type ScoreOption struct {
	Service     string
	logistics   string
	Description string
}
type ImageOption struct {
	ID     int
	status string
}

type CommentService interface {
	SaveComment(comment *Comment) error
	GetCommentById(id int) (*Comment, error)
	GetCommentByTargetId(targetId int32, from, limit int) (*Comment, error)
	GetCommentList(targetId int32) ([]*Comment, error)
	UpdateComment(comment *Comment) error
	DeleteComment(comment *Comment) error
}
