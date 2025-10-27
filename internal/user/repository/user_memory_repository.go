package user_repository

import (
	"errors"
	"sync"

	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	port_user_repository "github.com/williamkoller/system-education/internal/user/domain/port/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserMemoryRepository struct {
	mu    sync.RWMutex
	store map[string]*user_entity.User
}

var _ port_user_repository.UserRepository = &UserMemoryRepository{}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		store: make(map[string]*user_entity.User),
	}
}

func (r *UserMemoryRepository) Save(u *user_entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[u.ID] = u
	return nil
}

func (r *UserMemoryRepository) FindByID(id string) (*user_entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, exists := r.store[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (r *UserMemoryRepository) FindAll() ([]*user_entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := make([]*user_entity.User, 0, len(r.store))

	for _, u := range r.store {
		users = append(users, u)
	}

	return users, nil
}

func (r *UserMemoryRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.store[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.store, id)
	return nil
}
