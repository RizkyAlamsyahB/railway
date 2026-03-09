-- 000027_address_shipping_integration.down.sql
-- Reverts the address table extensions.

DROP INDEX IF EXISTS idx_addresses_user_default;

ALTER TABLE addresses
    DROP COLUMN IF EXISTS province_name,
    DROP COLUMN IF EXISTS city_name,
    DROP COLUMN IF EXISTS district_name,
    DROP COLUMN IF EXISTS subdistrict_id,
    DROP COLUMN IF EXISTS subdistrict_name,
    DROP COLUMN IF EXISTS notes,
    DROP COLUMN IF EXISTS latitude,
    DROP COLUMN IF EXISTS longitude,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_at;
