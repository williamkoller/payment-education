package user_usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	port_cryptography "github.com/williamkoller/system-education/internal/user/port/cryptography"
	port_event "github.com/williamkoller/system-education/internal/user/port/event"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
	port_user_usecase "github.com/williamkoller/system-education/internal/user/port/usecase"
	"github.com/williamkoller/system-education/internal/user/presentation/dtos"
)

type UserUsecase struct {
	repo   port_user_repository.UserRepository
	crypto port_cryptography.Bcrypt
	event  port_event.Dispacther
}

func NewUserUsecase(repo port_user_repository.UserRepository, crypto port_cryptography.Bcrypt, event port_event.Dispacther) *UserUsecase {
	return &UserUsecase{repo: repo, crypto: crypto, event: event}
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
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	for _, domainEvent := range newUser.PullDomainEvents() {
		u.event.Dispatch(domainEvent)
	}

	return user, nil

}

func (u *UserUsecase) FindAll() ([]*user_entity.User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find all users: %w", err)
	}

	return users, nil
}

func (u *UserUsecase) FindByID(id string) (*user_entity.User, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	user, err := u.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID %s: %w", id, err)
	}

	return user, nil
}

func (u *UserUsecase) Update(id string, input dtos.UpdateUserDto) (*user_entity.User, error) {
	if id == "" {
		return nil, errors.New("user ID cannot be empty")
	}

	log.Printf("Updating user with ID: %s, Data: %+v", id, input)

	userExists, err := u.FindByID(id)
	if err != nil {
		return nil, err
	}
	if userExists == nil {
		return nil, errors.New("user not found")
	}

	if input.Password != nil && *input.Password != "" {
		hash, err := u.crypto.Hash(*input.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		input.Password = &hash
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
	if id == "" {
		return errors.New("user ID cannot be empty")
	}

	userExists, err := u.FindByID(id)
	if err != nil {
		return err
	}

	err = u.repo.Delete(userExists.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
