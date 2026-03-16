ALTER TABLE refunds
    ADD COLUMN IF NOT EXISTS requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

CREATE INDEX IF NOT EXISTS idx_refunds_requested_at ON refunds (requested_at DESC);
