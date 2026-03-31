CREATE TABLE vendor_vouchers (
    id UUID PRIMARY KEY NOT NULL,
    vendor_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(120) NOT NULL,
    description TEXT,
    discount_amount NUMERIC(18,2) NOT NULL,
    quota_total INT NOT NULL,
    quota_used INT NOT NULL DEFAULT 0,
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_vendor_vouchers_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id) ON DELETE CASCADE,
    CONSTRAINT ck_vendor_vouchers_discount_positive CHECK (discount_amount > 0),
    CONSTRAINT ck_vendor_vouchers_quota_valid CHECK (quota_total > 0 AND quota_used >= 0 AND quota_used <= quota_total),
    CONSTRAINT ck_vendor_vouchers_period_valid CHECK (ends_at >= starts_at)
);

CREATE UNIQUE INDEX uq_vendor_vouchers_vendor_code
    ON vendor_vouchers (vendor_id, lower(code));

CREATE INDEX idx_vendor_vouchers_vendor_period_active
    ON vendor_vouchers (vendor_id, is_active, starts_at, ends_at);

CREATE TABLE vendor_voucher_products (
    voucher_id UUID NOT NULL,
    product_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT pk_vendor_voucher_products PRIMARY KEY (voucher_id, product_id),
    CONSTRAINT fk_vendor_voucher_products_voucher FOREIGN KEY (voucher_id) REFERENCES vendor_vouchers (id) ON DELETE CASCADE,
    CONSTRAINT fk_vendor_voucher_products_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
);

CREATE INDEX idx_vendor_voucher_products_product
    ON vendor_voucher_products (product_id);
