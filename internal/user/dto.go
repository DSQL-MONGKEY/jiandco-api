package user

import (
	"github.com/google/uuid"
)

type CreateRequest struct {
	Name     string    `json:"name" validate:"required,min=2,max=100"`
	Username string    `json:"username" validate:"required,min=3,max=50"`
	Email    string    `json:"email" validate:"required,email,max=255"`
	Password string    `json:"password" validate:"required,min=8,max=72"`
	Phone    string    `json:"phone" validate:"max=30"`
	RoleID   uuid.UUID `json:"role_id" validate:"required"`
	IsActive *bool     `json:"is_active"`
}

type UpdateRequest struct {
	Name     string    `json:"name" validate:"required,min=2,max=100"`
	Email    string    `json:"email" validate:"required,email,max=255"`
	Phone    string    `json:"phone" validate:"max=30"`
	RoleID   uuid.UUID `json:"role_id" validate:"required"`
	IsActive *bool     `json:"is_active"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=8,max=72"`
}

type ListFilter struct {
	Page     int
	PerPage  int
	Search   string
	RoleID   *uuid.UUID
	IsActive *bool
}
