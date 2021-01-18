package agency

import "gorm.io/gorm"

type Agency struct {



	*gorm.Model
}

func (m *Agency) TableName() string {
	return "agency"
}
