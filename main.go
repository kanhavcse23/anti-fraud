package main

import (
	account_middleware_v1 "anti-fraud/account-service/middleware/v1"
	operation_middleware_v1 "anti-fraud/operation-service/middleware/v1"
	transaction_middleware_v1 "anti-fraud/transaction-service/middleware/v1"
	"fmt"
	"log"
	"net/http"

	operationClientV1Package "anti-fraud/mediator-service/operation-service-client"
	dbConnPackage "anti-fraud/utils-server/utils/v1"

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
	//configure account service
	accountMiddlewareV1 := account_middleware_v1.NewAccountMiddleware(db, router, logger)
	accountMiddlewareV1.Init()

	//configure transaction service
	transactionMiddlewareV1 := transaction_middleware_v1.NewTransactionMiddleware(db, router, logger, operationClient)
	transactionMiddlewareV1.Init()

	//configure operation service
	operationMiddlewareV1 := operation_middleware_v1.NewOperationMiddleware(logger)
	operationMiddlewareV1.Init()
	operationMiddlewareV1.ConfigureClient(operationClient)

	logger.Errorf("All components has been wired.")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
