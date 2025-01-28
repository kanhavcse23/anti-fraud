package operation_repo_v1

import (
	constantPackage "anti-fraud/constants/operation"
	entityDbV1Package "anti-fraud/operation-service/entity/db/v1"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IOperationRepository defines methods interface for performing database operations
// related to "operations".
type IOperationRepository interface {

	// GetOperation retrieves the operation record by its unique ID.
	GetOperation(operationId int, tx *gorm.DB) (*entityDbV1Package.Operation, error)
}

// OperationRepository implements IOperationRepository
type OperationRepository struct {
	logger *logrus.Logger
}

// NewOperationRepository creates new instance of OperationRepository
func NewOperationRepository(logger *logrus.Logger) *OperationRepository {
	return &OperationRepository{logger: logger}
}

// GetOperation finds an operation record based on the provided operationId.
//
// Steps:
//  1. Query table (constantPackage.TABLE_NAME) to
//     first record matching the given operationId (primary key).
//  3. If no record is found, returns a "not found" error.
//  4. Otherwise returns the operation record and any error encountered.
//
// Parameters:
//   - operationId: The unique identifier of the operation to retrieve.
//   - tx:          The gorm.DB transaction context to execute this query.
//
// Returns:
//   - A pointer to the retrieved Operation entity.
//   - An error if the record is not found or if any DB error occurs.
func (repo *OperationRepository) GetOperation(operationId int, tx *gorm.DB) (*entityDbV1Package.Operation, error) {
	repo.logger.Info("GetOperation method called in operation repo layer.")
	var operation entityDbV1Package.Operation
	result := tx.Table(constantPackage.TABLE_NAME).First(&operation, operationId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &operation, fmt.Errorf("operation id: %d not found in database", operationId)
	}
	return &operation, result.Error
}
