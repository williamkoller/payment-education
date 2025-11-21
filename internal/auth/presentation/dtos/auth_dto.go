package auth_dtos

type AuthDto struct {
	Email    string `json:"email" validate:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" validate:"required,min=8" example:"strongPassword123"`
}
