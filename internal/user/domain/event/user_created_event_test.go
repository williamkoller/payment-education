package event_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/williamkoller/system-education/internal/user/domain/event"
)

func TestNewUserCreatedEvent(t *testing.T) {
	e := event.NewUserCreatedEvent("123", "William", "william@example.com")

	assert.NotNil(t, e, "Evento não deve ser nulo")
	assert.Equal(t, "123", e.UserID)
	assert.Equal(t, "William", e.Name)
	assert.Equal(t, "william@example.com", e.Email)

	assert.True(t, e.Date.IsZero(), "A data deve estar zero se não for setada manualmente")
}

func TestUserCreatedEvent_EventName(t *testing.T) {
	e := event.NewUserCreatedEvent("1", "Will", "will@mail.com")

	assert.Equal(t, "user.created", e.EventName())
}

func TestUserCreatedEvent_OccurredOn(t *testing.T) {
	now := time.Now()
	e := &event.UserCreatedEvent{
		UserID: "10",
		Name:   "Alice",
		Email:  "alice@example.com",
		Date:   now,
	}

	assert.Equal(t, now, e.OccurredOn())
	assert.WithinDuration(t, now, e.OccurredOn(), time.Millisecond)
}
