-- ============================================================
-- Seed data for COMPLETE Checkout flow testing
-- Self-contained: cleans existing data and re-seeds everything.
-- Assumes roles & couriers were created by migrations/app startup.
--
-- Run:
--   PGPASSWORD=postgres psql -h localhost -U postgres haji_umroh_store \
--     -f scripts/seed_checkout_data.sql
-- ============================================================

BEGIN;

-- ============================================================
-- 0. CLEANUP — truncate all seed-managed tables (CASCADE)
-- ============================================================
DELETE FROM vendor_couriers;
DELETE FROM couriers;
TRUNCATE TABLE
  cart_items, carts,
  product_images, product_variants, products,
  vendor_bank_accounts, vendor_balances, vendors,
  addresses, categories,
  users
CASCADE;

-- Re-seed couriers (auto-increment IDs reset)
INSERT INTO couriers (id, code, name) VALUES
  (1,  'jne',     'JNE'),
  (2,  'sicepat', 'SiCepat'),
  (3,  'ide',     'IDExpress'),
  (4,  'sap',     'SAP Express'),
  (5,  'ninja',   'Ninja'),
  (6,  'jnt',     'J&T Express'),
  (7,  'tiki',    'TIKI'),
  (8,  'wahana',  'Wahana Express'),
  (9,  'pos',     'POS Indonesia'),
  (10, 'sentral', 'Sentral Cargo'),
  (11, 'lion',    'Lion Parcel'),
  (12, 'rex',     'Royal Express Asia');
SELECT setval('couriers_id_seq', 12);

-- ============================================================
-- REFERENCE IDs (fixed for reproducible dev/test)
-- ============================================================
-- Couriers (seeded by migration, NOT truncated):
--   1=jne 2=sicepat 3=ide 4=sap 5=ninja 6=jnt
--   7=tiki 8=wahana 9=pos 10=sentral 11=lion 12=rex
--
-- Password hash for 'Test1234!' (shared across all seed users):
--   $2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y

-- ============================================================
-- 1. ALL USERS (admin + umkm×3 + customer×2 + cs×3 + finance)
-- ============================================================
INSERT INTO users (id, email, full_name, phone, password_hash, role_id, status, email_verified_at) VALUES
  -- Admin
  ('7907d953-ba02-40a4-b403-07b68d8471e8',
   'admin@example.com', 'Admin', '08120000000',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   1, 'active', NOW()),
  -- UMKM 1 (Vendor 1 owner)
  ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d',
   'umkm@dev.local', 'Dev UMKM', '081200000001',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   2, 'active', NOW()),
  -- UMKM 2 (Vendor 2 owner)
  ('aa000002-0000-0000-0000-000000000001',
   'umkm2@dev.local', 'Toko Oleh-Oleh Tanah Suci', '081299000002',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   2, 'active', NOW()),
  -- UMKM 3 (Vendor 3 owner)
  ('aa000002-0000-0000-0000-000000000002',
   'umkm3@dev.local', 'Perlengkapan Ibadah Barokah', '081299000003',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   2, 'active', NOW()),
  -- Customer 1
  ('7d281d37-8f4a-48a8-9b33-e47c76b3d13c',
   'customer@dev.local', 'Dev Customer', '081200000002',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   3, 'active', NOW()),
  -- Customer 2
  ('aa000002-0000-0000-0000-000000000003',
   'customer2@dev.local', 'Siti Aisyah', '081299000004',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   3, 'active', NOW()),
  -- CS 1
  ('cc000001-0000-0000-0000-000000000001',
   'cs@dev.local', 'Dev Customer Service 1', '081200000003',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   4, 'active', NOW()),
  -- CS 2
  ('cc000001-0000-0000-0000-000000000002',
   'cs2@dev.local', 'Dev Customer Service 2', '081200000005',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   4, 'active', NOW()),
  -- CS 3
  ('cc000001-0000-0000-0000-000000000003',
   'cs3@dev.local', 'Dev Customer Service 3', '081200000006',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   4, 'active', NOW()),
  -- Finance
  ('ff000001-0000-0000-0000-000000000001',
   'finance@dev.local', 'Dev Finance', '081200000004',
   '$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y',
   5, 'active', NOW());

