package model

import (
	"gorm.io/gorm"
	"time"
)

// api 资源集合
type Resources struct {
	ID uint `gorm:"primarykey"`

	Name        string `gorm:"not null"`
	Path        string `gorm:"not null"`
	Action      string `gorm:"not null"`
	Description string `gorm:"not null"`
	OpUser      string `gorm:"not null"`       //操作人

	Roles []Roles `gorm:"many2many:role_resource;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
