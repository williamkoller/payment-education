package infra_cryptography

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	port_cryptography "github.com/williamkoller/system-education/internal/auth/port/cryptography"
)

type JWTTokenManager struct {
	secretKey string
	expiresIn time.Duration
}

var _ port_cryptography.TokenManager = &JWTTokenManager{}

func NewJWTTokenManager(secret string, expiresIn time.Duration) *JWTTokenManager {
	return &JWTTokenManager{secretKey: secret, expiresIn: expiresIn}
}

func (j *JWTTokenManager) Verify(tokenStr string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(j.secretKey), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		result := make(map[string]interface{})
		for k, v := range claims {
			result[k] = v
		}
		return result, nil
	}

	return nil, errors.New("could not parse claims")
}

func (j *JWTTokenManager) Sign(data map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{}

	for k, v := range data {
		claims[k] = v
	}

	claims["exp"] = time.Now().Add(j.expiresIn).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}
