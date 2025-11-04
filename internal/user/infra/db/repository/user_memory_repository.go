package user_repository

import (
	"sync"

	userEntity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"github.com/williamkoller/system-education/internal/user/port/repository"
)

type UserMemoryRepository struct {
	mu    sync.RWMutex
	store map[string]*userEntity.User
}

var _ port_user_repository.UserRepository = &UserMemoryRepository{}

func NewUserMemoryRepository() *UserMemoryRepository {
	return &UserMemoryRepository{
		store: make(map[string]*userEntity.User),
	}
}

func (r *UserMemoryRepository) Save(u *userEntity.User) (*userEntity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.store[u.ID] = u
	return u, nil
}

func (r *UserMemoryRepository) FindByID(id string) (*userEntity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, exists := r.store[id]
	if !exists {
		return nil, port_user_repository.ErrUserNotFound
	}

	return user, nil
}

func (r *UserMemoryRepository) FindAll() ([]*userEntity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := make([]*userEntity.User, 0, len(r.store))

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

func (r *UserMemoryRepository) FindByEmail(email string) (*userEntity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.store {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, port_user_repository.ErrUserNotFound
}
