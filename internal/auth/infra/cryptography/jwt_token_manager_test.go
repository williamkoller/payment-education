package infra_cryptography

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTTokenManager(t *testing.T) {
	manager := NewJWTTokenManager("secret", time.Hour)

	data := map[string]interface{}{"user_id": 123}
	token, err := manager.Sign(data)
	require.NoError(t, err)

	parsedData, err := manager.Verify(token)
	require.NoError(t, err)
	require.Equal(t, float64(123), parsedData["user_id"])
}

func TestJWTTokenManager_Sign_Success(t *testing.T) {
	manager := NewJWTTokenManager("secret-key", time.Hour)

	data := map[string]interface{}{
		"user_id": "user-123",
		"email":   "test@example.com",
		"modules": []string{"admin", "user"},
	}

	token, err := manager.Sign(data)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify the token can be parsed
	parsedData, err := manager.Verify(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", parsedData["user_id"])
	assert.Equal(t, "test@example.com", parsedData["email"])
}

func TestJWTTokenManager_Verify_InvalidToken(t *testing.T) {
	manager := NewJWTTokenManager("secret-key", time.Hour)

	invalidToken := "invalid.jwt.token"

	parsedData, err := manager.Verify(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, parsedData)
	assert.Equal(t, "invalid token", err.Error())
}

func TestJWTTokenManager_Verify_ExpiredToken(t *testing.T) {
	manager := NewJWTTokenManager("secret-key", 1*time.Millisecond)

	data := map[string]interface{}{"user_id": "user-123"}
	token, err := manager.Sign(data)
	require.NoError(t, err)

	// Wait for token to expire
	time.Sleep(10 * time.Millisecond)

	parsedData, err := manager.Verify(token)

	assert.Error(t, err)
	assert.Nil(t, parsedData)
	assert.Equal(t, "invalid token", err.Error())
}

func TestJWTTokenManager_Verify_WrongSecret(t *testing.T) {
	manager1 := NewJWTTokenManager("secret-key-1", time.Hour)
	manager2 := NewJWTTokenManager("secret-key-2", time.Hour)

	data := map[string]interface{}{"user_id": "user-123"}
	token, err := manager1.Sign(data)
	require.NoError(t, err)

	// Try to verify with different secret
	parsedData, err := manager2.Verify(token)

	assert.Error(t, err)
	assert.Nil(t, parsedData)
	assert.Equal(t, "invalid token", err.Error())
}

func TestJWTTokenManager_Verify_MalformedToken(t *testing.T) {
	manager := NewJWTTokenManager("secret-key", time.Hour)

	malformedToken := "not.a.valid.jwt.token.format"

	parsedData, err := manager.Verify(malformedToken)

	assert.Error(t, err)
	assert.Nil(t, parsedData)
	assert.Equal(t, "invalid token", err.Error())
}

func TestJWTTokenManager_Verify_EmptyToken(t *testing.T) {
	manager := NewJWTTokenManager("secret-key", time.Hour)

	parsedData, err := manager.Verify("")

	assert.Error(t, err)
	assert.Nil(t, parsedData)
	assert.Equal(t, "invalid token", err.Error())
}

func TestJWTTokenManager_Verify_UnexpectedSigningMethod(t *testing.T) {
	manager := NewJWTTokenManager("secret-key", time.Hour)

	// Create a token with RSA (not HMAC)
	claims := jwt.MapClaims{
		"user_id": "user-123",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	// Generate RSA key pair for testing
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	parsedData, err := manager.Verify(tokenString)

	// Should fail because we expect HMAC but got None
	assert.Error(t, err)
	assert.Nil(t, parsedData)
	assert.Equal(t, "invalid token", err.Error())
}

func TestJWTTokenManager_Sign_WithMultipleFields(t *testing.T) {
	manager := NewJWTTokenManager("secret-key", time.Hour)

	data := map[string]interface{}{
		"user_id":    "user-123",
		"email":      "test@example.com",
		"name":       "John Doe",
		"modules":    []string{"admin", "user"},
		"actions":    []string{"read", "write"},
		"created_at": time.Now().Unix(),
	}

	token, err := manager.Sign(data)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify all fields are preserved
	parsedData, err := manager.Verify(token)
	assert.NoError(t, err)
	assert.Equal(t, "user-123", parsedData["user_id"])
	assert.Equal(t, "test@example.com", parsedData["email"])
	assert.Equal(t, "John Doe", parsedData["name"])
	assert.NotNil(t, parsedData["modules"])
	assert.NotNil(t, parsedData["actions"])
	assert.NotNil(t, parsedData["exp"]) // exp should be added automatically
}

func TestNewJWTTokenManager(t *testing.T) {
	secret := "test-secret"
	expiresIn := 2 * time.Hour

	manager := NewJWTTokenManager(secret, expiresIn)

	assert.NotNil(t, manager)
	assert.Equal(t, secret, manager.secretKey)
	assert.Equal(t, expiresIn, manager.expiresIn)
}
