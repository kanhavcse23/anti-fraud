package transaction_core_v1

import (
	entityCoreV1Package "anti-fraud/transaction-service/entity/core/v1"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"
	mapperV1Package "anti-fraud/transaction-service/mapper/v1"
	repoV1Package "anti-fraud/transaction-service/repository/v1"

	operationClientPackageV1 "anti-fraud/mediator-service/operation-service-client"

	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TransactionCore struct {
	repoV1          *repoV1Package.TransactionRepository
	logger          *logrus.Logger
	operationClient *operationClientPackageV1.OperationClient
}

func NewTransactionCore(repoV1 *repoV1Package.TransactionRepository, logger *logrus.Logger, operationClient *operationClientPackageV1.OperationClient) *TransactionCore {
	return &TransactionCore{repoV1: repoV1, logger: logger, operationClient: operationClient}
}

func (core *TransactionCore) FinalTransactionAmount(amount float64, operationTypeID int, tx *gorm.DB) (float64, error) {
	//Business logic to compute amount by operation type id
	coef, err := core.operationClient.GetOperationCoefficient(operationTypeID, tx)
	if err != nil {
		return amount, err
	}
	return (math.Abs(amount) * float64(coef)), nil
}

func (core *TransactionCore) CreateTransaction(transactionPayload *entityCoreV1Package.CreateTransactionPayload, tx *gorm.DB) (*entityDbV1Package.Transaction, error) {
	core.logger.Info("CreateTransaction method called in transaction core layer.")
	transaction := mapperV1Package.TransactionMapper(transactionPayload)
	amount, err := core.FinalTransactionAmount(transaction.Amount, transaction.OperationTypeId, tx)
	if err != nil {

		return transaction, err
	}
	transaction.Amount = amount
	err = core.repoV1.CreateTransaction(transaction, tx)
	return transaction, err
}
