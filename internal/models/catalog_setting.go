package models

import (
	"time"

	"github.com/google/uuid"
)

type CatalogSetting struct {
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Key                  string    `gorm:"type:varchar(20);not null;uniqueIndex:idx_catalog_settings_key;default:'default'"`
	PromoBadge           string    `gorm:"type:varchar(100)"`
	PromoTitle           string    `gorm:"type:varchar(150)"`
	PromoDescription     string    `gorm:"type:text"`
	PromoNote            string    `gorm:"type:text"`
	PromoCtaLabel        string    `gorm:"type:varchar(100)"`
	PageEyebrow          string    `gorm:"type:varchar(100)"`
	PageTitle            string    `gorm:"type:varchar(150)"`
	PageTitleAccent      string    `gorm:"type:varchar(150)"`
	PageDescription      string    `gorm:"type:text"`
	PageFootnote         string    `gorm:"type:text"`
	CtaEyebrow           string    `gorm:"type:varchar(100)"`
	CtaTitle             string    `gorm:"type:varchar(150)"`
	CtaTitleAccent       string    `gorm:"type:varchar(150)"`
	CtaDescription       string    `gorm:"type:text"`
	CtaLabel             string    `gorm:"type:varchar(100)"`
	CtaImageURL          string    `gorm:"type:varchar(255)"`
	CtaSecondaryImageURL string    `gorm:"type:varchar(255)"`
	ShopeeStoreURL       string    `gorm:"type:varchar(255)"`
	CreatedAt            time.Time `gorm:"type:timestamptz;not null;default:now()"`
	UpdatedAt            time.Time `gorm:"type:timestamptz;not null;default:now()"`
}