-- ============================================================
-- 2. VENDORS (3 vendors — different types & cities)
-- ============================================================
INSERT INTO vendors (id, owner_user_id, vendor_type, legal_name, display_name,
  responsible_person_name, description, status, approved_by, approved_at) VALUES
  -- Vendor 1: existing umkm@dev.local → Jakarta
  ('da000001-0000-0000-0000-000000000001',
   '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d',
   'hajj_souvenir_store',
   'PT Ihram Jaya Nusantara',
   'Ihram Jaya',
   'Dev UMKM',
   'Toko perlengkapan haji & umrah terlengkap di Jakarta. Menyediakan kain ihram, sajadah, dan perlengkapan ibadah lainnya dengan kualitas premium.',
   'active',
   '7907d953-ba02-40a4-b403-07b68d8471e8', NOW()),
  -- Vendor 2: umkm2@dev.local → Surabaya
  ('da000001-0000-0000-0000-000000000002',
   'aa000002-0000-0000-0000-000000000001',
   'umrah_souvenir_store',
   'CV Oleh-Oleh Tanah Suci',
   'Oleh-Oleh Tanah Suci',
   'Toko Oleh-Oleh Tanah Suci',
   'Pusat oleh-oleh haji dan umrah langsung dari Tanah Suci. Kurma premium, minyak zaitun asli, dan souvenir eksklusif Mekkah-Madinah.',
   'active',
   '7907d953-ba02-40a4-b403-07b68d8471e8', NOW()),
  -- Vendor 3: umkm3@dev.local → Bandung
  ('da000001-0000-0000-0000-000000000003',
   'aa000002-0000-0000-0000-000000000002',
   'general_souvenir_store',
   'UD Barokah Store',
   'Barokah Store',
   'Perlengkapan Ibadah Barokah',
   'Menyediakan berbagai perlengkapan ibadah berkualitas dengan harga terjangkau. Mulai dari mukena, tasbih, Al-Quran, hingga hampers haji umrah.',
   'active',
   '7907d953-ba02-40a4-b403-07b68d8471e8', NOW());

-- ============================================================
-- 3. VENDOR BALANCES
-- ============================================================
INSERT INTO vendor_balances (vendor_id) VALUES
  ('da000001-0000-0000-0000-000000000001'),
  ('da000001-0000-0000-0000-000000000002'),
  ('da000001-0000-0000-0000-000000000003');

-- ============================================================
-- 4. VENDOR BANK ACCOUNTS (1 per vendor, pending verification)
-- ============================================================
INSERT INTO vendor_bank_accounts (id, vendor_id, bank_name, account_number,
  account_holder_name, verification_status) VALUES
  ('db000001-0000-0000-0000-000000000001',
   'da000001-0000-0000-0000-000000000001',
   'Bank Syariah Indonesia', '7788001122', 'PT Ihram Jaya Nusantara', 'pending'),
  ('db000001-0000-0000-0000-000000000002',
   'da000001-0000-0000-0000-000000000002',
   'Bank Muamalat', '1122334455', 'CV Oleh-Oleh Tanah Suci', 'pending'),
  ('db000001-0000-0000-0000-000000000003',
   'da000001-0000-0000-0000-000000000003',
   'Bank BCA Syariah', '5566001122', 'UD Barokah Store', 'pending');

