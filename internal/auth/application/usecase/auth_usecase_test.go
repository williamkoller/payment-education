package auth_usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestAuthUsecase_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockTokenManager, mockBcrypt)

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
	mockTokenManager.On("Sign", mock.MatchedBy(func(data map[string]interface{}) bool {
		return data["email"] == email && data["name"] == "John"
	})).Return(expectedToken, nil)

	token, err := usecase.Login(email, password)

	assert.NoError(t, err)
	assert.Equal(t, expectedToken, token)
	mockRepo.AssertExpectations(t)
	mockBcrypt.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockTokenManager, mockBcrypt)

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
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockTokenManager, mockBcrypt)

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
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockTokenManager, mockBcrypt)

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
	mockTokenManager.On("Sign", mock.Anything).Return("", errors.New("token generation failed"))

	token, err := usecase.Login(email, password)

	assert.Error(t, err)
	assert.Equal(t, "", token)
	assert.Equal(t, "error in generate token", err.Error())
	mockRepo.AssertExpectations(t)
	mockBcrypt.AssertExpectations(t)
	mockTokenManager.AssertExpectations(t)
}

func TestAuthUsecase_Profile_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockTokenManager, mockBcrypt)

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
	mockTokenManager := new(MockTokenManager)
	mockBcrypt := new(MockBcrypt)

	usecase := NewAuthUsecase(mockRepo, mockTokenManager, mockBcrypt)

	email := "nonexistent@example.com"

	mockRepo.On("FindByEmail", email).Return(nil, errors.New("not found"))

	result, err := usecase.Profile(email)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, "user not found", err.Error())
	mockRepo.AssertExpectations(t)
}
