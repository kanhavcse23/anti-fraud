package operation_repo_v1

import (
	constantPackage "anti-fraud/constants/operation"
	entityDbV1Package "anti-fraud/operation-service/entity/db/v1"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OperationRepository struct {
	logger *logrus.Logger
}

func NewOperationRepository(logger *logrus.Logger) *OperationRepository {
	return &OperationRepository{logger: logger}
}

func (repo *OperationRepository) GetOperation(operationId int, tx *gorm.DB) (*entityDbV1Package.Operation, error) {
	repo.logger.Info("GetOperation method called in operation repo layer.")
	var operation entityDbV1Package.Operation
	result := tx.Table(constantPackage.TABLE_NAME).First(&operation, operationId)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		return &operation, fmt.Errorf("operation id: %d not found in database", operationId)
	}
	return &operation, result.Error
}
