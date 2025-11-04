package user_repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	userEntity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/infra/db/repository"
	portUserRepository "github.com/williamkoller/system-education/internal/user/port/repository"
)

func createTestUser(id string) *userEntity.User {
	return &userEntity.User{
		ID:       id,
		Name:     "Test",
		Surname:  "User",
		Nickname: "tester",
		Age:      25,
		Email:    id + "@example.com",
		Password: "password123",
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

func TestUserMemoryRepository_SaveAndFindByEmail(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()
	user := createTestUser("u1")

	_, err := repo.Save(user)
	assert.NoError(t, err)

	stored, err := repo.FindByEmail("u1@example.com")
	assert.NoError(t, err)
	assert.Equal(t, "Test", stored.GetName())
	assert.Equal(t, "u1", stored.GetID())
}

func TestUserMemoryRepository_FindByID_NotFound(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()

	_, err := repo.FindByID("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, portUserRepository.ErrUserNotFound, err)
}

func TestUserMemoryRepository_FindByEmail_NotFound(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()

	_, err := repo.FindByEmail("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, portUserRepository.ErrUserNotFound, err)
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
	assert.Equal(t, portUserRepository.ErrUserNotFound, err)
}

func TestUserMemoryRepository_Delete_NotFound(t *testing.T) {
	repo := user_repository.NewUserMemoryRepository()
	user := &userEntity.User{
		ID:       "u1",
		Name:     "Test",
		Surname:  "User",
		Nickname: "tester",
		Age:      25,
		Email:    "u1@example.com",
		Password: "password123",
	}

	_, err := repo.Save(user)
	assert.NoError(t, err)

	err = repo.Delete("u1")
	assert.NoError(t, err)

	err = repo.Delete("u1")
	assert.Error(t, err)
	assert.Equal(t, portUserRepository.ErrUserNotFound, err)
}
