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

	// Make sure your Account entity is migrated to the DB
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

	// Create a sample account
	acc := &entityDbV1Package.Account{
		DocumentNumber: "123456789",
	}

	err := repo.CreateAccount(acc, db)
	assert.NoError(t, err, "expected no error on account creation")
	assert.NotZero(t, acc.ID, "expected the newly created account to have a non-zero ID")

	// Verify in the DB that the record is actually created
	var found entityDbV1Package.Account
	result := db.Table(constantPackage.TABLE_NAME).First(&found, acc.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, acc.DocumentNumber, found.DocumentNumber)
}

func TestGetAccount_Success(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	// Create and insert an account manually
	acc := &entityDbV1Package.Account{
		DocumentNumber: "987654321",
	}
	db.Create(acc)

	// Retrieve via repository
	got, err := repo.GetAccount(int(acc.ID), db)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, acc.ID, got.ID)
	assert.Equal(t, "987654321", got.DocumentNumber)
}

func TestGetAccount_NotFound(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	// Attempt to get an account that doesn't exist
	got, err := repo.GetAccount(9999, db)

	// The repository code returns nil error for ErrRecordNotFound,
	// but an empty struct is returned. Let's verify that logic.
	assert.NoError(t, err, "repository returns nil error for not found")
	assert.NotNil(t, got)
	assert.Equal(t, uint(0), got.ID, "ID should be zero when not found")
}

func TestCheckDuplicateAccount_NoDuplicate(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	// There's no account with DocumentNumber "unique" yet
	got, err := repo.CheckDuplicateAccount("unique", db)
	assert.NoError(t, err, "no error should occur if record not found")
	assert.Equal(t, uint(0), got.ID, "ID should be zero if not found")
}

func TestCheckDuplicateAccount_Found(t *testing.T) {
	db := setupTestDB(t)
	logger := logrus.New()
	repo := NewAccountRepository(logger)

	// Create an account with "duplicate" document number
	acc := &entityDbV1Package.Account{
		DocumentNumber: "duplicate",
	}
	db.Create(acc)

	// Check duplicate
	got, err := repo.CheckDuplicateAccount("duplicate", db)
	assert.NoError(t, err)
	assert.NotNil(t, got)
	assert.Equal(t, acc.ID, got.ID)
	assert.Equal(t, "duplicate", got.DocumentNumber)
}
