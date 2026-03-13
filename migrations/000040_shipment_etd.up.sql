-- Add etd (estimated time of delivery) column to shipments table.
-- Stores the ETD string from RajaOngkir at checkout time (e.g. "1-2 day", "3 day").
ALTER TABLE shipments ADD COLUMN IF NOT EXISTS etd VARCHAR(30) DEFAULT '';
