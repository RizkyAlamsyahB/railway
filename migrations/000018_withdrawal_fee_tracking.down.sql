DROP INDEX IF EXISTS idx_vendor_withdrawals_xendit_payout_id;
DROP INDEX IF EXISTS idx_vendor_withdrawals_fee_status;

ALTER TABLE vendor_withdrawals
    DROP CONSTRAINT IF EXISTS ck_vendor_withdrawals_fee_status,
    DROP CONSTRAINT IF EXISTS ck_vendor_withdrawals_total_deducted,
    DROP CONSTRAINT IF EXISTS ck_vendor_withdrawals_amount_net,
    DROP CONSTRAINT IF EXISTS ck_vendor_withdrawals_fee_actual,
    DROP CONSTRAINT IF EXISTS ck_vendor_withdrawals_fee_estimated;

ALTER TABLE vendor_withdrawals
    DROP COLUMN IF EXISTS xendit_transaction_id,
    DROP COLUMN IF EXISTS fee_status,
    DROP COLUMN IF EXISTS total_deducted,
    DROP COLUMN IF EXISTS amount_net,
    DROP COLUMN IF EXISTS fee_actual,
    DROP COLUMN IF EXISTS fee_estimated;
