package main

import (
	account_manager_v1 "anti-fraud/account-service/manager/v1"
	operation_manager_v1 "anti-fraud/operation-service/manager/v1"
	transaction_manager_v1 "anti-fraud/transaction-service/manager/v1"
	"fmt"
	"log"
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

	db, err := dbConnPackage.EstablishPostgresqlDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	// Configure Client

	//Operation Client
	operationClient := operationClientV1Package.NewOperationClient()
	//Operation Client
	accountClient := accountClientV1Package.NewAccountClient()
	//configure account service
	accountManagerV1 := account_manager_v1.NewAccountManager(db, router, logger)
	accountManagerV1.Init()
	accountManagerV1.ConfigureClient(accountClient)

	//configure transaction service
	transactionManagerV1 := transaction_manager_v1.NewTransactionManager(db, router, logger, operationClient, accountClient)
	transactionManagerV1.Init()

	//configure operation service
	operationManagerV1 := operation_manager_v1.NewOperationManager(logger)
	operationManagerV1.Init()
	operationManagerV1.ConfigureClient(operationClient)

	logger.Info("All components has been wired.")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
