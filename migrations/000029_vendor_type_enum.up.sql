ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_vendor_type;

UPDATE vendors
SET vendor_type = 'souvenir_store'
WHERE vendor_type IN (
    'general_souvenir_store',
    'umrah_souvenir_store',
    'hajj_souvenir_store'
);

ALTER TABLE vendors
    ADD CONSTRAINT ck_vendors_vendor_type CHECK (
        vendor_type IN (
            'souvenir_store',
            'ppiu',
            'hajj_dormitory'
        )
    );

ALTER TABLE vendors
    ALTER COLUMN vendor_type DROP DEFAULT;
