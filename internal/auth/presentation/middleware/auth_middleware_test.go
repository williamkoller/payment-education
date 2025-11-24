package auth_middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestAuthMiddleware_Success_WithAllClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)
	token := "valid.jwt.token"

	claims := map[string]interface{}{
		"email":   "user@example.com",
		"user_id": "user-123",
		"modules": []interface{}{"admin", "user"},
	}

	mockJWT.On("Verify", token).Return(claims, nil)

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		// Verify all claims were set in context
		email, _ := c.Get("userEmail")
		userID, _ := c.Get("userID")
		modules, _ := c.Get("modules")

		c.JSON(http.StatusOK, gin.H{
			"email":   email,
			"user_id": userID,
			"modules": modules,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user@example.com")
	assert.Contains(t, w.Body.String(), "user-123")
	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_Success_WithoutOptionalClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)
	token := "valid.jwt.token"

	// Only email claim (minimal valid token)
	claims := map[string]interface{}{
		"email": "user@example.com",
	}

	mockJWT.On("Verify", token).Return(claims, nil)

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		email, emailExists := c.Get("userEmail")
		_, userIDExists := c.Get("userID")
		_, modulesExist := c.Get("modules")

		c.JSON(http.StatusOK, gin.H{
			"email":         email,
			"email_exists":  emailExists,
			"userid_exists": userIDExists,
			"modules_exist": modulesExist,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user@example.com")
	assert.Contains(t, w.Body.String(), `"email_exists":true`)
	assert.Contains(t, w.Body.String(), `"userid_exists":false`)
	assert.Contains(t, w.Body.String(), `"modules_exist":false`)
	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// No Authorization header
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing or invalid token")
	mockJWT.AssertNotCalled(t, "Verify")
}

func TestAuthMiddleware_InvalidAuthorizationHeader_TooShort(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bear")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing or invalid token")
	mockJWT.AssertNotCalled(t, "Verify")
}

func TestAuthMiddleware_InvalidAuthorizationHeader_WrongPrefix(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Basic token123")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "missing or invalid token")
	mockJWT.AssertNotCalled(t, "Verify")
}

func TestAuthMiddleware_TokenVerificationFails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)
	token := "invalid.jwt.token"

	mockJWT.On("Verify", token).Return(nil, errors.New("invalid signature"))

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)
	token := "expired.jwt.token"

	mockJWT.On("Verify", token).Return(nil, errors.New("token expired"))

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_EmptyToken(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)

	// Mock should expect empty string and return error
	mockJWT.On("Verify", "").Return(nil, errors.New("empty token"))

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer ")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_ContextPropagation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)
	token := "valid.jwt.token"

	claims := map[string]interface{}{
		"email":   "test@example.com",
		"user_id": "user-456",
		"modules": []interface{}{"reports", "analytics"},
	}

	mockJWT.On("Verify", token).Return(claims, nil)

	var capturedEmail interface{}
	var capturedUserID interface{}
	var capturedModules interface{}

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		capturedEmail, _ = c.Get("userEmail")
		capturedUserID, _ = c.Get("userID")
		capturedModules, _ = c.Get("modules")
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test@example.com", capturedEmail)
	assert.Equal(t, "user-456", capturedUserID)
	assert.NotNil(t, capturedModules)
	mockJWT.AssertExpectations(t)
}

func TestAuthMiddleware_WithActionsInClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJWT := new(MockTokenManager)
	token := "valid.jwt.token"

	claims := map[string]interface{}{
		"email":   "user@example.com",
		"user_id": "user-123",
		"modules": []interface{}{"admin", "user"},
		"actions": []interface{}{"read", "delete", "update"},
	}

	mockJWT.On("Verify", token).Return(claims, nil)

	router := gin.New()
	router.GET("/test", AuthMiddleware(mockJWT), func(c *gin.Context) {
		// Verify all claims were set in context
		email, _ := c.Get("userEmail")
		userID, _ := c.Get("userID")
		modules, _ := c.Get("modules")
		actions, actionsExist := c.Get("actions")

		c.JSON(http.StatusOK, gin.H{
			"email":         email,
			"user_id":       userID,
			"modules":       modules,
			"actions":       actions,
			"actions_exist": actionsExist,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user@example.com")
	assert.Contains(t, w.Body.String(), "user-123")
	assert.Contains(t, w.Body.String(), `"actions_exist":true`)
	assert.Contains(t, w.Body.String(), "read")
	assert.Contains(t, w.Body.String(), "delete")
	assert.Contains(t, w.Body.String(), "update")
	mockJWT.AssertExpectations(t)
}
