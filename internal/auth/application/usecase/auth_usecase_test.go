package auth_usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	permissionEntity "github.com/williamkoller/system-education/internal/permission/domain/entity"
	userEntity "github.com/williamkoller/system-education/internal/user/domain/entity"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user *userEntity.User) (*userEntity.User, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*userEntity.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*userEntity.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) Update(id string, user *userEntity.User) (*userEntity.User, error) {
	args := m.Called(id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userEntity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll() ([]*userEntity.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*userEntity.User), args.Error(1)
}

type MockTokenManager struct {
	mock.Mock
}

func (m *MockTokenManager) Sign(data map[string]interface{}) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockTokenManager) Verify(token string) (map[string]interface{}, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

type MockBcrypt struct {
	mock.Mock
}

func (m *MockBcrypt) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockBcrypt) HashComparer(password, hash string) (bool, error) {
	args := m.Called(password, hash)
	return args.Bool(0), args.Error(1)
}

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) Save(permission *permissionEntity.Permission) (*permissionEntity.Permission, error) {
	args := m.Called(permission)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionEntity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindByID(id string) (*permissionEntity.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionEntity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) FindPermissionByUserID(userID string) ([]*permissionEntity.Permission, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permissionEntity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Update(id string, permission *permissionEntity.Permission) (*permissionEntity.Permission, error) {
	args := m.Called(id, permission)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*permissionEntity.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPermissionRepository) FindAll() ([]*permissionEntity.Permission, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*permissionEntity.Permission), args.Error(1)
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "test@example.com"
	password := "password123"
	hashedPassword := "$2a$10$hashedpassword"

	user := &userEntity.User{
		ID:        "user-123",
		Name:      "John",
		Surname:   "Doe",
		Nickname:  "johndoe",
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedToken := "jwt.token.here"

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockBcrypt.On("HashComparer", password, hashedPassword).Return(true, nil)
	mockPermissionRepo.On("FindPermissionByUserID", "user-123").Return([]*permissionEntity.Permission{}, nil)
	mockTokenManager.On("Sign", mock.MatchedBy(func(data map[string]interface{}) bool {
		return data["email"] == email && data["name"] == "John" && data["user_id"] == "user-123"
	})).Return(expectedToken, nil)

	token, err := usecase.Login(email, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockBcrypt.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "nonexistent@example.com"
	password := "password123"

	mockRepo.On("FindByEmail", email).Return(nil, errors.New("not found"))

	token, err := usecase.Login(email, password)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "test@example.com"
	password := "wrongpassword"
	hashedPassword := "$2a$10$hashedpassword"

	user := &userEntity.User{
		ID:        "user-123",
		Name:      "John",
		Surname:   "Doe",
		Nickname:  "johndoe",
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockBcrypt.On("HashComparer", password, hashedPassword).Return(false, errors.New("password mismatch"))

	// Act
	token, err := usecase.Login(email, password)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, "invalid credentials", err.Error())
	mockRepo.AssertExpectations(t)
	mockBcrypt.AssertExpectations(t)
}

func TestAuthUsecase_Login_TokenGenerationError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "test@example.com"
	password := "password123"
	hashedPassword := "$2a$10$hashedpassword"

	user := &userEntity.User{
		ID:        "user-123",
		Name:      "John",
		Surname:   "Doe",
		Nickname:  "johndoe",
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockBcrypt.On("HashComparer", password, hashedPassword).Return(true, nil)
	mockPermissionRepo.On("FindPermissionByUserID", "user-123").Return([]*permissionEntity.Permission{}, nil)
	mockTokenManager.On("Sign", mock.Anything).Return("", errors.New("token generation failed"))

	token, err := usecase.Login(email, password)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, "error in generate token", err.Error())
	mockRepo.AssertExpectations(t)
	mockBcrypt.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}

func TestAuthUsecase_Profile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "test@example.com"

	user := &userEntity.User{
		ID:        "user-123",
		Name:      "John",
		Surname:   "Doe",
		Nickname:  "johndoe",
		Email:     email,
		Password:  "$2a$10$hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByEmail", email).Return(user, nil)

	result, err := usecase.Profile(email)

	assert.NoError(t, err)
	assert.Equal(t, email, result)
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Profile_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "nonexistent@example.com"

	mockRepo.On("FindByEmail", email).Return(nil, errors.New("not found"))

	result, err := usecase.Profile(email)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestAuthUsecase_Login_WithPermissionsAndModules(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "test@example.com"
	password := "password123"
	hashedPassword := "$2a$10$hashedpassword"

	user := &userEntity.User{
		ID:        "user-123",
		Name:      "John",
		Surname:   "Doe",
		Nickname:  "johndoe",
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create permissions with modules
	permission1 := &permissionEntity.Permission{}
	permission1.UserID = "user-123"
	permission1.Modules = []string{"admin", "user"}

	permission2 := &permissionEntity.Permission{}
	permission2.UserID = "user-123"
	permission2.Modules = []string{"reports"}

	permissions := []*permissionEntity.Permission{permission1, permission2}

	expectedToken := "jwt.token.with.modules"

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockBcrypt.On("HashComparer", password, hashedPassword).Return(true, nil)
	mockPermissionRepo.On("FindPermissionByUserID", "user-123").Return(permissions, nil)
	mockTokenManager.On("Sign", mock.MatchedBy(func(data map[string]interface{}) bool {
		modules, ok := data["modules"].([]string)
		if !ok {
			return false
		}
		// Should have all modules from both permissions
		return len(modules) == 3 &&
			data["email"] == email &&
			data["user_id"] == "user-123"
	})).Return(expectedToken, nil)

	token, err := usecase.Login(email, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockBcrypt.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}

func TestAuthUsecase_Login_PermissionFetchError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPermissionRepo := new(MockPermissionRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockPermissionRepo, mockTokenManager, mockBcrypt)

	email := "test@example.com"
	password := "password123"
	hashedPassword := "$2a$10$hashedpassword"

	user := &userEntity.User{
		ID:        "user-123",
		Name:      "John",
		Surname:   "Doe",
		Nickname:  "johndoe",
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	expectedToken := "jwt.token.no.modules"

	mockRepo.On("FindByEmail", email).Return(user, nil)
	mockBcrypt.On("HashComparer", password, hashedPassword).Return(true, nil)
	// Permission fetch fails, but login should still succeed with empty modules
	mockPermissionRepo.On("FindPermissionByUserID", "user-123").Return(nil, errors.New("permission db error"))
	mockTokenManager.On("Sign", mock.MatchedBy(func(data map[string]interface{}) bool {
		modules, ok := data["modules"].([]string)
		if !ok {
			return false
		}
		// Should have empty modules array when permission fetch fails
		return len(modules) == 0 &&
			data["email"] == email &&
			data["user_id"] == "user-123"
	})).Return(expectedToken, nil)

	token, err := usecase.Login(email, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockBcrypt.AssertExpectations(t)
	mockPermissionRepo.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}
