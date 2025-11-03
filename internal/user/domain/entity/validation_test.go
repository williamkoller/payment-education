package user_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationUser(t *testing.T) {
	tests := []struct {
		name    string
		input   *User
		wantErr bool
	}{
		{
			name: "valid user",
			input: &User{
				Name:        "John",
				Surname:     "Doe",
				Age:         30,
				Email:       "john@example.com",
				Nickname:    "jdoe",
				Password:    "secret123",
			},
			wantErr: false,
		},
		{
			name: "missing fields",
			input: &User{
				Name:        "",
				Surname:     "",
				Age:         -1,
				Email:       "invalidemail",
				Nickname:    "",
				Password:    "",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			input: &User{
				ID:          "123",
				Name:        "Jane",
				Surname:     "Doe",
				Age:         25,
				Email:       "jane[at]email.com",
				Nickname:    "janed",
				Password:    "pass",
			},
			wantErr: true,
		},
		{
			name: "invalid email is required",
			input: &User{
				ID:          "123",
				Name:        "Jane",
				Surname:     "Doe",
				Age:         25,
				Email:       "",
				Nickname:    "janed",
				Password:    "pass",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidationUser(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidationUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	ve := &ValidationError{
		Errors: []string{"name is required", "email is invalid"},
	}

	expected := "validation failed: name is required, email is invalid"
	assert.Equal(t, expected, ve.Error())
}
