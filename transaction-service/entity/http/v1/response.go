package transaction_entity_http_v1

import "time"

type CreateTransactionResponse struct {
	TransactionID   int       `json:"transaction_id"`
	AccountId       int       `json:"account_id"`
	OperationTypeId int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}
