package account_controller_v1

import (
	entityCoreV1Package "anti-fraud/account-service/entity/core/v1"
	entityDbV1Package "anti-fraud/account-service/entity/db/v1"

	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//---------------------------//
//  1. Mock IAccountCore
//---------------------------//

type MockAccountCore struct {
	mock.Mock
}

func (m *MockAccountCore) CreateAccount(logger *logrus.Entry, payload *entityCoreV1Package.CreateAccountPayload, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	args := m.Called(payload, tx)
	account, _ := args.Get(0).(*entityDbV1Package.Account)
	return account, args.Error(1)
}

func (m *MockAccountCore) GetAccount(logger *logrus.Entry, accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	args := m.Called(accountId, tx)
	account, _ := args.Get(0).(*entityDbV1Package.Account)
	return account, args.Error(1)
}

//--------------------------------//
//  2. Helper: Create Test DB
//--------------------------------//

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	return db
}

//--------------------------------//
//  3. Tests for CreateAccount
//--------------------------------//

func TestCreateAccount_Success(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	mockCore.On("CreateAccount", mock.Anything, mock.Anything).
		Return(&entityDbV1Package.Account{
			Model:          gorm.Model{ID: 1},
			DocumentNumber: "123456789",
		}, nil)
	requestBody := `{"document_number":"123456789"}`
	req := httptest.NewRequest("POST", "/accounts/v1", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller.CreateAccount(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK")
	assert.Contains(t, rr.Body.String(), `"success":true`)

	mockCore.AssertExpectations(t)
}

func TestCreateAccount_InvalidPayload(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	requestBody := `{"document_number": 123}`
	req := httptest.NewRequest("POST", "/accounts/v1", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller.CreateAccount(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error decoding request body")
	mockCore.AssertNotCalled(t, "CreateAccount", mock.Anything, mock.Anything)
}

func TestCreateAccount_MissingDocumentNumber(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	requestBody := `{"document_number": ""}`
	req := httptest.NewRequest("POST", "/accounts/v1", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller.CreateAccount(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "document number should not be empty")
	mockCore.AssertNotCalled(t, "CreateAccount", mock.Anything, mock.Anything)
}

func TestCreateAccount_CommitError(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	requestBody := `{"document_number":"123456789"}`
	req := httptest.NewRequest("POST", "/accounts/v1", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockCore.On("CreateAccount", mock.Anything, mock.Anything).
		Return(&entityDbV1Package.Account{
			Model:          gorm.Model{ID: 1},
			DocumentNumber: "123456789",
		}, nil)

	tx := db.Begin()
	defer tx.Rollback()

	tx.AddError(errors.New("commit failed"))

	controller.db = tx
	controller.CreateAccount(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "commit failed")

	mockCore.AssertExpectations(t)
}
func TestCreateAccount_CoreError(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	requestBody := `{"document_number":"123456789"}`
	req := httptest.NewRequest("POST", "/accounts/v1", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	mockCore.On("CreateAccount", mock.Anything, mock.Anything).
		Return(nil, errors.New("some core error"))

	controller.CreateAccount(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "An internal error occurred")

	mockCore.AssertExpectations(t)
}

//-------------------------------------//
//  4. Tests for GetAccountDetails
//-------------------------------------//

func TestGetAccountDetails_Success(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	accountID := 1
	mockCore.On("GetAccount", accountID, mock.Anything).
		Return(&entityDbV1Package.Account{Model: gorm.Model{ID: 1}, DocumentNumber: "987654321"}, nil)

	req := httptest.NewRequest("GET", "/accounts/v1/1", nil)
	rr := httptest.NewRecorder()
	vars := map[string]string{
		"accountId": "1",
	}
	req = mux.SetURLVars(req, vars)
	controller.GetAccountDetails(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"success":true`)
	assert.Contains(t, rr.Body.String(), `"document_number":"987654321"`)

	mockCore.AssertExpectations(t)
}

func TestGetAccountDetails_InvalidPathParam(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	req := httptest.NewRequest("GET", "/accounts/v1/abc", nil)
	rr := httptest.NewRecorder()
	vars := map[string]string{
		"accountId": "abc",
	}
	req = mux.SetURLVars(req, vars)
	controller.GetAccountDetails(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error converting string to int")
	mockCore.AssertNotCalled(t, "GetAccount", mock.Anything, mock.Anything)
}

func TestGetAccountDetails_CoreError(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	accountID := 5
	mockCore.On("GetAccount", accountID, mock.Anything).
		Return(nil, errors.New("core error"))

	req := httptest.NewRequest("GET", "/accounts/v1/5", nil)
	rr := httptest.NewRecorder()
	vars := map[string]string{
		"accountId": "5",
	}
	req = mux.SetURLVars(req, vars)

	controller.GetAccountDetails(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "An internal error occurred")

	mockCore.AssertExpectations(t)
}
func TestGetAccountDetails_CommitError(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	accountID := 5
	mockCore.On("GetAccount", accountID, mock.Anything).
		Return(&entityDbV1Package.Account{Model: gorm.Model{ID: 1}, DocumentNumber: "987654321"}, nil)

	req := httptest.NewRequest("GET", "/accounts/v1/5", nil)
	rr := httptest.NewRecorder()
	vars := map[string]string{
		"accountId": "5",
	}
	req = mux.SetURLVars(req, vars)
	tx := db.Begin()
	defer tx.Rollback()

	tx.AddError(errors.New("commit failed"))

	controller.db = tx
	controller.GetAccountDetails(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "commit failed")

	mockCore.AssertExpectations(t)
}
