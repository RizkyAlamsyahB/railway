-- Drop the status column and its constraint/index from chat_conversations.
-- The column was always 'open' and never transitioned to 'closed'.

DROP INDEX IF EXISTS idx_chat_conv_status;
ALTER TABLE chat_conversations DROP CONSTRAINT IF EXISTS ck_chat_conv_status;
ALTER TABLE chat_conversations DROP COLUMN IF EXISTS status;
