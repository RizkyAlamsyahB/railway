-- Add escrow balance to vendor balances.
ALTER TABLE vendor_balances
    ADD COLUMN escrow_balance NUMERIC(18,2) NOT NULL DEFAULT 0;

ALTER TABLE vendor_balances
    ADD CONSTRAINT ck_vendor_balances_escrow CHECK (escrow_balance >= 0);

-- Extend order status constraints to include received.
ALTER TABLE orders
    DROP CONSTRAINT ck_orders_order_status;

ALTER TABLE orders
    ADD CONSTRAINT ck_orders_order_status CHECK (
        order_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'received', 'completed', 'canceled', 'refunded')
    );

ALTER TABLE order_status_history
    DROP CONSTRAINT ck_order_status_history_old_status;

ALTER TABLE order_status_history
    DROP CONSTRAINT ck_order_status_history_new_status;

ALTER TABLE order_status_history
    ADD CONSTRAINT ck_order_status_history_old_status CHECK (
        old_status IS NULL OR old_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'received', 'completed', 'canceled', 'refunded')
    );

ALTER TABLE order_status_history
    ADD CONSTRAINT ck_order_status_history_new_status CHECK (
        new_status IN ('pending_payment', 'paid', 'packed', 'shipped', 'received', 'completed', 'canceled', 'refunded')
    );
