ALTER TABLE vendor_documents ADD COLUMN rejection_reason TEXT;
ALTER TABLE vendor_documents ADD COLUMN verification_status VARCHAR(16) NOT NULL DEFAULT 'pending';
ALTER TABLE vendor_documents ADD CONSTRAINT ck_vendor_documents_verification_status
    CHECK (verification_status IN ('pending', 'verified', 'rejected'));
ALTER TABLE vendor_documents ADD CONSTRAINT ck_vendor_documents_verified_consistency CHECK (
    (verification_status = 'verified' AND verified_by IS NOT NULL AND verified_at IS NOT NULL)
    OR
    (verification_status IN ('pending', 'rejected') AND verified_by IS NULL AND verified_at IS NULL)
);
CREATE INDEX idx_vendor_documents_vendor_status ON vendor_documents (vendor_id, verification_status);
