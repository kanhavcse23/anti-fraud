package account_controller_v1

import (
	coreV1Package "anti-fraud/account-service/core/v1"
	entityHttpV1Package "anti-fraud/account-service/entity/http/v1"
	mapperV1Package "anti-fraud/account-service/mapper/v1"
	repoV1Package "anti-fraud/account-service/repository/v1"
	"strconv"

	"github.com/sirupsen/logrus"

	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AccountController struct {
	repoV1 *repoV1Package.AccountRepository
	coreV1 *coreV1Package.AccountCore
	db     *gorm.DB
	logger *logrus.Logger
}

func NewAccountController(repoV1 *repoV1Package.AccountRepository, coreV1 *coreV1Package.AccountCore, db *gorm.DB, logger *logrus.Logger) *AccountController {
	return &AccountController{repoV1: repoV1, coreV1: coreV1, db: db, logger: logger}
}

func (controller *AccountController) CreateAccount(w http.ResponseWriter, r *http.Request) {
	controller.logger.Info("CreateAccount endpoint called")
	var accountReq entityHttpV1Package.CreateAccountRequest
	err := json.NewDecoder(r.Body).Decode(&accountReq)
	if err != nil {
		controller.logger.Warningf("Error decoding request body: %v", err)
		http.Error(w, "Error decoding request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	err = accountReq.Validate()
	if err != nil {
		controller.logger.Warningf("Error: %v", err)
		http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	tx := controller.db.Begin()
	defer tx.Rollback()

	account, err := controller.coreV1.CreateAccount(mapperV1Package.CreateAccountPayloadMapper(&accountReq), tx)
	if err != nil {
		controller.logger.Errorf("Error creating account: %v", err)
		http.Error(w, "An internal error occurred"+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Commit().Error; err != nil {
		controller.logger.Errorf("Error committing transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	controller.logger.Infof("Account created successfully: %v", account)
	response := map[string]interface{}{
		"success": true,
		"account": mapperV1Package.AccountDetailsResponseMapper(account),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (controller *AccountController) GetAccountDetails(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	accountIdStr := params["accountId"]
	controller.logger.Infof("GetAccountDetails endpoint called for accountId: %v", accountIdStr)

	accountId, err := strconv.Atoi(accountIdStr)
	if err != nil {
		controller.logger.Errorf("Error converting string to int: %v\n", err)
		http.Error(w, "Error converting string to int: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tx := controller.db.Begin()
	defer tx.Rollback()

	account, err := controller.coreV1.GetAccount(accountId, tx)
	if err != nil {
		controller.logger.Errorf("Error fetching account details: %v", err)
		http.Error(w, "An internal error occurred"+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Commit().Error; err != nil {
		controller.logger.Errorf("Error commiting transaction: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"account": mapperV1Package.AccountDetailsResponseMapper(account),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
