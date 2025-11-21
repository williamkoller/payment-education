package user_repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	user_model "github.com/williamkoller/system-education/internal/user/infra/db/model"
	user_repository "github.com/williamkoller/system-education/internal/user/infra/db/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&user_model.User{})
	assert.NoError(t, err)

	return db
}

func TestUserGormRepository_SaveAndFindByEmail_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// preparar entidade
	u := &user_entity.User{
		ID:       "id‑123",
		Name:     "Test",
		Surname:  "User",
		Nickname: "testuser",
		Age:      25,
		Email:    "test@example.com",
		Password: "pass123",
	}

	// salvar
	saved, err := repo.Save(u)
	assert.NoError(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, u.Email, saved.Email)

	// buscar por email
	found, err := repo.FindByEmail(u.Email)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, u.Email, found.Email)
	// você pode também validar Name, Nickname etc se quiser
}

func TestUserGormRepository_FindByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// buscar email que não existe
	found, err := repo.FindByEmail("unknown@example.com")
	// dependendo de como você trata “não encontrado”, GORM retorna gorm.ErrRecordNotFound
	// você pode verificar com errors.Is(err, gorm.ErrRecordNotFound)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUserGormRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// buscar email que não existe
	found, err := repo.FindByID("id‑123")
	// dependendo de como você trata “não encontrado”, GORM retorna gorm.ErrRecordNotFound
	// você pode verificar com errors.Is(err, gorm.ErrRecordNotFound)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func TestUserGormRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// inserir 2 usuários
	u1 := &user_entity.User{ID: "id1", Name: "A", Surname: "B", Nickname: "a", Age: 20, Email: "a@example.com", Password: "p1"}
	u2 := &user_entity.User{ID: "id2", Name: "C", Surname: "D", Nickname: "c", Age: 30, Email: "c@example.com", Password: "p2"}
	_, err := repo.Save(u1)
	assert.NoError(t, err)
	_, err = repo.Save(u2)
	assert.NoError(t, err)

	all, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, all, 2)
}

func TestUserGormRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	u1 := &user_entity.User{ID: "id1", Name: "A", Surname: "B", Nickname: "a", Age: 20, Email: "a@example.com", Password: "p1"}

	_, _ = repo.Save(u1)
	err := repo.Delete("id1")
	assert.NoError(t, err)

}

func TestUserGormRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	u1 := &user_entity.User{ID: "id1", Name: "A", Surname: "B", Nickname: "a", Age: 20, Email: "a@example.com", Password: "p1"}

	_, _ = repo.Save(u1)
	updateU := &user_entity.User{ID: u1.ID, Name: "AA", Surname: "Bbb", Nickname: "aaa", Age: 21, Email: "aa@example.com", Password: "p1"}

	userUpdate, err := repo.Update("id1", updateU)
	assert.NoError(t, err)
	assert.NotNil(t, userUpdate)

}

func TestUserGormRepository_Save_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	u := &user_entity.User{
		ID:       "test-id-123",
		Name:     "John",
		Surname:  "Doe",
		Nickname: "johndoe",
		Age:      30,
		Email:    "john@example.com",
		Password: "hashed-password",
	}

	saved, err := repo.Save(u)

	assert.NoError(t, err)
	assert.NotNil(t, saved)
	assert.Equal(t, u.Email, saved.Email)
	assert.Equal(t, u.Name, saved.Name)
}

func TestUserGormRepository_FindByID_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// First save a user
	u := &user_entity.User{
		ID:       "find-id-123",
		Name:     "Jane",
		Surname:  "Smith",
		Nickname: "janesmith",
		Age:      25,
		Email:    "jane@example.com",
		Password: "pass123",
	}
	_, err := repo.Save(u)
	assert.NoError(t, err)

	// Now find it
	found, err := repo.FindByID("find-id-123")

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, u.ID, found.ID)
	assert.Equal(t, u.Email, found.Email)
}

func TestUserGormRepository_FindAll_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Save multiple users
	u1 := &user_entity.User{ID: "user1", Name: "User", Surname: "One", Nickname: "u1", Age: 20, Email: "user1@example.com", Password: "p1"}
	u2 := &user_entity.User{ID: "user2", Name: "User", Surname: "Two", Nickname: "u2", Age: 30, Email: "user2@example.com", Password: "p2"}
	u3 := &user_entity.User{ID: "user3", Name: "User", Surname: "Three", Nickname: "u3", Age: 40, Email: "user3@example.com", Password: "p3"}

	_, err := repo.Save(u1)
	assert.NoError(t, err)
	_, err = repo.Save(u2)
	assert.NoError(t, err)
	_, err = repo.Save(u3)
	assert.NoError(t, err)

	// Find all
	all, err := repo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, all)
	assert.GreaterOrEqual(t, len(all), 3)
}

