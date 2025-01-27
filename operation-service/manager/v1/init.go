package operation_manager_v1

import (
	coreV1Package "anti-fraud/operation-service/core/v1"
	repoV1Package "anti-fraud/operation-service/repository/v1"

	clientV1Package "anti-fraud/mediator-service/operation-service-client"

	"github.com/sirupsen/logrus"
)

type OperationManager struct {
	logger *logrus.Logger
	coreV1 *coreV1Package.OperationCore
}

func NewOperationManager(logger *logrus.Logger) *OperationManager {

	return &OperationManager{logger: logger}
}

func (mw *OperationManager) Init() {
	repoV1 := repoV1Package.NewOperationRepository(mw.logger)
	mw.coreV1 = coreV1Package.NewOperationCore(repoV1, mw.logger)
}

func (mw *OperationManager) ConfigureClient(client *clientV1Package.OperationClient) {
	client.SetupCore(mw.coreV1)
}
