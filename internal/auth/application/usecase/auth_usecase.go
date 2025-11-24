package auth_usecase

import (
	"errors"

	port_auth_cryptography "github.com/williamkoller/system-education/internal/auth/port/cryptography"
	port_auth_usecase "github.com/williamkoller/system-education/internal/auth/port/usecase"
	port_permission_repository "github.com/williamkoller/system-education/internal/permission/port/repository"
	port_cryptography "github.com/williamkoller/system-education/internal/user/port/cryptography"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
)

type AuthUsecase struct {
	repo            port_user_repository.UserRepository
	permissionRepo  port_permission_repository.PermissionRepository
	jwtTokenManager port_auth_cryptography.TokenManager
	passwordHasher  port_cryptography.Bcrypt
}

var _ port_auth_usecase.AuthUsecase = &AuthUsecase{}

func NewAuthUsecase(
	repo port_user_repository.UserRepository,
	permissionRepo port_permission_repository.PermissionRepository,
	jwtTokenManager port_auth_cryptography.TokenManager,
	passwordHasher port_cryptography.Bcrypt,
) *AuthUsecase {
	return &AuthUsecase{
		repo:            repo,
		permissionRepo:  permissionRepo,
		jwtTokenManager: jwtTokenManager,
		passwordHasher:  passwordHasher,
	}
}

func (a *AuthUsecase) Login(email, password string) (string, error) {
	user, err := a.repo.FindByEmail(email)

	if err != nil {
		return "", errors.New("user not found")
	}

	if _, err := a.passwordHasher.HashComparer(password, user.Password); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Fetch user permissions and extract modules
	modules := []string{}
	permissions, err := a.permissionRepo.FindPermissionByUserID(user.ID)
	if err == nil && len(permissions) > 0 {
		for _, permission := range permissions {
			modules = append(modules, permission.GetModules()...)
		}
	}

	userSign := map[string]interface{}{
		"user_id":    user.ID,
		"name":       user.Name,
		"nick_name":  user.Nickname,
		"email":      user.Email,
		"modules":    modules,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	token, err := a.jwtTokenManager.Sign(userSign)

	if err != nil {
		return "", errors.New("error in generate token")
	}

	return token, nil
}

func (a *AuthUsecase) Profile(email string) (string, error) {
	user, err := a.repo.FindByEmail(email)

	if err != nil {
		return "", errors.New("user not found")
	}

	return user.Email, nil
}
