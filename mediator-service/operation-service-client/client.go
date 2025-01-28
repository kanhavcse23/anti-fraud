package mediator_ops_client_v1

import (
	coreV1Package "anti-fraud/operation-service/core/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IOperationClient interface {
	SetupCore(operationCoreV1 coreV1Package.IOperationCore)
	GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error)
}
type OperationClient struct {
	operationCoreV1 coreV1Package.IOperationCore
	logger          *logrus.Logger
}

func NewOperationClient(logger *logrus.Logger) *OperationClient {

	return &OperationClient{logger: logger}
}
func (client *OperationClient) SetupCore(operationCoreV1 coreV1Package.IOperationCore) {
	client.operationCoreV1 = operationCoreV1
}
func (client *OperationClient) GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error) {
	client.logger.Info("GetOperationCoefficient method called in mediator-service for operation client.")
	return client.operationCoreV1.GetOperationCoefficient(operationId, tx)
}
