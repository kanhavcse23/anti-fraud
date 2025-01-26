package mediator_ops_client_v1

import (
	controllerV1Path "anti-fraud/operation-service/core/v1"

	"gorm.io/gorm"
)

type OperationClient struct {
	operationCoreV1 *controllerV1Path.OperationCore
}

func NewOperationClient() *OperationClient {

	return &OperationClient{}
}
func (client *OperationClient) SetupCore(operationCoreV1 *controllerV1Path.OperationCore) {
	client.operationCoreV1 = operationCoreV1
}
func (client *OperationClient) GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error) {
	return client.operationCoreV1.GetOperationCoefficient(operationId, tx)
}
