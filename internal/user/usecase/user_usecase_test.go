package user_usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/dtos"
	"github.com/williamkoller/system-education/internal/user/port/repository"
	user_usecase "github.com/williamkoller/system-education/internal/user/usecase"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(user *user_entity.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id string) (*user_entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*user_entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*user_entity.User, error) {
	args := m.Called(email)
	user, ok := args.Get(0).(*user_entity.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockUserRepository) FindAll() ([]*user_entity.User, error) {
	args := m.Called()
	return args.Get(0).([]*user_entity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := user_usecase.NewUserUsecase(mockRepo)

	input := dtos.AddUserDto{
		Name:        "Alice",
		Surname:     "Silva",
		Nickname:    "ali",
		Age:         30,
		Email:       "alice@example.com",
		Password:    "secure123",
		Roles:       []string{"admin"},
		Permissions: []string{"create"},
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, port_user_repository.ErrUserNotFound)
	mockRepo.On("Save", mock.MatchedBy(func(u *user_entity.User) bool {
		return u.Email == input.Email
	})).Return(nil)

	user, err := usecase.Create(input)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, input.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestCreate_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := user_usecase.NewUserUsecase(mockRepo)

	existing := &user_entity.User{Email: "alice@example.com"}
	mockRepo.On("FindByEmail", "alice@example.com").Return(existing, nil)

	_, err := usecase.Create(dtos.AddUserDto{Email: "alice@example.com"})

	assert.Error(t, err)
	assert.Equal(t, "user with this email already exists", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestCreate_FindByEmailFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := user_usecase.NewUserUsecase(mockRepo)

	mockRepo.On("FindByEmail", "error@example.com").Return(nil, errors.New("db error"))

	_, err := usecase.Create(dtos.AddUserDto{Email: "error@example.com"})

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
}

func TestCreate_InvalidUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := user_usecase.NewUserUsecase(mockRepo)

	// forçar falha em NewUser (simulado por omitir campos obrigatórios)
	input := dtos.AddUserDto{
		Email: "invalid@example.com",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, port_user_repository.ErrUserNotFound)

	_, err := usecase.Create(input)

	assert.Error(t, err)
	assert.Equal(t, "invalid user data", err.Error())
}

func TestCreate_SaveFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	usecase := user_usecase.NewUserUsecase(mockRepo)

	input := dtos.AddUserDto{
		Name:        "Alice",
		Surname:     "Silva",
		Nickname:    "ali",
		Age:         30,
		Email:       "alice@example.com",
		Password:    "secure123",
		Roles:       []string{"admin"},
		Permissions: []string{"create"},
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, port_user_repository.ErrUserNotFound)
	mockRepo.On("Save", mock.MatchedBy(func(u *user_entity.User) bool {
		return u.Email == input.Email
	})).Return(errors.New("save error"))

	_, err := usecase.Create(input)

	assert.Error(t, err)
	assert.Equal(t, "cannot create new user", err.Error())
}
