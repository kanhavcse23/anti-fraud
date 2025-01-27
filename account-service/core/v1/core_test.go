package account_core_v1

// import (
// 	entityCoreV1Package "anti-fraud/account-service/entity/core/v1"
// 	entityDbV1Package "anti-fraud/account-service/entity/db/v1"
// 	repoV1Package "anti-fraud/account-service/repository/v1"
// 	"testing"

// 	"github.com/sirupsen/logrus"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// 	"gorm.io/gorm"
// )

// // Define the interface that matches AccountRepository
// type AccountRepositoryInterface interface {
// 	CheckDuplicateAccount(documentNumber string, tx *gorm.DB) (*entityDbV1Package.Account, error)
// 	CreateAccount(account *entityDbV1Package.Account, tx *gorm.DB) error
// 	GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error)
// }

// // MockAccountRepository is a mock for AccountRepository
// type MockAccountRepository struct {
// 	mock.Mock
// 	repoV1Package.AccountRepository // Embed the original repository to satisfy the interface
// }

// func (m *MockAccountRepository) CheckDuplicateAccount(documentNumber string, tx *gorm.DB) (*entityDbV1Package.Account, error) {
// 	args := m.Called(documentNumber, tx)
// 	return args.Get(0).(*entityDbV1Package.Account), args.Error(1)
// }

// func (m *MockAccountRepository) CreateAccount(account *entityDbV1Package.Account, tx *gorm.DB) error {
// 	args := m.Called(account, tx)
// 	return args.Error(0)
// }

// func (m *MockAccountRepository) GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
// 	args := m.Called(accountId, tx)
// 	return args.Get(0).(*entityDbV1Package.Account), args.Error(1)
// }

// func TestCreateAccount(t *testing.T) {
// 	// Setup
// 	mockRepo := &MockAccountRepository{}
// 	logger := logrus.New()
// 	core := NewAccountCore(mockRepo, logger)
// 	tx := &gorm.DB{}

// 	tests := []struct {
// 		name           string
// 		payload        *entityCoreV1Package.CreateAccountPayload
// 		setupMock      func()
// 		expectedError  bool
// 		expectedResult *entityDbV1Package.Account
// 	}{
// 		{
// 			name: "Success - Create New Account",
// 			payload: &entityCoreV1Package.CreateAccountPayload{
// 				DocumentNumber: "12345",
// 			},
// 			setupMock: func() {
// 				mockRepo.On("CheckDuplicateAccount", "12345", tx).Return(&entityDbV1Package.Account{}, nil)
// 				mockRepo.On("CreateAccount", mock.AnythingOfType("*entityDbV1Package.Account"), tx).Return(nil)
// 			},
// 			expectedError: false,
// 			expectedResult: &entityDbV1Package.Account{
// 				DocumentNumber: "12345",
// 			},
// 		},
// 		{
// 			name: "Failure - Duplicate Account",
// 			payload: &entityCoreV1Package.CreateAccountPayload{
// 				DocumentNumber: "12345",
// 			},
// 			setupMock: func() {
// 				mockRepo.On("CheckDuplicateAccount", "12345", tx).Return(&entityDbV1Package.Account{
// 					DocumentNumber: "12345",
// 				}, nil)
// 			},
// 			expectedError: true,
// 			expectedResult: &entityDbV1Package.Account{
// 				DocumentNumber: "12345",
// 			},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mocks for this test case
// 			mockRepo.ExpectedCalls = nil
// 			tt.setupMock()

// 			// Execute
// 			result, err := core.CreateAccount(tt.payload, tx)

// 			// Assert
// 			if tt.expectedError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			assert.Equal(t, tt.expectedResult.DocumentNumber, result.DocumentNumber)
// 			mockRepo.AssertExpectations(t)
// 		})
// 	}
// }

// func TestGetAccount(t *testing.T) {
// 	// Setup
// 	mockRepo := &MockAccountRepository{}
// 	logger := logrus.New()
// 	core := NewAccountCore(mockRepo, logger)
// 	tx := &gorm.DB{}

// 	tests := []struct {
// 		name           string
// 		accountID      int
// 		setupMock      func()
// 		expectedError  bool
// 		expectedResult *entityDbV1Package.Account
// 	}{
// 		{
// 			name:      "Success - Get Existing Account",
// 			accountID: 1,
// 			setupMock: func() {
// 				mockRepo.On("GetAccount", 1, tx).Return(&entityDbV1Package.Account{
// 					DocumentNumber: "12345",
// 				}, nil)
// 			},
// 			expectedError: false,
// 			expectedResult: &entityDbV1Package.Account{
// 				DocumentNumber: "12345",
// 			},
// 		},
// 		{
// 			name:      "Failure - Account Not Found",
// 			accountID: 999,
// 			setupMock: func() {
// 				mockRepo.On("GetAccount", 999, tx).Return(&entityDbV1Package.Account{}, gorm.ErrRecordNotFound)
// 			},
// 			expectedError:  true,
// 			expectedResult: &entityDbV1Package.Account{},
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Setup mocks for this test case
// 			mockRepo.ExpectedCalls = nil
// 			tt.setupMock()

// 			// Execute
// 			result, err := core.GetAccount(tt.accountID, tx)

// 			// Assert
// 			if tt.expectedError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expectedResult.DocumentNumber, result.DocumentNumber)
// 			}
// 			mockRepo.AssertExpectations(t)
// 		})
// 	}
// }
