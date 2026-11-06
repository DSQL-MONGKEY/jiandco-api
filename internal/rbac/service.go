package rbac

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/models"
	"jianco-api/pkg/apperror"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListRoles(ctx context.Context) ([]models.Role, error) {
	roles, err := s.repo.ListRoles(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return roles, nil
}

func (s *Service) ListPermissions(ctx context.Context) ([]models.Permission, error) {
	permissions, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return permissions, nil
}

func (s *Service) RolePermissions(ctx context.Context, roleID uuid.UUID) (*models.Role, []models.Permission, error) {
	role, err := s.repo.FindRole(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apperror.NotFound("role not found")
		}
		return nil, nil, apperror.Internal(err)
	}

	permissions, err := s.repo.RolePermissions(ctx, roleID)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}

	return role, permissions, nil
}

func (s *Service) UpdateRolePermissions(ctx context.Context, actor *models.User, roleID uuid.UUID, permissionIDs []uuid.UUID) (*models.Role, []models.Permission, error) {
	if actor == nil {
		return nil, nil, apperror.Unauthorized("authentication is required")
	}
	if actor.Role.Level > models.RoleLevelSuperAdmin {
		return nil, nil, apperror.Forbidden("only developer and super admin can change role permissions")
	}

	role, err := s.repo.FindRole(ctx, roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apperror.NotFound("role not found")
		}
		return nil, nil, apperror.Internal(err)
	}

	if role.Level <= models.RoleLevelSuperAdmin {
		return nil, nil, apperror.BadRequest("developer and super admin roles always have full access")
	}

	uniqueIDs := make([]uuid.UUID, 0, len(permissionIDs))
	seen := make(map[uuid.UUID]struct{}, len(permissionIDs))
	for _, permissionID := range permissionIDs {
		if _, exists := seen[permissionID]; exists {
			continue
		}
		seen[permissionID] = struct{}{}
		uniqueIDs = append(uniqueIDs, permissionID)
	}

	count, err := s.repo.CountPermissions(ctx, uniqueIDs)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}
	if count != int64(len(uniqueIDs)) {
		return nil, nil, apperror.BadRequest("one or more permissions do not exist")
	}

	if err := s.repo.ReplaceRolePermissions(ctx, roleID, uniqueIDs); err != nil {
		return nil, nil, apperror.Internal(err)
	}

	permissions, err := s.repo.RolePermissions(ctx, roleID)
	if err != nil {
		return nil, nil, apperror.Internal(err)
	}

	return role, permissions, nil
}
