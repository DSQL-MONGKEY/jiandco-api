package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/models"
	"jianco-api/internal/storage"
	"jianco-api/pkg/apperror"
	"jianco-api/pkg/database"
	"jianco-api/pkg/utils"
)

type Service struct {
	repo    *Repository
	storage *storage.Client
}

func NewService(repo *Repository, storageClient *storage.Client) *Service {
	return &Service{repo: repo, storage: storageClient}
}

func (s *Service) ListCategories(ctx context.Context, search string) ([]models.Category, error) {
	categories, err := s.repo.ListCategories(ctx, search)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return categories, nil
}

func (s *Service) GetCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	category, err := s.repo.FindCategory(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("category not found")
		}
		return nil, apperror.Internal(err)
	}
	return category, nil
}

func (s *Service) CreateCategory(ctx context.Context, request CategoryRequest) (*models.Category, error) {
	slug, err := s.uniqueCategorySlug(ctx, request.Name, request.Slug, nil)
	if err != nil {
		return nil, err
	}

	category := &models.Category{
		Name:        strings.TrimSpace(request.Name),
		Slug:        slug,
		Description: strings.TrimSpace(request.Description),
	}
	if err := s.repo.CreateCategory(ctx, category); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("category name or slug is already used")
		}
		return nil, apperror.Internal(err)
	}
	return category, nil
}

func (s *Service) UpdateCategory(ctx context.Context, id uuid.UUID, request CategoryRequest) (*models.Category, error) {
	category, err := s.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	slug, err := s.uniqueCategorySlug(ctx, request.Name, request.Slug, &category.ID)
	if err != nil {
		return nil, err
	}

	err = s.repo.UpdateCategory(ctx, id, map[string]any{
		"name":        strings.TrimSpace(request.Name),
		"slug":        slug,
		"description": strings.TrimSpace(request.Description),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("category name or slug is already used")
		}
		return nil, apperror.Internal(err)
	}
	return s.GetCategory(ctx, id)
}

func (s *Service) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetCategory(ctx, id); err != nil {
		return err
	}
	if err := s.repo.DeleteCategory(ctx, id); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) ListSizes(ctx context.Context) ([]models.Size, error) {
	sizes, err := s.repo.ListSizes(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return sizes, nil
}

func (s *Service) GetSize(ctx context.Context, id uuid.UUID) (*models.Size, error) {
	size, err := s.repo.FindSize(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("size not found")
		}
		return nil, apperror.Internal(err)
	}
	return size, nil
}

func (s *Service) CreateSize(ctx context.Context, request SizeRequest) (*models.Size, error) {
	size := &models.Size{
		Name:      strings.TrimSpace(request.Name),
		SortOrder: request.SortOrder,
	}
	if err := s.repo.CreateSize(ctx, size); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("size name is already used")
		}
		return nil, apperror.Internal(err)
	}
	return size, nil
}

func (s *Service) UpdateSize(ctx context.Context, id uuid.UUID, request SizeRequest) (*models.Size, error) {
	if _, err := s.GetSize(ctx, id); err != nil {
		return nil, err
	}
	err := s.repo.UpdateSize(ctx, id, map[string]any{
		"name":       strings.TrimSpace(request.Name),
		"sort_order": request.SortOrder,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("size name is already used")
		}
		return nil, apperror.Internal(err)
	}
	return s.GetSize(ctx, id)
}

func (s *Service) DeleteSize(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetSize(ctx, id); err != nil {
		return err
	}
	if err := s.repo.DeleteSize(ctx, id); err != nil {
		if database.IsForeignKeyViolation(err) {
			return apperror.Conflict("size is used by product variants")
		}
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) ListColors(ctx context.Context) ([]models.Color, error) {
	colors, err := s.repo.ListColors(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return colors, nil
}

func (s *Service) GetColor(ctx context.Context, id uuid.UUID) (*models.Color, error) {
	color, err := s.repo.FindColor(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("color not found")
		}
		return nil, apperror.Internal(err)
	}
	return color, nil
}

func (s *Service) CreateColor(ctx context.Context, request ColorRequest) (*models.Color, error) {
	hexCode := strings.TrimSpace(request.HexCode)
	if hexCode != "" && (len(hexCode) != 7 || !strings.HasPrefix(hexCode, "#")) {
		return nil, apperror.BadRequest("hex_code must use #RRGGBB format")
	}

	color := &models.Color{
		Name:    strings.TrimSpace(request.Name),
		HexCode: strings.ToUpper(hexCode),
	}
	if err := s.repo.CreateColor(ctx, color); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("color name is already used")
		}
		return nil, apperror.Internal(err)
	}
	return color, nil
}

