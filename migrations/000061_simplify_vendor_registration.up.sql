ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_status;

-- Normalize legacy rows before tightening the allowed status values.
UPDATE vendor_onboardings
SET status = 'otp_verified'
WHERE status = 'password_set';

ALTER TABLE vendor_onboardings
    ADD CONSTRAINT ck_vendor_onboardings_status CHECK (
        status IN ('otp_verified', 'completed')
    );

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_vendor_type;

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_business_legal_type;

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_document_id_type;

ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_business_legal_type;

ALTER TABLE vendor_onboardings
    DROP COLUMN IF EXISTS store_name,
    DROP COLUMN IF EXISTS vendor_type,
    DROP COLUMN IF EXISTS business_legal_type,
    DROP COLUMN IF EXISTS document_id_type,
    DROP COLUMN IF EXISTS nik,
    DROP COLUMN IF EXISTS owner_name,
    DROP COLUMN IF EXISTS birth_date,
    DROP COLUMN IF EXISTS document_id_object_key,
    DROP COLUMN IF EXISTS nib,
    DROP COLUMN IF EXISTS company_name,
    DROP COLUMN IF EXISTS established_date,
    DROP COLUMN IF EXISTS registered_address,
    DROP COLUMN IF EXISTS nib_document_object_key;

ALTER TABLE vendors
    DROP COLUMN IF EXISTS business_legal_type;
