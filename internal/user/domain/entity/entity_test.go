package user_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("should create user when valid", func(t *testing.T) {
		input := &User{
			ID:          "123",
			Name:        "Alice",
			Surname:     "Smith",
			Nickname:    "alices",
			Age:         28,
			Email:       "alice@example.com",
			Password:    "secret",
			Roles:       []string{"user"},
			Permissions: []string{"read"},
		}



		user := NewUser(input)
		assert.NotNil(t, user)
		assert.Equal(t, "Alice", user.Name)
		assert.Equal(t, "alices", user.Nickname)
	})

	t.Run("should return nil when invalid", func(t *testing.T) {
		input := &User{
			ID:          "456",
			Name:        "",
			Surname:     "",
			Nickname:    "",
			Age:         -1,
			Email:       "invalid-email",
			Password:    "",
			Roles:       []string{},
			Permissions: []string{},
		}

		user := NewUser(input)
		assert.Nil(t, user)
	})
}

func TestUser_Getters(t *testing.T) {
	user := &User{
		ID:          "abc123",
		Name:        "Carlos",
		Surname:     "Silva",
		Nickname:    "casilva",
		Age:         35,
		Email:       "carlos@example.com",
		Password:    "secret123",
		Roles:       []string{"admin"},
		Permissions: []string{"read", "write"},
	}

	assert.Equal(t, "abc123", user.GetID())
	assert.Equal(t, "Carlos", user.GetName())
	assert.Equal(t, "Silva", user.GetSurname())
	assert.Equal(t, int32(35), user.GetAge())
	assert.Equal(t, "carlos@example.com", user.GetEmail())
	assert.Equal(t, "secret123", user.GetPassword())
	assert.Equal(t, []string{"admin"}, user.GetRoles())
	assert.Equal(t, []string{"read", "write"}, user.GetPermissions())
}