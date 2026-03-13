-- Revert order status constraints to exclude received.
ALTER TABLE order_status_history
    DROP CONSTRAINT ck_order_status_history_old_status;

ALTER TABLE order_status_history
    DROP CONSTRAINT ck_order_status_history_new_status;

ALTER TABLE order_status_history
    ADD CONSTRAINT ck_order_status_history_old_status CHECK (
        old_status IS NULL OR old_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'completed', 'canceled', 'refunded')
    );

ALTER TABLE order_status_history
    ADD CONSTRAINT ck_order_status_history_new_status CHECK (
        new_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'completed', 'canceled', 'refunded')
    );

ALTER TABLE orders
    DROP CONSTRAINT ck_orders_order_status;

ALTER TABLE orders
    ADD CONSTRAINT ck_orders_order_status CHECK (
        order_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'completed', 'canceled', 'refunded')
    );

-- Remove escrow balance column.
ALTER TABLE vendor_balances
    DROP CONSTRAINT ck_vendor_balances_escrow;

ALTER TABLE vendor_balances
    DROP COLUMN escrow_balance;
