-- Add 'processing' to order_status_history CHECK constraints (missed in 000037).
ALTER TABLE order_status_history DROP CONSTRAINT IF EXISTS ck_order_status_history_old_status;
ALTER TABLE order_status_history ADD CONSTRAINT ck_order_status_history_old_status
    CHECK (old_status IS NULL OR old_status IN ('pending_payment', 'paid', 'processing', 'packed', 'shipped', 'completed', 'canceled', 'refunded'));

ALTER TABLE order_status_history DROP CONSTRAINT IF EXISTS ck_order_status_history_new_status;
ALTER TABLE order_status_history ADD CONSTRAINT ck_order_status_history_new_status
    CHECK (new_status IN ('pending_payment', 'paid', 'processing', 'packed', 'shipped', 'completed', 'canceled', 'refunded'));
