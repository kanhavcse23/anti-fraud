package transaction_controller_v1

import (
	coreV1Package "anti-fraud/transaction-service/core/v1"
	entityHttpV1Package "anti-fraud/transaction-service/entity/http/v1"
	mapperV1Package "anti-fraud/transaction-service/mapper/v1"
	repoV1Package "anti-fraud/transaction-service/repository/v1"

	utilV1 "anti-fraud/utils-server/middleware/v1"

	"github.com/sirupsen/logrus"

	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

// ITransactionController defines methods interface for HTTP handler.
type ITransactionController interface {

	// CreateTransaction handles an HTTP request to create a new transaction.
	CreateTransaction(w http.ResponseWriter, r *http.Request)
}

// TransactionController implements ITransactionController interface.
type TransactionController struct {
	repoV1 repoV1Package.ITransactionRepository
	coreV1 coreV1Package.ITransactionCore
	db     *gorm.DB
	logger *logrus.Logger
}

// NewTransactionController creates and returns new TransactionController instance.
func NewTransactionController(repoV1 repoV1Package.ITransactionRepository, coreV1 coreV1Package.ITransactionCore, db *gorm.DB, logger *logrus.Logger) *TransactionController {
	return &TransactionController{repoV1: repoV1, coreV1: coreV1, db: db, logger: logger}
}

// CreateTransaction handles the HTTP request for creating a new transaction.
//
// Workflow:
//  1. Parse the JSON request body into a CreateTransactionRequest struct.
//  2. Validate the request data.
//  3. Start a new db txn.
//  4. Delegate to the core layer to create the transaction (business logic).
//  5. Commit db txn.
//  6. Return http response with the newly created transaction.
func (controller *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestID := utilV1.GetRequestID(ctx)
	logger := controller.logger.WithField("request_id", requestID)
	logger.Info("CreateTransaction endpoint called")
	var transactionReq entityHttpV1Package.CreateTransactionRequest

	// 1. Decode HTTP input payload.
	err := json.NewDecoder(r.Body).Decode(&transactionReq)
	if err != nil {
		logger.Errorf("Error decoding request body: %v", err)
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2. Validate payload.
	err = transactionReq.Validate()
	if err != nil {
		logger.Errorf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 3. Begin db txn.
	tx := controller.db.Begin()
	defer tx.Rollback()

	// 4. Create a new transaction via the core layer.
	transaction, err := controller.coreV1.CreateTransaction(logger, mapperV1Package.CreateTransactionPayloadMapper(&transactionReq), tx)
	if err != nil {
		logger.Errorf("Error creating transaction: %v", err)
		http.Error(w, "An internal error occurred"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Commit db txn
	if err := tx.Commit().Error; err != nil {
		logger.Errorf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Build and send http response.
	logger.Infof("Transaction created successfully: %v", transaction)
	response := map[string]interface{}{
		"success":     true,
		"transaction": mapperV1Package.TransactionDetailsResponseMapper(transaction),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
