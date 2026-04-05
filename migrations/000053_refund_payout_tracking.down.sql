DROP INDEX IF EXISTS idx_refunds_payout_status;
DROP INDEX IF EXISTS uq_refunds_payout_reference_id;

ALTER TABLE refunds
    DROP CONSTRAINT IF EXISTS fk_refunds_payout_requested_by;

ALTER TABLE refunds
    DROP COLUMN IF EXISTS payout_completed_at,
    DROP COLUMN IF EXISTS payout_requested_at,
    DROP COLUMN IF EXISTS payout_requested_by,
    DROP COLUMN IF EXISTS payout_failed_reason,
    DROP COLUMN IF EXISTS payout_status,
    DROP COLUMN IF EXISTS payout_id,
    DROP COLUMN IF EXISTS payout_destination_account_last4,
    DROP COLUMN IF EXISTS payout_destination_account_holder_name,
    DROP COLUMN IF EXISTS payout_destination_bank_name,
    DROP COLUMN IF EXISTS payout_channel_code,
    DROP COLUMN IF EXISTS payout_reference_id,
    DROP COLUMN IF EXISTS payout_strategy;
