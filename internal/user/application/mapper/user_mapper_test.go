package user_mapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	userEntity "github.com/williamkoller/system-education/internal/user/domain/entity"
)

func TestToUser(t *testing.T) {
	user := &userEntity.User{
		ID:       "u1",
		Name:     "Alice",
		Surname:  "Silva",
		Nickname: "ali",
		Email:    "alice@example.com",
		Age:      28,
	}

	tu := ToUser(user)

	ur := &UserResponse{
		ID:        tu.ID,
		Name:      tu.Name,
		Surname:   tu.Surname,
		Nickname:  tu.Nickname,
		Email:     tu.Email,
		Age:       tu.Age,
		CreatedAt: tu.CreatedAt,
		UpdatedAt: tu.UpdatedAt,
	}

	assert.NotNil(t, ur)
	assert.Equal(t, user.ID, ur.ID)
	assert.Equal(t, user.Name, ur.Name)
	assert.Equal(t, user.Surname, ur.Surname)
	assert.Equal(t, user.Nickname, ur.Nickname)
	assert.Equal(t, user.Email, ur.Email)
	assert.Equal(t, user.Age, ur.Age)
	assert.Equal(t, user.CreatedAt, ur.CreatedAt)
	assert.Equal(t, user.UpdatedAt, ur.UpdatedAt)
}
