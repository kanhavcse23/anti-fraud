package mediator_account_client_v1

import (
	coreV1Package "anti-fraud/account-service/core/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IAccountClient defines methods interface for interacting with the account core service via a mediator pattern.
type IAccountClient interface {
	// SetupCore allows for the injection of IAccountCore, enabling this client
	// to delegate account operations without directly depending on repository logic.
	SetupCore(accountCoreV1 coreV1Package.IAccountCore)

	// GetAccount retrieves an account by its ID, returning a local Account struct.
	GetAccount(logger *logrus.Entry, accountId int, tx *gorm.DB) (*Account, error)
}

// AccountClient implements IAccountClient(interface)
type AccountClient struct {
	accountCoreV1 coreV1Package.IAccountCore
	logger        *logrus.Logger
}

// NewAccountClient create new instance of AccountClient.
func NewAccountClient(logger *logrus.Logger) *AccountClient {

	return &AccountClient{logger: logger}
}

// SetupCore injects the IAccountCore dependency, enabling the client to call account service core methods.
func (client *AccountClient) SetupCore(accountCoreV1 coreV1Package.IAccountCore) {
	client.accountCoreV1 = accountCoreV1
}

// GetAccount calls the core's GetAccount method to get account data.
//
// Steps:
//  1. Invoke the accountCoreV1.GetAccount to fetch the account record from account-service.
//  2. If an error occurs, return an empty Account struct and the error.
//  3. Otherwise, map the fetched account to a mediator-level Account struct.
//
// Parameters:
//   - accountId: The unique ID of the account to fetch.
//   - tx:        db txn.
//
// Returns:
//   - *Account: The mediator-level account struct (e.g. DocumentNumber).
//   - error:    an encountered Error.
func (client *AccountClient) GetAccount(logger *logrus.Entry, accountId int, tx *gorm.DB) (*Account, error) {
	logger.Info("GetAccount method called in mediator-service for account client.")

	account, err := client.accountCoreV1.GetAccount(logger, accountId, tx)
	if err != nil {
		logger.Errorf("Error occured while fetching account data via account service: %s", err.Error())
		return &Account{}, err
	}
	return &Account{Id: int(account.ID), DocumentNumber: account.DocumentNumber}, nil
}
