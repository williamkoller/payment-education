package port_user_repository

import user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"

type UserRepository interface {
	Save(u *user_entity.User) error
	FindByID(id string) (*user_entity.User, error)
	FindAll() ([]*user_entity.User, error)
	Delete(id string) error
}
