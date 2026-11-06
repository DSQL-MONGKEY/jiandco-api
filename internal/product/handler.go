package product

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

func (h *Handler) ListCategories(c fiber.Ctx) error {
	categories, err := h.service.ListCategories(c.Context(), c.Query("search"))
	if err != nil {
		return err
	}
	return response.OK(c, ToCategoryResponses(categories))
}

func (h *Handler) CreateCategory(c fiber.Ctx) error {
	var request CategoryRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	category, err := h.service.CreateCategory(c.Context(), request)
	if err != nil {
		return err
	}
	return response.Created(c, ToCategoryResponse(category))
}

func (h *Handler) GetCategory(c fiber.Ctx) error {
	id, err := parseID(c, "id", "category")
	if err != nil {
		return err
	}
	category, err := h.service.GetCategory(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, ToCategoryResponse(category))
}

func (h *Handler) UpdateCategory(c fiber.Ctx) error {
	id, err := parseID(c, "id", "category")
	if err != nil {
		return err
	}

	var request CategoryRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	category, err := h.service.UpdateCategory(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToCategoryResponse(category))
}

func (h *Handler) DeleteCategory(c fiber.Ctx) error {
	id, err := parseID(c, "id", "category")
	if err != nil {
		return err
	}
	if err := h.service.DeleteCategory(c.Context(), id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "category deleted"})
}

func (h *Handler) ListSizes(c fiber.Ctx) error {
	sizes, err := h.service.ListSizes(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, ToSizeResponses(sizes))
}

func (h *Handler) CreateSize(c fiber.Ctx) error {
	var request SizeRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	size, err := h.service.CreateSize(c.Context(), request)
	if err != nil {
		return err
	}
	return response.Created(c, ToSizeResponse(size))
}

func (h *Handler) GetSize(c fiber.Ctx) error {
	id, err := parseID(c, "id", "size")
	if err != nil {
		return err
	}
	size, err := h.service.GetSize(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, ToSizeResponse(size))
}

func (h *Handler) UpdateSize(c fiber.Ctx) error {
	id, err := parseID(c, "id", "size")
	if err != nil {
		return err
	}

	var request SizeRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	size, err := h.service.UpdateSize(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToSizeResponse(size))
}

func (h *Handler) DeleteSize(c fiber.Ctx) error {
	id, err := parseID(c, "id", "size")
	if err != nil {
		return err
	}
	if err := h.service.DeleteSize(c.Context(), id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "size deleted"})
}

func (h *Handler) ListColors(c fiber.Ctx) error {
	colors, err := h.service.ListColors(c.Context())
	if err != nil {
		return err
	}
	return response.OK(c, ToColorResponses(colors))
}

func (h *Handler) CreateColor(c fiber.Ctx) error {
	var request ColorRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	color, err := h.service.CreateColor(c.Context(), request)
	if err != nil {
		return err
	}
	return response.Created(c, ToColorResponse(color))
}

func (h *Handler) GetColor(c fiber.Ctx) error {
	id, err := parseID(c, "id", "color")
	if err != nil {
		return err
	}
	color, err := h.service.GetColor(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, ToColorResponse(color))
}

func (h *Handler) UpdateColor(c fiber.Ctx) error {
	id, err := parseID(c, "id", "color")
	if err != nil {
		return err
	}

	var request ColorRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	color, err := h.service.UpdateColor(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToColorResponse(color))
}

func (h *Handler) DeleteColor(c fiber.Ctx) error {
	id, err := parseID(c, "id", "color")
	if err != nil {
		return err
	}
	if err := h.service.DeleteColor(c.Context(), id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "color deleted"})
}

func (h *Handler) ListProducts(c fiber.Ctx) error {
	pagination := utils.ParsePagination(c)
	filter := ProductListFilter{
		Page:    pagination.Page,
		PerPage: pagination.PerPage,
		Search:  strings.TrimSpace(c.Query("search")),
		Status:  strings.TrimSpace(c.Query("status")),
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		parsed, err := uuid.Parse(categoryID)
		if err != nil {
			return apperror.BadRequest("invalid category_id")
		}
		filter.CategoryID = &parsed
	}
	if filter.Status != "" && filter.Status != models.ProductStatusDraft && filter.Status != models.ProductStatusActive && filter.Status != models.ProductStatusArchived {
		return apperror.BadRequest("invalid status")
	}

	products, total, err := h.service.ListProducts(c.Context(), filter)
	if err != nil {
		return err
	}

	items := make([]ProductResponse, 0, len(products))
	for index := range products {
		items = append(items, ToProductResponse(&products[index], h.storage))
	}

	return response.OK(c, response.List{
		Items:      items,
		Pagination: pagination.Meta(total),
	})
}

