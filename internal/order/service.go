package order

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"jianco-api/internal/customer"
	"jianco-api/internal/models"
	"jianco-api/pkg/apperror"
)

type Service struct {
	repo      *Repository
	customers *customer.Repository
	location  *time.Location
}

func NewService(repo *Repository, customers *customer.Repository, location *time.Location) *Service {
	return &Service{repo: repo, customers: customers, location: location}
}

func (s *Service) List(ctx context.Context, filter ListFilter) ([]models.Order, int64, error) {
	orders, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, apperror.Internal(err)
	}
	return orders, total, nil
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("order not found")
		}
		return nil, apperror.Internal(err)
	}
	return order, nil
}

func (s *Service) Create(ctx context.Context, actor *models.User, request CreateRequest) (*models.Order, error) {
	customerModel, err := s.customers.FindByID(ctx, request.CustomerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.BadRequest("customer not found")
		}
		return nil, apperror.Internal(err)
	}

	items, subtotal, err := s.buildItems(ctx, request.Items)
	if err != nil {
		return nil, err
	}
	if request.DiscountAmount > subtotal {
		return nil, apperror.BadRequest("discount_amount cannot exceed subtotal")
	}

	total := subtotal - request.DiscountAmount + request.ShippingCost
	if total < 0 {
		return nil, apperror.BadRequest("total_amount cannot be negative")
	}

	recipientName := firstNonEmpty(request.RecipientName, customerModel.Name)
	recipientPhone := firstNonEmpty(request.RecipientPhone, customerModel.Phone)
	shippingAddress := firstNonEmpty(request.ShippingAddress, customerModel.Address)
	if recipientName == "" {
		return nil, apperror.BadRequest("recipient_name is required")
	}
	if recipientPhone == "" {
		return nil, apperror.BadRequest("recipient_phone is required")
	}
	if shippingAddress == "" {
		return nil, apperror.BadRequest("shipping_address is required")
	}

	order := &models.Order{
		CustomerID:         customerModel.ID,
		Status:             models.OrderStatusPending,
		PaymentStatus:      models.PaymentStatusUnpaid,
		PaymentMethod:      strings.TrimSpace(request.PaymentMethod),
		ShippingStatus:     models.ShippingStatusPending,
		Courier:            strings.TrimSpace(request.Courier),
		RecipientName:      recipientName,
		RecipientPhone:     recipientPhone,
		ShippingAddress:    shippingAddress,
		ShippingCity:       firstNonEmpty(request.ShippingCity, customerModel.City),
		ShippingProvince:   firstNonEmpty(request.ShippingProvince, customerModel.Province),
		ShippingPostalCode: firstNonEmpty(request.ShippingPostalCode, customerModel.PostalCode),
		Subtotal:           subtotal,
		DiscountAmount:     request.DiscountAmount,
		ShippingCost:       request.ShippingCost,
		TotalAmount:        total,
		Notes:              strings.TrimSpace(request.Notes),
	}
	if actor != nil {
		order.CreatedByID = &actor.ID
	}

	var lastErr error
	for attempt := 0; attempt < 5; attempt++ {
		orderNumber, err := generateOrderNumber(s.location)
		if err != nil {
			return nil, apperror.Internal(err)
		}

		candidate := *order
		candidate.ID = uuid.Nil
		candidate.OrderNumber = orderNumber
		if err := s.repo.Create(ctx, &candidate, items); err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				lastErr = err
				continue
			}
			return nil, apperror.Internal(err)
		}
		return s.Get(ctx, candidate.ID)
	}

	return nil, apperror.Internal(fmt.Errorf("failed to generate unique order number: %w", lastErr))
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, request UpdateRequest) (*models.Order, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	customerModel, err := s.customers.FindByID(ctx, request.CustomerID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.BadRequest("customer not found")
		}
		return nil, apperror.Internal(err)
	}

	if request.DiscountAmount > current.Subtotal {
		return nil, apperror.BadRequest("discount_amount cannot exceed subtotal")
	}
	total := current.Subtotal - request.DiscountAmount + request.ShippingCost
	if total < 0 {
		return nil, apperror.BadRequest("total_amount cannot be negative")
	}

	err = s.repo.UpdateHeader(ctx, id, map[string]any{
		"customer_id":          customerModel.ID,
		"recipient_name":       strings.TrimSpace(request.RecipientName),
		"recipient_phone":      strings.TrimSpace(request.RecipientPhone),
		"shipping_address":     strings.TrimSpace(request.ShippingAddress),
		"shipping_city":        strings.TrimSpace(request.ShippingCity),
		"shipping_province":    strings.TrimSpace(request.ShippingProvince),
		"shipping_postal_code": strings.TrimSpace(request.ShippingPostalCode),
		"payment_method":       strings.TrimSpace(request.PaymentMethod),
		"courier":              strings.TrimSpace(request.Courier),
		"discount_amount":      request.DiscountAmount,
		"shipping_cost":        request.ShippingCost,
		"total_amount":         total,
		"notes":                strings.TrimSpace(request.Notes),
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, id)
}

