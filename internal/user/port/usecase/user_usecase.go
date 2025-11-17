package port_user_usecase

import (
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/presentation/dtos"
)

type UserUsecase interface {
	Create(input dtos.AddUserDto) (*user_entity.User, error)
	FindAll() ([]*user_entity.User, error)
	FindByID(id string) (*user_entity.User, error)
	Update(id string, input dtos.UpdateUserDto) (*user_entity.User, error)
	Delete(id string) error
}