func (h *Handler) CreateProduct(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	var request CreateProductRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	product, err := h.service.CreateProduct(c.Context(), actor, request)
	if err != nil {
		return err
	}
	return response.Created(c, ToProductResponse(product, h.storage))
}

func (h *Handler) GetProduct(c fiber.Ctx) error {
	id, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	product, err := h.service.GetProduct(c.Context(), id)
	if err != nil {
		return err
	}
	return response.OK(c, ToProductResponse(product, h.storage))
}

func (h *Handler) UpdateProduct(c fiber.Ctx) error {
	id, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}

	var request UpdateProductRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	product, err := h.service.UpdateProduct(c.Context(), id, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToProductResponse(product, h.storage))
}

func (h *Handler) DeleteProduct(c fiber.Ctx) error {
	id, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	if err := h.service.DeleteProduct(c.Context(), id); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "product deleted"})
}

func (h *Handler) CreateVariant(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}

	var request CreateVariantRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	variant, err := h.service.CreateVariant(c.Context(), actor, productID, request)
	if err != nil {
		return err
	}
	return response.Created(c, ToVariantResponse(variant))
}

func (h *Handler) UpdateVariant(c fiber.Ctx) error {
	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	variantID, err := parseID(c, "variantId", "variant")
	if err != nil {
		return err
	}

	var request UpdateVariantRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	variant, err := h.service.UpdateVariant(c.Context(), productID, variantID, request)
	if err != nil {
		return err
	}
	return response.OK(c, ToVariantResponse(variant))
}

func (h *Handler) DeleteVariant(c fiber.Ctx) error {
	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	variantID, err := parseID(c, "variantId", "variant")
	if err != nil {
		return err
	}
	if err := h.service.DeleteVariant(c.Context(), productID, variantID); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "variant deleted"})
}

func (h *Handler) UpdateStock(c fiber.Ctx) error {
	actor, ok := auth.CurrentUser(c)
	if !ok {
		return apperror.Unauthorized("authentication is required")
	}

	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	variantID, err := parseID(c, "variantId", "variant")
	if err != nil {
		return err
	}

	var request UpdateStockRequest
	if err := c.Bind().Body(&request); err != nil {
		return apperror.BadRequest("invalid request body")
	}
	if err := utils.Validate(&request); err != nil {
		return err
	}

	variant, err := h.service.UpdateStock(c.Context(), actor, productID, variantID, *request.Stock, request.Note)
	if err != nil {
		return err
	}
	return response.OK(c, ToVariantResponse(variant))
}

func (h *Handler) ListStockMovements(c fiber.Ctx) error {
	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	variantID, err := parseID(c, "variantId", "variant")
	if err != nil {
		return err
	}

	pagination := utils.ParsePagination(c)
	movements, total, err := h.service.ListStockMovements(c.Context(), productID, variantID, pagination.Page, pagination.PerPage)
	if err != nil {
		return err
	}

	return response.OK(c, response.List{
		Items:      ToStockMovementResponses(movements),
		Pagination: pagination.Meta(total),
	})
}

func (h *Handler) UploadImages(c fiber.Ctx) error {
	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}

	form, err := c.MultipartForm()
	if err != nil {
		return apperror.BadRequest("invalid multipart form")
	}
	files := form.File["images"]
	if len(files) == 0 {
		files = form.File["file"]
	}

	images, err := h.service.UploadImages(c.Context(), productID, files)
	if err != nil {
		return err
	}

	items := make([]ImageResponse, 0, len(images))
	for index := range images {
		items = append(items, ToImageResponse(&images[index], h.storage))
	}
	return response.Created(c, items)
}

func (h *Handler) DeleteImage(c fiber.Ctx) error {
	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	imageID, err := parseID(c, "imageId", "image")
	if err != nil {
		return err
	}
	if err := h.service.DeleteImage(c.Context(), productID, imageID); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "image deleted"})
}

func (h *Handler) SetPrimaryImage(c fiber.Ctx) error {
	productID, err := parseID(c, "id", "product")
	if err != nil {
		return err
	}
	imageID, err := parseID(c, "imageId", "image")
	if err != nil {
		return err
	}
	if err := h.service.SetPrimaryImage(c.Context(), productID, imageID); err != nil {
		return err
	}
	return response.OK(c, map[string]any{"message": "primary image updated"})
}

func parseID(c fiber.Ctx, param, label string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(param))
	if err != nil {
		return uuid.Nil, apperror.BadRequest("invalid " + label + " id")
	}
	return id, nil
}
