-- Create "product_colors" table
CREATE TABLE "product_colors" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "product_id" uuid NOT NULL,
  "name" character varying(50) NOT NULL,
  "hex_code" character varying(7) NULL,
  "image_key" character varying(255) NULL,
  "sort_order" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_colors_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create indexes for product colors
CREATE INDEX "idx_product_colors_product_id" ON "product_colors" ("product_id");
CREATE UNIQUE INDEX "idx_product_colors_product_name" ON "product_colors" ("product_id", "name");
CREATE INDEX "idx_product_colors_sort_order" ON "product_colors" ("sort_order");

-- Allow legacy variants to coexist while product colors are introduced.
ALTER TABLE "product_variants" ALTER COLUMN "color_id" DROP NOT NULL;
ALTER TABLE "product_variants" ADD COLUMN "product_color_id" uuid NULL;
CREATE INDEX "idx_product_variants_product_color_id" ON "product_variants" ("product_color_id");
CREATE UNIQUE INDEX "idx_product_variants_product_color_combo" ON "product_variants" ("product_id", "size_id", "product_color_id");
ALTER TABLE "product_variants" ADD CONSTRAINT "fk_product_variants_product_color" FOREIGN KEY ("product_color_id") REFERENCES "product_colors" ("id") ON UPDATE CASCADE ON DELETE RESTRICT;

-- Backfill existing product/color combinations.
INSERT INTO "product_colors" ("id", "product_id", "name", "hex_code", "sort_order")
SELECT gen_random_uuid(), pv."product_id", c."name", c."hex_code", 0
FROM "product_variants" pv
JOIN "colors" c ON c."id" = pv."color_id"
GROUP BY pv."product_id", c."name", c."hex_code";

UPDATE "product_variants" pv
SET "product_color_id" = pc."id"
FROM "product_colors" pc
JOIN "colors" c ON c."id" = pv."color_id"
WHERE pc."product_id" = pv."product_id"
  AND pc."name" = c."name";
