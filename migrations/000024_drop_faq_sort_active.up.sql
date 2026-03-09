-- Remove sort_order and is_active from faqs table
DROP INDEX IF EXISTS idx_faqs_active;
ALTER TABLE faqs DROP COLUMN IF EXISTS sort_order;
ALTER TABLE faqs DROP COLUMN IF EXISTS is_active;
