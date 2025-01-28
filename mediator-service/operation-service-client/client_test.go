package mediator_ops_client_v1

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

//-------------------------------------------//
// Mock for IOperationCore
//-------------------------------------------//

type MockOperationCore struct {
	mock.Mock
}

func (m *MockOperationCore) GetOperationCoefficient(operationId int, tx *gorm.DB) (int, error) {
	args := m.Called(operationId, tx)
	coef, _ := args.Get(0).(int)
	return coef, args.Error(1)
}

//-------------------------------------------//
// Unit Tests for OperationClient
//-------------------------------------------//

func TestOperationClient_GetOperationCoefficient_Success(t *testing.T) {
	client := NewOperationClient()
	mockCore := new(MockOperationCore)

	client.SetupCore(mockCore)

	mockCore.On("GetOperationCoefficient", 10, mock.Anything).
		Return(3, nil)

	coef, err := client.GetOperationCoefficient(10, &gorm.DB{})
	assert.NoError(t, err)
	assert.Equal(t, 3, coef)

	mockCore.AssertExpectations(t)
}

func TestOperationClient_GetOperationCoefficient_Error(t *testing.T) {
	client := NewOperationClient()
	mockCore := new(MockOperationCore)

	client.SetupCore(mockCore)

	mockCore.On("GetOperationCoefficient", 20, mock.Anything).
		Return(0, errors.New("some error"))

	coef, err := client.GetOperationCoefficient(20, &gorm.DB{})
	assert.Error(t, err)
	assert.Equal(t, 0, coef)

	mockCore.AssertExpectations(t)
}

func TestOperationClient_SetupCore(t *testing.T) {
	client := NewOperationClient()
	mockCore := new(MockOperationCore)
	client.SetupCore(mockCore)

	mockCore.On("GetOperationCoefficient", 1, mock.Anything).Return(5, nil)

	c, err := client.GetOperationCoefficient(1, &gorm.DB{})
	assert.NoError(t, err)
	assert.Equal(t, 5, c)

	mockCore.AssertExpectations(t)
}
