package operation_core_v1

import (
	repoV1Package "anti-fraud/operation-service/repository/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IOperationCore defines the core interface for operation-related business logic.
type IOperationCore interface {

	// GetOperationCoefficient returns the coefficient for a given operation ID.
	// If the operation does not exist or a repository error occurs, returns an error.
	GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error)
}

// OperationCore implements IOperationCore
type OperationCore struct {
	repoV1 repoV1Package.IOperationRepository
	logger *logrus.Logger
}

// NewOperationCore return new OperationCore instance.
func NewOperationCore(repoV1 repoV1Package.IOperationRepository, logger *logrus.Logger) *OperationCore {
	return &OperationCore{repoV1: repoV1, logger: logger}
}

// GetOperationCoefficient retrieves the coefficient for the specified operationId.
//
// Workflow:
//  1. Retrieves the operation record from the repository by operationId.
//  2. If the repository call returns an error (e.g., not found), return Error.
//  3. Returns the Coefficient field from the retrieved operation.
//
// Parameters:
//   - operationId: The unique id to find operation from db.
//   - tx:          db txn.
//
// Returns:
//   - int:   The coefficient associated with the operation (+1 / -1). +1 for credit txn, -1 for debit txn.
//   - error: an encountered Error.
func (core *OperationCore) GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error) {
	core.logger.Info("GetOperationCoefficient method called in operation core layer.")
	operation, err := core.repoV1.GetOperation(operationId, tx)
	if err != nil {
		return 0, err
	}
	return operation.Coefficient, nil
}
