package accounting

import (
	"context"
	"strings"
	"time"

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

func (r *Repository) List(ctx context.Context, filter ListFilter, start, end *time.Time) ([]models.Expense, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Expense{})
	if search := strings.TrimSpace(filter.Search); search != "" {
		pattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR vendor ILIKE ?", pattern, pattern)
	}
	if filter.Category != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if start != nil {
		query = query.Where("expense_date >= ?", start.Format("2006-01-02"))
	}
	if end != nil {
		query = query.Where("expense_date < ?", end.Format("2006-01-02"))
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var expenses []models.Expense
	err := query.
		Preload("Items").
		Order("expense_date DESC, created_at DESC").
		Offset((filter.Page - 1) * filter.PerPage).
		Limit(filter.PerPage).
		Find(&expenses).Error
	return expenses, total, err
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Expense, error) {
	var expense models.Expense
	if err := r.db.WithContext(ctx).Preload("Items").First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}

func (r *Repository) Create(ctx context.Context, expense *models.Expense, items []models.ExpenseItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(expense).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		for index := range items {
			items[index].ExpenseID = expense.ID
		}
		return tx.Create(&items).Error
	})
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, values map[string]any, items []models.ExpenseItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Expense{}).Where("id = ?", id).Updates(values).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.ExpenseItem{}, "expense_id = ?", id).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		for index := range items {
			items[index].ExpenseID = id
		}
		return tx.Create(&items).Error
	})
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Expense{}, "id = ?", id).Error
}

func (r *Repository) UpdateReceipt(ctx context.Context, id uuid.UUID, objectKey string) error {
	return r.db.WithContext(ctx).Model(&models.Expense{}).Where("id = ?", id).Update("receipt_key", objectKey).Error
}

func (r *Repository) SumExpenses(ctx context.Context, start, end time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.Expense{}).
		Where("expense_date >= ? AND expense_date < ?", start.Format("2006-01-02"), end.Format("2006-01-02")).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *Repository) SumByCategory(ctx context.Context, start, end time.Time) ([]CategoryTotal, error) {
	totals := make([]CategoryTotal, 0)
	err := r.db.WithContext(ctx).
		Model(&models.Expense{}).
		Select("category, COALESCE(SUM(total_amount), 0) AS total").
		Where("expense_date >= ? AND expense_date < ?", start.Format("2006-01-02"), end.Format("2006-01-02")).
		Group("category").
		Order("category ASC").
		Scan(&totals).Error
	return totals, err
}

func (r *Repository) Recent(ctx context.Context, start, end time.Time, limit int) ([]models.Expense, error) {
	var expenses []models.Expense
	err := r.db.WithContext(ctx).
		Preload("Items").
		Where("expense_date >= ? AND expense_date < ?", start.Format("2006-01-02"), end.Format("2006-01-02")).
		Order("expense_date DESC, created_at DESC").
		Limit(limit).
		Find(&expenses).Error
	return expenses, err
}
