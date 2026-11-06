-- Create "stock_movements" table
CREATE TABLE "stock_movements" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "variant_id" uuid NULL,
  "product_id" uuid NULL,
  "sku" character varying(64) NOT NULL,
  "variant_label" character varying(100) NOT NULL,
  "type" character varying(30) NOT NULL,
  "quantity" integer NOT NULL,
  "stock_before" integer NOT NULL,
  "stock_after" integer NOT NULL,
  "reference_type" character varying(30) NULL,
  "reference_id" uuid NULL,
  "note" text NULL,
  "created_by_id" uuid NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_stock_movements_created_by" FOREIGN KEY ("created_by_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_stock_movements_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_stock_movements_variant" FOREIGN KEY ("variant_id") REFERENCES "product_variants" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "chk_stock_movements_quantity" CHECK (quantity <> 0),
  CONSTRAINT "chk_stock_movements_stock_after" CHECK (stock_after >= 0),
  CONSTRAINT "chk_stock_movements_stock_before" CHECK (stock_before >= 0),
  CONSTRAINT "chk_stock_movements_type" CHECK ((type)::text = ANY ((ARRAY['initial'::character varying, 'adjustment'::character varying, 'order'::character varying, 'cancel'::character varying, 'return'::character varying])::text[]))
);
-- Create index "idx_stock_movements_created_by" to table: "stock_movements"
CREATE INDEX "idx_stock_movements_created_by" ON "stock_movements" ("created_by_id");
-- Create index "idx_stock_movements_product_id" to table: "stock_movements"
CREATE INDEX "idx_stock_movements_product_id" ON "stock_movements" ("product_id");
-- Create index "idx_stock_movements_reference" to table: "stock_movements"
CREATE INDEX "idx_stock_movements_reference" ON "stock_movements" ("reference_type", "reference_id");
-- Create index "idx_stock_movements_type" to table: "stock_movements"
CREATE INDEX "idx_stock_movements_type" ON "stock_movements" ("type");
-- Create index "idx_stock_movements_variant_id" to table: "stock_movements"
CREATE INDEX "idx_stock_movements_variant_id" ON "stock_movements" ("variant_id");
