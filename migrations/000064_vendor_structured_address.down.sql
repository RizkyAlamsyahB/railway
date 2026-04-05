ALTER TABLE vendors
    DROP COLUMN IF EXISTS province_id,
    DROP COLUMN IF EXISTS province_name,
    DROP COLUMN IF EXISTS city_id,
    DROP COLUMN IF EXISTS city_name,
    DROP COLUMN IF EXISTS district_id,
    DROP COLUMN IF EXISTS district_name,
    DROP COLUMN IF EXISTS subdistrict_id,
    DROP COLUMN IF EXISTS subdistrict_name,
    DROP COLUMN IF EXISTS postal_code,
    DROP COLUMN IF EXISTS address_line,
    ADD COLUMN IF NOT EXISTS registered_address TEXT;
