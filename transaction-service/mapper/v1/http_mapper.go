package transaction_mapper_v1

import (
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"
	entityHttpV1Package "anti-fraud/transaction-service/entity/http/v1"
)

func TransactionDetailsResponseMapper(transaction *entityDbV1Package.Transaction) *entityHttpV1Package.CreateTransactionResponse {
	return &entityHttpV1Package.CreateTransactionResponse{
		TransactionID:   int(transaction.ID),
		AccountId:       transaction.AccountId,
		OperationTypeId: transaction.OperationTypeId,
		Amount:          transaction.Amount,
		EventDate:       transaction.CreatedAt,
	}
}