func (s *Service) UpdateColor(ctx context.Context, id uuid.UUID, request ColorRequest) (*models.Color, error) {
	if _, err := s.GetColor(ctx, id); err != nil {
		return nil, err
	}

	hexCode := strings.TrimSpace(request.HexCode)
	if hexCode != "" && (len(hexCode) != 7 || !strings.HasPrefix(hexCode, "#")) {
		return nil, apperror.BadRequest("hex_code must use #RRGGBB format")
	}

	err := s.repo.UpdateColor(ctx, id, map[string]any{
		"name":     strings.TrimSpace(request.Name),
		"hex_code": strings.ToUpper(hexCode),
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("color name is already used")
		}
		return nil, apperror.Internal(err)
	}
	return s.GetColor(ctx, id)
}

func (s *Service) DeleteColor(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetColor(ctx, id); err != nil {
		return err
	}
	if err := s.repo.DeleteColor(ctx, id); err != nil {
		if database.IsForeignKeyViolation(err) {
			return apperror.Conflict("color is used by product variants")
		}
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) ListProductColors(ctx context.Context, productID uuid.UUID) ([]models.ProductColor, error) {
	if _, err := s.GetProduct(ctx, productID); err != nil {
		return nil, err
	}
	colors, err := s.repo.ListProductColors(ctx, productID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return colors, nil
}

func (s *Service) CreateProductColor(ctx context.Context, productID uuid.UUID, request ProductColorRequest) (*models.ProductColor, error) {
	if _, err := s.GetProduct(ctx, productID); err != nil {
		return nil, err
	}
	hexCode, err := normalizeHexCode(request.HexCode)
	if err != nil {
		return nil, err
	}
	color := &models.ProductColor{
		ProductID: productID,
		Name:      strings.TrimSpace(request.Name),
		HexCode:   hexCode,
		SortOrder: request.SortOrder,
	}
	if err := s.repo.CreateProductColor(ctx, color); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("color name is already used for this product")
		}
		return nil, apperror.Internal(err)
	}
	return color, nil
}

func (s *Service) UpdateProductColor(ctx context.Context, productID, colorID uuid.UUID, request ProductColorRequest) (*models.ProductColor, error) {
	if _, err := s.repo.FindProductColor(ctx, productID, colorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("product color not found")
		}
		return nil, apperror.Internal(err)
	}
	hexCode, err := normalizeHexCode(request.HexCode)
	if err != nil {
		return nil, err
	}
	err = s.repo.UpdateProductColor(ctx, productID, colorID, map[string]any{
		"name":       strings.TrimSpace(request.Name),
		"hex_code":   hexCode,
		"sort_order": request.SortOrder,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("color name is already used for this product")
		}
		return nil, apperror.Internal(err)
	}
	return s.repo.FindProductColor(ctx, productID, colorID)
}

func (s *Service) DeleteProductColor(ctx context.Context, productID, colorID uuid.UUID) error {
	color, err := s.repo.FindProductColor(ctx, productID, colorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("product color not found")
		}
		return apperror.Internal(err)
	}
	used, err := s.repo.ProductColorUsedByVariants(ctx, color.ID)
	if err != nil {
		return apperror.Internal(err)
	}
	if used {
		return apperror.Conflict("color is used by product variants")
	}
	if err := s.repo.DeleteProductColor(ctx, productID, colorID); err != nil {
		return apperror.Internal(err)
	}
	_ = s.storage.Delete(ctx, color.ImageKey)
	return nil
}

