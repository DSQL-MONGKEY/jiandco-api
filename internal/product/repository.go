package product

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jianco-api/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListCategories(ctx context.Context, search string) ([]models.Category, error) {
	query := r.db.WithContext(ctx).Order("name ASC")
	if search = strings.TrimSpace(search); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	var categories []models.Category
	err := query.Find(&categories).Error
	return categories, err
}

func (r *Repository) FindCategory(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var category models.Category
	if err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *Repository) CreateCategory(ctx context.Context, category *models.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *Repository) UpdateCategory(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Category{}).Where("id = ?", id).Updates(values).Error
}

func (r *Repository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Product{}).Where("category_id = ?", id).Update("category_id", nil).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Category{}, "id = ?", id).Error
	})
}

func (r *Repository) CategorySlugExists(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.Category{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) ListSizes(ctx context.Context) ([]models.Size, error) {
	var sizes []models.Size
	err := r.db.WithContext(ctx).Order("sort_order ASC, name ASC").Find(&sizes).Error
	return sizes, err
}

func (r *Repository) FindSize(ctx context.Context, id uuid.UUID) (*models.Size, error) {
	var size models.Size
	if err := r.db.WithContext(ctx).First(&size, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &size, nil
}

func (r *Repository) CreateSize(ctx context.Context, size *models.Size) error {
	return r.db.WithContext(ctx).Create(size).Error
}

func (r *Repository) UpdateSize(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Size{}).Where("id = ?", id).Updates(values).Error
}

func (r *Repository) DeleteSize(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Size{}, "id = ?", id).Error
}

func (r *Repository) ListColors(ctx context.Context) ([]models.Color, error) {
	var colors []models.Color
	err := r.db.WithContext(ctx).Order("name ASC").Find(&colors).Error
	return colors, err
}

func (r *Repository) FindColor(ctx context.Context, id uuid.UUID) (*models.Color, error) {
	var color models.Color
	if err := r.db.WithContext(ctx).First(&color, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &color, nil
}

func (r *Repository) CreateColor(ctx context.Context, color *models.Color) error {
	return r.db.WithContext(ctx).Create(color).Error
}

func (r *Repository) UpdateColor(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Color{}).Where("id = ?", id).Updates(values).Error
}

func (r *Repository) DeleteColor(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Color{}, "id = ?", id).Error
}

func (r *Repository) ListProducts(ctx context.Context, filter ProductListFilter) ([]models.Product, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Product{})

	if filter.Search != "" {
		pattern := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ? OR slug ILIKE ?", pattern, pattern)
	}
	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.Product
	err := preloadProduct(query).
		Order("created_at DESC").
		Offset((filter.Page - 1) * filter.PerPage).
		Limit(filter.PerPage).
		Find(&products).Error
	return products, total, err
}

func (r *Repository) FindProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := preloadProduct(r.db.WithContext(ctx)).First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product *models.Product, colors []models.ProductColor, variants []models.ProductVariant, movements []models.StockMovement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(product).Error; err != nil {
			return err
		}
		for index := range colors {
			colors[index].ProductID = product.ID
		}
		if len(colors) > 0 {
			if err := tx.Create(&colors).Error; err != nil {
				return err
			}
		}
		if len(variants) > 0 {
			for index := range variants {
				variants[index].ProductID = product.ID
			}
			if err := tx.Create(&variants).Error; err != nil {
				return err
			}
		}
		if len(movements) == 0 {
			return nil
		}

		variantIDs := make(map[string]uuid.UUID, len(variants))
		for _, variant := range variants {
			variantIDs[variant.SKU] = variant.ID
		}
		for index := range movements {
			movements[index].ProductID = &product.ID
			movements[index].ReferenceID = &product.ID
			if variantID, ok := variantIDs[movements[index].SKU]; ok {
				movements[index].VariantID = &variantID
			}
		}
		return tx.Create(&movements).Error
	})
}

func (r *Repository) UpdateProduct(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Product{}).Where("id = ?", id).Updates(values).Error
}

func (r *Repository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, "id = ?", id).Error
}

