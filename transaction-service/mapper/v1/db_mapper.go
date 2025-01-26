package transaction_mapper_v1

import (
	entityCoreV1Path "anti-fraud/transaction-service/entity/core/v1"
	entityDbV1Path "anti-fraud/transaction-service/entity/db/v1"
)

func TransactionMapper(transactionPayload *entityCoreV1Path.CreateTransactionPayload) *entityDbV1Path.Transaction {
	return &entityDbV1Path.Transaction{
		AccountId:       transactionPayload.AccountId,
		OperationTypeId: transactionPayload.OperationTypeId,
		Amount:          transactionPayload.Amount,
	}
}
