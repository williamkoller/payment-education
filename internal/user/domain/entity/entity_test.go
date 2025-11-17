package user_entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sharedEvent "github.com/williamkoller/system-education/shared/domain/event"
)

func TestNewUser(t *testing.T) {
	t.Run("should create user when valid", func(t *testing.T) {
		input := &User{
			ID:       "123",
			Name:     "Alice",
			Surname:  "Smith",
			Nickname: "alices",
			Age:      28,
			Email:    "alice@example.com",
			Password: "secret",
		}

		user := NewUser(input)
		assert.NotNil(t, user)
		assert.Equal(t, "Alice", user.Name)
		assert.Equal(t, "alices", user.Nickname)
	})

	t.Run("should return nil when invalid", func(t *testing.T) {
		input := &User{
			ID:       "456",
			Name:     "",
			Surname:  "",
			Nickname: "",
			Age:      -1,
			Email:    "invalid-email",
			Password: "",
		}

		user := NewUser(input)
		assert.Nil(t, user)
	})
}

func TestUser_Getters(t *testing.T) {
	user := &User{
		ID:       "abc123",
		Name:     "Carlos",
		Surname:  "Silva",
		Nickname: "casilva",
		Age:      35,
		Email:    "carlos@example.com",
		Password: "secret123",
	}

	assert.Equal(t, "abc123", user.GetID())
	assert.Equal(t, "Carlos", user.GetName())
	assert.Equal(t, "Silva", user.GetSurname())
	assert.Equal(t, int32(35), user.GetAge())
	assert.Equal(t, "carlos@example.com", user.GetEmail())
	assert.Equal(t, "secret123", user.GetPassword())
}

func TestUpdateUser(t *testing.T) {
	t.Run("should update user successfully when valid values provided", func(t *testing.T) {
		user := &User{
			ID:       "abc123",
			Name:     "Carlos",
			Surname:  "Silva",
			Nickname: "casilva",
			Age:      35,
			Email:    "carlos@example.com",
			Password: "secret123",
		}

		newName := "Joao"
		newNickname := "C"
		newEmail := "joao.dev@mail.com"
		newPassword := "secret123456"
		newAge := int32(36)

		updated, err := user.UpdateUser(&newName, &newNickname, &newEmail, &newPassword, &newAge)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, newName, updated.Name)
		assert.Equal(t, newNickname, updated.Nickname)
		assert.Equal(t, newEmail, updated.Email)
		assert.Equal(t, newPassword, updated.Password)
		assert.Equal(t, newAge, updated.Age)
	})

}

type MockEvent struct{}

func (m MockEvent) EventName() string {
	return "MockEvent"
}

func (m MockEvent) OccurredOn() time.Time {
	return time.Now()
}

func TestPullDomainEvents(t *testing.T) {
	user := &User{
		ID:           "1",
		Name:         "Alice",
		Email:        "alice@example.com",
		domainEvents: []sharedEvent.Event{MockEvent{}},
	}

	events := user.PullDomainEvents()

	assert.Len(t, events, 1)
	assert.Equal(t, "MockEvent", events[0].EventName())
	assert.Empty(t, user.PullDomainEvents())
}
