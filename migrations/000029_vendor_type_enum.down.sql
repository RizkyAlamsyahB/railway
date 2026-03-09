ALTER TABLE vendors
    DROP CONSTRAINT IF EXISTS ck_vendors_vendor_type;

UPDATE vendors
SET vendor_type = 'general_souvenir_store'
WHERE vendor_type IN ('souvenir_store', 'ppiu', 'hajj_dormitory');

ALTER TABLE vendors
    ADD CONSTRAINT ck_vendors_vendor_type CHECK (
        vendor_type IN (
            'umrah_souvenir_store',
            'hajj_souvenir_store',
            'general_souvenir_store'
        )
    );

ALTER TABLE vendors
    ALTER COLUMN vendor_type SET DEFAULT 'general_souvenir_store';
