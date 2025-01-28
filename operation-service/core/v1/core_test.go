package operation_core_v1

import (
	entityDbV1Package "anti-fraud/operation-service/entity/db/v1"

	"errors"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockOperationRepository struct {
	mock.Mock
}

func (m *MockOperationRepository) GetOperation(logger *logrus.Entry, operationId int, tx *gorm.DB) (*entityDbV1Package.Operation, error) {
	args := m.Called(operationId, tx)
	op, _ := args.Get(0).(*entityDbV1Package.Operation)
	return op, args.Error(1)
}

func setupTestCore() (*OperationCore, *MockOperationRepository) {
	logger := logrus.New()
	mockRepo := new(MockOperationRepository)

	core := NewOperationCore(mockRepo, logger)

	return core, mockRepo
}

//---------------------------//
// Tests
//---------------------------//

func TestGetOperationCoefficient_Success(t *testing.T) {
	core, mockRepo := setupTestCore()

	mockRepo.On("GetOperation", 1, mock.Anything).
		Return(&entityDbV1Package.Operation{Model: gorm.Model{ID: 1}, Coefficient: 3}, nil)

	coef, err := core.GetOperationCoefficient(logrus.NewEntry(logrus.New()), 1, &gorm.DB{})
	assert.NoError(t, err)
	assert.Equal(t, 3, coef)

	mockRepo.AssertExpectations(t)
}

func TestGetOperationCoefficient_NotFoundError(t *testing.T) {
	core, mockRepo := setupTestCore()

	mockRepo.On("GetOperation", 99, mock.Anything).
		Return((*entityDbV1Package.Operation)(nil), errors.New("operation id: 99 not found in database"))

	coef, err := core.GetOperationCoefficient(logrus.NewEntry(logrus.New()), 99, &gorm.DB{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "99 not found")
	assert.Equal(t, 0, coef)

	mockRepo.AssertExpectations(t)
}

func TestGetOperationCoefficient_RepoError(t *testing.T) {
	core, mockRepo := setupTestCore()

	mockRepo.On("GetOperation", 2, mock.Anything).
		Return((*entityDbV1Package.Operation)(nil), errors.New("db error"))

	coef, err := core.GetOperationCoefficient(logrus.NewEntry(logrus.New()), 2, &gorm.DB{})
	assert.Error(t, err)
	assert.Equal(t, 0, coef)
	assert.Contains(t, err.Error(), "db error")

	mockRepo.AssertExpectations(t)
}
