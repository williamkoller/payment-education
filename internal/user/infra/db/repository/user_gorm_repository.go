package user_repository

import (
	userEntity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/infra/db/model"
	portUserRepository "github.com/williamkoller/system-education/internal/user/port/repository"
	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

var _ portUserRepository.UserRepository = &UserGormRepository{}

func NewUserGormRepository(db *gorm.DB) *UserGormRepository {
	return &UserGormRepository{db: db}
}

func (r *UserGormRepository) Save(u *userEntity.User) (*userEntity.User, error) {
	model := user_model.FromEntity(u)
	if err := r.db.Create(&model).Error; err != nil {
		return nil, err
	}

	return user_model.ToEntity(model), nil
}

func (r *UserGormRepository) FindByID(id string) (*userEntity.User, error) {
	var user *userEntity.User

	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserGormRepository) FindAll() ([]*userEntity.User, error) {
	var users []*userEntity.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserGormRepository) Delete(id string) error {
	return r.db.Delete(&user_model.User{}, "id = ?", id).Error
}

func (r *UserGormRepository) FindByEmail(email string) (*userEntity.User, error) {
	var user *userEntity.User
	model := user_model.FromEntity(user)

	if err := r.db.First(&model, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return user_model.ToEntity(model), nil
}
