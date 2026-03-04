UPDATE payout_batches
SET status = CASE status
    WHEN 'schedule' THEN 'approved'
    WHEN 'completed' THEN 'paid'
    WHEN 'complete' THEN 'paid'
    WHEN 'ready' THEN 'draft'
    WHEN 'on_hold' THEN 'draft'
    WHEN 'on hold' THEN 'draft'
    ELSE status
END;

ALTER TABLE payout_batches
    DROP CONSTRAINT IF EXISTS ck_payout_batches_status;

ALTER TABLE payout_batches
    ADD CONSTRAINT ck_payout_batches_status CHECK (
        status IN ('draft', 'approved', 'paid', 'failed')
    );
