package user_usecase_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	user_usecase "github.com/williamkoller/system-education/internal/user/application/usecase"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
	"github.com/williamkoller/system-education/internal/user/presentation/dtos"
	shared_event "github.com/williamkoller/system-education/shared/domain/event"
)

type MockUserRepository struct {
	mock.Mock
}

type MockBcryptAdapter struct {
	mock.Mock
}

type MockEvent struct {
	mock.Mock
}

func (m *MockEvent) Register(eventName string, handler shared_event.Handler) {
	m.Called(eventName, handler)
}

func (m *MockEvent) Dispatch(event interface{}) {
	m.Called(event)
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

func (m *MockUserRepository) Update(id string, user *user_entity.User) (*user_entity.User, error) {
	args := m.Called(id, user)
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

func TestCreate_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)
	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

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
	mockEvent.On("Dispatch", mock.MatchedBy(func(e interface{}) bool {
		return e != nil
	})).Return()

	user, err := usecase.Create(input)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, input.Email, user.Email)
	assert.Equal(t, "mocked:secure123", user.Password)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestCreate_AlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	existing := &user_entity.User{Email: "alice@example.com"}
	mockRepo.On("FindByEmail", "alice@example.com").Return(existing, nil)

	_, err := usecase.Create(dtos.AddUserDto{Email: "alice@example.com"})

	assert.ErrorIs(t, err, port_user_repository.ErrUserAlreadyExists)
	mockRepo.AssertExpectations(t)
}

func TestCreate_SaveFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

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

	mockEvent.On("Dispatch", mock.MatchedBy(func(e interface{}) bool {
		return e != nil
	})).Return()

	_, err := usecase.Create(input)

	assert.Error(t, err)
	assert.Equal(t, "cannot create new user", err.Error())
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestCreate_HashFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	input := dtos.AddUserDto{
		Name:     "Alice",
		Surname:  "Silva",
		Nickname: "ali",
		Age:      30,
		Email:    "alice@example.com",
		Password: "secure123",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, port_user_repository.ErrUserNotFound)
	mockCrypto.On("Hash", input.Password).Return("", errors.New("hashing error"))

	user, err := usecase.Create(input)

	assert.Nil(t, user)
	assert.Error(t, err)
	assert.Equal(t, "hashing error", err.Error())
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestCreate_InvalidUserData(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	input := dtos.AddUserDto{
		Name:     "Alice",
		Surname:  "Silva",
		Nickname: "ali",
		Age:      30,
		Email:    "",
		Password: "secure123",
	}

	mockRepo.On("FindByEmail", input.Email).Return(nil, port_user_repository.ErrUserNotFound)
	mockCrypto.On("Hash", input.Password).Return("mocked:secure123", nil)

	user, err := usecase.Create(input)

	assert.Nil(t, user)
	assert.EqualError(t, err, "invalid user data")
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func TestFindAll_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	expectedUsers := []*user_entity.User{
		{Name: "Alice"}, {Name: "Bob"},
	}

	mockRepo.On("FindAll").Return(expectedUsers, nil)

	users, err := usecase.FindAll()

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "Alice", users[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestFindByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	expectedUser := &user_entity.User{ID: "123", Name: "Alice"}
	mockRepo.On("FindByID", "123").Return(expectedUser, nil)

	user, err := usecase.FindByID("123")

	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
	mockRepo.AssertExpectations(t)
}

func TestDelete_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	user := &user_entity.User{ID: "123"}
	mockRepo.On("FindByID", "123").Return(user, nil)
	mockRepo.On("Delete", "123").Return(nil)

	err := usecase.Delete("123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDelete_FailDelete(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	user := &user_entity.User{ID: "123"}
	mockRepo.On("FindByID", "123").Return(user, nil)
	mockRepo.On("Delete", "123").Return(errors.New("db error"))

	err := usecase.Delete("123")

	assert.EqualError(t, err, "error in delete user")
	mockRepo.AssertExpectations(t)
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockCrypto := new(MockBcryptAdapter)
	mockEvent := new(MockEvent)

	usecase := user_usecase.NewUserUsecase(mockRepo, mockCrypto, mockEvent)

	id := "123"
	existingUser := &user_entity.User{ID: id, Email: "old@example.com"}
	mockRepo.On("FindByID", id).Return(existingUser, nil)

	input := dtos.UpdateUserDto{
		Name:     strPtr("Updated"),
		Email:    strPtr("new@example.com"),
		Password: strPtr("newpass"),
	}

	mockCrypto.On("Hash", "newpass").Return("hashed:newpass", nil)

	mockRepo.On("Update", id, mock.MatchedBy(func(u *user_entity.User) bool {
		return u.Email == "new@example.com" && u.Password == "hashed:newpass"
	})).Return(existingUser, nil)

	user, err := usecase.Update(id, input)

	assert.NoError(t, err)
	assert.Equal(t, "new@example.com", user.Email)
	mockRepo.AssertExpectations(t)
	mockCrypto.AssertExpectations(t)
}

func strPtr(s string) *string {
	return &s
}
