package user_cryptography_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	user_cryptography "github.com/williamkoller/system-education/internal/user/infra/cryptography"
)

func TestBcryptAdapter_HashAndCompare(t *testing.T) {
	hasher := user_cryptography.NewBcryptHasher(0)

	password := "supersecret"

	hash, err := hasher.Hash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	isMatch, err := hasher.HashComparer(password, hash)
	assert.NoError(t, err)
	assert.True(t, isMatch)

	isWrong, err := hasher.HashComparer("wrongpassword", hash)
	assert.Error(t, err)
	assert.False(t, isWrong)
}
