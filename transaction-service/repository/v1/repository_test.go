package transaction_repo_v1

import (
	constantPackage "anti-fraud/constants/transaction"
	entityDbV1Package "anti-fraud/transaction-service/entity/db/v1"

	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory DB: %v", err)
	}
	err = db.AutoMigrate(&entityDbV1Package.Transaction{})
	if err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func TestCreateTransaction_Success(t *testing.T) {
	logger := logrus.New()
	repo := NewTransactionRepository(logger)
	db := setupTestDB(t)

	txModel := &entityDbV1Package.Transaction{
		AccountId:       123,
		OperationTypeId: 1,
		Amount:          100.50,
	}

	err := repo.CreateTransaction(logrus.NewEntry(logrus.New()), txModel, db)
	assert.NoError(t, err, "expected no error inserting transaction")
	assert.NotZero(t, txModel.ID, "expected transaction to have a generated ID after insert")

	var found entityDbV1Package.Transaction
	result := db.Table(constantPackage.TABLE_NAME).First(&found, txModel.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, txModel.AccountId, found.AccountId)
	assert.Equal(t, txModel.Amount, found.Amount)
}

func TestCreateTransaction_DBError(t *testing.T) {
	logger := logrus.New()
	repo := NewTransactionRepository(logger)

	db := setupTestDB(t)
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	txModel := &entityDbV1Package.Transaction{
		AccountId:       999,
		OperationTypeId: 2,
		Amount:          500,
	}

	err = repo.CreateTransaction(logrus.NewEntry(logrus.New()), txModel, db)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database is closed")
}
