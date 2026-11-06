package overview

import (
	"context"
	"time"

	"gorm.io/gorm"

	"jianco-api/internal/models"
)

type PaidOrder struct {
	PaidAt      time.Time
	TotalAmount int64
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Revenue(ctx context.Context, start, end time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("payment_status = ? AND paid_at >= ? AND paid_at < ?", models.PaymentStatusPaid, start, end).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *Repository) PaidOrders(ctx context.Context, start, end time.Time) ([]PaidOrder, error) {
	var orders []PaidOrder
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Select("paid_at, total_amount").
		Where("payment_status = ? AND paid_at >= ? AND paid_at < ?", models.PaymentStatusPaid, start, end).
		Order("paid_at ASC").
		Scan(&orders).Error
	return orders, err
}

func (r *Repository) CountOrdersCreated(ctx context.Context, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("created_at >= ? AND created_at < ?", start, end).
		Count(&count).Error
	return count, err
}

func (r *Repository) CountPaidOrders(ctx context.Context, start, end time.Time) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Where("payment_status = ? AND paid_at >= ? AND paid_at < ?", models.PaymentStatusPaid, start, end).
		Count(&count).Error
	return count, err
}

func (r *Repository) CountOrdersByStatus(ctx context.Context, start, end time.Time) ([]OrderStatusCount, error) {
	var counts []OrderStatusCount
	err := r.db.WithContext(ctx).
		Model(&models.Order{}).
		Select("status, COUNT(*) AS count").
		Where("created_at >= ? AND created_at < ?", start, end).
		Group("status").
		Order("status ASC").
		Scan(&counts).Error
	return counts, err
}

func (r *Repository) Expenses(ctx context.Context, start, end time.Time) (int64, error) {
	var total int64
	err := r.db.WithContext(ctx).
		Model(&models.Expense{}).
		Where("expense_date >= ? AND expense_date < ?", start.Format("2006-01-02"), end.Format("2006-01-02")).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total).Error
	return total, err
}
