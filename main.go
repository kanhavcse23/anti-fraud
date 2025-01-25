package main

import (
	account_middleware_v1 "anti-fraud/account-service/middleware/v1"
	transaction_middleware_v1 "anti-fraud/transaction-service/middleware/v1"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	router := mux.NewRouter()
	dsn := "host=localhost user=postgres password=postgres dbname=account port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	// TODO: Improve. should be better than middleware
	accountMiddlewareV1 := account_middleware_v1.NewAccountMiddleware(db, router)
	accountMiddlewareV1.Init()

	// MongoDB connection URI
	uri := "mongodb://localhost:27017"

	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}

	mongoDb := mongoClient.Database("transaction")
	transactionMiddlewareV1 := transaction_middleware_v1.NewTransactionMiddleware(mongoDb)
	transactionMiddlewareV1.Init()

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
