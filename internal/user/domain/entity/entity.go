package user_entity

import (
	user_event "github.com/williamkoller/system-education/internal/user/domain/event"
	shared_event "github.com/williamkoller/system-education/shared/domain/event"
)

type User struct {
	ID          string
	Name        string
	Surname     string
	Nickname    string
	Age         int32
	Email       string
	Password    string
	Roles       []string
	Permissions []string

	domainEvents []shared_event.Event
}

func NewUser(u *User) *User {
	vu, err := ValidationUser(u)

	if err != nil {
		return nil
	}

	user := &User{
		ID:          vu.ID,
		Name:        vu.Name,
		Surname:     vu.Surname,
		Nickname:    vu.Nickname,
		Age:         vu.Age,
		Email:       vu.Email,
		Password:    vu.Password,
		Roles:       vu.Roles,
		Permissions: vu.Permissions,
	}

	user.AddDomainEvent(user_event.NewUserCreatedEvent(user.ID, user.Name, user.Email))

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

func (u *User) GetRoles() []string {
	return u.Roles
}

func (u *User) GetPermissions() []string {
	return u.Permissions
}

func (u *User) AddDomainEvent(e shared_event.Event) {
	u.domainEvents = append(u.domainEvents, e)
}

func (u *User) PullDomainEvents() []shared_event.Event {
	if u == nil {
		return nil
	}
	events := u.domainEvents
	u.domainEvents = []shared_event.Event{}
	return events
}
