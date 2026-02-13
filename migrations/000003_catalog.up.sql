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

CREATE TABLE shipping_services (
    id UUID PRIMARY KEY NOT NULL,
    code VARCHAR(40) UNIQUE NOT NULL,
    name VARCHAR(80) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_shipping_services (
    product_id UUID NOT NULL,
    shipping_service_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT pk_product_shipping_services PRIMARY KEY (product_id, shipping_service_id),
    CONSTRAINT fk_product_shipping_services_product FOREIGN KEY (product_id) REFERENCES products (id),
    CONSTRAINT fk_product_shipping_services_shipping_service FOREIGN KEY (shipping_service_id) REFERENCES shipping_services (id)
);

CREATE TABLE product_images (
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    image_url TEXT NOT NULL,
    mime_type VARCHAR(100),
    file_size_bytes INT,
    is_primary BOOLEAN NOT NULL DEFAULT false,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_product_images_sort_order_non_negative CHECK (sort_order >= 0),
    CONSTRAINT fk_product_images_product FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE TABLE product_variants (
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    sku VARCHAR(80) NOT NULL,
    variant_name VARCHAR(120) NOT NULL,
    price NUMERIC(18,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    stock_on_hand INT NOT NULL,
    weight_gram INT,
    is_default BOOLEAN NOT NULL DEFAULT false,
    is_active BOOLEAN NOT NULL DEFAULT true,
    CONSTRAINT uq_product_variants_product_sku UNIQUE (product_id, sku),
    CONSTRAINT ck_product_variants_price_positive CHECK (price > 0),
    CONSTRAINT ck_product_variants_stock_non_negative CHECK (stock_on_hand >= 0),
    CONSTRAINT ck_product_variants_weight_non_negative CHECK (weight_gram IS NULL OR weight_gram >= 0),
    CONSTRAINT fk_product_variants_product FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE UNIQUE INDEX uq_product_variants_default_per_product
    ON product_variants (product_id)
    WHERE is_default;

CREATE UNIQUE INDEX uq_product_images_primary_per_product
    ON product_images (product_id)
    WHERE is_primary;

CREATE OR REPLACE FUNCTION ensure_published_product_has_shipping_services()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.status = 'published' THEN
        IF NOT EXISTS (
            SELECT 1
            FROM product_shipping_services pss
            WHERE pss.product_id = NEW.id
        ) THEN
            RAISE EXCEPTION 'published product must have at least one shipping service';
        END IF;
    END IF;

    RETURN NEW;
END;
$$;

CREATE CONSTRAINT TRIGGER trg_products_published_requires_shipping
AFTER INSERT OR UPDATE OF status ON products
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION ensure_published_product_has_shipping_services();

CREATE OR REPLACE FUNCTION ensure_published_product_keeps_shipping_services()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    target_product_id UUID;
    product_status VARCHAR(16);
BEGIN
    target_product_id := COALESCE(OLD.product_id, NEW.product_id);

    SELECT p.status
    INTO product_status
    FROM products p
    WHERE p.id = target_product_id;

    IF product_status = 'published' THEN
        IF NOT EXISTS (
            SELECT 1
            FROM product_shipping_services pss
            WHERE pss.product_id = target_product_id
        ) THEN
            RAISE EXCEPTION 'published product must have at least one shipping service';
        END IF;
    END IF;

    RETURN COALESCE(NEW, OLD);
END;
$$;

CREATE CONSTRAINT TRIGGER trg_product_shipping_services_published_requires_shipping
AFTER DELETE OR UPDATE OF product_id ON product_shipping_services
DEFERRABLE INITIALLY DEFERRED
FOR EACH ROW
EXECUTE FUNCTION ensure_published_product_keeps_shipping_services();

CREATE OR REPLACE FUNCTION enforce_product_image_limit()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    current_count INT;
BEGIN
    SELECT COUNT(*)
    INTO current_count
    FROM product_images pi
    WHERE pi.product_id = NEW.product_id
      AND (TG_OP <> 'UPDATE' OR pi.id <> OLD.id);

    IF current_count >= 10 THEN
        RAISE EXCEPTION 'maximum 10 images per product is allowed';
    END IF;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_product_images_max_10
BEFORE INSERT OR UPDATE OF product_id ON product_images
FOR EACH ROW
EXECUTE FUNCTION enforce_product_image_limit();
