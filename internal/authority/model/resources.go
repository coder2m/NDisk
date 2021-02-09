package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

// api 资源集合
type Resources struct {
	ID uint `gorm:"primarykey"`

	Name        string `gorm:"not null,unique"`
	Path        string `gorm:"not null"`
	Action      string `gorm:"not null"`
	Description string `gorm:"not null"`

	Roles []Roles `gorm:"many2many:role_resource;"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (m *Resources) TableName() string {
	return "resources"
}

func (m *Resources) Add(ctx context.Context) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (m *Resources) Adds(ctx context.Context, data *[]Resources) (count int64, err error) {
	tx := MainDB().Table(m.TableName()).WithContext(ctx).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Resources) Del(ctx context.Context, ids []uint32) (count int64, err error) {
	tx := MainDB().Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return 0, err
	}

	var list []Resources
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

	for _, re := range list {
		if err := tx.Table(new(CasbinRule).TableName()).Where("v1 = ? and v2 = ?", re.Path, re.Action).Unscoped().Delete(&CasbinRule{}).Error; err != nil {
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

func (m *Resources) GetAll(ctx context.Context, data *[]Resources, wheres map[string][]interface{}, related bool, IgnoreDel bool) (err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if !IgnoreDel {
		db = db.Unscoped()
	}
	if related {
		db = db.Preload("Roles")
	}
	err = db.Find(&data).Error
	return
}

func (m *Resources) Get(ctx context.Context, start int, size int, data *[]Resources, wheres map[string][]interface{}, isDelete bool, related bool) (total int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	if related {
		db = db.Preload("Roles")
	}
	tx := db.Limit(size).Offset(start).Find(data)
	err = tx.Error
	total, err = m.Count(ctx, wheres, isDelete)
	return
}

func (m *Resources) GetById(ctx context.Context, IgnoreDel, related bool) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	if related {
		db = db.Preload("Roles")
	}
	return db.First(m).Error
}

func (m *Resources) GetByWhere(ctx context.Context, wheres map[string][]interface{}, related bool, args ...interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if related {
		db = db.Preload("Roles", args...)
	}
	return db.First(m).Error
}

func (m *Resources) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *Resources) UpdatesWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *Resources) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}

func (m Resources) UpdatesWhereById(ctx context.Context, id uint) error {
	tx := MainDB().Begin().WithContext(ctx)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var r Resources

	if err := tx.Table(m.TableName()).First(&r, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table(m.TableName()).Where("id = ?", id).Updates(m).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Table(new(CasbinRule).TableName()).
		Where("v1 = ? and v2 = ?", r.Path, r.Action).
		Updates(&CasbinRule{
			V1: m.Path,
			V2: m.Action,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return tx.Error
}


func (m *Resources) Count(ctx context.Context, wheres map[string][]interface{}, isDelete bool) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	tx := db.Count(&count)
	return count, tx.Error
}