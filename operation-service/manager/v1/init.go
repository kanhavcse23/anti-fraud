package operation_manager_v1

import (
	clientV1Package "anti-fraud/mediator-service/operation-service-client"
	coreV1Package "anti-fraud/operation-service/core/v1"
	repoV1Package "anti-fraud/operation-service/repository/v1"

	"github.com/sirupsen/logrus"
)

// OperationManager wires all components required to run account-service.
type OperationManager struct {
	logger *logrus.Logger
	coreV1 coreV1Package.IOperationCore
}

// NewOperationManager create and return new instance of OperationManager.
func NewOperationManager(logger *logrus.Logger) *OperationManager {

	return &OperationManager{logger: logger}
}

// Init wire all components, register routes for operation-service.
func (mw *OperationManager) Init() {
	repoV1 := repoV1Package.NewOperationRepository(mw.logger)
	mw.coreV1 = coreV1Package.NewOperationCore(repoV1, mw.logger)
}

// ConfigureClient configure core instance of operation service in operation-client.
func (mw *OperationManager) ConfigureClient(client clientV1Package.IOperationClient) {
	client.SetupCore(mw.coreV1)
}
