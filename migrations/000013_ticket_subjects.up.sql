-- ============================================================
-- MIGRATION: Ticket Subjects (enum from DB)
-- ============================================================

-- Lookup table for predefined ticket subjects
CREATE TABLE ticket_subjects (
    id         SERIAL       PRIMARY KEY,
    label      VARCHAR(100) NOT NULL UNIQUE,
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Seed default subjects
INSERT INTO ticket_subjects (label) VALUES
    ('Barang Tidak Sampai'),
    ('Barang Rusak'),
    ('Pengembalian Dana'),
    ('Barang Tidak Sesuai'),
    ('Pengiriman Terlambat'),
    ('Pembatalan Pesanan'),
    ('Kesalahan Produk'),
    ('Akun Bermasalah'),
    ('Pembayaran Gagal'),
    ('Voucher/Promo Tidak Berlaku'),
    ('Pertanyaan Umum'),
    ('Lainnya');

-- Add subject_id column to tickets table (nullable for backward compat with existing rows)
ALTER TABLE tickets ADD COLUMN subject_id INT REFERENCES ticket_subjects(id);

-- Backfill existing tickets: try to match subject text → subject_id
UPDATE tickets t
SET subject_id = ts.id
FROM ticket_subjects ts
WHERE lower(t.subject) = lower(ts.label);

-- For any remaining unmatched tickets, assign "Lainnya"
UPDATE tickets
SET subject_id = (SELECT id FROM ticket_subjects WHERE label = 'Lainnya')
WHERE subject_id IS NULL;

-- Now make subject_id NOT NULL for future inserts
ALTER TABLE tickets ALTER COLUMN subject_id SET NOT NULL;

-- Index for report queries (top kendala)
CREATE INDEX idx_tickets_subject_id ON tickets (subject_id);
