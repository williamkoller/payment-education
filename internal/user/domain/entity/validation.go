package user_entity

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

func ValidationUser(u *User) (*User, error) {
	var errs []string

	if strings.TrimSpace(u.Name) == "" {
		errs = append(errs, "name is required")
	}

	if strings.TrimSpace(u.Surname) == "" {
		errs = append(errs, "surname is required")
	}

	if u.Age < 0 {
		errs = append(errs, "age cannot be negative")
	}

	if strings.TrimSpace(u.Email) == "" {
		errs = append(errs, "email is required")
	} else if !strings.Contains(u.Email, "@") {
		errs = append(errs, "email is invalid")
	}

	if strings.TrimSpace(u.Nickname) == "" {
		errs = append(errs, "nickname is required")
	}

	if strings.TrimSpace(u.Password) == "" {
		errs = append(errs, "password is required")
	}

	if len(u.Roles) == 0 {
		errs = append(errs, "roles must not be empty")
	}

	if len(u.Permissions) == 0 {
		errs = append(errs, "permissions must not be empty")
	}

	if len(errs) > 0 {
		return nil, &ValidationError{Errors: errs}
	}

	return u, nil
}
