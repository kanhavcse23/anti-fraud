package operation_middleware_v1

import (
	coreV1Package "anti-fraud/operation-service/core/v1"
	repoV1Package "anti-fraud/operation-service/repository/v1"

	clientV1Package "anti-fraud/mediator-service/operation-service-client"

	"github.com/sirupsen/logrus"
)

type OperationMiddleware struct {
	logger *logrus.Logger
	coreV1 *coreV1Package.OperationCore
}

func NewOperationMiddleware(logger *logrus.Logger) *OperationMiddleware {

	return &OperationMiddleware{logger: logger}
}

func (mw *OperationMiddleware) Init() {
	repoV1 := repoV1Package.NewOperationRepository(mw.logger)
	mw.coreV1 = coreV1Package.NewOperationCore(repoV1, mw.logger)
}

func (mw *OperationMiddleware) ConfigureClient(client *clientV1Package.OperationClient) {
	client.SetupCore(mw.coreV1)
}
