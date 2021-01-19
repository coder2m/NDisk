package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type Agency struct {
	ID         uint   `gorm:"primarykey"`             // 主键
	ParentId   uint   `gorm:"parent_id,DEFAULT:0"   ` // 上级机构
	Name       string `gorm:"name,unique"`            // 部门/11111
	Remark     string `gorm:"remark"`                 // 机构描述
	Status     uint   `gorm:"status;DEFAULT:1"`       // 是否启用/1,启用,2,禁用
	CreateUId  uint   `gorm:"create_uid"`             // 创建者
	CreateUser User   `gorm:"foreignKey:CreateUId"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (m *Agency) TableName() string {
	return "agency"
}

func (m *Agency) Add(ctx context.Context) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (m *Agency) Adds(ctx context.Context, data *[]Agency) (count int64, err error) {
	tx := MainDB().Table(m.TableName()).WithContext(ctx).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Agency) Del(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Delete(m)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Agency) GetAll(ctx context.Context, data *[]Agency, wheres map[string][]interface{}, related bool, IgnoreDel bool) (err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if !IgnoreDel {
		db = db.Unscoped()
	}
	if related {
		db = db.Preload("CreateUser")
	}
	err = db.Find(&data).Error
	return
}

func (m *Agency) Get(ctx context.Context, start int, size int, data *[]Agency, wheres map[string][]interface{}, isDelete bool, related bool) (total int64, err error) {
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
		db = db.Preload("CreateUser")
	}
	tx := db.Limit(size).Offset(start).Find(data)
	total = tx.RowsAffected
	err = tx.Error
	return
}

func (m *Agency) GetById(ctx context.Context, IgnoreDel, related bool) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	if related {
		db = db.Preload("CreateUser")
	}
	return db.First(m).Error
}

func (m *Agency) GetByWhere(ctx context.Context, wheres map[string][]interface{}, related bool, args ...interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if related {
		db = db.Preload("CreateUser", args...)
	}
	return db.First(m).Error
}

func (m *Agency) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *Agency) UpdatesWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *Agency) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}

func (m *Agency) UpdateStatus(ctx context.Context, status uint32) error {
	return MainDB().Table(m.TableName()).WithContext(ctx).Where("id=?", m.ID).Update("status", status).Error
}

func (m *Agency) DelRes(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Update("deleted_at", nil)
	err = tx.Error
	count = tx.RowsAffected
	return
}
