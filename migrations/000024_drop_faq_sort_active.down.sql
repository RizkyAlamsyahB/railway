-- Restore sort_order and is_active on faqs table
ALTER TABLE faqs ADD COLUMN sort_order INT NOT NULL DEFAULT 0;
ALTER TABLE faqs ADD COLUMN is_active  BOOLEAN NOT NULL DEFAULT TRUE;
CREATE INDEX idx_faqs_active ON faqs (is_active) WHERE is_active;
