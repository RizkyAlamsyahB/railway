CREATE TABLE product_promotions (
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    promo_name VARCHAR(120) NOT NULL,
    promo_price NUMERIC(18,2) NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_product_promotions_price_positive CHECK (promo_price > 0),
    CONSTRAINT ck_product_promotions_period_valid CHECK (ends_at >= starts_at),
    CONSTRAINT fk_product_promotions_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE INDEX idx_product_promotions_product_active_period
    ON product_promotions (product_id, is_active, starts_at, ends_at);
