-- Modify "products" table
ALTER TABLE "products" ADD CONSTRAINT "chk_products_badge_tone" CHECK ((badge_tone)::text = ANY ((ARRAY['light'::character varying, 'accent'::character varying])::text[])), ADD CONSTRAINT "chk_products_shopee_price" CHECK (shopee_price >= 0), ADD COLUMN "shopee_url" character varying(255) NULL, ADD COLUMN "shopee_price" bigint NOT NULL DEFAULT 0, ADD COLUMN "details" jsonb NULL, ADD COLUMN "badge" character varying(50) NULL, ADD COLUMN "badge_tone" character varying(20) NOT NULL DEFAULT 'light', ADD COLUMN "sort_order" integer NOT NULL DEFAULT 0;
-- Create index "idx_products_sort_order" to table: "products"
CREATE INDEX "idx_products_sort_order" ON "products" ("sort_order");
-- Create "catalog_settings" table
CREATE TABLE "catalog_settings" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "key" character varying(20) NOT NULL DEFAULT 'default',
  "promo_badge" character varying(100) NULL,
  "promo_title" character varying(150) NULL,
  "promo_description" text NULL,
  "promo_note" text NULL,
  "promo_cta_label" character varying(100) NULL,
  "page_eyebrow" character varying(100) NULL,
  "page_title" character varying(150) NULL,
  "page_title_accent" character varying(150) NULL,
  "page_description" text NULL,
  "page_footnote" text NULL,
  "cta_eyebrow" character varying(100) NULL,
  "cta_title" character varying(150) NULL,
  "cta_title_accent" character varying(150) NULL,
  "cta_description" text NULL,
  "cta_label" character varying(100) NULL,
  "cta_image_url" character varying(255) NULL,
  "cta_secondary_image_url" character varying(255) NULL,
  "shopee_store_url" character varying(255) NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_catalog_settings_key" to table: "catalog_settings"
CREATE UNIQUE INDEX "idx_catalog_settings_key" ON "catalog_settings" ("key");
