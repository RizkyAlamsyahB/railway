CREATE TABLE categories (
    id UUID PRIMARY KEY NOT NULL,
    parent_id UUID,
    name VARCHAR(80) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES categories (id)
);

CREATE TABLE products (
    id UUID PRIMARY KEY NOT NULL,
    vendor_id UUID NOT NULL,
    category_id UUID NOT NULL,
    name VARCHAR(180) NOT NULL,
    slug VARCHAR(220) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    halal_ai_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    halal_ai_notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_products_status CHECK (status IN ('draft', 'published', 'blocked', 'archived')),
    CONSTRAINT ck_products_halal_ai_status CHECK (halal_ai_status IN ('pending', 'passed', 'failed')),
    CONSTRAINT fk_products_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id),
    CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories (id)
);

CREATE TABLE product_images (
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    image_url TEXT NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    CONSTRAINT fk_product_images_product FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    sku VARCHAR(80) UNIQUE NOT NULL,
    variant_name VARCHAR(120) NOT NULL,
    price NUMERIC(18,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    stock_on_hand INT NOT NULL,
    weight_gram INT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    CONSTRAINT fk_product_variants_product FOREIGN KEY (product_id) REFERENCES products (id)
);
