package account_controller_v1

import (
	coreV1Package "anti-fraud/account-service/core/v1"
	entityHttpV1Package "anti-fraud/account-service/entity/http/v1"
	mapperV1Package "anti-fraud/account-service/mapper/v1"
	repoV1Package "anti-fraud/account-service/repository/v1"
	utilV1 "anti-fraud/utils-server/middleware/v1"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"encoding/json"
	"net/http"
)

// IAccountController defines the methods interface for account-related HTTP handlers.
type IAccountController interface {

	// CreateAccount handles the creation of a new account.
	CreateAccount(w http.ResponseWriter, r *http.Request)

	// GetAccountDetails retrieves the details of an existing account by its ID.
	GetAccountDetails(w http.ResponseWriter, r *http.Request)
}

// AccountController implements IAccountController interface and
// provides HTTP handlers for account-related operations.
type AccountController struct {
	repoV1 repoV1Package.IAccountRepository // Repository for lower-level DB operations
	coreV1 coreV1Package.IAccountCore       // Core layer providing business logic.
	db     *gorm.DB
	logger *logrus.Logger
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
//  3. Begin db txn.
//  4. Invoke the core layer to create the account (business logic).
//  5. Commit the txn on success (or rollback on error).
//  6. Return a JSON response with the newly created account details.
func (controller *AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := utilV1.GetRequestID(ctx)
	logger := controller.logger.WithField("request_id", requestID)

	// 1. Decode JSON request body.
	var accountReq entityHttpV1Package.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&accountReq)
	if err != nil {
		logger.Errorf("Error decoding request body: %v", err)
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	logger.Infof("CreateAccount endpoint called input payload: %v", accountReq)

	// 2. Validate HTTP input payload.
	err = accountReq.Validate()
	if err != nil {
		logger.Errorf("Validation failed: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Begin new db txn.
	tx := controller.db.Begin()
	defer tx.Rollback() // Rollback if we exit prematurely.

	// 4. Create account using the core layer’s business logic.
	accountPayload := mapperV1Package.CreateAccountPayloadMapper(&accountReq)
	account, err := controller.coreV1.CreateAccount(logger, accountPayload, tx)
	if err != nil {
		logger.Errorf("Error creating account: %v", err)
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Commit txn.
	if err := tx.Commit().Error; err != nil {
		logger.Errorf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Infof("Account created successfully: %v", account)

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
//  2. Convert the accountId from string to an int.
//  3. Begin db txn.
//  4. Retrieve the account from the core layer.
//  5. Commit txn on success (or rollback on error).
//  6. Return a JSON response with the account details.
func (controller *AccountController) GetAccountDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := utilV1.GetRequestID(ctx)
	logger := controller.logger.WithField("request_id", requestID)
	// 1. Extract the "accountId" from URL params.
	params := mux.Vars(r)
	accountIdStr := params["accountId"]
	logger.Infof("GetAccountDetails endpoint called for accountId: %v", accountIdStr)

	// 2. Convert the ID to an integer.
	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		logger.Errorf("Error converting string to int: %v", err)
		http.Error(w, "Error converting string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Begin a db txn.
	tx := controller.db.Begin()
	defer tx.Rollback() // Rollback if we exit prematurely.

	// 4. Fetch the account via core layer.
	account, err := controller.coreV1.GetAccount(logger, accountId, tx)
	if err != nil {
		logger.Errorf("Error fetching account details: %v", err)
		http.Error(w, "An internal error occurred: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5. Commit txn.
	if err := tx.Commit().Error; err != nil {
		logger.Errorf("Error committing transaction: %v", err)
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
