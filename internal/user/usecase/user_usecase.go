package user_usecase

import (
	"errors"

	"github.com/google/uuid"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	port_user_repository "github.com/williamkoller/system-education/internal/user/domain/port/repository"
	port_user_usecase "github.com/williamkoller/system-education/internal/user/domain/port/usecase"
	"github.com/williamkoller/system-education/internal/user/dtos"
)

type UserUsecase struct {
	repo port_user_repository.UserRepository
}

func NewUserUsecase(repo port_user_repository.UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

var _ port_user_usecase.UserUsecase = &UserUsecase{}

func (u *UserUsecase) Create(input dtos.AddUserDto) (*user_entity.User, error) {
	existingUser, err := u.repo.FindByEmail(input.Email)

	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	if err != nil && !errors.Is(err, port_user_repository.ErrUserNotFound) {
		return nil, err
	}

	newUser := user_entity.NewUser(&user_entity.User{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Surname:     input.Surname,
		Nickname:    input.Nickname,
		Age:         input.Age,
		Email:       input.Email,
		Password:    input.Password,
		Roles:       input.Roles,
		Permissions: input.Permissions,
	})

	if newUser == nil {
		return nil, errors.New("invalid user data")
	}

	err = u.repo.Save(newUser)

	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return newUser, nil

}

func (u *UserUsecase) FindAll() ([]*user_entity.User, error) {
	users, err := u.repo.FindAll()

	if err != nil {
		return nil, errors.New("error in findAll users")
	}

	return users, nil
}