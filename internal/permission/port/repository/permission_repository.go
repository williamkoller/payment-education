package port_permission_repository

import permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"

type PermissionRepository interface {
	Create(p *permission_entity.Permission) (*permission_entity.Permission, error)
	FindAll() ([]*permission_entity.Permission, error)
	FindPermissionByUserID(id string) (*permission_entity.Permission, error)
	Update(p *permission_entity.Permission) (*permission_entity.Permission, error)
	Delete(id string) error
}