func (s *Service) UpdateItems(ctx context.Context, id uuid.UUID, request UpdateItemsRequest) (*models.Order, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if current.Status != models.OrderStatusPending || current.StockDeducted {
		return nil, apperror.Conflict("order items can only be edited while the order is pending")
	}

	items, subtotal, err := s.buildItems(ctx, request.Items)
	if err != nil {
		return nil, err
	}
	if current.DiscountAmount > subtotal {
		return nil, apperror.BadRequest("discount_amount cannot exceed subtotal")
	}
	total := subtotal - current.DiscountAmount + current.ShippingCost
	if total < 0 {
		return nil, apperror.BadRequest("total_amount cannot be negative")
	}

	err = s.repo.ReplaceItems(ctx, id, current.Status, items, subtotal, current.DiscountAmount, current.ShippingCost, total)
	if err != nil {
		if errors.Is(err, ErrOrderChanged) {
			return nil, apperror.Conflict("order status changed, refresh and try again")
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("order not found")
		}
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, id)
}

func (s *Service) UpdateStatus(ctx context.Context, actor *models.User, id uuid.UUID, newStatus string) (*models.Order, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canTransitionOrder(current.Status, newStatus) {
		return nil, apperror.BadRequest("invalid status transition from " + current.Status + " to " + newStatus)
	}
	if current.Status == newStatus {
		return current, nil
	}

	var actorID *uuid.UUID
	if actor != nil {
		actorID = &actor.ID
	}
	err = s.repo.UpdateStatus(ctx, id, current.Status, newStatus, actorID, time.Now())
	if err != nil {
		if errors.Is(err, ErrOrderChanged) {
			return nil, apperror.Conflict("order status changed, refresh and try again")
		}
		if errors.Is(err, ErrInsufficientStock) {
			return nil, apperror.Conflict("insufficient stock: " + strings.TrimPrefix(err.Error(), ErrInsufficientStock.Error()+": "))
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("order not found")
		}
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, id)
}

func (s *Service) UpdatePayment(ctx context.Context, id uuid.UUID, request UpdatePaymentRequest) (*models.Order, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canTransitionPayment(current.PaymentStatus, request.PaymentStatus) {
		return nil, apperror.BadRequest("invalid payment status transition from " + current.PaymentStatus + " to " + request.PaymentStatus)
	}

	if err := s.repo.UpdatePayment(ctx, id, request.PaymentStatus, strings.TrimSpace(request.PaymentMethod), time.Now()); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("order not found")
		}
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, id)
}

func (s *Service) UpdateShipping(ctx context.Context, actor *models.User, id uuid.UUID, request UpdateShippingRequest) (*models.Order, error) {
	current, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if !canTransitionShipping(current.ShippingStatus, request.ShippingStatus) {
		return nil, apperror.BadRequest("invalid shipping status transition from " + current.ShippingStatus + " to " + request.ShippingStatus)
	}

	var actorID *uuid.UUID
	if actor != nil {
		actorID = &actor.ID
	}
	err = s.repo.UpdateShipping(ctx, id, request.ShippingStatus, strings.TrimSpace(request.Courier), strings.TrimSpace(request.TrackingNumber), actorID, time.Now())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("order not found")
		}
		return nil, apperror.Internal(err)
	}
	return s.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, actor *models.User, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	var actorID *uuid.UUID
	if actor != nil {
		actorID = &actor.ID
	}
	if err := s.repo.Delete(ctx, id, actorID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("order not found")
		}
		return apperror.Internal(err)
	}
	return nil
}

func (s *Service) buildItems(ctx context.Context, requests []CreateItemRequest) ([]models.OrderItem, int64, error) {
	items := make([]models.OrderItem, 0, len(requests))
	seen := make(map[uuid.UUID]struct{}, len(requests))
	var subtotal int64

	for _, request := range requests {
		if _, exists := seen[request.VariantID]; exists {
			return nil, 0, apperror.BadRequest("duplicate variant in order items")
		}
		seen[request.VariantID] = struct{}{}

		variant, err := s.repo.FindVariantForOrder(ctx, request.VariantID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, 0, apperror.BadRequest("variant not found")
			}
			return nil, 0, apperror.Internal(err)
		}
		if !variant.IsActive || variant.Product.ID == uuid.Nil || variant.Product.Status != models.ProductStatusActive {
			return nil, 0, apperror.BadRequest("variant is not available")
		}
		if variant.Stock < request.Quantity {
			return nil, 0, apperror.Conflict("insufficient stock for " + variant.SKU)
		}

		lineSubtotal := variant.Price * int64(request.Quantity)
		productID := variant.ProductID
		variantID := variant.ID
		items = append(items, models.OrderItem{
			ProductID:    &productID,
			VariantID:    &variantID,
			ProductName:  variant.Product.Name,
			VariantLabel: variant.Size.Name + " / " + variant.Color.Name,
			SKU:          variant.SKU,
			Price:        variant.Price,
			Quantity:     request.Quantity,
			Subtotal:     lineSubtotal,
		})
		subtotal += lineSubtotal
	}

	return items, subtotal, nil
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}