-- ============================================================
-- 5. VENDOR ADDRESSES (origin for RajaOngkir shipping cost)
--    district_id is CRITICAL for cost calculation
-- ============================================================
INSERT INTO addresses (id, user_id, label, recipient_name, phone,
  province_id, province_name, city_id, city_name, district_id, district_name,
  postal_code, address_line, is_default, notes) VALUES
  -- Vendor 1 (Jakarta Selatan, Kebayoran Baru, district 2096)
  ('ab000001-0000-0000-0000-000000000001',
   '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d',
   'Gudang Utama', 'Ihram Jaya', '021-7654321',
   '6', 'DKI Jakarta', '152', 'Jakarta Selatan', '2096', 'Kebayoran Baru',
   '12160', 'Jl. Fatmawati Raya No. 88, Gedung D Lt. 1', true,
   'Gudang utama, buka Senin-Sabtu 08:00-17:00'),
  -- Vendor 2 (Surabaya, Gubeng, district 3980)
  ('ab000001-0000-0000-0000-000000000002',
   'aa000002-0000-0000-0000-000000000001',
   'Toko Pusat', 'Oleh-Oleh Tanah Suci', '031-5551234',
   '11', 'Jawa Timur', '444', 'Surabaya', '3980', 'Gubeng',
   '60281', 'Jl. Raya Gubeng No. 23, Ruko Blok A5', true,
   'Depan Stasiun Gubeng'),
  -- Vendor 3 (Bandung, Coblong, district 832)
  ('ab000001-0000-0000-0000-000000000003',
   'aa000002-0000-0000-0000-000000000002',
   'Toko & Gudang', 'Barokah Store', '022-4201234',
   '9', 'Jawa Barat', '23', 'Bandung', '832', 'Coblong',
   '40132', 'Jl. Dipatiukur No. 45, Dekat ITB', true,
   'Samping minimarket Alfamart');

-- ============================================================
-- 6. VENDOR COURIERS (each vendor different combos)
-- ============================================================
INSERT INTO vendor_couriers (id, vendor_id, courier_id, is_active) VALUES
  -- Vendor 1 (Jakarta): JNE(1), SiCepat(2), J&T(6), Lion(11)
  ('ac000001-0000-0000-0000-000000000001', 'da000001-0000-0000-0000-000000000001', 1, true),
  ('ac000001-0000-0000-0000-000000000002', 'da000001-0000-0000-0000-000000000001', 2, true),
  ('ac000001-0000-0000-0000-000000000003', 'da000001-0000-0000-0000-000000000001', 6, true),
  ('ac000001-0000-0000-0000-000000000004', 'da000001-0000-0000-0000-000000000001', 11, true),
  -- Vendor 2 (Surabaya): JNE(1), TIKI(7), POS(9), Wahana(8)
  ('ac000001-0000-0000-0000-000000000005', 'da000001-0000-0000-0000-000000000002', 1, true),
  ('ac000001-0000-0000-0000-000000000006', 'da000001-0000-0000-0000-000000000002', 7, true),
  ('ac000001-0000-0000-0000-000000000007', 'da000001-0000-0000-0000-000000000002', 9, true),
  ('ac000001-0000-0000-0000-000000000008', 'da000001-0000-0000-0000-000000000002', 8, true),
  -- Vendor 3 (Bandung): SiCepat(2), JNE(1), IDExpress(3), Ninja(5)
  ('ac000001-0000-0000-0000-000000000009', 'da000001-0000-0000-0000-000000000003', 2, true),
  ('ac000001-0000-0000-0000-000000000010', 'da000001-0000-0000-0000-000000000003', 1, true),
  ('ac000001-0000-0000-0000-000000000011', 'da000001-0000-0000-0000-000000000003', 3, true),
  ('ac000001-0000-0000-0000-000000000012', 'da000001-0000-0000-0000-000000000003', 5, true);

