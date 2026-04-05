-- Widen refunds.status from varchar(16) to varchar(30) to accommodate 'awaiting_destination' (20 chars)
ALTER TABLE refunds ALTER COLUMN status TYPE VARCHAR(30);
