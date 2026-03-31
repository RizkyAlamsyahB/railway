-- Add 'processing' and 'received' order status to the CHECK constraints.
ALTER TABLE orders DROP CONSTRAINT IF EXISTS ck_orders_order_status;
ALTER TABLE orders ADD CONSTRAINT ck_orders_order_status
    CHECK (order_status IN ('pending_payment', 'paid', 'processing', 'packed', 'shipped', 'received', 'completed', 'canceled', 'refunded'));

ALTER TABLE order_status_history DROP CONSTRAINT IF EXISTS ck_order_status_history_old_status;
ALTER TABLE order_status_history ADD CONSTRAINT ck_order_status_history_old_status
    CHECK (old_status IS NULL OR old_status IN ('pending_payment', 'paid', 'processing', 'packed', 'shipped', 'received', 'completed', 'canceled', 'refunded'));

ALTER TABLE order_status_history DROP CONSTRAINT IF EXISTS ck_order_status_history_new_status;
ALTER TABLE order_status_history ADD CONSTRAINT ck_order_status_history_new_status
    CHECK (new_status IN ('pending_payment', 'paid', 'processing', 'packed', 'shipped', 'received', 'completed', 'canceled', 'refunded'));
