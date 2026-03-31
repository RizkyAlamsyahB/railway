DROP INDEX IF EXISTS idx_vendor_documents_vendor_status;
ALTER TABLE vendor_documents DROP CONSTRAINT IF EXISTS ck_vendor_documents_verified_consistency;
ALTER TABLE vendor_documents DROP CONSTRAINT IF EXISTS ck_vendor_documents_verification_status;
ALTER TABLE vendor_documents DROP COLUMN IF EXISTS verification_status;
ALTER TABLE vendor_documents DROP COLUMN IF EXISTS rejection_reason;
