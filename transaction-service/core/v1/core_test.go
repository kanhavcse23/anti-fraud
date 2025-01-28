package transaction_core_v1

import (
	accountCoreV1Package "anti-fraud/account-service/core/v1"
	accountClientPackageV1 "anti-fraud/mediator-service/account-service-client"
	opsCoreV1Package "anti-fraud/operation-service/core/v1"
	entityCoreV1Package "anti-fraud/transaction-service/entity/core/v1"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"

	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(logger *logrus.Entry, transaction *entityDbV1Package.Transaction, tx *gorm.DB) error {
	args := m.Called(transaction, tx)
	return args.Error(0)
}

type MockOperationClient struct {
	mock.Mock
}

func (m *MockOperationClient) GetOperationCoefficient(logger *logrus.Entry, operationTypeID int, tx *gorm.DB) (int, error) {
	args := m.Called(operationTypeID, tx)
	coef, _ := args.Get(0).(int)
	return coef, args.Error(1)
}
func (m *MockOperationClient) SetupCore(core opsCoreV1Package.IOperationCore) {
	m.Called(core)
}

type MockAccountClient struct {
	mock.Mock
}

func (m *MockAccountClient) GetAccount(logger *logrus.Entry, accountId int, tx *gorm.DB) (*accountClientPackageV1.Account, error) {
	args := m.Called(accountId, tx)
	acc, _ := args.Get(0).(*accountClientPackageV1.Account)
	return acc, args.Error(1)
}

func (m *MockAccountClient) SetupCore(core accountCoreV1Package.IAccountCore) {
	m.Called(core)
}

//-------------------------------------------//
// 2. Setup Helpers
//-------------------------------------------//

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

func setupTestCore(t *testing.T) (*TransactionCore, *MockTransactionRepository, *MockOperationClient, *MockAccountClient, *gorm.DB) {
	logger := logrus.New()
	db := setupTestDB(t)

	repoMock := new(MockTransactionRepository)
	opMock := new(MockOperationClient)
	accMock := new(MockAccountClient)

	core := NewTransactionCore(repoMock, logger, opMock, accMock)

	return core, repoMock, opMock, accMock, db
}

//-------------------------------------------//
// 3. Test: FinalTransactionAmount
//-------------------------------------------//

func TestFinalTransactionAmount_Success(t *testing.T) {
	core, _, opMock, _, db := setupTestCore(t)

	opMock.On("GetOperationCoefficient", 1, mock.Anything).
		Return(3, nil)

	amount, err := core.FinalTransactionAmount(logrus.NewEntry(logrus.New()), 100, 1, db)
	assert.NoError(t, err)
	assert.Equal(t, 300.0, amount)

	opMock.AssertExpectations(t)
}

func TestFinalTransactionAmount_Error(t *testing.T) {
	core, _, opMock, _, db := setupTestCore(t)

	opMock.On("GetOperationCoefficient", 2, mock.Anything).
		Return(0.0, errors.New("operation client error"))

	amount, err := core.FinalTransactionAmount(logrus.NewEntry(logrus.New()), 100, 2, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "operation client error")
	assert.Equal(t, 100.0, amount)

	opMock.AssertExpectations(t)
}

//-------------------------------------------//
// 4. Test: CheckAccountIdExist
//-------------------------------------------//

func TestCheckAccountIdExist_Success(t *testing.T) {
	core, _, _, accMock, db := setupTestCore(t)

	accMock.On("GetAccount", 123, mock.Anything).
		Return(&accountClientPackageV1.Account{Id: 123}, nil)

	err := core.CheckAccountIdExist(logrus.NewEntry(logrus.New()), 123, db)
	assert.NoError(t, err)

	accMock.AssertExpectations(t)
}

func TestCheckAccountIdExist_NotFound(t *testing.T) {
	core, _, _, accMock, db := setupTestCore(t)

	accMock.On("GetAccount", 456, mock.Anything).
		Return(&accountClientPackageV1.Account{Id: 0}, nil)

	err := core.CheckAccountIdExist(logrus.NewEntry(logrus.New()), 456, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account_id: 456 not found")

	accMock.AssertExpectations(t)
}

func TestCheckAccountIdExist_Error(t *testing.T) {
	core, _, _, accMock, db := setupTestCore(t)

	accMock.On("GetAccount", 789, mock.Anything).
		Return((*accountClientPackageV1.Account)(nil), errors.New("db error"))

	err := core.CheckAccountIdExist(logrus.NewEntry(logrus.New()), 789, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	accMock.AssertExpectations(t)
}

//-------------------------------------------//
// 5. Test: CreateTransaction
//-------------------------------------------//

func TestCreateTransaction_Success(t *testing.T) {
	core, repoMock, opMock, accMock, db := setupTestCore(t)

	payload := &entityCoreV1Package.CreateTransactionPayload{
		AccountId:       111,
		OperationTypeId: 2,
		Amount:          1000.0,
	}

	accMock.On("GetAccount", 111, mock.Anything).Return(&accountClientPackageV1.Account{Id: 111}, nil)

	opMock.On("GetOperationCoefficient", 2, mock.Anything).Return(1, nil)

	repoMock.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			tr := args.Get(0).(*entityDbV1Package.Transaction)
			assert.Equal(t, 1000.0, tr.Amount, "expected final transaction amount to be 1000")
			assert.Equal(t, 2, tr.OperationTypeId)
			assert.Equal(t, 111, tr.AccountId)
		})

	tx := db.Begin()
	defer tx.Rollback()

	transaction, err := core.CreateTransaction(logrus.NewEntry(logrus.New()), payload, tx)
	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, 1000.0, transaction.Amount)

	repoMock.AssertExpectations(t)
	opMock.AssertExpectations(t)
	accMock.AssertExpectations(t)
}

