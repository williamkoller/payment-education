package user_mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
)

func TestToUser(t *testing.T) {
	user := &user_entity.User{
		ID:          "u1",
		Name:        "Alice",
		Surname:     "Silva",
		Nickname:    "ali",
		Email:       "alice@example.com",
		Age:         28,
		Roles:       []string{"admin"},
		Permissions: []string{"create", "update"},
	}

	resp := ToUser(user)

	assert.NotNil(t, resp)
	assert.Equal(t, user.ID, resp.ID)
	assert.Equal(t, user.Name, resp.Name)
	assert.Equal(t, user.Surname, resp.Surname)
	assert.Equal(t, user.Nickname, resp.Nickname)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, user.Age, resp.Age)
	assert.ElementsMatch(t, user.Roles, resp.Roles)
	assert.ElementsMatch(t, user.Permissions, resp.Permissions)
}
