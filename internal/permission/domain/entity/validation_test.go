package permission_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePermission(t *testing.T) {
	t.Run("should validate permission successfully", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			UserID:  "user-123",
			Modules: []string{"module1"},
			Actions: []string{"read"},
			Level:   "admin",
		}

		validatedPermission, err := ValidatePermission(p)

		assert.NoError(t, err)
		assert.NotNil(t, validatedPermission)
		assert.Equal(t, p, validatedPermission)
	})

	t.Run("should return error when id is empty", func(t *testing.T) {
		p := &Permission{
			UserID:  "user-123",
			Modules: []string{"module1"},
			Actions: []string{"read"},
			Level:   "admin",
		}

		validatedPermission, err := ValidatePermission(p)

		assert.Error(t, err)
		assert.Nil(t, validatedPermission)
		assert.IsType(t, &ValidationError{}, err)
		assert.Contains(t, err.Error(), "id is required")
	})

	t.Run("should return error when user id is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			Modules: []string{"module1"},
			Actions: []string{"read"},
			Level:   "admin",
		}

		validatedPermission, err := ValidatePermission(p)

		assert.Error(t, err)
		assert.Nil(t, validatedPermission)
		assert.IsType(t, &ValidationError{}, err)
		assert.Contains(t, err.Error(), "user id is required")
	})

	t.Run("should return error when modules is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			UserID:  "user-123",
			Actions: []string{"read"},
			Level:   "admin",
		}

		validatedPermission, err := ValidatePermission(p)

		assert.Error(t, err)
		assert.Nil(t, validatedPermission)
		assert.IsType(t, &ValidationError{}, err)
		assert.Contains(t, err.Error(), "modules is required")
	})

	t.Run("should return error when actions is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			UserID:  "user-123",
			Modules: []string{"module1"},
			Level:   "admin",
		}

		validatedPermission, err := ValidatePermission(p)

		assert.Error(t, err)
		assert.Nil(t, validatedPermission)
		assert.IsType(t, &ValidationError{}, err)
		assert.Contains(t, err.Error(), "actions is required")
	})

	t.Run("should return error when level is empty", func(t *testing.T) {
		p := &Permission{
			ID:      "123",
			UserID:  "user-123",
			Modules: []string{"module1"},
			Actions: []string{"read"},
		}

		validatedPermission, err := ValidatePermission(p)

		assert.Error(t, err)
		assert.Nil(t, validatedPermission)
		assert.IsType(t, &ValidationError{}, err)
		assert.Contains(t, err.Error(), "level is required")
	})

	t.Run("should return multiple errors", func(t *testing.T) {
		p := &Permission{}

		validatedPermission, err := ValidatePermission(p)

		assert.Error(t, err)
		assert.Nil(t, validatedPermission)
		assert.IsType(t, &ValidationError{}, err)
		assert.Contains(t, err.Error(), "id is required")
		assert.Contains(t, err.Error(), "user id is required")
		assert.Contains(t, err.Error(), "modules is required")
		assert.Contains(t, err.Error(), "actions is required")
		assert.Contains(t, err.Error(), "level is required")
	})
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Errors: []string{"error1", "error2"},
	}

	assert.Equal(t, "validation failed: error1, error2", err.Error())
}

func TestValidationUpdatePermission(t *testing.T) {
	t.Run("should validate update permission successfully", func(t *testing.T) {
		p := &Permission{
			Level: "admin",
		}

		validatedP, err := ValidationUpdatePermission(p)

		assert.NoError(t, err)
		assert.Equal(t, p, validatedP)
	})

	t.Run("should return error when level is empty", func(t *testing.T) {
		p := &Permission{
			Level: "",
		}

		validatedP, err := ValidationUpdatePermission(p)

		assert.Error(t, err)
		assert.Nil(t, validatedP)
		assert.Contains(t, err.Error(), "level cannot be empty")
	})
}
