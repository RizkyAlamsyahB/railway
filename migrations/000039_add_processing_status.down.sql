-- Revert: remove 'processing' from order_status CHECK constraint.
ALTER TABLE orders DROP CONSTRAINT IF EXISTS ck_orders_order_status;
ALTER TABLE orders ADD CONSTRAINT ck_orders_order_status
    CHECK (order_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'received', 'completed', 'canceled', 'refunded'));
