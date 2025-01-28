package account_core_v1

import (
	entityCoreV1Package "anti-fraud/account-service/entity/core/v1"
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"
	mapperV1Package "anti-fraud/account-service/mapper/v1"
	repoV1Package "anti-fraud/account-service/repository/v1"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IAccountCore defines the methods interface for handling account operations.
type IAccountCore interface {

	// CreateAccount creates a new account if the given document number is not already registered.
	CreateAccount(accountPayload *entityCoreV1Package.CreateAccountPayload, tx *gorm.DB) (*entityDbV1Package.Account, error)

	// GetAccount retrieves an account by its unique ID.
	GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error)
}

// AccountCore implements the IAccountCore interface, containing business logic for account operations.
type AccountCore struct {
	repoV1 repoV1Package.IAccountRepository
	logger *logrus.Logger
}

// NewAccountCore cretae new AccountCore instance.
func NewAccountCore(repoV1 repoV1Package.IAccountRepository, logger *logrus.Logger) *AccountCore {
	return &AccountCore{repoV1: repoV1, logger: logger}
}

// CreateAccount handles the creation of a new account.
//
// Steps:
//   1. Checks if an account with the same document number already exists via the repository.
//   2. If a duplicate is found, it returns that existing account and an error indicating a duplicate.
//   3. Otherwise, maps the request payload to a DB entity and creates a new account record.
//   4. Returns the created account and Error if occured.
//
// Parameters:
//   - accountPayload: Holds the new account details.
//   - tx: db txn.
//
// Returns:
//   - A pointer to the newly created Account entity (or the duplicate if found).
//   - An encountered Error.

func (core *AccountCore) CreateAccount(accountPayload *entityCoreV1Package.CreateAccountPayload, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	core.logger.Info("CreateAccount method called in account core layer.")

	// 1. Check for an existing account with the same document number
	accountFound, err := core.repoV1.CheckDuplicateAccount(accountPayload.DocumentNumber, tx)
	if err != nil {
		return &entityDbV1Package.Account{}, err
	}

	// 2. If a duplicate exists, return it along with an error
	if accountFound.ID > 0 {
		return accountFound, fmt.Errorf("Duplicate account found with document_number: %s", accountPayload.DocumentNumber)
	}

	// 3. Map the incoming payload to a DB entity
	account := mapperV1Package.AccountMapper(accountPayload)

	// 4. Create the new account record in the DB
	err = core.repoV1.CreateAccount(account, tx)
	return account, err

}

// GetAccount retrieves an account by its ID.
//
// Steps:
//  1. Delegates to the repository to fetch the account.
//  2. Returns the account if found, or an error if not found or if any issue occurs.
//
// Parameters:
//   - accountId: ID of the account to retrieve.
//   - tx: db txn.
//
// Returns:
//   - db entity Account.
//   - An encountered Error.
func (core *AccountCore) GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	core.logger.Info("GetAccount method called in account core layer.")
	account, err := core.repoV1.GetAccount(accountId, tx)
	return account, err
}
