package permission_entity

import (
	"time"

	permission_event "github.com/williamkoller/system-education/internal/permission/domain/event"
	sharedEvent "github.com/williamkoller/system-education/shared/domain/event"
)

type Permission struct {
	sharedEvent.AggregateRoot
	ID          string
	UserID      string
	Modules     []string
	Actions     []string
	Level       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

func NewPermission(p *Permission) (*Permission, error) {
	vp, err := ValidatePermission(p)
	if err != nil {
		return nil, err
	}
	permission := &Permission{
		ID:          vp.ID,
		UserID:      vp.UserID,
		Modules:     vp.Modules,
		Actions:     vp.Actions,
		Level:       vp.Level,
		Description: vp.Description,
		CreatedAt:   vp.CreatedAt,
		UpdatedAt:   vp.UpdatedAt,
		DeletedAt:   vp.DeletedAt,
	}

	permission.AddDomainEvent(permission_event.NewPermissionCreatedEvent(permission.ID, permission.UserID, permission.Modules, permission.Actions, permission.Level, permission.Description))

	return permission, nil
}

func (p *Permission) GetID() string {
	return p.ID
}

func (p *Permission) GetUserID() string {
	return p.UserID
}

func (p *Permission) GetModules() []string {
	return p.Modules
}

func (p *Permission) GetActions() []string {
	return p.Actions
}

func (p *Permission) GetLevel() string {
	return p.Level
}

func (p *Permission) GetDescription() string {
	return p.Description
}

func (p *Permission) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Permission) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

func (p *Permission) GetDeletedAt() time.Time {
	return p.DeletedAt
}

func (p *Permission) PullDomainEvents() []sharedEvent.Event {
	if p == nil {
		return nil
	}
	return p.AggregateRoot.PullDomainEvents()
}
