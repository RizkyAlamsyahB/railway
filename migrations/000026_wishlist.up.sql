CREATE TABLE wishlist_items (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL,
    product_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_wishlist_items_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_wishlist_items_product FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE UNIQUE INDEX uq_wishlist_items_user_product
    ON wishlist_items (user_id, product_id);

CREATE INDEX idx_wishlist_items_user_created_at
    ON wishlist_items (user_id, created_at DESC, id DESC);

CREATE INDEX idx_wishlist_items_product
    ON wishlist_items (product_id);
