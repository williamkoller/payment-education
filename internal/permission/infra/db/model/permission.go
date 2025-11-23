package permission_model

import (
	"fmt"

	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	UserID      string
	Modules     []string `gorm:"serializer:json"`
	Actions     []string `gorm:"serializer:json"`
	Level       string
	Description string
}

func (Permission) TableName() string {
	return "permissions"
}

func FromEntity(p *permission_entity.Permission) *Permission {
	if p == nil {
		return nil
	}
	return &Permission{
		Model:       gorm.Model{ID: stringToUint(p.ID)},
		UserID:      p.UserID,
		Modules:     p.Modules,
		Actions:     p.Actions,
		Level:       p.Level,
		Description: p.Description,
	}
}

func stringToUint(s string) uint {
	if s == "" {
		return 0
	}
	var id uint
	fmt.Sscanf(s, "%d", &id)
	return id
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
		ID:          fmt.Sprintf("%d", p.ID),
		UserID:      p.UserID,
		Modules:     p.Modules,
		Actions:     p.Actions,
		Level:       p.Level,
		Description: p.Description,
	}
}

func ToEntities(ps []*Permission) []*permission_entity.Permission {
	entities := make([]*permission_entity.Permission, 0, len(ps))
	for _, p := range ps {
		entities = append(entities, ToEntity(p))
	}
	return entities
}
