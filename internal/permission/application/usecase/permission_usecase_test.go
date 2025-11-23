package permission_usecase

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
	permission_dtos "github.com/williamkoller/system-education/internal/permission/presentation/dtos"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Save(p *permission_entity.Permission) (*permission_entity.Permission, error) {
	args := m.Called(p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission_entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindAll() ([]*permission_entity.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission_entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindPermissionByUserID(userID string) ([]*permission_entity.Permission, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permission_entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Update(id string, p *permission_entity.Permission) (*permission_entity.Permission, error) {
	args := m.Called(id, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission_entity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestPermissionUsecase_Create(t *testing.T) {
	t.Run("should create permission successfully", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		input := permission_dtos.AddPermissionDto{
			UserID:      "user-1",
			Modules:     []string{"module1"},
			Actions:     []string{"read"},
			Level:       "admin",
			Description: "test permission",
		}

		mockRepo.On("Save", mock.AnythingOfType("*permission_entity.Permission")).Return(&permission_entity.Permission{
			ID:          "123",
			UserID:      input.UserID,
			Modules:     input.Modules,
			Actions:     input.Actions,
			Level:       input.Level,
			Description: input.Description,
		}, nil)

		permission, err := usecase.Create(input)

		assert.NoError(t, err)
		assert.NotNil(t, permission)
		assert.Equal(t, input.UserID, permission.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when validation fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		input := permission_dtos.AddPermissionDto{
			UserID: "", // Invalid
		}

		permission, err := usecase.Create(input)

		assert.Error(t, err)
		assert.Nil(t, permission)
		mockRepo.AssertNotCalled(t, "Save")
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		input := permission_dtos.AddPermissionDto{
			UserID:      "user-1",
			Modules:     []string{"module1"},
			Actions:     []string{"read"},
			Level:       "admin",
			Description: "test permission",
		}

		mockRepo.On("Save", mock.AnythingOfType("*permission_entity.Permission")).Return(nil, errors.New("db error"))

		permission, err := usecase.Create(input)

		assert.Error(t, err)
		assert.Nil(t, permission)
		mockRepo.AssertExpectations(t)
	})
}

func TestPermissionUsecase_FindAll(t *testing.T) {
	t.Run("should return all permissions", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		expectedPermissions := []*permission_entity.Permission{
			{ID: "1", UserID: "user-1"},
			{ID: "2", UserID: "user-2"},
		}

		mockRepo.On("FindAll").Return(expectedPermissions, nil)

		permissions, err := usecase.FindAll()

		assert.NoError(t, err)
		assert.Equal(t, expectedPermissions, permissions)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		mockRepo.On("FindAll").Return(nil, errors.New("db error"))

		permissions, err := usecase.FindAll()

		assert.Error(t, err)
		assert.Nil(t, permissions)
		mockRepo.AssertExpectations(t)
	})
}
