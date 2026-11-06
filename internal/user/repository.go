package user

import (
	"context"
	"strings"
	"time"

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

func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.User, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.User{})

	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where(
			"name ILIKE ? OR username ILIKE ? OR email ILIKE ?",
			pattern, pattern, pattern,
		)
	}
	if filter.RoleID != nil {
		query = query.Where("role_id = ?", *filter.RoleID)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	err := query.
		Preload("Role").
		Order("created_at DESC").
		Offset((filter.Page - 1) * filter.PerPage).
		Limit(filter.PerPage).
		Find(&users).Error
	return users, total, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Preload("Role").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindRole(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	if err := r.db.WithContext(ctx).First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *Repository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(values).Error
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.User{}, "id = ?", id).Error; err != nil {
			return err
		}
		return tx.Model(&models.AuthSession{}).
			Where("user_id = ? AND revoked_at IS NULL", id).
			Update("revoked_at", now).Error
	})
}

func (r *Repository) UpdatePassword(ctx context.Context, id uuid.UUID, hash string, now time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", id).Update("password_hash", hash).Error; err != nil {
			return err
		}
		return tx.Model(&models.AuthSession{}).
			Where("user_id = ? AND revoked_at IS NULL", id).
			Update("revoked_at", now).Error
	})
}

func (r *Repository) RevokeSessions(ctx context.Context, id uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.AuthSession{}).
		Where("user_id = ? AND revoked_at IS NULL", id).
		Update("revoked_at", now).Error
}
