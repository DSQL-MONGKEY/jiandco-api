package catalog

import (
	"context"
	"errors"
	"sort"
	"strings"

	"gorm.io/gorm"

	"jianco-api/internal/models"
	"jianco-api/internal/storage"
	"jianco-api/pkg/apperror"
)

type Service struct {
	repo    *Repository
	storage *storage.Client
}

func NewService(repo *Repository, storageClient *storage.Client) *Service {
	return &Service{repo: repo, storage: storageClient}
}

func (s *Service) PublicCatalog(ctx context.Context) (*PublicCatalog, error) {
	products, err := s.repo.ListActiveProducts(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	collections, err := s.repo.ListActiveCollections(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	setting, err := s.getSettingsOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	publicProducts := make([]PublicProduct, 0, len(products))
	for index := range products {
		publicProducts = append(publicProducts, s.mapProduct(&products[index]))
	}

	return &PublicCatalog{
		Collections:    mapCollections(collections),
		Products:       publicProducts,
		Promo:          mapPromo(setting),
		Page:           mapPage(setting),
		CTA:            mapCTA(setting),
		ShopeeStoreURL: setting.ShopeeStoreURL,
	}, nil
}

func (s *Service) PublicProduct(ctx context.Context, slug string) (*PublicProductResponse, error) {
	product, err := s.repo.FindActiveProductBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("catalog product not found")
		}
		return nil, apperror.Internal(err)
	}
	collections, err := s.repo.ListActiveCollections(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	setting, err := s.getSettingsOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	return &PublicProductResponse{
		Product:        s.mapProduct(product),
		Collections:    mapCollections(collections),
		Promo:          mapPromo(setting),
		Page:           mapPage(setting),
		CTA:            mapCTA(setting),
		ShopeeStoreURL: setting.ShopeeStoreURL,
	}, nil
}

func (s *Service) Settings(ctx context.Context) (*SettingsResponse, error) {
	setting, err := s.getSettingsOrDefault(ctx)
	if err != nil {
		return nil, err
	}
	return mapSettingsResponse(setting), nil
}

func (s *Service) UpdateSettings(ctx context.Context, request SettingsRequest) (*SettingsResponse, error) {
	values := map[string]any{
		"promo_badge":             strings.TrimSpace(request.PromoBadge),
		"promo_title":             strings.TrimSpace(request.PromoTitle),
		"promo_description":       strings.TrimSpace(request.PromoDescription),
		"promo_note":              strings.TrimSpace(request.PromoNote),
		"promo_cta_label":         strings.TrimSpace(request.PromoCtaLabel),
		"page_eyebrow":            strings.TrimSpace(request.PageEyebrow),
		"page_title":              strings.TrimSpace(request.PageTitle),
		"page_title_accent":       strings.TrimSpace(request.PageTitleAccent),
		"page_description":        strings.TrimSpace(request.PageDescription),
		"page_footnote":           strings.TrimSpace(request.PageFootnote),
		"cta_eyebrow":             strings.TrimSpace(request.CtaEyebrow),
		"cta_title":               strings.TrimSpace(request.CtaTitle),
		"cta_title_accent":        strings.TrimSpace(request.CtaTitleAccent),
		"cta_description":         strings.TrimSpace(request.CtaDescription),
		"cta_label":               strings.TrimSpace(request.CtaLabel),
		"cta_image_url":           strings.TrimSpace(request.CtaImageURL),
		"cta_secondary_image_url": strings.TrimSpace(request.CtaSecondaryImageURL),
		"shopee_store_url":        strings.TrimSpace(request.ShopeeStoreURL),
	}

	_, err := s.repo.GetSettings(ctx)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Internal(err)
		}
		setting := &models.CatalogSetting{Key: "default"}
		if err := s.repo.CreateSettings(ctx, setting); err != nil {
			return nil, apperror.Internal(err)
		}
	}

	if err := s.repo.UpdateSettings(ctx, values); err != nil {
		return nil, apperror.Internal(err)
	}

	setting, err := s.repo.GetSettings(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	return mapSettingsResponse(setting), nil
}

