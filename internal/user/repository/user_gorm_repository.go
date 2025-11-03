package user_repository

import (
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	user_model "github.com/williamkoller/system-education/internal/user/model"
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
	var user *user_entity.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserGormRepository) FindAll() ([]*user_entity.User, error) {
	var users []*user_entity.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserGormRepository) Delete(id string) error {
	return r.db.Delete(&user_model.User{}, "id = ?", id).Error
}

func (r *UserGormRepository) FindByEmail(email string) (*user_entity.User, error) {
	var user *user_entity.User
	model := user_model.FromEntity(user)

	if err := r.db.First(&model, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user_model.ToEntity(model), nil
}
