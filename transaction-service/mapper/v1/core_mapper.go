package transaction_mapper_v1

import (
	entityCoreV1Path "anti-fraud/transaction-service/entity/core/v1"
	entityHttpV1Path "anti-fraud/transaction-service/entity/http/v1"
)

func CreateTransactionPayloadMapper(transactionCreationRequest *entityHttpV1Path.CreateTransactionRequest) *entityCoreV1Path.CreateTransactionPayload {
	return &entityCoreV1Path.CreateTransactionPayload{
		AccountId:       transactionCreationRequest.AccountId,
		OperationTypeId: transactionCreationRequest.OperationTypeId,
		Amount:          transactionCreationRequest.Amount,
	}
}
