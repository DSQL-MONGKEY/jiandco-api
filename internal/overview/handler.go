package overview

import (
	"github.com/gofiber/fiber/v3"

	"jianco-api/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Revenue(c fiber.Ctx) error {
	result, err := h.service.Revenue(c.Context(), c.Query("month"))
	if err != nil {
		return err
	}
	return response.OK(c, result)
}
