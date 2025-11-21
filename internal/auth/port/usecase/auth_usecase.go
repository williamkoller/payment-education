package port_auth_usecase

type AuthUsecase interface {
	Login(email, password string) (string, error)
	Profile(email string) (string, error)
}
