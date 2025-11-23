package permission_repository

import (
	permission_entity "github.com/williamkoller/system-education/internal/permission/domain/entity"
	permission_model "github.com/williamkoller/system-education/internal/permission/infra/db/model"
	port_permission_repository "github.com/williamkoller/system-education/internal/permission/port/repository"
	"gorm.io/gorm"
)

type PermissionGormRepository struct {
	DB *gorm.DB
}

func NewPermissionGormRepository(db *gorm.DB) *PermissionGormRepository {
	return &PermissionGormRepository{DB: db}
}

var _ port_permission_repository.PermissionRepository = &PermissionGormRepository{}

func (r *PermissionGormRepository) Create(p *permission_entity.Permission) (*permission_entity.Permission, error) {
	model := permission_model.FromEntity(p)
	if err := r.DB.Create(&model).Error; err != nil {
		return nil, err
	}
	return permission_model.ToEntity(model), nil
}

func (r *PermissionGormRepository) FindByID(id string) (*permission_entity.Permission, error) {
	var permission *permission_entity.Permission
	model := permission_model.FromEntity(permission)
	if err := r.DB.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return permission_model.ToEntity(model), nil
}

func (r *PermissionGormRepository) FindAll() ([]*permission_entity.Permission, error) {
	var permissions []*permission_entity.Permission
	model := permission_model.FromEntities(permissions)
	if err := r.DB.Find(&model).Error; err != nil {
		return nil, err
	}
	return permission_model.ToEntities(model), nil
}

func (r *PermissionGormRepository) Update(id string, p *permission_entity.Permission) (*permission_entity.Permission, error) {
	model := permission_model.FromEntity(p)
	if err := r.DB.Model(&permission_model.Permission{}).
		Where("id = ?", id).
		Updates(&model).Error; err != nil {
		return nil, err
	}
	return permission_model.ToEntity(model), nil
}

func (r *PermissionGormRepository) Delete(id string) error {
	return r.DB.Unscoped().Delete(&permission_model.Permission{}, "id = ?", id).Error
}

func (r *PermissionGormRepository) FindPermissionByUserID(userID string) ([]*permission_entity.Permission, error) {
	var permissions []*permission_entity.Permission
	model := permission_model.FromEntities(permissions)
	if err := r.DB.Where("user_id = ?", userID).Find(&model).Error; err != nil {
		return nil, err
	}
	return permission_model.ToEntities(model), nil
}