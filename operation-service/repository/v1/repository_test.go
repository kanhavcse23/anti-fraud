package operation_repo_v1

import (
	constantPackage "anti-fraud/constants/operation"
	entityDbV1Package "anti-fraud/operation-service/entity/db/v1"
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

	if err := db.AutoMigrate(&entityDbV1Package.Operation{}); err != nil {
		t.Fatalf("failed to migrate Operation schema: %v", err)
	}

	return db
}

func TestGetOperation_Success(t *testing.T) {
	logger := logrus.New()
	repo := NewOperationRepository(logger)
	db := setupTestDB(t)

	op := entityDbV1Package.Operation{
		Coefficient: 5,
	}
	result := db.Table(constantPackage.TABLE_NAME).Create(&op)
	assert.NoError(t, result.Error, "failed to insert operation")

	found, err := repo.GetOperation(logrus.NewEntry(logrus.New()), int(op.ID), db)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, op.ID, found.ID)
	assert.Equal(t, 5, found.Coefficient)
}

func TestGetOperation_NotFound(t *testing.T) {
	logger := logrus.New()
	repo := NewOperationRepository(logger)
	db := setupTestDB(t)

	found, err := repo.GetOperation(logrus.NewEntry(logrus.New()), 999, db)

	assert.Error(t, err)
	assert.Equal(t, 0, int(found.ID))
	assert.Contains(t, err.Error(), "operation id: 999 not found in database")
}

func TestGetOperation_DBError(t *testing.T) {
	logger := logrus.New()
	repo := NewOperationRepository(logger)

	db := setupTestDB(t)
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	found, err := repo.GetOperation(logrus.NewEntry(logrus.New()), 1, db)
	assert.Error(t, err)
	assert.Equal(t, 0, int(found.ID))

}
