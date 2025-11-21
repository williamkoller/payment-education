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
				Name:     "John",
				Surname:  "Doe",
				Age:      30,
				Email:    "john@example.com",
				Nickname: "jdoe",
				Password: "secret123",
			},
			wantErr: false,
		},
		{
			name: "missing fields",
			input: &User{
				Name:     "",
				Surname:  "",
				Age:      -1,
				Email:    "invalidemail",
				Nickname: "",
				Password: "",
			},
			wantErr: true,
		},
		{
			name: "invalid email format",
			input: &User{
				ID:       "123",
				Name:     "Jane",
				Surname:  "Doe",
				Age:      25,
				Email:    "jane[at]email.com",
				Nickname: "janed",
				Password: "pass",
			},
			wantErr: true,
		},
		{
			name: "invalid email is required",
			input: &User{
				ID:       "123",
				Name:     "Jane",
				Surname:  "Doe",
				Age:      25,
				Email:    "",
				Nickname: "janed",
				Password: "pass",
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

func TestValidationUpdateUser(t *testing.T) {
	tests := []struct {
		name    string
		input   *User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid update user",
			input: &User{
				ID:       "123",
				Name:     "John",
				Surname:  "Doe",
				Age:      30,
				Email:    "john@example.com",
				Nickname: "jdoe",
				Password: "secret123",
			},
			wantErr: false,
		},
		{
			name: "negative age",
			input: &User{
				ID:       "123",
				Name:     "John",
				Surname:  "Doe",
				Age:      -5,
				Email:    "john@example.com",
				Nickname: "jdoe",
				Password: "secret123",
			},
			wantErr: true,
			errMsg:  "age cannot be negative",
		},
		{
			name: "invalid email format",
			input: &User{
				ID:       "123",
				Name:     "Jane",
				Surname:  "Doe",
				Age:      25,
				Email:    "invalid-email",
				Nickname: "janed",
				Password: "pass",
			},
			wantErr: true,
			errMsg:  "email is invalid",
		},
		{
			name: "both age and email invalid",
			input: &User{
				ID:       "123",
				Name:     "Bob",
				Surname:  "Smith",
				Age:      -10,
				Email:    "bademail",
				Nickname: "bobs",
				Password: "pass",
			},
			wantErr: true,
			errMsg:  "age cannot be negative",
		},
		{
			name: "age zero is valid",
			input: &User{
				ID:       "123",
				Name:     "Baby",
				Surname:  "User",
				Age:      0,
				Email:    "baby@example.com",
				Nickname: "baby",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "email with @ is valid",
			input: &User{
				ID:       "123",
				Name:     "Valid",
				Surname:  "User",
				Age:      20,
				Email:    "valid@test.com",
				Nickname: "valid",
				Password: "pass",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidationUpdateUser(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.input, result)
			}
		})
	}
}
