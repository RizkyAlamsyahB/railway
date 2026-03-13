CREATE TABLE IF NOT EXISTS vendor_onboardings (
    id UUID PRIMARY KEY NOT NULL,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(32) NOT NULL,
    password_hash TEXT,
    store_name VARCHAR(120),
    vendor_type VARCHAR(32),
    business_legal_type VARCHAR(32),
    document_id_type VARCHAR(16),
    nik VARCHAR(32),
    owner_name VARCHAR(120),
    birth_date DATE,
    document_id_object_key TEXT,
    otp_verified_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_vendor_onboardings_email UNIQUE (email),
    CONSTRAINT ck_vendor_onboardings_status CHECK (
        status IN ('otp_verified', 'password_set', 'store_info_completed', 'completed')
    ),
    CONSTRAINT ck_vendor_onboardings_vendor_type CHECK (
        vendor_type IS NULL OR vendor_type IN ('souvenir_store', 'ppiu', 'hajj_dormitory')
    ),
    CONSTRAINT ck_vendor_onboardings_business_legal_type CHECK (
        business_legal_type IS NULL OR business_legal_type IN ('perorangan', 'korporasi')
    ),
    CONSTRAINT ck_vendor_onboardings_document_id_type CHECK (
        document_id_type IS NULL OR document_id_type IN ('ktp', 'passport')
    )
);
