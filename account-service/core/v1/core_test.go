package account_core_v1

import (
	"errors"
	"testing"

	entityCoreV1Package "anti-fraud/account-service/entity/core/v1"
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) CheckDuplicateAccount(documentNumber string, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	args := m.Called(documentNumber, tx)
	account, _ := args.Get(0).(*entityDbV1Package.Account)
	return account, args.Error(1)
}

func (m *MockAccountRepository) CreateAccount(account *entityDbV1Package.Account, tx *gorm.DB) error {
	args := m.Called(account, tx)
	return args.Error(0)
}

func (m *MockAccountRepository) GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	args := m.Called(accountId, tx)
	account, _ := args.Get(0).(*entityDbV1Package.Account)
	return account, args.Error(1)
}

// ------------------//
// Utility Mappers  //
// ------------------//
func init() {
}

//---------------------//
//   Unit Test Setup   //
//---------------------//

func setupTest() (*MockAccountRepository, *AccountCore) {
	mockRepo := new(MockAccountRepository)
	logger := logrus.New()

	accountCore := NewAccountCore(mockRepo, logger)

	return mockRepo, accountCore
}

//-------------------------------//
// Tests for CreateAccount Method //
//-------------------------------//

func TestCreateAccount_Success(t *testing.T) {
	mockRepo, accountCore := setupTest()

	payload := &entityCoreV1Package.CreateAccountPayload{
		DocumentNumber: "123456789",
	}

	// 1. CheckDuplicateAccount -> returns nil account, no error
	mockRepo.On("CheckDuplicateAccount", payload.DocumentNumber, mock.Anything).Return(&entityDbV1Package.Account{}, nil)
	// 2. CreateAccount -> returns no error
	mockRepo.On("CreateAccount", mock.Anything, mock.Anything).Return(nil)

	account, err := accountCore.CreateAccount(payload, &gorm.DB{})

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, payload.DocumentNumber, account.DocumentNumber)

	mockRepo.AssertExpectations(t)
}

func TestCreateAccount_DuplicateAccount(t *testing.T) {
	mockRepo, accountCore := setupTest()

	payload := &entityCoreV1Package.CreateAccountPayload{
		DocumentNumber: "123456789",
	}

	existingAccount := &entityDbV1Package.Account{
		Model:          gorm.Model{ID: 1},
		DocumentNumber: "123456789",
	}

	mockRepo.On("CheckDuplicateAccount", payload.DocumentNumber, mock.Anything).Return(existingAccount, nil)

	account, err := accountCore.CreateAccount(payload, &gorm.DB{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Duplicate account found")
	assert.NotNil(t, account)
	assert.Equal(t, existingAccount, account)

	mockRepo.AssertExpectations(t)
}

func TestCreateAccount_RepoError(t *testing.T) {
	mockRepo, accountCore := setupTest()

	payload := &entityCoreV1Package.CreateAccountPayload{
		DocumentNumber: "123456789",
	}

	mockRepo.On("CheckDuplicateAccount", payload.DocumentNumber, mock.Anything).
		Return(&entityDbV1Package.Account{}, errors.New("some repo error"))

	account, err := accountCore.CreateAccount(payload, &gorm.DB{})

	assert.Error(t, err)
	assert.Equal(t, 0, int(account.ID))

	mockRepo.AssertExpectations(t)
}

//----------------------------//
// Tests for GetAccount Method //
//----------------------------//

func TestGetAccount_Success(t *testing.T) {
	mockRepo, accountCore := setupTest()

	// Mock data
	accountID := 1
	accountDB := &entityDbV1Package.Account{
		Model:          gorm.Model{ID: 1},
		DocumentNumber: "123456789",
	}

	mockRepo.On("GetAccount", accountID, mock.Anything).Return(accountDB, nil)

	account, err := accountCore.GetAccount(accountID, &gorm.DB{})

	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, accountID, int(account.ID))
	assert.Equal(t, "123456789", account.DocumentNumber)

	mockRepo.AssertExpectations(t)
}

func TestGetAccount_NotFound(t *testing.T) {
	mockRepo, accountCore := setupTest()

	accountID := 1
	mockRepo.On("GetAccount", accountID, mock.Anything).Return((*entityDbV1Package.Account)(nil), gorm.ErrRecordNotFound)

	account, err := accountCore.GetAccount(accountID, &gorm.DB{})

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Equal(t, gorm.ErrRecordNotFound, err)

	mockRepo.AssertExpectations(t)
}

func TestGetAccount_RepoError(t *testing.T) {
	mockRepo, accountCore := setupTest()

	accountID := 2
	mockRepo.On("GetAccount", accountID, mock.Anything).Return((*entityDbV1Package.Account)(nil), errors.New("db error"))

	account, err := accountCore.GetAccount(accountID, &gorm.DB{})

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "db error")

	mockRepo.AssertExpectations(t)
}
