-- Add 'processing' order status to the CHECK constraint on orders table.
ALTER TABLE orders DROP CONSTRAINT IF EXISTS ck_orders_order_status;
ALTER TABLE orders ADD CONSTRAINT ck_orders_order_status
    CHECK (order_status IN ('pending_payment', 'paid', 'processing', 'packed', 'shipped', 'completed', 'canceled', 'refunded'));
