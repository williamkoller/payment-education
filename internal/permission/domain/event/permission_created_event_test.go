package permission_event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPermissionCreatedEvent(t *testing.T) {
	permissionID := "123"
	userID := "user-123"
	modules := []string{"module1"}
	actions := []string{"read"}
	level := "admin"
	description := "test permission"

	event := NewPermissionCreatedEvent(permissionID, userID, modules, actions, level, description)

	assert.NotNil(t, event)
	assert.Equal(t, permissionID, event.PermissionID)
	assert.Equal(t, userID, event.UserID)
	assert.Equal(t, modules, event.Modules)
	assert.Equal(t, actions, event.Actions)
	assert.Equal(t, level, event.Level)
	assert.Equal(t, description, event.Description)
	assert.WithinDuration(t, time.Now(), event.Date, time.Second)
}

func TestPermissionCreatedEvent_EventName(t *testing.T) {
	event := &PermissionCreatedEvent{}
	assert.Equal(t, "permission.created", event.EventName())
}

func TestPermissionCreatedEvent_OccurredOn(t *testing.T) {
	now := time.Now()
	event := &PermissionCreatedEvent{Date: now}
	assert.Equal(t, now, event.OccurredOn())
}
