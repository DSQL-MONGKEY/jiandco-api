package catalog

import (
	"strings"

	"github.com/gofiber/fiber/v3"

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

func (h *Handler) PublicCatalog(c fiber.Ctx) error {
	result, err := h.service.PublicCatalog(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

func (h *Handler) PublicProduct(c fiber.Ctx) error {
	result, err := h.service.PublicProduct(c.Context(), c.Params("slug"))
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

func (h *Handler) PublicMedia(c fiber.Ctx) error {
	key := strings.TrimPrefix(c.Path(), "/api/v1/public/media/")
	if key == "" {
		return apperror.NotFound("media not found")
	}
	if !allowedMediaKey(key) {
		return apperror.NotFound("media not found")
	}

	data, contentType, err := h.storage.GetObject(c.Context(), key)
	if err != nil {
		return apperror.NotFound("media not found")
	}

	c.Set(fiber.HeaderContentType, contentType)
	c.Set(fiber.HeaderCacheControl, "public, max-age=86400")
	return c.Send(data)
}

func (h *Handler) GetSettings(c fiber.Ctx) error {
	result, err := h.service.Settings(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

func (h *Handler) UpdateSettings(c fiber.Ctx) error {
	var request SettingsRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	result, err := h.service.UpdateSettings(c.Context(), request)
	if err != nil {
		return err
	}
	return response.OK(c, result)
}

func allowedMediaKey(key string) bool {
	key = strings.TrimLeft(key, "/")
	return strings.HasPrefix(key, "catalog/") ||
		strings.HasPrefix(key, "user/") ||
		strings.HasPrefix(key, "accounting/")
}
