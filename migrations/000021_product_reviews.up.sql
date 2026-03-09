CREATE TABLE product_reviews (
    id UUID PRIMARY KEY NOT NULL,
    product_id UUID NOT NULL,
    order_id UUID NOT NULL,
    order_item_id UUID NOT NULL,
    user_id UUID NOT NULL,
    rating SMALLINT NOT NULL,
    review_text TEXT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'published',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_product_reviews_rating_range CHECK (rating BETWEEN 1 AND 5),
    CONSTRAINT ck_product_reviews_status CHECK (status IN ('published', 'hidden')),
    CONSTRAINT uq_product_reviews_user_order_item UNIQUE (user_id, order_item_id),
    CONSTRAINT fk_product_reviews_product FOREIGN KEY (product_id) REFERENCES products (id),
    CONSTRAINT fk_product_reviews_order FOREIGN KEY (order_id) REFERENCES orders (id),
    CONSTRAINT fk_product_reviews_order_item FOREIGN KEY (order_item_id) REFERENCES order_items (id),
    CONSTRAINT fk_product_reviews_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE product_review_images (
    id UUID PRIMARY KEY NOT NULL,
    review_id UUID NOT NULL,
    object_key TEXT NOT NULL,
    mime_type VARCHAR(100),
    file_size_bytes INT,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_product_review_images_sort_non_negative CHECK (sort_order >= 0),
    CONSTRAINT fk_product_review_images_review FOREIGN KEY (review_id) REFERENCES product_reviews (id) ON DELETE CASCADE
);

CREATE TABLE product_review_stats (
    product_id UUID PRIMARY KEY NOT NULL,
    total_reviews BIGINT NOT NULL DEFAULT 0,
    total_stars BIGINT NOT NULL DEFAULT 0,
    average_rating NUMERIC(4,2) NOT NULL DEFAULT 0,
    star_0_count BIGINT NOT NULL DEFAULT 0,
    star_1_count BIGINT NOT NULL DEFAULT 0,
    star_2_count BIGINT NOT NULL DEFAULT 0,
    star_3_count BIGINT NOT NULL DEFAULT 0,
    star_4_count BIGINT NOT NULL DEFAULT 0,
    star_5_count BIGINT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_product_review_stats_non_negative CHECK (
        total_reviews >= 0 AND total_stars >= 0 AND
        star_0_count >= 0 AND star_1_count >= 0 AND star_2_count >= 0 AND
        star_3_count >= 0 AND star_4_count >= 0 AND star_5_count >= 0
    ),
    CONSTRAINT fk_product_review_stats_product FOREIGN KEY (product_id) REFERENCES products (id)
);

CREATE INDEX idx_product_reviews_product_created_at ON product_reviews (product_id, created_at DESC);
CREATE INDEX idx_product_reviews_product_rating ON product_reviews (product_id, rating);
CREATE INDEX idx_product_reviews_user_created_at ON product_reviews (user_id, created_at DESC);
CREATE INDEX idx_product_review_images_review_sort ON product_review_images (review_id, sort_order);
