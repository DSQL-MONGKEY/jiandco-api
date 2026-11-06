package rbac

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.WithContext(ctx).Order("level ASC").Find(&roles).Error
	return roles, err
}

func (r *Repository) FindRole(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.WithContext(ctx).
		Order("page ASC, action ASC").
		Find(&permissions).Error
	return permissions, err
}

func (r *Repository) RolePermissions(ctx context.Context, roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.WithContext(ctx).
		Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.page ASC, permissions.action ASC").
		Find(&permissions).Error
	return permissions, err
}

func (r *Repository) HasPermission(ctx context.Context, roleID uuid.UUID, page, action string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.RolePermission{}).
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ? AND permissions.page = ? AND permissions.action = ?", roleID, page, action).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CountPermissions(ctx context.Context, ids []uuid.UUID) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Permission{}).
		Where("id IN ?", ids).
		Count(&count).Error
	return count, err
}

func (r *Repository) ReplaceRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}
		if len(permissionIDs) == 0 {
			return nil
		}

		rows := make([]models.RolePermission, 0, len(permissionIDs))
		for _, permissionID := range permissionIDs {
			rows = append(rows, models.RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
			})
		}
		return tx.Create(&rows).Error
	})
}

func (r *Repository) PermissionKeys(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	var keys []string
	err := r.db.WithContext(ctx).
		Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Pluck("permissions.key", &keys).Error
	return keys, err
}
