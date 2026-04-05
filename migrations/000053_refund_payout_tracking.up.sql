ALTER TABLE refunds
    ADD COLUMN IF NOT EXISTS payout_strategy VARCHAR(40),
    ADD COLUMN IF NOT EXISTS payout_reference_id VARCHAR(80),
    ADD COLUMN IF NOT EXISTS payout_channel_code VARCHAR(40),
    ADD COLUMN IF NOT EXISTS payout_destination_bank_name VARCHAR(80),
    ADD COLUMN IF NOT EXISTS payout_destination_account_holder_name VARCHAR(120),
    ADD COLUMN IF NOT EXISTS payout_destination_account_last4 VARCHAR(8),
    ADD COLUMN IF NOT EXISTS payout_id VARCHAR(100),
    ADD COLUMN IF NOT EXISTS payout_status VARCHAR(40),
    ADD COLUMN IF NOT EXISTS payout_failed_reason TEXT,
    ADD COLUMN IF NOT EXISTS payout_requested_by UUID,
    ADD COLUMN IF NOT EXISTS payout_requested_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS payout_completed_at TIMESTAMPTZ;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_refunds_payout_requested_by'
    ) THEN
        ALTER TABLE refunds
            ADD CONSTRAINT fk_refunds_payout_requested_by
            FOREIGN KEY (payout_requested_by) REFERENCES users(id);
    END IF;
END $$;

CREATE UNIQUE INDEX IF NOT EXISTS uq_refunds_payout_reference_id
    ON refunds (payout_reference_id)
    WHERE payout_reference_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_refunds_payout_status
    ON refunds (payout_status);
