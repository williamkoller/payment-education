package user_repository

import (
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	user_model "github.com/williamkoller/system-education/internal/user/models"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

var _ port_user_repository.UserRepository = &UserGormRepository{}

func NewUserGormRepository(db *gorm.DB) *UserGormRepository {
	return &UserGormRepository{db: db}
}

func (r *UserGormRepository) Save(u *user_entity.User) (*user_entity.User, error) {
	model := user_model.FromEntity(u)

	if err := r.db.Create(&model).Error; err != nil {
		return nil, err
	}

	return user_model.ToEntity(model), nil
}

func (r *UserGormRepository) FindByID(id string) (*user_entity.User, error) {
	var model user_model.User

	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return user_model.ToEntity(&model), nil
}

func (r *UserGormRepository) FindAll() ([]*user_entity.User, error) {
	var models []*user_model.User

	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}

	return user_model.ToEntities(models), nil
}

func (r *UserGormRepository) Delete(id string) error {
	return r.db.Delete(&user_model.User{}, "id = ?", id).Error
}

func (r *UserGormRepository) FindByEmail(email string) (*user_entity.User, error) {
	var model user_model.User

	if err := r.db.First(&model, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user_model.ToEntity(&model), nil
}
