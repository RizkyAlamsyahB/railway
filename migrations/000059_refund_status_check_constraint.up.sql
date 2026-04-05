-- Add 'awaiting_destination' and 'processing' to refunds.status check constraint
ALTER TABLE refunds DROP CONSTRAINT ck_refunds_status;
ALTER TABLE refunds ADD CONSTRAINT ck_refunds_status
  CHECK (status IN ('requested', 'approved', 'rejected', 'processed', 'awaiting_destination', 'processing'));
