package event

import "time"

type UserCreatedEvent struct {
	UserID string
	Name   string
	Email  string
	Date   time.Time
}

func NewUserCreatedEvent(id, name, email string) *UserCreatedEvent {
	return &UserCreatedEvent{UserID: id, Name: name, Email: email}
}

func(e *UserCreatedEvent) EventName() string {
	return "user.created"
}

func (e *UserCreatedEvent) OccurredOn() time.Time {
	return e.Date
}