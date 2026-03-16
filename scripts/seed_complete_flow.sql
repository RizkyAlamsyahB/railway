-- =============================================================
-- Seed: 1 complete data row for ALL tables
-- Customer : 7d281d37-8f4a-48a8-9b33-e47c76b3d13c
-- Vendor   : 8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d
--
-- Tables already seeded by migrations (skip insert):
--   roles, ledger_accounts, ticket_subjects, couriers
-- =============================================================

BEGIN;

DO $$
DECLARE
    -- ===== existing actors =====
    v_customer_id    UUID := '7d281d37-8f4a-48a8-9b33-e47c76b3d13c';
    v_vendor_id      UUID := '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d';

    -- ===== new actors =====
    v_umkm_user_id   UUID := gen_random_uuid();
    v_cs_id          UUID;
    v_admin_id       UUID;

    -- ===== role IDs =====
    v_role_umkm      SMALLINT;
    v_role_cs        SMALLINT;
    v_role_admin     SMALLINT;

    -- ===== generated IDs =====
    v_category_id    UUID := gen_random_uuid();
    v_product_id     UUID := gen_random_uuid();
    v_variant_id     UUID := gen_random_uuid();
    v_image_id       UUID := gen_random_uuid();
    v_address_id     UUID := gen_random_uuid();
    v_cart_id        UUID := gen_random_uuid();
    v_cart_item_id   UUID := gen_random_uuid();
    v_order_id       UUID := gen_random_uuid();
    v_order_item_id  UUID := gen_random_uuid();
    v_invoice_id     UUID := gen_random_uuid();
    v_event_id       UUID := gen_random_uuid();
    v_shipment_id    UUID := gen_random_uuid();
    v_refund_id      UUID := gen_random_uuid();
    v_payout_id      UUID := gen_random_uuid();
    v_payout_item_id UUID := gen_random_uuid();
    v_journal_id     UUID := gen_random_uuid();
    v_line_d_id      UUID := gen_random_uuid();
    v_line_c_id      UUID := gen_random_uuid();
    v_ticket_id      UUID := gen_random_uuid();
    v_tmsg_id        UUID := gen_random_uuid();
    v_tattach_id     UUID := gen_random_uuid();
    v_conv_id        UUID := gen_random_uuid();
    v_cmsg_id        UUID := gen_random_uuid();
    v_tpl_id         UUID := gen_random_uuid();
    v_notif_id       UUID := gen_random_uuid();
    v_review_id      UUID := gen_random_uuid();
    v_rimg_id        UUID := gen_random_uuid();
    v_banner_id      UUID := gen_random_uuid();
    v_vbanner_id     UUID := gen_random_uuid();
    v_wishlist_id    UUID := gen_random_uuid();
    v_evtoken_id     UUID := gen_random_uuid();
    v_otp_id         UUID := gen_random_uuid();
    v_vonboard_id    UUID := gen_random_uuid();
    v_vdoc_id        UUID := gen_random_uuid();
    v_vbank_id       UUID := gen_random_uuid();
    v_vcourier_id    UUID := gen_random_uuid();

    -- ===== ledger account IDs =====
    v_acct_recv_id   UUID;
    v_acct_payable_id UUID;

    -- ===== constants =====
    v_unit_price     NUMERIC(18,2) := 185000;
    v_grand_total    NUMERIC(18,2) := 202500; -- 185000 + 15000 + 2500
    v_now            TIMESTAMPTZ := NOW();
    v_order_no       VARCHAR(30) := 'ORD-20260313-S001';