-- ============================================================
-- 7. CATEGORIES (hierarchical: 3 parents + 8 children)
-- ============================================================
INSERT INTO categories (id, parent_id, name, slug, is_active) VALUES
  -- Parent categories
  ('ca000001-0000-0000-0000-000000000001', NULL, 'Perlengkapan Haji & Umrah', 'perlengkapan-haji-umrah', true),
  ('ca000001-0000-0000-0000-000000000002', NULL, 'Oleh-oleh & Souvenir', 'oleh-oleh-souvenir', true),
  ('ca000001-0000-0000-0000-000000000003', NULL, 'Fashion Muslim', 'fashion-muslim', true),
  -- Children of "Perlengkapan Haji & Umrah"
  ('ca000001-0000-0000-0000-000000000011', 'ca000001-0000-0000-0000-000000000001', 'Pakaian Ihram', 'pakaian-ihram', true),
  ('ca000001-0000-0000-0000-000000000012', 'ca000001-0000-0000-0000-000000000001', 'Sajadah & Alat Shalat', 'sajadah-alat-shalat', true),
  ('ca000001-0000-0000-0000-000000000013', 'ca000001-0000-0000-0000-000000000001', 'Tasbih & Hampers', 'tasbih-hampers', true),
  -- Children of "Oleh-oleh & Souvenir"
  ('ca000001-0000-0000-0000-000000000021', 'ca000001-0000-0000-0000-000000000002', 'Kurma', 'kurma', true),
  ('ca000001-0000-0000-0000-000000000022', 'ca000001-0000-0000-0000-000000000002', 'Minyak & Herbal', 'minyak-herbal', true),
  ('ca000001-0000-0000-0000-000000000023', 'ca000001-0000-0000-0000-000000000002', 'Air Zamzam', 'air-zamzam', true),
  -- Children of "Fashion Muslim"
  ('ca000001-0000-0000-0000-000000000031', 'ca000001-0000-0000-0000-000000000003', 'Mukena', 'mukena', true),
  ('ca000001-0000-0000-0000-000000000032', 'ca000001-0000-0000-0000-000000000003', 'Gamis & Jubah', 'gamis-jubah', true);

-- ============================================================
-- 8. PRODUCTS (8 products across 3 vendors, various categories)
-- ============================================================

-- ───── Vendor 1: Ihram Jaya (Jakarta) ─────
INSERT INTO products (id, vendor_id, category_id, name, slug, description,
  status, halal_ai_status) VALUES
  ('b0000001-0000-0000-0000-000000000001',
   'da000001-0000-0000-0000-000000000001',
   'ca000001-0000-0000-0000-000000000011',
   'Kain Ihram Pria Premium',
   'kain-ihram-pria-premium',
   'Kain ihram pria berbahan katun premium grade A, lembut dan nyaman untuk ibadah haji dan umrah. Jahitan rapi, tidak mudah kusut. Tersedia dalam ukuran Standard dan Extra Large.',
   'draft', 'passed'),
  ('b0000001-0000-0000-0000-000000000002',
   'da000001-0000-0000-0000-000000000001',
   'ca000001-0000-0000-0000-000000000012',
   'Sajadah Travel Lipat Premium',
   'sajadah-travel-lipat-premium',
   'Sajadah travel lipat praktis dengan tas penyimpanan waterproof. Bahan beludru soft-touch, ringan hanya 300 gram. Ideal untuk dibawa haji dan umrah.',
   'draft', 'passed'),
  ('b0000001-0000-0000-0000-000000000003',
   'da000001-0000-0000-0000-000000000001',
   'ca000001-0000-0000-0000-000000000013',
   'Tasbih Digital Premium 33 Biji',
   'tasbih-digital-premium-33',
   'Tasbih digital elektronik dengan counter otomatis. Baterai tahan hingga 6 bulan pemakaian normal. Material ABS anti-slip, nyaman digenggam.',
   'draft', 'passed');

-- ───── Vendor 2: Oleh-Oleh Tanah Suci (Surabaya) ─────
INSERT INTO products (id, vendor_id, category_id, name, slug, description,
  status, halal_ai_status) VALUES
  ('b0000001-0000-0000-0000-000000000004',
   'da000001-0000-0000-0000-000000000002',
   'ca000001-0000-0000-0000-000000000021',
   'Kurma Ajwa Madinah Asli',
   'kurma-ajwa-madinah-asli',
   'Kurma Ajwa asli dari kebun Al-Madinah Al-Munawwarah. Tekstur lembut, rasa manis alami. Sering disebut sebagai kurma Nabi. Kemasan vacuum seal untuk menjaga kesegaran.',
   'draft', 'passed'),
  ('b0000001-0000-0000-0000-000000000005',
   'da000001-0000-0000-0000-000000000002',
   'ca000001-0000-0000-0000-000000000022',
   'Minyak Zaitun Extra Virgin Asli',
   'minyak-zaitun-extra-virgin-asli',
   'Minyak zaitun extra virgin 100% asli impor dari Timur Tengah. Cold-pressed, tanpa campuran. Cocok untuk kesehatan dan kecantikan.',
   'draft', 'passed'),
  ('b0000001-0000-0000-0000-000000000006',
   'da000001-0000-0000-0000-000000000002',
   'ca000001-0000-0000-0000-000000000023',
   'Air Zamzam Asli 5 Liter',
   'air-zamzam-asli-5-liter',
   'Air Zamzam asli dari sumur Zamzam, Mekkah. Dikemas dalam jeriken food-grade kedap udara. Sertifikat keaslian tersedia.',
   'draft', 'passed');

