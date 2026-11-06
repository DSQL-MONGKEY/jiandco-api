-- Create "roles" table
CREATE TABLE "roles" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(50) NOT NULL,
  "level" smallint NOT NULL,
  "description" character varying(255) NULL,
  "is_system" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "chk_roles_level" CHECK ((level >= 0) AND (level <= 3))
);
-- Create index "idx_roles_level" to table: "roles"
CREATE UNIQUE INDEX "idx_roles_level" ON "roles" ("level");
-- Create index "idx_roles_name" to table: "roles"
CREATE UNIQUE INDEX "idx_roles_name" ON "roles" ("name");
-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "role_id" uuid NOT NULL,
  "name" character varying(100) NOT NULL,
  "username" character varying(50) NOT NULL,
  "email" character varying(255) NOT NULL,
  "password_hash" character varying(255) NOT NULL,
  "phone" character varying(30) NULL,
  "avatar_key" character varying(255) NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  "last_login_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE CASCADE ON DELETE RESTRICT
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_is_active" to table: "users"
CREATE INDEX "idx_users_is_active" ON "users" ("is_active");
-- Create index "idx_users_role_id" to table: "users"
CREATE INDEX "idx_users_role_id" ON "users" ("role_id");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- Create "auth_sessions" table
CREATE TABLE "auth_sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "refresh_token_hash" character(64) NOT NULL,
  "user_agent" character varying(255) NULL,
  "ip_address" character varying(45) NULL,
  "expires_at" timestamptz NOT NULL,
  "revoked_at" timestamptz NULL,
  "last_used_at" timestamptz NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_auth_sessions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_auth_sessions_expires_at" to table: "auth_sessions"
CREATE INDEX "idx_auth_sessions_expires_at" ON "auth_sessions" ("expires_at");
-- Create index "idx_auth_sessions_refresh_token_hash" to table: "auth_sessions"
CREATE UNIQUE INDEX "idx_auth_sessions_refresh_token_hash" ON "auth_sessions" ("refresh_token_hash");
-- Create index "idx_auth_sessions_user_id" to table: "auth_sessions"
CREATE INDEX "idx_auth_sessions_user_id" ON "auth_sessions" ("user_id");
-- Create "expenses" table
CREATE TABLE "expenses" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "title" character varying(150) NOT NULL,
  "category" character varying(50) NOT NULL DEFAULT 'operational',
  "vendor" character varying(150) NULL,
  "expense_date" date NOT NULL,
  "payment_method" character varying(30) NULL,
  "total_amount" bigint NOT NULL DEFAULT 0,
  "receipt_key" character varying(255) NULL,
  "notes" text NULL,
  "created_by_id" uuid NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_expenses_created_by" FOREIGN KEY ("created_by_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "chk_expenses_total_amount" CHECK (total_amount >= 0)
);
-- Create index "idx_expenses_category" to table: "expenses"
CREATE INDEX "idx_expenses_category" ON "expenses" ("category");
-- Create index "idx_expenses_created_by" to table: "expenses"
CREATE INDEX "idx_expenses_created_by" ON "expenses" ("created_by_id");
-- Create index "idx_expenses_deleted_at" to table: "expenses"
CREATE INDEX "idx_expenses_deleted_at" ON "expenses" ("deleted_at");
-- Create index "idx_expenses_expense_date" to table: "expenses"
CREATE INDEX "idx_expenses_expense_date" ON "expenses" ("expense_date");
-- Create "expense_items" table
CREATE TABLE "expense_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "expense_id" uuid NOT NULL,
  "name" character varying(150) NOT NULL,
  "quantity" numeric(12,2) NOT NULL DEFAULT 1,
  "unit" character varying(20) NULL DEFAULT 'pcs',
  "unit_price" bigint NOT NULL DEFAULT 0,
  "subtotal" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_expense_items_expense" FOREIGN KEY ("expense_id") REFERENCES "expenses" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "chk_expense_items_quantity" CHECK (quantity > (0)::numeric),
  CONSTRAINT "chk_expense_items_subtotal" CHECK (subtotal >= 0),
  CONSTRAINT "chk_expense_items_unit_price" CHECK (unit_price >= 0)
);
-- Create index "idx_expense_items_expense_id" to table: "expense_items"
CREATE INDEX "idx_expense_items_expense_id" ON "expense_items" ("expense_id");
-- Create "customers" table
CREATE TABLE "customers" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "name" character varying(100) NOT NULL,
  "phone" character varying(30) NOT NULL,
  "email" character varying(255) NULL,
  "address" text NULL,
  "city" character varying(100) NULL,
  "province" character varying(100) NULL,
  "postal_code" character varying(10) NULL,
  "notes" text NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_customers_deleted_at" to table: "customers"
