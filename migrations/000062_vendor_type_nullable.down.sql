UPDATE vendors
SET vendor_type = 'souvenir_store'
WHERE vendor_type IS NULL;

ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_vendor_type;

ALTER TABLE vendors
    ALTER COLUMN vendor_type SET NOT NULL;

ALTER TABLE vendors
    ADD CONSTRAINT ck_vendors_vendor_type CHECK (
        vendor_type IN (
            'souvenir_store',
            'ppiu',
            'hajj_dormitory'
        )
    );
