ALTER TABLE refunds DROP COLUMN IF EXISTS user_bank_account_id;
ALTER TABLE refunds DROP COLUMN IF EXISTS xendit_refund_id;
ALTER TABLE refunds DROP COLUMN IF EXISTS refund_method;

DROP INDEX IF EXISTS idx_user_bank_accounts_default;
DROP INDEX IF EXISTS idx_user_bank_accounts_user_id;
DROP TABLE IF EXISTS user_bank_accounts;
