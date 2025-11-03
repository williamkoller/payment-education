package user_handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/dtos"
	user_handler "github.com/williamkoller/system-education/internal/user/handler"
	"github.com/williamkoller/system-education/shared/middleware"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Create(input dtos.AddUserDto) (*user_entity.User, error) {
	args := m.Called(input)
	user, _ := args.Get(0).(*user_entity.User)
	return user, args.Error(1)
}

func (m *MockUserUsecase) FindAll() ([]*user_entity.User, error) {
	args := m.Called()
	users, _ := args.Get(0).([]*user_entity.User)
	return users, args.Error(1)
}

func setupRouter(h *user_handler.UserHandler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.GlobalErrorHandler())
	r.POST("/users", h.CreateUser)
	return r
}

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockUserUsecase)
	h := user_handler.NewUserHandler(mockUsecase)
	router := setupRouter(h)

	body := []byte(`{
		"name": "Alice",
		"surname": "Silva",
		"nickname": "ali",
		"age": 30,
		"email": "alice@example.com",
		"password": "secure123",
		"roles": ["admin"],
		"permissions": ["create"]
	}`)

	expectedUser := &user_entity.User{
		ID:       "u1",
		Name:     "Alice",
		Surname:  "Silva",
		Nickname: "ali",
		Age:      30,
		Email:    "alice@example.com",
	}

	mockUsecase.On("Create", mock.Anything).Return(expectedUser, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"email":"alice@example.com"`)
	mockUsecase.AssertExpectations(t)
}

func TestCreateUser_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockUserUsecase)
	h := user_handler.NewUserHandler(mockUsecase)
	router := setupRouter(h)

	body := []byte(`{invalid json}`)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error"`)
}

func TestCreateUser_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := new(MockUserUsecase)
	h := user_handler.NewUserHandler(mockUsecase)
	router := setupRouter(h)

	body := []byte(`{
		"name": "Alice",
		"surname": "Silva",
		"nickname": "ali",
		"age": 30,
		"email": "alice@example.com",
		"password": "secure123",
		"roles": ["admin"],
		"permissions": ["create"]
	}`)

	mockUsecase.On("Create", mock.Anything).Return(nil, errors.New("email already exists"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"email already exists"`)
}
