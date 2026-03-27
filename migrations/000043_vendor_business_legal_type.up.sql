ALTER TABLE vendors
    ADD COLUMN IF NOT EXISTS business_legal_type VARCHAR(32);

ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_business_legal_type;

ALTER TABLE vendors
    ADD CONSTRAINT ck_vendors_business_legal_type CHECK (
        business_legal_type IS NULL OR business_legal_type IN ('perorangan', 'korporasi')
    );