-- ───── Vendor 3: Barokah Store (Bandung) ─────
INSERT INTO products (id, vendor_id, category_id, name, slug, description,
  status, halal_ai_status) VALUES
  ('b0000001-0000-0000-0000-000000000007',
   'da000001-0000-0000-0000-000000000003',
   'ca000001-0000-0000-0000-000000000031',
   'Mukena Katun Jepang Premium',
   'mukena-katun-jepang-premium',
   'Mukena bahan katun Jepang super halus dan adem. Tidak menerawang, jatuh sempurna. Dilengkapi tas travel matching. Tersedia warna Putih dan Dusty Pink.',
   'draft', 'passed'),
  ('b0000001-0000-0000-0000-000000000008',
   'da000001-0000-0000-0000-000000000003',
   'ca000001-0000-0000-0000-000000000032',
   'Gamis Pria Al-Haramain',
   'gamis-pria-al-haramain',
   'Gamis pria model Al-Haramain lengan panjang. Bahan toyobo premium anti kusut, adem, dan tidak menerawang. Cocok untuk shalat dan acara formal.',
   'draft', 'passed');

-- ============================================================
-- 9. PRODUCT VARIANTS (17 variants total, realistic SKUs & weights)
-- ============================================================

-- Kain Ihram Pria (2 sizes)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000001', 'b0000001-0000-0000-0000-000000000001',
   'IHR-PRM-STD', 'Standard (115x220 cm)', 75000, 'IDR', 150, 400, true, true),
  ('c0000001-0000-0000-0000-000000000002', 'b0000001-0000-0000-0000-000000000001',
   'IHR-PRM-XL',  'Extra Large (130x240 cm)', 95000, 'IDR', 80, 500, false, true);

-- Sajadah Travel (1 variant)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000003', 'b0000001-0000-0000-0000-000000000002',
   'SJD-TRV-GRN', 'Hijau Tosca', 45000, 'IDR', 200, 300, true, true);

-- Tasbih Digital (2 colors)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000004', 'b0000001-0000-0000-0000-000000000003',
   'TSB-DGT-BLK', 'Hitam', 25000, 'IDR', 300, 50, true, true),
  ('c0000001-0000-0000-0000-000000000005', 'b0000001-0000-0000-0000-000000000003',
   'TSB-DGT-GLD', 'Gold', 35000, 'IDR', 200, 55, false, true);

-- Kurma Ajwa (3 sizes)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000006', 'b0000001-0000-0000-0000-000000000004',
   'KRM-AJW-250', '250 gram',    65000, 'IDR', 120, 300, false, true),
  ('c0000001-0000-0000-0000-000000000007', 'b0000001-0000-0000-0000-000000000004',
   'KRM-AJW-500', '500 gram',    120000, 'IDR', 80, 550, true, true),
  ('c0000001-0000-0000-0000-000000000008', 'b0000001-0000-0000-0000-000000000004',
   'KRM-AJW-1KG', '1 Kilogram',  220000, 'IDR', 40, 1050, false, true);

-- Minyak Zaitun (2 sizes)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000009', 'b0000001-0000-0000-0000-000000000005',
   'MZT-EV-250', '250 ml',  85000, 'IDR', 100, 350, true, true),
  ('c0000001-0000-0000-0000-000000000010', 'b0000001-0000-0000-0000-000000000005',
   'MZT-EV-500', '500 ml',  155000, 'IDR', 60, 650, false, true);

