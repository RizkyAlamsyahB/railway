-- ============================================================
-- MIGRATION 000011: ticket attachment + reply_template shortcut/category
-- ============================================================

-- 1. Tambah kolom attachment ke tickets
ALTER TABLE tickets
    ADD COLUMN IF NOT EXISTS attachment_url          TEXT,
    ADD COLUMN IF NOT EXISTS attachment_content_type VARCHAR(50);

-- 2. Tambah shortcut & category ke reply_templates
ALTER TABLE reply_templates
    ADD COLUMN IF NOT EXISTS shortcut VARCHAR(100),
    ADD COLUMN IF NOT EXISTS category VARCHAR(50);

-- Set default untuk baris lama (jika ada) supaya NOT NULL bisa diterapkan
UPDATE reply_templates SET shortcut = '/other/' || id WHERE shortcut IS NULL;
UPDATE reply_templates SET category = 'other'           WHERE category IS NULL;

-- Terapkan NOT NULL + UNIQUE + CHECK setelah data bersih
ALTER TABLE reply_templates
    ALTER COLUMN shortcut SET NOT NULL,
    ALTER COLUMN category  SET NOT NULL;

ALTER TABLE reply_templates
    ADD CONSTRAINT uq_reply_templates_shortcut UNIQUE (shortcut);

ALTER TABLE reply_templates
    ADD CONSTRAINT ck_reply_templates_category CHECK (category IN (
        'general','order','payment','refund','shipping','product',
        'account','complaint','return','promo','voucher','technical',
        'verification','vendor','stock','cancellation','delivery',
        'pickup','subscription','review','fraud','other'
    ));

CREATE INDEX IF NOT EXISTS idx_reply_templates_category ON reply_templates (category);
CREATE INDEX IF NOT EXISTS idx_reply_templates_shortcut ON reply_templates (shortcut);
