CREATE TABLE vendors (
    id UUID PRIMARY KEY NOT NULL,
    owner_user_id UUID NOT NULL,
    vendor_type VARCHAR(32) NOT NULL DEFAULT 'souvenir_store',
    legal_name VARCHAR(160) NOT NULL,
    display_name VARCHAR(120) NOT NULL,
    description TEXT,
    status VARCHAR(16) NOT NULL DEFAULT 'draft',
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_vendors_vendor_type CHECK (vendor_type IN ('souvenir_store')),
    CONSTRAINT ck_vendors_status CHECK (status IN ('draft', 'submitted', 'active', 'rejected', 'blocked')),
    CONSTRAINT fk_vendors_owner FOREIGN KEY (owner_user_id) REFERENCES users (id),
    CONSTRAINT fk_vendors_approved_by FOREIGN KEY (approved_by) REFERENCES users (id)
);

CREATE TABLE vendor_documents (
    id UUID PRIMARY KEY NOT NULL,
    vendor_id UUID NOT NULL,
    doc_type VARCHAR(40) NOT NULL,
    file_url TEXT NOT NULL,
    verification_status VARCHAR(16) NOT NULL DEFAULT 'pending',
    verified_by UUID,
    verified_at TIMESTAMPTZ,
    CONSTRAINT ck_vendor_documents_verification_status CHECK (verification_status IN ('pending', 'verified', 'rejected')),
    CONSTRAINT fk_vendor_documents_vendor FOREIGN KEY (vendor_id) REFERENCES vendors (id),
    CONSTRAINT fk_vendor_documents_verified_by FOREIGN KEY (verified_by) REFERENCES users (id)
);
