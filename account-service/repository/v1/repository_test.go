package account_repo_v1

import (
	"testing"

	entityDbV1Package "anti-fraud/account-service/entity/db/v1"
	constantPackage "anti-fraud/constants/account"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database, applies migrations,
// and returns a gorm.DB pointer for testing.
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to in-memory database: %v", err)
	}

	err = db.AutoMigrate(&entityDbV1Package.Account{})
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestCreateAccount_Success(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	acc := &entityDbV1Package.Account{
		DocumentNumber: "123456789",
	}

	err := repo.CreateAccount(logrus.NewEntry(logrus.New()), acc, db)
	assert.NoError(t, err, "expected no error on account creation")
	assert.NotZero(t, acc.ID, "expected the newly created account to have a non-zero ID")

	var found entityDbV1Package.Account
	result := db.Table(constantPackage.TABLE_NAME).First(&found, acc.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, acc.DocumentNumber, found.DocumentNumber)
}

func TestGetAccount_Success(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	acc := &entityDbV1Package.Account{
		DocumentNumber: "987654321",
	}
	db.Create(acc)

	got, err := repo.GetAccount(logrus.NewEntry(logrus.New()), int(acc.ID), db)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, acc.ID, got.ID)
	assert.Equal(t, "987654321", got.DocumentNumber)
}

func TestGetAccount_NotFound(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	got, err := repo.GetAccount(logrus.NewEntry(logrus.New()), 9999, db)

	assert.NoError(t, err, "repository returns nil error for not found")
	assert.NotNil(t, got)
	assert.Equal(t, uint(0), got.ID, "ID should be zero when not found")
}

func TestCheckDuplicateAccount_NoDuplicate(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	got, err := repo.CheckDuplicateAccount(logrus.NewEntry(logrus.New()), "unique", db)
	assert.NoError(t, err, "no error should occur if record not found")
	assert.Equal(t, uint(0), got.ID, "ID should be zero if not found")
}

func TestCheckDuplicateAccount_Found(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	acc := &entityDbV1Package.Account{
		DocumentNumber: "duplicate",
	}
	db.Create(acc)

	got, err := repo.CheckDuplicateAccount(logrus.NewEntry(logrus.New()), "duplicate", db)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, acc.ID, got.ID)
	assert.Equal(t, "duplicate", got.DocumentNumber)
}
