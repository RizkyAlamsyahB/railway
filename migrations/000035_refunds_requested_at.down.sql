DROP INDEX IF EXISTS idx_refunds_requested_at;

ALTER TABLE refunds
    DROP COLUMN IF EXISTS requested_at;