func TestCreateTransaction_AccountNotFound(t *testing.T) {
	core, repoMock, _, accMock, db := setupTestCore(t)

	payload := &entityCoreV1Package.CreateTransactionPayload{
		AccountId:       999,
		OperationTypeId: 2,
		Amount:          100.0,
	}

	accMock.On("GetAccount", 999, mock.Anything).
		Return(&accountClientPackageV1.Account{Id: 0}, nil)

	tx := db.Begin()
	defer tx.Rollback()

	transaction, err := core.CreateTransaction(logrus.NewEntry(logrus.New()), payload, tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "account_id: 999 not found")
	assert.Equal(t, 0.0, transaction.Amount, "expect zero transaction returned or partial data")

	repoMock.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything)

	accMock.AssertExpectations(t)
}

func TestCreateTransaction_FinalAmountError(t *testing.T) {
	core, repoMock, opMock, accMock, db := setupTestCore(t)

	payload := &entityCoreV1Package.CreateTransactionPayload{
		AccountId:       222,
		OperationTypeId: 3,
		Amount:          50.0,
	}

	accMock.On("GetAccount", 222, mock.Anything).
		Return(&accountClientPackageV1.Account{Id: 222}, nil)

	opMock.On("GetOperationCoefficient", 3, mock.Anything).
		Return(0.0, errors.New("coef error"))

	tx := db.Begin()
	defer tx.Rollback()

	_, err := core.CreateTransaction(logrus.NewEntry(logrus.New()), payload, tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "coef error")

	repoMock.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything)

	accMock.AssertExpectations(t)
	opMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
}

func TestCreateTransaction_RepoError(t *testing.T) {
	core, repoMock, opMock, accMock, db := setupTestCore(t)

	payload := &entityCoreV1Package.CreateTransactionPayload{
		AccountId:       333,
		OperationTypeId: 4,
		Amount:          250.0,
	}

	accMock.On("GetAccount", 333, mock.Anything).Return(&accountClientPackageV1.Account{Id: 333}, nil)
	opMock.On("GetOperationCoefficient", 4, mock.Anything).Return(1.0, nil)

	repoMock.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(errors.New("repo create error"))

	tx := db.Begin()
	defer tx.Rollback()

	_, err := core.CreateTransaction(logrus.NewEntry(logrus.New()), payload, tx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "repo create error")

	accMock.AssertExpectations(t)
	opMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
}
