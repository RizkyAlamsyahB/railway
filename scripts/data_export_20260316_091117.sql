--
-- PostgreSQL database dump
--

\restrict uwmddg8yAn0SoGa7WelfBPBZoSrKDsVXNxdSEnzzUItuEqOEAlQJcdtA0JzD0Hg

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

INSERT INTO public.roles (id, code, name) VALUES (1, 'admin', 'Administrator');
INSERT INTO public.roles (id, code, name) VALUES (2, 'umkm', 'UMKM');
INSERT INTO public.roles (id, code, name) VALUES (3, 'customer', 'Customer');
INSERT INTO public.roles (id, code, name) VALUES (4, 'cs', 'Customer Service');
INSERT INTO public.roles (id, code, name) VALUES (5, 'finance', 'Finance');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('aa000002-0000-0000-0000-000000000001', 'umkm2@dev.local', 'Toko Oleh-Oleh Tanah Suci', NULL, '081299000002', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('aa000002-0000-0000-0000-000000000002', 'umkm3@dev.local', 'Perlengkapan Ibadah Barokah', NULL, '081299000003', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('aa000002-0000-0000-0000-000000000003', 'customer2@dev.local', 'Siti Aisyah', NULL, '081299000004', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 3, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('7907d953-ba02-40a4-b403-07b68d8471e8', 'admin@example.com', 'Admin', NULL, '08120000000', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 1, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('cc000001-0000-0000-0000-000000000002', 'cs2@dev.local', 'Dev Customer Service 2', NULL, '081200000005', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('cc000001-0000-0000-0000-000000000003', 'cs3@dev.local', 'Dev Customer Service 3', NULL, '081200000006', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'ah.nursufyantsauri@gmail.com', 'Dev UMKM', NULL, '081200000001', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'rizkyalamsyah.dev@gmail.com', 'Dev Customer', NULL, '081200000002', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 3, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('cc000001-0000-0000-0000-000000000001', 'adkhawildanrizqia@gmail.com', 'Dev Customer Service 1', NULL, '081200000003', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('ff000001-0000-0000-0000-000000000001', 'galangarsandy@gmail.com', 'Dev Finance', NULL, '081200000004', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 5, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at, image_url) VALUES ('e166571c-3c00-488b-9844-ec536a54f332', 'vendor-seed@dev.local', 'Vendor Seed Owner', NULL, '081300000099', '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890ABCDEF0123456', 2, 'active', '2026-02-11 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', NULL);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.addresses (id, user_id, label, recipient_name, phone, province_id, city_id, district_id, postal_code, address_line, is_default, province_name, city_name, district_name, subdistrict_id, subdistrict_name, notes, latitude, longitude, created_at, updated_at) VALUES ('e081dc95-612f-4d08-ab35-21fe7b7b998a', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'TOKO WAREHOUSE SEDATI - SIDOARJO', 'TOKO WAREHOUSE SEDATI - SIDOARJO', '062916302649', '18', '583', '6001', '61253', 'WAREHOUSE SEDATI', true, 'JAWA TIMUR', ' SIDOARJO', 'SEDATI', '70995', 'SEDATI GEDE', NULL, NULL, NULL, '2026-03-12 04:11:45.061823+00', '2026-03-12 04:12:09.361232+00');
INSERT INTO public.addresses (id, user_id, label, recipient_name, phone, province_id, city_id, district_id, postal_code, address_line, is_default, province_name, city_name, district_name, subdistrict_id, subdistrict_name, notes, latitude, longitude, created_at, updated_at) VALUES ('a0794254-012a-4381-8f22-4bab020ed420', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'ALAMAT PEMBELI - WIYUNG', 'RIZKY', '098765433210', '18', '577', '5901', '60228', 'Jl.xxxx.xxx No.01', true, 'JAWA TIMUR', ' SURABAYA', 'WIYUNG', '69354', 'WIYUNG', NULL, NULL, NULL, '2026-03-12 04:31:09.785581+00', '2026-03-12 04:31:09.785581+00');
INSERT INTO public.addresses (id, user_id, label, recipient_name, phone, province_id, city_id, district_id, postal_code, address_line, is_default, province_name, city_name, district_name, subdistrict_id, subdistrict_name, notes, latitude, longitude, created_at, updated_at) VALUES ('399e90f7-0e8c-41f8-8032-1eda354f36d3', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Rumah', 'Dev Customer', '081200000002', '18', '577', '5901', '60228', 'Jl. Kebon Jeruk No.10, RT 005/RW 003', true, 'JAWA TIMUR', ' SURABAYA', 'WIYUNG', '69354', 'WIYUNG', NULL, NULL, NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: admin_contacts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.admin_contacts (id, content, created_at, updated_at) VALUES (1, 'Hubungi admin di admin@hajjstore.id atau WhatsApp 08123456789', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: banners; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.banners (id, title, image_url, created_at, updated_at) VALUES ('4190f3f4-45aa-482b-9121-1836b127304c', 'Selamat Datang di Hajj Store', 'https://placehold.co/1200x400?text=Hajj+Store+Banner', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.carts (id, user_id, status, created_at, updated_at) VALUES ('38f3e138-58d9-4941-b948-ff7c252a8b0b', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-12 04:29:36.202609+00', '2026-03-12 04:31:46.662441+00');
INSERT INTO public.carts (id, user_id, status, created_at, updated_at) VALUES ('b81c7b33-534e-4b8e-a50d-e7758fdfb3ab', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'active', '2026-03-12 04:54:53.927372+00', '2026-03-12 04:54:53.927372+00');
INSERT INTO public.carts (id, user_id, status, created_at, updated_at) VALUES ('2a64b014-85a2-4c73-8f07-97cee7a652d1', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000001', NULL, 'Perlengkapan Haji & Umrah', 'perlengkapan-haji-umrah', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000002', NULL, 'Oleh-oleh & Souvenir', 'oleh-oleh-souvenir', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000003', NULL, 'Fashion Muslim', 'fashion-muslim', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000011', 'ca000001-0000-0000-0000-000000000001', 'Pakaian Ihram', 'pakaian-ihram', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000012', 'ca000001-0000-0000-0000-000000000001', 'Sajadah & Alat Shalat', 'sajadah-alat-shalat', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000013', 'ca000001-0000-0000-0000-000000000001', 'Tasbih & Hampers', 'tasbih-hampers', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000021', 'ca000001-0000-0000-0000-000000000002', 'Kurma', 'kurma', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000022', 'ca000001-0000-0000-0000-000000000002', 'Minyak & Herbal', 'minyak-herbal', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000023', 'ca000001-0000-0000-0000-000000000002', 'Air Zamzam', 'air-zamzam', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000031', 'ca000001-0000-0000-0000-000000000003', 'Mukena', 'mukena', true);
INSERT INTO public.categories (id, parent_id, name, slug, is_active) VALUES ('ca000001-0000-0000-0000-000000000032', 'ca000001-0000-0000-0000-000000000003', 'Gamis & Jubah', 'gamis-jubah', true);


--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendors (id, owner_user_id, vendor_type, legal_name, display_name, responsible_person_name, description, status, approved_by, approved_at, status_reason, xendit_account_id, created_at, updated_at) VALUES ('9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'souvenir_store', 'PT Nursufyan Sauri', 'TOKO Nursufyan Sauri', 'Nursufyan Sauri', 'Toko milik Nursufyan Sauri', 'active', '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-12 04:07:38.085199+00', NULL, 'dummy_xendit_123', '2026-03-12 04:05:34.329709+00', '2026-03-12 04:05:34.329709+00');
INSERT INTO public.vendors (id, owner_user_id, vendor_type, legal_name, display_name, responsible_person_name, description, status, approved_by, approved_at, status_reason, xendit_account_id, created_at, updated_at) VALUES ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'e166571c-3c00-488b-9844-ec536a54f332', 'souvenir_store', 'PT Oleh-Oleh Haji Berkah', 'Toko Berkah Haji', 'Vendor Seed Owner', 'Toko perlengkapan haji dan oleh-oleh berkualitas.', 'active', '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-02-16 07:24:21.226459+00', NULL, 'dummy_xendit_123', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.products (id, vendor_id, category_id, name, slug, description, status, halal_ai_status, halal_ai_notes, created_at, updated_at) VALUES ('0f0ce953-d4a0-4dcb-a621-a53656b69092', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000001', 'Paket Haji Premium', 'paket-haji-premium', 'Paket perlengkapan haji lengkap dan premium.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products (id, vendor_id, category_id, name, slug, description, status, halal_ai_status, halal_ai_notes, created_at, updated_at) VALUES ('2aa62691-ea2c-491e-8456-031984b79686', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000002', 'Souvenir Mekkah', 'souvenir-mekkah', 'Souvenir khas Mekkah untuk oleh-oleh.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products (id, vendor_id, category_id, name, slug, description, status, halal_ai_status, halal_ai_notes, created_at, updated_at) VALUES ('72f634dd-9f2e-429e-9e81-b5703f8b06f8', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000003', 'Baju Muslim Pria', 'baju-muslim-pria', 'Baju muslim pria modern dan nyaman.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products (id, vendor_id, category_id, name, slug, description, status, halal_ai_status, halal_ai_notes, created_at, updated_at) VALUES ('8bb433f1-55c1-40ce-a50b-158150e5cab9', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000011', 'Pakaian Ihram Dewasa', 'pakaian-ihram-dewasa', 'Pakaian ihram dewasa bahan premium.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products (id, vendor_id, category_id, name, slug, description, status, halal_ai_status, halal_ai_notes, created_at, updated_at) VALUES ('b2999c9e-df33-48ef-a109-38c77e51f25d', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'ca000001-0000-0000-0000-000000000001', 'Kain Ihram Premium', 'kain-ihram-premium-seed', 'Kain ihram premium berbahan katun lembut, nyaman dipakai saat ibadah haji dan umroh.', 'published', 'passed', NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_variants; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.product_variants (id, product_id, sku, variant_name, price, currency, stock_on_hand, weight_gram, is_default, is_active) VALUES ('747b1fc2-edaf-424b-9358-92526ac05a6e', '0f0ce953-d4a0-4dcb-a621-a53656b69092', 'SKU-HAJI-001', 'Default', 1500000.00, 'IDR', 10, 1200, true, true);
INSERT INTO public.product_variants (id, product_id, sku, variant_name, price, currency, stock_on_hand, weight_gram, is_default, is_active) VALUES ('14218407-1b19-40bd-9136-d611d927e573', '72f634dd-9f2e-429e-9e81-b5703f8b06f8', 'SKU-BAJU-001', 'Default', 250000.00, 'IDR', 20, 400, true, true);
INSERT INTO public.product_variants (id, product_id, sku, variant_name, price, currency, stock_on_hand, weight_gram, is_default, is_active) VALUES ('5f65a8ed-3c03-4cae-822f-4c97980db7f8', '8bb433f1-55c1-40ce-a50b-158150e5cab9', 'SKU-IHRAM-001', 'Default', 350000.00, 'IDR', 15, 800, true, true);
INSERT INTO public.product_variants (id, product_id, sku, variant_name, price, currency, stock_on_hand, weight_gram, is_default, is_active) VALUES ('c9f199bd-ea95-442b-a16c-c9fd0769a32f', '2aa62691-ea2c-491e-8456-031984b79686', 'SKU-SOUV-001', 'Default', 50000.00, 'IDR', 49, 200, true, true);
INSERT INTO public.product_variants (id, product_id, sku, variant_name, price, currency, stock_on_hand, weight_gram, is_default, is_active) VALUES ('dcc2c0cd-35f9-4744-a7d1-542e54b2743f', 'b2999c9e-df33-48ef-a109-38c77e51f25d', 'IHRAM-WHT-001', 'Kain Ihram Premium', 185000.00, 'IDR', 50, 500, true, true);


--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.cart_items (id, cart_id, product_variant_id, qty, created_at) VALUES ('208e50c9-2003-4eb5-94d1-2b423d5e64dd', '38f3e138-58d9-4941-b948-ff7c252a8b0b', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 1, '2026-03-12 04:29:36.216061+00');
INSERT INTO public.cart_items (id, cart_id, product_variant_id, qty, created_at) VALUES ('c28b28c8-e1e1-45dc-aead-d42be0ee5ecf', '2a64b014-85a2-4c73-8f07-97cee7a652d1', 'dcc2c0cd-35f9-4744-a7d1-542e54b2743f', 1, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: chat_conversations; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.chat_conversations (id, initiator_id, participant_id, last_message_at, created_at, updated_at) VALUES ('45ec9ef1-bd88-4dd0-8208-3b5c0e5e6eef', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'cc000001-0000-0000-0000-000000000002', '2026-03-13 06:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.chat_messages (id, conversation_id, sender_id, message, attachment_url, is_read, read_at, created_at) VALUES ('ae8f5ec0-04ed-4308-a818-b600da71ec8e', '45ec9ef1-bd88-4dd0-8208-3b5c0e5e6eef', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Halo, saya mau tanya soal pesanan saya.', NULL, true, '2026-03-13 06:34:21.226459+00', '2026-03-13 06:24:21.226459+00');


--
-- Data for Name: couriers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (1, 'jne', 'JNE', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (2, 'sicepat', 'SiCepat', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (3, 'ide', 'IDExpress', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (4, 'sap', 'SAP Express', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (5, 'ninja', 'Ninja', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (6, 'jnt', 'J&T Express', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (7, 'tiki', 'TIKI', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (8, 'wahana', 'Wahana Express', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (9, 'pos', 'POS Indonesia', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (10, 'sentral', 'Sentral Cargo', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (11, 'lion', 'Lion Parcel', NULL, true, '2026-03-12 03:55:18.668697+00');
INSERT INTO public.couriers (id, code, name, logo_url, is_active, created_at) VALUES (12, 'rex', 'Royal Express Asia', NULL, true, '2026-03-12 03:55:18.668697+00');


--
-- Data for Name: email_verification_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.email_verification_tokens (id, user_id, email, token_hash, expires_at, consumed_at, invalidated_at, created_at) VALUES ('779944a5-b096-4417-a389-8e28921bbd47', 'e166571c-3c00-488b-9844-ec536a54f332', 'vendor-seed@dev.local', 'seed_token_hash_0c538b83df4bdf8e4c3cd73df13b09b5', '2026-02-12 07:24:21.226459+00', '2026-02-11 07:24:21.226459+00', NULL, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: faqs; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.faqs (id, category, question, answer, created_at, updated_at) VALUES (1, 'Umum', 'Bagaimana cara memesan perlengkapan haji?', 'Anda bisa memesan melalui aplikasi atau website kami. Pilih produk, masukkan ke keranjang, lalu lakukan checkout.', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: ledger_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ledger_accounts (id, code, name, account_type, normal_side, is_active) VALUES ('60825dd3-e090-458e-a67b-9393b9834b10', '1100', 'Payment Gateway Receivable', 'asset', 'D', true);
INSERT INTO public.ledger_accounts (id, code, name, account_type, normal_side, is_active) VALUES ('8e9d7291-ced8-4ca5-86ae-c4739e0fefa4', '2100', 'Vendor Payable', 'liability', 'C', true);
INSERT INTO public.ledger_accounts (id, code, name, account_type, normal_side, is_active) VALUES ('ccf126d7-36dd-4d53-845f-a35da76acaa7', '4100', 'Platform Fee Revenue', 'revenue', 'C', true);
INSERT INTO public.ledger_accounts (id, code, name, account_type, normal_side, is_active) VALUES ('240bbd21-9b1f-42ac-b975-9ddff62ead79', '4200', 'Admin Fee Revenue', 'revenue', 'C', true);


--
-- Data for Name: ledger_journals; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ledger_journals (id, journal_no, source_type, source_id, event_time, description, status, created_by, created_at) VALUES ('f9f10cdd-475e-4ab6-9b63-8b5be5f4bae8', 'JRN-SEED-001', 'payment_invoice', 'a8a18b80-5aa5-4237-884d-35fd3f92c386', '2026-03-08 07:24:21.226459+00', 'Payment received for order ORD-20260313-S001', 'posted', '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: ledger_lines; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ledger_lines (id, journal_id, account_id, debit, credit, currency, reference) VALUES ('3a008fb2-5a0b-4134-8f14-cba10a1fcbec', 'f9f10cdd-475e-4ab6-9b63-8b5be5f4bae8', '60825dd3-e090-458e-a67b-9393b9834b10', 202500.00, 0.00, 'IDR', 'ORD-20260313-S001');
INSERT INTO public.ledger_lines (id, journal_id, account_id, debit, credit, currency, reference) VALUES ('46fb3921-36c5-4247-bee8-51677fe5198f', 'f9f10cdd-475e-4ab6-9b63-8b5be5f4bae8', '8e9d7291-ced8-4ca5-86ae-c4739e0fefa4', 0.00, 202500.00, 'IDR', 'ORD-20260313-S001');


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.notifications (id, user_id, type, title, message, created_at, read_at) VALUES ('f2bf6aa4-b496-442c-8fe8-2bcdc3e0d765', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'order', 'Pesanan Selesai', 'Pesanan ORD-20260313-S001 telah selesai. Terima kasih telah berbelanja!', '2026-03-13 07:24:21.226459+00', NULL);


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.orders (id, order_no, user_id, vendor_id, shipping_address_snapshot, order_status, payment_status, subtotal, shipping_fee, platform_fee, grand_total, placed_at, created_at, updated_at) VALUES ('382762cd-f805-41a0-ad67-07c6a33020d2', 'ORD-20260312-5F560F3E', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '{"label": "ALAMAT PEMBELI - WIYUNG", "phone": "098765433210", "city_name": " SURABAYA", "address_id": "a0794254-012a-4381-8f22-4bab020ed420", "is_default": true, "postal_code": "60228", "address_line": "Jl.xxxx.xxx No.01", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY", "subdistrict_name": "WIYUNG"}', 'shipped', 'paid', 50000.00, 8000.00, 6000.00, 64000.00, '2026-03-12 04:31:46.178518+00', '2026-03-12 04:31:46.178518+00', '2026-03-12 04:42:01.370921+00');
INSERT INTO public.orders (id, order_no, user_id, vendor_id, shipping_address_snapshot, order_status, payment_status, subtotal, shipping_fee, platform_fee, grand_total, placed_at, created_at, updated_at) VALUES ('80432903-81b1-48d6-82ba-670d2798a0d5', 'ORD-20260313-S001', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '{"city": "Jakarta Barat", "name": "Dev Customer", "phone": "081200000002", "address": "Jl. Kebon Jeruk No.10", "province": "DKI Jakarta", "postal_code": "11530"}', 'completed', 'paid', 185000.00, 15000.00, 2500.00, 202500.00, '2026-03-08 07:24:21.226459+00', '2026-03-08 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.order_items (id, order_id, product_variant_id, product_name_snapshot, sku_snapshot, qty, unit_price, line_total) VALUES ('28633b60-c64d-4859-a68d-16aba21ac12a', '382762cd-f805-41a0-ad67-07c6a33020d2', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 'Souvenir Mekkah', 'SKU-SOUV-001', 1, 50000.00, 50000.00);
INSERT INTO public.order_items (id, order_id, product_variant_id, product_name_snapshot, sku_snapshot, qty, unit_price, line_total) VALUES ('6ca279d6-8b8b-4ded-b035-3ea182014a7e', '80432903-81b1-48d6-82ba-670d2798a0d5', 'dcc2c0cd-35f9-4744-a7d1-542e54b2743f', 'Kain Ihram Premium', 'IHRAM-WHT-001', 1, 185000.00, 185000.00);


--
-- Data for Name: order_status_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('35825d8a-438f-4cdb-a0c5-4c910d3f83d5', '382762cd-f805-41a0-ad67-07c6a33020d2', NULL, 'pending_payment', NULL, '2026-03-12 04:31:46.641578+00', NULL);
INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('1e3dace3-5f22-413a-a742-fa2819d4244e', '382762cd-f805-41a0-ad67-07c6a33020d2', 'paid', 'processing', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-03-12 04:41:55.457721+00', 'Order accepted by vendor');
INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('b8e9dab5-12ed-4dfa-a524-0abb9506cf00', '382762cd-f805-41a0-ad67-07c6a33020d2', 'processing', 'shipped', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-03-12 04:42:01.371236+00', 'Order shipped by vendor, tracking: JX5670179152 (jnt EZ)');
INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('b83ae863-a173-4b57-948e-10af6c4d5951', '80432903-81b1-48d6-82ba-670d2798a0d5', 'pending_payment', 'paid', NULL, '2026-03-08 07:24:21.226459+00', 'Payment confirmed via Xendit');
INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('5c3f1d52-74e0-4622-bcab-6822382b0ee6', '80432903-81b1-48d6-82ba-670d2798a0d5', 'paid', 'processing', NULL, '2026-03-09 07:24:21.226459+00', 'Vendor memproses pesanan');
INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('330a18f4-193c-4d9a-a5d5-0dfd792c38fb', '80432903-81b1-48d6-82ba-670d2798a0d5', 'processing', 'packed', NULL, '2026-03-10 07:24:21.226459+00', 'Pesanan dikemas');
INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('664f4cb6-055f-4b44-be0f-bbc76ab188d8', '80432903-81b1-48d6-82ba-670d2798a0d5', 'packed', 'shipped', NULL, '2026-03-11 07:24:21.226459+00', 'Dikirim via JNE REG');
INSERT INTO public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) VALUES ('8625417a-9fc0-409d-b886-6c34e9344539', '80432903-81b1-48d6-82ba-670d2798a0d5', 'shipped', 'completed', NULL, '2026-03-12 07:24:21.226459+00', 'Pesanan selesai');


--
-- Data for Name: otp_codes; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.otp_codes (id, user_id, email, purpose, channel, code_hash, expires_at, attempt_count, consumed_at, invalidated_at, created_at) VALUES ('87229a9b-97d1-47e5-b99b-04d0f10cd25c', 'e166571c-3c00-488b-9844-ec536a54f332', 'vendor-seed@dev.local', 'vendor_onboarding', 'email', 'e10adc3949ba59abbe56e057f20f883e', '2026-02-12 07:24:21.226459+00', 0, '2026-02-11 07:24:21.226459+00', NULL, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: payment_invoices; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.payment_invoices (id, order_id, gateway, xendit_invoice_id, external_invoice_id, invoice_url, payment_method, payment_channel, amount, currency, status, expires_at, paid_at, raw_payload, created_at, updated_at) VALUES ('768ef5ce-15f8-476b-a4de-8d6c9b89538f', '382762cd-f805-41a0-ad67-07c6a33020d2', 'xendit', 'noop-inv-efc97e93', 'INV-ORD-20260312-5F560F3E-13850065', 'https://checkout-bypass.example.com/noop-inv-efc97e93', NULL, NULL, 64000.00, 'IDR', 'paid', '2026-03-13 04:31:46+00', '2026-03-12 04:36:09.603615+00', NULL, '2026-03-12 04:31:46.178518+00', '2026-03-12 04:31:46.178518+00');
INSERT INTO public.payment_invoices (id, order_id, gateway, xendit_invoice_id, external_invoice_id, invoice_url, payment_method, payment_channel, amount, currency, status, expires_at, paid_at, raw_payload, created_at, updated_at) VALUES ('a8a18b80-5aa5-4237-884d-35fd3f92c386', '80432903-81b1-48d6-82ba-670d2798a0d5', 'xendit', 'xnd_seed_001', 'INV-SEED-20260313-001', 'https://checkout-staging.xendit.co/v2/seed-001', 'QRIS', 'QRIS', 202500.00, 'IDR', 'paid', '2026-03-09 07:24:21.226459+00', '2026-03-08 07:34:21.226459+00', '{"id": "xnd_seed_001", "status": "PAID"}', '2026-03-08 07:24:21.226459+00', '2026-03-08 07:24:21.226459+00');


--
-- Data for Name: payment_events; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.payment_events (id, payment_invoice_id, event_type, external_event_id, payload, received_at) VALUES ('a4abe5cc-6431-475e-a139-49449198b14c', 'a8a18b80-5aa5-4237-884d-35fd3f92c386', 'invoice.paid', 'evt_seed_001', '{"id": "xnd_seed_001", "event": "invoice.paid", "status": "PAID"}', '2026-03-08 07:34:21.226459+00');


--
-- Data for Name: payout_batches; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.payout_batches (id, vendor_id, period_start, period_end, status, total_gross, total_fee, total_net, paid_at, created_by, xendit_payout_id, channel_code, description, xendit_status, created_at, updated_at) VALUES ('d7660a6f-1334-447b-bab7-7ab9e9066125', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-02-11', '2026-03-12', 'completed', 185000.00, 2500.00, 182500.00, '2026-03-12 07:24:21.226459+00', '7907d953-ba02-40a4-b403-07b68d8471e8', NULL, NULL, NULL, NULL, '2026-03-11 07:24:21.226459+00', '2026-03-12 07:24:21.226459+00');


--
-- Data for Name: payout_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.payout_items (id, payout_batch_id, order_id, gross_amount, platform_fee_amount, net_amount) VALUES ('1ab40eae-e67d-4454-ba37-1dc14345c377', 'd7660a6f-1334-447b-bab7-7ab9e9066125', '80432903-81b1-48d6-82ba-670d2798a0d5', 185000.00, 2500.00, 182500.00);


--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.product_images (id, product_id, image_url, mime_type, file_size_bytes, is_primary, sort_order, created_at) VALUES ('cafd4659-1e16-467d-a1af-a57165343cbb', 'b2999c9e-df33-48ef-a109-38c77e51f25d', 'https://placehold.co/600x600?text=Kain+Ihram', 'image/png', NULL, true, 0, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_reviews; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.product_reviews (id, product_id, order_id, order_item_id, user_id, rating, review_text, status, created_at, updated_at) VALUES ('f21d4561-2de3-4d9a-812e-2582c69330af', 'b2999c9e-df33-48ef-a109-38c77e51f25d', '80432903-81b1-48d6-82ba-670d2798a0d5', '6ca279d6-8b8b-4ded-b035-3ea182014a7e', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 4, 'Kainnya cukup bagus dan lembut, tapi warna agak berbeda dari foto. Overall puas.', 'published', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_review_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.product_review_images (id, review_id, object_key, mime_type, file_size_bytes, sort_order, created_at) VALUES ('7dc0b2b0-1ffb-4b66-9277-f8e99130f459', 'f21d4561-2de3-4d9a-812e-2582c69330af', 'reviews/seed-review-001.jpg', 'image/jpeg', 153600, 0, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_review_stats; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.product_review_stats (product_id, total_reviews, total_stars, average_rating, star_0_count, star_1_count, star_2_count, star_3_count, star_4_count, star_5_count, updated_at) VALUES ('b2999c9e-df33-48ef-a109-38c77e51f25d', 1, 4, 4.00, 0, 0, 0, 0, 1, 0, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: refunds; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.refunds (id, order_id, payment_invoice_id, amount, reason, status, requested_by, processed_by, processed_at) VALUES ('c4b049f4-91b0-4b22-87ac-6df46f01ee6d', '80432903-81b1-48d6-82ba-670d2798a0d5', 'a8a18b80-5aa5-4237-884d-35fd3f92c386', 50000.00, 'Barang tidak sesuai - refund partial', 'processed', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-12 07:24:21.226459+00');


--
-- Data for Name: reply_templates; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.reply_templates (id, title, content, is_active, created_by, updated_by, created_at, updated_at, shortcut, category) VALUES ('c0e89a64-456f-4df8-bf77-2b8acd4b1b25', 'Salam Pembuka', 'Assalamualaikum, terima kasih telah menghubungi kami. Ada yang bisa kami bantu?', true, 'cc000001-0000-0000-0000-000000000002', NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '/salam', 'general');


--
-- Data for Name: return_reasons; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.return_reasons (id, reason, created_at, updated_at) VALUES (1, 'Barang tidak sesuai deskripsi', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.schema_migrations (version, dirty) VALUES (42, false);


--
-- Data for Name: shipments; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.shipments (id, order_id, courier_code, service_type, tracking_no, shipment_status, shipped_at, delivered_at, etd) VALUES ('04c23c0f-1620-400c-b27a-949a27a5309e', '382762cd-f805-41a0-ad67-07c6a33020d2', 'anteraja', 'ECO', '11002960818873', 'shipped', '2026-03-12 04:42:01.3673+00', NULL, '2-3 day');
INSERT INTO public.shipments (id, order_id, courier_code, service_type, tracking_no, shipment_status, shipped_at, delivered_at, etd) VALUES ('de7cabce-2c88-4a0a-a7dc-e986fa75268e', '80432903-81b1-48d6-82ba-670d2798a0d5', 'jne', 'REG', 'JNESEED0000001', 'delivered', '2026-03-11 07:24:21.226459+00', '2026-03-12 07:24:21.226459+00', '2-3 hari');


--
-- Data for Name: ticket_subjects; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (1, 'Barang Tidak Sampai', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (2, 'Barang Rusak', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (3, 'Pengembalian Dana', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (4, 'Barang Tidak Sesuai', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (5, 'Pengiriman Terlambat', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (6, 'Pembatalan Pesanan', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (7, 'Kesalahan Produk', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (8, 'Akun Bermasalah', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (9, 'Pembayaran Gagal', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (10, 'Voucher/Promo Tidak Berlaku', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (11, 'Pertanyaan Umum', true, '2026-03-12 03:55:18.534203+00');
INSERT INTO public.ticket_subjects (id, label, is_active, created_at) VALUES (12, 'Lainnya', true, '2026-03-12 03:55:18.534203+00');


--
-- Data for Name: tickets; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.tickets (id, ticket_number, customer_id, assigned_cs_id, order_number, phone, reporter_name, subject, detail, status, source, closed_at, created_at, updated_at, attachment_url, attachment_content_type, subject_id) VALUES ('7bbf0cff-c50c-4a81-a116-6381d0fcec70', 'TKT-20260313-0001', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'cc000001-0000-0000-0000-000000000002', 'ORD-20260313-S001', '081200000002', 'Dev Customer', 'Barang Tidak Sesuai', 'Barang yang diterima tidak sesuai dengan deskripsi produk. Warna dan ukuran berbeda dari yang dipesan.', 'closed', 'web', '2026-03-13 07:24:21.226459+00', '2026-03-12 19:24:21.226459+00', '2026-03-13 07:24:21.226459+00', NULL, NULL, 4);


--
-- Data for Name: ticket_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ticket_messages (id, ticket_id, sender_id, message, is_from_cs, is_internal_note, created_at) VALUES ('406b9695-60be-4313-a6e8-c7e3cbb515eb', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Barang yang saya terima tidak sesuai dengan deskripsi. Mohon bantuan pengembaliannya.', false, false, '2026-03-12 19:24:21.226459+00');


--
-- Data for Name: ticket_attachments; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ticket_attachments (id, ticket_id, message_id, file_url, file_name, file_type, file_size, created_at) VALUES ('0f113c48-113c-4d01-8ce9-f7bb143fb24d', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', '406b9695-60be-4313-a6e8-c7e3cbb515eb', 'https://placehold.co/800x600?text=Bukti+Foto', 'bukti-foto.jpg', 'image/jpeg', 204800, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: ticket_status_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.ticket_status_logs (id, ticket_id, changed_by, old_status, new_status, notes, created_at) VALUES ('8ded3edf-4fd9-457c-8aee-aefc0f782ae0', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', 'cc000001-0000-0000-0000-000000000002', 'open', 'on_progress', 'CS mengambil tiket', '2026-03-13 01:24:21.226459+00');
INSERT INTO public.ticket_status_logs (id, ticket_id, changed_by, old_status, new_status, notes, created_at) VALUES ('4f40cca2-b3d2-4c03-a060-69641718f153', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', 'cc000001-0000-0000-0000-000000000002', 'on_progress', 'closed', 'Masalah telah diselesaikan', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_balances; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_balances (vendor_id, available_balance, pending_balance, total_earned, total_withdrawn, updated_at, escrow_balance) VALUES ('9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 1000000.00, 0.00, 1000000.00, 0.00, '2026-03-12 04:17:46.056983+00', 0.00);
INSERT INTO public.vendor_balances (vendor_id, available_balance, pending_balance, total_earned, total_withdrawn, updated_at, escrow_balance) VALUES ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 500000.00, 0.00, 1500000.00, 1000000.00, '2026-03-13 07:24:21.226459+00', 0.00);


--
-- Data for Name: vendor_bank_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_bank_accounts (id, vendor_id, bank_name, account_number, account_holder_name, verification_status, rejection_reason, verified_by, verified_at, created_at, updated_at) VALUES ('7b2aee15-39da-4080-bf21-df2feb5c791e', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'BANK Dummy', '1234567890', 'Nursufyan Sauri', 'verified', NULL, '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-12 04:17:26.124526+00', '2026-03-12 04:17:26.124526+00', '2026-03-12 04:17:26.124526+00');
INSERT INTO public.vendor_bank_accounts (id, vendor_id, bank_name, account_number, account_holder_name, verification_status, rejection_reason, verified_by, verified_at, created_at, updated_at) VALUES ('841c5c84-cc53-4196-97f0-bd919bea6179', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'Bank Syariah Indonesia', '7788990011', 'PT Oleh-Oleh Haji Berkah', 'verified', NULL, '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-02-16 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_banners; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_banners (id, vendor_id, title, image_url, created_at, updated_at) VALUES ('dd107057-bd24-4da5-860a-6de0b6d7cdb5', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'Promo Ramadhan', 'https://placehold.co/1200x400?text=Promo+Ramadhan', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_couriers; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('04737655-d896-4a11-aa1d-567769010fcb', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 3, true, '2026-03-12 04:12:25.593572+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('609cc811-d7c1-4ab9-a8d7-c45fb66a03da', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 6, true, '2026-03-12 04:12:25.594108+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('ee3d321e-42f7-4b6a-b99b-fb52fef0193c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 1, true, '2026-03-12 04:12:25.594364+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('83f69d0e-27bf-43d5-895c-597814d0f7ec', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 11, true, '2026-03-12 04:12:25.594595+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('f56186f5-f341-45e3-81ca-1a7b471af801', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 5, true, '2026-03-12 04:12:25.594792+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('bfe88d97-d13b-4ae5-9367-aad749d0634c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 9, true, '2026-03-12 04:12:25.59502+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('5203e452-d7a8-4f60-8983-ee01d23b33b9', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 12, true, '2026-03-12 04:12:25.595229+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('813352d8-5fec-4e1d-aa99-a73bc4fb1eb5', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 4, true, '2026-03-12 04:12:25.595467+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('8cf8a32c-f3e8-492f-9b0f-dd9889fc8ed5', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 10, true, '2026-03-12 04:12:25.595702+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('7fe72ce5-4807-4dee-b9d2-f2892c8cbe51', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 2, true, '2026-03-12 04:12:25.595913+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('91ce5cf3-18ba-4fb0-bac8-fcbc9fe471ec', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 7, true, '2026-03-12 04:12:25.596118+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('f371e602-a0a2-429a-b962-256df2ee211c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 8, true, '2026-03-12 04:12:25.596361+00');
INSERT INTO public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) VALUES ('cdc48954-fe6d-4680-997f-0914bbdc4874', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 1, true, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_documents; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_documents (id, vendor_id, doc_type, file_url, mime_type, file_size_bytes, file_checksum, uploaded_by, verification_status, rejection_reason, verified_by, verified_at, created_at, updated_at) VALUES ('a0726b95-ee4e-4ac3-b3e8-1d6ceeb6564b', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'owner_document_id', 'https://placehold.co/600x400?text=KTP', 'image/jpeg', NULL, NULL, 'e166571c-3c00-488b-9844-ec536a54f332', 'verified', NULL, '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-02-16 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_onboardings; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_onboardings (id, email, status, password_hash, store_name, vendor_type, business_legal_type, document_id_type, nik, owner_name, birth_date, document_id_object_key, otp_verified_at, completed_at, created_at, updated_at) VALUES ('fad8946b-1911-4917-bb7c-788a37a66954', 'vendor-seed@dev.local', 'completed', '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890ABCDEF0123456', 'Toko Berkah Haji', 'souvenir_store', 'perorangan', 'ktp', '3201010101900001', 'Vendor Seed Owner', '1990-01-01', NULL, '2026-02-10 07:24:21.226459+00', '2026-02-11 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_withdrawals; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.vendor_withdrawals (id, vendor_id, amount, channel_code, status, xendit_payout_id, xendit_status, description, failed_reason, created_at, updated_at, fee_estimated, fee_actual, amount_net, total_deducted, fee_status, xendit_transaction_id) VALUES ('90900fc3-880d-49e3-98ce-534da7bfa040', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 500000.00, 'ID_BCA', 'completed', NULL, NULL, 'Penarikan saldo vendor', NULL, '2026-03-10 07:24:21.226459+00', '2026-03-10 07:24:21.226459+00', 2500.00, NULL, NULL, 502500.00, 'resolved', NULL);


--
-- Data for Name: wishlist_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

INSERT INTO public.wishlist_items (id, user_id, product_id, created_at) VALUES ('f105e107-7125-4d58-9125-9286c3543080', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'b2999c9e-df33-48ef-a109-38c77e51f25d', '2026-03-13 07:24:21.226459+00');


--
-- Name: admin_contacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.admin_contacts_id_seq', 1, true);


--
-- Name: couriers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.couriers_id_seq', 12, true);


--
-- Name: faqs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.faqs_id_seq', 1, true);


--
-- Name: return_reasons_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.return_reasons_id_seq', 1, true);


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

\unrestrict uwmddg8yAn0SoGa7WelfBPBZoSrKDsVXNxdSEnzzUItuEqOEAlQJcdtA0JzD0Hg

