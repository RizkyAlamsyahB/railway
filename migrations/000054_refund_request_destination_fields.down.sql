DROP INDEX IF EXISTS idx_refunds_destination_channel_code;

ALTER TABLE refunds
    DROP COLUMN IF EXISTS destination_account_last4,
    DROP COLUMN IF EXISTS destination_account_holder_name,
    DROP COLUMN IF EXISTS destination_account_number,
    DROP COLUMN IF EXISTS destination_bank_name,
    DROP COLUMN IF EXISTS destination_channel_code,
    DROP COLUMN IF EXISTS reason_detail;
