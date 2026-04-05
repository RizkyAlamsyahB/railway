DROP INDEX IF EXISTS uq_refund_evidences_single_video;
DROP INDEX IF EXISTS idx_refund_evidences_media_type;
DROP INDEX IF EXISTS idx_refund_evidences_refund_sort;

DROP TABLE IF EXISTS refund_evidences;

ALTER TABLE refunds
    DROP CONSTRAINT IF EXISTS fk_refunds_return_reason;

DROP INDEX IF EXISTS idx_refunds_return_reason_id;

ALTER TABLE refunds
    DROP COLUMN IF EXISTS return_reason_id;
