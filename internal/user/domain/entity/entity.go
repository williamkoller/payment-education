package user_entity

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
}

func NewUser(u *User) *User {
	user, err := ValidationUser(u)

	if err != nil {
		return nil
	}

	return &User{
		ID:          user.ID,
		Name:        user.Name,
		Surname:     user.Surname,
		Nickname: user.Nickname,
		Age:         user.Age,
		Email:       user.Email,
		Password:    user.Password,
		Roles:       user.Roles,
		Permissions: user.Permissions,
	}
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
