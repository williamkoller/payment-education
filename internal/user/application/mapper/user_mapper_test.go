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

func TestToUsers(t *testing.T) {
	users := []*userEntity.User{
		{ID: "1", Name: "Alice", Email: "alice@example.com"},
		{ID: "2", Name: "Bob", Email: "bob@example.com"},
	}

	responses := ToUsers(users)

	assert.Len(t, responses, 2)

	assert.Equal(t, "1", responses[0].ID)
	assert.Equal(t, "Alice", responses[0].Name)
	assert.Equal(t, "alice@example.com", responses[0].Email)

	assert.Equal(t, "2", responses[1].ID)
	assert.Equal(t, "Bob", responses[1].Name)
	assert.Equal(t, "bob@example.com", responses[1].Email)
}
