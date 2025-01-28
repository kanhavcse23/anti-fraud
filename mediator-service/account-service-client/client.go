package mediator_account_client_v1

import (
	coreV1Package "anti-fraud/account-service/core/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IAccountClient interface {
	SetupCore(accountCoreV1 coreV1Package.IAccountCore)
	GetAccount(accountId int, tx *gorm.DB) (*Account, error)
}
type AccountClient struct {
	accountCoreV1 coreV1Package.IAccountCore
	logger        *logrus.Logger
}

func NewAccountClient(logger *logrus.Logger) *AccountClient {

	return &AccountClient{logger: logger}
}
func (client *AccountClient) SetupCore(accountCoreV1 coreV1Package.IAccountCore) {
	client.accountCoreV1 = accountCoreV1
}
func (client *AccountClient) GetAccount(accountId int, tx *gorm.DB) (*Account, error) {
	client.logger.Info("GetAccount method called in mediator-service for account client.")

	account, err := client.accountCoreV1.GetAccount(accountId, tx)
	if err != nil {
		return &Account{}, err
	}
	return &Account{Id: int(account.ID), DocumentNumber: account.DocumentNumber}, nil
}
