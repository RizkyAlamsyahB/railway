-- Add Xendit payout tracking columns to payout_batches
ALTER TABLE payout_batches
    ADD COLUMN IF NOT EXISTS xendit_payout_id VARCHAR(128),
    ADD COLUMN IF NOT EXISTS channel_code VARCHAR(32),
    ADD COLUMN IF NOT EXISTS description TEXT,
    ADD COLUMN IF NOT EXISTS xendit_status VARCHAR(32),
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE payout_batches
    ADD CONSTRAINT ck_payout_batches_xendit_status CHECK (
        xendit_status IN ('ready', 'schedule', 'complete', 'failed', 'on hold')
    );
