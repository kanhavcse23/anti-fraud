package account_repo_v1

import (
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"

	constantPackage "anti-fraud/constants/account"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IAccountRepository defines methods interface for account-related db operations.
type IAccountRepository interface {

	// CreateAccount persists a new account record to the db.
	CreateAccount(account *entityDbV1Package.Account, tx *gorm.DB) error

	// GetAccount retrieves an account by its unique ID.
	GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error)

	// CheckDuplicateAccount checks if an account with the given document number already exists.
	CheckDuplicateAccount(documentNumber string, tx *gorm.DB) (*entityDbV1Package.Account, error)
}

// AccountRepository implements IAccountRepository methods.
type AccountRepository struct {
	logger *logrus.Logger
}

// NewAccountRepository returns a new AccountRepository instance.
func NewAccountRepository(logger *logrus.Logger) *AccountRepository {
	return &AccountRepository{logger: logger}
}

// CreateAccount inserts a new account record into the database.
//
// Steps:
//  1. Perform an INSERT operation on the account table.
//  3. Returns any error encountered during the insertion.
//
// Parameters:
//   - account: account db entity.
//   - tx: db txn
//
// Returns:
//   - An error if the insert fails, otherwise nil.
func (repo *AccountRepository) CreateAccount(account *entityDbV1Package.Account, tx *gorm.DB) error {
	repo.logger.Info("CreateAccount method called in account repo layer.")
	result := tx.Table(constantPackage.TABLE_NAME).Create(account)
	if result.Error != nil {
		repo.logger.Errorf("Failed to create account: %v", result.Error)
	}
	return result.Error
}

// GetAccount fetches an account record by account id (p.k).
//
// Steps:
//  1. Executes a SELECT query using the given accountId as a primary key lookup.
//  2. If the record is not found, it returns an empty account and nil error
//  3. Otherwise, returns the account data.
//
// Parameters:
//   - accountId: account id to find record
//   - tx: db txn.
//
// Returns:
//   - db entity account.
//   - Encountered Error.
func (repo *AccountRepository) GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	repo.logger.Info("GetAccount method called in account repo layer.")
	var account entityDbV1Package.Account
	result := tx.Table(constantPackage.TABLE_NAME).First(&account, accountId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		repo.logger.Error("Failed to find account")
		return &account, nil
	}
	return &account, result.Error
}

// CheckDuplicateAccount determines if an account with the specified document number already exists.
//
// Steps:
//   1. Searches account that has similar `document_number`.
//   2. If no record is found, returns an empty account object and nil.
//   3. Otherwise, returns the found account and Error encountered.
//
// Parameters:
//   - documentNumber: filter used to find account.
//   - tx: db txn.
//
// Returns:
//   - db entity account.
//   - Encountered Error.

func (repo *AccountRepository) CheckDuplicateAccount(documentNumber string, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	repo.logger.Info("CheckDuplicateAccount method called in account repo layer.")
	var account entityDbV1Package.Account
	result := tx.Table(constantPackage.TABLE_NAME).
		Where("document_number = ?", documentNumber).First(&account)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound { // account doesn't exist with `documentNumber`
		repo.logger.Errorf("Failed to find account")
		return &account, nil
	}
	return &account, result.Error
}
