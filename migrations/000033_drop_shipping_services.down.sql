-- Recreate shipping_services table
CREATE TABLE shipping_services (
    id UUID PRIMARY KEY NOT NULL,
    code VARCHAR(40) UNIQUE NOT NULL,
    name VARCHAR(80) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Recreate product_shipping_services join table
CREATE TABLE product_shipping_services (
    product_id UUID NOT NULL,
    shipping_service_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT pk_product_shipping_services PRIMARY KEY (product_id, shipping_service_id),
    CONSTRAINT fk_product_shipping_services_product FOREIGN KEY (product_id) REFERENCES products (id),
    CONSTRAINT fk_product_shipping_services_shipping_service FOREIGN KEY (shipping_service_id) REFERENCES shipping_services (id)
);

-- Recreate indexes
CREATE INDEX idx_product_shipping_services_product_id ON product_shipping_services (product_id);
CREATE INDEX idx_product_shipping_services_shipping_service_id ON product_shipping_services (shipping_service_id);

-- Recreate trigger functions and triggers
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
    SELECT p.status INTO product_status FROM products p WHERE p.id = target_product_id;
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
