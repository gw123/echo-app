package models

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"AUTO_INCREMENT"`
	Name      string `gorm:"size:255"`
	Age       int    `gorm:"type:varchar(10)"`
	Phone     string `gorm:"type:varchar(16);unique"`
	Pwd       string `gorm:"type:varchar(16);not null"`
	RoleID    string `gorm:"type:varchar(12)" json:"role_id"`
	Avatar    string `gorm:"type:varchar(56)"`
	Email     string `gorm:"type:varchar(30)"`
	Score     int    `gorm:"type:varchar(16);not null"`
	CreatedAT time.Time
	UpdatedAT time.Time
}
type Params struct {
	UserId int `json:"user_id"`
	Score  int `json:"score"`
}
type UserRoleParam struct {
	Phone string `json:"phone"`
}

type User_Has_Role struct {
	UserID   uint   `json:"user_id"`
	RoleID   uint   `json:"role_id"`
	UserType string `gorm:"size:255" json:"user_type"`
}

type Role struct {
	ID        uint   `gorm:"AUTO_INCREMENT"`
	Key       string `gorm:"size:255"`
	Name      string `gorm:"size:255；unique"`
	GuardName string `json:"guard_name;not null"`
	CreatedAT time.Time
	UpdatedAT time.Time
}

type Permission struct {
	ID        uint
	Name      string
	GuardName string `json:"guard_name"`
	CreatedAT time.Time
	UpdatedAT time.Time
}
type RoleandermissionParam struct {
	Role       string `json:"role"`
	Permission string `json:"permission"`
}
type Role_Has_Permission struct {
	RoleID       uint `json:"role_id"`
	PermissionID uint `json:"permission_id"`
}
