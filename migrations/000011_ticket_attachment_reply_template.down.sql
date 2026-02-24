-- ============================================================
-- ROLLBACK 000011
-- ============================================================

ALTER TABLE reply_templates DROP CONSTRAINT IF EXISTS ck_reply_templates_category;
ALTER TABLE reply_templates DROP CONSTRAINT IF EXISTS uq_reply_templates_shortcut;
DROP INDEX IF EXISTS idx_reply_templates_category;
DROP INDEX IF EXISTS idx_reply_templates_shortcut;
ALTER TABLE reply_templates
    DROP COLUMN IF EXISTS shortcut,
    DROP COLUMN IF EXISTS category;

ALTER TABLE tickets
    DROP COLUMN IF EXISTS attachment_url,
    DROP COLUMN IF EXISTS attachment_content_type;