-- Air Zamzam (1 variant, heavy item)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000011', 'b0000001-0000-0000-0000-000000000006',
   'ZMZ-5L', '5 Liter', 175000, 'IDR', 30, 5200, true, true);

-- Mukena Katun Jepang (2 colors)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000012', 'b0000001-0000-0000-0000-000000000007',
   'MKN-KTJ-WHT', 'Putih Polos',  185000, 'IDR', 60, 450, true, true),
  ('c0000001-0000-0000-0000-000000000013', 'b0000001-0000-0000-0000-000000000007',
   'MKN-KTJ-PNK', 'Dusty Pink',   195000, 'IDR', 45, 460, false, true);

-- Gamis Pria (3 sizes)
INSERT INTO product_variants (id, product_id, sku, variant_name, price,
  currency, stock_on_hand, weight_gram, is_default, is_active) VALUES
  ('c0000001-0000-0000-0000-000000000014', 'b0000001-0000-0000-0000-000000000008',
   'GMS-ALH-M', 'M (Lingkar Dada 104cm)', 250000, 'IDR', 40, 380, false, true),
  ('c0000001-0000-0000-0000-000000000015', 'b0000001-0000-0000-0000-000000000008',
   'GMS-ALH-L', 'L (Lingkar Dada 110cm)', 250000, 'IDR', 55, 400, true, true),
  ('c0000001-0000-0000-0000-000000000016', 'b0000001-0000-0000-0000-000000000008',
   'GMS-ALH-XL', 'XL (Lingkar Dada 118cm)', 265000, 'IDR', 35, 420, false, true);

-- ============================================================
-- 10. PRODUCT IMAGES (1-2 images per product, 13 total)
-- ============================================================
INSERT INTO product_images (id, product_id, image_url, mime_type, is_primary, sort_order) VALUES
  -- Vendor 1 products
  ('a1b00001-0000-0000-0000-000000000001', 'b0000001-0000-0000-0000-000000000001', 'products/kain-ihram-pria-1.jpg', 'image/jpeg', true,  0),
  ('a1b00001-0000-0000-0000-000000000002', 'b0000001-0000-0000-0000-000000000001', 'products/kain-ihram-pria-2.jpg', 'image/jpeg', false, 1),
  ('a1b00001-0000-0000-0000-000000000003', 'b0000001-0000-0000-0000-000000000002', 'products/sajadah-travel-1.jpg',  'image/jpeg', true,  0),
  ('a1b00001-0000-0000-0000-000000000004', 'b0000001-0000-0000-0000-000000000003', 'products/tasbih-digital-1.jpg',  'image/jpeg', true,  0),
  ('a1b00001-0000-0000-0000-000000000005', 'b0000001-0000-0000-0000-000000000003', 'products/tasbih-digital-2.jpg',  'image/jpeg', false, 1),
  -- Vendor 2 products
  ('a1b00001-0000-0000-0000-000000000006', 'b0000001-0000-0000-0000-000000000004', 'products/kurma-ajwa-1.jpg',      'image/jpeg', true,  0),
  ('a1b00001-0000-0000-0000-000000000007', 'b0000001-0000-0000-0000-000000000004', 'products/kurma-ajwa-2.jpg',      'image/jpeg', false, 1),
  ('a1b00001-0000-0000-0000-000000000008', 'b0000001-0000-0000-0000-000000000005', 'products/minyak-zaitun-1.jpg',   'image/jpeg', true,  0),
  ('a1b00001-0000-0000-0000-000000000009', 'b0000001-0000-0000-0000-000000000006', 'products/air-zamzam-1.jpg',      'image/jpeg', true,  0),
  -- Vendor 3 products
  ('a1b00001-0000-0000-0000-000000000010', 'b0000001-0000-0000-0000-000000000007', 'products/mukena-katun-1.jpg',    'image/jpeg', true,  0),
  ('a1b00001-0000-0000-0000-000000000011', 'b0000001-0000-0000-0000-000000000007', 'products/mukena-katun-2.jpg',    'image/jpeg', false, 1),
  ('a1b00001-0000-0000-0000-000000000012', 'b0000001-0000-0000-0000-000000000008', 'products/gamis-pria-1.jpg',      'image/jpeg', true,  0),
  ('a1b00001-0000-0000-0000-000000000013', 'b0000001-0000-0000-0000-000000000008', 'products/gamis-pria-2.jpg',      'image/jpeg', false, 1);

