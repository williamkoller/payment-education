package auth_usecase

import (
	"errors"

	port_auth_cryptography "github.com/williamkoller/system-education/internal/auth/port/cryptography"
	port_auth_usecase "github.com/williamkoller/system-education/internal/auth/port/usecase"
	port_cryptography "github.com/williamkoller/system-education/internal/user/port/cryptography"
	port_user_repository "github.com/williamkoller/system-education/internal/user/port/repository"
)

type AuthUsecase struct {
	repo            port_user_repository.UserRepository
	jwtTokenManager port_auth_cryptography.TokenManager
	passwordHasher  port_cryptography.Bcrypt
}

var _ port_auth_usecase.AuthUsecase = &AuthUsecase{}

func NewAuthUsecase(
	repo port_user_repository.UserRepository,
	jwtTokenManager port_auth_cryptography.TokenManager,
	passwordHasher port_cryptography.Bcrypt,
) *AuthUsecase {
	return &AuthUsecase{
		repo:            repo,
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

	userSign := map[string]interface{}{
		"name":       user.Name,
		"nick_name":  user.Nickname,
		"email":      user.Email,
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