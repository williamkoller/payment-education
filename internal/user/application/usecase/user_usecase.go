package user_usecase

import (
	"errors"

	"github.com/google/uuid"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	port_cryptography "github.com/williamkoller/system-education/internal/user/port/cryptography"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
	port_user_usecase "github.com/williamkoller/system-education/internal/user/port/usecase"
	"github.com/williamkoller/system-education/internal/user/presentation/dtos"
)

type UserUsecase struct {
	repo   port_user_repository.UserRepository
	crypto port_cryptography.Bcrypt
}

func NewUserUsecase(repo port_user_repository.UserRepository, crypto port_cryptography.Bcrypt) *UserUsecase {
	return &UserUsecase{repo: repo, crypto: crypto}
}

var _ port_user_usecase.UserUsecase = &UserUsecase{}

func (u *UserUsecase) Create(input dtos.AddUserDto) (*user_entity.User, error) {
	existingUser, err := u.repo.FindByEmail(input.Email)

	if err == nil && existingUser != nil {
		return nil, port_user_repository.ErrUserAlreadyExists
	}

	hash, err := u.crypto.Hash(input.Password)

	if err != nil {
		return nil, err
	}

	newUser := user_entity.NewUser(&user_entity.User{
		ID:       uuid.New().String(),
		Name:     input.Name,
		Surname:  input.Surname,
		Nickname: input.Nickname,
		Age:      input.Age,
		Email:    input.Email,
		Password: hash,
	})

	if newUser == nil {
		return nil, errors.New("invalid user data")
	}

	user, err := u.repo.Save(newUser)

	if err != nil {
		return nil, errors.New("cannot create new user")
	}

	return user, nil

}

func (u *UserUsecase) FindAll() ([]*user_entity.User, error) {
	users, err := u.repo.FindAll()

	if err != nil {
		return nil, errors.New("error in findAll users")
	}

	return users, nil
}
