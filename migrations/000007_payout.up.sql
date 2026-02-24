CREATE TABLE payout_batches (
    id UUID PRIMARY KEY NOT NULL,
    vendor_id UUID NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    total_gross NUMERIC(18,2) NOT NULL,
    total_fee NUMERIC(18,2) NOT NULL,
    total_net NUMERIC(18,2) NOT NULL,
    paid_at TIMESTAMPTZ,
    created_by UUID NOT NULL,
    CONSTRAINT ck_payout_batches_status CHECK (status IN ('draft', 'approved', 'paid', 'failed')),
    CONSTRAINT fk_payout_batches_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id),
    CONSTRAINT fk_payout_batches_created_by FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE payout_items (
    id UUID PRIMARY KEY NOT NULL,
    payout_batch_id UUID NOT NULL,
    order_id UUID NOT NULL,
    gross_amount NUMERIC(18,2) NOT NULL,
    platform_fee_amount NUMERIC(18,2) NOT NULL,
    net_amount NUMERIC(18,2) NOT NULL,
    CONSTRAINT fk_payout_items_batch FOREIGN KEY (payout_batch_id) REFERENCES payout_batches (id),
    CONSTRAINT fk_payout_items_order FOREIGN KEY (order_id) REFERENCES orders (id)
);
