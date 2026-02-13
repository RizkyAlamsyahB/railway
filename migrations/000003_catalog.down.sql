DROP TRIGGER IF EXISTS trg_product_images_max_10 ON product_images;
DROP FUNCTION IF EXISTS enforce_product_image_limit();

DROP TRIGGER IF EXISTS trg_products_published_requires_shipping ON products;
DROP FUNCTION IF EXISTS ensure_published_product_has_shipping_services();
DROP TRIGGER IF EXISTS trg_product_shipping_services_published_requires_shipping ON product_shipping_services;
DROP FUNCTION IF EXISTS ensure_published_product_keeps_shipping_services();

DROP INDEX IF EXISTS uq_product_images_primary_per_product;
DROP INDEX IF EXISTS uq_product_variants_default_per_product;

DROP TABLE IF EXISTS product_shipping_services;
DROP TABLE IF EXISTS shipping_services;
DROP TABLE IF EXISTS product_images;
DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
