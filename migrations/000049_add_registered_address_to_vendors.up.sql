ALTER TABLE vendors
    ADD COLUMN IF NOT EXISTS registered_address TEXT;
