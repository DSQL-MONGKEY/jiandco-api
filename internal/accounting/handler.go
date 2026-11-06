package accounting

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"jianco-api/internal/auth"
	"jianco-api/internal/models"
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
		Page:     pagination.Page,
		PerPage:  pagination.PerPage,
		Search:   strings.TrimSpace(c.Query("search")),
		Category: strings.TrimSpace(c.Query("category")),
		Month:    strings.TrimSpace(c.Query("month")),
	}
	if filter.Category != "" && !validCategory(filter.Category) {
		return apperror.BadRequest("invalid category")
	}

	expenses, total, err := h.service.List(c.Context(), filter)
	if err != nil {
		return err
	}

	return response.OK(c, response.List{
		Items:      ToResponses(expenses, h.storage),
		Pagination: pagination.Meta(total),
	})
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

	expense, err := h.service.Create(c.Context(), actor, request)
	if err != nil {
		return err
	}
	return response.Created(c, ToResponse(expense, h.storage))
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid expense id")
	}
	expense, err := h.service.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(expense, h.storage))
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid expense id")
	}

	var request UpdateRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	expense, err := h.service.Update(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(expense, h.storage))
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid expense id")
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "expense deleted"})
}

func (h *Handler) UploadReceipt(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid expense id")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return apperror.BadRequest("receipt file is required")
	}

	result, err := h.service.UploadReceipt(c.Context(), id, file)
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

func (h *Handler) DeleteReceipt(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid expense id")
	}
	if err := h.service.DeleteReceipt(c.Context(), id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "receipt deleted"})
}

func (h *Handler) Overview(c fiber.Ctx) error {
	overview, err := h.service.Overview(c.Context(), c.Query("month"))
	if err != nil {
		return err
	}
	return response.OK(c, overview)
}

func validCategory(value string) bool {
	switch value {
	case models.ExpenseCategoryOperational, models.ExpenseCategoryProduction, models.ExpenseCategoryMarketing,
		models.ExpenseCategoryRestock, models.ExpenseCategoryOther:
		return true
	default:
		return false
	}
}
