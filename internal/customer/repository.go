package customer

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/models"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Customer, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Customer{})
	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR phone ILIKE ? OR email ILIKE ?", pattern, pattern, pattern)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var customers []models.Customer
	err := query.
		Order("created_at DESC").
		Offset((filter.Page - 1) * filter.PerPage).
		Limit(filter.PerPage).
		Find(&customers).Error
	return customers, total, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.db.WithContext(ctx).First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *Repository) Create(ctx context.Context, customer *models.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Customer{}).Where("id = ?", id).Updates(values).Error
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Customer{}, "id = ?", id).Error
}

func (r *Repository) CountOrders(ctx context.Context, customerID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Order{}).Where("customer_id = ?", customerID).Count(&count).Error
	return count, err
}
