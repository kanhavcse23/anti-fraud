package transaction_controller_v1

import (
	coreV1Package "anti-fraud/transaction-service/core/v1"
	entityHttpV1Package "anti-fraud/transaction-service/entity/http/v1"
	mapperV1Package "anti-fraud/transaction-service/mapper/v1"
	repoV1Package "anti-fraud/transaction-service/repository/v1"

	"github.com/sirupsen/logrus"

	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type ITransactionController interface {
	CreateTransaction(w http.ResponseWriter, r *http.Request)
}
type TransactionController struct {
	repoV1 *repoV1Package.TransactionRepository
	coreV1 coreV1Package.ITransactionCore
	db     *gorm.DB
	logger *logrus.Logger
}

func NewTransactionController(repoV1 *repoV1Package.TransactionRepository, coreV1 coreV1Package.ITransactionCore, db *gorm.DB, logger *logrus.Logger) *TransactionController {
	return &TransactionController{repoV1: repoV1, coreV1: coreV1, db: db, logger: logger}
}

func (controller *TransactionController) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	controller.logger.Info("CreateTransaction endpoint called")
	var transactionReq entityHttpV1Package.CreateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&transactionReq)
	if err != nil {
		controller.logger.Warningf("Error decoding request body: %v", err)
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = transactionReq.Validate()
	if err != nil {
		controller.logger.Warningf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := controller.db.Begin()
	defer tx.Rollback()

	transaction, err := controller.coreV1.CreateTransaction(mapperV1Package.CreateTransactionPayloadMapper(&transactionReq), tx)
	if err != nil {
		controller.logger.Errorf("Error creating transaction: %v", err)
		http.Error(w, "An internal error occurred"+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Commit().Error; err != nil {
		controller.logger.Errorf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	controller.logger.Infof("Transaction created successfully: %v", transaction)
	response := map[string]interface{}{
		"success":     true,
		"transaction": mapperV1Package.TransactionDetailsResponseMapper(transaction),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
