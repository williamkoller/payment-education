package user_repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	user_model "github.com/williamkoller/system-education/internal/user/model"
	user_repository "github.com/williamkoller/system-education/internal/user/repository"
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
