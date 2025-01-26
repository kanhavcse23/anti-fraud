package transaction_core_v1

import (
	entityCoreV1Path "anti-fraud/transaction-service/entity/core/v1"
	entityDbV1Path "anti-fraud/transaction-service/entity/db/v1"
	mapperV1Path "anti-fraud/transaction-service/mapper/v1"
	repoV1Path "anti-fraud/transaction-service/repository/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionCore struct {
	repoV1 *repoV1Path.TransactionRepository
	logger *logrus.Logger
}

func NewTransactionCore(repoV1 *repoV1Path.TransactionRepository, logger *logrus.Logger) *TransactionCore {
	return &TransactionCore{repoV1: repoV1, logger: logger}
}

func (core *TransactionCore) CreateTransaction(transactionPayload *entityCoreV1Path.CreateTransactionPayload, tx *gorm.DB) (*entityDbV1Path.Transaction, error) {
	core.logger.Info("CreateTransaction method called in transaction core layer.")
	transaction := mapperV1Path.TransactionMapper(transactionPayload)
	err := core.repoV1.CreateTransaction(transaction, tx)
	return transaction, err
}
