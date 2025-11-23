package permission_model

import (
	"time"

	"github.com/lib/pq"
	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
	"gorm.io/gorm"
)

type Permission struct {
	ID          string `gorm:"primaryKey;type:uuid"`
	UserID      string
	Modules     pq.StringArray `gorm:"type:text[]"`
	Actions     pq.StringArray `gorm:"type:text[]"`
	Level       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (Permission) TableName() string {
	return "permissions"
}

func FromEntity(p *permission_entity.Permission) *Permission {
	if p == nil {
		return nil
	}
	return &Permission{
		ID:          p.ID,
		UserID:      p.UserID,
		Modules:     pq.StringArray(p.Modules),
		Actions:     pq.StringArray(p.Actions),
		Level:       p.Level,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func FromEntities(ps []*permission_entity.Permission) []*Permission {
	models := make([]*Permission, 0, len(ps))
	for _, p := range ps {
		models = append(models, FromEntity(p))
	}
	return models
}

func ToEntity(p *Permission) *permission_entity.Permission {
	if p == nil {
		return nil
	}
	return &permission_entity.Permission{
		ID:          p.ID,
		UserID:      p.UserID,
		Modules:     []string(p.Modules),
		Actions:     []string(p.Actions),
		Level:       p.Level,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToEntities(ps []*Permission) []*permission_entity.Permission {
	entities := make([]*permission_entity.Permission, 0, len(ps))
	for _, p := range ps {
		entities = append(entities, ToEntity(p))
	}
	return entities
}
