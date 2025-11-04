package user_model

import (
	userEntity "github.com/williamkoller/system-education/internal/user/domain/entity"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       string
	Name     string
	Surname  string
	Nickname string
	Age      int32
	Email    string `gorm:"uniqueIndex"`
	Password string
}

func (User) TableName() string {
	return "users"
}

func FromEntity(u *userEntity.User) *User {
	if u == nil {
		return nil
	}
	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		Nickname: u.Nickname,
		Age:      u.Age,
		Email:    u.Email,
		Password: u.Password,
	}
}

func FromEntities(us []*userEntity.User) []*User {
	models := make([]*User, 0, len(us))
	for _, u := range us {
		models = append(models, FromEntity(u))
	}
	return models
}

func ToEntity(u *User) *userEntity.User {
	if u == nil {
		return nil
	}
	return &userEntity.User{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		Nickname: u.Nickname,
		Age:      u.Age,
		Email:    u.Email,
		Password: u.Password,
	}
}

func ToEntities(us []*User) []*userEntity.User {
	entities := make([]*userEntity.User, 0, len(us))
	for _, u := range us {
		entities = append(entities, ToEntity(u))
	}
	return entities
}
