package user_mapper

import user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"

type UserResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email"`
	Age         int32    `json:"age"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func ToResponse(d *user_entity.User) *UserResponse {
	return &UserResponse{
		ID: d.ID,
		Name: d.Name,
		Surname: d.Surname,
		Nickname: d.Nickname,
		Email: d.Email,
		Age: d.Age,
		Roles: d.Roles,
		Permissions: d.Permissions,
	}
}
