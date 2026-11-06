package models

import (
	"time"

	"github.com/google/uuid"
)

type StockMovement struct {
	ID            uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	VariantID     *uuid.UUID      `gorm:"type:uuid;index:idx_stock_movements_variant_id"`
	Variant       *ProductVariant `gorm:"foreignKey:VariantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProductID     *uuid.UUID      `gorm:"type:uuid;index:idx_stock_movements_product_id"`
	Product       *Product        `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	SKU           string          `gorm:"type:varchar(64);not null"`
	VariantLabel  string          `gorm:"type:varchar(100);not null"`
	Type          string          `gorm:"type:varchar(30);not null;index:idx_stock_movements_type;check:chk_stock_movements_type,type IN ('initial','adjustment','order','cancel','return')"`
	Quantity      int             `gorm:"type:integer;not null;check:chk_stock_movements_quantity,quantity <> 0"`
	StockBefore   int             `gorm:"type:integer;not null;check:chk_stock_movements_stock_before,stock_before >= 0"`
	StockAfter    int             `gorm:"type:integer;not null;check:chk_stock_movements_stock_after,stock_after >= 0"`
	ReferenceType string          `gorm:"type:varchar(30);index:idx_stock_movements_reference"`
	ReferenceID   *uuid.UUID      `gorm:"type:uuid;index:idx_stock_movements_reference"`
	Note          string          `gorm:"type:text"`
	CreatedByID   *uuid.UUID      `gorm:"type:uuid;index:idx_stock_movements_created_by"`
	CreatedBy     *User           `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt     time.Time       `gorm:"type:timestamptz;not null;default:now()"`
}
