package catalog

import (
	"testing"

	"github.com/google/uuid"

	"jianco-api/internal/models"
)

func TestMapProductPublicShape(t *testing.T) {
	category := &models.Category{Name: "Chongsam", Slug: "chongsam"}
	colorID := uuid.New()
	product := models.Product{
		Name:        "Cocoa Lace Vest",
		Slug:        "cocoa-lace-vest",
		Category:    category,
		ShopeePrice: 0,
		Details:     []string{"Lace vest"},
		Images: []models.ProductImage{
			{ObjectKey: "catalog/products/product/image.jpg"},
		},
		Colors: []models.ProductColor{
			{Base: models.Base{ID: colorID}, Name: "Cocoa", HexCode: "#8B5E3C"},
		},
		Variants: []models.ProductVariant{
			{Price: 150000, IsActive: true, Size: models.Size{Name: "M", SortOrder: 2}, ProductColorID: &colorID, ProductColor: &models.ProductColor{Name: "Cocoa"}},
			{Price: 100000, IsActive: false, Size: models.Size{Name: "S", SortOrder: 1}, ProductColorID: &colorID, ProductColor: &models.ProductColor{Name: "Cocoa"}},
			{Price: 120000, IsActive: true, Size: models.Size{Name: "S", SortOrder: 1}, ProductColorID: &colorID, ProductColor: &models.ProductColor{Name: "Cocoa"}},
		},
	}

	got := (&Service{}).mapProduct(&product)

	if got.Collection != "chongsam" {
		t.Fatalf("collection = %q, want chongsam", got.Collection)
	}
	if got.Color != "Cocoa" {
		t.Fatalf("color = %q, want Cocoa", got.Color)
	}
	if got.WhatsAppPrice != 120000 {
		t.Fatalf("whatsapp price = %d, want 120000", got.WhatsAppPrice)
	}
	if got.ShopeePrice != 120000 {
		t.Fatalf("shopee price = %d, want fallback 120000", got.ShopeePrice)
	}
	if len(got.Sizes) != 2 || got.Sizes[0] != "S" || got.Sizes[1] != "M" {
		t.Fatalf("sizes = %#v, want [S M]", got.Sizes)
	}
	if len(got.Images) != 1 || got.Images[0].Src != "/api/v1/public/media/catalog/products/product/image.jpg" {
		t.Fatalf("images = %#v", got.Images)
	}
	if len(got.Colors) != 1 || got.Colors[0].Name != "Cocoa" {
		t.Fatalf("colors = %#v", got.Colors)
	}
	if len(got.Colors[0].Sizes) != 2 || got.Colors[0].Sizes[0] != "S" || got.Colors[0].Sizes[1] != "M" {
		t.Fatalf("color sizes = %#v, want [S M]", got.Colors[0].Sizes)
	}
}

func TestMapCollections(t *testing.T) {
	got := mapCollections([]models.Category{{Name: "Blush", Slug: "blush", Description: "Soft pieces"}})
	if len(got) != 1 || got[0].ID != "blush" || got[0].Label != "Blush" {
		t.Fatalf("collections = %#v", got)
	}
}