BEGIN
    -- ===== lookup role IDs =====
    SELECT id INTO v_role_umkm FROM roles WHERE code = 'umkm';
    SELECT id INTO v_role_cs   FROM roles WHERE code = 'cs';
    SELECT id INTO v_role_admin FROM roles WHERE code = 'admin';

    -- ===== lookup existing CS & admin =====
    SELECT u.id INTO v_cs_id FROM users u WHERE u.role_id = v_role_cs LIMIT 1;
    SELECT u.id INTO v_admin_id FROM users u WHERE u.role_id = v_role_admin LIMIT 1;

    -- ===== lookup ledger accounts =====
    SELECT id INTO v_acct_recv_id FROM ledger_accounts WHERE code = '1100';
    SELECT id INTO v_acct_payable_id FROM ledger_accounts WHERE code = '2100';

    -- =========================================================
    -- 1. USERS (umkm user for vendor owner)
    -- =========================================================
    INSERT INTO users (id, email, full_name, phone, password_hash, role_id, status, email_verified_at)
    VALUES (
        v_umkm_user_id, 'vendor-seed@dev.local', 'Vendor Seed Owner', '081300000099',
        '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890ABCDEF0123456',
        v_role_umkm, 'active', v_now - INTERVAL '30 days'
    );

    -- =========================================================
    -- 2. EMAIL VERIFICATION TOKEN
    -- =========================================================
    INSERT INTO email_verification_tokens (id, user_id, email, token_hash, expires_at, consumed_at)
    VALUES (
        v_evtoken_id, v_umkm_user_id, 'vendor-seed@dev.local',
        'seed_token_hash_' || md5(random()::text),
        v_now - INTERVAL '29 days',
        v_now - INTERVAL '30 days'
    );

    -- =========================================================
    -- 3. ADDRESS (customer)
    -- =========================================================
    INSERT INTO addresses (id, user_id, label, recipient_name, phone, province_id, city_id, district_id, postal_code, address_line, is_default, province_name, city_name, district_name)
    VALUES (
        v_address_id, v_customer_id, 'Rumah', 'Dev Customer', '081200000002',
        '6', '151', '2087', '11530',
        'Jl. Kebon Jeruk No.10, RT 005/RW 003',
        TRUE, 'DKI Jakarta', 'Jakarta Barat', 'Kebon Jeruk'
    );

    -- =========================================================
    -- 4. VENDOR + DOCUMENTS + BANK + BALANCE + BANNER + COURIER
    -- =========================================================
    INSERT INTO vendors (id, owner_user_id, vendor_type, legal_name, display_name, responsible_person_name, description, status, approved_by, approved_at)
    VALUES (
        v_vendor_id, v_umkm_user_id, 'souvenir_store',
        'PT Oleh-Oleh Haji Berkah', 'Toko Berkah Haji', 'Vendor Seed Owner',
        'Toko perlengkapan haji dan oleh-oleh berkualitas.',
        'active', v_admin_id, v_now - INTERVAL '25 days'
    );

    INSERT INTO vendor_documents (id, vendor_id, doc_type, file_url, mime_type, verification_status, uploaded_by, verified_by, verified_at)
    VALUES (v_vdoc_id, v_vendor_id, 'owner_document_id', 'https://placehold.co/600x400?text=KTP', 'image/jpeg', 'verified', v_umkm_user_id, v_admin_id, v_now - INTERVAL '25 days');

    INSERT INTO vendor_bank_accounts (id, vendor_id, bank_name, account_number, account_holder_name, verification_status, verified_by, verified_at)
    VALUES (v_vbank_id, v_vendor_id, 'Bank Syariah Indonesia', '7788990011', 'PT Oleh-Oleh Haji Berkah', 'verified', v_admin_id, v_now - INTERVAL '25 days');

    INSERT INTO vendor_balances (vendor_id, available_balance, pending_balance, total_earned, total_withdrawn, escrow_balance)
    VALUES (v_vendor_id, 500000, 0, 1500000, 1000000, 0);

    INSERT INTO vendor_banners (id, vendor_id, title, image_url)
    VALUES (v_vbanner_id, v_vendor_id, 'Promo Ramadhan', 'https://placehold.co/1200x400?text=Promo+Ramadhan');

    INSERT INTO vendor_couriers (id, vendor_id, courier_id, is_active)
    VALUES (v_vcourier_id, v_vendor_id, (SELECT id FROM couriers WHERE code = 'jne'), TRUE);

    -- =========================================================
    -- 5. CATEGORY + PRODUCT + VARIANT + IMAGE
    -- =========================================================
    -- reuse existing category if any
    SELECT id INTO v_category_id FROM categories WHERE is_active = TRUE LIMIT 1;
    IF v_category_id IS NULL THEN
        v_category_id := gen_random_uuid();
        INSERT INTO categories (id, parent_id, name, slug, is_active)
        VALUES (v_category_id, NULL, 'Perlengkapan Haji', 'perlengkapan-haji', TRUE);
    END IF;

    INSERT INTO products (id, vendor_id, category_id, name, slug, description, status, halal_ai_status)
    VALUES (
        v_product_id, v_vendor_id, v_category_id,
        'Kain Ihram Premium', 'kain-ihram-premium-seed',
        'Kain ihram premium berbahan katun lembut, nyaman dipakai saat ibadah haji dan umroh.',
        'published', 'passed'
    );

    INSERT INTO product_variants (id, product_id, sku, variant_name, price, currency, stock_on_hand, weight_gram, is_default, is_active)
    VALUES (v_variant_id, v_product_id, 'IHRAM-WHT-001', 'Kain Ihram Premium', v_unit_price, 'IDR', 50, 500, TRUE, TRUE);

    INSERT INTO product_images (id, product_id, image_url, mime_type, is_primary, sort_order)
    VALUES (v_image_id, v_product_id, 'https://placehold.co/600x600?text=Kain+Ihram', 'image/png', TRUE, 0);

    -- =========================================================
    -- 6. PRODUCT REVIEW + IMAGES + STATS
    -- =========================================================
    -- (review depends on order+order_item, so we insert order first, then review later)

    -- =========================================================
    -- 7. CART + CART ITEMS (converted cart)
    -- =========================================================
    INSERT INTO carts (id, user_id, status) VALUES (v_cart_id, v_customer_id, 'converted');

    INSERT INTO cart_items (id, cart_id, product_variant_id, qty) VALUES (v_cart_item_id, v_cart_id, v_variant_id, 1);

    -- =========================================================
    -- 8. ORDER + ORDER ITEMS + STATUS HISTORY
    -- =========================================================
    INSERT INTO orders (id, order_no, user_id, vendor_id, shipping_address_snapshot, order_status, payment_status, subtotal, shipping_fee, platform_fee, grand_total, placed_at, created_at, updated_at)
    VALUES (
        v_order_id, v_order_no, v_customer_id, v_vendor_id,
        '{"name":"Dev Customer","phone":"081200000002","address":"Jl. Kebon Jeruk No.10","city":"Jakarta Barat","province":"DKI Jakarta","postal_code":"11530"}'::jsonb,
        'completed', 'paid',
        v_unit_price, 15000, 2500, v_grand_total,
        v_now - INTERVAL '5 days', v_now - INTERVAL '5 days', v_now
    );

    INSERT INTO order_items (id, order_id, product_variant_id, product_name_snapshot, sku_snapshot, qty, unit_price, line_total)
    VALUES (v_order_item_id, v_order_id, v_variant_id, 'Kain Ihram Premium', 'IHRAM-WHT-001', 1, v_unit_price, v_unit_price);

    INSERT INTO order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES
        (gen_random_uuid(), v_order_id, 'pending_payment', 'paid',      NULL, v_now - INTERVAL '5 days',  'Payment confirmed via Xendit'),
        (gen_random_uuid(), v_order_id, 'paid',            'processing', NULL, v_now - INTERVAL '4 days',  'Vendor memproses pesanan'),
        (gen_random_uuid(), v_order_id, 'processing',      'packed',    NULL, v_now - INTERVAL '3 days',  'Pesanan dikemas'),
        (gen_random_uuid(), v_order_id, 'packed',          'shipped',   NULL, v_now - INTERVAL '2 days',  'Dikirim via JNE REG'),
        (gen_random_uuid(), v_order_id, 'shipped',         'completed', NULL, v_now - INTERVAL '1 day',   'Pesanan selesai');

    -- =========================================================
    -- 9. PAYMENT INVOICE + EVENT
    -- =========================================================
    INSERT INTO payment_invoices (id, order_id, gateway, xendit_invoice_id, external_invoice_id, invoice_url, payment_method, payment_channel, amount, currency, status, expires_at, paid_at, raw_payload, created_at, updated_at)
    VALUES (
        v_invoice_id, v_order_id, 'xendit', 'xnd_seed_001', 'INV-SEED-20260313-001',
        'https://checkout-staging.xendit.co/v2/seed-001', 'QRIS', 'QRIS',
        v_grand_total, 'IDR', 'paid',
        v_now - INTERVAL '4 days', v_now - INTERVAL '5 days' + INTERVAL '10 minutes',
        '{"id":"xnd_seed_001","status":"PAID"}'::jsonb,
        v_now - INTERVAL '5 days', v_now - INTERVAL '5 days'
    );

    INSERT INTO payment_events (id, payment_invoice_id, event_type, external_event_id, payload, received_at)
    VALUES (
        v_event_id, v_invoice_id, 'invoice.paid', 'evt_seed_001',
        '{"id":"xnd_seed_001","event":"invoice.paid","status":"PAID"}'::jsonb,
        v_now - INTERVAL '5 days' + INTERVAL '10 minutes'
    );

    -- =========================================================
    -- 10. SHIPMENT
    -- =========================================================
    INSERT INTO shipments (id, order_id, courier_code, service_type, tracking_no, shipment_status, shipped_at, delivered_at, etd)
    VALUES (v_shipment_id, v_order_id, 'jne', 'REG', 'JNESEED0000001', 'delivered', v_now - INTERVAL '2 days', v_now - INTERVAL '1 day', '2-3 hari');

    -- =========================================================
    -- 11. REFUND (completed refund)
    -- =========================================================
    INSERT INTO refunds (id, order_id, payment_invoice_id, amount, reason, status, requested_by, processed_by, processed_at)
    VALUES (v_refund_id, v_order_id, v_invoice_id, 50000, 'Barang tidak sesuai - refund partial', 'processed', v_customer_id, v_admin_id, v_now - INTERVAL '1 day');

    -- =========================================================
    -- 12. PAYOUT BATCH + ITEM
    -- =========================================================
    INSERT INTO payout_batches (id, vendor_id, period_start, period_end, status, total_gross, total_fee, total_net, paid_at, created_by, created_at, updated_at)
    VALUES (
        v_payout_id, v_vendor_id, (v_now - INTERVAL '30 days')::date, (v_now - INTERVAL '1 day')::date,
        'completed', v_unit_price, 2500, v_unit_price - 2500,
        v_now - INTERVAL '1 day', COALESCE(v_admin_id, v_customer_id),
        v_now - INTERVAL '2 days', v_now - INTERVAL '1 day'
    );

    INSERT INTO payout_items (id, payout_batch_id, order_id, gross_amount, platform_fee_amount, net_amount)
    VALUES (v_payout_item_id, v_payout_id, v_order_id, v_unit_price, 2500, v_unit_price - 2500);

    -- =========================================================
    -- 13. LEDGER JOURNAL + LINES
    -- =========================================================
    INSERT INTO ledger_journals (id, journal_no, source_type, source_id, event_time, description, status, created_by)
    VALUES (v_journal_id, 'JRN-SEED-001', 'payment_invoice', v_invoice_id, v_now - INTERVAL '5 days', 'Payment received for order ' || v_order_no, 'posted', v_admin_id);

    INSERT INTO ledger_lines (id, journal_id, account_id, debit, credit, currency, reference) VALUES
        (v_line_d_id, v_journal_id, v_acct_recv_id, v_grand_total, 0, 'IDR', v_order_no),
        (v_line_c_id, v_journal_id, v_acct_payable_id, 0, v_grand_total, 'IDR', v_order_no);

    -- =========================================================
    -- 14. TICKET + MESSAGE + ATTACHMENT + STATUS LOG
    -- =========================================================
    INSERT INTO tickets (id, ticket_number, customer_id, assigned_cs_id, order_number, phone, reporter_name, subject_id, subject, detail, status, source, closed_at, created_at, updated_at)
    VALUES (
        v_ticket_id, 'TKT-20260313-0001', v_customer_id, v_cs_id,
        v_order_no, '081200000002', 'Dev Customer',
        4, 'Barang Tidak Sesuai',
        'Barang yang diterima tidak sesuai dengan deskripsi produk. Warna dan ukuran berbeda dari yang dipesan.',
        'closed', 'web', v_now, v_now - INTERVAL '12 hours', v_now
    );

    INSERT INTO ticket_messages (id, ticket_id, sender_id, message, is_from_cs, is_internal_note, created_at)
    VALUES (v_tmsg_id, v_ticket_id, v_customer_id, 'Barang yang saya terima tidak sesuai dengan deskripsi. Mohon bantuan pengembaliannya.', FALSE, FALSE, v_now - INTERVAL '12 hours');

    INSERT INTO ticket_attachments (id, ticket_id, message_id, file_url, file_name, file_type, file_size)
    VALUES (v_tattach_id, v_ticket_id, v_tmsg_id, 'https://placehold.co/800x600?text=Bukti+Foto', 'bukti-foto.jpg', 'image/jpeg', 204800);

    INSERT INTO ticket_status_logs (id, ticket_id, changed_by, old_status, new_status, notes, created_at) VALUES
        (gen_random_uuid(), v_ticket_id, COALESCE(v_cs_id, v_customer_id), 'open',        'on_progress', 'CS mengambil tiket',         v_now - INTERVAL '6 hours'),
        (gen_random_uuid(), v_ticket_id, COALESCE(v_cs_id, v_customer_id), 'on_progress', 'closed',      'Masalah telah diselesaikan',  v_now);

    -- =========================================================
    -- 15. CHAT CONVERSATION + MESSAGE
    -- =========================================================
    IF v_cs_id IS NOT NULL THEN
        INSERT INTO chat_conversations (id, initiator_id, participant_id, last_message_at)
        VALUES (v_conv_id, v_customer_id, v_cs_id, v_now - INTERVAL '1 hour');

        INSERT INTO chat_messages (id, conversation_id, sender_id, message, is_read, read_at, created_at)
        VALUES (v_cmsg_id, v_conv_id, v_customer_id, 'Halo, saya mau tanya soal pesanan saya.', TRUE, v_now - INTERVAL '50 minutes', v_now - INTERVAL '1 hour');
    END IF;

    -- =========================================================
    -- 16. REPLY TEMPLATE
    -- =========================================================
    IF v_cs_id IS NOT NULL THEN
        INSERT INTO reply_templates (id, title, content, is_active, created_by, shortcut, category)
        VALUES (v_tpl_id, 'Salam Pembuka', 'Assalamualaikum, terima kasih telah menghubungi kami. Ada yang bisa kami bantu?', TRUE, v_cs_id, '/salam', 'general');
    ELSE
        INSERT INTO reply_templates (id, title, content, is_active, created_by, shortcut, category)
        VALUES (v_tpl_id, 'Salam Pembuka', 'Assalamualaikum, terima kasih telah menghubungi kami. Ada yang bisa kami bantu?', TRUE, v_customer_id, '/salam', 'general');
    END IF;

    -- =========================================================
    -- 17. FAQ
    -- =========================================================
    INSERT INTO faqs (category, question, answer)
    VALUES ('Umum', 'Bagaimana cara memesan perlengkapan haji?', 'Anda bisa memesan melalui aplikasi atau website kami. Pilih produk, masukkan ke keranjang, lalu lakukan checkout.')
    ON CONFLICT DO NOTHING;

    -- =========================================================
    -- 18. NOTIFICATION
    -- =========================================================
    INSERT INTO notifications (id, user_id, type, title, message)
    VALUES (v_notif_id, v_customer_id, 'order', 'Pesanan Selesai', 'Pesanan ' || v_order_no || ' telah selesai. Terima kasih telah berbelanja!');

    -- =========================================================
    -- 19. PRODUCT REVIEW + IMAGE + STATS
    -- =========================================================
    INSERT INTO product_reviews (id, product_id, order_id, order_item_id, user_id, rating, review_text, status)
    VALUES (v_review_id, v_product_id, v_order_id, v_order_item_id, v_customer_id, 4, 'Kainnya cukup bagus dan lembut, tapi warna agak berbeda dari foto. Overall puas.', 'published');

    INSERT INTO product_review_images (id, review_id, object_key, mime_type, file_size_bytes, sort_order)
    VALUES (v_rimg_id, v_review_id, 'reviews/seed-review-001.jpg', 'image/jpeg', 153600, 0);

    INSERT INTO product_review_stats (product_id, total_reviews, total_stars, average_rating, star_0_count, star_1_count, star_2_count, star_3_count, star_4_count, star_5_count)
    VALUES (v_product_id, 1, 4, 4.00, 0, 0, 0, 0, 1, 0);

    -- =========================================================
    -- 20. BANNER
    -- =========================================================
    INSERT INTO banners (id, title, image_url)
    VALUES (v_banner_id, 'Selamat Datang di Hajj Store', 'https://placehold.co/1200x400?text=Hajj+Store+Banner');

    -- =========================================================
    -- 21. RETURN REASONS
    -- =========================================================
    INSERT INTO return_reasons (reason) VALUES ('Barang tidak sesuai deskripsi') ON CONFLICT DO NOTHING;

    -- =========================================================
    -- 22. ADMIN CONTACTS
    -- =========================================================
    INSERT INTO admin_contacts (content) VALUES ('Hubungi admin di admin@hajjstore.id atau WhatsApp 08123456789') ON CONFLICT DO NOTHING;

    -- =========================================================
    -- 23. WISHLIST ITEM
    -- =========================================================
    INSERT INTO wishlist_items (id, user_id, product_id) VALUES (v_wishlist_id, v_customer_id, v_product_id);

    -- =========================================================
    -- 24. OTP CODE (expired, consumed)
    -- =========================================================
    INSERT INTO otp_codes (id, user_id, email, purpose, channel, code_hash, expires_at, consumed_at)
    VALUES (v_otp_id, v_umkm_user_id, 'vendor-seed@dev.local', 'vendor_onboarding', 'email', md5('123456'), v_now - INTERVAL '29 days', v_now - INTERVAL '30 days');

    -- =========================================================
    -- 25. VENDOR ONBOARDING (completed)
    -- =========================================================
    INSERT INTO vendor_onboardings (id, email, status, password_hash, store_name, vendor_type, business_legal_type, document_id_type, nik, owner_name, birth_date, otp_verified_at, completed_at)
    VALUES (
        v_vonboard_id, 'vendor-seed@dev.local', 'completed',
        '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890ABCDEF0123456',
        'Toko Berkah Haji', 'souvenir_store', 'perorangan', 'ktp',
        '3201010101900001', 'Vendor Seed Owner', '1990-01-01',
        v_now - INTERVAL '31 days', v_now - INTERVAL '30 days'
    );

    -- =========================================================
    -- 26. VENDOR WITHDRAWAL
    -- =========================================================
    INSERT INTO vendor_withdrawals (id, vendor_id, amount, channel_code, status, description, fee_estimated, total_deducted, fee_status, created_at, updated_at)
    VALUES (
        gen_random_uuid(), v_vendor_id, 500000, 'ID_BCA', 'completed',
        'Penarikan saldo vendor', 2500, 502500,
        'resolved', v_now - INTERVAL '3 days', v_now - INTERVAL '3 days'
    );

    RAISE NOTICE '';
    RAISE NOTICE '============================================';
    RAISE NOTICE '  SEED COMPLETE - All tables populated!';
    RAISE NOTICE '============================================';
    RAISE NOTICE '  UMKM User  : %', v_umkm_user_id;
    RAISE NOTICE '  Vendor     : %', v_vendor_id;
    RAISE NOTICE '  Product    : Kain Ihram Premium (%)', v_product_id;
    RAISE NOTICE '  Order      : %', v_order_no;
    RAISE NOTICE '  Payment    : INV-SEED-20260313-001 (QRIS, paid)';
    RAISE NOTICE '  Shipment   : JNESEED0000001 (JNE REG, delivered)';
    RAISE NOTICE '  Ticket     : TKT-20260313-0001 (closed)';
    RAISE NOTICE '  Review     : 4 stars';
    RAISE NOTICE '============================================';
END $$;

COMMIT;
