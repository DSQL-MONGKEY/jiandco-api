package user

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"jianco-api/internal/auth"
	"jianco-api/internal/storage"
	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
	"jianco-api/pkg/utils"
)

type Handler struct {
	service *Service
	storage *storage.Client
}

func NewHandler(service *Service, storageClient *storage.Client) *Handler {
	return &Handler{service: service, storage: storageClient}
}

func (h *Handler) List(c fiber.Ctx) error {
	pagination := utils.ParsePagination(c)

	filter := ListFilter{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Search:  strings.TrimSpace(c.Query("search")),
	}
	if roleID := c.Query("role_id"); roleID != "" {
		parsed, err := uuid.Parse(roleID)
		if err != nil {
			return apperror.BadRequest("invalid role_id")
		}
		filter.RoleID = &parsed
	}
	if isActive := c.Query("is_active"); isActive != "" {
		parsed, err := strconv.ParseBool(isActive)
		if err != nil {
			return apperror.BadRequest("invalid is_active")
		}
		filter.IsActive = &parsed
	}

	users, total, err := h.service.List(c.Context(), filter)
	if err != nil {
		return err
	}

	items := make([]auth.UserResponse, 0, len(users))
	for index := range users {
		items = append(items, auth.ToUserResponse(&users[index], h.storage))
	}

	return response.OK(c, response.List{
		Items:      items,
		Pagination: pagination.Meta(total),
	})
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid user id")
	}

	user, err := h.service.Get(c.Context(), id)
	if err != nil {
		return err
	}

	return response.OK(c, auth.ToUserResponse(user, h.storage))
}

func (h *Handler) Create(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	var request CreateRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	created, err := h.service.Create(c.Context(), actor, request)
	if err != nil {
		return err
	}

	return response.Created(c, auth.ToUserResponse(created, h.storage))
}

func (h *Handler) Update(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid user id")
	}

	var request UpdateRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	updated, err := h.service.Update(c.Context(), actor, id, request)
	if err != nil {
		return err
	}

	return response.OK(c, auth.ToUserResponse(updated, h.storage))
}

func (h *Handler) Delete(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid user id")
	}

	if err := h.service.Delete(c.Context(), actor, id); err != nil {
		return err
	}

	return response.OK(c, map[string]any{"message": "user deleted"})
}

func (h *Handler) ResetPassword(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid user id")
	}

	var request ResetPasswordRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	if err := h.service.ResetPassword(c.Context(), actor, id, request.NewPassword); err != nil {
		return err
	}

	return response.OK(c, map[string]any{"message": "password reset"})
}
