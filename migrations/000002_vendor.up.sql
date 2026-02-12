CREATE TABLE vendors (
    id UUID PRIMARY KEY NOT NULL,
    owner_user_id UUID NOT NULL,
    vendor_type VARCHAR(32) NOT NULL DEFAULT 'general_souvenir_store',
    legal_name VARCHAR(160),
    display_name VARCHAR(120) NOT NULL,
    responsible_person_name VARCHAR(120) NOT NULL,
    description TEXT,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    status_reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_vendors_vendor_type CHECK (
        vendor_type IN (
            'umrah_souvenir_store',
            'hajj_souvenir_store',
            'general_souvenir_store'
        )
    ),
    CONSTRAINT ck_vendors_status CHECK (status IN ('draft', 'submitted', 'active', 'rejected', 'blocked')),
    CONSTRAINT uq_vendors_owner_user UNIQUE (owner_user_id),
    CONSTRAINT fk_vendors_owner FOREIGN KEY (owner_user_id) REFERENCES users (id),
    CONSTRAINT fk_vendors_approved_by FOREIGN KEY (approved_by) REFERENCES users (id)
);

CREATE TABLE vendor_documents (
    id UUID PRIMARY KEY NOT NULL,
    vendor_id UUID NOT NULL,
    doc_type VARCHAR(40) NOT NULL,
    file_url TEXT NOT NULL,
    mime_type VARCHAR(100),
    file_size_bytes INT,
    file_checksum VARCHAR(128),
    uploaded_by UUID,
    verification_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    rejection_reason TEXT,
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_vendor_documents_doc_type CHECK (
        doc_type IN (
            'owner_ktp',
            'owner_passport',
            'business_npwp',
            'store_photo',
            'bank_account_proof',
            'business_logo',
            'business_banner'
        )
    ),
    CONSTRAINT ck_vendor_documents_verification_status CHECK (verification_status IN ('pending', 'verified', 'rejected')),
    CONSTRAINT ck_vendor_documents_verified_consistency CHECK (
        (verification_status = 'verified' AND verified_by IS NOT NULL AND verified_at IS NOT NULL)
        OR
        (verification_status IN ('pending', 'rejected') AND verified_by IS NULL AND verified_at IS NULL)
    ),
    CONSTRAINT fk_vendor_documents_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id),
    CONSTRAINT fk_vendor_documents_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES users (id),
    CONSTRAINT fk_vendor_documents_verified_by FOREIGN KEY (verified_by) REFERENCES users (id)
);

CREATE UNIQUE INDEX uq_vendor_documents_vendor_doc_type ON vendor_documents (vendor_id, doc_type);
CREATE INDEX idx_vendor_documents_vendor_status ON vendor_documents (vendor_id, verification_status);

CREATE TABLE vendor_bank_accounts (
    id UUID PRIMARY KEY NOT NULL,
    vendor_id UUID NOT NULL UNIQUE,
    bank_name VARCHAR(120) NOT NULL,
    account_number VARCHAR(60) NOT NULL,
    account_holder_name VARCHAR(160) NOT NULL,
    verification_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    rejection_reason TEXT,
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_vendor_bank_accounts_verification_status CHECK (verification_status IN ('pending', 'verified', 'rejected')),
    CONSTRAINT ck_vendor_bank_accounts_verified_consistency CHECK (
        (verification_status = 'verified' AND verified_by IS NOT NULL AND verified_at IS NOT NULL)
        OR
        (verification_status IN ('pending', 'rejected') AND verified_by IS NULL AND verified_at IS NULL)
    ),
    CONSTRAINT fk_vendor_bank_accounts_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id),
    CONSTRAINT fk_vendor_bank_accounts_verified_by FOREIGN KEY (verified_by) REFERENCES users (id)
);

CREATE INDEX idx_vendor_bank_accounts_status ON vendor_bank_accounts (verification_status);
