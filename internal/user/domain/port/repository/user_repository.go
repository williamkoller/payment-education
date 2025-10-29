package port_user_repository

import (
	"errors"

	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
)

type UserRepository interface {
	Save(u *user_entity.User) error
	FindByID(id string) (*user_entity.User, error)
	FindAll() ([]*user_entity.User, error)
	Delete(id string) error
	FindByEmail(email string) (*user_entity.User, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)
