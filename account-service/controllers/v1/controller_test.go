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

func (m *MockAccountCore) CreateAccount(payload *entityCoreV1Package.CreateAccountPayload, tx *gorm.DB) (*entityDbV1Package.Account, error) {
	args := m.Called(payload, tx)
	account, _ := args.Get(0).(*entityDbV1Package.Account)
	return account, args.Error(1)
}

func (m *MockAccountCore) GetAccount(accountId int, tx *gorm.DB) (*entityDbV1Package.Account, error) {
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
	// Arrange
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	// We don't need the repository here since the controller calls the core, not the repo directly
	controller := NewAccountController(nil, mockCore, db, logger)

	// Mock the core layer behavior
	mockCore.On("CreateAccount", mock.Anything, mock.Anything).
		Return(&entityDbV1Package.Account{
			Model:          gorm.Model{ID: 1},
			DocumentNumber: "123456789",
		}, nil)
	// Prepare HTTP request
	requestBody := `{"document_number":"123456789"}`
	req := httptest.NewRequest("POST", "/accounts/v1", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	controller.CreateAccount(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK")

	// Optionally, parse and check JSON response
	// (simple check just for demonstration)
	assert.Contains(t, rr.Body.String(), `"success":true`)

	mockCore.AssertExpectations(t)
}

func TestCreateAccount_InvalidPayload(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	// No need to mock CreateAccount; the request will fail before it hits the core

	requestBody := `{"document_number": 123}` // invalid type if your struct expects a string
	req := httptest.NewRequest("POST", "/accounts/v1", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	controller.CreateAccount(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error decoding request body")
	mockCore.AssertNotCalled(t, "CreateAccount", mock.Anything, mock.Anything)
}

func TestCreateAccount_CoreError(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	// Valid JSON but core returns error
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

	// Mock the core layer's GetAccount
	accountID := 1
	mockCore.On("GetAccount", accountID, mock.Anything).
		Return(&entityDbV1Package.Account{Model: gorm.Model{ID: 1}, DocumentNumber: "987654321"}, nil)

	// Create test request
	req := httptest.NewRequest("GET", "/accounts/v1/1", nil)
	rr := httptest.NewRecorder()
	// Simulate Gorilla Mux behavior by adding route variables to the request context
	vars := map[string]string{
		"accountId": "1",
	}
	req = mux.SetURLVars(req, vars)
	// Call handler directly
	controller.GetAccountDetails(rr, req)

	// Assert
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

	// Set up a request that has a non-integer path param
	req := httptest.NewRequest("GET", "/accounts/v1/abc", nil)
	rr := httptest.NewRecorder()
	// Simulate Gorilla Mux behavior by adding route variables to the request context
	vars := map[string]string{
		"accountId": "abc",
	}
	req = mux.SetURLVars(req, vars)

	controller.GetAccountDetails(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Error converting string to int")

	// The core shouldn't be called at all in this scenario
	mockCore.AssertNotCalled(t, "GetAccount", mock.Anything, mock.Anything)
}

func TestGetAccountDetails_CoreError(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()

	mockCore := new(MockAccountCore)
	controller := NewAccountController(nil, mockCore, db, logger)

	// Mock the core to return an error
	accountID := 5
	mockCore.On("GetAccount", accountID, mock.Anything).
		Return(nil, errors.New("core error"))

	req := httptest.NewRequest("GET", "/accounts/v1/5", nil)
	rr := httptest.NewRecorder()
	// Simulate Gorilla Mux behavior by adding route variables to the request context
	vars := map[string]string{
		"accountId": "5",
	}
	req = mux.SetURLVars(req, vars)

	controller.GetAccountDetails(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "An internal error occurred")

	mockCore.AssertExpectations(t)
}
