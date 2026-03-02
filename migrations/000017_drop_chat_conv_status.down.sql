-- Re-add the status column to chat_conversations.
ALTER TABLE chat_conversations ADD COLUMN status VARCHAR(10) NOT NULL DEFAULT 'open';
ALTER TABLE chat_conversations ADD CONSTRAINT ck_chat_conv_status CHECK (status IN ('open', 'closed'));
CREATE INDEX idx_chat_conv_status ON chat_conversations (status);
