package order

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"jianco-api/internal/models"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrOrderChanged      = errors.New("order status changed")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for index := range items {
			items[index].OrderID = order.ID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := r.db.WithContext(ctx).
		Preload("Customer").
		Preload("Items").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) List(ctx context.Context, filter ListFilter) ([]models.Order, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Order{})

	if filter.Search != "" {
		pattern := "%" + filter.Search + "%"
		query = query.Where(
			"order_number ILIKE ? OR recipient_name ILIKE ? OR recipient_phone ILIKE ? OR tracking_number ILIKE ?",
			pattern, pattern, pattern, pattern,
		)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.PaymentStatus != "" {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}
	if filter.ShippingStatus != "" {
		query = query.Where("shipping_status = ?", filter.ShippingStatus)
	}
	if filter.CustomerID != nil {
		query = query.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.From != nil {
		query = query.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("created_at < ?", *filter.To)
	}

	var total int64
	if err := query.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []models.Order
	err := query.
		Preload("Customer").
		Preload("Items").
		Order("created_at DESC").
		Offset((filter.Page - 1) * filter.PerPage).
		Limit(filter.PerPage).
		Find(&orders).Error
	return orders, total, err
}

func (r *Repository) UpdateHeader(ctx context.Context, id uuid.UUID, values map[string]any) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Updates(values).Error
}

func (r *Repository) ReplaceItems(ctx context.Context, orderID uuid.UUID, expectedStatus string, items []models.OrderItem, subtotal, discount, shipping, total int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&order, "id = ?", orderID).Error; err != nil {
			return err
		}
		if order.Status != expectedStatus || order.StockDeducted {
			return ErrOrderChanged
		}

		if err := tx.Delete(&models.OrderItem{}, "order_id = ?", orderID).Error; err != nil {
			return err
		}
		for index := range items {
			items[index].OrderID = orderID
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return tx.Model(&models.Order{}).Where("id = ?", orderID).Updates(map[string]any{
			"subtotal":        subtotal,
			"discount_amount": discount,
			"shipping_cost":   shipping,
			"total_amount":    total,
		}).Error
	})
}

func (r *Repository) UpdateStatus(ctx context.Context, orderID uuid.UUID, expectedStatus, newStatus string, actorID *uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Items").
			First(&order, "id = ?", orderID).Error; err != nil {
			return err
		}
		if order.Status != expectedStatus {
			return ErrOrderChanged
		}

		values := map[string]any{"status": newStatus}
		if newStatus == models.OrderStatusConfirmed && !order.StockDeducted {
			if err := deductStock(tx, &order, order.Items, actorID); err != nil {
				return err
			}
			values["stock_deducted"] = true
		}
		if newStatus == models.OrderStatusCancelled && order.StockDeducted {
			if err := restoreStock(tx, &order, order.Items, actorID, models.StockMovementTypeCancel, "order cancelled"); err != nil {
				return err
			}
			values["stock_deducted"] = false
		}

		return tx.Model(&models.Order{}).Where("id = ?", orderID).Updates(values).Error
	})
}

func (r *Repository) UpdatePayment(ctx context.Context, orderID uuid.UUID, status, method string, now time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&order, "id = ?", orderID).Error; err != nil {
			return err
		}

		values := map[string]any{"payment_status": status}
		if method != "" {
			values["payment_method"] = method
		}
		switch status {
		case models.PaymentStatusPaid:
			if order.PaidAt == nil {
				values["paid_at"] = now
			}
		case models.PaymentStatusUnpaid:
			values["paid_at"] = nil
		}

		return tx.Model(&models.Order{}).Where("id = ?", orderID).Updates(values).Error
	})
}

