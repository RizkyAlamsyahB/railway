ALTER TABLE vendor_bank_accounts
    DROP CONSTRAINT IF EXISTS vendor_bank_accounts_vendor_id_key;

ALTER TABLE vendor_bank_accounts
    ADD COLUMN IF NOT EXISTS channel_code VARCHAR(40) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS account_last4 VARCHAR(4) NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS is_default BOOLEAN NOT NULL DEFAULT false;

UPDATE vendor_bank_accounts
SET is_default = true
WHERE is_default = false;

CREATE INDEX IF NOT EXISTS idx_vendor_bank_accounts_vendor_id ON vendor_bank_accounts (vendor_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_vendor_bank_accounts_vendor_default
    ON vendor_bank_accounts (vendor_id)
    WHERE is_default = true;
