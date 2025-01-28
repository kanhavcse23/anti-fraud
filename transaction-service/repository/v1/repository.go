package transaction_repo_v1

import (
	constantPackage "anti-fraud/constants/transaction"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ITransactionRepository defines method interface for performing operations related to transactions in the db.
type ITransactionRepository interface {

	// CreateTransaction persists a Transaction entity to the db.
	CreateTransaction(transaction *entityDbV1Package.Transaction, tx *gorm.DB) error
}

// TransactionRepository implements the ITransactionRepository interface,
type TransactionRepository struct {
	logger *logrus.Logger
}

// NewTransactionRepository return new TransactionRepository instance.
func NewTransactionRepository(logger *logrus.Logger) *TransactionRepository {

	return &TransactionRepository{logger: logger}
}

// CreateTransaction inserts a new transaction record into the db.
//
// Steps:
//  1. Perform an INSERT on the table specified by constantPackage.TABLE_NAME.
//  2. Returns Error.
//
// Parameters:
//   - transaction: entity db transaction.
//   - tx:          db txn.
//
// Returns:
//   - error: If the record fails to insert. else return nil.
func (repo *TransactionRepository) CreateTransaction(transaction *entityDbV1Package.Transaction, tx *gorm.DB) error {
	repo.logger.Info("CreateTransaction method called in transaction repo layer.")
	result := tx.Table(constantPackage.TABLE_NAME).Create(transaction)
	return result.Error
}
