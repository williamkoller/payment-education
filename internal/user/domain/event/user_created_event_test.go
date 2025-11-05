package user_event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUserCreatedEvent(t *testing.T) {
	e := NewUserCreatedEvent("123", "William", "william@example.com")

	assert.NotNil(t, e)
	assert.Equal(t, "123", e.UserID)
	assert.Equal(t, "William", e.Name)
	assert.Equal(t, "william@example.com", e.Email)

	assert.False(t, e.Date.IsZero(), "A data não deve estar zero — deve ser inicializada em NewUserCreatedEvent")
}

func TestUserCreatedEvent_EventName(t *testing.T) {
	e := NewUserCreatedEvent("1", "Will", "will@mail.com")

	assert.Equal(t, "user.created", e.EventName())
}

func TestUserCreatedEvent_OccurredOn(t *testing.T) {
	now := time.Now()
	e := &UserCreatedEvent{
		UserID: "10",
		Name:   "Alice",
		Email:  "alice@example.com",
		Date:   now,
	}

	assert.Equal(t, now, e.OccurredOn())
	assert.WithinDuration(t, now, e.OccurredOn(), time.Millisecond)
}
