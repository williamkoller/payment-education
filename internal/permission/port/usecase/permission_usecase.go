package port_permission_usecase

import (
	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
	permission_dtos "github.com/williamkoller/system-education/internal/permission/presentation/dtos"
)

type PermissionUsecase interface {
	Create(input permission_dtos.AddPermissionDto) (*permission_entity.Permission, error)
	FindAll() ([]*permission_entity.Permission, error)
}
