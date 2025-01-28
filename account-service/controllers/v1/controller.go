package account_controller_v1

import (
	coreV1Package "anti-fraud/account-service/core/v1"
	entityHttpV1Package "anti-fraud/account-service/entity/http/v1"
	mapperV1Package "anti-fraud/account-service/mapper/v1"
	repoV1Package "anti-fraud/account-service/repository/v1"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"encoding/json"
	"net/http"
)

// IAccountController defines the contract for account-related HTTP handlers.
type IAccountController interface {
	// CreateAccount handles the creation of a new account.
	CreateAccount(w http.ResponseWriter, r *http.Request)
	// GetAccountDetails retrieves the details of an existing account by its ID.
	GetAccountDetails(w http.ResponseWriter, r *http.Request)
}

// AccountController provides HTTP handlers for account-related operations.
// It delegates business logic to IAccountCore and persistence tasks to IAccountRepository.
type AccountController struct {
	repoV1 repoV1Package.IAccountRepository // Repository for lower-level DB operations
	coreV1 coreV1Package.IAccountCore       // Core layer providing business logic.
	db     *gorm.DB                         // GORM DB instance for transaction management.
	logger *logrus.Logger                   // Logger for capturing logs and debug information.
}

// NewAccountController creates and returns a new AccountController initialized.
func NewAccountController(
	repoV1 repoV1Package.IAccountRepository,
	coreV1 coreV1Package.IAccountCore,
	db *gorm.DB,
	logger *logrus.Logger,
) *AccountController {
	return &AccountController{
		repoV1: repoV1,
		coreV1: coreV1,
		db:     db,
		logger: logger,
	}
}

// CreateAccount is an HTTP handler that creates a new account record in db.
//
// Workflow:
//  1. Decode the incoming JSON payload into CreateAccountRequest.
//  2. Validate the request payload.
//  3. Begin a database transaction.
//  4. Invoke the core layer to create the account (business logic).
//  5. Commit the transaction on success (or rollback on error, deferred).
//  6. Return a JSON response with the newly created account details.
func (controller *AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	controller.logger.Info("CreateAccount endpoint called")

	// 1. Decode JSON request body.
	var accountReq entityHttpV1Package.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&accountReq)
	if err != nil {
		controller.logger.Warningf("Error decoding request body: %v", err)
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Validate HTTP input payload.
	err = accountReq.Validate()
	if err != nil {
		controller.logger.Warningf("Validation failed: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Begin new db txn.
	tx := controller.db.Begin()
	defer tx.Rollback() // Rollback if we exit prematurely.

	// 4. Create account using the core layer’s business logic.
	accountPayload := mapperV1Package.CreateAccountPayloadMapper(&accountReq)
	account, err := controller.coreV1.CreateAccount(accountPayload, tx)
	if err != nil {
		controller.logger.Errorf("Error creating account: %v", err)
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Commit txn.
	if err := tx.Commit().Error; err != nil {
		controller.logger.Errorf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	controller.logger.Infof("Account created successfully: %v", account)

	// 6. Build and send JSON response.
	response := map[string]interface{}{
		"success": true,
		"account": mapperV1Package.AccountDetailsResponseMapper(account),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAccountDetails is an HTTP handler that retrieves account details by its ID.
//
// Workflow:
//  1. Extract the "accountId" from the URL path using Gorilla Mux.
//  2. Convert the ID to an integer.
//  3. Begin a database transaction.
//  4. Retrieve the account from the core layer.
//  5. Commit the transaction on success (or rollback on error, deferred).
//  6. Return a JSON response with the account details.
func (controller *AccountController) GetAccountDetails(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the "accountId" from URL params.
	params := mux.Vars(r)
	accountIdStr := params["accountId"]
	controller.logger.Infof("GetAccountDetails endpoint called for accountId: %v", accountIdStr)

	// 2. Convert the ID to an integer.
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		controller.logger.Errorf("Error converting string to int: %v", err)
		http.Error(w, "Error converting string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Begin a db txn.
	tx := controller.db.Begin()
	defer tx.Rollback() // Rollback if we exit prematurely.

	// 4. Fetch the account via core layer.
	account, err := controller.coreV1.GetAccount(accountId, tx)
	if err != nil {
		controller.logger.Errorf("Error fetching account details: %v", err)
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Commit txn.
	if err := tx.Commit().Error; err != nil {
		controller.logger.Errorf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Build and send the JSON response.
	response := map[string]interface{}{
		"success": true,
		"account": mapperV1Package.AccountDetailsResponseMapper(account),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