CREATE INDEX "idx_customers_deleted_at" ON "customers" ("deleted_at");
-- Create index "idx_customers_email" to table: "customers"
CREATE INDEX "idx_customers_email" ON "customers" ("email");
-- Create index "idx_customers_name" to table: "customers"
CREATE INDEX "idx_customers_name" ON "customers" ("name");
-- Create index "idx_customers_phone" to table: "customers"
CREATE INDEX "idx_customers_phone" ON "customers" ("phone");
-- Create "orders" table
CREATE TABLE "orders" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "order_number" character varying(32) NOT NULL,
  "customer_id" uuid NOT NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "payment_status" character varying(20) NOT NULL DEFAULT 'unpaid',
  "payment_method" character varying(30) NULL,
  "paid_at" timestamptz NULL,
  "shipping_status" character varying(20) NOT NULL DEFAULT 'pending',
  "courier" character varying(50) NULL,
  "tracking_number" character varying(100) NULL,
  "shipped_at" timestamptz NULL,
  "delivered_at" timestamptz NULL,
  "recipient_name" character varying(100) NOT NULL,
  "recipient_phone" character varying(30) NOT NULL,
  "shipping_address" text NOT NULL,
  "shipping_city" character varying(100) NULL,
  "shipping_province" character varying(100) NULL,
  "shipping_postal_code" character varying(10) NULL,
  "subtotal" bigint NOT NULL DEFAULT 0,
  "discount_amount" bigint NOT NULL DEFAULT 0,
  "shipping_cost" bigint NOT NULL DEFAULT 0,
  "total_amount" bigint NOT NULL DEFAULT 0,
  "notes" text NULL,
  "created_by_id" uuid NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_orders_created_by" FOREIGN KEY ("created_by_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_orders_customer" FOREIGN KEY ("customer_id") REFERENCES "customers" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "chk_orders_discount_amount" CHECK (discount_amount >= 0),
  CONSTRAINT "chk_orders_payment_status" CHECK ((payment_status)::text = ANY ((ARRAY['unpaid'::character varying, 'paid'::character varying, 'refunded'::character varying])::text[])),
  CONSTRAINT "chk_orders_shipping_cost" CHECK (shipping_cost >= 0),
  CONSTRAINT "chk_orders_shipping_status" CHECK ((shipping_status)::text = ANY ((ARRAY['pending'::character varying, 'picked_up'::character varying, 'in_transit'::character varying, 'delivered'::character varying, 'returned'::character varying])::text[])),
  CONSTRAINT "chk_orders_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'confirmed'::character varying, 'processing'::character varying, 'shipped'::character varying, 'completed'::character varying, 'cancelled'::character varying])::text[])),
  CONSTRAINT "chk_orders_subtotal" CHECK (subtotal >= 0),
  CONSTRAINT "chk_orders_total_amount" CHECK (total_amount >= 0)
);
-- Create index "idx_orders_created_by" to table: "orders"
CREATE INDEX "idx_orders_created_by" ON "orders" ("created_by_id");
-- Create index "idx_orders_customer_id" to table: "orders"
CREATE INDEX "idx_orders_customer_id" ON "orders" ("customer_id");
-- Create index "idx_orders_deleted_at" to table: "orders"
CREATE INDEX "idx_orders_deleted_at" ON "orders" ("deleted_at");
-- Create index "idx_orders_order_number" to table: "orders"
CREATE UNIQUE INDEX "idx_orders_order_number" ON "orders" ("order_number");
-- Create index "idx_orders_paid_at" to table: "orders"
CREATE INDEX "idx_orders_paid_at" ON "orders" ("paid_at");
-- Create index "idx_orders_payment_status" to table: "orders"
CREATE INDEX "idx_orders_payment_status" ON "orders" ("payment_status");
-- Create index "idx_orders_shipping_status" to table: "orders"
CREATE INDEX "idx_orders_shipping_status" ON "orders" ("shipping_status");
-- Create index "idx_orders_status" to table: "orders"
CREATE INDEX "idx_orders_status" ON "orders" ("status");
-- Create index "idx_orders_tracking_number" to table: "orders"
CREATE INDEX "idx_orders_tracking_number" ON "orders" ("tracking_number");
-- Create "categories" table
CREATE TABLE "categories" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "name" character varying(100) NOT NULL,
  "slug" character varying(120) NOT NULL,
  "description" text NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_categories_deleted_at" to table: "categories"
