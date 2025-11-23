package user_entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	t.Run("should create user when valid", func(t *testing.T) {
		input := &User{
			ID:       "123",
			Name:     "Alice",
			Surname:  "Smith",
			Nickname: "alices",
			Age:      28,
			Email:    "alice@example.com",
			Password: "secret",
		}

		user := NewUser(input)
		assert.NotNil(t, user)
		assert.Equal(t, "Alice", user.Name)
		assert.Equal(t, "alices", user.Nickname)
	})

	t.Run("should return nil when invalid", func(t *testing.T) {
		input := &User{
			ID:       "456",
			Name:     "",
			Surname:  "",
			Nickname: "",
			Age:      -1,
			Email:    "invalid-email",
			Password: "",
		}

		user := NewUser(input)
		assert.Nil(t, user)
	})
}

type MockEvent struct{}

func (m MockEvent) EventName() string {
	return "MockEvent"
}

func (m MockEvent) OccurredOn() time.Time {
	return time.Now()
}

func TestPullDomainEvents(t *testing.T) {
	tests := []struct {
		name           string
		user           *User
		setup          func(*User)
		expectedEvents int
		expectedName   string
	}{
		{
			name: "should return events and clear them",
			user: &User{
				ID:    "1",
				Name:  "Alice",
				Email: "alice@example.com",
			},
			setup: func(u *User) {
				u.AddDomainEvent(MockEvent{})
			},
			expectedEvents: 1,
			expectedName:   "MockEvent",
		},
		{
			name: "should return empty when no events",
			user: &User{
				ID:    "1",
				Name:  "Alice",
				Email: "alice@example.com",
			},
			setup:          func(u *User) {},
			expectedEvents: 0,
		},
		{
			name:           "should return nil when user is nil",
			user:           nil,
			setup:          func(u *User) {},
			expectedEvents: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil && tt.user != nil {
				tt.setup(tt.user)
			}

			events := tt.user.PullDomainEvents()

			if tt.user == nil {
				assert.Nil(t, events)
			} else {
				assert.Len(t, events, tt.expectedEvents)
				if tt.expectedEvents > 0 {
					assert.Equal(t, tt.expectedName, events[0].EventName())
					assert.Empty(t, tt.user.PullDomainEvents())
				}
			}
		})
	}
}

func TestUser_Getters(t *testing.T) {
	user := &User{
		ID:       "abc123",
		Name:     "Carlos",
		Surname:  "Silva",
		Nickname: "casilva",
		Age:      35,
		Email:    "carlos@example.com",
		Password: "secret123",
	}

	assert.Equal(t, "abc123", user.GetID())
	assert.Equal(t, "Carlos", user.GetName())
	assert.Equal(t, "Silva", user.GetSurname())
	assert.Equal(t, int32(35), user.GetAge())
	assert.Equal(t, "carlos@example.com", user.GetEmail())
	assert.Equal(t, "secret123", user.GetPassword())
}

func TestUpdateUser(t *testing.T) {
	name := "New Name"
	nickname := "new_nick"
	email := "new@example.com"
	password := "new_pass"
	age := int32(30)
	invalidEmail := "invalid"
	negativeAge := int32(-1)

	tests := []struct {
		name  string
		user  *User
		input struct {
			name     *string
			nickname *string
			email    *string
			password *string
			age      *int32
		}
		expected    *User
		expectError bool
		errMessage  string
	}{
		{
			name: "should update all fields successfully",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				name:     &name,
				nickname: &nickname,
				email:    &email,
				password: &password,
				age:      &age,
			},
			expected: &User{
				ID:       "123",
				Name:     name,
				Surname:  "User",
				Nickname: nickname,
				Age:      age,
				Email:    email,
				Password: password,
			},
			expectError: false,
		},
		{
			name: "should update only name",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				name: &name,
			},
			expected: &User{
				ID:       "123",
				Name:     name,
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			expectError: false,
		},
		{
			name: "should update only nickname",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				nickname: &nickname,
			},
			expected: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: nickname,
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			expectError: false,
		},
		{
			name: "should update only email",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				email: &email,
			},
			expected: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    email,
				Password: "oldpass",
			},
			expectError: false,
		},
		{
			name: "should update only password",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				password: &password,
			},
			expected: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: password,
			},
			expectError: false,
		},
		{
			name: "should update only age",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				age: &age,
			},
			expected: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      age,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			expectError: false,
		},
		{
			name: "should not update any fields when all inputs are nil",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{},
			expected: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			expectError: false,
		},
		{
			name: "should return error when email is invalid",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				email: &invalidEmail,
			},
			expected:    nil,
			expectError: true,
			errMessage:  "email is invalid",
		},
		{
			name: "should return error when age is negative",
			user: &User{
				ID:       "123",
				Name:     "Old",
				Surname:  "User",
				Nickname: "old",
				Age:      20,
				Email:    "old@example.com",
				Password: "oldpass",
			},
			input: struct {
				name     *string
				nickname *string
				email    *string
				password *string
				age      *int32
			}{
				age: &negativeAge,
			},
			expected:    nil,
			expectError: true,
			errMessage:  "age cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := tt.user.UpdateUser(tt.input.name, tt.input.nickname, tt.input.email, tt.input.password, tt.input.age)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, updated)
				assert.Contains(t, err.Error(), tt.errMessage)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, updated)
				assert.Equal(t, tt.expected.Name, updated.Name)
				assert.Equal(t, tt.expected.Nickname, updated.Nickname)
				assert.Equal(t, tt.expected.Email, updated.Email)
				assert.Equal(t, tt.expected.Password, updated.Password)
				assert.Equal(t, tt.expected.Age, updated.Age)
			}
		})
	}
}
