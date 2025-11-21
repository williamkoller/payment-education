package permission_event

import "time"

type PermissionCreatedEvent struct {
	PermissionID string
	UserID       string
	Modules      []string
	Actions      []string
	Level        string
	Description  string
	Date         time.Time
}

func NewPermissionCreatedEvent(permissionID string, userID string, modules []string, actions []string, level string, description string) *PermissionCreatedEvent {
	return &PermissionCreatedEvent{
		PermissionID: permissionID,
		UserID:       userID,
		Modules:      modules,
		Actions:      actions,
		Level:        level,
		Description:  description,
		Date:         time.Now(),
	}
}

func (e *PermissionCreatedEvent) EventName() string {
	return "permission.created"
}

func (e *PermissionCreatedEvent) OccurredOn() time.Time {
	return e.Date
}
