package account_repo_v1

import (
	entityDbV1Path "anti-fraud/account-service/entity/db/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountRepository struct {
	logger *logrus.Logger
}

func NewAccountRepository(logger *logrus.Logger) *AccountRepository {
	return &AccountRepository{logger: logger}
}
func (repo *AccountRepository) CreateAccount(account *entityDbV1Path.Account, tx *gorm.DB) error {
	result := tx.Create(account)
	return result.Error
}
func (repo *AccountRepository) GetAccount(accountId string, tx *gorm.DB) (*entityDbV1Path.Account, error) {
	var account entityDbV1Path.Account
	result := tx.First(&account, accountId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &account, nil
	}
	return &account, result.Error
}
