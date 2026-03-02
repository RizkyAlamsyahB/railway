-- FAQ / Pusat Bantuan
CREATE TABLE faqs (
    id            SERIAL       PRIMARY KEY,
    category      VARCHAR(60)  NOT NULL,
    question      TEXT         NOT NULL,
    answer        TEXT         NOT NULL,
    sort_order    INT          NOT NULL DEFAULT 0,
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX idx_faqs_category ON faqs (category);
CREATE INDEX idx_faqs_active   ON faqs (is_active) WHERE is_active;
