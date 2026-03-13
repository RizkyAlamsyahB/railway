--
-- PostgreSQL database dump
--

\restrict X2QszTK48syofOONPN8pyTEjpmRpWHpFmM0v3d4eOzc8nsI6HC9ma11eR7chIZm

-- Dumped from database version 17.9
-- Dumped by pg_dump version 17.9

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.roles VALUES (1, 'admin', 'Administrator');
INSERT INTO public.roles VALUES (2, 'umkm', 'UMKM');
INSERT INTO public.roles VALUES (3, 'customer', 'Customer');
INSERT INTO public.roles VALUES (4, 'cs', 'Customer Service');
INSERT INTO public.roles VALUES (5, 'finance', 'Finance');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users VALUES ('aa000002-0000-0000-0000-000000000001', 'umkm2@dev.local', 'Toko Oleh-Oleh Tanah Suci', NULL, '081299000002', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('aa000002-0000-0000-0000-000000000002', 'umkm3@dev.local', 'Perlengkapan Ibadah Barokah', NULL, '081299000003', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('aa000002-0000-0000-0000-000000000003', 'customer2@dev.local', 'Siti Aisyah', NULL, '081299000004', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 3, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('7907d953-ba02-40a4-b403-07b68d8471e8', 'admin@example.com', 'Admin', NULL, '08120000000', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 1, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('cc000001-0000-0000-0000-000000000002', 'cs2@dev.local', 'Dev Customer Service 2', NULL, '081200000005', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('cc000001-0000-0000-0000-000000000003', 'cs3@dev.local', 'Dev Customer Service 3', NULL, '081200000006', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'ah.nursufyantsauri@gmail.com', 'Dev UMKM', NULL, '081200000001', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'rizkyalamsyah.dev@gmail.com', 'Dev Customer', NULL, '081200000002', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 3, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('cc000001-0000-0000-0000-000000000001', 'adkhawildanrizqia@gmail.com', 'Dev Customer Service 1', NULL, '081200000003', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('ff000001-0000-0000-0000-000000000001', 'galangarsandy@gmail.com', 'Dev Finance', NULL, '081200000004', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 5, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.addresses VALUES ('e081dc95-612f-4d08-ab35-21fe7b7b998a', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'TOKO WAREHOUSE SEDATI - SIDOARJO', 'TOKO WAREHOUSE SEDATI - SIDOARJO', '062916302649', '18', '583', '6001', '61253', 'WAREHOUSE SEDATI', true, 'JAWA TIMUR', ' SIDOARJO', 'SEDATI', '70995', 'SEDATI GEDE', NULL, NULL, NULL, '2026-03-12 04:11:45.061823+00', '2026-03-12 04:12:09.361232+00');
INSERT INTO public.addresses VALUES ('a0794254-012a-4381-8f22-4bab020ed420', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'ALAMAT PEMBELI - WIYUNG', 'RIZKY', '098765433210', '18', '577', '5901', '60228', 'Jl.xxxx.xxx No.01', true, 'JAWA TIMUR', ' SURABAYA', 'WIYUNG', '69354', 'WIYUNG', NULL, NULL, NULL, '2026-03-12 04:31:09.785581+00', '2026-03-12 04:31:09.785581+00');


--
-- Data for Name: admin_contacts; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: banners; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.carts VALUES ('38f3e138-58d9-4941-b948-ff7c252a8b0b', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-12 04:29:36.202609+00', '2026-03-12 04:31:46.662441+00');


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000001', NULL, 'Perlengkapan Haji & Umrah', 'perlengkapan-haji-umrah', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000002', NULL, 'Oleh-oleh & Souvenir', 'oleh-oleh-souvenir', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000003', NULL, 'Fashion Muslim', 'fashion-muslim', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000011', 'ca000001-0000-0000-0000-000000000001', 'Pakaian Ihram', 'pakaian-ihram', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000012', 'ca000001-0000-0000-0000-000000000001', 'Sajadah & Alat Shalat', 'sajadah-alat-shalat', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000013', 'ca000001-0000-0000-0000-000000000001', 'Tasbih & Hampers', 'tasbih-hampers', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000021', 'ca000001-0000-0000-0000-000000000002', 'Kurma', 'kurma', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000022', 'ca000001-0000-0000-0000-000000000002', 'Minyak & Herbal', 'minyak-herbal', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000023', 'ca000001-0000-0000-0000-000000000002', 'Air Zamzam', 'air-zamzam', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000031', 'ca000001-0000-0000-0000-000000000003', 'Mukena', 'mukena', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000032', 'ca000001-0000-0000-0000-000000000003', 'Gamis & Jubah', 'gamis-jubah', true);


--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendors VALUES ('9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'souvenir_store', 'PT Nursufyan Sauri', 'TOKO Nursufyan Sauri', 'Nursufyan Sauri', 'Toko milik Nursufyan Sauri', 'active', '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-12 04:07:38.085199+00', NULL, 'dummy_xendit_123', '2026-03-12 04:05:34.329709+00', '2026-03-12 04:05:34.329709+00');


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.products VALUES ('0f0ce953-d4a0-4dcb-a621-a53656b69092', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000001', 'Paket Haji Premium', 'paket-haji-premium', 'Paket perlengkapan haji lengkap dan premium.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products VALUES ('2aa62691-ea2c-491e-8456-031984b79686', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000002', 'Souvenir Mekkah', 'souvenir-mekkah', 'Souvenir khas Mekkah untuk oleh-oleh.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products VALUES ('72f634dd-9f2e-429e-9e81-b5703f8b06f8', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000003', 'Baju Muslim Pria', 'baju-muslim-pria', 'Baju muslim pria modern dan nyaman.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products VALUES ('8bb433f1-55c1-40ce-a50b-158150e5cab9', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000011', 'Pakaian Ihram Dewasa', 'pakaian-ihram-dewasa', 'Pakaian ihram dewasa bahan premium.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');


--
-- Data for Name: product_variants; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.product_variants VALUES ('747b1fc2-edaf-424b-9358-92526ac05a6e', '0f0ce953-d4a0-4dcb-a621-a53656b69092', 'SKU-HAJI-001', 'Default', 1500000.00, 'IDR', 10, 1200, true, true);
INSERT INTO public.product_variants VALUES ('14218407-1b19-40bd-9136-d611d927e573', '72f634dd-9f2e-429e-9e81-b5703f8b06f8', 'SKU-BAJU-001', 'Default', 250000.00, 'IDR', 20, 400, true, true);
INSERT INTO public.product_variants VALUES ('5f65a8ed-3c03-4cae-822f-4c97980db7f8', '8bb433f1-55c1-40ce-a50b-158150e5cab9', 'SKU-IHRAM-001', 'Default', 350000.00, 'IDR', 15, 800, true, true);
INSERT INTO public.product_variants VALUES ('c9f199bd-ea95-442b-a16c-c9fd0769a32f', '2aa62691-ea2c-491e-8456-031984b79686', 'SKU-SOUV-001', 'Default', 50000.00, 'IDR', 49, 200, true, true);


--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.cart_items VALUES ('208e50c9-2003-4eb5-94d1-2b423d5e64dd', '38f3e138-58d9-4941-b948-ff7c252a8b0b', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 1, '2026-03-12 04:29:36.216061+00');


--
-- Data for Name: chat_conversations; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: couriers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.couriers VALUES (1, 'jne', 'JNE', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (2, 'sicepat', 'SiCepat', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (3, 'ide', 'IDExpress', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (4, 'sap', 'SAP Express', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (5, 'ninja', 'Ninja', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (6, 'jnt', 'J&T Express', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (7, 'tiki', 'TIKI', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (8, 'wahana', 'Wahana Express', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (9, 'pos', 'POS Indonesia', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (10, 'sentral', 'Sentral Cargo', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (11, 'lion', 'Lion Parcel', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers VALUES (12, 'rex', 'Royal Express Asia', NULL, true, '2026-03-12 03:55:18.668697+00');


--
-- Data for Name: email_verification_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: faqs; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: ledger_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ledger_accounts VALUES ('60825dd3-e090-458e-a67b-9393b9834b10', '1100', 'Payment Gateway Receivable', 'asset', 'D', true);
INSERT INTO public.ledger_accounts VALUES ('8e9d7291-ced8-4ca5-86ae-c4739e0fefa4', '2100', 'Vendor Payable', 'liability', 'C', true);
INSERT INTO public.ledger_accounts VALUES ('ccf126d7-36dd-4d53-845f-a35da76acaa7', '4100', 'Platform Fee Revenue', 'revenue', 'C', true);
INSERT INTO public.ledger_accounts VALUES ('240bbd21-9b1f-42ac-b975-9ddff62ead79', '4200', 'Admin Fee Revenue', 'revenue', 'C', true);


--
-- Data for Name: ledger_journals; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: ledger_lines; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.orders VALUES ('382762cd-f805-41a0-ad67-07c6a33020d2', 'ORD-20260312-5F560F3E', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '{"label": "ALAMAT PEMBELI - WIYUNG", "phone": "098765433210", "city_name": " SURABAYA", "address_id": "a0794254-012a-4381-8f22-4bab020ed420", "is_default": true, "postal_code": "60228", "address_line": "Jl.xxxx.xxx No.01", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY", "subdistrict_name": "WIYUNG"}', 'pending_payment', 'unpaid', 50000.00, 8000.00, 6000.00, 64000.00, '2026-03-12 04:31:46.178518+00', '2026-03-12 04:31:46.178518+00', '2026-03-12 04:31:46.178518+00');


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.order_items VALUES ('28633b60-c64d-4859-a68d-16aba21ac12a', '382762cd-f805-41a0-ad67-07c6a33020d2', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 'Souvenir Mekkah', 'SKU-SOUV-001', 1, 50000.00, 50000.00);


--
-- Data for Name: order_status_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.order_status_history VALUES ('35825d8a-438f-4cdb-a0c5-4c910d3f83d5', '382762cd-f805-41a0-ad67-07c6a33020d2', NULL, 'pending_payment', NULL, '2026-03-12 04:31:46.641578+00', NULL);


--
-- Data for Name: payment_invoices; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.payment_invoices VALUES ('768ef5ce-15f8-476b-a4de-8d6c9b89538f', '382762cd-f805-41a0-ad67-07c6a33020d2', 'xendit', 'noop-inv-efc97e93', 'INV-ORD-20260312-5F560F3E-13850065', 'https://checkout-bypass.example.com/noop-inv-efc97e93', NULL, NULL, 64000.00, 'IDR', 'pending', '2026-03-13 04:31:46+00', NULL, NULL, '2026-03-12 04:31:46.178518+00', '2026-03-12 04:31:46.178518+00');


--
-- Data for Name: payment_events; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: payout_batches; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: payout_items; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: product_reviews; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: product_review_images; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: product_review_stats; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: refunds; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: reply_templates; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: return_reasons; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.schema_migrations VALUES (40, false);
INSERT INTO public.schema_migrations VALUES (33, false);


--
-- Data for Name: shipments; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.shipments VALUES ('04c23c0f-1620-400c-b27a-949a27a5309e', '382762cd-f805-41a0-ad67-07c6a33020d2', 'jnt', 'EZ', '', 'waiting_pickup', NULL, NULL, '');


--
-- Data for Name: ticket_subjects; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ticket_subjects VALUES (1, 'Barang Tidak Sampai', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (2, 'Barang Rusak', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (3, 'Pengembalian Dana', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (4, 'Barang Tidak Sesuai', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (5, 'Pengiriman Terlambat', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (6, 'Pembatalan Pesanan', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (7, 'Kesalahan Produk', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (8, 'Akun Bermasalah', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (9, 'Pembayaran Gagal', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (10, 'Voucher/Promo Tidak Berlaku', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (11, 'Pertanyaan Umum', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects VALUES (12, 'Lainnya', true, '2026-03-12 03:55:18.534203+00');


--
-- Data for Name: tickets; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: ticket_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: ticket_attachments; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: ticket_status_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: vendor_balances; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_balances VALUES ('9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 1000000.00, 0.00, 1000000.00, 0.00, '2026-03-12 04:17:46.056983+00', 0.00);


--
-- Data for Name: vendor_bank_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_bank_accounts VALUES ('7b2aee15-39da-4080-bf21-df2feb5c791e', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'BANK Dummy', '1234567890', 'Nursufyan Sauri', 'verified', NULL, '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-12 04:17:26.124526+00', '2026-03-12 04:17:26.124526+00', '2026-03-12 04:17:26.124526+00');


--
-- Data for Name: vendor_banners; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: vendor_couriers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_couriers VALUES ('04737655-d896-4a11-aa1d-567769010fcb', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 3, true, '2026-03-12 04:12:25.593572+00');
INSERT INTO public.vendor_couriers VALUES ('609cc811-d7c1-4ab9-a8d7-c45fb66a03da', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 6, true, '2026-03-12 04:12:25.594108+00');
INSERT INTO public.vendor_couriers VALUES ('ee3d321e-42f7-4b6a-b99b-fb52fef0193c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 1, true, '2026-03-12 04:12:25.594364+00');
INSERT INTO public.vendor_couriers VALUES ('83f69d0e-27bf-43d5-895c-597814d0f7ec', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 11, true, '2026-03-12 04:12:25.594595+00');
INSERT INTO public.vendor_couriers VALUES ('f56186f5-f341-45e3-81ca-1a7b471af801', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 5, true, '2026-03-12 04:12:25.594792+00');
INSERT INTO public.vendor_couriers VALUES ('bfe88d97-d13b-4ae5-9367-aad749d0634c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 9, true, '2026-03-12 04:12:25.59502+00');
INSERT INTO public.vendor_couriers VALUES ('5203e452-d7a8-4f60-8983-ee01d23b33b9', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 12, true, '2026-03-12 04:12:25.595229+00');
INSERT INTO public.vendor_couriers VALUES ('813352d8-5fec-4e1d-aa99-a73bc4fb1eb5', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 4, true, '2026-03-12 04:12:25.595467+00');
INSERT INTO public.vendor_couriers VALUES ('8cf8a32c-f3e8-492f-9b0f-dd9889fc8ed5', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 10, true, '2026-03-12 04:12:25.595702+00');
INSERT INTO public.vendor_couriers VALUES ('7fe72ce5-4807-4dee-b9d2-f2892c8cbe51', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 2, true, '2026-03-12 04:12:25.595913+00');
INSERT INTO public.vendor_couriers VALUES ('91ce5cf3-18ba-4fb0-bac8-fcbc9fe471ec', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 7, true, '2026-03-12 04:12:25.596118+00');
INSERT INTO public.vendor_couriers VALUES ('f371e602-a0a2-429a-b962-256df2ee211c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 8, true, '2026-03-12 04:12:25.596361+00');


--
-- Data for Name: vendor_documents; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: vendor_withdrawals; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Data for Name: wishlist_items; Type: TABLE DATA; Schema: public; Owner: postgres
--



--
-- Name: admin_contacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.admin_contacts_id_seq', 1, false);


--
-- Name: couriers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.couriers_id_seq', 12, true);


--
-- Name: faqs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.faqs_id_seq', 1, false);


--
-- Name: return_reasons_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.return_reasons_id_seq', 1, false);


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.roles_id_seq', 5, true);


--
-- Name: ticket_subjects_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.ticket_subjects_id_seq', 12, true);


--
-- PostgreSQL database dump complete
--

\unrestrict X2QszTK48syofOONPN8pyTEjpmRpWHpFmM0v3d4eOzc8nsI6HC9ma11eR7chIZm

