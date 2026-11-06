package rbac

import (
	"github.com/google/uuid"

	"jianco-api/internal/models"
)

type RoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Level       int       `json:"level"`
	Description string    `json:"description"`
	IsSystem    bool      `json:"is_system"`
}

type PermissionResponse struct {
	ID     uuid.UUID `json:"id"`
	Key    string    `json:"key"`
	Page   string    `json:"page"`
	Action string    `json:"action"`
	Label  string    `json:"label"`
}

type UpdateRolePermissionsRequest struct {
	PermissionIDs []uuid.UUID `json:"permission_ids"`
}

func ToRoleResponse(role *models.Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Level:       role.Level,
		Description: role.Description,
		IsSystem:    role.IsSystem,
	}
}

func ToRoleResponses(roles []models.Role) []RoleResponse {
	result := make([]RoleResponse, 0, len(roles))
	for index := range roles {
		result = append(result, ToRoleResponse(&roles[index]))
	}
	return result
}

func ToPermissionResponse(permission *models.Permission) PermissionResponse {
	return PermissionResponse{
		ID:     permission.ID,
		Key:    permission.Key,
		Page:   permission.Page,
		Action: permission.Action,
		Label:  permission.Label,
	}
}

func ToPermissionResponses(permissions []models.Permission) []PermissionResponse {
	result := make([]PermissionResponse, 0, len(permissions))
	for index := range permissions {
		result = append(result, ToPermissionResponse(&permissions[index]))
	}
	return result
}
