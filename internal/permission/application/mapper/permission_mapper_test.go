package permission_mapper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
)

func TestToPermission(t *testing.T) {
	now := time.Now()
	permission := &permission_entity.Permission{
		ID:          "123",
		UserID:      "user-1",
		Modules:     []string{"module1"},
		Actions:     []string{"read"},
		Level:       "admin",
		Description: "test description",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := ToPermission(permission)

	assert.Equal(t, permission.ID, response.ID)
	assert.Equal(t, permission.UserID, response.UserID)
	assert.Equal(t, permission.Modules, response.Modules)
	assert.Equal(t, permission.Actions, response.Actions)
	assert.Equal(t, permission.Level, response.Level)
	assert.Equal(t, permission.Description, response.Description)
	assert.Equal(t, permission.CreatedAt, response.CreatedAt)
	assert.Equal(t, permission.UpdatedAt, response.UpdatedAt)
}

func TestToPermissions(t *testing.T) {
	now := time.Now()
	permissions := []*permission_entity.Permission{
		{
			ID:          "123",
			UserID:      "user-1",
			Modules:     []string{"module1"},
			Actions:     []string{"read"},
			Level:       "admin",
			Description: "test description 1",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          "456",
			UserID:      "user-2",
			Modules:     []string{"module2"},
			Actions:     []string{"write"},
			Level:       "user",
			Description: "test description 2",
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	responses := ToPermissions(permissions)

	assert.Len(t, responses, 2)
	assert.Equal(t, permissions[0].ID, responses[0].ID)
	assert.Equal(t, permissions[1].ID, responses[1].ID)
}