func (r *Repository) ProductSlugExists(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.Product{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) FindVariant(ctx context.Context, productID, variantID uuid.UUID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Size").
		Preload("Color").
		Preload("ProductColor").
		First(&variant, "id = ? AND product_id = ?", variantID, productID).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (r *Repository) VariantComboExists(ctx context.Context, productID, sizeID, productColorID uuid.UUID, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.ProductVariant{}).
		Where("product_id = ? AND size_id = ? AND product_color_id = ?", productID, sizeID, productColorID)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) VariantSKUExists(ctx context.Context, sku string, excludeID *uuid.UUID) (bool, error) {
	query := r.db.WithContext(ctx).Model(&models.ProductVariant{}).Where("sku = ?", sku)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CreateVariant(ctx context.Context, variant *models.ProductVariant, movement *models.StockMovement) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(variant).Error; err != nil {
			return err
		}
		if movement == nil {
			return nil
		}
		movement.VariantID = &variant.ID
		movement.ProductID = &variant.ProductID
		movement.ReferenceID = &variant.ProductID
		return tx.Create(movement).Error
	})
}

func (r *Repository) UpdateVariant(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.ProductVariant{}).Where("id = ?", id).Updates(values).Error
}

func (r *Repository) DeleteVariant(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ProductVariant{}, "id = ?", id).Error
}

func (r *Repository) UpdateStockWithMovement(ctx context.Context, variantID uuid.UUID, stock int, actorID *uuid.UUID, note string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var variant models.ProductVariant
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Size").
			Preload("Color").
			Preload("ProductColor").
			First(&variant, "id = ?", variantID).Error; err != nil {
			return err
		}
		if variant.Stock == stock {
			return nil
		}

		if err := tx.Model(&models.ProductVariant{}).Where("id = ?", variantID).Update("stock", stock).Error; err != nil {
			return err
		}

		movement := models.StockMovement{
			VariantID:     &variant.ID,
			ProductID:     &variant.ProductID,
			SKU:           variant.SKU,
			VariantLabel:  variant.Size.Name + " / " + variantColorName(&variant),
			Type:          models.StockMovementTypeAdjustment,
			Quantity:      stock - variant.Stock,
			StockBefore:   variant.Stock,
			StockAfter:    stock,
			ReferenceType: models.StockReferenceTypeManual,
			Note:          note,
			CreatedByID:   actorID,
		}
		return tx.Create(&movement).Error
	})
}

func (r *Repository) ListStockMovements(ctx context.Context, productID, variantID uuid.UUID, page, perPage int) ([]models.StockMovement, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.StockMovement{}).
		Where("product_id = ? AND variant_id = ?", productID, variantID)

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var movements []models.StockMovement
	err := query.
		Preload("CreatedBy").
		Order("created_at DESC").
		Offset((page - 1) * perPage).
		Limit(perPage).
		Find(&movements).Error
	return movements, total, err
}

func (r *Repository) AddImages(ctx context.Context, images []models.ProductImage) error {
	return r.db.WithContext(ctx).Create(&images).Error
}

func (r *Repository) FindImage(ctx context.Context, productID, imageID uuid.UUID) (*models.ProductImage, error) {
	var image models.ProductImage
	if err := r.db.WithContext(ctx).First(&image, "id = ? AND product_id = ?", imageID, productID).Error; err != nil {
		return nil, err
	}
	return &image, nil
}

func (r *Repository) DeleteImage(ctx context.Context, productID, imageID uuid.UUID) (string, error) {
	var objectKey string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var image models.ProductImage
		if err := tx.First(&image, "id = ? AND product_id = ?", imageID, productID).Error; err != nil {
			return err
		}
		objectKey = image.ObjectKey
		if err := tx.Delete(&models.ProductImage{}, "id = ?", imageID).Error; err != nil {
			return err
		}
		if !image.IsPrimary {
			return nil
		}

		var next models.ProductImage
		err := tx.Where("product_id = ?", productID).Order("sort_order ASC, created_at ASC").First(&next).Error
		if err == nil {
			return tx.Model(&models.ProductImage{}).Where("id = ?", next.ID).Update("is_primary", true).Error
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	})
	return objectKey, err
}

