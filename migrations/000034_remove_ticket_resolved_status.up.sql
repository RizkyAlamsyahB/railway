-- Migrate existing "resolved" tickets to "closed".
-- Copy resolved_at → closed_at for tickets that were resolved but not yet closed.
UPDATE tickets
SET    status    = 'closed',
       closed_at = COALESCE(closed_at, resolved_at, NOW()),
       updated_at = NOW()
WHERE  status = 'resolved';

-- Drop the resolved_at column.
ALTER TABLE tickets DROP COLUMN IF EXISTS resolved_at;

-- Update ticket_status_logs that reference "resolved".
UPDATE ticket_status_logs SET new_status = 'closed' WHERE new_status = 'resolved';
UPDATE ticket_status_logs SET old_status = 'closed' WHERE old_status = 'resolved';
