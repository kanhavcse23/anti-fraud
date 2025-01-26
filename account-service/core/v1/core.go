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

type AccountCore struct {
	repoV1 *repoV1Package.AccountRepository
	logger *logrus.Logger
}

func NewAccountCore(repoV1 *repoV1Package.AccountRepository, logger *logrus.Logger) *AccountCore {
	return &AccountCore{repoV1: repoV1, logger: logger}
}

func (core *AccountCore) CreateAccount(accountPayload *entityCoreV1Package.CreateAccountPayload, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	core.logger.Info("CreateAccount method called in account core layer.")
	accountFound, err := core.repoV1.CheckDuplicateAccount(accountPayload.DocumentNumber, tx)
	if err != nil {
		return &entityDbV1Package.Account{}, err
	}
	if accountFound.ID > 0 {
		return accountFound, fmt.Errorf("Duplicate account found with document_number: %s", accountPayload.DocumentNumber)
	}
	account := mapperV1Package.AccountMapper(accountPayload)
	err = core.repoV1.CreateAccount(account, tx)
	return account, err

}

func (core *AccountCore) GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	core.logger.Info("GetAccount method called in account core layer.")
	account, err := core.repoV1.GetAccount(accountId, tx)
	return account, err
}
