-- Reverse: remove subject_id from tickets and drop ticket_subjects table

DROP INDEX IF EXISTS idx_tickets_subject_id;

ALTER TABLE tickets DROP COLUMN IF EXISTS subject_id;

DROP TABLE IF EXISTS ticket_subjects;
