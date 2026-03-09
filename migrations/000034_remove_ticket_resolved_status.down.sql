-- Re-add the resolved_at column.
ALTER TABLE tickets ADD COLUMN resolved_at TIMESTAMPTZ;
