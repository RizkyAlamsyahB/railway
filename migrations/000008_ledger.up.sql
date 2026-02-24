CREATE TABLE ledger_accounts (
    id UUID PRIMARY KEY NOT NULL,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(120) NOT NULL,
    account_type VARCHAR(16) NOT NULL,
    normal_side VARCHAR(1) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    CONSTRAINT ck_ledger_accounts_account_type CHECK (account_type IN ('asset', 'liability', 'equity', 'revenue', 'expense')),
    CONSTRAINT ck_ledger_accounts_normal_side CHECK (normal_side IN ('D', 'C'))
);

CREATE TABLE ledger_journals (
    id UUID PRIMARY KEY NOT NULL,
    journal_no VARCHAR(40) UNIQUE NOT NULL,
    source_type VARCHAR(24) NOT NULL,
    source_id UUID NOT NULL,
    event_time TIMESTAMPTZ NOT NULL,
    description TEXT,
    status VARCHAR(16) NOT NULL DEFAULT 'posted',
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_ledger_journals_source_type CHECK (source_type IN ('order', 'payment_invoice', 'refund', 'payout_batch', 'manual')),
    CONSTRAINT ck_ledger_journals_status CHECK (status IN ('posted', 'reversed')),
    CONSTRAINT fk_ledger_journals_created_by FOREIGN KEY (created_by) REFERENCES users (id)
);

CREATE TABLE ledger_lines (
    id UUID PRIMARY KEY NOT NULL,
    journal_id UUID NOT NULL,
    account_id UUID NOT NULL,
    debit NUMERIC(18,2) NOT NULL DEFAULT 0,
    credit NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    reference VARCHAR(80),
    CONSTRAINT fk_ledger_lines_journal FOREIGN KEY (journal_id) REFERENCES ledger_journals (id),
    CONSTRAINT fk_ledger_lines_account FOREIGN KEY (account_id) REFERENCES ledger_accounts (id)
);
