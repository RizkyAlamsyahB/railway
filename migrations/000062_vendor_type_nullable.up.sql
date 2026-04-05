ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_vendor_type;

ALTER TABLE vendors
    ALTER COLUMN vendor_type DROP NOT NULL;

ALTER TABLE vendors
    ADD CONSTRAINT ck_vendors_vendor_type CHECK (
        vendor_type IS NULL OR vendor_type IN (
            'souvenir_store',
            'ppiu',
            'hajj_dormitory'
        )
    );
