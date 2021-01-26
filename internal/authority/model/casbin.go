package model

import "context"

type CasbinRule struct {
	PType string `json:"p_type" gorm:"size:100;"`
	V0    string `json:"v0" gorm:"size:100;"`
	V1    string `json:"v1" gorm:"size:100;"`
	V2    string `json:"v2" gorm:"size:100;"`
	V3    string `json:"v3" gorm:"size:100;"`
	V4    string `json:"v4" gorm:"size:100;"`
	V5    string `json:"v5" gorm:"size:100;"`
}

func (CasbinRule) TableName() string {
	return "sys_casbin"
}

type UsersRoles struct {
	Uid   uint32
	Roles string
}

//获取用户权限
func (m CasbinRule) GetUsersRoles(ctx context.Context, uid []uint32) ([]UsersRoles, error) {
	var (
		sql  = `SELECT v0 AS "uid",GROUP_CONCAT(v1) AS "roles" FROM sys_casbin WHERE p_type = "g" AND v0 IN ( ? ) GROUP BY v0`
		data []UsersRoles
	)
	tx := MainDB().WithContext(ctx).Raw(sql, uid).Find(&data)
	return data, tx.Error
}
