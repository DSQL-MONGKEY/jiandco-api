package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/auth"
	"jianco-api/internal/models"
	"jianco-api/pkg/apperror"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.User, int64, error) {
	users, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return users, total, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("user not found")
		}
		return nil, apperror.Internal(err)
	}
	return user, nil
}

func (s *Service) Create(ctx context.Context, actor *models.User, request CreateRequest) (*models.User, error) {
	role, err := s.repo.FindRole(ctx, request.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.BadRequest("role not found")
		}
		return nil, apperror.Internal(err)
	}
	if !canManage(actor, role.Level) {
		return nil, apperror.Forbidden("you cannot assign this role")
	}

	hash, err := auth.HashPassword(request.Password)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	isActive := true
	if request.IsActive != nil {
		isActive = *request.IsActive
	}

	newUser := &models.User{
		RoleID:       role.ID,
		Name:         strings.TrimSpace(request.Name),
		Username:     strings.ToLower(strings.TrimSpace(request.Username)),
		Email:        strings.ToLower(strings.TrimSpace(request.Email)),
		PasswordHash: hash,
		Phone:        strings.TrimSpace(request.Phone),
		IsActive:     isActive,
	}
	if err := s.repo.Create(ctx, newUser); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("username or email is already used")
		}
		return nil, apperror.Internal(err)
	}

	created, err := s.repo.FindByID(ctx, newUser.ID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return created, nil
}

func (s *Service) Update(ctx context.Context, actor *models.User, id uuid.UUID, request UpdateRequest) (*models.User, error) {
	target, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("user not found")
		}
		return nil, apperror.Internal(err)
	}
	if !canManage(actor, target.Role.Level) {
		return nil, apperror.Forbidden("you cannot manage this user")
	}

	role, err := s.repo.FindRole(ctx, request.RoleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.BadRequest("role not found")
		}
		return nil, apperror.Internal(err)
	}
	if !canManage(actor, role.Level) {
		return nil, apperror.Forbidden("you cannot assign this role")
	}

	isActive := target.IsActive
	if request.IsActive != nil {
		isActive = *request.IsActive
	}

	if target.ID == actor.ID {
		if role.ID != target.RoleID {
			return nil, apperror.BadRequest("you cannot change your own role")
		}
		if !isActive {
			return nil, apperror.BadRequest("you cannot deactivate your own account")
		}
	}

	values := map[string]any{
		"name":      strings.TrimSpace(request.Name),
		"email":     strings.ToLower(strings.TrimSpace(request.Email)),
		"phone":     strings.TrimSpace(request.Phone),
		"role_id":   role.ID,
		"is_active": isActive,
	}
	if err := s.repo.Update(ctx, id, values); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("email is already used")
		}
		return nil, apperror.Internal(err)
	}

	if !isActive {
		if err := s.repo.RevokeSessions(ctx, id, time.Now()); err != nil {
			return nil, apperror.Internal(err)
		}
	}

	updated, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return updated, nil
}

func (s *Service) Delete(ctx context.Context, actor *models.User, id uuid.UUID) error {
	target, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("user not found")
		}
		return apperror.Internal(err)
	}
	if target.ID == actor.ID {
		return apperror.BadRequest("you cannot delete your own account")
	}
	if !canManage(actor, target.Role.Level) {
		return apperror.Forbidden("you cannot manage this user")
	}

	if err := s.repo.Delete(ctx, id, time.Now()); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) ResetPassword(ctx context.Context, actor *models.User, id uuid.UUID, newPassword string) error {
	target, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("user not found")
		}
		return apperror.Internal(err)
	}
	if !canManage(actor, target.Role.Level) {
		return apperror.Forbidden("you cannot manage this user")
	}

	hash, err := auth.HashPassword(newPassword)
	if err != nil {
		return apperror.Internal(err)
	}
	if err := s.repo.UpdatePassword(ctx, id, hash, time.Now()); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func canManage(actor *models.User, targetLevel int) bool {
	if actor == nil {
		return false
	}
	if actor.Role.Level == models.RoleLevelDeveloper {
		return true
	}
	return targetLevel > actor.Role.Level
}
