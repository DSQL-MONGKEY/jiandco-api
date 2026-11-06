package product

import (
	"time"

	"github.com/google/uuid"

	"jianco-api/internal/models"
	"jianco-api/internal/storage"
)

type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"omitempty,max=120"`
	Description string `json:"description"`
}

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SizeRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=20"`
	SortOrder int    `json:"sort_order"`
}

type SizeResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ColorRequest struct {
	Name    string `json:"name" validate:"required,min=1,max=50"`
	HexCode string `json:"hex_code" validate:"omitempty,max=7"`
}

type ColorResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	HexCode   string    `json:"hex_code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateVariantRequest struct {
	SizeID   uuid.UUID `json:"size_id" validate:"required"`
	ColorID  uuid.UUID `json:"color_id" validate:"required"`
	SKU      string    `json:"sku" validate:"required,max=64"`
	Price    int64     `json:"price" validate:"min=0"`
	Stock    int       `json:"stock" validate:"min=0"`
	IsActive *bool     `json:"is_active"`
}

type UpdateVariantRequest struct {
	SizeID   uuid.UUID `json:"size_id" validate:"required"`
	ColorID  uuid.UUID `json:"color_id" validate:"required"`
	SKU      string    `json:"sku" validate:"required,max=64"`
	Price    int64     `json:"price" validate:"min=0"`
	IsActive *bool     `json:"is_active"`
}

type UpdateStockRequest struct {
	Stock *int   `json:"stock" validate:"required,min=0"`
	Note  string `json:"note" validate:"max=255"`
}

type CreateProductRequest struct {
	CategoryID  *uuid.UUID             `json:"category_id"`
	Name        string                 `json:"name" validate:"required,min=2,max=150"`
	Slug        string                 `json:"slug" validate:"omitempty,max=180"`
	Description string                 `json:"description"`
	Status      string                 `json:"status" validate:"omitempty,oneof=draft active archived"`
	ShopeeURL   string                 `json:"shopee_url" validate:"omitempty,url,max=255"`
	ShopeePrice int64                  `json:"shopee_price" validate:"min=0"`
	Details     []string               `json:"details"`
	Badge       string                 `json:"badge" validate:"max=50"`
	BadgeTone   string                 `json:"badge_tone" validate:"omitempty,oneof=light accent"`
	SortOrder   int                    `json:"sort_order"`
	Variants    []CreateVariantRequest `json:"variants" validate:"omitempty,dive"`
}

type UpdateProductRequest struct {
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name" validate:"required,min=2,max=150"`
	Slug        string     `json:"slug" validate:"omitempty,max=180"`
	Description string     `json:"description"`
	Status      string     `json:"status" validate:"required,oneof=draft active archived"`
	ShopeeURL   string     `json:"shopee_url" validate:"omitempty,url,max=255"`
	ShopeePrice int64      `json:"shopee_price" validate:"min=0"`
	Details     []string   `json:"details"`
	Badge       string     `json:"badge" validate:"max=50"`
	BadgeTone   string     `json:"badge_tone" validate:"omitempty,oneof=light accent"`
	SortOrder   int        `json:"sort_order"`
}

type ProductListFilter struct {
	Page       int
	PerPage    int
	Search     string
	CategoryID *uuid.UUID
	Status     string
}

type ImageResponse struct {
	ID        uuid.UUID `json:"id"`
	ObjectKey string    `json:"object_key"`
	URL       string    `json:"url"`
	AltText   string    `json:"alt_text"`
	SortOrder int       `json:"sort_order"`
	IsPrimary bool      `json:"is_primary"`
	CreatedAt time.Time `json:"created_at"`
}

