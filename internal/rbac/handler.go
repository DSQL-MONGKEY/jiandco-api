package rbac

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"jianco-api/internal/auth"
	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListRoles(c fiber.Ctx) error {
	roles, err := h.service.ListRoles(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, ToRoleResponses(roles))
}

func (h *Handler) ListPermissions(c fiber.Ctx) error {
	permissions, err := h.service.ListPermissions(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, ToPermissionResponses(permissions))
}

func (h *Handler) GetRolePermissions(c fiber.Ctx) error {
	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid role id")
	}

	role, permissions, err := h.service.RolePermissions(c.Context(), roleID)
	if err != nil {
		return err
	}

	return response.OK(c, map[string]any{
		"role":        ToRoleResponse(role),
		"permissions": ToPermissionResponses(permissions),
	})
}

func (h *Handler) UpdateRolePermissions(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	roleID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid role id")
	}

	var request UpdateRolePermissionsRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}

	role, permissions, err := h.service.UpdateRolePermissions(c.Context(), actor, roleID, request.PermissionIDs)
	if err != nil {
		return err
	}

	return response.OK(c, map[string]any{
		"role":        ToRoleResponse(role),
		"permissions": ToPermissionResponses(permissions),
	})
}
