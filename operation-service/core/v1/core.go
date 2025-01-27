package operation_core_v1

import (
	repoV1Package "anti-fraud/operation-service/repository/v1"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IOperationCore interface {
	GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error)
}
type OperationCore struct {
	repoV1 repoV1Package.IOperationRepository
	logger *logrus.Logger
}

func NewOperationCore(repoV1 repoV1Package.IOperationRepository, logger *logrus.Logger) *OperationCore {
	return &OperationCore{repoV1: repoV1, logger: logger}
}

func (core *OperationCore) GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error) {
	core.logger.Info("GetOperationCoefficient method called in operation core layer.")
	operation, err := core.repoV1.GetOperation(operationId, tx)
	if err != nil {
		return 0, err
	}
	return operation.Coefficient, nil
}