func (s *Service) getSettingsOrDefault(ctx context.Context) (*models.CatalogSetting, error) {
	setting, err := s.repo.GetSettings(ctx)
	if err == nil {
		return setting, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperror.Internal(err)
	}
	return DefaultSettings(), nil
}

func (s *Service) mapProduct(product *models.Product) PublicProduct {
	variants := activeVariants(product.Variants)
	whatsappPrice := minPrice(variants, product.Variants)
	shopeePrice := product.ShopeePrice
	if shopeePrice <= 0 {
		shopeePrice = whatsappPrice
	}

	images := make([]PublicImage, 0, len(product.Images))
	for index := range product.Images {
		images = append(images, PublicImage{
			Src: s.publicImageURL(product.Images[index].ObjectKey),
			Alt: firstNonEmpty(product.Images[index].AltText, product.Name),
		})
	}

	collection := ""
	if product.Category != nil {
		collection = product.Category.Slug
	}

	badgeTone := product.BadgeTone
	if badgeTone == "" {
		badgeTone = "light"
	}

	return PublicProduct{
		ID:            product.ID.String(),
		Slug:          product.Slug,
		Name:          product.Name,
		Collection:    collection,
		Color:         variantColor(variants, product.Variants),
		WhatsAppPrice: whatsappPrice,
		ShopeePrice:   shopeePrice,
		Images:        images,
		Description:   product.Description,
		Details:       nonNilStrings(product.Details),
		Sizes:         variantSizes(variants, product.Variants),
		ShopeeURL:     product.ShopeeURL,
		Badge:         product.Badge,
		BadgeTone:     badgeTone,
	}
}

func (s *Service) publicImageURL(key string) string {
	url := s.storage.URL(key)
	if strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://") {
		return url
	}
	return "/api/v1/public/media/" + strings.TrimLeft(key, "/")
}

func activeVariants(variants []models.ProductVariant) []models.ProductVariant {
	result := make([]models.ProductVariant, 0, len(variants))
	for _, variant := range variants {
		if variant.IsActive {
			result = append(result, variant)
		}
	}
	if len(result) == 0 {
		return variants
	}
	return result
}

func minPrice(primary, fallback []models.ProductVariant) int64 {
	if len(primary) == 0 {
		primary = fallback
	}
	if len(primary) == 0 {
		return 0
	}
	price := primary[0].Price
	for _, variant := range primary[1:] {
		if variant.Price < price {
			price = variant.Price
		}
	}
	return price
}

func variantColor(primary, fallback []models.ProductVariant) string {
	if len(primary) == 0 {
		primary = fallback
	}
	if len(primary) == 0 {
		return ""
	}
	return primary[0].Color.Name
}

func variantSizes(primary, fallback []models.ProductVariant) []string {
	if len(primary) == 0 {
		primary = fallback
	}
	seen := make(map[string]struct{}, len(primary))
	sizes := make([]models.Size, 0, len(primary))
	for _, variant := range primary {
		if _, exists := seen[variant.Size.Name]; exists {
			continue
		}
		seen[variant.Size.Name] = struct{}{}
		sizes = append(sizes, variant.Size)
	}
	sort.SliceStable(sizes, func(i, j int) bool {
		if sizes[i].SortOrder == sizes[j].SortOrder {
			return sizes[i].Name < sizes[j].Name
		}
		return sizes[i].SortOrder < sizes[j].SortOrder
	})
	result := make([]string, 0, len(sizes))
	for _, size := range sizes {
		result = append(result, size.Name)
	}
	return result
}

func mapCollections(categories []models.Category) []PublicCollection {
	result := make([]PublicCollection, 0, len(categories))
	for index := range categories {
		result = append(result, PublicCollection{
			ID:          categories[index].Slug,
			Label:       categories[index].Name,
			Description: categories[index].Description,
		})
	}
	return result
}

