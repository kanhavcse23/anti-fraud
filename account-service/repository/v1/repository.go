package account_repo_v1

import (
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"

	constantPackage "anti-fraud/constants/account"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AccountRepository struct {
	logger *logrus.Logger
}

func NewAccountRepository(logger *logrus.Logger) *AccountRepository {
	return &AccountRepository{logger: logger}
}
func (repo *AccountRepository) CreateAccount(account *entityDbV1Package.Account, tx *gorm.DB) error {
	repo.logger.Info("CreateAccount method called in account repo layer.")
	result := tx.Table(constantPackage.TABLE_NAME).Create(account)
	return result.Error
}
func (repo *AccountRepository) GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	repo.logger.Info("GetAccount method called in account repo layer.")
	var account entityDbV1Package.Account
	result := tx.Table(constantPackage.TABLE_NAME).First(&account, accountId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &account, nil
	}
	return &account, result.Error
}
func (repo *AccountRepository) CheckDuplicateAccount(documentNumber string, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	repo.logger.Info("CheckDuplicateAccount method called in account repo layer.")
	var account entityDbV1Package.Account
	result := tx.Table(constantPackage.TABLE_NAME).
		Where("document_number = ?", documentNumber).First(&account)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &account, nil
	}
	return &account, result.Error
}
