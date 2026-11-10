package catalog

type PublicImage struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

type PublicCollection struct {
	ID          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type PublicColor struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	HexCode string      `json:"hexCode"`
	Image   PublicImage `json:"image"`
	Sizes   []string    `json:"sizes"`
}

type PublicProduct struct {
	ID            string        `json:"id"`
	Slug          string        `json:"slug"`
	Name          string        `json:"name"`
	Collection    string        `json:"collection"`
	Color         string        `json:"color"`
	Colors        []PublicColor `json:"colors"`
	WhatsAppPrice int64         `json:"whatsappPrice"`
	ShopeePrice   int64         `json:"shopeePrice"`
	Images        []PublicImage `json:"images"`
	Description   string        `json:"description"`
	Details       []string      `json:"details"`
	Sizes         []string      `json:"sizes"`
	ShopeeURL     string        `json:"shopeeUrl"`
	Badge         string        `json:"badge,omitempty"`
	BadgeTone     string        `json:"badgeTone,omitempty"`
}

type PublicPromo struct {
	Badge       string `json:"badge"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Note        string `json:"note"`
	CtaLabel    string `json:"ctaLabel"`
}

type PublicPage struct {
	Eyebrow     string `json:"eyebrow"`
	Title       string `json:"title"`
	TitleAccent string `json:"titleAccent"`
	Description string `json:"description"`
	Footnote    string `json:"footnote"`
}

type PublicCTA struct {
	Eyebrow           string `json:"eyebrow"`
	Title             string `json:"title"`
	TitleAccent       string `json:"titleAccent"`
	Description       string `json:"description"`
	Label             string `json:"label"`
	ImageURL          string `json:"imageUrl"`
	SecondaryImageURL string `json:"secondaryImageUrl"`
}

type PublicCatalog struct {
	Collections    []PublicCollection `json:"collections"`
	Products       []PublicProduct    `json:"products"`
	Promo          PublicPromo        `json:"promo"`
	Page           PublicPage         `json:"page"`
	CTA            PublicCTA          `json:"cta"`
	ShopeeStoreURL string             `json:"shopeeStoreUrl"`
}

type PublicProductResponse struct {
	Product        PublicProduct      `json:"product"`
	Collections    []PublicCollection `json:"collections"`
	Promo          PublicPromo        `json:"promo"`
	Page           PublicPage         `json:"page"`
	CTA            PublicCTA          `json:"cta"`
	ShopeeStoreURL string             `json:"shopeeStoreUrl"`
}

type SettingsRequest struct {
	PromoBadge           string `json:"promo_badge" validate:"max=100"`
	PromoTitle           string `json:"promo_title" validate:"max=150"`
	PromoDescription     string `json:"promo_description"`
	PromoNote            string `json:"promo_note"`
	PromoCtaLabel        string `json:"promo_cta_label" validate:"max=100"`
	PageEyebrow          string `json:"page_eyebrow" validate:"max=100"`
	PageTitle            string `json:"page_title" validate:"max=150"`
	PageTitleAccent      string `json:"page_title_accent" validate:"max=150"`
	PageDescription      string `json:"page_description"`
	PageFootnote         string `json:"page_footnote"`
	CtaEyebrow           string `json:"cta_eyebrow" validate:"max=100"`
	CtaTitle             string `json:"cta_title" validate:"max=150"`
	CtaTitleAccent       string `json:"cta_title_accent" validate:"max=150"`
	CtaDescription       string `json:"cta_description"`
	CtaLabel             string `json:"cta_label" validate:"max=100"`
	CtaImageURL          string `json:"cta_image_url" validate:"max=255"`
	CtaSecondaryImageURL string `json:"cta_secondary_image_url" validate:"max=255"`
	ShopeeStoreURL       string `json:"shopee_store_url" validate:"omitempty,url,max=255"`
}

type SettingsResponse struct {
	PromoBadge           string `json:"promo_badge"`
	PromoTitle           string `json:"promo_title"`
	PromoDescription     string `json:"promo_description"`
	PromoNote            string `json:"promo_note"`
	PromoCtaLabel        string `json:"promo_cta_label"`
	PageEyebrow          string `json:"page_eyebrow"`
	PageTitle            string `json:"page_title"`
	PageTitleAccent      string `json:"page_title_accent"`
	PageDescription      string `json:"page_description"`
	PageFootnote         string `json:"page_footnote"`
	CtaEyebrow           string `json:"cta_eyebrow"`
	CtaTitle             string `json:"cta_title"`
	CtaTitleAccent       string `json:"cta_title_accent"`
	CtaDescription       string `json:"cta_description"`
	CtaLabel             string `json:"cta_label"`
	CtaImageURL          string `json:"cta_image_url"`
	CtaSecondaryImageURL string `json:"cta_secondary_image_url"`
	ShopeeStoreURL       string `json:"shopee_store_url"`
}
