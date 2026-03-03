DROP INDEX IF EXISTS idx_carts_user_status_updated_at;
DROP INDEX IF EXISTS uq_carts_user_active;

ALTER TABLE carts
    ADD CONSTRAINT carts_user_id_key UNIQUE (user_id);
