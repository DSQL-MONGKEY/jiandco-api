package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Base
	RoleID       uuid.UUID  `gorm:"type:uuid;not null;index:idx_users_role_id"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Name         string     `gorm:"type:varchar(100);not null"`
	Username     string     `gorm:"type:varchar(50);not null;uniqueIndex:idx_users_username"`
	Email        string     `gorm:"type:varchar(255);not null;uniqueIndex:idx_users_email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	Phone        string     `gorm:"type:varchar(30)"`
	AvatarKey    string     `gorm:"type:varchar(255)"`
	IsActive     bool       `gorm:"not null;default:true;index:idx_users_is_active"`
	LastLoginAt  *time.Time `gorm:"type:timestamptz"`
	SoftDelete
}

type AuthSession struct {
	ID               uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index:idx_auth_sessions_user_id"`
	User             User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	RefreshTokenHash string     `gorm:"type:char(64);not null;uniqueIndex:idx_auth_sessions_refresh_token_hash"`
	UserAgent        string     `gorm:"type:varchar(255)"`
	IPAddress        string     `gorm:"type:varchar(45)"`
	ExpiresAt        time.Time  `gorm:"type:timestamptz;not null;index:idx_auth_sessions_expires_at"`
	RevokedAt        *time.Time `gorm:"type:timestamptz"`
	LastUsedAt       *time.Time `gorm:"type:timestamptz"`
	CreatedAt        time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}