func (r *Repository) SetPrimaryImage(ctx context.Context, productID, imageID uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var image models.ProductImage
		if err := tx.First(&image, "id = ? AND product_id = ?", imageID, productID).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.ProductImage{}).Where("product_id = ?", productID).Update("is_primary", false).Error; err != nil {
			return err
		}
		return tx.Model(&models.ProductImage{}).Where("id = ?", imageID).Update("is_primary", true).Error
	})
}

func (r *Repository) ImageStats(ctx context.Context, productID uuid.UUID) (maxSortOrder int, hasPrimary bool, err error) {
	var maxSort int
	if err := r.db.WithContext(ctx).
		Model(&models.ProductImage{}).
		Where("product_id = ?", productID).
		Select("COALESCE(MAX(sort_order), -1)").
		Scan(&maxSort).Error; err != nil {
		return 0, false, err
	}

	var primaryCount int64
	if err := r.db.WithContext(ctx).
		Model(&models.ProductImage{}).
		Where("product_id = ? AND is_primary = ?", productID, true).
		Count(&primaryCount).Error; err != nil {
		return 0, false, err
	}

	return maxSort, primaryCount > 0, nil
}

func (r *Repository) ListProductColors(ctx context.Context, productID uuid.UUID) ([]models.ProductColor, error) {
	var colors []models.ProductColor
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("sort_order ASC, created_at ASC").
		Find(&colors).Error
	return colors, err
}

func (r *Repository) FindProductColor(ctx context.Context, productID, colorID uuid.UUID) (*models.ProductColor, error) {
	var color models.ProductColor
	if err := r.db.WithContext(ctx).
		Where("id = ? AND product_id = ?", colorID, productID).
		First(&color).Error; err != nil {
		return nil, err
	}
	return &color, nil
}

func (r *Repository) FindProductColorByName(ctx context.Context, productID uuid.UUID, name string) (*models.ProductColor, error) {
	var color models.ProductColor
	if err := r.db.WithContext(ctx).
		Where("product_id = ? AND lower(name) = lower(?)", productID, strings.TrimSpace(name)).
		First(&color).Error; err != nil {
		return nil, err
	}
	return &color, nil
}

func (r *Repository) CreateProductColor(ctx context.Context, color *models.ProductColor) error {
	return r.db.WithContext(ctx).Create(color).Error
}

func (r *Repository) UpdateProductColor(ctx context.Context, productID, colorID uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&models.ProductColor{}).
		Where("id = ? AND product_id = ?", colorID, productID).
		Updates(values).Error
}

func (r *Repository) DeleteProductColor(ctx context.Context, productID, colorID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Delete(&models.ProductColor{}, "id = ? AND product_id = ?", colorID, productID).Error
}

func (r *Repository) ProductColorUsedByVariants(ctx context.Context, colorID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&models.ProductVariant{}).
		Where("product_color_id = ?", colorID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) SetProductColorImage(ctx context.Context, productID, colorID uuid.UUID, imageKey string) (string, error) {
	var previous string
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var color models.ProductColor
		if err := tx.Where("id = ? AND product_id = ?", colorID, productID).First(&color).Error; err != nil {
			return err
		}
		previous = color.ImageKey
		return tx.Model(&models.ProductColor{}).Where("id = ?", colorID).Update("image_key", imageKey).Error
	})
	return previous, err
}

func variantColorName(variant *models.ProductVariant) string {
	if variant.ProductColor != nil {
		return variant.ProductColor.Name
	}
	if variant.Color != nil {
		return variant.Color.Name
	}
	return ""
}

func preloadProduct(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Category").
		Preload("Images", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("sort_order ASC, created_at ASC")
		}).
		Preload("Colors", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("sort_order ASC, created_at ASC")
		}).
		Preload("Variants", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("created_at ASC")
		}).
		Preload("Variants.Size").
		Preload("Variants.Color").
		Preload("Variants.ProductColor")
}
