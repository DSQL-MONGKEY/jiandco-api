package models

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_roles_name"`
	Level       int       `gorm:"type:smallint;not null;uniqueIndex:idx_roles_level;check:chk_roles_level,level >= 0 AND level <= 3"`
	Description string    `gorm:"type:varchar(255)"`
	IsSystem    bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

type Permission struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Key       string    `gorm:"column:key;type:varchar(100);not null;uniqueIndex:idx_permissions_key"`
	Page      string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_permissions_page_action"`
	Action    string    `gorm:"type:varchar(20);not null;uniqueIndex:idx_permissions_page_action"`
	Label     string    `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

type RolePermission struct {
	RoleID       uuid.UUID  `gorm:"type:uuid;primaryKey"`
	PermissionID uuid.UUID  `gorm:"type:uuid;primaryKey;index:idx_role_permissions_permission_id"`
	Role         Role       `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Permission   Permission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt    time.Time  `gorm:"type:timestamptz;not null;default:now()"`
}
