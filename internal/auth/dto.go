package auth

import (
	"time"

	"github.com/google/uuid"

	"jianco-api/internal/models"
	"jianco-api/internal/storage"
)

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email,max=255"`
	Phone string `json:"phone" validate:"max=30"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=72"`
}

type RoleResponse struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Level int       `json:"level"`
}

type UserResponse struct {
	ID          uuid.UUID    `json:"id"`
	RoleID      uuid.UUID    `json:"role_id"`
	Role        RoleResponse `json:"role"`
	Name        string       `json:"name"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Phone       string       `json:"phone"`
	AvatarKey   string       `json:"avatar_key"`
	AvatarURL   string       `json:"avatar_url"`
	IsActive    bool         `json:"is_active"`
	LastLoginAt *time.Time   `json:"last_login_at"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type MeResponse struct {
	User        UserResponse `json:"user"`
	Permissions []string     `json:"permissions"`
}

type AuthResponse struct {
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
	ExpiresIn    int          `json:"expiresIn"`
	User         UserResponse `json:"user"`
}

func ToUserResponse(user *models.User, storageClient *storage.Client) UserResponse {
	if user == nil {
		return UserResponse{}
	}

	return UserResponse{
		ID:     user.ID,
		RoleID: user.RoleID,
		Role: RoleResponse{
			ID:    user.Role.ID,
			Name:  user.Role.Name,
			Level: user.Role.Level,
		},
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		Phone:       user.Phone,
		AvatarKey:   user.AvatarKey,
		AvatarURL:   storageClient.URL(user.AvatarKey),
		IsActive:    user.IsActive,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
