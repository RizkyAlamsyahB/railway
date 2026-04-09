ALTER TABLE users
    DROP COLUMN IF EXISTS deletion_requested_at,
    DROP COLUMN IF EXISTS deletion_reason,
    DROP CONSTRAINT ck_users_status,
    ADD CONSTRAINT ck_users_status CHECK (status IN ('pending', 'active', 'blocked'));
