package user_usecase

import (
	"errors"
	"fmt"
	"log"

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

func (u *UserUsecase) FindByID(id string) (*user_entity.User, error) {
	user, err := u.repo.FindByID(id)

	if err != nil {
		return nil, errors.New("error in findByID user")
	}

	return user, nil
}

func (u *UserUsecase) Update(id string, input dtos.UpdateUserDto) (*user_entity.User, error) {
	log.Printf("Updating user with ID: %s, Data: %+v", id, input)

	userExists, err := u.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	if userExists == nil {
		return nil, errors.New("user not found")
	}

	if input.Password != nil {
		if *input.Password == "" {
		} else {
			hash, err := u.crypto.Hash(*input.Password)
			if err != nil {
				return nil, err
			}
			input.Password = &hash
		}
	}

	userExists.UpdateUser(
		input.Name,
		input.Nickname,
		input.Email,
		input.Password,
		input.Age,
	)

	updatedUser, err := u.repo.Update(userExists.ID, userExists)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

func (u *UserUsecase) Delete(id string) error {
	userExists, err := u.FindByID(id)

	if err != nil {
		return errors.New("user not dound")
	}

	err = u.repo.Delete(userExists.ID)

	if err != nil {
		return errors.New("error in delete user")
	}

	return nil
}
