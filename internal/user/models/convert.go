package user_model

import user_entity "github.com/williamkoller/system-education/internal/user/domain/entity"

func FromEntity(u *user_entity.User) *User {
	return &User{
		ID:          u.ID,
		Name:        u.Name,
		Surname:     u.Surname,
		Nickname:    u.Nickname,
		Age:         u.Age,
		Email:       u.Email,
		Password:    u.Password,
		Roles:       u.Roles,
		Permissions: u.Permissions,
	}
}

func FromEntities(us []*user_entity.User) []*User {
	models := make([]*User, 0, len(us))
	for _, u := range us {
		models = append(models, FromEntity(u))
	}
	return models
}

func ToEntity(u *User) *user_entity.User {
	return &user_entity.User{
		ID:          u.ID,
		Name:        u.Name,
		Surname:     u.Surname,
		Nickname:    u.Nickname,
		Age:         u.Age,
		Email:       u.Email,
		Password:    u.Password,
		Roles:       u.Roles,
		Permissions: u.Permissions,
	}
}

func ToEntities(us []*User) []*user_entity.User {
	entities := make([]*user_entity.User, 0, len(us))
	for _, u := range us {
		entities = append(entities, ToEntity(u))
	}
	return entities
}
