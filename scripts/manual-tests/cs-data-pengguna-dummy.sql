-- ============================================================
-- DUMMY DATA: CS Data Pengguna — Manual Testing
-- ============================================================
-- Jalankan SETELAH semua migrasi diterapkan.
-- Pastikan tabel `roles` dan `ticket_subjects` sudah terisi (dari migrasi).
--
-- Data yg dibuat:
--   • 1 CS agent (cs@dev.local)
--   • 1 Vendor/UMKM owner (vendor@dev.local) + vendor entity
--   • 5 Customer (customer01–05@dev.local)
--   • 1 Kategori produk
--   • 1 Produk + 1 varian
--   • 5 Order (1 per customer) dengan payment_invoices
--   • 7 Tiket bantuan (beberapa customer punya >1 tiket)
--   • Beberapa ticket_messages
--
-- Password semua akun: Password123!
-- bcrypt hash ($2a$10$...) di bawah adalah hash dari "Password123!"
-- ============================================================

BEGIN;

-- ════════════════════════════════════════════════════════════
-- 1. USERS
-- ════════════════════════════════════════════════════════════

-- CS Agent
INSERT INTO users (id, email, full_name, phone, password_hash, role_id, status, email_verified_at)
VALUES (
  'a0000000-0000-0000-0000-000000000001',
  'cs@dev.local',
  'CS Agent Satu',
  '081200000001',
  '$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012', -- Password123!
  4, -- cs
  'active',
  NOW()
) ON CONFLICT (id) DO NOTHING;

-- Vendor/UMKM owner
INSERT INTO users (id, email, full_name, phone, password_hash, role_id, status, email_verified_at)
VALUES (
  'b0000000-0000-0000-0000-000000000001',
  'vendor@dev.local',
  'Toko Oleh-Oleh Haji',
  '081200000002',
  '$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012',
  2, -- umkm
  'active',
  NOW()
) ON CONFLICT (id) DO NOTHING;

-- Customers
INSERT INTO users (id, email, full_name, phone, password_hash, role_id, status, email_verified_at, birth_date) VALUES
  ('c0000000-0000-0000-0000-000000000001', 'customer01@dev.local', 'Budi Santoso',       '081300000001', '$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012', 3, 'active',  NOW(), '1990-05-15'),
  ('c0000000-0000-0000-0000-000000000002', 'customer02@dev.local', 'Siti Nurhaliza',      '081300000002', '$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012', 3, 'active',  NOW(), '1988-11-22'),
  ('c0000000-0000-0000-0000-000000000003', 'customer03@dev.local', 'Ahmad Fauzi',         '081300000003', '$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012', 3, 'active',  NOW(), '1995-03-08'),
  ('c0000000-0000-0000-0000-000000000004', 'customer04@dev.local', 'Dewi Ratnasari',      '081300000004', '$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012', 3, 'pending', NULL,  '1992-07-30'),
  ('c0000000-0000-0000-0000-000000000005', 'customer05@dev.local', 'Muhammad Rizky Alif', '081300000005', '$2a$10$abcdefghijklmnopqrstuuABCDEFGHIJKLMNOPQRSTUVWXYZ012', 3, 'active',  NOW(), '1997-01-12')
ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 2. VENDOR
-- ════════════════════════════════════════════════════════════
INSERT INTO vendors (id, owner_user_id, vendor_type, display_name, responsible_person_name, status, approved_at)
VALUES (
  'd0000000-0000-0000-0000-000000000001',
  'b0000000-0000-0000-0000-000000000001',
  'umrah_souvenir_store',
  'Toko Oleh-Oleh Haji Mabrur',
  'Haji Mabrur',
  'active',
  NOW()
) ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 3. CATEGORY + PRODUCT + VARIANT
-- ════════════════════════════════════════════════════════════
INSERT INTO categories (id, parent_id, name, slug, is_active)
VALUES ('e0000000-0000-0000-0000-000000000001', NULL, 'Perlengkapan Haji', 'perlengkapan-haji', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, vendor_id, category_id, name, slug, description, status)
VALUES (
  'f0000000-0000-0000-0000-000000000001',
  'd0000000-0000-0000-0000-000000000001',
  'e0000000-0000-0000-0000-000000000001',
  'Kain Ihram Premium',
  'kain-ihram-premium-dummy',
  'Kain ihram premium berbahan lembut dan nyaman dipakai.',
  'draft'
) ON CONFLICT (id) DO NOTHING;

