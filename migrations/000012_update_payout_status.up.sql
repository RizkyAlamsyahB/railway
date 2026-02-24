UPDATE payout_batches
SET status = CASE status
    WHEN 'draft' THEN 'ready'
    WHEN 'approved' THEN 'schedule'
    WHEN 'paid' THEN 'complete'
    ELSE status
END;

ALTER TABLE payout_batches
    DROP CONSTRAINT IF EXISTS ck_payout_batches_status;

ALTER TABLE payout_batches
    ADD CONSTRAINT ck_payout_batches_status CHECK (
        status IN ('ready', 'schedule', 'complete', 'failed', 'on hold')
    );
