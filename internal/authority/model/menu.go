package model

import (
	"gorm.io/gorm"
	"time"
)

//菜单
type Menu struct {
	ID uint `gorm:"primarykey"`

	ParentId    uint   `gorm:"null;default:0"` //	父级id
	Path        string `gorm:"not null"`       //	路径
	Name        string `gorm:"not null"`       //	菜单名字
	Description string `gorm:"null"`           //	描述
	IconClass   string `gorm:"null"`           //	图标class
	OpUser      string `gorm:"not null"`       //	操作人

	Roles []Roles `gorm:"many2many:role_menu;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