func (r *Repository) UpdateShipping(ctx context.Context, orderID uuid.UUID, status, courier, trackingNumber string, actorID *uuid.UUID, now time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Items").
			First(&order, "id = ?", orderID).Error; err != nil {
			return err
		}

		values := map[string]any{"shipping_status": status}
		if courier != "" {
			values["courier"] = courier
		}
		if trackingNumber != "" {
			values["tracking_number"] = trackingNumber
		}

		switch status {
		case models.ShippingStatusPickedUp, models.ShippingStatusInTransit:
			if order.ShippedAt == nil {
				values["shipped_at"] = now
			}
		case models.ShippingStatusDelivered:
			if order.DeliveredAt == nil {
				values["delivered_at"] = now
			}
		case models.ShippingStatusReturned:
			if order.StockDeducted {
				if err := restoreStock(tx, &order, order.Items, actorID, models.StockMovementTypeReturn, "order returned"); err != nil {
					return err
				}
				values["stock_deducted"] = false
			}
		}

		return tx.Model(&models.Order{}).Where("id = ?", orderID).Updates(values).Error
	})
}

func (r *Repository) Delete(ctx context.Context, orderID uuid.UUID, actorID *uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order models.Order
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Preload("Items").
			First(&order, "id = ?", orderID).Error; err != nil {
			return err
		}
		if order.StockDeducted {
			if err := restoreStock(tx, &order, order.Items, actorID, models.StockMovementTypeCancel, "order deleted"); err != nil {
				return err
			}
		}
		return tx.Delete(&models.Order{}, "id = ?", orderID).Error
	})
}

func (r *Repository) FindVariantForOrder(ctx context.Context, variantID uuid.UUID) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	err := r.db.WithContext(ctx).
		Preload("Product").
		Preload("Size").
		Preload("Color").
		First(&variant, "id = ?", variantID).Error
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func deductStock(tx *gorm.DB, order *models.Order, items []models.OrderItem, actorID *uuid.UUID) error {
	for _, item := range items {
		if item.VariantID == nil {
			continue
		}

		var variant models.ProductVariant
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&variant, "id = ?", *item.VariantID).Error; err != nil {
			return err
		}
		if variant.Stock < item.Quantity {
			return fmt.Errorf("%w: %s", ErrInsufficientStock, strings.TrimSpace(item.SKU))
		}

		stockAfter := variant.Stock - item.Quantity
		if err := tx.Model(&models.ProductVariant{}).Where("id = ?", variant.ID).Update("stock", stockAfter).Error; err != nil {
			return err
		}

		productID := variant.ProductID
		if item.ProductID != nil {
			productID = *item.ProductID
		}
		movement := models.StockMovement{
			VariantID:     &variant.ID,
			ProductID:     &productID,
			SKU:           item.SKU,
			VariantLabel:  item.VariantLabel,
			Type:          models.StockMovementTypeOrder,
			Quantity:      -item.Quantity,
			StockBefore:   variant.Stock,
			StockAfter:    stockAfter,
			ReferenceType: models.StockReferenceTypeOrder,
			ReferenceID:   &order.ID,
			Note:          "order " + order.OrderNumber,
			CreatedByID:   actorID,
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
	}
	return nil
}

func restoreStock(tx *gorm.DB, order *models.Order, items []models.OrderItem, actorID *uuid.UUID, movementType, note string) error {
	for _, item := range items {
		if item.VariantID == nil {
			continue
		}

		var variant models.ProductVariant
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&variant, "id = ?", *item.VariantID).Error; err != nil {
			return err
		}

		stockAfter := variant.Stock + item.Quantity
		if err := tx.Model(&models.ProductVariant{}).Where("id = ?", variant.ID).Update("stock", stockAfter).Error; err != nil {
			return err
		}

		productID := variant.ProductID
		if item.ProductID != nil {
			productID = *item.ProductID
		}
		movement := models.StockMovement{
			VariantID:     &variant.ID,
			ProductID:     &productID,
			SKU:           item.SKU,
			VariantLabel:  item.VariantLabel,
			Type:          movementType,
			Quantity:      item.Quantity,
			StockBefore:   variant.Stock,
			StockAfter:    stockAfter,
			ReferenceType: models.StockReferenceTypeOrder,
			ReferenceID:   &order.ID,
			Note:          note,
			CreatedByID:   actorID,
		}
		if err := tx.Create(&movement).Error; err != nil {
			return err
		}
	}
	return nil
}
