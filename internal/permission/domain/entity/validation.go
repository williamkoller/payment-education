package permission_entity

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(v.Errors, ", "))
}

func ValidatePermission(permission *Permission) (*Permission, error) {
	var errs []string

	if permission.ID == "" {
		errs = append(errs, "id is required")
	}

	if permission.UserID == "" {
		errs = append(errs, "user id is required")
	}

	if len(permission.Modules) == 0 {
		errs = append(errs, "modules is required")
	}

	if len(permission.Actions) == 0 {
		errs = append(errs, "actions is required")
	}

	if permission.Level == "" {
		errs = append(errs, "level is required")
	}

	if len(errs) > 0 {
		return nil, &ValidationError{Errors: errs}
	}

	return permission, nil
}
