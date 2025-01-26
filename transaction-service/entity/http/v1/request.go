package transaction_entity_http_v1

import "errors"

type CreateTransactionRequest struct {
	AccountId       int     `json:"account_id"`
	OperationTypeId int     `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func (createAccountRequest *CreateTransactionRequest) Validate() error {
	if createAccountRequest.AccountId == 0 {
		return errors.New("account_id is mandatory")
	}
	if createAccountRequest.OperationTypeId == 0 {
		return errors.New("operation_type_id is mandatory")
	}
	if createAccountRequest.Amount == 0.0 {
		return errors.New("amount should be non-zero")
	}
	return nil
}
