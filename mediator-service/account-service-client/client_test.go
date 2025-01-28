package mediator_account_client_v1

import (
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"
	"errors"
	"testing"

	entityCoreV1Package "anti-fraud/account-service/entity/core/v1"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

//-------------------------------------------//
// Mock for IAccountCore
//-------------------------------------------//

type MockAccountCore struct {
	mock.Mock
}

func (m *MockAccountCore) CreateAccount(logger *logrus.Entry, payload *entityCoreV1Package.CreateAccountPayload, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	args := m.Called(payload, tx)
	acc, _ := args.Get(0).(*entityDbV1Package.Account)
	return acc, args.Error(1)
}

func (m *MockAccountCore) GetAccount(logger *logrus.Entry, accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	args := m.Called(accountId, tx)
	acc, _ := args.Get(0).(*entityDbV1Package.Account)
	return acc, args.Error(1)
}

//-------------------------------------------//
// Unit Tests for AccountClient
//-------------------------------------------//

func TestAccountClient_GetAccount_Success(t *testing.T) {
	logger := logrus.New()
	client := NewAccountClient(logger)
	mockCore := new(MockAccountCore)

	client.SetupCore(mockCore)

	accountId := 123
	mockCore.On("GetAccount", accountId, mock.Anything).
		Return(&entityDbV1Package.Account{Model: gorm.Model{ID: 123}, DocumentNumber: "ABC123"}, nil)

	result, err := client.GetAccount(logrus.NewEntry(logrus.New()), accountId, &gorm.DB{})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 123, result.Id)
	assert.Equal(t, "ABC123", result.DocumentNumber)

	mockCore.AssertExpectations(t)
}

func TestAccountClient_GetAccount_Error(t *testing.T) {
	logger := logrus.New()
	client := NewAccountClient(logger)
	mockCore := new(MockAccountCore)

	client.SetupCore(mockCore)

	accountId := 999
	mockCore.On("GetAccount", accountId, mock.Anything).
		Return((*entityDbV1Package.Account)(nil), errors.New("db error"))

	result, err := client.GetAccount(logrus.NewEntry(logrus.New()), accountId, &gorm.DB{})
	assert.Error(t, err)
	assert.Equal(t, int(result.Id), 0)

	mockCore.AssertExpectations(t)
}

func TestAccountClient_SetupCore(t *testing.T) {
	logger := logrus.New()
	client := NewAccountClient(logger)
	mockCore := new(MockAccountCore)

	client.SetupCore(mockCore)
	mockCore.On("GetAccount", 1, mock.Anything).
		Return(&entityDbV1Package.Account{Model: gorm.Model{ID: 1}}, nil)

	res, err := client.GetAccount(logrus.NewEntry(logrus.New()), 1, &gorm.DB{})
	assert.NoError(t, err)
	assert.Equal(t, 1, res.Id)

	mockCore.AssertExpectations(t)
}
