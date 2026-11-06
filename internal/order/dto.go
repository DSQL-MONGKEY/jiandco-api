package order

import (
	"time"

	"github.com/google/uuid"

	"jianco-api/internal/customer"
	"jianco-api/internal/models"
)

type CreateItemRequest struct {
	VariantID uuid.UUID `json:"variant_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type CreateRequest struct {
	CustomerID         uuid.UUID           `json:"customer_id" validate:"required"`
	Items              []CreateItemRequest `json:"items" validate:"required,min=1,dive"`
	RecipientName      string              `json:"recipient_name" validate:"max=100"`
	RecipientPhone     string              `json:"recipient_phone" validate:"max=30"`
	ShippingAddress    string              `json:"shipping_address"`
	ShippingCity       string              `json:"shipping_city" validate:"max=100"`
	ShippingProvince   string              `json:"shipping_province" validate:"max=100"`
	ShippingPostalCode string              `json:"shipping_postal_code" validate:"max=10"`
	PaymentMethod      string              `json:"payment_method" validate:"max=30"`
	Courier            string              `json:"courier" validate:"max=50"`
	DiscountAmount     int64               `json:"discount_amount" validate:"min=0"`
	ShippingCost       int64               `json:"shipping_cost" validate:"min=0"`
	Notes              string              `json:"notes"`
}

type UpdateRequest struct {
	CustomerID         uuid.UUID `json:"customer_id" validate:"required"`
	RecipientName      string    `json:"recipient_name" validate:"required,max=100"`
	RecipientPhone     string    `json:"recipient_phone" validate:"required,max=30"`
	ShippingAddress    string    `json:"shipping_address" validate:"required"`
	ShippingCity       string    `json:"shipping_city" validate:"max=100"`
	ShippingProvince   string    `json:"shipping_province" validate:"max=100"`
	ShippingPostalCode string    `json:"shipping_postal_code" validate:"max=10"`
	PaymentMethod      string    `json:"payment_method" validate:"max=30"`
	Courier            string    `json:"courier" validate:"max=50"`
	DiscountAmount     int64     `json:"discount_amount" validate:"min=0"`
	ShippingCost       int64     `json:"shipping_cost" validate:"min=0"`
	Notes              string    `json:"notes"`
}

type UpdateItemsRequest struct {
	Items []CreateItemRequest `json:"items" validate:"required,min=1,dive"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=pending confirmed processing shipped completed cancelled"`
}

type UpdatePaymentRequest struct {
	PaymentStatus string `json:"payment_status" validate:"required,oneof=unpaid paid refunded"`
	PaymentMethod string `json:"payment_method" validate:"max=30"`
}

type UpdateShippingRequest struct {
	ShippingStatus string `json:"shipping_status" validate:"required,oneof=pending picked_up in_transit delivered returned"`
	Courier        string `json:"courier" validate:"max=50"`
	TrackingNumber string `json:"tracking_number" validate:"max=100"`
}

type ListFilter struct {
	Page           int
	PerPage        int
	Search         string
	Status         string
	PaymentStatus  string
	ShippingStatus string
	CustomerID     *uuid.UUID
	From           *time.Time
	To             *time.Time
}

type ItemResponse struct {
	ID           uuid.UUID  `json:"id"`
	ProductID    *uuid.UUID `json:"product_id"`
	VariantID    *uuid.UUID `json:"variant_id"`
	ProductName  string     `json:"product_name"`
	VariantLabel string     `json:"variant_label"`
	SKU          string     `json:"sku"`
	Price        int64      `json:"price"`
	Quantity     int        `json:"quantity"`
	Subtotal     int64      `json:"subtotal"`
}

type Response struct {
	ID                 uuid.UUID          `json:"id"`
	OrderNumber        string             `json:"order_number"`
	CustomerID         uuid.UUID          `json:"customer_id"`
	Customer           *customer.Response `json:"customer"`
	Status             string             `json:"status"`
	PaymentStatus      string             `json:"payment_status"`
	PaymentMethod      string             `json:"payment_method"`
	PaidAt             *time.Time         `json:"paid_at"`
	ShippingStatus     string             `json:"shipping_status"`
	Courier            string             `json:"courier"`
	TrackingNumber     string             `json:"tracking_number"`
	ShippedAt          *time.Time         `json:"shipped_at"`
	DeliveredAt        *time.Time         `json:"delivered_at"`
	RecipientName      string             `json:"recipient_name"`
	RecipientPhone     string             `json:"recipient_phone"`
	ShippingAddress    string             `json:"shipping_address"`
	ShippingCity       string             `json:"shipping_city"`
	ShippingProvince   string             `json:"shipping_province"`
	ShippingPostalCode string             `json:"shipping_postal_code"`
	Subtotal           int64              `json:"subtotal"`
	DiscountAmount     int64              `json:"discount_amount"`
	ShippingCost       int64              `json:"shipping_cost"`
	TotalAmount        int64              `json:"total_amount"`
	Notes              string             `json:"notes"`
	StockDeducted      bool               `json:"stock_deducted"`
	CreatedByID        *uuid.UUID         `json:"created_by_id"`
	Items              []ItemResponse     `json:"items"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

func ToResponse(order *models.Order) Response {
	items := make([]ItemResponse, 0, len(order.Items))
	for index := range order.Items {
		item := order.Items[index]
		items = append(items, ItemResponse{
			ID:           item.ID,
			ProductID:    item.ProductID,
			VariantID:    item.VariantID,
			ProductName:  item.ProductName,
			VariantLabel: item.VariantLabel,
			SKU:          item.SKU,
			Price:        item.Price,
			Quantity:     item.Quantity,
			Subtotal:     item.Subtotal,
		})
	}

	var customerResponse *customer.Response
	if order.Customer.ID != uuid.Nil {
		value := customer.ToResponse(&order.Customer)
		customerResponse = &value
	}

	return Response{
		ID:                 order.ID,
		OrderNumber:        order.OrderNumber,
		CustomerID:         order.CustomerID,
		Customer:           customerResponse,
		Status:             order.Status,
		PaymentStatus:      order.PaymentStatus,
		PaymentMethod:      order.PaymentMethod,
		PaidAt:             order.PaidAt,
		ShippingStatus:     order.ShippingStatus,
		Courier:            order.Courier,
		TrackingNumber:     order.TrackingNumber,
		ShippedAt:          order.ShippedAt,
		DeliveredAt:        order.DeliveredAt,
		RecipientName:      order.RecipientName,
		RecipientPhone:     order.RecipientPhone,
		ShippingAddress:    order.ShippingAddress,
		ShippingCity:       order.ShippingCity,
		ShippingProvince:   order.ShippingProvince,
		ShippingPostalCode: order.ShippingPostalCode,
		Subtotal:           order.Subtotal,
		DiscountAmount:     order.DiscountAmount,
		ShippingCost:       order.ShippingCost,
		TotalAmount:        order.TotalAmount,
		Notes:              order.Notes,
		StockDeducted:      order.StockDeducted,
		CreatedByID:        order.CreatedByID,
		Items:              items,
		CreatedAt:          order.CreatedAt,
		UpdatedAt:          order.UpdatedAt,
	}
}

func ToResponses(orders []models.Order) []Response {
	result := make([]Response, 0, len(orders))
	for index := range orders {
		result = append(result, ToResponse(&orders[index]))
	}
	return result
}
