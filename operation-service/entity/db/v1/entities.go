package account_entity_db_v1

import "gorm.io/gorm"

type Operation struct {
	gorm.Model
	Description string `json:"description"`
	Coefficient int    `json:"coefficient"`
}
