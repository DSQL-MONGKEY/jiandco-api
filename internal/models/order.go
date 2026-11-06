package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	Base
	Name       string `gorm:"type:varchar(100);not null;index:idx_customers_name"`
	Phone      string `gorm:"type:varchar(30);not null;index:idx_customers_phone"`
	Email      string `gorm:"type:varchar(255);index:idx_customers_email"`
	Address    string `gorm:"type:text"`
	City       string `gorm:"type:varchar(100)"`
	Province   string `gorm:"type:varchar(100)"`
	PostalCode string `gorm:"type:varchar(10)"`
	Notes      string `gorm:"type:text"`
	SoftDelete
}

type Order struct {
	Base
	OrderNumber        string      `gorm:"type:varchar(32);not null;uniqueIndex:idx_orders_order_number"`
	CustomerID         uuid.UUID   `gorm:"type:uuid;not null;index:idx_orders_customer_id"`
	Customer           Customer    `gorm:"foreignKey:CustomerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Status             string      `gorm:"type:varchar(20);not null;default:'pending';index:idx_orders_status;check:chk_orders_status,status IN ('pending','confirmed','processing','shipped','completed','cancelled')"`
	PaymentStatus      string      `gorm:"type:varchar(20);not null;default:'unpaid';index:idx_orders_payment_status;check:chk_orders_payment_status,payment_status IN ('unpaid','paid','refunded')"`
	PaymentMethod      string      `gorm:"type:varchar(30)"`
	PaidAt             *time.Time  `gorm:"type:timestamptz;index:idx_orders_paid_at"`
	ShippingStatus     string      `gorm:"type:varchar(20);not null;default:'pending';index:idx_orders_shipping_status;check:chk_orders_shipping_status,shipping_status IN ('pending','picked_up','in_transit','delivered','returned')"`
	Courier            string      `gorm:"type:varchar(50)"`
	TrackingNumber     string      `gorm:"type:varchar(100);index:idx_orders_tracking_number"`
	ShippedAt          *time.Time  `gorm:"type:timestamptz"`
	DeliveredAt        *time.Time  `gorm:"type:timestamptz"`
	RecipientName      string      `gorm:"type:varchar(100);not null"`
	RecipientPhone     string      `gorm:"type:varchar(30);not null"`
	ShippingAddress    string      `gorm:"type:text;not null"`
	ShippingCity       string      `gorm:"type:varchar(100)"`
	ShippingProvince   string      `gorm:"type:varchar(100)"`
	ShippingPostalCode string      `gorm:"type:varchar(10)"`
	Subtotal           int64       `gorm:"type:bigint;not null;default:0;check:chk_orders_subtotal,subtotal >= 0"`
	DiscountAmount     int64       `gorm:"type:bigint;not null;default:0;check:chk_orders_discount_amount,discount_amount >= 0"`
	ShippingCost       int64       `gorm:"type:bigint;not null;default:0;check:chk_orders_shipping_cost,shipping_cost >= 0"`
	TotalAmount        int64       `gorm:"type:bigint;not null;default:0;check:chk_orders_total_amount,total_amount >= 0"`
	Notes              string      `gorm:"type:text"`
	CreatedByID        *uuid.UUID  `gorm:"type:uuid;index:idx_orders_created_by"`
	CreatedBy          *User       `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	StockDeducted      bool        `gorm:"not null;default:false"`
	Items              []OrderItem `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SoftDelete
}

type OrderItem struct {
	ID           uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	OrderID      uuid.UUID       `gorm:"type:uuid;not null;index:idx_order_items_order_id"`
	Order        Order           `gorm:"foreignKey:OrderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ProductID    *uuid.UUID      `gorm:"type:uuid;index:idx_order_items_product_id"`
	Product      *Product        `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	VariantID    *uuid.UUID      `gorm:"type:uuid;index:idx_order_items_variant_id"`
	Variant      *ProductVariant `gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProductName  string          `gorm:"type:varchar(150);not null"`
	VariantLabel string          `gorm:"type:varchar(100);not null"`
	SKU          string          `gorm:"type:varchar(64);not null"`
	Price        int64           `gorm:"type:bigint;not null;default:0;check:chk_order_items_price,price >= 0"`
	Quantity     int             `gorm:"type:integer;not null;default:1;check:chk_order_items_quantity,quantity > 0"`
	Subtotal     int64           `gorm:"type:bigint;not null;default:0;check:chk_order_items_subtotal,subtotal >= 0"`
	CreatedAt    time.Time       `gorm:"type:timestamptz;not null;default:now()"`
}
