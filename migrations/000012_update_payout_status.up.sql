UPDATE payout_batches
SET status = CASE status
    WHEN 'draft' THEN 'schedule'
    WHEN 'approved' THEN 'schedule'
    WHEN 'ready' THEN 'schedule'
    WHEN 'paid' THEN 'completed'
    WHEN 'complete' THEN 'completed'
    WHEN 'on hold' THEN 'on_hold'
    ELSE status
END;

ALTER TABLE payout_batches
    DROP CONSTRAINT IF EXISTS ck_payout_batches_status;

ALTER TABLE payout_batches
    ADD CONSTRAINT ck_payout_batches_status CHECK (
        status IN ('schedule', 'completed', 'failed', 'on_hold')
    );
