package transaction_controller_v1

import (
	entityCoreV1Package "anti-fraud/transaction-service/entity/core/v1"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"
	entityHttpV1Package "anti-fraud/transaction-service/entity/http/v1"

	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//--------------------------------//
// Mock for ITransactionCore
//--------------------------------//

type MockTransactionCore struct {
	mock.Mock
}

func (m *MockTransactionCore) CreateTransaction(logger *logrus.Entry, payload *entityCoreV1Package.CreateTransactionPayload, tx *gorm.DB) (*entityDbV1Package.Transaction, error) {
	args := m.Called(payload, tx)
	transaction, _ := args.Get(0).(*entityDbV1Package.Transaction)
	return transaction, args.Error(1)
}

func (m *MockTransactionCore) FinalTransactionAmount(logger *logrus.Entry, amount float64, operationTypeID int, tx *gorm.DB) (float64, error) {
	args := m.Called(amount, operationTypeID, tx)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockTransactionCore) CheckAccountIdExist(logger *logrus.Entry, accountId int, tx *gorm.DB) error {
	args := m.Called(accountId, tx)
	return args.Error(0)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}
	return db
}

//----------------------------------------------//
// Test Helpers
//----------------------------------------------//

func setupTestController(t *testing.T) (*TransactionController, *MockTransactionCore, *gorm.DB) {
	logger := logrus.New()
	db := setupTestDB(t)
	mockCore := new(MockTransactionCore)

	controller := NewTransactionController(
		nil,
		mockCore,
		db,
		logger,
	)

	return controller, mockCore, db
}

//------------------------------------------------//
// 1) TestCreateTransaction_Success
//------------------------------------------------//

func TestCreateTransaction_Success(t *testing.T) {
	controller, mockCore, _ := setupTestController(t)

	validPayload := entityHttpV1Package.CreateTransactionRequest{
		AccountId:       123,
		OperationTypeId: 1,
		Amount:          200.0,
	}
	bodyBytes, _ := json.Marshal(validPayload)

	req := httptest.NewRequest(http.MethodPost, "/transactions/v1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockCore.On(
		"CreateTransaction",
		mock.Anything,
		mock.Anything,
	).Return(&entityDbV1Package.Transaction{
		Model:           gorm.Model{ID: 1},
		AccountId:       123,
		OperationTypeId: 1,
		Amount:          200.0,
	}, nil)

	controller.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"success":true`)
	assert.Contains(t, rr.Body.String(), `"amount":200`)

	mockCore.AssertExpectations(t)
}

//------------------------------------------------//
// 2) TestCreateTransaction_BadJSON
//------------------------------------------------//

func TestCreateTransaction_BadJSON(t *testing.T) {
	controller, mockCore, _ := setupTestController(t)

	req := httptest.NewRequest(http.MethodPost, "/transactions/v1", bytes.NewReader([]byte(`{invalid-json`)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error decoding request body")

	mockCore.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything)
}

//------------------------------------------------//
// 3) TestCreateTransaction_ValidationError
//------------------------------------------------//

func TestCreateTransaction_ValidationError(t *testing.T) {
	controller, mockCore, _ := setupTestController(t)

	payload := entityHttpV1Package.CreateTransactionRequest{
		AccountId: 123,
		Amount:    1000.0,
	}
	bodyBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/transactions/v1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error:")

	mockCore.AssertNotCalled(t, "CreateTransaction", mock.Anything, mock.Anything)
}

//------------------------------------------------//
// 4) TestCreateTransaction_CoreError
//------------------------------------------------//

func TestCreateTransaction_CoreError(t *testing.T) {
	controller, mockCore, _ := setupTestController(t)

	payload := entityHttpV1Package.CreateTransactionRequest{
		AccountId:       123,
		OperationTypeId: 1,
		Amount:          500,
	}
	bodyBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/transactions/v1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockCore.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(nil, errors.New("some core error"))

	controller.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "An internal error occurred")

	mockCore.AssertExpectations(t)
}

// ------------------------------------------------//
// 5) TestCreateTransaction_CommitError
// ------------------------------------------------//
func TestCreateTransaction_CommitError(t *testing.T) {
	controller, mockCore, db := setupTestController(t)

	payload := entityHttpV1Package.CreateTransactionRequest{
		AccountId:       123,
		OperationTypeId: 1,
		Amount:          500,
	}
	bodyBytes, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/transactions/v1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockCore.On("CreateTransaction", mock.Anything, mock.Anything).
		Return(&entityDbV1Package.Transaction{Model: gorm.Model{ID: 1}, Amount: 500}, nil)

	tx := db.Begin()
	defer tx.Rollback()
	tx.AddError(errors.New("commit failed"))

	controller.db = tx

	controller.CreateTransaction(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "commit failed")

	mockCore.AssertExpectations(t)
}
