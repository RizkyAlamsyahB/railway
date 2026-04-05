ALTER TABLE refunds
    ADD COLUMN IF NOT EXISTS return_reason_id INTEGER;

CREATE INDEX IF NOT EXISTS idx_refunds_return_reason_id
    ON refunds (return_reason_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_refunds_return_reason'
    ) THEN
        ALTER TABLE refunds
            ADD CONSTRAINT fk_refunds_return_reason
            FOREIGN KEY (return_reason_id)
            REFERENCES return_reasons (id)
            ON DELETE SET NULL;
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS refund_evidences (
    id UUID PRIMARY KEY NOT NULL,
    refund_id UUID NOT NULL,
    object_key TEXT NOT NULL,
    mime_type VARCHAR(120),
    file_size_bytes BIGINT,
    media_type VARCHAR(16) NOT NULL,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_refund_evidences_refund FOREIGN KEY (refund_id) REFERENCES refunds (id) ON DELETE CASCADE,
    CONSTRAINT ck_refund_evidences_media_type CHECK (media_type IN ('image', 'video')),
    CONSTRAINT ck_refund_evidences_file_size_non_negative CHECK (file_size_bytes IS NULL OR file_size_bytes >= 0)
);

CREATE INDEX IF NOT EXISTS idx_refund_evidences_refund_sort
    ON refund_evidences (refund_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_refund_evidences_media_type
    ON refund_evidences (media_type);

CREATE UNIQUE INDEX IF NOT EXISTS uq_refund_evidences_single_video
    ON refund_evidences (refund_id)
    WHERE media_type = 'video';