CREATE INDEX "idx_categories_deleted_at" ON "categories" ("deleted_at");
-- Create index "idx_categories_name" to table: "categories"
CREATE UNIQUE INDEX "idx_categories_name" ON "categories" ("name");
-- Create index "idx_categories_slug" to table: "categories"
CREATE UNIQUE INDEX "idx_categories_slug" ON "categories" ("slug");
-- Create "products" table
CREATE TABLE "products" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "category_id" uuid NULL,
  "name" character varying(150) NOT NULL,
  "slug" character varying(180) NOT NULL,
  "description" text NULL,
  "status" character varying(20) NOT NULL DEFAULT 'draft',
  "created_by_id" uuid NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_products_category" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_products_created_by" FOREIGN KEY ("created_by_id") REFERENCES "users" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "chk_products_status" CHECK ((status)::text = ANY ((ARRAY['draft'::character varying, 'active'::character varying, 'archived'::character varying])::text[]))
);
-- Create index "idx_products_category_id" to table: "products"
CREATE INDEX "idx_products_category_id" ON "products" ("category_id");
-- Create index "idx_products_created_by" to table: "products"
CREATE INDEX "idx_products_created_by" ON "products" ("created_by_id");
-- Create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "products" ("deleted_at");
-- Create index "idx_products_name" to table: "products"
CREATE INDEX "idx_products_name" ON "products" ("name");
-- Create index "idx_products_slug" to table: "products"
CREATE UNIQUE INDEX "idx_products_slug" ON "products" ("slug");
-- Create index "idx_products_status" to table: "products"
CREATE INDEX "idx_products_status" ON "products" ("status");
-- Create "colors" table
CREATE TABLE "colors" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "name" character varying(50) NOT NULL,
  "hex_code" character varying(7) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_colors_name" to table: "colors"
