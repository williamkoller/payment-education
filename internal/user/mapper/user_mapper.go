package user_mapper

import (
	"time"

	user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"
)

type UserResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Surname     string   `json:"surname"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email"`
	Age         int32    `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToUser(d *user_entity.User) *UserResponse {
	return &UserResponse{
		ID:          d.ID,
		Name:        d.Name,
		Surname:     d.Surname,
		Nickname:    d.Nickname,
		Email:       d.Email,
		Age:         d.Age,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func ToUsers(users []*user_entity.User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))

	for _, u := range users {
		responses = append(responses, ToUser(u))
	}

	return responses
}
