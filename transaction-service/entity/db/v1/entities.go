package transaction_entity_db_v1

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	AccountId       int     `json:"account_id"`
	OperationTypeId int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}
