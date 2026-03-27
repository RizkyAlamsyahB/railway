ALTER TABLE vendor_onboardings
    ADD COLUMN IF NOT EXISTS nib TEXT,
    ADD COLUMN IF NOT EXISTS company_name VARCHAR(120),
    ADD COLUMN IF NOT EXISTS established_date DATE,
    ADD COLUMN IF NOT EXISTS registered_address TEXT,
    ADD COLUMN IF NOT EXISTS nib_document_object_key TEXT;
