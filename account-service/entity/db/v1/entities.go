package account_entity_db_v1

import "gorm.io/gorm"

type Account struct {
	gorm.Model
	DocumentNumber string `json:"document_number"`
}

func (Account) TableName() string {
	return "account"
}
