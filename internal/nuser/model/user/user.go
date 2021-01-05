/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/5 15:26
 **/
package user

import (
	"github.com/myxy99/ndisk/internal/nuser/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Name     string `gorm:"not null;unique_index;"`
	Alias    string `gorm:"not null"`
	Tel      string `gorm:"type:varchar(11);unique_index;"`
	Email    string `gorm:"type:varchar(100);unique_index;not null"`
	Password string `gorm:"not null"`
	Status   int    `gorm:"DEFAULT:1;not null"`

	*gorm.Model
}

func (m *User) TableName() string {
	return "user"
}

func (m *User) Add() error {
	return model.MainDB().Table(m.TableName()).Create(m).Error
}
func (m *User) Del(wheres map[string]interface{}) error {
	return model.MainDB().Table(m.TableName()).Where(wheres).Delete(m).Error
}
func (m *User) GetAll(data *[]User, wheres map[string][]interface{}) (err error) {
	db := model.MainDB().Table(m.TableName())
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	err = db.Find(&data).Error
	return
}
func (m *User) Get(start int, size int, data *[]User, wheres map[string][]interface{}, isDelete bool) (total int64, err error) {
	db := model.MainDB().Table(m.TableName())
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	if isDelete {
		db = db.Unscoped().Where("interface.deleted_at is not null")
	} else {
		db = db.Where(map[string]interface{}{"deleted_at": nil})
	}
	err = db.Limit(size).Offset(start).Find(data).Error
	err = db.Count(&total).Error
	return
}

func (m *User) GetById() error {
	return model.MainDB().Table(m.TableName()).First(m).Error
}

func (m *User) GetByWhere(wheres map[string][]interface{}) error {
	db := model.MainDB().Table(m.TableName())
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.First(m).Error
}

func (m *User) UpdateWhere(wheres map[string][]interface{}) error {
	db := model.MainDB().Table(m.TableName())
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Updates(m).Error
}

func (m *User) UpdateStatus(status int) error {
	return model.MainDB().Table(m.TableName()).Where("id=?", m.ID).Update("status", status).Error
}

func (m *User) DelRes(wheres map[string][]interface{}) error {
	db := model.MainDB().Table(m.TableName())
	for s, i := range wheres {
		db = db.Where(s, i...)
	}
	return db.Update("deleted_at", nil).Error
}

func (m *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(password))
	return err == nil
}

func (m *User) SetPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	m.Password = string(bytes)
	return nil
}
