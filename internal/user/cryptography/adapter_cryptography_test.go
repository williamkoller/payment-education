package cryptography_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/williamkoller/system-education/internal/user/cryptography"
)

func TestBcryptAdapter_HashAndCompare(t *testing.T) {
	hasher := cryptography.NewBcryptHasher(0) // usa default cost

	password := "supersecret"

	hash, err := hasher.Hash(password)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	isMatch := hasher.HashComparer(password, hash)
	assert.True(t, isMatch)

	isWrong := hasher.HashComparer("wrongpassword", hash)
	assert.False(t, isWrong)
}
