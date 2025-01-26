package operation_middleware_v1

import (
	coreV1Path "anti-fraud/operation-service/core/v1"
	repoV1Path "anti-fraud/operation-service/repository/v1"

	clientPathV1 "anti-fraud/mediator-service/operation-service-client"

	"github.com/sirupsen/logrus"
)

type OperationMiddleware struct {
	logger *logrus.Logger
	coreV1 *coreV1Path.OperationCore
}

func NewOperationMiddleware(logger *logrus.Logger) *OperationMiddleware {

	return &OperationMiddleware{logger: logger}
}

func (mw *OperationMiddleware) Init() {
	repoV1 := repoV1Path.NewOperationRepository(mw.logger)
	mw.coreV1 = coreV1Path.NewOperationCore(repoV1, mw.logger)
}

func (mw *OperationMiddleware) ConfigureClient(client *clientPathV1.OperationClient) {
	client.SetupCore(mw.coreV1)
}
