package port_user_repository

import (
	"errors"

	userEntity "github.com/williamkoller/system-education/internal/user/domain/entity"
)

type UserRepository interface {
	Save(u *userEntity.User) (*userEntity.User, error)
	FindByID(id string) (*userEntity.User, error)
	FindAll() ([]*userEntity.User, error)
	Delete(id string) error
	FindByEmail(email string) (*userEntity.User, error)
	Update(id string, u *userEntity.User) (*userEntity.User, error)
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)
