package user_usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/dtos"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
	user_usecase "github.com/williamkoller/system-education/internal/user/usecase"
)

type MockUserRepository struct {
	mock.Mock
}

type MockBcryptAdapter struct {
	cost int
	mock.Mock
}

func (m *MockBcryptAdapter) Hash(plaintext string) (string, error) {
	args := m.Called(plaintext)
	return args.String(0), args.Error(1)
}

func (m *MockBcryptAdapter) HashComparer(plaintext string, hashed string) bool {
	expected := "mocked:" + plaintext
	return hashed == expected
}

func (m *MockUserRepository) Save(user *user_entity.User) (*user_entity.User, error) {
	args := m.Called(user)
	result, _ := args.Get(0).(*user_entity.User)
	return result, args.Error(1)
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

// ✅ Teste de criação com sucesso
func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto)

	input := dtos.AddUserDto{
		Name:     "Alice",
		Surname:  "Silva",
		Nickname: "ali",
		Age:      30,
		Email:    "alice@example.com",
		Password: "secure123",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, port_user_repository.ErrUserNotFound)
	mockCrypto.On("Hash", input.Password).Return("mocked:secure123", nil)
	mockRepo.On("Save", mock.MatchedBy(func(u *user_entity.User) bool {
		return u.Email == input.Email && u.Password == "mocked:secure123"
	})).Return(&user_entity.User{
		Name:     input.Name,
		Surname:  input.Surname,
		Nickname: input.Nickname,
		Age:      input.Age,
		Email:    input.Email,
		Password: "mocked:secure123",
	}, nil)

	user, err := usecase.Create(input)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, input.Email, user.Email)
	assert.Equal(t, "mocked:secure123", user.Password)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

// ✅ Teste: Usuário já existe
func TestCreate_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto)

	existing := &user_entity.User{Email: "alice@example.com"}
	mockRepo.On("FindByEmail", "alice@example.com").Return(existing, nil)

	_, err := usecase.Create(dtos.AddUserDto{Email: "alice@example.com"})

	assert.ErrorIs(t, err, port_user_repository.ErrUserAlreadyExists)
	mockRepo.AssertExpectations(t)
}

// ✅ Teste: erro ao salvar
func TestCreate_SaveFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto)

	input := dtos.AddUserDto{
		Name:     "Alice",
		Surname:  "Silva",
		Nickname: "ali",
		Age:      30,
		Email:    "alice@example.com",
		Password: "secure123",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, port_user_repository.ErrUserNotFound)
	mockCrypto.On("Hash", input.Password).Return("mocked:secure123", nil)
	mockRepo.On("Save", mock.MatchedBy(func(u *user_entity.User) bool {
		return u.Email == input.Email
	})).Return(nil, errors.New("save error"))

	_, err := usecase.Create(input)

	assert.Error(t, err)
	assert.Equal(t, "cannot create new user", err.Error())
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}
