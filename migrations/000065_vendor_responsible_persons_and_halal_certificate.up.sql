CREATE TABLE vendor_responsible_persons (
    id UUID PRIMARY KEY NOT NULL,
    vendor_id UUID NOT NULL UNIQUE,
    user_id UUID NOT NULL UNIQUE,
    nik TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_vendor_responsible_persons_vendor FOREIGN KEY (vendor_id) REFERENCES vendors(id) ON DELETE CASCADE,
    CONSTRAINT fk_vendor_responsible_persons_user FOREIGN KEY (user_id) REFERENCES users(id)
);

ALTER TABLE vendor_documents
    DROP CONSTRAINT IF EXISTS ck_vendor_documents_doc_type;

ALTER TABLE vendor_documents
    ADD CONSTRAINT ck_vendor_documents_doc_type CHECK (
        doc_type IN (
            'owner_document_id',
            'business_nib',
            'business_npwp',
            'halal_certificate',
            'store_photo',
            'bank_account_proof',
            'business_logo',
            'business_banner'
        )
    );
