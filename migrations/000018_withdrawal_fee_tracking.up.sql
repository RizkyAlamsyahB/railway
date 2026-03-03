ALTER TABLE vendor_withdrawals
    ADD COLUMN IF NOT EXISTS fee_estimated NUMERIC(18,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS fee_actual NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS amount_net NUMERIC(18,2),
    ADD COLUMN IF NOT EXISTS total_deducted NUMERIC(18,2) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS fee_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    ADD COLUMN IF NOT EXISTS xendit_transaction_id VARCHAR(128);

-- Initialize new fields for existing rows without changing historical balances.
UPDATE vendor_withdrawals
SET amount_net = COALESCE(amount_net, amount),
    total_deducted = CASE WHEN total_deducted = 0 THEN amount ELSE total_deducted END,
    fee_status = CASE
        WHEN status IN ('completed', 'failed') THEN 'resolved'
        ELSE fee_status
    END;

ALTER TABLE vendor_withdrawals
    ADD CONSTRAINT ck_vendor_withdrawals_fee_estimated CHECK (fee_estimated >= 0),
    ADD CONSTRAINT ck_vendor_withdrawals_fee_actual CHECK (fee_actual IS NULL OR fee_actual >= 0),
    ADD CONSTRAINT ck_vendor_withdrawals_amount_net CHECK (amount_net IS NULL OR amount_net >= 0),
    ADD CONSTRAINT ck_vendor_withdrawals_total_deducted CHECK (total_deducted >= 0),
    ADD CONSTRAINT ck_vendor_withdrawals_fee_status CHECK (fee_status IN ('pending', 'resolved', 'failed'));

CREATE INDEX IF NOT EXISTS idx_vendor_withdrawals_fee_status ON vendor_withdrawals (fee_status);
CREATE INDEX IF NOT EXISTS idx_vendor_withdrawals_xendit_payout_id ON vendor_withdrawals (xendit_payout_id);
