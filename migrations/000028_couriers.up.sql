-- 000028_couriers.up.sql
-- Master couriers table + vendor_couriers join table.

CREATE TABLE IF NOT EXISTS couriers (
    id          SERIAL PRIMARY KEY,
    code        VARCHAR(20)  NOT NULL UNIQUE,
    name        VARCHAR(100) NOT NULL,
    logo_url    VARCHAR(255),
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Seed all RajaOngkir Komerce supported domestic couriers.
INSERT INTO couriers (code, name) VALUES
    ('jne',     'JNE'),
    ('sicepat', 'SiCepat'),
    ('ide',     'IDExpress'),
    ('sap',     'SAP Express'),
    ('ninja',   'Ninja'),
    ('jnt',     'J&T Express'),
    ('tiki',    'TIKI'),
    ('wahana',  'Wahana Express'),
    ('pos',     'POS Indonesia'),
    ('sentral', 'Sentral Cargo'),
    ('lion',    'Lion Parcel'),
    ('rex',     'Royal Express Asia')
ON CONFLICT (code) DO NOTHING;

-- Join table: which couriers a vendor supports.
CREATE TABLE IF NOT EXISTS vendor_couriers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id   UUID    NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    courier_id  INT     NOT NULL REFERENCES couriers(id) ON DELETE CASCADE,
    is_active   BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (vendor_id, courier_id)
);

CREATE INDEX idx_vendor_couriers_vendor ON vendor_couriers (vendor_id);
