DROP INDEX IF EXISTS idx_vendor_bank_accounts_vendor_default;
DROP INDEX IF EXISTS idx_vendor_bank_accounts_vendor_id;

WITH ranked AS (
    SELECT
        id,
        vendor_id,
        ROW_NUMBER() OVER (
            PARTITION BY vendor_id
            ORDER BY is_default DESC, created_at DESC, id DESC
        ) AS rn
    FROM vendor_bank_accounts
)
DELETE FROM vendor_bank_accounts v
USING ranked r
WHERE v.id = r.id
  AND r.rn > 1;

ALTER TABLE vendor_bank_accounts
    DROP COLUMN IF EXISTS is_default,
    DROP COLUMN IF EXISTS account_last4,
    DROP COLUMN IF EXISTS channel_code;

ALTER TABLE vendor_bank_accounts
    ADD CONSTRAINT vendor_bank_accounts_vendor_id_key UNIQUE (vendor_id);
