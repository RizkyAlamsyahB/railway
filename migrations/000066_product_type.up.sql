CREATE TYPE product_type AS ENUM ('single', 'variant', 'package');

ALTER TABLE products
    ADD COLUMN product_type product_type NOT NULL DEFAULT 'single';
