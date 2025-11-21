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

func TestPullDomainEvents_NilUser(t *testing.T) {
	var user *User
	events := user.PullDomainEvents()
	assert.Nil(t, events)
}

func TestUpdateUser_PartialUpdate(t *testing.T) {
	t.Run("update only name", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Original",
			Surname:  "User",
			Nickname: "original",
			Age:      25,
			Email:    "original@example.com",
			Password: "pass123",
		}

		newName := "Updated"
		updated, err := user.UpdateUser(&newName, nil, nil, nil, nil)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "Updated", updated.Name)
		assert.Equal(t, "original", updated.Nickname)          // unchanged
		assert.Equal(t, "original@example.com", updated.Email) // unchanged
	})

	t.Run("update only nickname", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Test",
			Surname:  "User",
			Nickname: "old",
			Age:      25,
			Email:    "test@example.com",
			Password: "pass123",
		}

		newNickname := "new"
		updated, err := user.UpdateUser(nil, &newNickname, nil, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, "new", updated.Nickname)
		assert.Equal(t, "Test", updated.Name) // unchanged
	})

	t.Run("update only email", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Test",
			Surname:  "User",
			Nickname: "test",
			Age:      25,
			Email:    "old@example.com",
			Password: "pass123",
		}

		newEmail := "new@example.com"
		updated, err := user.UpdateUser(nil, nil, &newEmail, nil, nil)

		assert.NoError(t, err)
		assert.Equal(t, "new@example.com", updated.Email)
	})

	t.Run("update only password", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Test",
			Surname:  "User",
			Nickname: "test",
			Age:      25,
			Email:    "test@example.com",
			Password: "oldpass",
		}

		newPassword := "newpass"
		updated, err := user.UpdateUser(nil, nil, nil, &newPassword, nil)

		assert.NoError(t, err)
		assert.Equal(t, "newpass", updated.Password)
	})

	t.Run("update only age", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Test",
			Surname:  "User",
			Nickname: "test",
			Age:      25,
			Email:    "test@example.com",
			Password: "pass123",
		}

		newAge := int32(30)
		updated, err := user.UpdateUser(nil, nil, nil, nil, &newAge)

		assert.NoError(t, err)
		assert.Equal(t, int32(30), updated.Age)
	})

	t.Run("update with validation error - invalid email", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Test",
			Surname:  "User",
			Nickname: "test",
			Age:      25,
			Email:    "test@example.com",
			Password: "pass123",
		}

		invalidEmail := "invalid-email"
		updated, err := user.UpdateUser(nil, nil, &invalidEmail, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, updated)
		assert.Contains(t, err.Error(), "email is invalid")
	})

	t.Run("update with validation error - negative age", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Test",
			Surname:  "User",
			Nickname: "test",
			Age:      25,
			Email:    "test@example.com",
			Password: "pass123",
		}

		negativeAge := int32(-5)
		updated, err := user.UpdateUser(nil, nil, nil, nil, &negativeAge)

		assert.Error(t, err)
		assert.Nil(t, updated)
		assert.Contains(t, err.Error(), "age cannot be negative")
	})

	t.Run("update all fields at once", func(t *testing.T) {
		user := &User{
			ID:       "123",
			Name:     "Old",
			Surname:  "User",
			Nickname: "old",
			Age:      20,
			Email:    "old@example.com",
			Password: "oldpass",
		}

		newName := "New"
		newNickname := "new"
		newEmail := "new@example.com"
		newPassword := "newpass"
		newAge := int32(25)

		updated, err := user.UpdateUser(&newName, &newNickname, &newEmail, &newPassword, &newAge)

		assert.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, "New", updated.Name)
		assert.Equal(t, "new", updated.Nickname)
		assert.Equal(t, "new@example.com", updated.Email)
		assert.Equal(t, "newpass", updated.Password)
		assert.Equal(t, int32(25), updated.Age)
	})
}
