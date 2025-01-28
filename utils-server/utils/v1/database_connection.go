package utils_db

import (
	configPackage "anti-fraud/utils-server/config"

	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func EstablishDBConnection() (*gorm.DB, error) {

	config, err := configPackage.LoadConfig()
	if err != nil {
		return nil, err
	}
	dbConfig := config.Database

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
