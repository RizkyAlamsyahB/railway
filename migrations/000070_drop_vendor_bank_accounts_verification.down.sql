ALTER TABLE vendor_bank_accounts
    ADD COLUMN IF NOT EXISTS verification_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS rejection_reason TEXT,
    ADD COLUMN IF NOT EXISTS verified_by UUID,
    ADD COLUMN IF NOT EXISTS verified_at TIMESTAMPTZ;

ALTER TABLE vendor_bank_accounts
    DROP CONSTRAINT IF EXISTS ck_vendor_bank_accounts_verification_status;

ALTER TABLE vendor_bank_accounts
    ADD CONSTRAINT ck_vendor_bank_accounts_verification_status
    CHECK (verification_status IN ('pending', 'verified', 'rejected'));

ALTER TABLE vendor_bank_accounts
    DROP CONSTRAINT IF EXISTS ck_vendor_bank_accounts_verified_consistency;

ALTER TABLE vendor_bank_accounts
    ADD CONSTRAINT ck_vendor_bank_accounts_verified_consistency CHECK (
        (verification_status = 'verified' AND verified_by IS NOT NULL AND verified_at IS NOT NULL)
        OR
        (verification_status IN ('pending', 'rejected') AND verified_by IS NULL AND verified_at IS NULL)
    );

ALTER TABLE vendor_bank_accounts
    DROP CONSTRAINT IF EXISTS fk_vendor_bank_accounts_verified_by;

ALTER TABLE vendor_bank_accounts
    ADD CONSTRAINT fk_vendor_bank_accounts_verified_by
    FOREIGN KEY (verified_by) REFERENCES users (id);

CREATE INDEX IF NOT EXISTS idx_vendor_bank_accounts_status ON vendor_bank_accounts (verification_status);
