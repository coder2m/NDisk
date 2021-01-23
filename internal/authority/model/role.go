package model

import (
	"gorm.io/gorm"
	"time"
)

type Roles struct {
	ID uint `gorm:"primarykey"`

	Name        string `gorm:"not null"`
	Description string `gorm:"not null"`
	OpUser      string `gorm:"not null"` //操作人

	Menus     []Menu      `gorm:"many2many:role_menu;"`
	Resources []Resources `gorm:"many2many:role_resource;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
