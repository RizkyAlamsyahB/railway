-- ============================================================
-- MIGRATION: Modul Customer Service
-- ============================================================


-- ============================================================
-- TICKETS
-- ============================================================
CREATE TABLE tickets (
    id              UUID         PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    ticket_number   VARCHAR(20)  NOT NULL UNIQUE,
    customer_id     UUID         NOT NULL,
    assigned_cs_id  UUID,
    order_number    VARCHAR(50)  NOT NULL,
    phone           VARCHAR(20)  NOT NULL,
    reporter_name   VARCHAR(120) NOT NULL,
    subject         VARCHAR(255) NOT NULL,
    detail          TEXT         NOT NULL,
    status          VARCHAR(20)  NOT NULL DEFAULT 'open',
    source          VARCHAR(10)  NOT NULL DEFAULT 'web',
    resolved_at     TIMESTAMPTZ,
    closed_at       TIMESTAMPTZ,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT ck_tickets_status   CHECK (status IN ('open', 'on_progress', 'resolved', 'closed')),
    CONSTRAINT ck_tickets_source   CHECK (source IN ('app', 'web')),
    CONSTRAINT fk_tickets_customer FOREIGN KEY (customer_id)    REFERENCES users (id),
    CONSTRAINT fk_tickets_cs       FOREIGN KEY (assigned_cs_id) REFERENCES users (id)
);

CREATE INDEX idx_tickets_customer_id ON tickets (customer_id);
CREATE INDEX idx_tickets_status      ON tickets (status);
CREATE INDEX idx_tickets_assigned_cs ON tickets (assigned_cs_id);
CREATE INDEX idx_tickets_created_at  ON tickets (created_at DESC);


-- ============================================================
-- TICKET MESSAGES
-- ============================================================
CREATE TABLE ticket_messages (
    id               UUID        PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    ticket_id        UUID        NOT NULL,
    sender_id        UUID        NOT NULL,
    message          TEXT        NOT NULL,
    is_from_cs       BOOLEAN     NOT NULL DEFAULT FALSE,
    is_internal_note BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_tmsg_ticket FOREIGN KEY (ticket_id) REFERENCES tickets (id) ON DELETE CASCADE,
    CONSTRAINT fk_tmsg_sender FOREIGN KEY (sender_id) REFERENCES users (id)
);

CREATE INDEX idx_ticket_messages_ticket_id  ON ticket_messages (ticket_id);
CREATE INDEX idx_ticket_messages_created_at ON ticket_messages (ticket_id, created_at ASC);


-- ============================================================
-- TICKET ATTACHMENTS
-- ============================================================
CREATE TABLE ticket_attachments (
    id          UUID         PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    ticket_id   UUID         NOT NULL,
    message_id  UUID,
    file_url    TEXT         NOT NULL,
    file_name   VARCHAR(255) NOT NULL,
    file_type   VARCHAR(50),
    file_size   INTEGER,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT fk_tattach_ticket   FOREIGN KEY (ticket_id)  REFERENCES tickets (id)         ON DELETE CASCADE,
    CONSTRAINT fk_tattach_message  FOREIGN KEY (message_id) REFERENCES ticket_messages (id) ON DELETE SET NULL
);

CREATE INDEX idx_ticket_attachments_ticket_id  ON ticket_attachments (ticket_id);
CREATE INDEX idx_ticket_attachments_message_id ON ticket_attachments (message_id);


-- ============================================================
-- TICKET STATUS LOGS (audit trail)
-- ============================================================
CREATE TABLE ticket_status_logs (
    id          UUID        PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    ticket_id   UUID        NOT NULL,
    changed_by  UUID        NOT NULL,
    old_status  VARCHAR(20),
    new_status  VARCHAR(20) NOT NULL,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_tslog_ticket FOREIGN KEY (ticket_id)  REFERENCES tickets (id) ON DELETE CASCADE,
    CONSTRAINT fk_tslog_user   FOREIGN KEY (changed_by) REFERENCES users (id)
);

CREATE INDEX idx_ticket_status_logs_ticket_id ON ticket_status_logs (ticket_id);


-- ============================================================
-- CHAT CONVERSATIONS
-- Generic 2-pihak: initiator & participant bisa role apapun.
-- Rules boleh/tidak boleh chat ditangani di application layer.
--
-- Rules:
--   customer → umkm ✅, cs ✅
--   umkm     → customer ✅, admin ✅, cs ✅, finance ✅
--   admin    → umkm ✅, admin ✅, cs ✅, finance ✅
--   cs       → semua ✅
--   finance  → umkm ✅, admin ✅, cs ✅, finance ✅
--
-- (A,B) dan (B,A) dianggap conversation yang sama,
-- handle di app layer saat create (cek keduanya sebelum insert).
-- ============================================================
CREATE TABLE chat_conversations (
    id              UUID        PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    initiator_id    UUID        NOT NULL,
    participant_id  UUID        NOT NULL,
    status          VARCHAR(10) NOT NULL DEFAULT 'open',
    last_message_at TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT ck_chat_conv_not_self    CHECK (initiator_id <> participant_id),
    CONSTRAINT ck_chat_conv_status      CHECK (status IN ('open', 'closed')),
    CONSTRAINT uq_chat_conv_pair        UNIQUE (initiator_id, participant_id),
    CONSTRAINT fk_chat_conv_initiator   FOREIGN KEY (initiator_id)  REFERENCES users (id),
    CONSTRAINT fk_chat_conv_participant FOREIGN KEY (participant_id) REFERENCES users (id)
);

CREATE INDEX idx_chat_conv_initiator_id   ON chat_conversations (initiator_id);
CREATE INDEX idx_chat_conv_participant_id ON chat_conversations (participant_id);
CREATE INDEX idx_chat_conv_status         ON chat_conversations (status);
CREATE INDEX idx_chat_conv_last_msg       ON chat_conversations (last_message_at DESC);


-- ============================================================
-- CHAT MESSAGES
-- ============================================================
CREATE TABLE chat_messages (
    id              UUID        PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    conversation_id UUID        NOT NULL,
    sender_id       UUID        NOT NULL,
    message         TEXT        NOT NULL,
    attachment_url  TEXT,
    is_read         BOOLEAN     NOT NULL DEFAULT FALSE,
    read_at         TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT fk_cmsg_conv   FOREIGN KEY (conversation_id) REFERENCES chat_conversations (id) ON DELETE CASCADE,
    CONSTRAINT fk_cmsg_sender FOREIGN KEY (sender_id)       REFERENCES users (id)
);

CREATE INDEX idx_chat_messages_conv_id ON chat_messages (conversation_id);
CREATE INDEX idx_chat_messages_created ON chat_messages (conversation_id, created_at ASC);
CREATE INDEX idx_chat_messages_unread  ON chat_messages (conversation_id, is_read) WHERE NOT is_read;


-- ============================================================
-- REPLY TEMPLATES (hanya untuk CS)
-- ============================================================
CREATE TABLE reply_templates (
    id          UUID         PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    title       VARCHAR(150) NOT NULL,
    content     TEXT         NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    created_by  UUID         NOT NULL,
    updated_by  UUID,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),

    CONSTRAINT fk_rtpl_creator FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT fk_rtpl_updater FOREIGN KEY (updated_by) REFERENCES users (id)
);

CREATE INDEX idx_reply_templates_active ON reply_templates (is_active) WHERE is_active = TRUE;
