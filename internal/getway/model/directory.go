package model

import (
	"context"
	"gorm.io/gorm"
)

type (
	Directory struct {
		gorm.Model
		Uid      uint `gorm:"not null"`
		FileId   uint `gorm:"null"`
		IsDir    bool //0 为文件 1 为文件夹
		Name     string `gorm:"null"`
		ParentID uint   `gorm:"DEFAULT:0;not null"`
	}
)

func (m *Directory) TableName() string {
	return "directory"
}

func (m *Directory) Add(ctx context.Context) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (m *Directory) Adds(ctx context.Context, data []Directory) (count int64, err error) {
	tx := MainDB().Table(m.TableName()).WithContext(ctx).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Directory) Del(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Delete(m)
	err = tx.Error
	count = tx.RowsAffected
	return
}
func (m *Directory) GetAll(ctx context.Context, data *[]Directory, wheres map[string][]interface{}) (err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	err = db.Find(&data).Error
	return
}
func (m *Directory) Get(ctx context.Context, start int, size int, data *[]Directory, wheres map[string][]interface{}, isDelete bool) (total int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	}
	tx := db.Limit(size).Offset(start).Find(data)
	err = tx.Error
	total, err = m.Count(ctx, wheres, isDelete)
	return
}

func (m *Directory) GetById(ctx context.Context, IgnoreDel bool) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	return db.First(m).Error
}

func (m *Directory) GetByWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.First(m).Error
}

func (m *Directory) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *Directory) UpdatesWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *Directory) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}

func (m *Directory) UpdateStatus(ctx context.Context, status uint32) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Where("id=?", m.ID).Update("status", status).Error
}

func (m *Directory) DelRes(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Update("deleted_at", nil)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Directory) Count(ctx context.Context, wheres map[string][]interface{}, isDelete bool) (count int64, err error) {
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
