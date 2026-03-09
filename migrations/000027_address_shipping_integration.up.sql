-- 000027_address_shipping_integration.up.sql
-- Extends addresses table for full RajaOngkir hierarchical location data
-- and adds timestamps.

ALTER TABLE addresses
    ADD COLUMN IF NOT EXISTS province_name   VARCHAR(100),
    ADD COLUMN IF NOT EXISTS city_name        VARCHAR(100),
    ADD COLUMN IF NOT EXISTS district_name    VARCHAR(100),
    ADD COLUMN IF NOT EXISTS subdistrict_id   VARCHAR(20),
    ADD COLUMN IF NOT EXISTS subdistrict_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS notes            TEXT,
    ADD COLUMN IF NOT EXISTS latitude         DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS longitude        DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    ADD COLUMN IF NOT EXISTS updated_at       TIMESTAMPTZ NOT NULL DEFAULT now();

-- Index for user's default address lookup
CREATE INDEX IF NOT EXISTS idx_addresses_user_default ON addresses (user_id, is_default) WHERE is_default = true;