CREATE UNIQUE INDEX "idx_colors_name" ON "colors" ("name");
-- Create "sizes" table
CREATE TABLE "sizes" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "name" character varying(20) NOT NULL,
  "sort_order" integer NOT NULL DEFAULT 0,
  PRIMARY KEY ("id")
);
-- Create index "idx_sizes_name" to table: "sizes"
CREATE UNIQUE INDEX "idx_sizes_name" ON "sizes" ("name");
-- Create "product_variants" table
CREATE TABLE "product_variants" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "product_id" uuid NOT NULL,
  "size_id" uuid NOT NULL,
  "color_id" uuid NOT NULL,
  "sku" character varying(64) NOT NULL,
  "price" bigint NOT NULL DEFAULT 0,
  "stock" integer NOT NULL DEFAULT 0,
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_variants_color" FOREIGN KEY ("color_id") REFERENCES "colors" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "fk_product_variants_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_product_variants_size" FOREIGN KEY ("size_id") REFERENCES "sizes" ("id") ON UPDATE CASCADE ON DELETE RESTRICT,
  CONSTRAINT "chk_product_variants_price" CHECK (price >= 0),
  CONSTRAINT "chk_product_variants_stock" CHECK (stock >= 0)
);
-- Create index "idx_product_variants_color_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_color_id" ON "product_variants" ("color_id");
-- Create index "idx_product_variants_combo" to table: "product_variants"
CREATE UNIQUE INDEX "idx_product_variants_combo" ON "product_variants" ("product_id", "size_id", "color_id");
-- Create index "idx_product_variants_product_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_product_id" ON "product_variants" ("product_id");
-- Create index "idx_product_variants_size_id" to table: "product_variants"
CREATE INDEX "idx_product_variants_size_id" ON "product_variants" ("size_id");
-- Create index "idx_product_variants_sku" to table: "product_variants"
CREATE UNIQUE INDEX "idx_product_variants_sku" ON "product_variants" ("sku");
-- Create "order_items" table
CREATE TABLE "order_items" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "order_id" uuid NOT NULL,
  "product_id" uuid NULL,
  "variant_id" uuid NULL,
  "product_name" character varying(150) NOT NULL,
  "variant_label" character varying(100) NOT NULL,
  "sku" character varying(64) NOT NULL,
  "price" bigint NOT NULL DEFAULT 0,
  "quantity" integer NOT NULL DEFAULT 1,
  "subtotal" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_order_items_order" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_order_items_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "fk_order_items_variant" FOREIGN KEY ("variant_id") REFERENCES "product_variants" ("id") ON UPDATE CASCADE ON DELETE SET NULL,
  CONSTRAINT "chk_order_items_price" CHECK (price >= 0),
  CONSTRAINT "chk_order_items_quantity" CHECK (quantity > 0),
  CONSTRAINT "chk_order_items_subtotal" CHECK (subtotal >= 0)
);
-- Create index "idx_order_items_order_id" to table: "order_items"
CREATE INDEX "idx_order_items_order_id" ON "order_items" ("order_id");
-- Create index "idx_order_items_product_id" to table: "order_items"
CREATE INDEX "idx_order_items_product_id" ON "order_items" ("product_id");
-- Create index "idx_order_items_variant_id" to table: "order_items"
CREATE INDEX "idx_order_items_variant_id" ON "order_items" ("variant_id");
-- Create "product_images" table
CREATE TABLE "product_images" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "product_id" uuid NOT NULL,
  "object_key" character varying(255) NOT NULL,
  "alt_text" character varying(150) NULL,
  "sort_order" integer NOT NULL DEFAULT 0,
  "is_primary" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_product_images_product" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_product_images_product_id" to table: "product_images"
CREATE INDEX "idx_product_images_product_id" ON "product_images" ("product_id");
-- Create "permissions" table
CREATE TABLE "permissions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "key" character varying(100) NOT NULL,
  "page" character varying(50) NOT NULL,
  "action" character varying(20) NOT NULL,
  "label" character varying(100) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id")
);
-- Create index "idx_permissions_key" to table: "permissions"
CREATE UNIQUE INDEX "idx_permissions_key" ON "permissions" ("key");
-- Create index "idx_permissions_page_action" to table: "permissions"
CREATE UNIQUE INDEX "idx_permissions_page_action" ON "permissions" ("page", "action");
-- Create "role_permissions" table
CREATE TABLE "role_permissions" (
  "role_id" uuid NOT NULL,
  "permission_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("role_id", "permission_id"),
  CONSTRAINT "fk_role_permissions_permission" FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_role_permissions_role" FOREIGN KEY ("role_id") REFERENCES "roles" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_role_permissions_permission_id" to table: "role_permissions"
CREATE INDEX "idx_role_permissions_permission_id" ON "role_permissions" ("permission_id");
