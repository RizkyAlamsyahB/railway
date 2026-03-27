ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_business_legal_type;

ALTER TABLE vendors
    DROP COLUMN IF EXISTS business_legal_type;
