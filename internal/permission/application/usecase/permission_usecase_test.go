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

func (m *MockPermissionRepository) FindByID(id string) (*permission_entity.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permission_entity.Permission), args.Error(1)
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

func TestPermissionUsecase_FindById(t *testing.T) {
	t.Run("should return permission by id", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		expectedPermission := &permission_entity.Permission{ID: "123", UserID: "user-1"}

		mockRepo.On("FindByID", "123").Return(expectedPermission, nil)

		permission, err := usecase.FindById("123")

		assert.NoError(t, err)
		assert.Equal(t, expectedPermission, permission)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		mockRepo.On("FindByID", "123").Return(nil, errors.New("db error"))

		permission, err := usecase.FindById("123")

		assert.Error(t, err)
		assert.Nil(t, permission)
		mockRepo.AssertExpectations(t)
	})
}

func TestPermissionUsecase_Update(t *testing.T) {
	t.Run("should update permission successfully", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		id := "123"
		modules := []string{"module2"}
		actions := []string{"write"}
		level := "user"
		description := "updated permission"

		input := permission_dtos.UpdatePermissionDto{
			Modules:     &modules,
			Actions:     &actions,
			Level:       &level,
			Description: &description,
		}

		existingPermission := &permission_entity.Permission{
			ID:          id,
			UserID:      "user-1",
			Modules:     []string{"module1"},
			Actions:     []string{"read"},
			Level:       "admin",
			Description: "test permission",
		}

		// Expected updated permission
		updatedPermission := &permission_entity.Permission{
			ID:          id,
			UserID:      "user-1",
			Modules:     *input.Modules,
			Actions:     *input.Actions,
			Level:       *input.Level,
			Description: *input.Description,
		}

		mockRepo.On("FindByID", id).Return(existingPermission, nil)
		// We can't easily match the exact object pointer because it's modified in place,
		// but we can match the content or just use mock.AnythingOfType
		mockRepo.On("Update", id, mock.AnythingOfType("*permission_entity.Permission")).Return(updatedPermission, nil)

		permission, err := usecase.Update(id, input)

		assert.NoError(t, err)
		assert.Equal(t, updatedPermission, permission)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when find by id fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		id := "123"
		input := permission_dtos.UpdatePermissionDto{}

		mockRepo.On("FindByID", id).Return(nil, permission_entity.ErrNotFound)

		permission, err := usecase.Update(id, input)

		assert.Error(t, err)
		assert.ErrorIs(t, err, permission_entity.ErrNotFound)
		assert.Nil(t, permission)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		id := "123"
		modules := []string{"module2"}
		input := permission_dtos.UpdatePermissionDto{
			Modules: &modules,
		}

		existingPermission := &permission_entity.Permission{
			ID:      id,
			Modules: []string{"module1"},
			Level:   "user",
		}

		mockRepo.On("FindByID", id).Return(existingPermission, nil)
		mockRepo.On("Update", id, mock.AnythingOfType("*permission_entity.Permission")).Return(nil, errors.New("db error"))

		permission, err := usecase.Update(id, input)

		assert.Error(t, err)
		assert.Nil(t, permission)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when validation fails during update", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		id := "123"
		level := ""
		input := permission_dtos.UpdatePermissionDto{
			Level: &level,
		}

		existingPermission := &permission_entity.Permission{
			ID:    id,
			Level: "admin",
		}

		mockRepo.On("FindByID", id).Return(existingPermission, nil)

		permission, err := usecase.Update(id, input)

		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.Contains(t, err.Error(), "level cannot be empty")
		mockRepo.AssertNotCalled(t, "Update")
	})
}

func TestPermissionUsecase_Delete(t *testing.T) {
	t.Run("should delete permission successfully", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		id := "123"

		// Mock FindByID first as Delete now calls it
		mockRepo.On("FindByID", id).Return(&permission_entity.Permission{ID: id}, nil)
		mockRepo.On("Delete", id).Return(nil)

		err := usecase.Delete(id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when find by id fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		id := "123"

		mockRepo.On("FindByID", id).Return(nil, permission_entity.ErrNotFound)

		err := usecase.Delete(id)

		assert.Error(t, err)
		assert.ErrorIs(t, err, permission_entity.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		id := "123"

		mockRepo.On("FindByID", id).Return(&permission_entity.Permission{ID: id}, nil)
		mockRepo.On("Delete", id).Return(errors.New("db error"))

		err := usecase.Delete(id)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestPermissionUsecase_FindPermissionByUserID(t *testing.T) {
	t.Run("should return permissions by user id", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		userID := "user-1"
		expectedPermissions := []*permission_entity.Permission{
			{ID: "1", UserID: userID},
		}

		mockRepo.On("FindPermissionByUserID", userID).Return(expectedPermissions, nil)

		permissions, err := usecase.FindPermissionByUserID(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPermissions, permissions)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo := new(MockPermissionRepository)
		usecase := NewPermissionUsecase(mockRepo)

		userID := "user-1"

		mockRepo.On("FindPermissionByUserID", userID).Return(nil, errors.New("db error"))

		permissions, err := usecase.FindPermissionByUserID(userID)

		assert.Error(t, err)
		assert.Nil(t, permissions)
		mockRepo.AssertExpectations(t)
	})
}
