package permission_entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	permission_event "github.com/williamkoller/system-education/internal/permission/domain/event"
)

func TestNewPermission(t *testing.T) {
	t.Run("should create a permission successfully", func(t *testing.T) {
		p := &Permission{
			ID:          "123",
			UserID:      "user-123",
			Modules:     []string{"module1"},
			Actions:     []string{"read"},
			Level:       "admin",
			Description: "test permission",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		permission, err := NewPermission(p)

		assert.NoError(t, err)
		assert.NotNil(t, permission)
		assert.Equal(t, p.ID, permission.ID)
		assert.Equal(t, p.UserID, permission.UserID)
		assert.Equal(t, p.Modules, permission.Modules)
		assert.Equal(t, p.Actions, permission.Actions)
		assert.Equal(t, p.Level, permission.Level)
		assert.Equal(t, p.Description, permission.Description)

		events := permission.PullDomainEvents()
		assert.Len(t, events, 1)
		assert.Equal(t, "permission.created", events[0].EventName())
	})

	t.Run("should return error when id is empty", func(t *testing.T) {
		p := &Permission{
			UserID:  "user-123",
			Modules: []string{"module1"},
			Actions: []string{"read"},
			Level:   "admin",
		}

		permission, err := NewPermission(p)

		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.Contains(t, err.Error(), "id is required")
	})

	t.Run("should return error when user id is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			Modules: []string{"module1"},
			Actions: []string{"read"},
			Level:   "admin",
		}

		permission, err := NewPermission(p)

		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.Contains(t, err.Error(), "user id is required")
	})

	t.Run("should return error when modules is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			UserID:  "user-123",
			Actions: []string{"read"},
			Level:   "admin",
		}

		permission, err := NewPermission(p)

		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.Contains(t, err.Error(), "modules is required")
	})

	t.Run("should return error when actions is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			UserID:  "user-123",
			Modules: []string{"module1"},
			Level:   "admin",
		}

		permission, err := NewPermission(p)

		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.Contains(t, err.Error(), "actions is required")
	})

	t.Run("should return error when level is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			UserID:  "user-123",
			Modules: []string{"module1"},
			Actions: []string{"read"},
		}

		permission, err := NewPermission(p)

		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.Contains(t, err.Error(), "level is required")
	})
	t.Run("should return error when permission is nil", func(t *testing.T) {
		permission, err := NewPermission(nil)

		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.Contains(t, err.Error(), "permission is required")
	})
}

func TestPermission_Getters(t *testing.T) {
	t.Run("should return values when permission is valid", func(t *testing.T) {
		now := time.Now()
		p := &Permission{
			ID:          "123",
			UserID:      "user-123",
			Modules:     []string{"module1"},
			Actions:     []string{"read"},
			Level:       "admin",
			Description: "test permission",
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		assert.Equal(t, "123", p.GetID())
		assert.Equal(t, "user-123", p.GetUserID())
		assert.Equal(t, []string{"module1"}, p.GetModules())
		assert.Equal(t, []string{"read"}, p.GetActions())
		assert.Equal(t, "admin", p.GetLevel())
		assert.Equal(t, "test permission", p.GetDescription())
		assert.Equal(t, now, p.GetCreatedAt())
		assert.Equal(t, now, p.GetUpdatedAt())
	})

	t.Run("should return zero values when permission is nil", func(t *testing.T) {
		var p *Permission

		assert.Equal(t, "", p.GetID())
		assert.Equal(t, "", p.GetUserID())
		assert.Nil(t, p.GetModules())
		assert.Nil(t, p.GetActions())
		assert.Equal(t, "", p.GetLevel())
		assert.Equal(t, "", p.GetDescription())
		assert.True(t, p.GetCreatedAt().IsZero())
		assert.True(t, p.GetUpdatedAt().IsZero())
	})
}

func TestPermission_PullDomainEvents(t *testing.T) {
	t.Run("should return events and clear them", func(t *testing.T) {
		p := &Permission{
			ID: "123",
		}
		event := permission_event.NewPermissionCreatedEvent("123", "user-123", []string{"module1"}, []string{"read"}, "admin", "test")
		p.AddDomainEvent(event)

		events := p.PullDomainEvents()

		assert.Len(t, events, 1)
		assert.Equal(t, event, events[0])
		assert.Empty(t, p.PullDomainEvents())
	})

	t.Run("should return nil when permission is nil", func(t *testing.T) {
		var p *Permission

		events := p.PullDomainEvents()

		assert.Nil(t, events)
	})
}
