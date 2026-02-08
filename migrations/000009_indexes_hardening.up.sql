CREATE INDEX idx_products_vendor_status ON products (vendor_id, status);
CREATE INDEX idx_product_variants_product_active ON product_variants (product_id, is_active);
CREATE UNIQUE INDEX uq_cart_items_cart_variant ON cart_items (cart_id, product_variant_id);
CREATE INDEX idx_orders_user_created_at ON orders (user_id, created_at);
CREATE INDEX idx_orders_vendor_created_at ON orders (vendor_id, created_at);
CREATE INDEX idx_orders_status_payment_status ON orders (order_status, payment_status);
CREATE INDEX idx_payment_invoices_order_status ON payment_invoices (order_id, status);
CREATE INDEX idx_payout_batches_vendor_status ON payout_batches (vendor_id, status);
CREATE UNIQUE INDEX uq_payout_items_batch_order ON payout_items (payout_batch_id, order_id);
CREATE INDEX idx_ledger_journals_source ON ledger_journals (source_type, source_id);
CREATE INDEX idx_ledger_lines_account ON ledger_lines (account_id);

ALTER TABLE cart_items
    ADD CONSTRAINT ck_cart_items_qty_positive CHECK (qty > 0);

ALTER TABLE order_items
    ADD CONSTRAINT ck_order_items_qty_positive CHECK (qty > 0);

ALTER TABLE refunds
    ADD CONSTRAINT ck_refunds_amount_positive CHECK (amount > 0);