func (s *Service) UploadProductColorImage(ctx context.Context, productID, colorID uuid.UUID, fileHeader *multipart.FileHeader) (*models.ProductColor, error) {
	if _, err := s.repo.FindProductColor(ctx, productID, colorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("product color not found")
		}
		return nil, apperror.Internal(err)
	}
	if !s.storage.Enabled() {
		return nil, apperror.New(503, "file storage is not configured")
	}
	upload, err := storage.OpenValidated(fileHeader, storage.ImageTypes, storage.MaxProductImageSize)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrEmptyFile):
			return nil, apperror.BadRequest("image file is empty")
		case errors.Is(err, storage.ErrFileTooLarge):
			return nil, apperror.BadRequest("image must be at most 10 MB")
		case errors.Is(err, storage.ErrUnsupportedType):
			return nil, apperror.BadRequest("image must be JPEG, PNG, or WebP")
		default:
			return nil, apperror.BadRequest("failed to read uploaded image")
		}
	}
	defer upload.File.Close()

	key := storage.BuildKey("catalog/products", productID.String()+"/colors/"+colorID.String(), "color"+upload.Extension)
	if _, err := s.storage.Upload(ctx, key, upload.File, upload.ContentType, upload.Size); err != nil {
		return nil, apperror.Wrap(502, "failed to upload color image", err)
	}
	previous, err := s.repo.SetProductColorImage(ctx, productID, colorID, key)
	if err != nil {
		_ = s.storage.Delete(ctx, key)
		return nil, apperror.Internal(err)
	}
	if previous != "" && previous != key {
		_ = s.storage.Delete(ctx, previous)
	}
	return s.repo.FindProductColor(ctx, productID, colorID)
}

func (s *Service) DeleteProductColorImage(ctx context.Context, productID, colorID uuid.UUID) error {
	color, err := s.repo.FindProductColor(ctx, productID, colorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("product color not found")
		}
		return apperror.Internal(err)
	}
	if _, err := s.repo.SetProductColorImage(ctx, productID, colorID, ""); err != nil {
		return apperror.Internal(err)
	}
	_ = s.storage.Delete(ctx, color.ImageKey)
	return nil
}

func normalizeHexCode(value string) (string, error) {
	hexCode := strings.TrimSpace(value)
	if hexCode == "" {
		return "", nil
	}
	if len(hexCode) != 7 || !strings.HasPrefix(hexCode, "#") {
		return "", apperror.BadRequest("hex_code must use #RRGGBB format")
	}
	return strings.ToUpper(hexCode), nil
}

func (s *Service) ListProducts(ctx context.Context, filter ProductListFilter) ([]models.Product, int64, error) {
	products, total, err := s.repo.ListProducts(ctx, filter)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return products, total, nil
}

func (s *Service) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := s.repo.FindProduct(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("product not found")
		}
		return nil, apperror.Internal(err)
	}
	return product, nil
}

func (s *Service) CreateProduct(ctx context.Context, actor *models.User, request CreateProductRequest) (*models.Product, error) {
	if request.CategoryID != nil {
		if _, err := s.GetCategory(ctx, *request.CategoryID); err != nil {
			return nil, err
		}
	}

	status := request.Status
	if status == "" {
		status = models.ProductStatusDraft
	}

	slug, err := s.uniqueProductSlug(ctx, request.Name, request.Slug, nil)
	if err != nil {
		return nil, err
	}

	colors, colorByID, legacyColorIDs, err := s.prepareProductColors(ctx, request.Colors, request.Variants)
	if err != nil {
		return nil, err
	}
	variants, labels, err := s.buildVariants(ctx, request.Variants, colorByID, legacyColorIDs)
	if err != nil {
		return nil, err
	}

	badgeTone := request.BadgeTone
	if badgeTone == "" {
		badgeTone = "light"
	}

	product := &models.Product{
		CategoryID:  request.CategoryID,
		Name:        strings.TrimSpace(request.Name),
		Slug:        slug,
		Description: strings.TrimSpace(request.Description),
		Status:      status,
		ShopeeURL:   strings.TrimSpace(request.ShopeeURL),
		ShopeePrice: request.ShopeePrice,
		Details:     cleanDetails(request.Details),
		Badge:       strings.TrimSpace(request.Badge),
		BadgeTone:   badgeTone,
		SortOrder:   request.SortOrder,
	}
	if actor != nil {
		product.CreatedByID = &actor.ID
	}

	var actorID *uuid.UUID
	if actor != nil {
		actorID = &actor.ID
	}
	movements := make([]models.StockMovement, 0, len(variants))
	for index, variant := range variants {
		if variant.Stock <= 0 {
			continue
		}
		movements = append(movements, models.StockMovement{
			SKU:           variant.SKU,
			VariantLabel:  labels[index],
			Type:          models.StockMovementTypeInitial,
			Quantity:      variant.Stock,
			StockBefore:   0,
			StockAfter:    variant.Stock,
			ReferenceType: models.StockReferenceTypeProduct,
			Note:          "initial stock",
			CreatedByID:   actorID,
		})
	}

	if err := s.repo.CreateProduct(ctx, product, colors, variants, movements); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("product slug, variant SKU, or variant combination is already used")
		}
		return nil, apperror.Internal(err)
	}

	return s.GetProduct(ctx, product.ID)
}

