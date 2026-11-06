package auth

import (
	"context"
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

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) CreateSession(ctx context.Context, session *models.AuthSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r *Repository) FindSessionByHash(ctx context.Context, hash string) (*models.AuthSession, error) {
	var session models.AuthSession
	err := r.db.WithContext(ctx).
		Preload("User.Role").
		Where("refresh_token_hash = ?", hash).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) FindSessionByID(ctx context.Context, id uuid.UUID) (*models.AuthSession, error) {
	var session models.AuthSession
	err := r.db.WithContext(ctx).
		Preload("User.Role").
		First(&session, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *Repository) RotateSession(ctx context.Context, sessionID uuid.UUID, refreshTokenHash string, usedAt time.Time) (bool, error) {
	result := r.db.WithContext(ctx).
		Model(&models.AuthSession{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Updates(map[string]any{
			"refresh_token_hash": refreshTokenHash,
			"last_used_at":       usedAt,
		})
	return result.RowsAffected > 0, result.Error
}

func (r *Repository) RevokeSession(ctx context.Context, sessionID uuid.UUID, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.AuthSession{}).
		Where("id = ? AND revoked_at IS NULL", sessionID).
		Update("revoked_at", revokedAt).Error
}

func (r *Repository) RevokeAllSessions(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.AuthSession{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", revokedAt).Error
}

func (r *Repository) RevokeAllSessionsExcept(ctx context.Context, userID, sessionID uuid.UUID, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.AuthSession{}).
		Where("user_id = ? AND id <> ? AND revoked_at IS NULL", userID, sessionID).
		Update("revoked_at", revokedAt).Error
}

func (r *Repository) UpdateLastLogin(ctx context.Context, userID uuid.UUID, at time.Time) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("last_login_at", at).Error
}

func (r *Repository) UpdateProfile(ctx context.Context, userID uuid.UUID, name, email, phone string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"name":  name,
			"email": email,
			"phone": phone,
		}).Error
}

func (r *Repository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("password_hash", passwordHash).Error
}

func (r *Repository) UpdateAvatar(ctx context.Context, userID uuid.UUID, objectKey string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("avatar_key", objectKey).Error
}

func (r *Repository) PermissionKeys(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	keys := make([]string, 0)
	err := r.db.WithContext(ctx).
		Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleID).
		Order("permissions.key ASC").
		Pluck("permissions.key", &keys).Error
	return keys, err
}

func (r *Repository) AllPermissionKeys(ctx context.Context) ([]string, error) {
	keys := make([]string, 0)
	err := r.db.WithContext(ctx).
		Model(&models.Permission{}).
		Order("key ASC").
		Pluck("key", &keys).Error
	return keys, err
}
