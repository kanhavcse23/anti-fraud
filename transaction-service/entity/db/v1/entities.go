package transaction_entity_db_v1

import (
	constantPackage "anti-fraud/constants/transaction"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	AccountId       int     `json:"account_id"`
	OperationTypeId int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func (Transaction) TableName() string {
	return constantPackage.TABLE_NAME
}
