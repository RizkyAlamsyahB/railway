ALTER TABLE refunds DROP CONSTRAINT ck_refunds_status;
ALTER TABLE refunds ADD CONSTRAINT ck_refunds_status
  CHECK (status IN ('requested', 'approved', 'rejected', 'processed'));