func (s *Service) UpdateProduct(ctx context.Context, id uuid.UUID, request UpdateProductRequest) (*models.Product, error) {
	product, err := s.GetProduct(ctx, id)
	if err != nil {
		return nil, err
	}
	if request.CategoryID != nil {
		if _, err := s.GetCategory(ctx, *request.CategoryID); err != nil {
			return nil, err
		}
	}

	slug, err := s.uniqueProductSlug(ctx, request.Name, request.Slug, &product.ID)
	if err != nil {
		return nil, err
	}

	var categoryID any
	if request.CategoryID != nil {
		categoryID = *request.CategoryID
	}

	badgeTone := request.BadgeTone
	if badgeTone == "" {
		badgeTone = "light"
	}
	detailsJSON, err := json.Marshal(cleanDetails(request.Details))
	if err != nil {
		return nil, apperror.Internal(err)
	}

	err = s.repo.UpdateProduct(ctx, id, map[string]any{
		"category_id":  categoryID,
		"name":         strings.TrimSpace(request.Name),
		"slug":         slug,
		"description":  strings.TrimSpace(request.Description),
		"status":       request.Status,
		"shopee_url":   strings.TrimSpace(request.ShopeeURL),
		"shopee_price": request.ShopeePrice,
		"details":      gorm.Expr("?::jsonb", string(detailsJSON)),
		"badge":        strings.TrimSpace(request.Badge),
		"badge_tone":   badgeTone,
		"sort_order":   request.SortOrder,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("product slug is already used")
		}
		return nil, apperror.Internal(err)
	}
	return s.GetProduct(ctx, id)
}

func (s *Service) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if _, err := s.GetProduct(ctx, id); err != nil {
		return err
	}
	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) CreateVariant(ctx context.Context, actor *models.User, productID uuid.UUID, request CreateVariantRequest) (*models.ProductVariant, error) {
	if _, err := s.GetProduct(ctx, productID); err != nil {
		return nil, err
	}
	size, err := s.GetSize(ctx, request.SizeID)
	if err != nil {
		return nil, err
	}
	color, err := s.resolveProductColor(ctx, productID, request.ProductColorID, request.ColorID)
	if err != nil {
		return nil, err
	}

	sku := strings.ToUpper(strings.TrimSpace(request.SKU))
	comboExists, err := s.repo.VariantComboExists(ctx, productID, request.SizeID, color.ID, nil)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if comboExists {
		return nil, apperror.Conflict("variant combination already exists for this product")
	}
	skuExists, err := s.repo.VariantSKUExists(ctx, sku, nil)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if skuExists {
		return nil, apperror.Conflict("variant SKU is already used")
	}

	isActive := true
	if request.IsActive != nil {
		isActive = *request.IsActive
	}
	variant := &models.ProductVariant{
		ProductID:      productID,
		SizeID:         request.SizeID,
		ProductColorID: &color.ID,
		SKU:            sku,
		Price:          request.Price,
		Stock:          request.Stock,
		IsActive:       isActive,
	}
	var movement *models.StockMovement
	if variant.Stock > 0 {
		var actorID *uuid.UUID
		if actor != nil {
			actorID = &actor.ID
		}
		movement = &models.StockMovement{
			SKU:           variant.SKU,
			VariantLabel:  size.Name + " / " + color.Name,
			Type:          models.StockMovementTypeInitial,
			Quantity:      variant.Stock,
			StockBefore:   0,
			StockAfter:    variant.Stock,
			ReferenceType: models.StockReferenceTypeProduct,
			Note:          "initial stock",
			CreatedByID:   actorID,
		}
	}

	if err := s.repo.CreateVariant(ctx, variant, movement); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("variant SKU or combination is already used")
		}
		return nil, apperror.Internal(err)
	}
	return s.repo.FindVariant(ctx, productID, variant.ID)
}

