package customer

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"jianco-api/pkg/apperror"
	"jianco-api/pkg/response"
	"jianco-api/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(c fiber.Ctx) error {
	pagination := utils.ParsePagination(c)
	filter := ListFilter{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Search:  strings.TrimSpace(c.Query("search")),
	}

	customers, total, err := h.service.List(c.Context(), filter)
	if err != nil {
		return err
	}

	return response.OK(c, response.List{
		Items:      ToResponses(customers),
		Pagination: pagination.Meta(total),
	})
}

func (h *Handler) Create(c fiber.Ctx) error {
	var request Request
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	customer, err := h.service.Create(c.Context(), request)
	if err != nil {
		return err
	}
	return response.Created(c, ToResponse(customer))
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid customer id")
	}
	customer, err := h.service.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(customer))
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid customer id")
	}

	var request Request
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	customer, err := h.service.Update(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToResponse(customer))
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return apperror.BadRequest("invalid customer id")
	}
	if err := h.service.Delete(c.Context(), id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "customer deleted"})
}
