package permission_mapper

import (
	"time"

	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
)

type PermissionResponse struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	Modules     []string  `json:"modules"`
	Actions     []string  `json:"actions"`
	Level       string    `json:"level"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func ToPermission(p *permission_entity.Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:          p.ID,
		UserID:      p.UserID,
		Modules:     p.Modules,
		Actions:     p.Actions,
		Level:       p.Level,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToPermissions(ps []*permission_entity.Permission) []*PermissionResponse {
	responses := make([]*PermissionResponse, 0, len(ps))
	for _, p := range ps {
		responses = append(responses, ToPermission(p))
	}
	return responses
}
