-- Remove Xendit payout tracking columns from payout_batches
ALTER TABLE payout_batches
    DROP CONSTRAINT IF EXISTS ck_payout_batches_xendit_status;

ALTER TABLE payout_batches
    DROP COLUMN IF EXISTS xendit_payout_id,
    DROP COLUMN IF EXISTS channel_code,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS xendit_status,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;
