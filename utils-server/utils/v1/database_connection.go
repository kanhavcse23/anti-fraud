package utils_db

import (
	configPath "anti-fraud/utils-server/config"
	"errors"

	"fmt"

	"context"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func EstablishPostgresqlDBConnection(serviceName string) (*gorm.DB, error) {

	config, err := configPath.LoadConfig()
	if err != nil {
		return nil, err
	}
	dbConfig, ok := config.Database[serviceName]
	if !ok {
		return nil, errors.New("DB Config is not configured")
	}
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port, dbConfig.SSLMode, dbConfig.TimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Use singular table names
		},
	})

	return db, err
}

func EstablishMongoDBConnection(serviceName string) (*mongo.Database, error) {
	config, err := configPath.LoadConfig()
	if err != nil {
		return nil, err
	}
	dbConfig, ok := config.Database[serviceName]
	if !ok {
		return nil, errors.New("DB Config is not configured")
	}

	clientOptions := options.Client().ApplyURI(dbConfig.Uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("MongoDB connection failed: %v", err))
	}

	db := mongoClient.Database(dbConfig.DBName)
	return db, nil

}
