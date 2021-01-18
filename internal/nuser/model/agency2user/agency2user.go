package agency2user

import "gorm.io/gorm"

type AgencyUser struct {
	*gorm.Model
}

func (m *AgencyUser) TableName() string {
	return "agency_user"
}
