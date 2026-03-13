CREATE TABLE IF NOT EXISTS otp_codes (
    id UUID PRIMARY KEY NOT NULL,
    user_id UUID,
    email VARCHAR(255) NOT NULL,
    purpose VARCHAR(50) NOT NULL,
    channel VARCHAR(20) NOT NULL,
    code_hash TEXT NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    attempt_count INT NOT NULL DEFAULT 0,
    consumed_at TIMESTAMPTZ,
    invalidated_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_otp_codes_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT ck_otp_codes_channel CHECK (channel IN ('email')),
    CONSTRAINT ck_otp_codes_attempt_count CHECK (attempt_count >= 0)
);

CREATE UNIQUE INDEX uq_otp_codes_active_email_purpose_channel
    ON otp_codes (lower(email), purpose, channel)
    WHERE consumed_at IS NULL AND invalidated_at IS NULL;

CREATE INDEX idx_otp_codes_expires_at ON otp_codes (expires_at);
CREATE INDEX idx_otp_codes_email_purpose_created_at ON otp_codes (lower(email), purpose, created_at DESC);
