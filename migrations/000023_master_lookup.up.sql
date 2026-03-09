-- Master lookup tables for admin settings

-- Return reasons (alasan return)
CREATE TABLE return_reasons (
    id         SERIAL       PRIMARY KEY,
    reason     VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Admin contacts (kontak admin)
CREATE TABLE admin_contacts (
    id         SERIAL       PRIMARY KEY,
    content    TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);
