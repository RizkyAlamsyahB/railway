--
-- PostgreSQL database dump
--

\restrict DB8L0OWFBExbXjJkiFzNg8z38hhHubmGCEwuaIY4bWU4a8WsW8Kw6qhvdOlWAM5

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
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.roles VALUES (1, 'admin', 'Administrator');
INSERT INTO public.roles VALUES (2, 'umkm', 'UMKM');
INSERT INTO public.roles VALUES (3, 'customer', 'Customer');
INSERT INTO public.roles VALUES (4, 'cs', 'Customer Service');
INSERT INTO public.roles VALUES (5, 'finance', 'Finance');


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.users VALUES ('aa000002-0000-0000-0000-000000000001', 'umkm2@dev.local', 'Toko Oleh-Oleh Tanah Suci', NULL, '081299000002', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('aa000002-0000-0000-0000-000000000002', 'umkm3@dev.local', 'Perlengkapan Ibadah Barokah', NULL, '081299000003', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('aa000002-0000-0000-0000-000000000003', 'customer2@dev.local', 'Siti Aisyah', NULL, '081299000004', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 3, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('7907d953-ba02-40a4-b403-07b68d8471e8', 'admin@example.com', 'Admin', NULL, '08120000000', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 1, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('cc000001-0000-0000-0000-000000000003', 'cs3@dev.local', 'Dev Customer Service 3', NULL, '081200000006', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'ah.nursufyantsauri@gmail.com', 'Dev UMKM', NULL, '081200000001', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 2, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'rizkyalamsyah.dev@gmail.com', 'Dev Customer', NULL, '081200000002', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 3, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('cc000001-0000-0000-0000-000000000001', 'adkhawildanrizqia@gmail.com', 'Dev Customer Service 1', NULL, '081200000003', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('ff000001-0000-0000-0000-000000000001', 'galangarsandy@gmail.com', 'Dev Finance', NULL, '081200000004', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 5, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('e166571c-3c00-488b-9844-ec536a54f332', 'vendor-seed@dev.local', 'Vendor Seed Owner', NULL, '081300000099', '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890ABCDEF0123456', 2, 'active', '2026-02-11 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', NULL);
INSERT INTO public.users VALUES ('cc000001-0000-0000-0000-000000000002', 'harundarat@gmail.com', 'Dev Customer Service 2', NULL, '081200000005', '$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq', 4, 'active', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', '2026-03-06 01:58:19.265177+00', NULL);
INSERT INTO public.users VALUES ('ae02bdd5-246c-4a1c-a619-f7eefbbe4db6', '21082010249@student.upnjatim.ac.id', 'Ahmad Subarkah', '1990-01-02', NULL, '$2a$10$gbPSerb5KETGyJCR/BGXK.Fg6L1tmpiCop.TMufeORmKZDooiFnV6', 2, 'active', '2026-03-16 04:13:27.285894+00', '2026-03-16 04:14:11.886223+00', '2026-03-16 04:14:11.886223+00', NULL);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.addresses VALUES ('e081dc95-612f-4d08-ab35-21fe7b7b998a', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'TOKO WAREHOUSE SEDATI - SIDOARJO', 'TOKO WAREHOUSE SEDATI - SIDOARJO', '062916302649', '18', '583', '6001', '61253', 'WAREHOUSE SEDATI', true, 'JAWA TIMUR', ' SIDOARJO', 'SEDATI', '70995', 'SEDATI GEDE', NULL, NULL, NULL, '2026-03-12 04:11:45.061823+00', '2026-03-12 04:12:09.361232+00');
INSERT INTO public.addresses VALUES ('399e90f7-0e8c-41f8-8032-1eda354f36d3', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Rumah', 'Dev Customer', '081200000002', '18', '577', '5901', '60228', 'Jl. Kebon Jeruk No.10, RT 005/RW 003', true, 'JAWA TIMUR', ' SURABAYA', 'WIYUNG', '69354', 'WIYUNG', NULL, NULL, NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');
INSERT INTO public.addresses VALUES ('a0794254-012a-4381-8f22-4bab020ed420', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'ALAMAT PEMBELI - WIYUNG', 'RIZKY', '098765433210', '18', '577', '5901', '60228', 'Jl.Pagesangangan No.01', true, 'JAWA TIMUR', ' SURABAYA', 'WIYUNG', '69354', 'WIYUNG', NULL, NULL, NULL, '2026-03-12 04:31:09.785581+00', '2026-03-12 04:31:09.785581+00');


--
-- Data for Name: admin_contacts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.admin_contacts VALUES (1, 'Hubungi admin di admin@hajjstore.id atau WhatsApp 08123456789', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: banners; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.banners VALUES ('4190f3f4-45aa-482b-9121-1836b127304c', 'Selamat Datang di Hajj Store', 'https://placehold.co/1200x400?text=Hajj+Store+Banner', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.carts VALUES ('38f3e138-58d9-4941-b948-ff7c252a8b0b', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-12 04:29:36.202609+00', '2026-03-12 04:31:46.662441+00');
INSERT INTO public.carts VALUES ('2a64b014-85a2-4c73-8f07-97cee7a652d1', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');
INSERT INTO public.carts VALUES ('b81c7b33-534e-4b8e-a50d-e7758fdfb3ab', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-12 04:54:53.927372+00', '2026-03-26 04:48:50.995336+00');
INSERT INTO public.carts VALUES ('e30314df-56be-4c00-9a52-b84a1b1c1391', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-26 05:03:46.078079+00', '2026-03-26 05:03:59.801365+00');
INSERT INTO public.carts VALUES ('15da7cc2-086b-4046-b1cd-709c36fabfea', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'converted', '2026-03-26 05:08:47.985175+00', '2026-03-26 05:09:06.583623+00');


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: -
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
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendors VALUES ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'e166571c-3c00-488b-9844-ec536a54f332', 'souvenir_store', 'PT Oleh-Oleh Haji Berkah', 'Toko Berkah Haji', 'Vendor Seed Owner', 'Toko perlengkapan haji dan oleh-oleh berkualitas.', 'active', NULL, NULL, NULL, '', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', NULL);
INSERT INTO public.vendors VALUES ('ef0ad9ea-13b6-46ac-9bac-7a81af8ee705', 'ae02bdd5-246c-4a1c-a619-f7eefbbe4db6', 'souvenir_store', NULL, 'Toko Oleh Oleh Haji Test', 'Ahmad Subarkah', NULL, 'draft', NULL, NULL, NULL, NULL, '2026-03-16 04:14:11.886223+00', '2026-03-16 04:14:11.886223+00', NULL);
INSERT INTO public.vendors VALUES ('9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'souvenir_store', 'PT Nursufyan Sauri', 'TOKO Nursufyan Sauri', 'Nursufyan Sauri', 'Toko milik Nursufyan Sauri', 'active', NULL, NULL, NULL, 'dummy', '2026-03-12 04:05:34.329709+00', '2026-03-12 04:05:34.329709+00', NULL);


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.products VALUES ('0f0ce953-d4a0-4dcb-a621-a53656b69092', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000001', 'Paket Haji Premium', 'paket-haji-premium', 'Paket perlengkapan haji lengkap dan premium.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products VALUES ('2aa62691-ea2c-491e-8456-031984b79686', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000002', 'Souvenir Mekkah', 'souvenir-mekkah', 'Souvenir khas Mekkah untuk oleh-oleh.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products VALUES ('72f634dd-9f2e-429e-9e81-b5703f8b06f8', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000003', 'Baju Muslim Pria', 'baju-muslim-pria', 'Baju muslim pria modern dan nyaman.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products VALUES ('8bb433f1-55c1-40ce-a50b-158150e5cab9', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'ca000001-0000-0000-0000-000000000011', 'Pakaian Ihram Dewasa', 'pakaian-ihram-dewasa', 'Pakaian ihram dewasa bahan premium.', 'published', 'passed', 'Produk halal', '2026-03-12 04:24:27.093679+00', '2026-03-12 04:24:27.093679+00');
INSERT INTO public.products VALUES ('b2999c9e-df33-48ef-a109-38c77e51f25d', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'ca000001-0000-0000-0000-000000000001', 'Kain Ihram Premium', 'kain-ihram-premium-seed', 'Kain ihram premium berbahan katun lembut, nyaman dipakai saat ibadah haji dan umroh.', 'published', 'passed', NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_variants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_variants VALUES ('14218407-1b19-40bd-9136-d611d927e573', '72f634dd-9f2e-429e-9e81-b5703f8b06f8', 'SKU-BAJU-001', 'Default', 250000.00, 'IDR', 20, 400, true, true);
INSERT INTO public.product_variants VALUES ('5f65a8ed-3c03-4cae-822f-4c97980db7f8', '8bb433f1-55c1-40ce-a50b-158150e5cab9', 'SKU-IHRAM-001', 'Default', 350000.00, 'IDR', 15, 800, true, true);
INSERT INTO public.product_variants VALUES ('dcc2c0cd-35f9-4744-a7d1-542e54b2743f', 'b2999c9e-df33-48ef-a109-38c77e51f25d', 'IHRAM-WHT-001', 'Kain Ihram Premium', 185000.00, 'IDR', 50, 500, true, true);
INSERT INTO public.product_variants VALUES ('747b1fc2-edaf-424b-9358-92526ac05a6e', '0f0ce953-d4a0-4dcb-a621-a53656b69092', 'SKU-HAJI-001', 'Default', 1500000.00, 'IDR', 8, 1200, true, true);
INSERT INTO public.product_variants VALUES ('c9f199bd-ea95-442b-a16c-c9fd0769a32f', '2aa62691-ea2c-491e-8456-031984b79686', 'SKU-SOUV-001', 'Default', 50000.00, 'IDR', 48, 200, true, true);


--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.cart_items VALUES ('208e50c9-2003-4eb5-94d1-2b423d5e64dd', '38f3e138-58d9-4941-b948-ff7c252a8b0b', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 1, '2026-03-12 04:29:36.216061+00');
INSERT INTO public.cart_items VALUES ('c28b28c8-e1e1-45dc-aead-d42be0ee5ecf', '2a64b014-85a2-4c73-8f07-97cee7a652d1', 'dcc2c0cd-35f9-4744-a7d1-542e54b2743f', 1, '2026-03-13 07:24:21.226459+00');
INSERT INTO public.cart_items VALUES ('4c28fa3d-b8c7-43fe-a6ce-48ce748e1d06', 'b81c7b33-534e-4b8e-a50d-e7758fdfb3ab', '747b1fc2-edaf-424b-9358-92526ac05a6e', 1, '2026-03-26 04:44:38.813053+00');
INSERT INTO public.cart_items VALUES ('b098095a-223a-4254-91ad-bf4d68e4bf5c', 'e30314df-56be-4c00-9a52-b84a1b1c1391', '747b1fc2-edaf-424b-9358-92526ac05a6e', 1, '2026-03-26 05:03:46.081755+00');
INSERT INTO public.cart_items VALUES ('b81c30eb-8837-49e7-96e2-7a68c4799f77', '15da7cc2-086b-4046-b1cd-709c36fabfea', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 1, '2026-03-26 05:08:47.996751+00');


--
-- Data for Name: chat_conversations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.chat_conversations VALUES ('45ec9ef1-bd88-4dd0-8208-3b5c0e5e6eef', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'cc000001-0000-0000-0000-000000000002', '2026-03-13 06:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.chat_messages VALUES ('ae8f5ec0-04ed-4308-a818-b600da71ec8e', '45ec9ef1-bd88-4dd0-8208-3b5c0e5e6eef', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Halo, saya mau tanya soal pesanan saya.', NULL, true, '2026-03-13 06:34:21.226459+00', '2026-03-13 06:24:21.226459+00');


--
-- Data for Name: couriers; Type: TABLE DATA; Schema: public; Owner: -
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
-- Data for Name: email_verification_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.email_verification_tokens VALUES ('779944a5-b096-4417-a389-8e28921bbd47', 'e166571c-3c00-488b-9844-ec536a54f332', 'vendor-seed@dev.local', 'seed_token_hash_0c538b83df4bdf8e4c3cd73df13b09b5', '2026-02-12 07:24:21.226459+00', '2026-02-11 07:24:21.226459+00', NULL, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: faqs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.faqs VALUES (1, 'Umum', 'Bagaimana cara memesan perlengkapan haji?', 'Anda bisa memesan melalui aplikasi atau website kami. Pilih produk, masukkan ke keranjang, lalu lakukan checkout.', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: ledger_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ledger_accounts VALUES ('60825dd3-e090-458e-a67b-9393b9834b10', '1100', 'Payment Gateway Receivable', 'asset', 'D', true);
INSERT INTO public.ledger_accounts VALUES ('8e9d7291-ced8-4ca5-86ae-c4739e0fefa4', '2100', 'Vendor Payable', 'liability', 'C', true);
INSERT INTO public.ledger_accounts VALUES ('ccf126d7-36dd-4d53-845f-a35da76acaa7', '4100', 'Platform Fee Revenue', 'revenue', 'C', true);
INSERT INTO public.ledger_accounts VALUES ('240bbd21-9b1f-42ac-b975-9ddff62ead79', '4200', 'Admin Fee Revenue', 'revenue', 'C', true);


--
-- Data for Name: ledger_journals; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ledger_journals VALUES ('f9f10cdd-475e-4ab6-9b63-8b5be5f4bae8', 'JRN-SEED-001', 'payment_invoice', 'a8a18b80-5aa5-4237-884d-35fd3f92c386', '2026-03-08 07:24:21.226459+00', 'Payment received for order ORD-20260313-S001', 'posted', '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-13 07:24:21.226459+00');
INSERT INTO public.ledger_journals VALUES ('14802771-711f-425e-a8eb-e4e6ab176c24', 'JRN-PAY-ORD-20260326-2A4E66F5-d0e1da02', 'payment_invoice', 'df7bf23e-4db0-4f1e-8ea0-23809ad0ce78', '2026-03-26 05:04:02.318631+00', 'Payment received for order ORD-20260326-2A4E66F5', 'posted', NULL, '2026-03-26 05:04:02.318631+00');
INSERT INTO public.ledger_journals VALUES ('cc39d886-3674-4cb3-b318-4df834ee2577', 'JRN-PAY-ORD-20260326-017AFEBB-16928071', 'payment_invoice', '950bcbb9-11af-48c8-9a20-2a3289e416fc', '2026-03-26 05:09:10.28312+00', 'Payment received for order ORD-20260326-017AFEBB', 'posted', NULL, '2026-03-26 05:09:10.28312+00');


--
-- Data for Name: ledger_lines; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ledger_lines VALUES ('3a008fb2-5a0b-4134-8f14-cba10a1fcbec', 'f9f10cdd-475e-4ab6-9b63-8b5be5f4bae8', '60825dd3-e090-458e-a67b-9393b9834b10', 202500.00, 0.00, 'IDR', 'ORD-20260313-S001');
INSERT INTO public.ledger_lines VALUES ('46fb3921-36c5-4247-bee8-51677fe5198f', 'f9f10cdd-475e-4ab6-9b63-8b5be5f4bae8', '8e9d7291-ced8-4ca5-86ae-c4739e0fefa4', 0.00, 202500.00, 'IDR', 'ORD-20260313-S001');
INSERT INTO public.ledger_lines VALUES ('0753a766-96db-45ad-90de-44025d71f722', '14802771-711f-425e-a8eb-e4e6ab176c24', '60825dd3-e090-458e-a67b-9393b9834b10', 1536000.00, 0.00, 'IDR', 'ORD-20260326-2A4E66F5');
INSERT INTO public.ledger_lines VALUES ('72887479-ba1b-44a1-9b72-5aa657895155', '14802771-711f-425e-a8eb-e4e6ab176c24', '8e9d7291-ced8-4ca5-86ae-c4739e0fefa4', 0.00, 1500000.00, 'IDR', 'ORD-20260326-2A4E66F5');
INSERT INTO public.ledger_lines VALUES ('7d521c63-73b0-4196-82cc-712824eae501', '14802771-711f-425e-a8eb-e4e6ab176c24', 'ccf126d7-36dd-4d53-845f-a35da76acaa7', 0.00, 1000.00, 'IDR', 'ORD-20260326-2A4E66F5');
INSERT INTO public.ledger_lines VALUES ('c62e094e-9a1c-4fd7-8671-43b427c2fa48', '14802771-711f-425e-a8eb-e4e6ab176c24', '240bbd21-9b1f-42ac-b975-9ddff62ead79', 0.00, 5000.00, 'IDR', 'ORD-20260326-2A4E66F5');
INSERT INTO public.ledger_lines VALUES ('46648d07-931f-4b69-bbc4-6be9664fcf41', 'cc39d886-3674-4cb3-b318-4df834ee2577', '60825dd3-e090-458e-a67b-9393b9834b10', 61000.00, 0.00, 'IDR', 'ORD-20260326-017AFEBB');
INSERT INTO public.ledger_lines VALUES ('3e589cf8-660d-4ef6-b90a-ba0d4e03fd92', 'cc39d886-3674-4cb3-b318-4df834ee2577', '8e9d7291-ced8-4ca5-86ae-c4739e0fefa4', 0.00, 50000.00, 'IDR', 'ORD-20260326-017AFEBB');
INSERT INTO public.ledger_lines VALUES ('50353f1a-3d79-4718-bce0-a651175fd901', 'cc39d886-3674-4cb3-b318-4df834ee2577', 'ccf126d7-36dd-4d53-845f-a35da76acaa7', 0.00, 1000.00, 'IDR', 'ORD-20260326-017AFEBB');
INSERT INTO public.ledger_lines VALUES ('ae3347ab-2e9e-4ba6-b862-d7645614a458', 'cc39d886-3674-4cb3-b318-4df834ee2577', '240bbd21-9b1f-42ac-b975-9ddff62ead79', 0.00, 5000.00, 'IDR', 'ORD-20260326-017AFEBB');


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.notifications VALUES ('f2bf6aa4-b496-442c-8fe8-2bcdc3e0d765', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'order', 'Pesanan Selesai', 'Pesanan ORD-20260313-S001 telah selesai. Terima kasih telah berbelanja!', '2026-03-13 07:24:21.226459+00', NULL);


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.orders VALUES ('382762cd-f805-41a0-ad67-07c6a33020d2', 'ORD-20260312-5F560F3E', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '{"label": "ALAMAT PEMBELI - WIYUNG", "phone": "098765433210", "city_name": " SURABAYA", "address_id": "a0794254-012a-4381-8f22-4bab020ed420", "is_default": true, "postal_code": "60228", "address_line": "Jl.xxxx.xxx No.01", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY", "subdistrict_name": "WIYUNG"}', 'shipped', 'paid', 50000.00, 8000.00, 6000.00, 64000.00, '2026-03-12 04:31:46.178518+00', '2026-03-12 04:31:46.178518+00', '2026-03-12 04:42:01.370921+00');
INSERT INTO public.orders VALUES ('80432903-81b1-48d6-82ba-670d2798a0d5', 'ORD-20260313-S001', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '{"city": "Jakarta Barat", "name": "Dev Customer", "phone": "081200000002", "address": "Jl. Kebon Jeruk No.10", "province": "DKI Jakarta", "postal_code": "11530"}', 'completed', 'paid', 185000.00, 15000.00, 2500.00, 202500.00, '2026-03-08 07:24:21.226459+00', '2026-03-08 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');
INSERT INTO public.orders VALUES ('a5c313aa-3464-48c3-baf6-ae6fd5befd62', 'ORD-20260326-03936280', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '{"label": "Rumah", "phone": "081200000002", "city_name": " SURABAYA", "address_id": "399e90f7-0e8c-41f8-8032-1eda354f36d3", "is_default": true, "postal_code": "60228", "address_line": "Jl. Kebon Jeruk No.10, RT 005/RW 003", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "Dev Customer", "subdistrict_name": "WIYUNG"}', 'pending_payment', 'unpaid', 1500000.00, 14000.00, 6000.00, 1520000.00, '2026-03-26 04:48:50.540824+00', '2026-03-26 04:48:50.540824+00', '2026-03-26 04:48:50.540824+00');
INSERT INTO public.orders VALUES ('99418086-e674-41ec-8dcc-6c4301202502', 'ORD-20260326-2A4E66F5', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '{"label": "Rumah", "phone": "081200000002", "city_name": " SURABAYA", "address_id": "399e90f7-0e8c-41f8-8032-1eda354f36d3", "is_default": true, "postal_code": "60228", "address_line": "Jl. Kebon Jeruk No.10, RT 005/RW 003", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "Dev Customer", "subdistrict_name": "WIYUNG"}', 'shipped', 'paid', 1500000.00, 30000.00, 6000.00, 1536000.00, '2026-03-26 05:03:59.310154+00', '2026-03-26 05:03:59.310154+00', '2026-03-26 05:07:14.128213+00');
INSERT INTO public.orders VALUES ('1e88413e-31ae-4698-b321-a6862044d31c', 'ORD-20260326-017AFEBB', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', '{"label": "Rumah", "phone": "081200000002", "city_name": " SURABAYA", "address_id": "399e90f7-0e8c-41f8-8032-1eda354f36d3", "is_default": true, "postal_code": "60228", "address_line": "Jl. Kebon Jeruk No.10, RT 005/RW 003", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "Dev Customer", "subdistrict_name": "WIYUNG"}', 'paid', 'paid', 50000.00, 5000.00, 6000.00, 61000.00, '2026-03-26 05:09:06.094561+00', '2026-03-26 05:09:06.094561+00', '2026-03-26 05:09:10.279433+00');


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.order_items VALUES ('28633b60-c64d-4859-a68d-16aba21ac12a', '382762cd-f805-41a0-ad67-07c6a33020d2', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 'Souvenir Mekkah', 'SKU-SOUV-001', 1, 50000.00, 50000.00);
INSERT INTO public.order_items VALUES ('6ca279d6-8b8b-4ded-b035-3ea182014a7e', '80432903-81b1-48d6-82ba-670d2798a0d5', 'dcc2c0cd-35f9-4744-a7d1-542e54b2743f', 'Kain Ihram Premium', 'IHRAM-WHT-001', 1, 185000.00, 185000.00);
INSERT INTO public.order_items VALUES ('2ceb6681-aea9-4d2a-b2ba-2cf6ef017785', 'a5c313aa-3464-48c3-baf6-ae6fd5befd62', '747b1fc2-edaf-424b-9358-92526ac05a6e', 'Paket Haji Premium', 'SKU-HAJI-001', 1, 1500000.00, 1500000.00);
INSERT INTO public.order_items VALUES ('160c07bb-ae7c-49e1-8121-8b2af10b3526', '99418086-e674-41ec-8dcc-6c4301202502', '747b1fc2-edaf-424b-9358-92526ac05a6e', 'Paket Haji Premium', 'SKU-HAJI-001', 1, 1500000.00, 1500000.00);
INSERT INTO public.order_items VALUES ('ed506dfc-97f5-4969-b8ce-096d08ab352f', '1e88413e-31ae-4698-b321-a6862044d31c', 'c9f199bd-ea95-442b-a16c-c9fd0769a32f', 'Souvenir Mekkah', 'SKU-SOUV-001', 1, 50000.00, 50000.00);


--
-- Data for Name: order_status_history; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.order_status_history VALUES ('35825d8a-438f-4cdb-a0c5-4c910d3f83d5', '382762cd-f805-41a0-ad67-07c6a33020d2', NULL, 'pending_payment', NULL, '2026-03-12 04:31:46.641578+00', NULL);
INSERT INTO public.order_status_history VALUES ('1e3dace3-5f22-413a-a742-fa2819d4244e', '382762cd-f805-41a0-ad67-07c6a33020d2', 'paid', 'processing', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-03-12 04:41:55.457721+00', 'Order accepted by vendor');
INSERT INTO public.order_status_history VALUES ('b8e9dab5-12ed-4dfa-a524-0abb9506cf00', '382762cd-f805-41a0-ad67-07c6a33020d2', 'processing', 'shipped', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-03-12 04:42:01.371236+00', 'Order shipped by vendor, tracking: JX5670179152 (jnt EZ)');
INSERT INTO public.order_status_history VALUES ('b83ae863-a173-4b57-948e-10af6c4d5951', '80432903-81b1-48d6-82ba-670d2798a0d5', 'pending_payment', 'paid', NULL, '2026-03-08 07:24:21.226459+00', 'Payment confirmed via Xendit');
INSERT INTO public.order_status_history VALUES ('5c3f1d52-74e0-4622-bcab-6822382b0ee6', '80432903-81b1-48d6-82ba-670d2798a0d5', 'paid', 'processing', NULL, '2026-03-09 07:24:21.226459+00', 'Vendor memproses pesanan');
INSERT INTO public.order_status_history VALUES ('330a18f4-193c-4d9a-a5d5-0dfd792c38fb', '80432903-81b1-48d6-82ba-670d2798a0d5', 'processing', 'packed', NULL, '2026-03-10 07:24:21.226459+00', 'Pesanan dikemas');
INSERT INTO public.order_status_history VALUES ('664f4cb6-055f-4b44-be0f-bbc76ab188d8', '80432903-81b1-48d6-82ba-670d2798a0d5', 'packed', 'shipped', NULL, '2026-03-11 07:24:21.226459+00', 'Dikirim via JNE REG');
INSERT INTO public.order_status_history VALUES ('8625417a-9fc0-409d-b886-6c34e9344539', '80432903-81b1-48d6-82ba-670d2798a0d5', 'shipped', 'completed', NULL, '2026-03-12 07:24:21.226459+00', 'Pesanan selesai');
INSERT INTO public.order_status_history VALUES ('823231ad-4a75-4e41-aa8f-21414f36710c', 'a5c313aa-3464-48c3-baf6-ae6fd5befd62', NULL, 'pending_payment', NULL, '2026-03-26 04:48:50.981964+00', NULL);
INSERT INTO public.order_status_history VALUES ('11b94751-51eb-4ea0-adfb-a7ee178b8541', '99418086-e674-41ec-8dcc-6c4301202502', NULL, 'pending_payment', NULL, '2026-03-26 05:03:59.790135+00', NULL);
INSERT INTO public.order_status_history VALUES ('a071e80b-94ce-4f47-96ec-f02c4fc33762', '99418086-e674-41ec-8dcc-6c4301202502', 'pending_payment', 'paid', NULL, '2026-03-26 05:04:02.31439+00', 'Payment received via SIMULATE/BYPASS');
INSERT INTO public.order_status_history VALUES ('cc64c41c-a92a-4ab8-8d23-8fcb4bf33b01', '99418086-e674-41ec-8dcc-6c4301202502', 'paid', 'processing', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-03-26 05:04:55.950359+00', 'Order accepted by vendor');
INSERT INTO public.order_status_history VALUES ('e670e1a4-aa87-48a1-94e4-6b900cce6d2e', '99418086-e674-41ec-8dcc-6c4301202502', 'processing', 'shipped', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-03-26 05:07:14.128639+00', 'Order shipped by vendor, tracking: 11002960818873 (anteraja ECO)');
INSERT INTO public.order_status_history VALUES ('c4037c26-eb55-4ca8-93ee-8eaab75fadbe', '1e88413e-31ae-4698-b321-a6862044d31c', NULL, 'pending_payment', NULL, '2026-03-26 05:09:06.569078+00', NULL);
INSERT INTO public.order_status_history VALUES ('2d06900d-c40f-492a-aa06-68690853ed18', '1e88413e-31ae-4698-b321-a6862044d31c', 'pending_payment', 'paid', NULL, '2026-03-26 05:09:10.279971+00', 'Payment received via SIMULATE/BYPASS');


--
-- Data for Name: otp_codes; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.otp_codes VALUES ('87229a9b-97d1-47e5-b99b-04d0f10cd25c', 'e166571c-3c00-488b-9844-ec536a54f332', 'vendor-seed@dev.local', 'vendor_onboarding', 'email', 'e10adc3949ba59abbe56e057f20f883e', '2026-02-12 07:24:21.226459+00', 0, '2026-02-11 07:24:21.226459+00', NULL, '2026-03-13 07:24:21.226459+00');
INSERT INTO public.otp_codes VALUES ('dbae13bc-b218-4cf8-a7d1-bcf0366ca83f', NULL, '21082010249@student.upnjatim.ac.id', 'email_verification', 'email', 'NFkfSRQ3inFxjQ610kZQP7VyCQs6LsRh0-Bn8J2jGv0', '2026-03-16 03:46:23.603302+00', 0, '2026-03-16 03:41:49.875082+00', NULL, '2026-03-16 03:41:23.603302+00');
INSERT INTO public.otp_codes VALUES ('9a6cc236-66ca-4247-8e4f-d0b463b6d2fe', NULL, '21082010249@student.upnjatim.ac.id', 'email_verification', 'email', 'b8AGgNAuQulGRC8styzGqeSkcq1YedWjPaX67MjiJDo', '2026-03-16 04:17:28.5027+00', 0, '2026-03-16 04:13:27.273311+00', NULL, '2026-03-16 04:12:28.5027+00');


--
-- Data for Name: payment_invoices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.payment_invoices VALUES ('a8a18b80-5aa5-4237-884d-35fd3f92c386', '80432903-81b1-48d6-82ba-670d2798a0d5', 'xendit', 'xnd_seed_001', 'INV-SEED-20260313-001', 'https://checkout-staging.xendit.co/v2/seed-001', 'QRIS', 'QRIS', 202500.00, 'IDR', 'paid', '2026-03-09 07:24:21.226459+00', '2026-03-08 07:34:21.226459+00', '{"id": "xnd_seed_001", "status": "PAID"}', '2026-03-08 07:24:21.226459+00', '2026-03-08 07:24:21.226459+00');
INSERT INTO public.payment_invoices VALUES ('768ef5ce-15f8-476b-a4de-8d6c9b89538f', '382762cd-f805-41a0-ad67-07c6a33020d2', 'xendit', 'noop-inv-efc97e93', 'INV-ORD-20260312-5F560F3E-13850065', 'https://checkout-bypass.example.com/noop-inv-efc97e93', 'QRIS', 'QRIS', 64000.00, 'IDR', 'paid', '2026-03-13 04:31:46+00', '2026-03-12 04:36:09.603615+00', NULL, '2026-03-12 04:31:46.178518+00', '2026-03-12 04:31:46.178518+00');
INSERT INTO public.payment_invoices VALUES ('b148cb4e-905d-41e3-9017-bab601385dfe', 'a5c313aa-3464-48c3-baf6-ae6fd5befd62', 'xendit', 'noop-inv-784d0512', 'INV-ORD-20260326-03936280-3fbfea5c', 'https://checkout-bypass.example.com/noop-inv-784d0512', NULL, NULL, 1520000.00, 'IDR', 'pending', '2026-03-27 04:48:50+00', NULL, NULL, '2026-03-26 04:48:50.540824+00', '2026-03-26 04:48:50.540824+00');
INSERT INTO public.payment_invoices VALUES ('df7bf23e-4db0-4f1e-8ea0-23809ad0ce78', '99418086-e674-41ec-8dcc-6c4301202502', 'xendit', 'noop-inv-e2775f6d', 'INV-ORD-20260326-2A4E66F5-d34839f5', 'https://checkout-bypass.example.com/noop-inv-e2775f6d', 'SIMULATE', 'BYPASS', 1536000.00, 'IDR', 'paid', '2026-03-27 05:03:59+00', '2026-03-26 05:04:02+00', '{"webhook_payload": {"id": "sim-d7bbdec3", "amount": 1536000, "status": "PAID", "paid_at": "2026-03-26T05:04:02Z", "user_id": "", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260326-2A4E66F5-d34839f5", "paid_amount": 1536000, "payment_method": "SIMULATE", "payment_channel": "BYPASS"}}', '2026-03-26 05:03:59.310154+00', '2026-03-26 05:04:02.310819+00');
INSERT INTO public.payment_invoices VALUES ('950bcbb9-11af-48c8-9a20-2a3289e416fc', '1e88413e-31ae-4698-b321-a6862044d31c', 'xendit', 'noop-inv-a33dc80c', 'INV-ORD-20260326-017AFEBB-9e30ec7e', 'https://checkout-bypass.example.com/noop-inv-a33dc80c', 'SIMULATE', 'BYPASS', 61000.00, 'IDR', 'paid', '2026-03-27 05:09:06+00', '2026-03-26 05:09:10+00', '{"webhook_payload": {"id": "sim-d107f60a", "amount": 61000, "status": "PAID", "paid_at": "2026-03-26T05:09:10Z", "user_id": "", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260326-017AFEBB-9e30ec7e", "paid_amount": 61000, "payment_method": "SIMULATE", "payment_channel": "BYPASS"}}', '2026-03-26 05:09:06.094561+00', '2026-03-26 05:09:10.277574+00');


--
-- Data for Name: payment_events; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.payment_events VALUES ('a4abe5cc-6431-475e-a139-49449198b14c', 'a8a18b80-5aa5-4237-884d-35fd3f92c386', 'invoice.paid', 'evt_seed_001', '{"id": "xnd_seed_001", "event": "invoice.paid", "status": "PAID"}', '2026-03-08 07:34:21.226459+00');
INSERT INTO public.payment_events VALUES ('21cf2f05-4211-47e4-a32a-b59a00f03f57', 'df7bf23e-4db0-4f1e-8ea0-23809ad0ce78', 'invoice.paid', 'sim-d7bbdec3:PAID', '{"id": "sim-d7bbdec3", "amount": 1536000, "status": "PAID", "paid_at": "2026-03-26T05:04:02Z", "external_id": "INV-ORD-20260326-2A4E66F5-d34839f5", "paid_amount": 1536000, "payment_method": "SIMULATE", "payment_channel": "BYPASS"}', '2026-03-26 05:04:02.307989+00');
INSERT INTO public.payment_events VALUES ('678bebab-22ba-4d00-bd19-4e7139626dc6', '950bcbb9-11af-48c8-9a20-2a3289e416fc', 'invoice.paid', 'sim-d107f60a:PAID', '{"id": "sim-d107f60a", "amount": 61000, "status": "PAID", "paid_at": "2026-03-26T05:09:10Z", "external_id": "INV-ORD-20260326-017AFEBB-9e30ec7e", "paid_amount": 61000, "payment_method": "SIMULATE", "payment_channel": "BYPASS"}', '2026-03-26 05:09:10.275502+00');


--
-- Data for Name: payout_batches; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.payout_batches VALUES ('d7660a6f-1334-447b-bab7-7ab9e9066125', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', '2026-02-11', '2026-03-12', 'completed', 185000.00, 2500.00, 182500.00, '2026-03-12 07:24:21.226459+00', '7907d953-ba02-40a4-b403-07b68d8471e8', NULL, NULL, NULL, NULL, '2026-03-11 07:24:21.226459+00', '2026-03-12 07:24:21.226459+00');


--
-- Data for Name: payout_items; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.payout_items VALUES ('1ab40eae-e67d-4454-ba37-1dc14345c377', 'd7660a6f-1334-447b-bab7-7ab9e9066125', '80432903-81b1-48d6-82ba-670d2798a0d5', 185000.00, 2500.00, 182500.00);


--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_images VALUES ('cafd4659-1e16-467d-a1af-a57165343cbb', 'b2999c9e-df33-48ef-a109-38c77e51f25d', 'https://placehold.co/600x600?text=Kain+Ihram', 'image/png', NULL, true, 0, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_reviews; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_reviews VALUES ('f21d4561-2de3-4d9a-812e-2582c69330af', 'b2999c9e-df33-48ef-a109-38c77e51f25d', '80432903-81b1-48d6-82ba-670d2798a0d5', '6ca279d6-8b8b-4ded-b035-3ea182014a7e', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 4, 'Kainnya cukup bagus dan lembut, tapi warna agak berbeda dari foto. Overall puas.', 'published', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_review_images; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_review_images VALUES ('7dc0b2b0-1ffb-4b66-9277-f8e99130f459', 'f21d4561-2de3-4d9a-812e-2582c69330af', 'reviews/seed-review-001.jpg', 'image/jpeg', 153600, 0, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: product_review_stats; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_review_stats VALUES ('b2999c9e-df33-48ef-a109-38c77e51f25d', 1, 4, 4.00, 0, 0, 0, 0, 1, 0, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: refunds; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.refunds VALUES ('c4b049f4-91b0-4b22-87ac-6df46f01ee6d', '80432903-81b1-48d6-82ba-670d2798a0d5', 'a8a18b80-5aa5-4237-884d-35fd3f92c386', 50000.00, 'Barang tidak sesuai - refund partial', 'processed', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-12 07:24:21.226459+00');


--
-- Data for Name: reply_templates; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.reply_templates VALUES ('c0e89a64-456f-4df8-bf77-2b8acd4b1b25', 'Salam Pembuka', 'Assalamualaikum, terima kasih telah menghubungi kami. Ada yang bisa kami bantu?', true, 'cc000001-0000-0000-0000-000000000002', NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '/salam', 'general');
INSERT INTO public.reply_templates VALUES ('db924680-cc08-4bf0-9a4c-d54605ff89c1', 'Selamat Datang di UmrahMart', 'Assalamu''alaikum, selamat datang di UmrahMart! 🕌

Kami hadir untuk memudahkan Anda mendapatkan perlengkapan haji dan umroh terbaik. Ada yang bisa kami bantu hari ini?', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/selamat-datang', 'general');
INSERT INTO public.reply_templates VALUES ('b7ec135c-0073-43ec-a7ee-65fecb5668b8', 'Terima Kasih Telah Menghubungi Kami', 'Terima kasih telah menghubungi tim Customer Service UmrahMart. 🙏

Kami akan segera membantu Anda. Mohon tunggu sebentar ya.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/terima-kasih-menghubungi', 'general');
INSERT INTO public.reply_templates VALUES ('9950affe-68af-4a87-bddd-57cb2b2d4522', 'Apakah Masih Ada Yang Bisa Dibantu?', 'Apakah masih ada yang bisa kami bantu? 😊

Jika masalah Anda sudah teratasi, jangan lupa berikan ulasan untuk membantu kami meningkatkan layanan. Semoga ibadah haji/umroh Anda berjalan lancar. Aamiin.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/ada-yang-bisa-dibantu', 'general');
INSERT INTO public.reply_templates VALUES ('2210ed2d-6243-47da-af10-cd789f2416b0', 'Jam Operasional Customer Service', 'Tim Customer Service UmrahMart melayani Anda setiap hari:

🕐 Senin – Jumat : 08.00 – 21.00 WIB
🕐 Sabtu – Minggu : 09.00 – 18.00 WIB

Di luar jam operasional, Anda tetap bisa meninggalkan pesan dan kami akan merespons saat jam kerja.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/jam-operasional', 'general');
INSERT INTO public.reply_templates VALUES ('fd06b343-9d5a-49bf-ad6a-427fa4b238b0', 'Konfirmasi Pesanan Diterima', 'Alhamdulillah, pesanan Anda telah kami terima! ✅

Nomor pesanan Anda: {order_number}
Status: Menunggu Pembayaran

Silakan selesaikan pembayaran sebelum {expired_at} agar pesanan dapat segera kami proses. Jazakallahu khairan.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/konfirmasi-pesanan-diterima', 'order');
INSERT INTO public.reply_templates VALUES ('2c22e30d-7e7a-484f-ab6e-bd3b6f90dd29', 'Pesanan Sedang Diproses', 'Kabar baik! Pesanan Anda dengan nomor {order_number} sedang kami proses. 📦

Estimasi pesanan siap dikirim: 1–2 hari kerja.
Kami akan menginformasikan nomor resi pengiriman segera setelah paket diserahkan ke kurir.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/pesanan-diproses', 'order');
INSERT INTO public.reply_templates VALUES ('0dcf4408-01fb-49de-b782-03c949175719', 'Pesanan Berhasil Dibatalkan', 'Pesanan Anda dengan nomor {order_number} telah berhasil dibatalkan.

Jika Anda sudah melakukan pembayaran, proses refund akan kami lakukan dalam 3–7 hari kerja ke metode pembayaran asal.

Apabila ada pertanyaan lebih lanjut, jangan ragu untuk menghubungi kami kembali.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/pesanan-dibatalkan', 'order');
INSERT INTO public.reply_templates VALUES ('09e4155c-ccaf-4466-99e4-6f09c12682c8', 'Detail Pesanan', 'Berikut detail pesanan Anda:

📋 Nomor Pesanan : {order_number}
📅 Tanggal       : {order_date}
💰 Total         : {total_amount}
📍 Status        : {order_status}

Untuk informasi lebih lengkap, silakan cek halaman "Pesanan Saya" di aplikasi.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/detail-pesanan', 'order');
INSERT INTO public.reply_templates VALUES ('039e2bf9-5a35-47fa-b5c6-df6550e1c4bd', 'Pembayaran Berhasil Diterima', 'Alhamdulillah, pembayaran Anda telah berhasil kami terima! 💚

Nomor pesanan : {order_number}
Jumlah        : {amount}
Metode        : {payment_method}

Pesanan Anda akan segera kami proses. Terima kasih.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/pembayaran-berhasil', 'payment');
INSERT INTO public.reply_templates VALUES ('279010e0-341c-42b1-8818-bf5e5b0d2bea', 'Menunggu Konfirmasi Pembayaran', 'Halo, kami melihat pesanan Anda belum terkonfirmasi pembayarannya.

Batas waktu pembayaran: {expired_at}

Jika Anda sudah melakukan transfer manual, mohon upload bukti pembayaran melalui aplikasi agar kami dapat segera memprosesnya.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/menunggu-konfirmasi', 'payment');
INSERT INTO public.reply_templates VALUES ('c5c27b1c-05d5-4ae6-83a0-7bb1c89c74f9', 'Pembayaran Gagal / Kedaluwarsa', 'Kami informasikan bahwa pembayaran untuk pesanan {order_number} telah gagal atau melewati batas waktu. ⚠️

Pesanan Anda telah kami batalkan secara otomatis. Anda dapat melakukan pemesanan ulang kapan saja.

Mohon pastikan saldo/limit mencukupi saat melakukan pembayaran berikutnya.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/pembayaran-gagal', 'payment');
INSERT INTO public.reply_templates VALUES ('98ccc880-d214-4c26-bb39-c6fc1cc4bfe3', 'Cara Melakukan Pembayaran', 'Berikut metode pembayaran yang tersedia di UmrahMart:

💳 Transfer Bank (BCA, Mandiri, BNI, BRI)
📱 E-Wallet (GoPay, OVO, DANA, ShopeePay)
🏪 Gerai Retail (Alfamart, Indomaret)
💵 QRIS

Pilih metode yang paling nyaman untuk Anda saat checkout.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/cara-pembayaran', 'payment');
INSERT INTO public.reply_templates VALUES ('022b8cc0-8129-45d4-a9b5-0601ce941461', 'Pengajuan Refund Diterima', 'Pengajuan refund Anda telah kami terima dan sedang dalam proses review. ✅

Nomor tiket refund: {ticket_number}

Tim kami akan memverifikasi dalam 1–3 hari kerja. Anda akan mendapat notifikasi setelah proses selesai.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/pengajuan-diterima', 'refund');
INSERT INTO public.reply_templates VALUES ('60cf1bc4-50e9-49ca-a6fa-dc79ac495678', 'Refund Sedang Diproses', 'Refund Anda sedang dalam proses pencairan. 🔄

Jumlah refund : {refund_amount}
Tujuan        : {refund_destination}
Estimasi      : 3–7 hari kerja

Mohon bersabar dan pastikan rekening/e-wallet tujuan masih aktif.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/sedang-diproses', 'refund');
INSERT INTO public.reply_templates VALUES ('83ec97f6-b6fc-4e70-9491-ac41192fe241', 'Refund Berhasil', 'Dana refund sebesar {refund_amount} telah berhasil kami kirimkan ke {refund_destination}. 🎉

Jika dalam 1×24 jam dana belum masuk, mohon hubungi kami kembali dengan menyertakan nomor tiket {ticket_number}.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/berhasil', 'refund');
INSERT INTO public.reply_templates VALUES ('c0aced63-92c1-4857-8db0-91187f1822ad', 'Pesanan Telah Dikirim', 'Pesanan Anda sudah dalam perjalanan! 🚚

Nomor Pesanan : {order_number}
Kurir         : {courier_name}
Nomor Resi    : {tracking_number}

Lacak pengiriman di: {tracking_url}

Estimasi tiba: {estimated_arrival}', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/pesanan-dikirim', 'shipping');
INSERT INTO public.reply_templates VALUES ('933272ef-be7b-44dd-9fb8-8eaa07c48592', 'Konfirmasi Penerimaan Paket', 'Halo, berdasarkan data kurir, paket Anda sudah dinyatakan terkirim pada {delivered_at}.

Apakah paket sudah Anda terima dengan kondisi baik? Mohon konfirmasi agar pesanan dapat diselesaikan. 📦✅', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/konfirmasi-penerimaan', 'shipping');
INSERT INTO public.reply_templates VALUES ('57dc278b-9cab-43bf-8836-901e1b4fe9cf', 'Paket Tertahan / Terlambat', 'Kami melihat terdapat kendala pada pengiriman paket Anda (no. resi: {tracking_number}). ⚠️

Tim kami sedang berkoordinasi dengan pihak kurir untuk menindaklanjuti hal ini. Kami akan segera menginformasikan perkembangannya. Mohon maaf atas ketidaknyamanannya.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/paket-tertahan', 'shipping');
INSERT INTO public.reply_templates VALUES ('0c4d32dc-44fc-402b-9fa5-44a181768d95', 'Informasi Ketersediaan Stok', 'Terima kasih atas minat Anda pada produk {product_name}.

Saat ini stok produk tersebut {stock_status}.

Anda dapat mengaktifkan notifikasi "Ingatkan Saya" pada halaman produk agar kami bisa memberitahu Anda saat stok tersedia kembali.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/ketersediaan-stok', 'product');
INSERT INTO public.reply_templates VALUES ('f2db8965-96ee-4957-a034-9bafb1981291', 'Detail Spesifikasi Produk', 'Berikut informasi lengkap mengenai {product_name}:

📌 Bahan    : {material}
📐 Ukuran   : {size}
🎨 Warna    : {color}
🏷️ Berat    : {weight}
✅ Halal    : {halal_status}

Apakah ada pertanyaan lain mengenai produk ini?', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/spesifikasi-produk', 'product');
INSERT INTO public.reply_templates VALUES ('d1a8034b-7a77-47f1-b549-440ff1a13d66', 'Verifikasi Berhasil', 'Selamat! Verifikasi akun Anda telah berhasil. ✅

Akun Anda kini sudah terverifikasi dan dapat menikmati semua fitur UmrahMart tanpa batasan. Terima kasih atas kerja samanya.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/verification/berhasil', 'verification');
INSERT INTO public.reply_templates VALUES ('31b16e1b-998a-4779-87a2-663ef09cffb4', 'Produk Tidak Lagi Tersedia', 'Mohon maaf, produk {product_name} saat ini sudah tidak tersedia di katalog kami. 😔

Kami merekomendasikan produk serupa yang mungkin sesuai dengan kebutuhan Anda:
👉 {alternative_product}

Silakan cek koleksi terbaru kami di aplikasi.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/produk-tidak-tersedia', 'product');
INSERT INTO public.reply_templates VALUES ('9d2a3755-e301-4552-a55f-d1ca99ee9aab', 'Bantuan Reset Password', 'Untuk mereset password akun Anda, ikuti langkah berikut:

1️⃣ Buka halaman Login
2️⃣ Klik "Lupa Password?"
3️⃣ Masukkan email terdaftar Anda
4️⃣ Cek inbox email untuk link reset (cek juga folder Spam)
5️⃣ Klik link dan buat password baru

Link reset berlaku selama 60 menit. Jika belum menerima email, hubungi kami kembali.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/reset-password', 'account');
INSERT INTO public.reply_templates VALUES ('52d1c3ad-e619-47c3-930c-2497057c746c', 'Verifikasi Email', 'Untuk memverifikasi email Anda:

1️⃣ Cek inbox email {email} (termasuk folder Spam/Junk)
2️⃣ Buka email dari UmrahMart dengan subjek "Verifikasi Akun"
3️⃣ Klik tombol "Verifikasi Sekarang"

Jika email tidak ditemukan, klik "Kirim Ulang Verifikasi" di halaman akun Anda.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/verifikasi-email', 'account');
INSERT INTO public.reply_templates VALUES ('849c3f8f-05e6-4ad7-b8a7-9b39914b99ed', 'Akun Diblokir / Dinonaktifkan', 'Kami melihat akun Anda saat ini tidak aktif. 🔒

Hal ini dapat terjadi karena:
• Aktivitas yang mencurigakan terdeteksi
• Pelanggaran syarat & ketentuan penggunaan
• Permintaan penonaktifan sebelumnya

Untuk mengajukan reaktivasi, mohon kirimkan data diri Anda ke email support@umrahmart.id.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/akun-diblokir', 'account');
INSERT INTO public.reply_templates VALUES ('c729c57e-34d3-475d-8d65-7d0ff314a51d', 'Keluhan Diterima dan Dicatat', 'Kami sangat menyesal mendengar pengalaman yang tidak menyenangkan ini. 🙏

Keluhan Anda telah kami catat dengan nomor tiket {ticket_number}. Tim kami akan menginvestigasi dan menghubungi Anda dalam 1×24 jam kerja.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/keluhan-diterima', 'complaint');
INSERT INTO public.reply_templates VALUES ('08411013-5cb3-497d-87a5-14b18b27361d', 'Keluhan Sedang Diinvestigasi', 'Kami tengah menginvestigasi keluhan yang Anda sampaikan terkait {complaint_subject}.

Proses investigasi membutuhkan waktu maksimal 3 hari kerja. Kami mohon kesabaran Anda dan akan segera memberikan informasi lebih lanjut.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/sedang-investigasi', 'complaint');
INSERT INTO public.reply_templates VALUES ('2be83748-504e-4be7-8895-47ef175d1689', 'Keluhan Telah Diselesaikan', 'Keluhan Anda dengan nomor tiket {ticket_number} telah kami selesaikan. ✅

Solusi yang diberikan: {resolution}

Kami memohon maaf atas pengalaman yang kurang menyenangkan ini dan berterima kasih atas masukan Anda untuk perbaikan layanan kami.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/telah-diselesaikan', 'complaint');
INSERT INTO public.reply_templates VALUES ('29ef4e0e-fff3-4d1d-998f-b4b8c16956f6', 'Prosedur Pengembalian Barang', 'Berikut prosedur pengembalian barang (return) di UmrahMart:

1️⃣ Ajukan return melalui menu "Pesanan Saya" → "Ajukan Return"
2️⃣ Pilih produk dan alasan pengembalian
3️⃣ Unggah foto kondisi barang
4️⃣ Tunggu konfirmasi persetujuan (1–2 hari kerja)
5️⃣ Kirim barang ke alamat yang tertera

Syarat: barang dalam kondisi asli, belum dipakai, dan beserta kemasan lengkap. Maksimal 7 hari setelah diterima.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/prosedur-pengembalian', 'return');
INSERT INTO public.reply_templates VALUES ('8c893df1-0794-4a42-99d1-dbc94ca8c470', 'Pengajuan Return Disetujui', 'Kabar baik! Pengajuan return Anda untuk pesanan {order_number} telah disetujui. ✅

Silakan kirimkan barang ke:
📍 {return_address}
Atas nama: UmrahMart Returns

Setelah barang kami terima dan verifikasi, dana akan dikembalikan dalam 3–5 hari kerja.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/disetujui', 'return');
INSERT INTO public.reply_templates VALUES ('a4203d8e-4615-4ef0-b3eb-1a90bc79924e', 'Pengajuan Return Ditolak', 'Mohon maaf, pengajuan return untuk pesanan {order_number} tidak dapat kami proses. ❌

Alasan: {rejection_reason}

Jika Anda merasa keberatan, Anda dapat mengajukan banding dengan menghubungi kami dan melampirkan bukti pendukung.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/ditolak', 'return');
INSERT INTO public.reply_templates VALUES ('81873554-d801-46b4-8c20-a68f0aa0b604', 'Informasi Promo Aktif', 'Halo! Berikut promo yang sedang berlangsung di UmrahMart 🎉:

🏷️ {promo_name}
💰 Diskon hingga {discount_value}
📅 Berlaku: {promo_start} – {promo_end}
📝 Syarat: {promo_terms}

Jangan sampai terlewat! Belanja sekarang sebelum promo berakhir.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/promo/info-promo-aktif', 'promo');
INSERT INTO public.reply_templates VALUES ('6dcdda2a-39dc-4762-b8f0-754619e8da71', 'Promo Tidak Berlaku untuk Pesanan Ini', 'Mohon maaf, promo yang Anda gunakan tidak berlaku untuk pesanan ini. ⚠️

Kemungkinan penyebab:
• Promo sudah berakhir
• Produk tidak termasuk dalam kategori promo
• Minimum pembelian tidak terpenuhi

Cek syarat & ketentuan promo di halaman "Promo" pada aplikasi.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/promo/tidak-berlaku', 'promo');
INSERT INTO public.reply_templates VALUES ('716a7e0f-9c10-404c-8fe1-a019a38adc7d', 'Cara Menggunakan Voucher', 'Berikut cara menggunakan voucher di UmrahMart:

1️⃣ Tambahkan produk ke keranjang
2️⃣ Masuk ke halaman Checkout
3️⃣ Klik "Gunakan Voucher"
4️⃣ Masukkan kode voucher: {voucher_code}
5️⃣ Klik "Terapkan" — diskon akan langsung terpotong

Pastikan pesanan Anda memenuhi minimum pembelian voucher.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/voucher/cara-menggunakan', 'voucher');
INSERT INTO public.reply_templates VALUES ('56959391-4b07-4cb8-89ab-a61abbef4748', 'Voucher Tidak Valid atau Kedaluwarsa', 'Kode voucher yang Anda masukkan tidak dapat digunakan. ❌

Kemungkinan penyebab:
• Kode voucher salah atau sudah digunakan
• Voucher sudah kedaluwarsa ({expired_date})
• Akun Anda tidak memenuhi syarat voucher

Untuk bantuan lebih lanjut, silakan kirimkan screenshot kode voucher Anda.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/voucher/tidak-valid', 'voucher');
INSERT INTO public.reply_templates VALUES ('7ca45115-0adc-4c23-a1a1-5924e6177e20', 'Masalah Teknis pada Aplikasi', 'Kami mohon maaf atas kendala teknis yang Anda alami. 🛠️

Beberapa langkah yang dapat dicoba:
1. Tutup dan buka kembali aplikasi
2. Pastikan koneksi internet stabil
3. Update aplikasi ke versi terbaru
4. Hapus cache aplikasi
5. Restart perangkat

Jika masalah berlanjut, mohon kirimkan screenshot error beserta tipe perangkat dan versi OS Anda.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/masalah-teknis-aplikasi', 'technical');
INSERT INTO public.reply_templates VALUES ('c139c11f-0159-4e62-b765-58dcc55b58eb', 'Tidak Bisa Login ke Akun', 'Kami memahami betapa frustrasinya tidak bisa masuk ke akun. Mari kami bantu! 🔑

Silakan coba langkah berikut:
1. Pastikan email dan password sudah benar
2. Aktifkan Caps Lock/periksa huruf kapital
3. Gunakan fitur "Lupa Password" untuk reset
4. Coba login dari perangkat atau browser berbeda

Apakah Anda menerima pesan error tertentu? Mohon informasikan agar kami dapat membantu lebih lanjut.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/tidak-bisa-login', 'technical');
INSERT INTO public.reply_templates VALUES ('82f514a7-153d-4f59-9327-229778c82585', 'Fitur Sedang Dalam Perbaikan', 'Kami menginformasikan bahwa fitur {feature_name} saat ini sedang dalam perbaikan/maintenance. 🔧

Estimasi selesai: {maintenance_end}

Kami mohon maaf atas ketidaknyamanannya. Tim teknis kami sedang bekerja keras untuk memulihkan layanan secepatnya.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/fitur-dalam-perbaikan', 'technical');
INSERT INTO public.reply_templates VALUES ('774c04c8-9d21-4d1a-8c86-3ab0ef2987a0', 'Verifikasi Identitas Diperlukan', 'Untuk keamanan akun Anda, kami perlu melakukan verifikasi identitas. 🔐

Dokumen yang diperlukan:
📄 KTP (foto depan, jelas dan tidak blur)
🤳 Selfie sambil memegang KTP

Kirimkan dokumen melalui email ke: verify@umrahmart.id
Subject: Verifikasi Identitas - {user_id}

Proses verifikasi membutuhkan 1–2 hari kerja.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/verification/identitas-diperlukan', 'verification');
INSERT INTO public.reply_templates VALUES ('967434ba-1ba8-459e-9514-f8c5dfdc3ec1', 'Informasi Pendaftaran Vendor', 'Terima kasih atas minat Anda untuk bergabung sebagai vendor di UmrahMart! 🏪

Persyaratan pendaftaran:
✅ KTP pemilik usaha
✅ Foto toko/tempat usaha
✅ Logo bisnis
✅ Banner bisnis
✅ Buku rekening bank
✅ NPWP (opsional)

Proses review membutuhkan 2–3 hari kerja setelah semua dokumen lengkap.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/info-pendaftaran', 'vendor');
INSERT INTO public.reply_templates VALUES ('96cb369d-1535-4e9f-b794-2db6a703901f', 'Status Vendor Sedang Direview', 'Pendaftaran vendor Anda sedang dalam proses review oleh tim kami. ⏳

Kami akan menginformasikan hasilnya melalui email dan notifikasi aplikasi dalam 2–3 hari kerja.

Pastikan data dan dokumen yang Anda unggah sudah lengkap dan terbaca dengan jelas.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/sedang-direview', 'vendor');
INSERT INTO public.reply_templates VALUES ('bf62a18f-4d7d-48f0-99a3-fa1fae0cf4fc', 'Vendor Berhasil Diaktifkan', 'Selamat! Akun vendor Anda telah berhasil diaktifkan! 🎉

Anda sekarang dapat mulai:
🛍️ Menambahkan produk ke katalog
📊 Mengelola stok dan harga
📦 Menerima dan memproses pesanan

Untuk panduan lengkap, kunjungi halaman "Panduan Vendor" di dashboard Anda. Semoga sukses!', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/berhasil-diaktifkan', 'vendor');
INSERT INTO public.reply_templates VALUES ('df27202f-3ef0-4ac9-84ea-e14506b78e01', 'Notifikasi Stok Tersedia', 'Kabar gembira! 🎉 Produk {product_name} yang Anda tunggu-tunggu kini sudah tersedia kembali.

Stok tersisa: {remaining_stock} pcs

Segera pesan sebelum kehabisan! Klik di sini: {product_url}', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/stock/stok-tersedia', 'stock');
INSERT INTO public.reply_templates VALUES ('30cf40bb-4b5a-42d1-a288-7b1b3ff5ceb5', 'Informasi Pre-Order', 'Produk {product_name} saat ini tersedia dalam mode Pre-Order. 📋

Detail Pre-Order:
📅 Batas pemesanan : {po_end_date}
🚚 Estimasi kirim  : {estimated_delivery}
💰 Harga PO        : {po_price}

Pesan sekarang untuk mendapatkan kepastian stok!', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/stock/info-pre-order', 'stock');
INSERT INTO public.reply_templates VALUES ('24e59108-be8b-4641-85f0-8f11128cf74c', 'Konfirmasi Pembatalan Pesanan', 'Kami telah menerima permintaan pembatalan untuk pesanan {order_number}. ℹ️

Untuk melanjutkan proses pembatalan, mohon konfirmasi alasan Anda:
a) Salah alamat pengiriman
b) Ingin ganti produk
c) Menemukan harga lebih murah
d) Berubah pikiran
e) Lainnya

Balas dengan huruf pilihan Anda.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/konfirmasi-pembatalan', 'cancellation');
INSERT INTO public.reply_templates VALUES ('7ffccfea-bef8-45d5-b8cb-c46798a0da55', 'Pesanan Tidak Dapat Dibatalkan', 'Mohon maaf, pesanan {order_number} tidak dapat dibatalkan karena status pesanan sudah dalam tahap {order_status}. ⚠️

Jika barang sudah diterima dan terdapat masalah, Anda masih bisa mengajukan return dalam 7 hari setelah penerimaan.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/tidak-dapat-dibatalkan', 'cancellation');
INSERT INTO public.reply_templates VALUES ('f04503ae-fb8b-4e54-b17d-b2473014537f', 'Pembatalan Berhasil', 'Pesanan {order_number} berhasil dibatalkan. ✅

Jika ada pembayaran yang perlu dikembalikan:
💰 Jumlah refund : {refund_amount}
⏱️ Estimasi dana masuk : 3–7 hari kerja

Terima kasih atas pengertian Anda.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/berhasil', 'cancellation');
INSERT INTO public.reply_templates VALUES ('45b72f3e-52d9-499e-8d31-faa17af981be', 'Jadwal Estimasi Pengiriman', 'Berikut estimasi pengiriman berdasarkan lokasi Anda:

🏙️ Jabodetabek    : 1–2 hari kerja
🗺️ Pulau Jawa     : 2–3 hari kerja
🏝️ Luar Jawa      : 3–5 hari kerja
🗾 Papua/Terpencil : 5–10 hari kerja

Estimasi dihitung sejak paket diserahkan ke kurir, tidak termasuk hari libur nasional.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/delivery/estimasi-pengiriman', 'delivery');
INSERT INTO public.reply_templates VALUES ('ca6a9164-b296-4daf-9131-6902d1a3aa6e', 'Pengiriman Tertunda', 'Kami mohon maaf atas keterlambatan pengiriman paket Anda (no. resi: {tracking_number}). 🙏

Kendala yang terjadi: {delay_reason}

Tim kami sedang berkoordinasi dengan pihak kurir {courier_name} untuk mempercepat proses pengiriman. Kami akan menginformasikan update terbaru sesegera mungkin.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/delivery/pengiriman-tertunda', 'delivery');
INSERT INTO public.reply_templates VALUES ('a217d188-4f32-4184-80a2-cadda1107c5d', 'Informasi Pengambilan di Toko (Pickup)', 'Pesanan Anda dapat diambil langsung di:

📍 {store_address}
🕐 Jam operasional: {store_hours}

Tunjukkan nomor pesanan {order_number} atau kode QR di aplikasi kepada petugas. Pesanan akan disimpan selama 3 hari sejak notifikasi siap diambil.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/pickup/info-pickup', 'pickup');
INSERT INTO public.reply_templates VALUES ('cb92a7d1-22cc-4fb6-afa0-3b7c0140eb93', 'Pesanan Siap Diambil', 'Pesanan Anda sudah siap untuk diambil! 🛍️

Nomor Pesanan : {order_number}
Lokasi        : {store_address}
Batas Ambil   : {pickup_deadline}

Jangan lupa bawa identitas diri saat pengambilan. Sampai jumpa!', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/pickup/siap-diambil', 'pickup');
INSERT INTO public.reply_templates VALUES ('a2940fd9-f18b-4024-8694-29be71b0cf23', 'Informasi Paket Langganan', 'Berikut informasi paket langganan UmrahMart Premium:

⭐ Paket Basic   : Rp 29.000/bulan — Gratis ongkir 2x/bulan
⭐ Paket Premium : Rp 59.000/bulan — Gratis ongkir unlimited + cashback 5%
⭐ Paket VIP     : Rp 99.000/bulan — Semua benefit + akses produk eksklusif

Daftar sekarang di menu "Langganan" pada aplikasi.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/subscription/info-langganan', 'subscription');
INSERT INTO public.reply_templates VALUES ('fec12945-f26b-4e5b-892d-e53b916e5ad3', 'Masa Aktif Langganan Hampir Habis', 'Halo! Masa aktif langganan UmrahMart Premium Anda akan berakhir pada {expiry_date}. ⏰

Perpanjang sekarang dan nikmati terus benefit:
✅ Gratis ongkos kirim
✅ Cashback eksklusif member
✅ Early access produk baru

Klik di sini untuk perpanjang: {renewal_url}', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/subscription/hampir-habis', 'subscription');
INSERT INTO public.reply_templates VALUES ('8ae2e394-0d0c-4fb7-8744-f13bbb7d448e', 'Permintaan Ulasan Produk', 'Assalamu''alaikum! Semoga pesanan Anda sudah diterima dengan baik. 😊

Kami akan sangat berterima kasih jika Anda meluangkan waktu untuk memberikan ulasan pada produk {product_name}.

Ulasan Anda sangat membantu calon pembeli lain dalam memilih produk yang tepat. Berikan ulasan di menu "Pesanan Saya" → "Beri Ulasan".', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/review/minta-ulasan', 'review');
INSERT INTO public.reply_templates VALUES ('011a94d2-8ce3-4594-87ba-63e9b46fe1eb', 'Terima Kasih Atas Ulasan Anda', 'Jazakallahu khairan atas ulasan yang telah Anda berikan! 🌟

Masukan Anda sangat berharga bagi kami untuk terus meningkatkan kualitas produk dan layanan. Semoga UmrahMart selalu bisa menjadi teman setia perjalanan ibadah Anda.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/review/terima-kasih-ulasan', 'review');
INSERT INTO public.reply_templates VALUES ('23b40e79-e1f3-4c17-a3ee-64ec201cd2ae', 'Laporan Penipuan Diterima', 'Laporan penipuan Anda telah kami terima dan kami tangani dengan sangat serius. 🚨

Nomor laporan: {report_number}

Tim keamanan kami akan menginvestigasi dalam 1×24 jam. Jangan lakukan transfer atau berikan data sensitif kepada pihak yang mencurigakan.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/fraud/laporan-penipuan-diterima', 'fraud');
INSERT INTO public.reply_templates VALUES ('53916ccc-b29b-40cb-a795-dc5f7de740e5', 'Tindakan Keamanan Akun', 'Kami mendeteksi aktivitas tidak biasa pada akun Anda. 🔐

Sebagai tindakan keamanan, akun Anda telah kami sementara kunci.

Langkah yang perlu Anda lakukan:
1. Segera ubah password Anda
2. Aktifkan verifikasi 2 langkah
3. Hubungi kami melalui {support_contact} untuk membuka kunci akun

Jangan pernah bagikan OTP atau password kepada siapapun, termasuk tim kami.', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/fraud/tindakan-keamanan-akun', 'fraud');
INSERT INTO public.reply_templates VALUES ('1a57d5dc-1751-43c9-bdf7-9ccfc1d9898c', 'Pertanyaan Tidak Dapat Kami Jawab Saat Ini', 'Terima kasih atas pertanyaan Anda. Pertanyaan ini membutuhkan penanganan dari tim khusus kami. 🔄

Kami akan meneruskan pertanyaan ini ke tim terkait dan akan menghubungi Anda kembali dalam 1×24 jam kerja.

Nomor tiket Anda: {ticket_number}', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/other/eskalasi-ke-tim-terkait', 'other');
INSERT INTO public.reply_templates VALUES ('9a10b61b-1dc3-407a-b5a5-b80062a404b4', 'Cara Menghubungi Kami', 'Ada beberapa cara untuk menghubungi tim UmrahMart:

💬 Live Chat  : Aplikasi UmrahMart (menu "Bantuan")
📧 Email      : support@umrahmart.id
📞 Telepon    : 021-XXXX-XXXX (08.00–21.00 WIB)
📱 WhatsApp   : wa.me/628XXXXXXXXX
📘 Instagram  : @umrahmart.id

Kami siap membantu Anda! 🙏', true, 'cc000001-0000-0000-0000-000000000001', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/other/cara-menghubungi-kami', 'other');


--
-- Data for Name: return_reasons; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.return_reasons VALUES (1, 'Barang tidak sesuai deskripsi', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.schema_migrations VALUES (47, false);


--
-- Data for Name: shipments; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.shipments VALUES ('04c23c0f-1620-400c-b27a-949a27a5309e', '382762cd-f805-41a0-ad67-07c6a33020d2', 'anteraja', 'ECO', '11002960818873', 'shipped', '2026-03-12 04:42:01.3673+00', NULL, '2-3 day');
INSERT INTO public.shipments VALUES ('de7cabce-2c88-4a0a-a7dc-e986fa75268e', '80432903-81b1-48d6-82ba-670d2798a0d5', 'jne', 'REG', 'JNESEED0000001', 'delivered', '2026-03-11 07:24:21.226459+00', '2026-03-12 07:24:21.226459+00', '2-3 hari');
INSERT INTO public.shipments VALUES ('b4318ee7-8458-49d7-be8b-56e2a14a8255', 'a5c313aa-3464-48c3-baf6-ae6fd5befd62', 'wahana', 'NextDay', '', 'waiting_pickup', NULL, NULL, '1 day');
INSERT INTO public.shipments VALUES ('d8faa35a-c7eb-4cb5-bef9-3be855695e54', '99418086-e674-41ec-8dcc-6c4301202502', 'anteraja', 'ECO', '11002960818873', 'shipped', '2026-03-26 05:07:14.126627+00', NULL, '3 day');
INSERT INTO public.shipments VALUES ('c19f3861-9019-49c0-baba-5d75a64f67bf', '1e88413e-31ae-4698-b321-a6862044d31c', 'wahana', 'Express', '', 'waiting_pickup', NULL, NULL, '2 day');


--
-- Data for Name: ticket_subjects; Type: TABLE DATA; Schema: public; Owner: -
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
-- Data for Name: tickets; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.tickets VALUES ('7bbf0cff-c50c-4a81-a116-6381d0fcec70', 'TKT-20260313-0001', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'cc000001-0000-0000-0000-000000000002', 'ORD-20260313-S001', '081200000002', 'Dev Customer', 'Barang Tidak Sesuai', 'Barang yang diterima tidak sesuai dengan deskripsi produk. Warna dan ukuran berbeda dari yang dipesan.', 'closed', 'web', '2026-03-13 07:24:21.226459+00', '2026-03-12 19:24:21.226459+00', '2026-03-13 07:24:21.226459+00', NULL, NULL, 4);
INSERT INTO public.tickets VALUES ('450c957b-f4a9-418e-ba12-a55a33fe7627', 'TKT-20260316-0001', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'cc000001-0000-0000-0000-000000000002', 'ORD-20260225-001', '081234567890', 'Test Customer', 'Pengiriman Terlambat', 'Sudah 7 hari tapi pesanan perlengkapan haji saya belum sampai. Mohon dicek.', 'on_progress', 'web', NULL, '2026-03-16 03:44:42.689908+00', '2026-03-16 03:45:38.189758+00', 'https://kemenhaj.s3.ap-southeast-1.amazonaws.com/ticket-attachments/2026/03/16/a4292972-521a-42bc-9944-4d934ecffa99', 'image/jpeg', 5);


--
-- Data for Name: ticket_messages; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ticket_messages VALUES ('406b9695-60be-4313-a6e8-c7e3cbb515eb', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Barang yang saya terima tidak sesuai dengan deskripsi. Mohon bantuan pengembaliannya.', false, false, '2026-03-12 19:24:21.226459+00');
INSERT INTO public.ticket_messages VALUES ('689921d7-ecb1-4825-acad-384c83bec881', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'Tes catatan', true, false, '2026-03-16 03:45:54.521898+00');
INSERT INTO public.ticket_messages VALUES ('f4587c2c-9de6-43f2-9074-a066ebe15d6c', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'Tes internal', true, true, '2026-03-16 03:46:38.261261+00');
INSERT INTO public.ticket_messages VALUES ('01f46127-cda7-4f38-b7dc-e618edf6e2b7', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'cek lagi', true, false, '2026-03-16 03:46:57.818314+00');
INSERT INTO public.ticket_messages VALUES ('1b3248af-7f14-4162-9d40-071da847f06b', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'coba ke 4', true, false, '2026-03-16 03:47:05.117601+00');
INSERT INTO public.ticket_messages VALUES ('0015a95c-a538-47a8-98c4-ecd6f12e47ce', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'Coba ke 5', true, false, '2026-03-16 03:47:09.340605+00');
INSERT INTO public.ticket_messages VALUES ('ef57170d-4e5d-45ed-a5d3-b3c903313986', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'Coba ke 6', true, false, '2026-03-16 03:47:14.955477+00');
INSERT INTO public.ticket_messages VALUES ('b620160c-161f-44bc-b0f8-af6178927d44', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'Coba ke 7', true, false, '2026-03-16 03:47:20.29548+00');
INSERT INTO public.ticket_messages VALUES ('6b2518d7-79d2-4f83-be49-f286c7f66296', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'Coba ke 8', true, false, '2026-03-16 03:47:24.673385+00');
INSERT INTO public.ticket_messages VALUES ('85003e1b-5349-45f4-8af6-4f2fe877bf75', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'Cek 9', true, false, '2026-03-16 03:49:17.518821+00');
INSERT INTO public.ticket_messages VALUES ('4a2abb72-d454-4329-9a3f-39cddb141c0c', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'tes 10', true, false, '2026-03-16 03:50:54.476421+00');
INSERT INTO public.ticket_messages VALUES ('7e12c108-0a00-41de-a569-406722590d78', '450c957b-f4a9-418e-ba12-a55a33fe7627', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Oke masuk', false, false, '2026-03-16 03:52:11.251387+00');
INSERT INTO public.ticket_messages VALUES ('146c36d3-0b93-406a-a336-c10f60d00003', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'mantap', true, false, '2026-03-16 03:53:51.573787+00');
INSERT INTO public.ticket_messages VALUES ('bde4acc2-c737-4c58-acc1-64c482121003', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'tes fetch cs 1', true, false, '2026-03-16 03:54:22.110705+00');
INSERT INTO public.ticket_messages VALUES ('1b00042f-0a04-4977-aaaa-7ecfa81afba9', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'lagi', true, false, '2026-03-16 03:54:36.963792+00');
INSERT INTO public.ticket_messages VALUES ('b5689879-e2ca-4eb2-8145-38e946545022', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'last', true, false, '2026-03-16 03:55:41.213226+00');
INSERT INTO public.ticket_messages VALUES ('9cf1bc0b-4ebd-4144-9d86-d351314568c9', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'last no fake', true, false, '2026-03-16 03:56:14.679372+00');
INSERT INTO public.ticket_messages VALUES ('8c42a3a4-b202-4dcd-a94a-5349f4ad9fc5', '450c957b-f4a9-418e-ba12-a55a33fe7627', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Oke masuk cek 10', false, false, '2026-03-16 03:57:11.492892+00');
INSERT INTO public.ticket_messages VALUES ('fb3d6a7d-7a3e-4dc1-b289-1813f882ffb2', '450c957b-f4a9-418e-ba12-a55a33fe7627', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Get messagenya tiap 5 menit klo ga salah', false, false, '2026-03-16 03:57:11.49536+00');
INSERT INTO public.ticket_messages VALUES ('971155d3-c50c-4a98-af50-79af5c973cdd', '450c957b-f4a9-418e-ba12-a55a33fe7627', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Masuk lagi', false, false, '2026-03-16 03:57:11.497077+00');
INSERT INTO public.ticket_messages VALUES ('2034a1c8-2176-477a-aed9-fa8f50f30802', '450c957b-f4a9-418e-ba12-a55a33fe7627', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Oke masuk test fetch 1', false, false, '2026-03-16 03:57:11.498492+00');
INSERT INTO public.ticket_messages VALUES ('5a7894bb-06af-4462-9b8e-d52aea20a9a2', '450c957b-f4a9-418e-ba12-a55a33fe7627', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'Okee, pesan muncul / GET tiap 5 menit', false, false, '2026-03-16 03:57:11.49983+00');


--
-- Data for Name: ticket_attachments; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ticket_attachments VALUES ('0f113c48-113c-4d01-8ce9-f7bb143fb24d', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', '406b9695-60be-4313-a6e8-c7e3cbb515eb', 'https://placehold.co/800x600?text=Bukti+Foto', 'bukti-foto.jpg', 'image/jpeg', 204800, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: ticket_status_logs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ticket_status_logs VALUES ('8ded3edf-4fd9-457c-8aee-aefc0f782ae0', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', 'cc000001-0000-0000-0000-000000000002', 'open', 'on_progress', 'CS mengambil tiket', '2026-03-13 01:24:21.226459+00');
INSERT INTO public.ticket_status_logs VALUES ('4f40cca2-b3d2-4c03-a060-69641718f153', '7bbf0cff-c50c-4a81-a116-6381d0fcec70', 'cc000001-0000-0000-0000-000000000002', 'on_progress', 'closed', 'Masalah telah diselesaikan', '2026-03-13 07:24:21.226459+00');
INSERT INTO public.ticket_status_logs VALUES ('d77b896b-7b33-4eec-9788-8c47c6ffd0c0', '450c957b-f4a9-418e-ba12-a55a33fe7627', 'cc000001-0000-0000-0000-000000000002', 'open', 'on_progress', NULL, '2026-03-16 03:45:38.189383+00');


--
-- Data for Name: vendor_balances; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_balances VALUES ('9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 1000000.00, 0.00, 1000000.00, 0.00, '2026-03-12 04:17:46.056983+00', 0.00);
INSERT INTO public.vendor_balances VALUES ('8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 500000.00, 0.00, 1500000.00, 1000000.00, '2026-03-13 07:24:21.226459+00', 0.00);
INSERT INTO public.vendor_balances VALUES ('ef0ad9ea-13b6-46ac-9bac-7a81af8ee705', 0.00, 0.00, 0.00, 0.00, '2026-03-16 04:14:11.907424+00', 0.00);


--
-- Data for Name: vendor_bank_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_bank_accounts VALUES ('7b2aee15-39da-4080-bf21-df2feb5c791e', '9851d7b4-7099-42a2-a7c8-4ce8b868fa33', 'BANK Dummy', '1234567890', 'Nursufyan Sauri', 'verified', NULL, '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-03-12 04:17:26.124526+00', '2026-03-12 04:17:26.124526+00', '2026-03-12 04:17:26.124526+00');
INSERT INTO public.vendor_bank_accounts VALUES ('841c5c84-cc53-4196-97f0-bd919bea6179', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'Bank Syariah Indonesia', '7788990011', 'PT Oleh-Oleh Haji Berkah', 'verified', NULL, '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-02-16 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_banners; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_banners VALUES ('dd107057-bd24-4da5-860a-6de0b6d7cdb5', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'Promo Ramadhan', 'https://placehold.co/1200x400?text=Promo+Ramadhan', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_couriers; Type: TABLE DATA; Schema: public; Owner: -
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
INSERT INTO public.vendor_couriers VALUES ('cdc48954-fe6d-4680-997f-0914bbdc4874', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 1, true, '2026-03-13 07:24:21.226459+00');


--
-- Data for Name: vendor_documents; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_documents VALUES ('a0726b95-ee4e-4ac3-b3e8-1d6ceeb6564b', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 'owner_document_id', 'https://placehold.co/600x400?text=KTP', 'image/jpeg', NULL, NULL, 'e166571c-3c00-488b-9844-ec536a54f332', 'verified', NULL, '7907d953-ba02-40a4-b403-07b68d8471e8', '2026-02-16 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00');
INSERT INTO public.vendor_documents VALUES ('4de3eaf7-647f-4720-8952-0dc47f7c77aa', 'ef0ad9ea-13b6-46ac-9bac-7a81af8ee705', 'owner_document_id', 'vendor-onboardings/0615ee4a-e3cd-4af7-8855-ba6fa88a06a1/documents/owner_document_id/4a6dab7c-78a0-42ed-a183-9c91ac3c1b6a', 'image/png', 3065, NULL, 'ae02bdd5-246c-4a1c-a619-f7eefbbe4db6', 'pending', NULL, NULL, NULL, '2026-03-16 04:14:11.886223+00', '2026-03-16 04:14:11.886223+00');


--
-- Data for Name: vendor_onboardings; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_onboardings VALUES ('fad8946b-1911-4917-bb7c-788a37a66954', 'vendor-seed@dev.local', 'completed', '$2a$10$abcdefghijklmnopqrstuvwxyz1234567890ABCDEF0123456', 'Toko Berkah Haji', 'souvenir_store', 'perorangan', 'ktp', '3201010101900001', 'Vendor Seed Owner', '1990-01-01', NULL, '2026-02-10 07:24:21.226459+00', '2026-02-11 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', NULL, NULL, NULL, NULL, NULL);
INSERT INTO public.vendor_onboardings VALUES ('0615ee4a-e3cd-4af7-8855-ba6fa88a06a1', '21082010249@student.upnjatim.ac.id', 'completed', '$2a$10$gbPSerb5KETGyJCR/BGXK.Fg6L1tmpiCop.TMufeORmKZDooiFnV6', 'Toko Oleh Oleh Haji Test', 'souvenir_store', 'perorangan', 'ktp', '3173000000000001', 'Ahmad Subarkah', '1990-01-02', 'vendor-onboardings/0615ee4a-e3cd-4af7-8855-ba6fa88a06a1/documents/owner_document_id/4a6dab7c-78a0-42ed-a183-9c91ac3c1b6a', '2026-03-16 04:13:27.285894+00', '2026-03-16 04:14:11.886223+00', '2026-03-16 03:41:49.87868+00', '2026-03-16 04:14:11.886223+00', NULL, NULL, NULL, NULL, NULL);


--
-- Data for Name: vendor_withdrawals; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_withdrawals VALUES ('90900fc3-880d-49e3-98ce-534da7bfa040', '8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d', 500000.00, 'ID_BCA', 'completed', NULL, NULL, 'Penarikan saldo vendor', NULL, '2026-03-10 07:24:21.226459+00', '2026-03-10 07:24:21.226459+00', 2500.00, NULL, NULL, 502500.00, 'resolved', NULL);


--
-- Data for Name: wishlist_items; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.wishlist_items VALUES ('f105e107-7125-4d58-9125-9286c3543080', '7d281d37-8f4a-48a8-9b33-e47c76b3d13c', 'b2999c9e-df33-48ef-a109-38c77e51f25d', '2026-03-13 07:24:21.226459+00');


--
-- Name: admin_contacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.admin_contacts_id_seq', 1, true);


--
-- Name: couriers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.couriers_id_seq', 12, true);


--
-- Name: faqs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.faqs_id_seq', 1, true);


--
-- Name: return_reasons_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.return_reasons_id_seq', 1, true);


--
-- Name: roles_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.roles_id_seq', 5, true);


--
-- Name: ticket_subjects_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.ticket_subjects_id_seq', 12, true);


--
-- PostgreSQL database dump complete
--

\unrestrict DB8L0OWFBExbXjJkiFzNg8z38hhHubmGCEwuaIY4bWU4a8WsW8Kw6qhvdOlWAM5

