ALTER TABLE carts
    DROP CONSTRAINT IF EXISTS carts_user_id_key;

CREATE UNIQUE INDEX uq_carts_user_active
    ON carts (user_id)
    WHERE status = 'active';

CREATE INDEX idx_carts_user_status_updated_at
    ON carts (user_id, status, updated_at DESC);
