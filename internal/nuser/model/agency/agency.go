package agency

import (
	"context"
	"github.com/myxy99/ndisk/internal/nuser/model"
	"github.com/myxy99/ndisk/internal/nuser/model/user"
	"gorm.io/gorm"
	"time"
)

type Agency struct {
	ID         uint        `gorm:"primarykey"`             // 主键
	ParentId   uint        `gorm:"parent_id,DEFAULT:0"   ` // 上级机构
	Name       string      `gorm:"name,unique"`            // 部门/11111
	Remark     string      `gorm:"remark"`                 // 机构描述
	Status     uint        `gorm:"status;DEFAULT:1"`       // 是否启用/1,启用,2,禁用
	CreateUId  uint        `gorm:"create_uid"`             // 创建者
	CreateUser user.User   `gorm:"foreignKey:CreateUId"`
	Users      []user.User `gorm:"many2many:agency_user;"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func (m *Agency) TableName() string {
	return "agency"
}

func (m *Agency) Add(ctx context.Context) error {
	return model.MainDB().Table(m.TableName()).WithContext(ctx).Create(m).Error
}

func (m *Agency) Adds(ctx context.Context, data *[]Agency) (count int64, err error) {
	tx := model.MainDB().Table(m.TableName()).WithContext(ctx).CreateInBatches(data, 200)
	err = tx.Error
	count = tx.RowsAffected
	return
}

func (m *Agency) Del(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Delete(m)
	err = tx.Error
	count = tx.RowsAffected
	return
}
func (m *Agency) GetAll(ctx context.Context, data *[]Agency, wheres map[string][]interface{}) (err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	err = db.Find(&data).Error
	return
}
func (m *Agency) Get(ctx context.Context, start int, size int, data *[]Agency, wheres map[string][]interface{}, isDelete bool) (total int64, err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	tx := db.Preload("CreateUser").Limit(size).Offset(start).Find(data)
	total = tx.RowsAffected
	err = tx.Error
	return
}

func (m *Agency) GetById(ctx context.Context, IgnoreDel bool) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	if !IgnoreDel {
		db = db.Unscoped()
	}
	return db.First(m).Error
}

func (m *Agency) GetByWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.First(m).Error
}

func (m *Agency) ExistWhere(ctx context.Context, wheres map[string][]interface{}) bool {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	first := db.First(m)
	return first.RowsAffected != 0
}

func (m *Agency) UpdatesWhere(ctx context.Context, wheres map[string][]interface{}) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *Agency) UpdateWhere(ctx context.Context, wheres map[string][]interface{}, column string, value interface{}) error {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update(column, value).Error
}

func (m *Agency) UpdateStatus(ctx context.Context, status uint32) error {
	return model.MainDB().Table(m.TableName()).WithContext(ctx).Where("id=?", m.ID).Update("status", status).Error
}

func (m *Agency) DelRes(ctx context.Context, wheres map[string][]interface{}) (count int64, err error) {
	db := model.MainDB().Table(m.TableName()).WithContext(ctx)
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	tx := db.Update("deleted_at", nil)
	err = tx.Error
	count = tx.RowsAffected
	return
}
