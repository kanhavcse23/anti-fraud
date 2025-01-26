package transaction_mapper_v1

import (
	entityDbV1Path "anti-fraud/transaction-service/entity/db/v1"
	entityHttpV1Path "anti-fraud/transaction-service/entity/http/v1"
)

func TransactionDetailsResponseMapper(transaction *entityDbV1Path.Transaction) *entityHttpV1Path.CreateTransactionResponse {
	return &entityHttpV1Path.CreateTransactionResponse{
		TransactionID:   int(transaction.ID),
		AccountId:       transaction.AccountId,
		OperationTypeId: transaction.OperationTypeId,
		Amount:          transaction.Amount,
		EventDate:       transaction.CreatedAt,
	}
}