type VariantResponse struct {
	ID        uuid.UUID     `json:"id"`
	ProductID uuid.UUID     `json:"product_id"`
	Size      SizeResponse  `json:"size"`
	Color     ColorResponse `json:"color"`
	SKU       string        `json:"sku"`
	Price     int64         `json:"price"`
	Stock     int           `json:"stock"`
	IsActive  bool          `json:"is_active"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type StockMovementResponse struct {
	ID            uuid.UUID  `json:"id"`
	VariantID     *uuid.UUID `json:"variant_id"`
	ProductID     *uuid.UUID `json:"product_id"`
	SKU           string     `json:"sku"`
	VariantLabel  string     `json:"variant_label"`
	Type          string     `json:"type"`
	Quantity      int        `json:"quantity"`
	StockBefore   int        `json:"stock_before"`
	StockAfter    int        `json:"stock_after"`
	ReferenceType string     `json:"reference_type"`
	ReferenceID   *uuid.UUID `json:"reference_id"`
	Note          string     `json:"note"`
	CreatedByID   *uuid.UUID `json:"created_by_id"`
	CreatedByName string     `json:"created_by_name"`
	CreatedAt     time.Time  `json:"created_at"`
}

type ProductResponse struct {
	ID          uuid.UUID         `json:"id"`
	CategoryID  *uuid.UUID        `json:"category_id"`
	Category    *CategoryResponse `json:"category"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	ShopeeURL   string            `json:"shopee_url"`
	ShopeePrice int64             `json:"shopee_price"`
	Details     []string          `json:"details"`
	Badge       string            `json:"badge"`
	BadgeTone   string            `json:"badge_tone"`
	SortOrder   int               `json:"sort_order"`
	CreatedByID *uuid.UUID        `json:"created_by_id"`
	Images      []ImageResponse   `json:"images"`
	Variants    []VariantResponse `json:"variants"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func ToCategoryResponse(category *models.Category) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func ToCategoryResponses(categories []models.Category) []CategoryResponse {
	result := make([]CategoryResponse, 0, len(categories))
	for index := range categories {
		result = append(result, ToCategoryResponse(&categories[index]))
	}
	return result
}

func ToSizeResponse(size *models.Size) SizeResponse {
	return SizeResponse{
		ID:        size.ID,
		Name:      size.Name,
		SortOrder: size.SortOrder,
		CreatedAt: size.CreatedAt,
		UpdatedAt: size.UpdatedAt,
	}
}

func ToSizeResponses(sizes []models.Size) []SizeResponse {
	result := make([]SizeResponse, 0, len(sizes))
	for index := range sizes {
		result = append(result, ToSizeResponse(&sizes[index]))
	}
	return result
}

func ToColorResponse(color *models.Color) ColorResponse {
	return ColorResponse{
		ID:        color.ID,
		Name:      color.Name,
		HexCode:   color.HexCode,
		CreatedAt: color.CreatedAt,
		UpdatedAt: color.UpdatedAt,
	}
}

func ToColorResponses(colors []models.Color) []ColorResponse {
	result := make([]ColorResponse, 0, len(colors))
	for index := range colors {
		result = append(result, ToColorResponse(&colors[index]))
	}
	return result
}

func ToImageResponse(image *models.ProductImage, storageClient *storage.Client) ImageResponse {
	return ImageResponse{
		ID:        image.ID,
		ObjectKey: image.ObjectKey,
		URL:       storageClient.URL(image.ObjectKey),
		AltText:   image.AltText,
		SortOrder: image.SortOrder,
		IsPrimary: image.IsPrimary,
		CreatedAt: image.CreatedAt,
	}
}

func ToVariantResponse(variant *models.ProductVariant) VariantResponse {
	return VariantResponse{
		ID:        variant.ID,
		ProductID: variant.ProductID,
		Size:      ToSizeResponse(&variant.Size),
		Color:     ToColorResponse(&variant.Color),
		SKU:       variant.SKU,
		Price:     variant.Price,
		Stock:     variant.Stock,
		IsActive:  variant.IsActive,
		CreatedAt: variant.CreatedAt,
		UpdatedAt: variant.UpdatedAt,
	}
}

func ToStockMovementResponse(movement *models.StockMovement) StockMovementResponse {
	createdByName := ""
	if movement.CreatedBy != nil {
		createdByName = movement.CreatedBy.Name
	}
	return StockMovementResponse{
		ID:            movement.ID,
		VariantID:     movement.VariantID,
		ProductID:     movement.ProductID,
		SKU:           movement.SKU,
		VariantLabel:  movement.VariantLabel,
		Type:          movement.Type,
		Quantity:      movement.Quantity,
		StockBefore:   movement.StockBefore,
		StockAfter:    movement.StockAfter,
		ReferenceType: movement.ReferenceType,
		ReferenceID:   movement.ReferenceID,
		Note:          movement.Note,
		CreatedByID:   movement.CreatedByID,
		CreatedByName: createdByName,
		CreatedAt:     movement.CreatedAt,
	}
}

func ToStockMovementResponses(movements []models.StockMovement) []StockMovementResponse {
	result := make([]StockMovementResponse, 0, len(movements))
	for index := range movements {
		result = append(result, ToStockMovementResponse(&movements[index]))
	}
	return result
}

func ToProductResponse(product *models.Product, storageClient *storage.Client) ProductResponse {
	images := make([]ImageResponse, 0, len(product.Images))
	for index := range product.Images {
		images = append(images, ToImageResponse(&product.Images[index], storageClient))
	}

	variants := make([]VariantResponse, 0, len(product.Variants))
	for index := range product.Variants {
		variants = append(variants, ToVariantResponse(&product.Variants[index]))
	}

	var category *CategoryResponse
	if product.Category != nil {
		value := ToCategoryResponse(product.Category)
		category = &value
	}

	return ProductResponse{
		ID:          product.ID,
		CategoryID:  product.CategoryID,
		Category:    category,
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Status:      product.Status,
		ShopeeURL:   product.ShopeeURL,
		ShopeePrice: product.ShopeePrice,
		Details:     product.Details,
		Badge:       product.Badge,
		BadgeTone:   product.BadgeTone,
		SortOrder:   product.SortOrder,
		CreatedByID: product.CreatedByID,
		Images:      images,
		Variants:    variants,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
