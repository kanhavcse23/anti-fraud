package mediator_ops_client_v1

import (
	coreV1Package "anti-fraud/operation-service/core/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// IOperationClient defines methods interface for operation-related data
// via the operation core service.
type IOperationClient interface {
	// SetupCore injects the IOperationCore dependency.
	SetupCore(operationCoreV1 coreV1Package.IOperationCore)

	// GetOperationCoefficient fetches the coefficient for a specific operation ID.
	GetOperationCoefficient(logger *logrus.Entry, operationId int, tx *gorm.DB) (int, error)
}

// OperationClient implements IOperationClient, acting as a mediator to the operation core service.
type OperationClient struct {
	operationCoreV1 coreV1Package.IOperationCore
	logger          *logrus.Logger
}

// NewOperationClient create new instance of OperationClient.
func NewOperationClient(logger *logrus.Logger) *OperationClient {

	return &OperationClient{logger: logger}
}

// SetupCore injects the IOperationCore into this client.
func (client *OperationClient) SetupCore(operationCoreV1 coreV1Package.IOperationCore) {
	client.operationCoreV1 = operationCoreV1
}

// GetOperationCoefficient retrieves the coefficient for the given operation ID.
//
// Steps:
//  1. Delegate to operation core layer to fetch the coefficient from the DB.
//  2. Return the coefficient or an encountered error.
//
// Parameters:
//   - operationId: Unique identifier for the operation.
//   - tx:          db txn.
//
// Returns:
//   - int:   The coefficient associated with the operation ID.
//   - error: an encountered Error.
func (client *OperationClient) GetOperationCoefficient(logger *logrus.Entry, operationId int, tx *gorm.DB) (int, error) {
	logger.Info("GetOperationCoefficient method called in mediator-service for operation client.")
	coef, err := client.operationCoreV1.GetOperationCoefficient(logger, operationId, tx)
	if err != nil {
		logger.Errorf("Error occured while fetching coefficient associated on operationId %d via operation service: %s", operationId, err.Error())
	}
	return coef, err
}
