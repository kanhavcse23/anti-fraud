package transaction_mapper_v1

import (
	entityCoreV1Package "anti-fraud/transaction-service/entity/core/v1"
	entityHttpV1Package "anti-fraud/transaction-service/entity/http/v1"
)

func CreateTransactionPayloadMapper(transactionCreationRequest *entityHttpV1Package.CreateTransactionRequest) *entityCoreV1Package.CreateTransactionPayload {
	return &entityCoreV1Package.CreateTransactionPayload{
		AccountId:       transactionCreationRequest.AccountId,
		OperationTypeId: transactionCreationRequest.OperationTypeId,
		Amount:          transactionCreationRequest.Amount,
	}
}