func (s *Service) UpdateVariant(ctx context.Context, productID, variantID uuid.UUID, request UpdateVariantRequest) (*models.ProductVariant, error) {
	variant, err := s.repo.FindVariant(ctx, productID, variantID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("variant not found")
		}
		return nil, apperror.Internal(err)
	}
	if _, err := s.GetSize(ctx, request.SizeID); err != nil {
		return nil, err
	}
	color, err := s.resolveProductColor(ctx, productID, request.ProductColorID, request.ColorID)
	if err != nil {
		return nil, err
	}

	sku := strings.ToUpper(strings.TrimSpace(request.SKU))
	comboExists, err := s.repo.VariantComboExists(ctx, productID, request.SizeID, color.ID, &variant.ID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if comboExists {
		return nil, apperror.Conflict("variant combination already exists for this product")
	}
	skuExists, err := s.repo.VariantSKUExists(ctx, sku, &variant.ID)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if skuExists {
		return nil, apperror.Conflict("variant SKU is already used")
	}

	isActive := variant.IsActive
	if request.IsActive != nil {
		isActive = *request.IsActive
	}

	err = s.repo.UpdateVariant(ctx, variantID, map[string]any{
		"size_id":          request.SizeID,
		"product_color_id": color.ID,
		"color_id":         nil,
		"sku":              sku,
		"price":            request.Price,
		"is_active":        isActive,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, apperror.Conflict("variant SKU or combination is already used")
		}
		return nil, apperror.Internal(err)
	}
	return s.repo.FindVariant(ctx, productID, variantID)
}

func (s *Service) DeleteVariant(ctx context.Context, productID, variantID uuid.UUID) error {
	if _, err := s.repo.FindVariant(ctx, productID, variantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("variant not found")
		}
		return apperror.Internal(err)
	}
	if err := s.repo.DeleteVariant(ctx, variantID); err != nil {
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) UpdateStock(ctx context.Context, actor *models.User, productID, variantID uuid.UUID, stock int, note string) (*models.ProductVariant, error) {
	if _, err := s.repo.FindVariant(ctx, productID, variantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("variant not found")
		}
		return nil, apperror.Internal(err)
	}
	var actorID *uuid.UUID
	if actor != nil {
		actorID = &actor.ID
	}
	if err := s.repo.UpdateStockWithMovement(ctx, variantID, stock, actorID, strings.TrimSpace(note)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("variant not found")
		}
		return nil, apperror.Internal(err)
	}
	return s.repo.FindVariant(ctx, productID, variantID)
}

func (s *Service) ListStockMovements(ctx context.Context, productID, variantID uuid.UUID, page, perPage int) ([]models.StockMovement, int64, error) {
	if _, err := s.repo.FindVariant(ctx, productID, variantID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, apperror.NotFound("variant not found")
		}
		return nil, 0, apperror.Internal(err)
	}
	movements, total, err := s.repo.ListStockMovements(ctx, productID, variantID, page, perPage)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return movements, total, nil
}

func (s *Service) UploadImages(ctx context.Context, productID uuid.UUID, files []*multipart.FileHeader) ([]models.ProductImage, error) {
	if _, err := s.GetProduct(ctx, productID); err != nil {
		return nil, err
	}
	if !s.storage.Enabled() {
		return nil, apperror.New(503, "file storage is not configured")
	}
	if len(files) == 0 {
		return nil, apperror.BadRequest("at least one image is required")
	}

	maxSort, hasPrimary, err := s.repo.ImageStats(ctx, productID)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	uploadedKeys := make([]string, 0, len(files))
	images := make([]models.ProductImage, 0, len(files))
	for index, fileHeader := range files {
		upload, err := storage.OpenValidated(fileHeader, storage.ImageTypes, storage.MaxProductImageSize)
		if err != nil {
			s.cleanupKeys(ctx, uploadedKeys)
			switch {
			case errors.Is(err, storage.ErrEmptyFile):
				return nil, apperror.BadRequest("image file is empty")
			case errors.Is(err, storage.ErrFileTooLarge):
				return nil, apperror.BadRequest("each image must be at most 10 MB")
			case errors.Is(err, storage.ErrUnsupportedType):
				return nil, apperror.BadRequest("images must be JPEG, PNG, or WebP")
			default:
				return nil, apperror.BadRequest("failed to read uploaded image")
			}
		}

		key := storage.BuildKey("catalog/products", productID.String(), "image"+upload.Extension)
		if _, err := s.storage.Upload(ctx, key, upload.File, upload.ContentType, upload.Size); err != nil {
			_ = upload.File.Close()
			s.cleanupKeys(ctx, uploadedKeys)
			return nil, apperror.Wrap(502, "failed to upload image", err)
		}
		_ = upload.File.Close()
		uploadedKeys = append(uploadedKeys, key)

		images = append(images, models.ProductImage{
			ProductID: productID,
			ObjectKey: key,
			SortOrder: maxSort + 1 + index,
			IsPrimary: !hasPrimary && index == 0,
		})
	}

	if err := s.repo.AddImages(ctx, images); err != nil {
		s.cleanupKeys(ctx, uploadedKeys)
		return nil, apperror.Internal(err)
	}
	return images, nil
}

