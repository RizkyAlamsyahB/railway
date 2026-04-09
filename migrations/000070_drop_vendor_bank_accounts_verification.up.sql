ALTER TABLE vendor_bank_accounts
    DROP CONSTRAINT IF EXISTS fk_vendor_bank_accounts_verified_by;

ALTER TABLE vendor_bank_accounts
    DROP CONSTRAINT IF EXISTS ck_vendor_bank_accounts_verified_consistency;

ALTER TABLE vendor_bank_accounts
    DROP CONSTRAINT IF EXISTS ck_vendor_bank_accounts_verification_status;

DROP INDEX IF EXISTS idx_vendor_bank_accounts_status;

ALTER TABLE vendor_bank_accounts
    DROP COLUMN IF EXISTS verified_at,
    DROP COLUMN IF EXISTS verified_by,
    DROP COLUMN IF EXISTS rejection_reason,
    DROP COLUMN IF EXISTS verification_status;
