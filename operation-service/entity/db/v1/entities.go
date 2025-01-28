package account_entity_db_v1

import (
	constantPackage "anti-fraud/constants/operation"

	"gorm.io/gorm"
)

type Operation struct {
	gorm.Model
	Description string `json:"description"`
	Coefficient int    `json:"coefficient"`
}

func (Operation) TableName() string {
	return constantPackage.TABLE_NAME
}
