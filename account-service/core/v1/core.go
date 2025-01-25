package account_core_v1

import (
	entityCoreV1Path "anti-fraud/account-service/entity/core/v1"
	entityDbV1Path "anti-fraud/account-service/entity/db/v1"
	mapperV1Path "anti-fraud/account-service/mapper/v1"
	repoV1Path "anti-fraud/account-service/repository/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountCore struct {
	repoV1 *repoV1Path.AccountRepository
	logger *logrus.Logger
}

func NewAccountCore(repoV1 *repoV1Path.AccountRepository, logger *logrus.Logger) *AccountCore {
	return &AccountCore{repoV1: repoV1, logger: logger}
}

func (core *AccountCore) CreateAccount(accountPayload *entityCoreV1Path.CreateAccountPayload, tx *gorm.DB) (*entityDbV1Path.Account, error) {

	account := mapperV1Path.AccountMapper(accountPayload)
	err := core.repoV1.CreateAccount(account, tx)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (core *AccountCore) GetAccount(accountId string, tx *gorm.DB) (*entityDbV1Path.Account, error) {

	account, err := core.repoV1.GetAccount(accountId, tx)
	if err != nil {
		return nil, err
	}
	return account, nil
}
