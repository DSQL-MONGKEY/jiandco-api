package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	Base
	Name        string `gorm:"type:varchar(100);not null;uniqueIndex:idx_categories_name"`
	Slug        string `gorm:"type:varchar(120);not null;uniqueIndex:idx_categories_slug"`
	Description string `gorm:"type:text"`
	SoftDelete
}

type Size struct {
	Base
	Name      string `gorm:"type:varchar(20);not null;uniqueIndex:idx_sizes_name"`
	SortOrder int    `gorm:"type:integer;not null;default:0"`
}

type Color struct {
	Base
	Name    string `gorm:"type:varchar(50);not null;uniqueIndex:idx_colors_name"`
	HexCode string `gorm:"type:varchar(7)"`
}

type ProductColor struct {
	Base
	ProductID uuid.UUID `gorm:"type:uuid;not null;index:idx_product_colors_product_id;uniqueIndex:idx_product_colors_product_name"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Name      string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_product_colors_product_name"`
	HexCode   string    `gorm:"type:varchar(7)"`
	ImageKey  string    `gorm:"type:varchar(255)"`
	SortOrder int       `gorm:"type:integer;not null;default:0;index:idx_product_colors_sort_order"`
}

type Product struct {
	Base
	CategoryID  *uuid.UUID       `gorm:"type:uuid;index:idx_products_category_id"`
	Category    *Category        `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name        string           `gorm:"type:varchar(150);not null;index:idx_products_name"`
	Slug        string           `gorm:"type:varchar(180);not null;uniqueIndex:idx_products_slug"`
	Description string           `gorm:"type:text"`
	Status      string           `gorm:"type:varchar(20);not null;default:'draft';index:idx_products_status;check:chk_products_status,status IN ('draft','active','archived')"`
	ShopeeURL   string           `gorm:"type:varchar(255)"`
	ShopeePrice int64            `gorm:"type:bigint;not null;default:0;check:chk_products_shopee_price,shopee_price >= 0"`
	Details     []string         `gorm:"type:jsonb;serializer:json"`
	Badge       string           `gorm:"type:varchar(50)"`
	BadgeTone   string           `gorm:"type:varchar(20);not null;default:'light';check:chk_products_badge_tone,badge_tone IN ('light','accent')"`
	SortOrder   int              `gorm:"type:integer;not null;default:0;index:idx_products_sort_order"`
	CreatedByID *uuid.UUID       `gorm:"type:uuid;index:idx_products_created_by"`
	CreatedBy   *User            `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Images      []ProductImage   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Colors      []ProductColor   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Variants    []ProductVariant `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SoftDelete
}

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProductID uuid.UUID `gorm:"type:uuid;not null;index:idx_product_images_product_id"`
	Product   Product   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ObjectKey string    `gorm:"type:varchar(255);not null"`
	AltText   string    `gorm:"type:varchar(150)"`
	SortOrder int       `gorm:"type:integer;not null;default:0"`
	IsPrimary bool      `gorm:"not null;default:false"`
	CreatedAt time.Time `gorm:"type:timestamptz;not null;default:now()"`
}

type ProductVariant struct {
	Base
	ProductID      uuid.UUID     `gorm:"type:uuid;not null;index:idx_product_variants_product_id;uniqueIndex:idx_product_variants_combo;uniqueIndex:idx_product_variants_product_color_combo"`
	Product        Product       `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SizeID         uuid.UUID     `gorm:"type:uuid;not null;index:idx_product_variants_size_id;uniqueIndex:idx_product_variants_combo;uniqueIndex:idx_product_variants_product_color_combo"`
	Size           Size          `gorm:"foreignKey:SizeID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ColorID        *uuid.UUID    `gorm:"type:uuid;index:idx_product_variants_color_id;uniqueIndex:idx_product_variants_combo"`
	Color          *Color        `gorm:"foreignKey:ColorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	ProductColorID *uuid.UUID    `gorm:"type:uuid;index:idx_product_variants_product_color_id;uniqueIndex:idx_product_variants_product_color_combo"`
	ProductColor   *ProductColor `gorm:"foreignKey:ProductColorID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	SKU            string        `gorm:"type:varchar(64);not null;uniqueIndex:idx_product_variants_sku"`
	Price          int64         `gorm:"type:bigint;not null;default:0;check:chk_product_variants_price,price >= 0"`
	Stock          int           `gorm:"type:integer;not null;default:0;check:chk_product_variants_stock,stock >= 0"`
	IsActive       bool          `gorm:"not null;default:true"`
}