func (s *Service) DeleteImage(ctx context.Context, productID, imageID uuid.UUID) error {
	if _, err := s.repo.FindImage(ctx, productID, imageID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("image not found")
		}
		return apperror.Internal(err)
	}

	objectKey, err := s.repo.DeleteImage(ctx, productID, imageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("image not found")
		}
		return apperror.Internal(err)
	}
	_ = s.storage.Delete(ctx, objectKey)
	return nil
}

func (s *Service) SetPrimaryImage(ctx context.Context, productID, imageID uuid.UUID) error {
	if err := s.repo.SetPrimaryImage(ctx, productID, imageID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("image not found")
		}
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) prepareProductColors(ctx context.Context, requests []ProductColorRequest, variants []CreateVariantRequest) ([]models.ProductColor, map[uuid.UUID]models.ProductColor, map[uuid.UUID]uuid.UUID, error) {
	colors := make([]models.ProductColor, 0, len(requests))
	byID := make(map[uuid.UUID]models.ProductColor, len(requests))
	byName := make(map[string]uuid.UUID, len(requests))
	legacyColorIDs := make(map[uuid.UUID]uuid.UUID)

	for _, request := range requests {
		id := uuid.New()
		if request.ID != nil && *request.ID != uuid.Nil {
			id = *request.ID
		}
		name := strings.TrimSpace(request.Name)
		if name == "" {
			return nil, nil, nil, apperror.BadRequest("color name is required")
		}
		normalizedName := strings.ToLower(name)
		if _, exists := byID[id]; exists {
			return nil, nil, nil, apperror.BadRequest("duplicate product color id")
		}
		if _, exists := byName[normalizedName]; exists {
			return nil, nil, nil, apperror.BadRequest("duplicate product color name")
		}
		hexCode, err := normalizeHexCode(request.HexCode)
		if err != nil {
			return nil, nil, nil, err
		}
		color := models.ProductColor{
			Base:      models.Base{ID: id},
			Name:      name,
			HexCode:   hexCode,
			SortOrder: request.SortOrder,
		}
		colors = append(colors, color)
		byID[id] = color
		byName[normalizedName] = id
	}

	// Backward compatibility for clients that still send global color_id.
	for _, request := range variants {
		if request.ProductColorID != nil || request.ColorID == nil {
			continue
		}
		if _, exists := legacyColorIDs[*request.ColorID]; exists {
			continue
		}
		globalColor, err := s.GetColor(ctx, *request.ColorID)
		if err != nil {
			return nil, nil, nil, err
		}
		normalizedName := strings.ToLower(globalColor.Name)
		if id, exists := byName[normalizedName]; exists {
			legacyColorIDs[*request.ColorID] = id
			continue
		}
		id := uuid.New()
		color := models.ProductColor{
			Base:    models.Base{ID: id},
			Name:    globalColor.Name,
			HexCode: globalColor.HexCode,
		}
		colors = append(colors, color)
		byID[id] = color
		byName[normalizedName] = id
		legacyColorIDs[*request.ColorID] = id
	}

	return colors, byID, legacyColorIDs, nil
}