-- ============================================================
-- 11. PUBLISH all 8 products
-- ============================================================
UPDATE products SET status = 'published' WHERE id IN (
  'b0000001-0000-0000-0000-000000000001',
  'b0000001-0000-0000-0000-000000000002',
  'b0000001-0000-0000-0000-000000000003',
  'b0000001-0000-0000-0000-000000000004',
  'b0000001-0000-0000-0000-000000000005',
  'b0000001-0000-0000-0000-000000000006',
  'b0000001-0000-0000-0000-000000000007',
  'b0000001-0000-0000-0000-000000000008'
);

-- ============================================================
-- 12. CUSTOMER ADDRESSES (2 for customer1 + 1 for customer2)
-- ============================================================
INSERT INTO addresses (id, user_id, label, recipient_name, phone,
  province_id, province_name, city_id, city_name, district_id, district_name,
  postal_code, address_line, is_default, notes) VALUES
  -- Customer 1 → Rumah (Jakarta Selatan, Kebayoran Baru, district 2096)
  ('ad000001-0000-0000-0000-000000000001',
   '7d281d37-8f4a-48a8-9b33-e47c76b3d13c',
   'Rumah', 'Customer Test', '081234567890',
   '6', 'DKI Jakarta', '152', 'Jakarta Selatan', '2096', 'Kebayoran Baru',
   '12160', 'Jl. Senopati No. 123, RT 001/RW 002, Kelurahan Senayan',
   true, 'Rumah warna putih, pagar hitam'),
  -- Customer 1 → Kantor (Jakarta Pusat, Menteng, district 2118)
  ('ad000001-0000-0000-0000-000000000002',
   '7d281d37-8f4a-48a8-9b33-e47c76b3d13c',
   'Kantor', 'Customer Test', '081234567891',
   '6', 'DKI Jakarta', '153', 'Jakarta Pusat', '2118', 'Menteng',
   '10310', 'Jl. Thamrin No. 55, Gedung Perkantoran Lt. 5',
   false, 'Lobby lantai 1, bilang ke resepsionis'),
  -- Customer 2 → Rumah (Bekasi, Bekasi Barat, district 752)
  ('ad000001-0000-0000-0000-000000000003',
   'aa000002-0000-0000-0000-000000000003',
   'Rumah', 'Siti Aisyah', '081200000004',
   '9', 'Jawa Barat', '54', 'Bekasi', '752', 'Bekasi Barat',
   '17134', 'Perumahan Grand Galaxy City Blok AA-12',
   true, 'Cluster depan gerbang utama')
;

-- ============================================================
-- 13. VERIFY — count all seeded tables
-- ============================================================
DO $$ BEGIN RAISE NOTICE '──── SEED VERIFICATION ────'; END $$;

SELECT tbl, cnt FROM (
  SELECT 'users'                  AS tbl, COUNT(*) AS cnt FROM users
  UNION ALL SELECT 'vendors',                COUNT(*) FROM vendors
  UNION ALL SELECT 'vendor_balances',        COUNT(*) FROM vendor_balances
  UNION ALL SELECT 'vendor_bank_accounts',   COUNT(*) FROM vendor_bank_accounts
  UNION ALL SELECT 'vendor_couriers',        COUNT(*) FROM vendor_couriers
  UNION ALL SELECT 'addresses (total)',      COUNT(*) FROM addresses
  UNION ALL SELECT 'categories',             COUNT(*) FROM categories
  UNION ALL SELECT 'products (published)',   COUNT(*) FROM products WHERE status = 'published'
  UNION ALL SELECT 'product_variants',       COUNT(*) FROM product_variants
  UNION ALL SELECT 'product_images',         COUNT(*) FROM product_images
) AS verification
ORDER BY tbl;

COMMIT;
