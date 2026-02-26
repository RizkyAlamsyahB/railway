CREATE TABLE payment_invoices (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID NOT NULL,
    gateway VARCHAR(20) NOT NULL DEFAULT 'xendit',
    xendit_invoice_id VARCHAR(100),
    external_invoice_id VARCHAR(100) UNIQUE NOT NULL,
    invoice_url TEXT,
    payment_method VARCHAR(30),
    payment_channel VARCHAR(30),
    amount NUMERIC(18,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'IDR',
    status VARCHAR(16) NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    raw_payload JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_payment_invoices_status CHECK (status IN ('pending', 'paid', 'expired', 'failed')),
    CONSTRAINT fk_payment_invoices_order FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE TABLE payment_events (
    id UUID PRIMARY KEY NOT NULL,
    payment_invoice_id UUID NOT NULL,
    event_type VARCHAR(50) NOT NULL,
    external_event_id VARCHAR(120) UNIQUE NOT NULL,
    payload JSONB NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_payment_events_invoice FOREIGN KEY (payment_invoice_id) REFERENCES payment_invoices (id)
);

CREATE TABLE refunds (
    id UUID PRIMARY KEY NOT NULL,
    order_id UUID NOT NULL,
    payment_invoice_id UUID NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    reason TEXT,
    status VARCHAR(16) NOT NULL DEFAULT 'requested',
    requested_by UUID NOT NULL,
    processed_by UUID,
    processed_at TIMESTAMPTZ,
    CONSTRAINT ck_refunds_status CHECK (status IN ('requested', 'approved', 'rejected', 'processed')),
    CONSTRAINT fk_refunds_order FOREIGN KEY (order_id) REFERENCES orders (id),
    CONSTRAINT fk_refunds_invoice FOREIGN KEY (payment_invoice_id) REFERENCES payment_invoices (id),
    CONSTRAINT fk_refunds_requested_by FOREIGN KEY (requested_by) REFERENCES users (id),
    CONSTRAINT fk_refunds_processed_by FOREIGN KEY (processed_by) REFERENCES users (id)
);