func mapPromo(setting *models.CatalogSetting) PublicPromo {
	return PublicPromo{
		Badge:       setting.PromoBadge,
		Title:       setting.PromoTitle,
		Description: setting.PromoDescription,
		Note:        setting.PromoNote,
		CtaLabel:    setting.PromoCtaLabel,
	}
}

func mapPage(setting *models.CatalogSetting) PublicPage {
	return PublicPage{
		Eyebrow:     setting.PageEyebrow,
		Title:       setting.PageTitle,
		TitleAccent: setting.PageTitleAccent,
		Description: setting.PageDescription,
		Footnote:    setting.PageFootnote,
	}
}

func mapCTA(setting *models.CatalogSetting) PublicCTA {
	return PublicCTA{
		Eyebrow:           setting.CtaEyebrow,
		Title:             setting.CtaTitle,
		TitleAccent:       setting.CtaTitleAccent,
		Description:       setting.CtaDescription,
		Label:             setting.CtaLabel,
		ImageURL:          setting.CtaImageURL,
		SecondaryImageURL: setting.CtaSecondaryImageURL,
	}
}

func mapSettingsResponse(setting *models.CatalogSetting) *SettingsResponse {
	return &SettingsResponse{
		PromoBadge:           setting.PromoBadge,
		PromoTitle:           setting.PromoTitle,
		PromoDescription:     setting.PromoDescription,
		PromoNote:            setting.PromoNote,
		PromoCtaLabel:        setting.PromoCtaLabel,
		PageEyebrow:          setting.PageEyebrow,
		PageTitle:            setting.PageTitle,
		PageTitleAccent:      setting.PageTitleAccent,
		PageDescription:      setting.PageDescription,
		PageFootnote:         setting.PageFootnote,
		CtaEyebrow:           setting.CtaEyebrow,
		CtaTitle:             setting.CtaTitle,
		CtaTitleAccent:       setting.CtaTitleAccent,
		CtaDescription:       setting.CtaDescription,
		CtaLabel:             setting.CtaLabel,
		CtaImageURL:          setting.CtaImageURL,
		CtaSecondaryImageURL: setting.CtaSecondaryImageURL,
		ShopeeStoreURL:       setting.ShopeeStoreURL,
	}
}

func DefaultSettings() *models.CatalogSetting {
	return &models.CatalogSetting{
		Key:              "default",
		PromoBadge:       "Harga Spesial WhatsApp",
		PromoTitle:       "Lebih hemat lewat WhatsApp.",
		PromoDescription: "Dapatkan harga spesial saat memesan langsung melalui WhatsApp. Checkout Shopee tetap tersedia untuk pembayaran instan via QRIS, ShopeePay, atau PayLater.",
		PromoNote:        "Harga dan ketersediaan dapat berubah. Konfirmasi melalui WhatsApp atau Shopee.",
		PromoCtaLabel:    "Chat Admin",
		PageEyebrow:      "Catalog — Chongsam & Blush",
		PageTitle:        "The",
		PageTitleAccent:  "collection.",
		PageDescription:  "Atasan pilihan yang dibuat untuk dipadukan dengan ritme harimu. Pilih ukuran, masukkan ke keranjang, lalu checkout dengan cara yang paling nyaman untukmu.",
		PageFootnote:     "Harga via WhatsApp adalah harga spesial. Harga Shopee mengikuti harga yang tertera di toko Shopee. Ketersediaan dikonfirmasi saat checkout.",
		CtaEyebrow:       "The Collection",
		CtaTitle:         "Temukan koleksi",
		CtaTitleAccent:   "Ji&Co.",
		CtaDescription:   "Jelajahi Chongsam dan Blush—potongan atasan yang dibuat untuk dipadukan dengan ritme harimu sendiri.",
		CtaLabel:         "Jelajahi Katalog",
		ShopeeStoreURL:   "https://shopee.co.id/",
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func nonNilStrings(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
