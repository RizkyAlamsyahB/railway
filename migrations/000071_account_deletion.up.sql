ALTER TABLE users
    DROP CONSTRAINT ck_users_status,
    ADD CONSTRAINT ck_users_status CHECK (status IN ('pending', 'active', 'blocked', 'deactivated')),
    ADD COLUMN deletion_reason VARCHAR(50),
    ADD COLUMN deletion_requested_at TIMESTAMPTZ;
