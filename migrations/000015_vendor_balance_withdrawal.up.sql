-- Vendor balances: tracks real-time available/pending balance per vendor.
CREATE TABLE vendor_balances (
    vendor_id       UUID PRIMARY KEY NOT NULL,
    available_balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    pending_balance   NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_earned      NUMERIC(18,2) NOT NULL DEFAULT 0,
    total_withdrawn   NUMERIC(18,2) NOT NULL DEFAULT 0,
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_vendor_balances_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id),
    CONSTRAINT ck_vendor_balances_available CHECK (available_balance >= 0),
    CONSTRAINT ck_vendor_balances_pending CHECK (pending_balance >= 0)
);

-- Vendor withdrawals: vendor-initiated payout requests (replaces admin-batch flow).
CREATE TABLE vendor_withdrawals (
    id               UUID PRIMARY KEY NOT NULL,
    vendor_id        UUID NOT NULL,
    amount           NUMERIC(18,2) NOT NULL,
    channel_code     VARCHAR(32) NOT NULL,
    status           VARCHAR(16) NOT NULL DEFAULT 'pending',
    xendit_payout_id VARCHAR(128),
    xendit_status    VARCHAR(32),
    description      TEXT,
    failed_reason    TEXT,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_vendor_withdrawals_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id),
    CONSTRAINT ck_vendor_withdrawals_status CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    CONSTRAINT ck_vendor_withdrawals_amount CHECK (amount >= 10000)
);

CREATE INDEX idx_vendor_withdrawals_vendor_id ON vendor_withdrawals (vendor_id);
CREATE INDEX idx_vendor_withdrawals_status ON vendor_withdrawals (status);

-- Initialize balance rows for all existing vendors.
INSERT INTO vendor_balances (vendor_id, available_balance, pending_balance, total_earned, total_withdrawn, updated_at)
SELECT id, 0, 0, 0, 0, NOW()
FROM vendors
ON CONFLICT (vendor_id) DO NOTHING;
