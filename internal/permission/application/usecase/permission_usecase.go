package permission_usecase

import (
	"fmt"

	"github.com/google/uuid"
	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
	port_permission_repository "github.com/williamkoller/system-education/internal/permission/port/repository"
	port_permission_usecase "github.com/williamkoller/system-education/internal/permission/port/usecase"
	permission_dtos "github.com/williamkoller/system-education/internal/permission/presentation/dtos"
)

type PermissionUsecase struct {
	permissionRepository port_permission_repository.PermissionRepository
}

func NewPermissionUsecase(permissionRepository port_permission_repository.PermissionRepository) *PermissionUsecase {
	return &PermissionUsecase{
		permissionRepository: permissionRepository,
	}
}

var _ port_permission_usecase.PermissionUsecase = &PermissionUsecase{}

func (p *PermissionUsecase) Create(input permission_dtos.AddPermissionDto) (*permission_entity.Permission, error) {
	newPermission, err := permission_entity.NewPermission(&permission_entity.Permission{
		ID:          uuid.New().String(),
		UserID:      input.UserID,
		Modules:     input.Modules,
		Actions:     input.Actions,
		Level:       input.Level,
		Description: input.Description,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	permission, err := p.permissionRepository.Save(newPermission)
	if err != nil {
		return nil, fmt.Errorf("failed to save permission: %w", err)
	}

	return permission, nil
}

func (p *PermissionUsecase) FindAll() ([]*permission_entity.Permission, error) {
	permissions, err := p.permissionRepository.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find all permissions: %w", err)
	}
	return permissions, nil
}

func (p *PermissionUsecase) FindById(id string) (*permission_entity.Permission, error) {
	permission, err := p.permissionRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find permission by id: %w", err)
	}
	return permission, nil
}

func (p *PermissionUsecase) Update(id string, input permission_dtos.UpdatePermissionDto) (*permission_entity.Permission, error) {
	permission, err := p.permissionRepository.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to find permission by id: %w", err)
	}

	permission, err = permission.UpdatePermission(input.Modules, input.Actions, input.Level, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	permission, err = p.permissionRepository.Update(id, permission)
	if err != nil {
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}
	return permission, nil
}

func (p *PermissionUsecase) Delete(id string) error {
	user, err := p.permissionRepository.FindByID(id)
	if err != nil {
		return fmt.Errorf("failed to find permission by id: %w", err)
	}
	return p.permissionRepository.Delete(user.ID)
}

func (p *PermissionUsecase) FindPermissionByUserID(userID string) ([]*permission_entity.Permission, error) {
	permissions, err := p.permissionRepository.FindPermissionByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find permission by user id: %w", err)
	}
	return permissions, nil
}
