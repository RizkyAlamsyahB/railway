UPDATE vendors
SET display_name = 'vendor-' || LEFT(id::text, 8)
WHERE display_name IS NULL;

ALTER TABLE vendors
    ALTER COLUMN display_name SET NOT NULL;
