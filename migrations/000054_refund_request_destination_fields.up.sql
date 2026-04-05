ALTER TABLE refunds
    ADD COLUMN IF NOT EXISTS reason_detail TEXT,
    ADD COLUMN IF NOT EXISTS destination_channel_code VARCHAR(40),
    ADD COLUMN IF NOT EXISTS destination_bank_name VARCHAR(80),
    ADD COLUMN IF NOT EXISTS destination_account_number TEXT,
    ADD COLUMN IF NOT EXISTS destination_account_holder_name VARCHAR(120),
    ADD COLUMN IF NOT EXISTS destination_account_last4 VARCHAR(8);

CREATE INDEX IF NOT EXISTS idx_refunds_destination_channel_code
    ON refunds (destination_channel_code);