INSERT INTO product_variants (id, product_id, sku, variant_name, price, stock_on_hand, weight_gram, is_default, is_active)
VALUES (
  'f1000000-0000-0000-0000-000000000001',
  'f0000000-0000-0000-0000-000000000001',
  'IHRAM-PRM-001',
  'Kain Ihram Premium - Putih',
  150000.00,
  100,
  500,
  TRUE,
  TRUE
) ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 4. ORDERS (1 per customer)
-- ════════════════════════════════════════════════════════════
INSERT INTO orders (id, order_no, user_id, vendor_id, shipping_address_snapshot, order_status, payment_status, subtotal, shipping_fee, platform_fee, grand_total) VALUES
  ('10000000-0000-0000-0000-000000000001', 'ORD-20260301-0001', 'c0000000-0000-0000-0000-000000000001', 'd0000000-0000-0000-0000-000000000001',
   '{"recipient_name":"Budi Santoso","phone":"081300000001","address_line":"Jl. Merdeka No.1, Jakarta","city":"Jakarta","province":"DKI Jakarta","postal_code":"10110"}',
   'completed', 'paid', 150000.00, 15000.00, 5000.00, 170000.00),
  ('10000000-0000-0000-0000-000000000002', 'ORD-20260302-0002', 'c0000000-0000-0000-0000-000000000002', 'd0000000-0000-0000-0000-000000000001',
   '{"recipient_name":"Siti Nurhaliza","phone":"081300000002","address_line":"Jl. Sudirman No.5, Bandung","city":"Bandung","province":"Jawa Barat","postal_code":"40111"}',
   'shipped', 'paid', 300000.00, 20000.00, 10000.00, 330000.00),
  ('10000000-0000-0000-0000-000000000003', 'ORD-20260303-0003', 'c0000000-0000-0000-0000-000000000003', 'd0000000-0000-0000-0000-000000000001',
   '{"recipient_name":"Ahmad Fauzi","phone":"081300000003","address_line":"Jl. Diponegoro No.10, Surabaya","city":"Surabaya","province":"Jawa Timur","postal_code":"60111"}',
   'paid', 'paid', 450000.00, 25000.00, 15000.00, 490000.00),
  ('10000000-0000-0000-0000-000000000004', 'ORD-20260304-0004', 'c0000000-0000-0000-0000-000000000004', 'd0000000-0000-0000-0000-000000000001',
   '{"recipient_name":"Dewi Ratnasari","phone":"081300000004","address_line":"Jl. Gatot Subroto No.20, Semarang","city":"Semarang","province":"Jawa Tengah","postal_code":"50111"}',
   'pending_payment', 'unpaid', 150000.00, 10000.00, 5000.00, 165000.00),
  ('10000000-0000-0000-0000-000000000005', 'ORD-20260305-0005', 'c0000000-0000-0000-0000-000000000005', 'd0000000-0000-0000-0000-000000000001',
   '{"recipient_name":"Muhammad Rizky Alif","phone":"081300000005","address_line":"Jl. Ahmad Yani No.15, Medan","city":"Medan","province":"Sumatera Utara","postal_code":"20111"}',
   'canceled', 'unpaid', 600000.00, 30000.00, 20000.00, 650000.00)
ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 5. ORDER ITEMS
-- ════════════════════════════════════════════════════════════
INSERT INTO order_items (id, order_id, product_variant_id, product_name_snapshot, sku_snapshot, qty, unit_price, line_total) VALUES
  ('11000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'f1000000-0000-0000-0000-000000000001', 'Kain Ihram Premium - Putih', 'IHRAM-PRM-001', 1, 150000.00, 150000.00),
  ('11000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', 'f1000000-0000-0000-0000-000000000001', 'Kain Ihram Premium - Putih', 'IHRAM-PRM-001', 2, 150000.00, 300000.00),
  ('11000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000003', 'f1000000-0000-0000-0000-000000000001', 'Kain Ihram Premium - Putih', 'IHRAM-PRM-001', 3, 150000.00, 450000.00),
  ('11000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000004', 'f1000000-0000-0000-0000-000000000001', 'Kain Ihram Premium - Putih', 'IHRAM-PRM-001', 1, 150000.00, 150000.00),
  ('11000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', 'f1000000-0000-0000-0000-000000000001', 'Kain Ihram Premium - Putih', 'IHRAM-PRM-001', 4, 150000.00, 600000.00)
ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 6. PAYMENT INVOICES
-- ════════════════════════════════════════════════════════════
INSERT INTO payment_invoices (id, order_id, gateway, external_invoice_id, payment_method, payment_channel, amount, status, paid_at) VALUES
  ('12000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', 'xendit', 'INV-DUMMY-0001', 'BANK_TRANSFER', 'BCA',  170000.00, 'paid', NOW() - INTERVAL '5 days'),
  ('12000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000002', 'xendit', 'INV-DUMMY-0002', 'EWALLET',       'OVO',  330000.00, 'paid', NOW() - INTERVAL '4 days'),
  ('12000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000003', 'xendit', 'INV-DUMMY-0003', 'BANK_TRANSFER', 'BNI',  490000.00, 'paid', NOW() - INTERVAL '3 days'),
  ('12000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000004', 'xendit', 'INV-DUMMY-0004', NULL,             NULL,   165000.00, 'pending', NULL),
  ('12000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000005', 'xendit', 'INV-DUMMY-0005', NULL,             NULL,   650000.00, 'expired', NULL)
ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 7. TICKETS
-- ════════════════════════════════════════════════════════════
-- ticket_subjects sudah terisi dari migrasi:
--   1=Barang Tidak Sampai, 2=Barang Rusak, 3=Pengembalian Dana,
--   4=Barang Tidak Sesuai, 5=Pengiriman Terlambat, 6=Pembatalan Pesanan,
--   7=Kesalahan Produk, 8=Akun Bermasalah, 9=Pembayaran Gagal,
--   10=Voucher/Promo Tidak Berlaku, 11=Pertanyaan Umum, 12=Lainnya

INSERT INTO tickets (id, ticket_number, customer_id, assigned_cs_id, order_number, phone, reporter_name, subject, subject_id, detail, status, source, closed_at) VALUES
  -- Customer 01: 2 tiket
  ('20000000-0000-0000-0000-000000000001', 'TKT-20260301-001',
   'c0000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001',
   'ORD-20260301-0001', '081300000001', 'Budi Santoso',
   'Barang Tidak Sampai', 1,
   'Pesanan saya sudah 7 hari tapi belum sampai. Mohon dicek.',
   'on_progress', 'app', NULL),

  ('20000000-0000-0000-0000-000000000002', 'TKT-20260305-002',
   'c0000000-0000-0000-0000-000000000001', NULL,
   'ORD-20260301-0001', '081300000001', 'Budi Santoso',
   'Pertanyaan Umum', 11,
   'Apakah bisa request packaging khusus?',
   'open', 'web', NULL),

  -- Customer 02: 2 tiket
  ('20000000-0000-0000-0000-000000000003', 'TKT-20260302-003',
   'c0000000-0000-0000-0000-000000000002', 'a0000000-0000-0000-0000-000000000001',
   'ORD-20260302-0002', '081300000002', 'Siti Nurhaliza',
   'Barang Rusak', 2,
   'Kain ihram yang diterima ada robekan di bagian pinggir.',
   'closed', 'app', NOW() - INTERVAL '1 day'),

  ('20000000-0000-0000-0000-000000000004', 'TKT-20260306-004',
   'c0000000-0000-0000-0000-000000000002', NULL,
   'ORD-20260302-0002', '081300000002', 'Siti Nurhaliza',
   'Pengembalian Dana', 3,
   'Saya ingin refund untuk pesanan barang yang rusak.',
   'open', 'web', NULL),

  -- Customer 03: 1 tiket
  ('20000000-0000-0000-0000-000000000005', 'TKT-20260303-005',
   'c0000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001',
   'ORD-20260303-0003', '081300000003', 'Ahmad Fauzi',
   'Pengiriman Terlambat', 5,
   'Paket sudah 5 hari di status "shipped" tapi belum sampai.',
   'on_progress', 'app', NULL),

  -- Customer 04: 1 tiket
  ('20000000-0000-0000-0000-000000000006', 'TKT-20260304-006',
   'c0000000-0000-0000-0000-000000000004', NULL,
   'ORD-20260304-0004', '081300000004', 'Dewi Ratnasari',
   'Pembayaran Gagal', 9,
   'Pembayaran via OVO gagal terus. Sudah coba 3 kali.',
   'open', 'web', NULL),

  -- Customer 05: 1 tiket — sudah closed
  ('20000000-0000-0000-0000-000000000007', 'TKT-20260305-007',
   'c0000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001',
   'ORD-20260305-0005', '081300000005', 'Muhammad Rizky Alif',
   'Pembatalan Pesanan', 6,
   'Saya ingin membatalkan pesanan karena salah pilih ukuran.',
   'closed', 'app', NOW() - INTERVAL '2 days')
ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 8. TICKET MESSAGES
-- ════════════════════════════════════════════════════════════
INSERT INTO ticket_messages (id, ticket_id, sender_id, message, is_from_cs, is_internal_note, created_at) VALUES
  -- Tiket 001: percakapan
  ('30000000-0000-0000-0000-000000000001',
   '20000000-0000-0000-0000-000000000001',
   'c0000000-0000-0000-0000-000000000001',
   'Halo, pesanan saya dengan nomor ORD-20260301-0001 sudah 7 hari belum sampai. Bisa dicek?',
   FALSE, FALSE, NOW() - INTERVAL '3 days'),

  ('30000000-0000-0000-0000-000000000002',
   '20000000-0000-0000-0000-000000000001',
   'a0000000-0000-0000-0000-000000000001',
   'Terima kasih sudah menghubungi kami. Kami sedang menghubungi pihak kurir untuk update pengiriman Anda.',
   TRUE, FALSE, NOW() - INTERVAL '3 days' + INTERVAL '1 hour'),

  ('30000000-0000-0000-0000-000000000003',
   '20000000-0000-0000-0000-000000000001',
   'a0000000-0000-0000-0000-000000000001',
   '[Internal] Sudah follow up ke JNE, resi JNE12345. Estimasi 2 hari lagi.',
   TRUE, TRUE, NOW() - INTERVAL '2 days'),

  ('30000000-0000-0000-0000-000000000004',
   '20000000-0000-0000-0000-000000000001',
   'c0000000-0000-0000-0000-000000000001',
   'Baik, terima kasih infonya. Saya tunggu ya.',
   FALSE, FALSE, NOW() - INTERVAL '2 days' + INTERVAL '2 hours'),

  -- Tiket 003: sudah selesai
  ('30000000-0000-0000-0000-000000000005',
   '20000000-0000-0000-0000-000000000003',
   'c0000000-0000-0000-0000-000000000002',
   'Kain ihram yang saya terima ada robekan. Lihat foto yang saya lampirkan.',
   FALSE, FALSE, NOW() - INTERVAL '4 days'),

  ('30000000-0000-0000-0000-000000000006',
   '20000000-0000-0000-0000-000000000003',
   'a0000000-0000-0000-0000-000000000001',
   'Mohon maaf atas ketidaknyamanannya. Kami akan mengirimkan pengganti. Tiket ini kami tutup setelah barang pengganti dikirim.',
   TRUE, FALSE, NOW() - INTERVAL '3 days'),

  -- Tiket 005: on_progress
  ('30000000-0000-0000-0000-000000000007',
   '20000000-0000-0000-0000-000000000005',
   'c0000000-0000-0000-0000-000000000003',
   'Paket saya sudah 5 hari di status shipped. Apakah ada masalah?',
   FALSE, FALSE, NOW() - INTERVAL '2 days'),

  ('30000000-0000-0000-0000-000000000008',
   '20000000-0000-0000-0000-000000000005',
   'a0000000-0000-0000-0000-000000000001',
   'Sedang kami tracking. Mohon tunggu 1x24 jam untuk update selanjutnya.',
   TRUE, FALSE, NOW() - INTERVAL '1 day'),

  -- Tiket 007: closed
  ('30000000-0000-0000-0000-000000000009',
   '20000000-0000-0000-0000-000000000007',
   'c0000000-0000-0000-0000-000000000005',
   'Tolong batalkan pesanan saya ORD-20260305-0005. Salah pilih ukuran.',
   FALSE, FALSE, NOW() - INTERVAL '4 days'),

  ('30000000-0000-0000-0000-000000000010',
   '20000000-0000-0000-0000-000000000007',
   'a0000000-0000-0000-0000-000000000001',
   'Pesanan sudah kami batalkan. Terima kasih, tiket ini kami tutup.',
   TRUE, FALSE, NOW() - INTERVAL '2 days')
ON CONFLICT (id) DO NOTHING;

-- ════════════════════════════════════════════════════════════
-- 9. TICKET STATUS LOGS (audit trail)
-- ════════════════════════════════════════════════════════════
INSERT INTO ticket_status_logs (id, ticket_id, changed_by, old_status, new_status, notes, created_at) VALUES
  -- Tiket 001: open → on_progress
  ('40000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', 'a0000000-0000-0000-0000-000000000001', 'open', 'on_progress', 'CS mengambil tiket', NOW() - INTERVAL '3 days'),

  -- Tiket 003: open → on_progress → closed
  ('40000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'open', 'on_progress', 'CS mengambil tiket', NOW() - INTERVAL '4 days'),
  ('40000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000003', 'a0000000-0000-0000-0000-000000000001', 'on_progress', 'closed', 'Barang pengganti dikirim', NOW() - INTERVAL '1 day'),

  -- Tiket 005: open → on_progress
  ('40000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000005', 'a0000000-0000-0000-0000-000000000001', 'open', 'on_progress', 'CS mengambil tiket', NOW() - INTERVAL '2 days'),

  -- Tiket 007: open → on_progress → closed
  ('40000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 'open', 'on_progress', 'CS mengambil tiket', NOW() - INTERVAL '4 days'),
  ('40000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000007', 'a0000000-0000-0000-0000-000000000001', 'on_progress', 'closed', 'Pesanan dibatalkan sesuai permintaan', NOW() - INTERVAL '2 days')
ON CONFLICT (id) DO NOTHING;

COMMIT;

-- ============================================================
-- CLEANUP (jalankan manual jika ingin menghapus dummy data)
-- ============================================================
-- BEGIN;
-- DELETE FROM ticket_status_logs WHERE id IN ('40000000-0000-0000-0000-000000000001','40000000-0000-0000-0000-000000000002','40000000-0000-0000-0000-000000000003','40000000-0000-0000-0000-000000000004','40000000-0000-0000-0000-000000000005','40000000-0000-0000-0000-000000000006');
-- DELETE FROM ticket_messages WHERE id LIKE '30000000-0000-0000-0000-%';
-- DELETE FROM tickets WHERE id LIKE '20000000-0000-0000-0000-%';
-- DELETE FROM payment_invoices WHERE id LIKE '12000000-0000-0000-0000-%';
-- DELETE FROM order_items WHERE id LIKE '11000000-0000-0000-0000-%';
-- DELETE FROM orders WHERE id LIKE '10000000-0000-0000-0000-%';
-- DELETE FROM product_variants WHERE id = 'f1000000-0000-0000-0000-000000000001';
-- DELETE FROM products WHERE id = 'f0000000-0000-0000-0000-000000000001';
-- DELETE FROM categories WHERE id = 'e0000000-0000-0000-0000-000000000001';
-- DELETE FROM vendors WHERE id = 'd0000000-0000-0000-0000-000000000001';
-- DELETE FROM users WHERE id IN ('a0000000-0000-0000-0000-000000000001','b0000000-0000-0000-0000-000000000001','c0000000-0000-0000-0000-000000000001','c0000000-0000-0000-0000-000000000002','c0000000-0000-0000-0000-000000000003','c0000000-0000-0000-0000-000000000004','c0000000-0000-0000-0000-000000000005');
-- COMMIT;
