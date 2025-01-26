package transaction_mapper_v1

import (
	entityCoreV1Package "anti-fraud/transaction-service/entity/core/v1"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"
)

func TransactionMapper(transactionPayload *entityCoreV1Package.CreateTransactionPayload) *entityDbV1Package.Transaction {
	return &entityDbV1Package.Transaction{
		AccountId:       transactionPayload.AccountId,
		OperationTypeId: transactionPayload.OperationTypeId,
		Amount:          transactionPayload.Amount,
	}
}
