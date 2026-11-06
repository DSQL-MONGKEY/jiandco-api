package catalog

import (
	"testing"

	"jianco-api/internal/models"
)

func TestMapProductPublicShape(t *testing.T) {
	category := &models.Category{Name: "Chongsam", Slug: "chongsam"}
	product := models.Product{
		Name:        "Cocoa Lace Vest",
		Slug:        "cocoa-lace-vest",
		Category:    category,
		ShopeePrice: 0,
		Details:     []string{"Lace vest"},
		Images: []models.ProductImage{
			{ObjectKey: "catalog/products/product/image.jpg"},
		},
		Variants: []models.ProductVariant{
			{Price: 150000, IsActive: true, Size: models.Size{Name: "M", SortOrder: 2}, Color: models.Color{Name: "Cocoa"}},
			{Price: 100000, IsActive: false, Size: models.Size{Name: "S", SortOrder: 1}, Color: models.Color{Name: "Cocoa"}},
			{Price: 120000, IsActive: true, Size: models.Size{Name: "S", SortOrder: 1}, Color: models.Color{Name: "Cocoa"}},
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
}

func TestMapCollections(t *testing.T) {
	got := mapCollections([]models.Category{{Name: "Blush", Slug: "blush", Description: "Soft pieces"}})
	if len(got) != 1 || got[0].ID != "blush" || got[0].Label != "Blush" {
		t.Fatalf("collections = %#v", got)
	}
}
