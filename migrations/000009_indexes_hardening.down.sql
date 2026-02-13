ALTER TABLE refunds
    DROP CONSTRAINT IF EXISTS ck_refunds_amount_positive;

ALTER TABLE order_items
    DROP CONSTRAINT IF EXISTS ck_order_items_qty_positive;

ALTER TABLE cart_items
    DROP CONSTRAINT IF EXISTS ck_cart_items_qty_positive;

DROP INDEX IF EXISTS idx_ledger_lines_account;
DROP INDEX IF EXISTS idx_ledger_journals_source;
DROP INDEX IF EXISTS uq_payout_items_batch_order;
DROP INDEX IF EXISTS idx_payout_batches_vendor_status;
DROP INDEX IF EXISTS idx_payment_invoices_order_status;
DROP INDEX IF EXISTS idx_orders_status_payment_status;
DROP INDEX IF EXISTS idx_orders_vendor_created_at;
DROP INDEX IF EXISTS idx_orders_user_created_at;
DROP INDEX IF EXISTS uq_cart_items_cart_variant;
DROP INDEX IF EXISTS idx_product_shipping_services_shipping_service_id;
DROP INDEX IF EXISTS idx_product_shipping_services_product_id;
DROP INDEX IF EXISTS idx_product_images_product_sort;
DROP INDEX IF EXISTS idx_product_variants_product_active;
DROP INDEX IF EXISTS idx_products_vendor_status;
