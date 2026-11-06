-- Modify "expense_items" table
ALTER TABLE "expense_items" DROP CONSTRAINT "fk_expense_items_expense", ADD CONSTRAINT "fk_expenses_items" FOREIGN KEY ("expense_id") REFERENCES "expenses" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "orders" table
ALTER TABLE "orders" ADD COLUMN "stock_deducted" boolean NOT NULL DEFAULT false;
-- Modify "order_items" table
ALTER TABLE "order_items" DROP CONSTRAINT "fk_order_items_order", ADD CONSTRAINT "fk_orders_items" FOREIGN KEY ("order_id") REFERENCES "orders" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "product_images" table
ALTER TABLE "product_images" DROP CONSTRAINT "fk_product_images_product", ADD CONSTRAINT "fk_products_images" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
-- Modify "product_variants" table
ALTER TABLE "product_variants" DROP CONSTRAINT "fk_product_variants_product", ADD CONSTRAINT "fk_products_variants" FOREIGN KEY ("product_id") REFERENCES "products" ("id") ON UPDATE CASCADE ON DELETE CASCADE;
