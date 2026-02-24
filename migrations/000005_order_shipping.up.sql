CREATE TABLE orders (
    id UUID PRIMARY KEY NOT NULL,
    order_no VARCHAR(30) UNIQUE NOT NULL,
    user_id UUID NOT NULL,
    vendor_id UUID NOT NULL,
    shipping_address_snapshot JSONB NOT NULL,
    order_status VARCHAR(24) NOT NULL DEFAULT 'pending_payment',
    payment_status VARCHAR(16) NOT NULL DEFAULT 'unpaid',
    subtotal NUMERIC(18,2) NOT NULL,
    shipping_fee NUMERIC(18,2) NOT NULL DEFAULT 0,
    platform_fee NUMERIC(18,2) NOT NULL DEFAULT 0,
    grand_total NUMERIC(18,2) NOT NULL,
    placed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_orders_order_status CHECK (order_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'completed', 'canceled', 'refunded')),
    CONSTRAINT ck_orders_payment_status CHECK (payment_status IN ('unpaid', 'paid_partial', 'paid', 'refunded')),
    CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT fk_orders_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id)
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID NOT NULL,
    product_variant_id UUID NOT NULL,
    product_name_snapshot VARCHAR(180) NOT NULL,
    sku_snapshot VARCHAR(80) NOT NULL,
    qty INT NOT NULL,
    unit_price NUMERIC(18,2) NOT NULL,
    line_total NUMERIC(18,2) NOT NULL,
    CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES orders (id),
    CONSTRAINT fk_order_items_product_variant FOREIGN KEY (product_variant_id) REFERENCES product_variants (id)
);

CREATE TABLE order_status_history (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID NOT NULL,
    old_status VARCHAR(24),
    new_status VARCHAR(24) NOT NULL,
    changed_by UUID,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    notes TEXT,
    CONSTRAINT ck_order_status_history_old_status CHECK (
        old_status IS NULL OR old_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'completed', 'canceled', 'refunded')
    ),
    CONSTRAINT ck_order_status_history_new_status CHECK (
        new_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'completed', 'canceled', 'refunded')
    ),
    CONSTRAINT fk_order_status_history_order FOREIGN KEY (order_id) REFERENCES orders (id),
    CONSTRAINT fk_order_status_history_changed_by FOREIGN KEY (changed_by) REFERENCES users (id)
);

CREATE TABLE shipments (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID UNIQUE NOT NULL,
    courier_code VARCHAR(30),
    service_type VARCHAR(30),
    tracking_no VARCHAR(80),
    shipment_status VARCHAR(24) NOT NULL DEFAULT 'waiting_pickup',
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    CONSTRAINT ck_shipments_status CHECK (shipment_status IN ('waiting_pickup', 'shipped', 'delivered', 'failed', 'returned')),
    CONSTRAINT fk_shipments_order FOREIGN KEY (order_id) REFERENCES orders (id)
);
