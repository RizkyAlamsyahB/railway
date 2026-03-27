ALTER TABLE vendor_onboardings
    DROP COLUMN IF EXISTS nib,
    DROP COLUMN IF EXISTS company_name,
    DROP COLUMN IF EXISTS established_date,
    DROP COLUMN IF EXISTS registered_address,
    DROP COLUMN IF EXISTS nib_document_object_key;