func (s *Service) buildVariants(ctx context.Context, requests []CreateVariantRequest, colorsByID map[uuid.UUID]models.ProductColor, legacyColorIDs map[uuid.UUID]uuid.UUID) ([]models.ProductVariant, []string, error) {
	variants := make([]models.ProductVariant, 0, len(requests))
	labels := make([]string, 0, len(requests))
	seenCombos := make(map[string]struct{}, len(requests))
	seenSKUs := make(map[string]struct{}, len(requests))

	for _, request := range requests {
		size, err := s.GetSize(ctx, request.SizeID)
		if err != nil {
			return nil, nil, err
		}

		colorID, err := s.resolveVariantColorID(request.ProductColorID, request.ColorID, colorsByID, legacyColorIDs)
		if err != nil {
			return nil, nil, err
		}
		color, exists := colorsByID[colorID]
		if !exists {
			return nil, nil, apperror.BadRequest("product color not found")
		}

		combo := fmt.Sprintf("%s:%s", request.SizeID, colorID)
		if _, exists := seenCombos[combo]; exists {
			return nil, nil, apperror.BadRequest("duplicate size and color combination in variants")
		}
		seenCombos[combo] = struct{}{}

		sku := strings.ToUpper(strings.TrimSpace(request.SKU))
		if _, exists := seenSKUs[sku]; exists {
			return nil, nil, apperror.BadRequest("duplicate SKU in variants")
		}
		seenSKUs[sku] = struct{}{}

		isActive := true
		if request.IsActive != nil {
			isActive = *request.IsActive
		}
		productColorID := colorID
		variants = append(variants, models.ProductVariant{
			SizeID:         request.SizeID,
			ProductColorID: &productColorID,
			SKU:            sku,
			Price:          request.Price,
			Stock:          request.Stock,
			IsActive:       isActive,
		})
		labels = append(labels, size.Name+" / "+color.Name)
	}

	return variants, labels, nil
}

func (s *Service) resolveVariantColorID(productColorID, legacyColorID *uuid.UUID, colorsByID map[uuid.UUID]models.ProductColor, legacyColorIDs map[uuid.UUID]uuid.UUID) (uuid.UUID, error) {
	if productColorID != nil {
		return *productColorID, nil
	}
	if legacyColorID != nil {
		if id, exists := legacyColorIDs[*legacyColorID]; exists {
			return id, nil
		}
	}
	return uuid.Nil, apperror.BadRequest("product color is required")
}

func (s *Service) resolveProductColor(ctx context.Context, productID uuid.UUID, productColorID, legacyColorID *uuid.UUID) (*models.ProductColor, error) {
	if productColorID != nil {
		color, err := s.repo.FindProductColor(ctx, productID, *productColorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperror.BadRequest("product color not found")
			}
			return nil, apperror.Internal(err)
		}
		return color, nil
	}

	if legacyColorID != nil {
		globalColor, err := s.GetColor(ctx, *legacyColorID)
		if err != nil {
			return nil, err
		}
		color, err := s.repo.FindProductColorByName(ctx, productID, globalColor.Name)
		if err == nil {
			return color, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Internal(err)
		}
		color = &models.ProductColor{
			ProductID: productID,
			Name:      globalColor.Name,
			HexCode:   globalColor.HexCode,
		}
		if err := s.repo.CreateProductColor(ctx, color); err != nil {
			return nil, apperror.Internal(err)
		}
		return color, nil
	}

	return nil, apperror.BadRequest("product color is required")
}

func (s *Service) uniqueProductSlug(ctx context.Context, name, requested string, excludeID *uuid.UUID) (string, error) {
	base := strings.TrimSpace(requested)
	if base == "" {
		base = name
	}
	slug, err := utils.UniqueSlug(base, func(candidate string) (bool, error) {
		return s.repo.ProductSlugExists(ctx, candidate, excludeID)
	})
	if err != nil {
		return "", apperror.Internal(err)
	}
	return slug, nil
}

func (s *Service) uniqueCategorySlug(ctx context.Context, name, requested string, excludeID *uuid.UUID) (string, error) {
	base := strings.TrimSpace(requested)
	if base == "" {
		base = name
	}
	slug, err := utils.UniqueSlug(base, func(candidate string) (bool, error) {
		return s.repo.CategorySlugExists(ctx, candidate, excludeID)
	})
	if err != nil {
		return "", apperror.Internal(err)
	}
	return slug, nil
}

func cleanDetails(details []string) []string {
	result := make([]string, 0, len(details))
	for _, detail := range details {
		detail = strings.TrimSpace(detail)
		if detail != "" {
			result = append(result, detail)
		}
	}
	return result
}

func (s *Service) cleanupKeys(ctx context.Context, keys []string) {
	for _, key := range keys {
		_ = s.storage.Delete(ctx, key)
	}
}
