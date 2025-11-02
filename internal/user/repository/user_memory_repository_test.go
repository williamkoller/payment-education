package user_repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/port/repository"
	user_repository "github.com/williamkoller/system-education/internal/user/repository"
)

func createTestUser(id string) *user_entity.User {
	return &user_entity.User{
		ID:       id,
		Name:     "Test",
		Surname:  "User",
		Nickname: "tester",
		Age:      25,
		Email:    id + "@example.com",
		Password: "password123",
		Roles:    []string{"user"},
	}
}

func TestUserMemoryRepository_SaveAndFindByID(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()
	user := createTestUser("u1")

	_, err := repo.Save(user)
	assert.NoError(t, err)

	stored, err := repo.FindByID("u1")
	assert.NoError(t, err)
	assert.Equal(t, "Test", stored.GetName())
	assert.Equal(t, "u1", stored.GetID())
}

func TestUserMemoryRepository_FindByID_NotFound(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()

	_, err := repo.FindByID("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, port_user_repository.ErrUserNotFound, err)
}

func TestUserMemoryRepository_FindAll(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()

	repo.Save(createTestUser("u1"))
	repo.Save(createTestUser("u2"))

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserMemoryRepository_Delete(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()

	repo.Save(createTestUser("u1"))

	err := repo.Delete("u1")
	assert.NoError(t, err)

	_, err = repo.FindByID("u1")
	assert.Error(t, err)
	assert.Equal(t, port_user_repository.ErrUserNotFound, err)
}

func TestUserMemoryRepository_Delete_NotFound(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()
	user := &user_entity.User{
		ID:       "u1",
		Name:     "Test",
		Surname:  "User",
		Nickname: "tester",
		Age:      25,
		Email:    "u1@example.com",
		Password: "password123",
		Roles:    []string{"user"},
	}

	_, err := repo.Save(user)
	assert.NoError(t, err)

	err = repo.Delete("u1")
	assert.NoError(t, err)

	err = repo.Delete("u1")
	assert.Error(t, err)
	assert.Equal(t, port_user_repository.ErrUserNotFound, err)
}
