CREATE TABLE carts (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID UNIQUE NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_carts_status CHECK (status IN ('active', 'converted', 'abandoned')),
    CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE cart_items (
    id UUID PRIMARY KEY NOT NULL,
    cart_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    qty INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES carts (id),
    CONSTRAINT fk_cart_items_product_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants (id)
);
