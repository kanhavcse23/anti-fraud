package transaction_repo_v1

import (
	entityDbV1Path "anti-fraud/transaction-service/entity/db/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionRepository struct {
	logger *logrus.Logger
}

func NewTransactionRepository(logger *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{logger: logger}
}

func (repo *TransactionRepository) CreateTransaction(transaction *entityDbV1Path.Transaction, tx *gorm.DB) error {
	repo.logger.Info("CreateTransaction method called in transaction repo layer.")
	result := tx.Create(transaction)
	return result.Error
}
