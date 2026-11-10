package catalog

import (
	"context"

	"gorm.io/gorm"

	"jianco-api/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func preloadCatalogProduct(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Category").
		Preload("Images", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("is_primary DESC, sort_order ASC, created_at ASC")
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

func (r *Repository) ListActiveProducts(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := preloadCatalogProduct(r.db.WithContext(ctx)).
		Where("status = ?", models.ProductStatusActive).
		Order("sort_order ASC, created_at DESC").
		Find(&products).Error
	return products, err
}

func (r *Repository) FindActiveProductBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	err := preloadCatalogProduct(r.db.WithContext(ctx)).
		Where("slug = ? AND status = ?", slug, models.ProductStatusActive).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) ListActiveCollections(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.WithContext(ctx).
		Model(&models.Category{}).
		Joins("JOIN products ON products.category_id = categories.id AND products.status = ? AND products.deleted_at IS NULL", models.ProductStatusActive).
		Group("categories.id").
		Order("categories.name ASC").
		Find(&categories).Error
	return categories, err
}

func (r *Repository) GetSettings(ctx context.Context) (*models.CatalogSetting, error) {
	var setting models.CatalogSetting
	if err := r.db.WithContext(ctx).Where("key = ?", "default").First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *Repository) CreateSettings(ctx context.Context, setting *models.CatalogSetting) error {
	return r.db.WithContext(ctx).Create(setting).Error
}

func (r *Repository) UpdateSettings(ctx context.Context, values map[string]any) error {
	return r.db.WithContext(ctx).
		Model(&models.CatalogSetting{}).
		Where("key = ?", "default").
		Updates(values).Error
}
