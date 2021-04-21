package model

import (
	"context"
	"github.com/coder2z/g-server/xtrace"
	"gorm.io/gorm"
)

type (
	Directory struct {
		gorm.Model
		Uid      uint   `gorm:"not null"`
		FileId   uint   `gorm:"null"`
		IsDir    bool   //0 为文件 1 为文件夹
		Name     string `gorm:"null"`
		ParentID uint   `gorm:"DEFAULT:0;not null"`
		Type     uint   `gorm:"DEFAULT:1;not null"` //1为个人文件夹 2为机构文件
	}
)

func (m *Directory) TableName() string {
	return "directory"
}

func (m *Directory) Add(ctx context.Context) error {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model Add")
	defer span.Finish()
	return MainDB().Table(m.TableName()).WithContext(ctx2).Create(m).Error
}

func (m *Directory) Adds(ctx context.Context, data []Directory) (count int64, err error) {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model Adds")
	defer span.Finish()
	tx := MainDB().Table(m.TableName()).WithContext(ctx2).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Directory) Del(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model Del")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Delete(m)
	err = tx.Error
	count = tx.RowsAffected
	return
}
func (m *Directory) GetAll(ctx context.Context, data *[]Directory, wheres map[string][]interface{}) (err error) {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model GetAll")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	err = db.Find(&data).Error
	return
}
func (m *Directory) Get(ctx context.Context, start int, size int, data *[]Directory, wheres map[string][]interface{}, isDelete bool) (total int64, err error) {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model Get")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
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
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model GetById")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	return db.First(m).Error
}

func (m *Directory) GetByWhere(ctx context.Context, wheres map[string][]interface{}) error {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model GetByWhere")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.First(m).Error
}

func (m *Directory) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model ExistWhere")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *Directory) UpdatesWhere(ctx context.Context, wheres map[string][]interface{}) error {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model UpdatesWhere")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *Directory) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model UpdateWhere")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}

func (m *Directory) UpdateStatus(ctx context.Context, status uint32) error {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model UpdateStatus")
	defer span.Finish()
	return MainDB().Table(m.TableName()).WithContext(ctx2).Where("id=?", m.ID).Update("status", status).Error
}

func (m *Directory) DelRes(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model DelRes")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Update("deleted_at", nil)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Directory) Count(ctx context.Context, wheres map[string][]interface{}, isDelete bool) (count int64, err error) {
	span, ctx2 := xtrace.StartSpanFromContext(ctx, "Directory Model Count")
	defer span.Finish()
	db := MainDB().Table(m.TableName()).WithContext(ctx2)
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
