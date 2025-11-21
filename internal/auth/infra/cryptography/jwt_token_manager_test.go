package infra_cryptography

import (
	"testing"
	"time"

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
