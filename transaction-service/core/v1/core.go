package transaction_core_v1

import (
	entityCoreV1Package "anti-fraud/transaction-service/entity/core/v1"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"
	mapperV1Package "anti-fraud/transaction-service/mapper/v1"
	repoV1Package "anti-fraud/transaction-service/repository/v1"
	"fmt"

	accountClientPackageV1 "anti-fraud/mediator-service/account-service-client"
	operationClientPackageV1 "anti-fraud/mediator-service/operation-service-client"

	"math"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ITransactionCore defines the methods for core business logic for transaction operations.
type ITransactionCore interface {
	// CreateTransaction creates a new transaction record in the db
	CreateTransaction(transactionPayload *entityCoreV1Package.CreateTransactionPayload, tx *gorm.DB) (*entityDbV1Package.Transaction, error)

	// FinalTransactionAmount applies business logic to compute the final transaction amount
	// based on the operationTypeID and coefficient retrieved from the operation service.
	FinalTransactionAmount(amount float64, operationTypeID int, tx *gorm.DB) (float64, error)

	// CheckAccountIdExist verifies whether the provided accountId exists by calling the account service.
	// Returns an error if the account is not found.
	CheckAccountIdExist(accountId int, tx *gorm.DB) error
}

// TransactionCore implements ITransactionCore.
type TransactionCore struct {
	repoV1          repoV1Package.ITransactionRepository
	logger          *logrus.Logger
	operationClient operationClientPackageV1.IOperationClient
	accountClient   accountClientPackageV1.IAccountClient
}

// NewTransactionCore return new TransactionCore instance.
func NewTransactionCore(repoV1 repoV1Package.ITransactionRepository, logger *logrus.Logger, operationClient operationClientPackageV1.IOperationClient, accountClient accountClientPackageV1.IAccountClient) *TransactionCore {
	return &TransactionCore{repoV1: repoV1, logger: logger, operationClient: operationClient, accountClient: accountClient}
}

// FinalTransactionAmount calculates the final amount for a transaction based on the operation type.
//
// Steps:
//  1. Retrieve the coefficient from the operation service using the provided operationTypeID.
//  2. Convert the original amount to its absolute value and multiply by the coefficient.
//  3. Return the computed final amount or the original amount plus an error if an external call fails.
//
// Parameters:
//   - amount: The initial transaction amount.
//   - operationTypeID: The ID representing the type of operation (e.g., debit, credit).
//   - tx: db txn.
//
// Returns:
//   - float64: The final computed transaction amount after applying the operation coefficient.
//   - error:   Encountered Error.
func (core *TransactionCore) FinalTransactionAmount(amount float64, operationTypeID int, tx *gorm.DB) (float64, error) {
	// Fetch coefficient from the operation service
	coef, err := core.operationClient.GetOperationCoefficient(operationTypeID, tx)
	if err != nil {
		return amount, err
	}

	// Compute final amount using absolute value and the retrieved coefficient
	return (math.Abs(amount) * float64(coef)), nil
}

// CheckAccountIdExist verifies that the provided accountId exists in the account service.
// Steps:
//  1. Calls the account service to retrieve an account by account id.
//  2. If no account is found (ID == 0) in db, returns an Error.
//  3. Otherwise, returns nil to indicate the account exists.
//
// Parameters:
//   - accountId: id of account.
//   - tx:        db txn.
//
// Returns:
//   - error: If the account is not found or if there's an error in the account service call.
func (core *TransactionCore) CheckAccountIdExist(accountId int, tx *gorm.DB) error {

	account, err := core.accountClient.GetAccount(accountId, tx)
	if err != nil {
		return err
	}
	if account.Id == 0 { // account id not found in database
		return fmt.Errorf("account_id: %d not found in database", accountId)
	}
	return nil
}

// CreateTransaction creates a new transaction record in the db after verifying the account and
// calculating the final amount.
//
// Steps:
//   1. Ensure the account ID is valid via CheckAccountIdExist. If invalid, return an error.
//   2. Map the incoming payload to a database transaction entity.
//   3. Calculate the final transaction amount using FinalTransactionAmount.
//   4. Persist the transaction in the DB
//
// Parameters:
//   - transactionPayload: Payload containing the base data needed to create a transaction (accountId, amount, etc.).
//   - tx:                 The GORM transaction in which all DB operations will run.
//
// Returns:
//   - A pointer to the newly created Transaction entity.
//   - An error if account validation fails, amount calculation fails, or the DB create operation fails.

func (core *TransactionCore) CreateTransaction(transactionPayload *entityCoreV1Package.CreateTransactionPayload, tx *gorm.DB) (*entityDbV1Package.Transaction, error) {
	core.logger.Info("CreateTransaction method called in transaction core layer.")

	// // 1. Validate the account_id exist in db
	err := core.CheckAccountIdExist(transactionPayload.AccountId, tx)
	if err != nil {
		return &entityDbV1Package.Transaction{}, err
	}

	// 2. Map the payload to a DB entity
	transaction := mapperV1Package.TransactionMapper(transactionPayload)

	// 3. Compute the final transaction amount
	amount, err := core.FinalTransactionAmount(transaction.Amount, transaction.OperationTypeId, tx)
	if err != nil {
		return transaction, err
	}
	transaction.Amount = amount

	// 4. Persist the transaction in the DB
	err = core.repoV1.CreateTransaction(transaction, tx)
	return transaction, err
}
