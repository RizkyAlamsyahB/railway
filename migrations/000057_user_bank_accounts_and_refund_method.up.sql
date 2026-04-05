-- user_bank_accounts: stores customer bank accounts for refund disbursements
CREATE TABLE IF NOT EXISTS user_bank_accounts (
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel_code    VARCHAR(40)  NOT NULL,          -- e.g. ID_BCA, ID_MANDIRI
    bank_name       VARCHAR(80)  NOT NULL,          -- display name: "BCA", "Mandiri"
    account_number  TEXT         NOT NULL,           -- AES-256-GCM encrypted
    account_holder_name VARCHAR(120) NOT NULL,
    account_last4   VARCHAR(4)   NOT NULL,
    is_default      BOOLEAN      NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_user_bank_accounts_user_id ON user_bank_accounts(user_id);

-- Ensure at most one default per user
CREATE UNIQUE INDEX idx_user_bank_accounts_default
    ON user_bank_accounts(user_id) WHERE is_default = true;

-- refund_method tracks how the refund will be processed
ALTER TABLE refunds ADD COLUMN IF NOT EXISTS refund_method VARCHAR(30);
-- xendit_refund_id is set when using Xendit Refund API (QR/E-Wallet)
ALTER TABLE refunds ADD COLUMN IF NOT EXISTS xendit_refund_id VARCHAR(100);
-- user_bank_account_id links to the saved bank account used for VA disbursement
ALTER TABLE refunds ADD COLUMN IF NOT EXISTS user_bank_account_id UUID REFERENCES user_bank_accounts(id);

-- Add awaiting_destination status support (for VA cancel/reject flow)
-- No DDL needed — status column is varchar, we just start using the new value.

-- Add payment_status 'refunded' support on orders
-- No DDL needed — payment_status is varchar.
