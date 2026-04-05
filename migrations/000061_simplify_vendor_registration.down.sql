ALTER TABLE vendors
    ADD COLUMN IF NOT EXISTS business_legal_type VARCHAR(32);

ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_business_legal_type;

ALTER TABLE vendors
    ADD CONSTRAINT ck_vendors_business_legal_type CHECK (
        business_legal_type IS NULL OR business_legal_type IN ('perorangan', 'korporasi')
    );

ALTER TABLE vendor_onboardings
    ADD COLUMN IF NOT EXISTS store_name VARCHAR(120),
    ADD COLUMN IF NOT EXISTS vendor_type VARCHAR(32),
    ADD COLUMN IF NOT EXISTS business_legal_type VARCHAR(32),
    ADD COLUMN IF NOT EXISTS document_id_type VARCHAR(16),
    ADD COLUMN IF NOT EXISTS nik VARCHAR(32),
    ADD COLUMN IF NOT EXISTS owner_name VARCHAR(120),
    ADD COLUMN IF NOT EXISTS birth_date DATE,
    ADD COLUMN IF NOT EXISTS document_id_object_key TEXT,
    ADD COLUMN IF NOT EXISTS nib TEXT,
    ADD COLUMN IF NOT EXISTS company_name VARCHAR(120),
    ADD COLUMN IF NOT EXISTS established_date DATE,
    ADD COLUMN IF NOT EXISTS registered_address TEXT,
    ADD COLUMN IF NOT EXISTS nib_document_object_key TEXT;

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_status;

ALTER TABLE vendor_onboardings
    ADD CONSTRAINT ck_vendor_onboardings_status CHECK (
        status IN ('otp_verified', 'password_set', 'completed')
    );

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_vendor_type;

ALTER TABLE vendor_onboardings
    ADD CONSTRAINT ck_vendor_onboardings_vendor_type CHECK (
        vendor_type IS NULL OR vendor_type IN ('souvenir_store', 'ppiu', 'hajj_dormitory')
    );

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_business_legal_type;

ALTER TABLE vendor_onboardings
    ADD CONSTRAINT ck_vendor_onboardings_business_legal_type CHECK (
        business_legal_type IS NULL OR business_legal_type IN ('perorangan', 'korporasi')
    );

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_document_id_type;

ALTER TABLE vendor_onboardings
    ADD CONSTRAINT ck_vendor_onboardings_document_id_type CHECK (
        document_id_type IS NULL OR document_id_type IN ('ktp', 'passport')
    );
