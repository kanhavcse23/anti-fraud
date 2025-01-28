package operation_core_v1

import (
	repoV1Package "anti-fraud/operation-service/repository/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IOperationCore defines the methods interface for operation-related core business logic.
type IOperationCore interface {

	// GetOperationCoefficient returns the coefficient for a given operation ID.
	GetOperationCoefficient(logger *logrus.Entry, operationId int, tx *gorm.DB) (int, error)
}

// OperationCore implements IOperationCore interface.
type OperationCore struct {
	repoV1 repoV1Package.IOperationRepository
	logger *logrus.Logger
}

// NewOperationCore creates and return new OperationCore instance.
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
//   - int:   The coefficient associated with the operation type.
//   - error: an encountered Error.
func (core *OperationCore) GetOperationCoefficient(logger *logrus.Entry, operationId int, tx *gorm.DB) (int, error) {
	logger.Info("GetOperationCoefficient method called in operation core layer.")
	operation, err := core.repoV1.GetOperation(logger, operationId, tx)
	if err != nil {
		logger.Errorf("Error occured while fetching coefficient associated with operationId (%d) : %s", operationId, err.Error())
		return 0, err
	}
	return operation.Coefficient, nil
}
