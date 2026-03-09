-- Drop indexes from hardening migration (000009)
DROP INDEX IF EXISTS idx_product_shipping_services_shipping_service_id;
DROP INDEX IF EXISTS idx_product_shipping_services_product_id;

-- Drop triggers that enforce shipping service constraints
DROP TRIGGER IF EXISTS trg_products_published_requires_shipping ON products;
DROP FUNCTION IF EXISTS ensure_published_product_has_shipping_services();

DROP TRIGGER IF EXISTS trg_product_shipping_services_published_requires_shipping ON product_shipping_services;
DROP FUNCTION IF EXISTS ensure_published_product_keeps_shipping_services();

-- Drop tables (order matters due to FK)
DROP TABLE IF EXISTS product_shipping_services;
DROP TABLE IF EXISTS shipping_services;
