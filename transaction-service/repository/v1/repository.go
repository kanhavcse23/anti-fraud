package transaction_repo_v1

import (
	constantPackage "anti-fraud/constants/transaction"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	CreateTransaction(transaction *entityDbV1Package.Transaction, tx *gorm.DB) error
}

type TransactionRepository struct {
	logger *logrus.Logger
}

func NewTransactionRepository(logger *logrus.Logger) *TransactionRepository {
	return &TransactionRepository{logger: logger}
}

func (repo *TransactionRepository) CreateTransaction(transaction *entityDbV1Package.Transaction, tx *gorm.DB) error {
	repo.logger.Info("CreateTransaction method called in transaction repo layer.")
	result := tx.Table(constantPackage.TABLE_NAME).Create(transaction)
	return result.Error
}
