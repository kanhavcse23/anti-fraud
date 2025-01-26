package main

import (
	account_middleware_v1 "anti-fraud/account-service/middleware/v1"
	transaction_middleware_v1 "anti-fraud/transaction-service/middleware/v1"
	"fmt"
	"log"
	"net/http"

	dbConnPath "anti-fraud/utils-server/utils/v1"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	router := mux.NewRouter()

	db, err := dbConnPath.EstablishPostgresqlDBConnection("account")
	if err != nil {
		log.Fatal(err)
	}

	accountMiddlewareV1 := account_middleware_v1.NewAccountMiddleware(db, router, logger)
	accountMiddlewareV1.Init()

	db2, err := dbConnPath.EstablishMongoDBConnection("transaction")
	if err != nil {
		log.Fatal(err)
	}
	transactionMiddlewareV1 := transaction_middleware_v1.NewTransactionMiddleware(db2)
	transactionMiddlewareV1.Init()

	logger.Errorf("All components has been wired.")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
