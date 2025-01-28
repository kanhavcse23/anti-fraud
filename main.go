package main

import (
	account_manager_v1 "anti-fraud/account-service/manager/v1"
	operation_manager_v1 "anti-fraud/operation-service/manager/v1"
	transaction_manager_v1 "anti-fraud/transaction-service/manager/v1"
	"net/http"

	operationClientV1Package "anti-fraud/mediator-service/operation-service-client"
	dbConnPackage "anti-fraud/utils-server/utils/v1"

	accountClientV1Package "anti-fraud/mediator-service/account-service-client"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	router := mux.NewRouter()

	// Establish db connection
	db, err := dbConnPackage.EstablishDBConnection()
	if err != nil {
		logger.Fatalf("Error: %v", err)
	}

	// Operation Client
	operationClient := operationClientV1Package.NewOperationClient(logger)

	// Account Client
	accountClient := accountClientV1Package.NewAccountClient(logger)

	// Account Service
	accountManagerV1 := account_manager_v1.NewAccountManager(db, router, logger)
	accountManagerV1.Init()
	accountManagerV1.ConfigureClient(accountClient)

	// Transaction Service
	transactionManagerV1 := transaction_manager_v1.NewTransactionManager(db, router, logger, operationClient, accountClient)
	transactionManagerV1.Init()

	// Operation Service
	operationManagerV1 := operation_manager_v1.NewOperationManager(logger)
	operationManagerV1.Init()
	operationManagerV1.ConfigureClient(operationClient)

	logger.Info("All components has been wired.")

	if err := http.ListenAndServe(":8080", router); err != nil {
		logger.Fatalf("Failed to start server: %v\n", err)
	}
}
