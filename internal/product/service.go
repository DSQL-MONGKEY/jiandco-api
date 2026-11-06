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

	variants, labels, err := s.buildVariants(ctx, request.Variants)
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

	if err := s.repo.CreateProduct(ctx, product, variants, movements); err != nil {
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
	color, err := s.GetColor(ctx, request.ColorID)
	if err != nil {
		return nil, err
	}

	sku := strings.ToUpper(strings.TrimSpace(request.SKU))
	comboExists, err := s.repo.VariantComboExists(ctx, productID, request.SizeID, request.ColorID, nil)
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
		ProductID: productID,
		SizeID:    request.SizeID,
		ColorID:   request.ColorID,
		SKU:       sku,
		Price:     request.Price,
		Stock:     request.Stock,
		IsActive:  isActive,
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
	if _, err := s.GetColor(ctx, request.ColorID); err != nil {
		return nil, err
	}

	sku := strings.ToUpper(strings.TrimSpace(request.SKU))
	comboExists, err := s.repo.VariantComboExists(ctx, productID, request.SizeID, request.ColorID, &variant.ID)
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
		"size_id":   request.SizeID,
		"color_id":  request.ColorID,
		"sku":       sku,
		"price":     request.Price,
		"is_active": isActive,
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
				return nil, apperror.BadRequest("each image must be at most 5 MB")
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

func (s *Service) buildVariants(ctx context.Context, requests []CreateVariantRequest) ([]models.ProductVariant, []string, error) {
	variants := make([]models.ProductVariant, 0, len(requests))
	labels := make([]string, 0, len(requests))
	seenCombos := make(map[string]struct{}, len(requests))
	seenSKUs := make(map[string]struct{}, len(requests))

	for _, request := range requests {
		size, err := s.GetSize(ctx, request.SizeID)
		if err != nil {
			return nil, nil, err
		}
		color, err := s.GetColor(ctx, request.ColorID)
		if err != nil {
			return nil, nil, err
		}

		combo := fmt.Sprintf("%s:%s", request.SizeID, request.ColorID)
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
		variants = append(variants, models.ProductVariant{
			SizeID:   request.SizeID,
			ColorID:  request.ColorID,
			SKU:      sku,
			Price:    request.Price,
			Stock:    request.Stock,
			IsActive: isActive,
		})
		labels = append(labels, size.Name+" / "+color.Name)
	}

	return variants, labels, nil
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
