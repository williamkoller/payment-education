package user_entity

import (
	"time"

	userEvent "github.com/williamkoller/system-education/internal/user/domain/event"
	sharedEvent "github.com/williamkoller/system-education/shared/domain/event"
)

type User struct {
	ID        string
	Name      string
	Surname   string
	Nickname  string
	Age       int32
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	sharedEvent.AggregateRoot
}

func NewUser(u *User) *User {
	vu, err := ValidationUser(u)

	if err != nil {
		return nil
	}

	user := &User{
		ID:       vu.ID,
		Name:     vu.Name,
		Surname:  vu.Surname,
		Nickname: vu.Nickname,
		Age:      vu.Age,
		Email:    vu.Email,
		Password: vu.Password,
	}

	user.AddDomainEvent(userEvent.NewUserCreatedEvent(user.ID, user.Name, user.Email))

	return user
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetSurname() string {
	return u.Surname
}

func (u *User) GetAge() int32 {
	return u.Age
}

func (u *User) GetEmail() string {
	return u.Email
}

func (u *User) GetPassword() string {
	return u.Password
}

func (u *User) PullDomainEvents() []sharedEvent.Event {
	if u == nil {
		return nil
	}
	return u.AggregateRoot.PullDomainEvents()
}

func (u *User) UpdateUser(name, nickname, email, password *string, age *int32) (*User, error) {
	if name != nil {
		u.Name = *name
	}
	if nickname != nil {
		u.Nickname = *nickname
	}
	if email != nil {
		u.Email = *email
	}
	if password != nil {
		u.Password = *password
	}
	if age != nil {
		u.Age = *age
	}

	uv, err := ValidationUpdateUser(u)
	if err != nil {
		return nil, err
	}

	return uv, nil
}
