package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Roles struct {
	ID uint `gorm:"primarykey"`

	Name        string `gorm:"not null,unique"`
	Description string `gorm:"not null"`

	Menus     []Menu      `gorm:"many2many:role_menu;"`
	Resources []Resources `gorm:"many2many:role_resource;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Roles) TableName() string {
	return "roles"
}

func (m *Roles) Add(ctx context.Context) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (m *Roles) Adds(ctx context.Context, data *[]Roles) (count int64, err error) {
	tx := MainDB().Table(m.TableName()).WithContext(ctx).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Roles) Del(ctx context.Context, ids []uint32) (count int64, err error) {
	tx := MainDB().Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, err
	}

	var list []Roles
	if err := tx.Table(m.TableName()).Where("id in (?)", ids).Find(&list).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if res := tx.Table(m.TableName()).Where("id in (?)", ids).Delete(m); res.Error != nil {
		tx.Rollback()
		return 0, err
	} else {
		count = res.RowsAffected
	}

	for _, roles := range list {
		if err := tx.Table(new(CasbinRule).TableName()).Where("v0 in (?)", roles.Name).Delete(&CasbinRule{}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}
	err = tx.Error

	return
}

func (m *Roles) GetAll(ctx context.Context, data *[]Roles, wheres map[string][]interface{}, related bool, IgnoreDel bool) (err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if !IgnoreDel {
		db = db.Unscoped()
	}
	if related {
		db = db.Preload("Menus").Preload("Resources")
	}
	err = db.Find(&data).Error
	return
}

func (m *Roles) Get(ctx context.Context, start int, size int, data *[]Roles, wheres map[string][]interface{}, isDelete, related bool) (total int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if related {
		db = db.Preload("Menus").Preload("Resources")
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	tx := db.Limit(size).Offset(start).Find(data)
	total = tx.RowsAffected
	err = tx.Error
	return
}

func (m *Roles) GetById(ctx context.Context, IgnoreDel, related bool) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	if related {
		db = db.Preload("Menus").Preload("Resources")
	}
	return db.First(m).Error
}

func (m *Roles) GetByWhere(ctx context.Context, wheres map[string][]interface{}, related bool) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if related {
		db = db.Preload("Menus").Preload("Resources")
	}
	return db.First(m).Error
}

func (m *Roles) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *Roles) UpdatesWhereById(ctx context.Context, id uint) (err error) {
	tx := MainDB().Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var r Roles

	if err := tx.Table(m.TableName()).First(&r, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table(m.TableName()).Where("id = ?", id).Updates(m).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table(new(CasbinRule).TableName()).Where("v0 = ?", r.Name).Update("v0", m.Name).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	err = tx.Error
	return
}

func (m *Roles) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}
