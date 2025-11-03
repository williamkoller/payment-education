package user_repository

import (
	"sync"

	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/port/repository"
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

func (r *UserMemoryRepository) Save(u *user_entity.User) (*user_entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[u.ID] = u
	return u, nil
}

func (r *UserMemoryRepository) FindByID(id string) (*user_entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, exists := r.store[id]
	if !exists {
		return nil, port_user_repository.ErrUserNotFound
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
		return port_user_repository.ErrUserNotFound
	}

	delete(r.store, id)
	return nil
}

func (r *UserMemoryRepository) FindByEmail(email string) (*user_entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.store {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, port_user_repository.ErrUserNotFound
}