func TestUserGormRepository_FindAll_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Find all on empty database
	all, err := repo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, all)
	assert.Len(t, all, 0)
}

func TestUserGormRepository_Update_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Save original user
	original := &user_entity.User{
		ID:       "update-id-123",
		Name:     "Original",
		Surname:  "Name",
		Nickname: "original",
		Age:      25,
		Email:    "original@example.com",
		Password: "oldpass",
	}
	_, err := repo.Save(original)
	assert.NoError(t, err)

	// Update the user
	updated := &user_entity.User{
		ID:       "update-id-123",
		Name:     "Updated",
		Surname:  "NewName",
		Nickname: "updated",
		Age:      26,
		Email:    "updated@example.com",
		Password: "newpass",
	}

	result, err := repo.Update("update-id-123", updated)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Updated", result.Name)
	assert.Equal(t, "updated@example.com", result.Email)
}


func TestUserGormRepository_Save_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Close the database to force an error
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	u := &user_entity.User{
		ID:       "error-test",
		Name:     "Test",
		Surname:  "User",
		Nickname: "test",
		Age:      25,
		Email:    "test@example.com",
		Password: "pass123",
	}

	saved, err := repo.Save(u)

	assert.Error(t, err)
	assert.Nil(t, saved)
}

func TestUserGormRepository_FindAll_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Close the database to force an error
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	users, err := repo.FindAll()

	assert.Error(t, err)
	assert.Nil(t, users)
}

func TestUserGormRepository_Update_Error(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Close the database to force an error
	sqlDB, err := db.DB()
	assert.NoError(t, err)
	sqlDB.Close()

	updated := &user_entity.User{
		ID:       "test-id",
		Name:     "Updated",
		Surname:  "User",
		Nickname: "updated",
		Age:      30,
		Email:    "updated@example.com",
		Password: "newpass",
	}

	result, err := repo.Update("test-id", updated)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// Additional success tests for better coverage

func TestUserGormRepository_Save_DuplicateID(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	u1 := &user_entity.User{
		ID:       "duplicate-id",
		Name:     "First",
		Surname:  "User",
		Nickname: "first",
		Age:      25,
		Email:    "first@example.com",
		Password: "pass123",
	}

	// Save first user
	_, err := repo.Save(u1)
	assert.NoError(t, err)

	// Try to save another user with same ID
	u2 := &user_entity.User{
		ID:       "duplicate-id",
		Name:     "Second",
		Surname:  "User",
		Nickname: "second",
		Age:      30,
		Email:    "second@example.com",
		Password: "pass456",
	}

	_, err = repo.Save(u2)
	// Should get error due to duplicate primary key
	assert.Error(t, err)
}

func TestUserGormRepository_FindAll_MultipleUsers(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Save multiple users
	users := []*user_entity.User{
		{ID: "u1", Name: "User1", Surname: "S1", Nickname: "n1", Age: 20, Email: "u1@test.com", Password: "p1"},
		{ID: "u2", Name: "User2", Surname: "S2", Nickname: "n2", Age: 25, Email: "u2@test.com", Password: "p2"},
		{ID: "u3", Name: "User3", Surname: "S3", Nickname: "n3", Age: 30, Email: "u3@test.com", Password: "p3"},
		{ID: "u4", Name: "User4", Surname: "S4", Nickname: "n4", Age: 35, Email: "u4@test.com", Password: "p4"},
	}

	for _, u := range users {
		_, err := repo.Save(u)
		assert.NoError(t, err)
	}

	// Find all
	all, err := repo.FindAll()

	assert.NoError(t, err)
	assert.NotNil(t, all)
	assert.GreaterOrEqual(t, len(all), 4)
}

func TestUserGormRepository_Update_NonExistentUser(t *testing.T) {
	db := setupTestDB(t)
	repo := user_repository.NewUserGormRepository(db)

	// Try to update a user that doesn't exist
	updated := &user_entity.User{
		ID:       "non-existent",
		Name:     "Updated",
		Surname:  "User",
		Nickname: "updated",
		Age:      30,
		Email:    "updated@example.com",
		Password: "newpass",
	}

	result, err := repo.Update("non-existent", updated)

	// GORM doesn't return error for updates that affect 0 rows
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
