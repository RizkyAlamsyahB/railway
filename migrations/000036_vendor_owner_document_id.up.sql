ALTER TABLE vendor_documents
    DROP CONSTRAINT IF EXISTS ck_vendor_documents_doc_type;

DELETE FROM vendor_documents
WHERE doc_type IN ('owner_ktp', 'owner_passport');

ALTER TABLE vendor_documents
    ADD CONSTRAINT ck_vendor_documents_doc_type CHECK (
        doc_type IN (
            'owner_document_id',
            'business_npwp',
            'store_photo',
            'bank_account_proof',
            'business_logo',
            'business_banner'
        )
    );
