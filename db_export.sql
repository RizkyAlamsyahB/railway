--
-- PostgreSQL database dump
--

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

INSERT INTO public.roles VALUES (1, 'admin', 'Administrator') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (2, 'umkm', 'UMKM') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (3, 'customer', 'Customer') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (4, 'cs', 'Customer Service') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (5, 'finance', 'Finance') ON CONFLICT DO NOTHING;


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.users VALUES ('2aea50ec-ff64-4f82-b729-c2f4085a5622', '404lamfound@gmail.com', 'Ahmad Subarkah', '1990-01-02', NULL, '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 2, 'active', '2026-03-31 00:37:46.560083+00', '2026-03-31 00:38:29.394538+00', '2026-03-31 00:38:29.394538+00', NULL);
INSERT INTO public.users VALUES ('577c7f1a-8d9d-4c94-9758-920cc72e81d7', 'harundarat@gmail.com', 'System Administrator', NULL, '0811111', '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 1, 'active', '2026-03-31 00:39:54.52106+00', '2026-03-31 00:39:54.52106+00', '2026-03-31 00:39:54.52106+00', NULL);
INSERT INTO public.users VALUES ('d1202987-63e0-4bd5-8087-0ac8190e326a', 'ah.nursufyantsauri@gmail.com', 'Dev Finance', NULL, '081200000004', '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 5, 'active', '2026-03-31 00:39:54.91783+00', '2026-03-31 00:39:54.91783+00', '2026-03-31 00:39:54.91783+00', NULL);
INSERT INTO public.users VALUES ('cf0add84-b149-40c7-ae68-fb9d02594cae', 'rizkyalamsyah.dev@gmail.com', 'Dev Customer', NULL, '081200000006', '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 3, 'active', '2026-03-31 00:39:54.850767+00', '2026-03-31 00:39:54.850767+00', '2026-03-31 00:39:54.850767+00', NULL);
INSERT INTO public.users VALUES ('789da1a5-b9e4-4fa0-aaca-6ed900e31b10', 'adkhawildanrizqia@gmail.com', 'Dev Customer Service 1', NULL, '081200000003', '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 4, 'active', '2026-03-31 00:39:54.719756+00', '2026-03-31 00:39:54.719756+00', '2026-03-31 00:39:54.719756+00', NULL);
INSERT INTO public.users VALUES ('93e8cd6a-75d6-4386-875c-6247b0680a17', 'galangarsandy@gmail.com', 'Dev Customer Service 2', NULL, '081200000005', '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 4, 'active', '2026-03-31 00:39:54.784821+00', '2026-03-31 00:39:54.784821+00', '2026-03-31 00:39:54.784821+00', NULL);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.addresses VALUES ('3658d6f0-78e0-46df-8f9d-691282aa2e54', '2aea50ec-ff64-4f82-b729-c2f4085a5622', 'Alamat Tenggilis', 'Admin warehouse A', '053950154920', '18', '577', '5900', '60292', 'JL. Tenggilis Mejoyo No 99 ', true, 'JAWA TIMUR', ' SURABAYA', 'TENGGILIS MEJOYO', '69350', 'TENGILIS MEJOYO', NULL, NULL, NULL, '2026-03-31 01:15:00.554788+00', '2026-03-31 01:15:00.554788+00');
INSERT INTO public.addresses VALUES ('a2db0a04-f57d-4d03-b4f2-1436a5887137', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'Rumah', 'RIZKY ', '085179733184', '18', '583', '6001', '61253', 'Jl. Raya Sedati Gede No 100', true, 'JAWA TIMUR', ' SIDOARJO', 'SEDATI', '70995', 'SEDATI GEDE', NULL, NULL, NULL, '2026-03-31 01:16:20.300846+00', '2026-03-31 01:16:20.300846+00');


--
-- Data for Name: admin_contacts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.admin_contacts VALUES (1, 'Hubungi admin di admin@hajjstore.id atau WhatsApp 08123456789', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: banners; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.banners VALUES ('fefdc1c3-c086-4497-b4ee-55516306f062', 'Promo Haji 2026', 'banners/fefdc1c3-c086-4497-b4ee-55516306f062/1774922027425', '2026-03-31 01:53:47.425733+00', '2026-03-31 01:55:19.146975+00');


--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.carts VALUES ('e3b41485-0f30-41bc-ad2d-bae4be7dc545', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'converted', '2026-03-31 01:16:25.480648+00', '2026-03-31 01:17:39.17205+00');
INSERT INTO public.carts VALUES ('938ea007-41c2-4b16-bcc1-97ff974b81aa', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'converted', '2026-03-31 01:18:13.726642+00', '2026-03-31 01:18:48.018397+00');
INSERT INTO public.carts VALUES ('b43a934d-889c-4cef-aeb2-d6643d63f346', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'converted', '2026-03-31 01:26:15.781617+00', '2026-03-31 01:26:25.991573+00');
INSERT INTO public.carts VALUES ('fc8523a3-f579-4d10-9440-79e2f031d760', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'converted', '2026-03-31 01:27:49.504713+00', '2026-03-31 01:28:08.060794+00');
INSERT INTO public.carts VALUES ('b4164544-5b99-4edd-8501-6b8f289d8bd1', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'converted', '2026-03-31 01:29:18.165575+00', '2026-03-31 01:29:30.59768+00');


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: -
--

-- Parent categories first (self-referencing)
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000001', 'ca000001-0000-0000-0000-000000000001', 'Perlengkapan Haji & Umrah', 'perlengkapan-haji-umrah', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000002', 'ca000001-0000-0000-0000-000000000002', 'Oleh-oleh & Souvenir', 'oleh-oleh-souvenir', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000003', 'ca000001-0000-0000-0000-000000000003', 'Fashion Muslim', 'fashion-muslim', true);
-- Child categories
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

INSERT INTO public.vendors VALUES ('297a0ab4-2168-4d54-9bf1-f526a476ac19', '2aea50ec-ff64-4f82-b729-c2f4085a5622', 'souvenir_store', NULL, 'Toko Oleh Oleh Haji Test', NULL, 'active', '577c7f1a-8d9d-4c94-9758-920cc72e81d7', '2026-03-31 00:45:53.171587+00', NULL, '69cb1940dd22af33b6ce22b5', '2026-03-31 00:38:29.394538+00', '2026-03-31 00:45:53.171587+00', 'perorangan', NULL);


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.products VALUES ('8051032a-f0a6-4c6f-8893-2429dcda5542', '297a0ab4-2168-4d54-9bf1-f526a476ac19', 'ca000001-0000-0000-0000-000000000021', 'Kurma Ajwa', 'kurma-ajwa', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 'published', 'pending', NULL, '2026-03-31 01:09:21.769077+00', '2026-03-31 01:09:21.769077+00');


--
-- Data for Name: product_variants; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_variants VALUES ('1b81cec0-bc14-4610-8237-c6daba3f90c8', '8051032a-f0a6-4c6f-8893-2429dcda5542', 'TOKOO-KURMA000-00', 'Default', 100000.00, 'IDR', 5, 500, true, true);


--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.cart_items VALUES ('96ff4561-72ca-499c-9900-086965aa1640', 'e3b41485-0f30-41bc-ad2d-bae4be7dc545', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 1, '2026-03-31 01:16:25.490708+00');
INSERT INTO public.cart_items VALUES ('8b3b8382-9001-44c6-b8f7-2bf13596267c', '938ea007-41c2-4b16-bcc1-97ff974b81aa', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 1, '2026-03-31 01:18:13.741329+00');
INSERT INTO public.cart_items VALUES ('fc92f9c4-afe3-45ea-9e0c-c1eb6649a967', 'b43a934d-889c-4cef-aeb2-d6643d63f346', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 1, '2026-03-31 01:26:15.794892+00');
INSERT INTO public.cart_items VALUES ('7f29a240-1d1b-4938-bc85-9efb97077a1c', 'fc8523a3-f579-4d10-9440-79e2f031d760', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 1, '2026-03-31 01:27:49.517334+00');
INSERT INTO public.cart_items VALUES ('da9e6a2e-356a-4074-b0ab-b678dcea30ab', 'b4164544-5b99-4edd-8501-6b8f289d8bd1', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 1, '2026-03-31 01:29:18.177049+00');


--
-- Data for Name: chat_conversations; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: couriers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.couriers VALUES (1, 'jne', 'JNE', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (2, 'sicepat', 'SiCepat', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (3, 'ide', 'IDExpress', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (4, 'sap', 'SAP Express', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (5, 'ninja', 'Ninja', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (6, 'jnt', 'J&T Express', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (7, 'tiki', 'TIKI', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (8, 'wahana', 'Wahana Express', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (9, 'pos', 'POS Indonesia', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (10, 'sentral', 'Sentral Cargo', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (11, 'lion', 'Lion Parcel', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (12, 'rex', 'Royal Express Asia', NULL, true, '2026-03-31 00:34:10.091535+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: email_verification_tokens; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: faqs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.faqs VALUES (1, 'Umum', 'Bagaimana cara memesan perlengkapan haji?', 'Anda bisa memesan melalui aplikasi atau website kami. Pilih produk, masukkan ke keranjang, lalu lakukan checkout.', '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: ledger_accounts; Type: TABLE DATA; Schema: public; Owner: -
-- (delete existing seeded data first, then re-insert with correct UUIDs)
--

DELETE FROM public.ledger_lines;
DELETE FROM public.ledger_journals;
DELETE FROM public.ledger_accounts;
INSERT INTO public.ledger_accounts VALUES ('461bdec8-e30f-4fc0-8613-254a5cf9ed1b', '1100', 'Payment Gateway Receivable', 'asset', 'D', true);
INSERT INTO public.ledger_accounts VALUES ('8a321b91-52a2-4574-b4b1-adaf5b4ee045', '2100', 'Vendor Payable', 'liability', 'C', true);
INSERT INTO public.ledger_accounts VALUES ('c7fa5483-7860-469a-b52e-da8f12bf331f', '4100', 'Platform Fee Revenue', 'revenue', 'C', true);
INSERT INTO public.ledger_accounts VALUES ('52a94155-2924-42bf-a426-e7ad56d66d82', '4200', 'Admin Fee Revenue', 'revenue', 'C', true);


--
-- Data for Name: ledger_journals; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ledger_journals VALUES ('831af351-5f0e-4647-b8b1-0297cdd09b10', 'JRN-PAY-ORD-20260331-0242A9B9-a428b23f', 'payment_invoice', '0327d632-8721-4742-9eea-63d1d24bef5f', '2026-03-31 01:19:03.180985+00', 'Payment received for order ORD-20260331-0242A9B9', 'posted', NULL, '2026-03-31 01:19:03.180985+00');
INSERT INTO public.ledger_journals VALUES ('63af8577-c25c-4504-a28b-5715c50881dd', 'JRN-PAY-ORD-20260331-21CC195D-636e6a19', 'payment_invoice', 'd2c66430-3f5a-4b26-924b-bebf22f8a431', '2026-03-31 01:26:35.397754+00', 'Payment received for order ORD-20260331-21CC195D', 'posted', NULL, '2026-03-31 01:26:35.397754+00');
INSERT INTO public.ledger_journals VALUES ('2974da14-d81c-4b84-8a5b-2432f83201d6', 'JRN-PAY-ORD-20260331-A36452DC-e5345fe5', 'payment_invoice', 'bd49105a-2bec-4b9a-8e67-4d8a16a255b6', '2026-03-31 01:28:18.105471+00', 'Payment received for order ORD-20260331-A36452DC', 'posted', NULL, '2026-03-31 01:28:18.105471+00');
INSERT INTO public.ledger_journals VALUES ('d2566270-b4fe-4e85-9927-602ee610efd1', 'JRN-PAY-ORD-20260331-84E3D283-9f7fbe0e', 'payment_invoice', 'ff4b8a75-e87f-4325-b6f5-68a26ba67187', '2026-03-31 01:29:39.194587+00', 'Payment received for order ORD-20260331-84E3D283', 'posted', NULL, '2026-03-31 01:29:39.194587+00');


--
-- Data for Name: ledger_lines; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ledger_lines VALUES ('82787462-9cb4-477a-b088-b216a1fdc337', '831af351-5f0e-4647-b8b1-0297cdd09b10', '461bdec8-e30f-4fc0-8613-254a5cf9ed1b', 146000.00, 0.00, 'IDR', 'ORD-20260331-0242A9B9');
INSERT INTO public.ledger_lines VALUES ('4ec7a2a5-8b36-4c55-8fe9-ae245e473305', '831af351-5f0e-4647-b8b1-0297cdd09b10', '8a321b91-52a2-4574-b4b1-adaf5b4ee045', 0.00, 100000.00, 'IDR', 'ORD-20260331-0242A9B9');
INSERT INTO public.ledger_lines VALUES ('6b566e32-de64-4c0f-b91b-a4eba7dafbbb', '831af351-5f0e-4647-b8b1-0297cdd09b10', 'c7fa5483-7860-469a-b52e-da8f12bf331f', 0.00, 1000.00, 'IDR', 'ORD-20260331-0242A9B9');
INSERT INTO public.ledger_lines VALUES ('ebd82f3c-4a59-4709-9bcd-7c8394f4b572', '831af351-5f0e-4647-b8b1-0297cdd09b10', '52a94155-2924-42bf-a426-e7ad56d66d82', 0.00, 5000.00, 'IDR', 'ORD-20260331-0242A9B9');
INSERT INTO public.ledger_lines VALUES ('b514e430-7ead-4a06-8095-963861b51a9a', '63af8577-c25c-4504-a28b-5715c50881dd', '461bdec8-e30f-4fc0-8613-254a5cf9ed1b', 113000.00, 0.00, 'IDR', 'ORD-20260331-21CC195D');
INSERT INTO public.ledger_lines VALUES ('3a1534f5-53b5-462c-aab1-187e3c789967', '63af8577-c25c-4504-a28b-5715c50881dd', '8a321b91-52a2-4574-b4b1-adaf5b4ee045', 0.00, 100000.00, 'IDR', 'ORD-20260331-21CC195D');
INSERT INTO public.ledger_lines VALUES ('6a779320-80eb-4a20-b447-b8936f67fb8f', '63af8577-c25c-4504-a28b-5715c50881dd', 'c7fa5483-7860-469a-b52e-da8f12bf331f', 0.00, 1000.00, 'IDR', 'ORD-20260331-21CC195D');
INSERT INTO public.ledger_lines VALUES ('126521f0-5eff-4287-b415-3728c640642e', '63af8577-c25c-4504-a28b-5715c50881dd', '52a94155-2924-42bf-a426-e7ad56d66d82', 0.00, 5000.00, 'IDR', 'ORD-20260331-21CC195D');
INSERT INTO public.ledger_lines VALUES ('91fb9c7b-a895-4cc5-acc8-c31ece6cfd24', '2974da14-d81c-4b84-8a5b-2432f83201d6', '461bdec8-e30f-4fc0-8613-254a5cf9ed1b', 113000.00, 0.00, 'IDR', 'ORD-20260331-A36452DC');
INSERT INTO public.ledger_lines VALUES ('b358c7df-f12e-4c27-9827-665d824160ea', '2974da14-d81c-4b84-8a5b-2432f83201d6', '8a321b91-52a2-4574-b4b1-adaf5b4ee045', 0.00, 100000.00, 'IDR', 'ORD-20260331-A36452DC');
INSERT INTO public.ledger_lines VALUES ('9d8d0cad-c2fa-4185-9efa-13944ee0426e', '2974da14-d81c-4b84-8a5b-2432f83201d6', 'c7fa5483-7860-469a-b52e-da8f12bf331f', 0.00, 1000.00, 'IDR', 'ORD-20260331-A36452DC');
INSERT INTO public.ledger_lines VALUES ('86dba8ff-29a8-414b-81c8-eed5e4ec2cc8', '2974da14-d81c-4b84-8a5b-2432f83201d6', '52a94155-2924-42bf-a426-e7ad56d66d82', 0.00, 5000.00, 'IDR', 'ORD-20260331-A36452DC');
INSERT INTO public.ledger_lines VALUES ('6ae71c26-fbf6-4284-831b-10d57c553726', 'd2566270-b4fe-4e85-9927-602ee610efd1', '461bdec8-e30f-4fc0-8613-254a5cf9ed1b', 113000.00, 0.00, 'IDR', 'ORD-20260331-84E3D283');
INSERT INTO public.ledger_lines VALUES ('364573da-253a-4d2d-879c-dc1e04f8490b', 'd2566270-b4fe-4e85-9927-602ee610efd1', '8a321b91-52a2-4574-b4b1-adaf5b4ee045', 0.00, 100000.00, 'IDR', 'ORD-20260331-84E3D283');
INSERT INTO public.ledger_lines VALUES ('8f25f3d4-9ae1-4d15-aaa0-ce28214988c1', 'd2566270-b4fe-4e85-9927-602ee610efd1', 'c7fa5483-7860-469a-b52e-da8f12bf331f', 0.00, 1000.00, 'IDR', 'ORD-20260331-84E3D283');
INSERT INTO public.ledger_lines VALUES ('b6b3b242-e63a-4155-b0ac-f7a111f11247', 'd2566270-b4fe-4e85-9927-602ee610efd1', '52a94155-2924-42bf-a426-e7ad56d66d82', 0.00, 5000.00, 'IDR', 'ORD-20260331-84E3D283');


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.orders VALUES ('d31dc0b1-0079-4f43-971a-d0f86fe64cfb', 'ORD-20260331-3744013F', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '297a0ab4-2168-4d54-9bf1-f526a476ac19', '{"label": "Rumah", "phone": "085179733184", "city_name": " SIDOARJO", "address_id": "a2db0a04-f57d-4d03-b4f2-1436a5887137", "is_default": true, "postal_code": "61253", "address_line": "Jl. Raya Sedati Gede No 100", "district_name": "SEDATI", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY ", "subdistrict_name": "SEDATI GEDE"}', 'pending_payment', 'unpaid', 100000.00, 7000.00, 6000.00, 113000.00, '2026-03-31 01:17:35.778637+00', '2026-03-31 01:17:35.778637+00', '2026-03-31 01:17:35.778637+00');
INSERT INTO public.orders VALUES ('e0ee921c-9ed6-41e7-b1c4-ae3c55fc6280', 'ORD-20260331-21CC195D', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '297a0ab4-2168-4d54-9bf1-f526a476ac19', '{"label": "Rumah", "phone": "085179733184", "city_name": " SIDOARJO", "address_id": "a2db0a04-f57d-4d03-b4f2-1436a5887137", "is_default": true, "postal_code": "61253", "address_line": "Jl. Raya Sedati Gede No 100", "district_name": "SEDATI", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY ", "subdistrict_name": "SEDATI GEDE"}', 'paid', 'paid', 100000.00, 7000.00, 6000.00, 113000.00, '2026-03-31 01:26:24.900206+00', '2026-03-31 01:26:24.900206+00', '2026-03-31 01:26:35.395992+00');
INSERT INTO public.orders VALUES ('4605d7ba-5a51-49d1-b793-29d82de1099d', 'ORD-20260331-0242A9B9', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '297a0ab4-2168-4d54-9bf1-f526a476ac19', '{"label": "Rumah", "phone": "085179733184", "city_name": " SIDOARJO", "address_id": "a2db0a04-f57d-4d03-b4f2-1436a5887137", "is_default": true, "postal_code": "61253", "address_line": "Jl. Raya Sedati Gede No 100", "district_name": "SEDATI", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY ", "subdistrict_name": "SEDATI GEDE"}', 'processing', 'paid', 100000.00, 40000.00, 6000.00, 146000.00, '2026-03-31 01:18:47.418869+00', '2026-03-31 01:18:47.418869+00', '2026-03-31 01:27:25.651314+00');
INSERT INTO public.orders VALUES ('3c6b4adf-7643-4336-9088-a6ba755e30d4', 'ORD-20260331-A36452DC', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '297a0ab4-2168-4d54-9bf1-f526a476ac19', '{"label": "Rumah", "phone": "085179733184", "city_name": " SIDOARJO", "address_id": "a2db0a04-f57d-4d03-b4f2-1436a5887137", "is_default": true, "postal_code": "61253", "address_line": "Jl. Raya Sedati Gede No 100", "district_name": "SEDATI", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY ", "subdistrict_name": "SEDATI GEDE"}', 'shipped', 'paid', 100000.00, 7000.00, 6000.00, 113000.00, '2026-03-31 01:28:06.922757+00', '2026-03-31 01:28:06.922757+00', '2026-03-31 01:28:47.513426+00');
INSERT INTO public.orders VALUES ('959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'ORD-20260331-84E3D283', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '297a0ab4-2168-4d54-9bf1-f526a476ac19', '{"label": "Rumah", "phone": "085179733184", "city_name": " SIDOARJO", "address_id": "a2db0a04-f57d-4d03-b4f2-1436a5887137", "is_default": true, "postal_code": "61253", "address_line": "Jl. Raya Sedati Gede No 100", "district_name": "SEDATI", "province_name": "JAWA TIMUR", "recipient_name": "RIZKY ", "subdistrict_name": "SEDATI GEDE"}', 'completed', 'paid', 100000.00, 7000.00, 6000.00, 113000.00, '2026-03-31 01:29:30.018078+00', '2026-03-31 01:29:30.018078+00', '2026-03-31 01:38:34.011887+00');


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.order_items VALUES ('776d416d-ac28-4998-8c60-6c3dc0487ac3', 'd31dc0b1-0079-4f43-971a-d0f86fe64cfb', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 'Kurma Ajwa', 'TOKOO-KURMA000-00', 1, 100000.00, 100000.00);
INSERT INTO public.order_items VALUES ('9592450c-6a44-47d1-999d-84345394295e', '4605d7ba-5a51-49d1-b793-29d82de1099d', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 'Kurma Ajwa', 'TOKOO-KURMA000-00', 1, 100000.00, 100000.00);
INSERT INTO public.order_items VALUES ('8749e151-ad4e-4ffd-9e1e-321c6448fe97', 'e0ee921c-9ed6-41e7-b1c4-ae3c55fc6280', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 'Kurma Ajwa', 'TOKOO-KURMA000-00', 1, 100000.00, 100000.00);
INSERT INTO public.order_items VALUES ('79f6915c-8ce7-4c7a-8223-9f568ae07ada', '3c6b4adf-7643-4336-9088-a6ba755e30d4', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 'Kurma Ajwa', 'TOKOO-KURMA000-00', 1, 100000.00, 100000.00);
INSERT INTO public.order_items VALUES ('ada434c2-5887-4e5c-bcf6-9e48ca5226d4', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', '1b81cec0-bc14-4610-8237-c6daba3f90c8', 'Kurma Ajwa', 'TOKOO-KURMA000-00', 1, 100000.00, 100000.00);


--
-- Data for Name: order_status_history; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.order_status_history VALUES ('80542489-edfb-4dcf-addd-b4f4b3e0dcb9', 'd31dc0b1-0079-4f43-971a-d0f86fe64cfb', NULL, 'pending_payment', NULL, '2026-03-31 01:17:35.862975+00', NULL);
INSERT INTO public.order_status_history VALUES ('989bd992-6b31-49b2-b27f-a7c0413a3007', '4605d7ba-5a51-49d1-b793-29d82de1099d', NULL, 'pending_payment', NULL, '2026-03-31 01:18:47.508046+00', NULL);
INSERT INTO public.order_status_history VALUES ('d3bac92e-05e5-422f-9fd9-b40a7c1c9781', '4605d7ba-5a51-49d1-b793-29d82de1099d', 'pending_payment', 'paid', NULL, '2026-03-31 01:19:03.178425+00', 'Payment received via BANK_TRANSFER/BCA');
INSERT INTO public.order_status_history VALUES ('e46519b8-7aac-46f3-8d5a-25e734b87436', 'e0ee921c-9ed6-41e7-b1c4-ae3c55fc6280', NULL, 'pending_payment', NULL, '2026-03-31 01:26:24.977129+00', NULL);
INSERT INTO public.order_status_history VALUES ('bae6a6ef-b7e8-4143-a218-dc95f6f7a0fa', 'e0ee921c-9ed6-41e7-b1c4-ae3c55fc6280', 'pending_payment', 'paid', NULL, '2026-03-31 01:26:35.396233+00', 'Payment received via BANK_TRANSFER/BNI');
INSERT INTO public.order_status_history VALUES ('c4a89b3c-8fdb-4304-bd28-77b7cd8d4842', '4605d7ba-5a51-49d1-b793-29d82de1099d', 'paid', 'processing', '2aea50ec-ff64-4f82-b729-c2f4085a5622', '2026-03-31 01:27:25.65231+00', 'Order accepted by vendor');
INSERT INTO public.order_status_history VALUES ('99208ac5-9160-491b-8846-43a6624c0fd4', '3c6b4adf-7643-4336-9088-a6ba755e30d4', NULL, 'pending_payment', NULL, '2026-03-31 01:28:06.997889+00', NULL);
INSERT INTO public.order_status_history VALUES ('618377ca-aa00-418b-840e-c91810c9aaf8', '3c6b4adf-7643-4336-9088-a6ba755e30d4', 'pending_payment', 'paid', NULL, '2026-03-31 01:28:18.102747+00', 'Payment received via BANK_TRANSFER/CIMB');
INSERT INTO public.order_status_history VALUES ('49026b50-d40d-47fe-822c-c5ec2eab4dff', '3c6b4adf-7643-4336-9088-a6ba755e30d4', 'paid', 'processing', '2aea50ec-ff64-4f82-b729-c2f4085a5622', '2026-03-31 01:28:35.563615+00', 'Order accepted by vendor');
INSERT INTO public.order_status_history VALUES ('9a935ceb-5621-47b1-a7ef-7f2ba51ebd1c', '3c6b4adf-7643-4336-9088-a6ba755e30d4', 'processing', 'shipped', '2aea50ec-ff64-4f82-b729-c2f4085a5622', '2026-03-31 01:28:47.514354+00', 'Order shipped by vendor, tracking: JX7629871362 (jnt EZ)');
INSERT INTO public.order_status_history VALUES ('0af1d3f4-07c3-4f60-83d4-6ab8f536c89d', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', NULL, 'pending_payment', NULL, '2026-03-31 01:29:30.098581+00', NULL);
INSERT INTO public.order_status_history VALUES ('94d07ed2-9283-42e8-b50c-7948bf90e447', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'pending_payment', 'paid', NULL, '2026-03-31 01:29:39.192103+00', 'Payment received via BANK_TRANSFER/BRI');
INSERT INTO public.order_status_history VALUES ('d87446b4-04d9-4603-bbab-88b6b45d5573', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'paid', 'processing', '2aea50ec-ff64-4f82-b729-c2f4085a5622', '2026-03-31 01:29:55.814853+00', 'Order accepted by vendor');
INSERT INTO public.order_status_history VALUES ('4ba922c1-b83d-4968-9317-b3d2658748ed', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'processing', 'shipped', '2aea50ec-ff64-4f82-b729-c2f4085a5622', '2026-03-31 01:30:04.966062+00', 'Order shipped by vendor, tracking: JX7629871362 (jnt EZ)');
INSERT INTO public.order_status_history VALUES ('af186b51-20c5-442b-9e10-2c9886c8acd3', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'shipped', 'received', NULL, '2026-03-31 01:30:07.262257+00', 'Auto-updated: courier confirmed delivery');
INSERT INTO public.order_status_history VALUES ('8c1beb00-835e-49f9-b3fe-f994b399499f', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'received', 'completed', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '2026-03-31 01:38:34.011887+00', 'Order completed by customer');


--
-- Data for Name: otp_codes; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.otp_codes VALUES ('e3343240-2269-450a-a3ae-b824db10a342', NULL, '404lamfound@gmail.com', 'email_verification', 'email', '733ugJuTiSfbCJzmgmUthwBSgY6uzzQ-BNHMB0nB_d8', '2026-03-31 00:41:41.275797+00', 0, '2026-03-31 00:37:46.55652+00', NULL, '2026-03-31 00:36:41.275797+00');


--
-- Data for Name: payment_invoices; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.payment_invoices VALUES ('6103e03b-02f2-45dc-8eef-da0625291e9e', 'd31dc0b1-0079-4f43-971a-d0f86fe64cfb', 'xendit', '69cb20b07cba7679601590f5', 'INV-ORD-20260331-3744013F-6ef8357b', 'https://checkout-staging.xendit.co/web/69cb20b07cba7679601590f5', NULL, NULL, 113000.00, 'IDR', 'pending', '2026-04-01 01:17:39.004+00', NULL, NULL, '2026-03-31 01:17:35.778637+00', '2026-03-31 01:17:35.778637+00');
INSERT INTO public.payment_invoices VALUES ('0327d632-8721-4742-9eea-63d1d24bef5f', '4605d7ba-5a51-49d1-b793-29d82de1099d', 'xendit', '69cb20f77cba76796015918e', 'INV-ORD-20260331-0242A9B9-dcfc237f', 'https://checkout-staging.xendit.co/web/69cb20f77cba76796015918e', 'BANK_TRANSFER', 'BCA', 146000.00, 'IDR', 'paid', '2026-04-01 01:18:47.654+00', '2026-03-31 01:19:02+00', '{"webhook_payload": {"id": "69cb20f77cba76796015918e", "amount": 146000, "status": "PAID", "paid_at": "2026-03-31T01:19:02.000Z", "user_id": "69cb1940dd22af33b6ce22b5", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260331-0242A9B9-dcfc237f", "paid_amount": 146000, "payment_method": "BANK_TRANSFER", "payment_channel": "BCA"}}', '2026-03-31 01:18:47.418869+00', '2026-03-31 01:19:03.174523+00');
INSERT INTO public.payment_invoices VALUES ('d2c66430-3f5a-4b26-924b-bebf22f8a431', 'e0ee921c-9ed6-41e7-b1c4-ae3c55fc6280', 'xendit', '69cb22c17cba767960159422', 'INV-ORD-20260331-21CC195D-9f55a414', 'https://checkout-staging.xendit.co/web/69cb22c17cba767960159422', 'BANK_TRANSFER', 'BNI', 113000.00, 'IDR', 'paid', '2026-04-01 01:26:25.711+00', '2026-03-31 01:26:34+00', '{"webhook_payload": {"id": "69cb22c17cba767960159422", "amount": 113000, "status": "PAID", "paid_at": "2026-03-31T01:26:34.000Z", "user_id": "69cb1940dd22af33b6ce22b5", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260331-21CC195D-9f55a414", "paid_amount": 113000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}}', '2026-03-31 01:26:24.900206+00', '2026-03-31 01:26:35.394199+00');
INSERT INTO public.payment_invoices VALUES ('bd49105a-2bec-4b9a-8e67-4d8a16a255b6', '3c6b4adf-7643-4336-9088-a6ba755e30d4', 'xendit', '69cb2327eea2af3427b89c5a', 'INV-ORD-20260331-A36452DC-aebf604f', 'https://checkout-staging.xendit.co/web/69cb2327eea2af3427b89c5a', 'BANK_TRANSFER', 'CIMB', 113000.00, 'IDR', 'paid', '2026-04-01 01:28:07.607+00', '2026-03-31 01:28:17+00', '{"webhook_payload": {"id": "69cb2327eea2af3427b89c5a", "amount": 113000, "status": "PAID", "paid_at": "2026-03-31T01:28:17.000Z", "user_id": "69cb1940dd22af33b6ce22b5", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260331-A36452DC-aebf604f", "paid_amount": 113000, "payment_method": "BANK_TRANSFER", "payment_channel": "CIMB"}}', '2026-03-31 01:28:06.922757+00', '2026-03-31 01:28:18.098738+00');
INSERT INTO public.payment_invoices VALUES ('ff4b8a75-e87f-4325-b6f5-68a26ba67187', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'xendit', '69cb237aeea2af3427b89ca8', 'INV-ORD-20260331-84E3D283-b4c10887', 'https://checkout-staging.xendit.co/web/69cb237aeea2af3427b89ca8', 'BANK_TRANSFER', 'BRI', 113000.00, 'IDR', 'paid', '2026-04-01 01:29:30.245+00', '2026-03-31 01:29:38+00', '{"webhook_payload": {"id": "69cb237aeea2af3427b89ca8", "amount": 113000, "status": "PAID", "paid_at": "2026-03-31T01:29:38.000Z", "user_id": "69cb1940dd22af33b6ce22b5", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260331-84E3D283-b4c10887", "paid_amount": 113000, "payment_method": "BANK_TRANSFER", "payment_channel": "BRI"}}', '2026-03-31 01:29:30.018078+00', '2026-03-31 01:29:39.187119+00');


--
-- Data for Name: payment_events; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.payment_events VALUES ('75be0502-3157-4976-9851-e8ad0169bfa7', '0327d632-8721-4742-9eea-63d1d24bef5f', 'invoice.paid', '69cb20f77cba76796015918e:PAID', '{"id": "69cb20f77cba76796015918e", "amount": 146000, "status": "PAID", "paid_at": "2026-03-31T01:19:02.000Z", "external_id": "INV-ORD-20260331-0242A9B9-dcfc237f", "paid_amount": 146000, "payment_method": "BANK_TRANSFER", "payment_channel": "BCA"}', '2026-03-31 01:19:03.165997+00');
INSERT INTO public.payment_events VALUES ('01b7e9fd-d6ab-4022-882f-a970d3e7e98e', 'd2c66430-3f5a-4b26-924b-bebf22f8a431', 'invoice.paid', '69cb22c17cba767960159422:PAID', '{"id": "69cb22c17cba767960159422", "amount": 113000, "status": "PAID", "paid_at": "2026-03-31T01:26:34.000Z", "external_id": "INV-ORD-20260331-21CC195D-9f55a414", "paid_amount": 113000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}', '2026-03-31 01:26:35.382228+00');
INSERT INTO public.payment_events VALUES ('e40fd520-ccf6-4f16-8797-8e6311b56307', 'bd49105a-2bec-4b9a-8e67-4d8a16a255b6', 'invoice.paid', '69cb2327eea2af3427b89c5a:PAID', '{"id": "69cb2327eea2af3427b89c5a", "amount": 113000, "status": "PAID", "paid_at": "2026-03-31T01:28:17.000Z", "external_id": "INV-ORD-20260331-A36452DC-aebf604f", "paid_amount": 113000, "payment_method": "BANK_TRANSFER", "payment_channel": "CIMB"}', '2026-03-31 01:28:18.085344+00');
INSERT INTO public.payment_events VALUES ('22cc0e41-f522-4d62-8a55-117340946c82', 'ff4b8a75-e87f-4325-b6f5-68a26ba67187', 'invoice.paid', '69cb237aeea2af3427b89ca8:PAID', '{"id": "69cb237aeea2af3427b89ca8", "amount": 113000, "status": "PAID", "paid_at": "2026-03-31T01:29:38.000Z", "external_id": "INV-ORD-20260331-84E3D283-b4c10887", "paid_amount": 113000, "payment_method": "BANK_TRANSFER", "payment_channel": "BRI"}', '2026-03-31 01:29:39.181725+00');


--
-- Data for Name: payout_batches; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: payout_items; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.product_images VALUES ('5aa11b32-25c7-4647-8a31-5f5c0a863582', '8051032a-f0a6-4c6f-8893-2429dcda5542', 'products/8051032a-f0a6-4c6f-8893-2429dcda5542/images/5aa11b32-25c7-4647-8a31-5f5c0a863582/product-main.jpg', 'image/jpeg', 11930, true, 0, '2026-03-31 01:09:21.769077+00');


--
-- Data for Name: product_promotions; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: product_reviews; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: product_review_images; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: product_review_stats; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: refunds; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: reply_templates; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.reply_templates VALUES ('c0e89a64-456f-4df8-bf77-2b8acd4b1b25', 'Salam Pembuka', 'Assalamualaikum, terima kasih telah menghubungi kami. Ada yang bisa kami bantu?', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '/salam', 'general');
INSERT INTO public.reply_templates VALUES ('db924680-cc08-4bf0-9a4c-d54605ff89c1', 'Selamat Datang di UmrahMart', 'Assalamu''alaikum, selamat datang di UmrahMart! 🕌

Kami hadir untuk memudahkan Anda mendapatkan perlengkapan haji dan umroh terbaik. Ada yang bisa kami bantu hari ini?', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/selamat-datang', 'general');
INSERT INTO public.reply_templates VALUES ('b7ec135c-0073-43ec-a7ee-65fecb5668b8', 'Terima Kasih Telah Menghubungi Kami', 'Terima kasih telah menghubungi tim Customer Service UmrahMart. 🙏

Kami akan segera membantu Anda. Mohon tunggu sebentar ya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/terima-kasih-menghubungi', 'general');
INSERT INTO public.reply_templates VALUES ('9950affe-68af-4a87-bddd-57cb2b2d4522', 'Apakah Masih Ada Yang Bisa Dibantu?', 'Apakah masih ada yang bisa kami bantu? 😊

Jika masalah Anda sudah teratasi, jangan lupa berikan ulasan untuk membantu kami meningkatkan layanan. Semoga ibadah haji/umroh Anda berjalan lancar. Aamiin.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/ada-yang-bisa-dibantu', 'general');
INSERT INTO public.reply_templates VALUES ('2210ed2d-6243-47da-af10-cd789f2416b0', 'Jam Operasional Customer Service', 'Tim Customer Service UmrahMart melayani Anda setiap hari:

🕐 Senin – Jumat : 08.00 – 21.00 WIB
🕐 Sabtu – Minggu : 09.00 – 18.00 WIB

Di luar jam operasional, Anda tetap bisa meninggalkan pesan dan kami akan merespons saat jam kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/jam-operasional', 'general');
INSERT INTO public.reply_templates VALUES ('fd06b343-9d5a-49bf-ad6a-427fa4b238b0', 'Konfirmasi Pesanan Diterima', 'Alhamdulillah, pesanan Anda telah kami terima! ✅

Nomor pesanan Anda: {order_number}
Status: Menunggu Pembayaran

Silakan selesaikan pembayaran sebelum {expired_at} agar pesanan dapat segera kami proses. Jazakallahu khairan.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/konfirmasi-pesanan-diterima', 'order');
INSERT INTO public.reply_templates VALUES ('2c22e30d-7e7a-484f-ab6e-bd3b6f90dd29', 'Pesanan Sedang Diproses', 'Kabar baik! Pesanan Anda dengan nomor {order_number} sedang kami proses. 📦

Estimasi pesanan siap dikirim: 1–2 hari kerja.
Kami akan menginformasikan nomor resi pengiriman segera setelah paket diserahkan ke kurir.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/pesanan-diproses', 'order');
INSERT INTO public.reply_templates VALUES ('0dcf4408-01fb-49de-b782-03c949175719', 'Pesanan Berhasil Dibatalkan', 'Pesanan Anda dengan nomor {order_number} telah berhasil dibatalkan.

Jika Anda sudah melakukan pembayaran, proses refund akan kami lakukan dalam 3–7 hari kerja ke metode pembayaran asal.

Apabila ada pertanyaan lebih lanjut, jangan ragu untuk menghubungi kami kembali.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/pesanan-dibatalkan', 'order');
INSERT INTO public.reply_templates VALUES ('09e4155c-ccaf-4466-99e4-6f09c12682c8', 'Detail Pesanan', 'Berikut detail pesanan Anda:

📋 Nomor Pesanan : {order_number}
📅 Tanggal       : {order_date}
💰 Total         : {total_amount}
📍 Status        : {order_status}

Untuk informasi lebih lengkap, silakan cek halaman "Pesanan Saya" di aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/detail-pesanan', 'order');
INSERT INTO public.reply_templates VALUES ('039e2bf9-5a35-47fa-b5c6-df6550e1c4bd', 'Pembayaran Berhasil Diterima', 'Alhamdulillah, pembayaran Anda telah berhasil kami terima! 💚

Nomor pesanan : {order_number}
Jumlah        : {amount}
Metode        : {payment_method}

Pesanan Anda akan segera kami proses. Terima kasih.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/pembayaran-berhasil', 'payment');
INSERT INTO public.reply_templates VALUES ('279010e0-341c-42b1-8818-bf5e5b0d2bea', 'Menunggu Konfirmasi Pembayaran', 'Halo, kami melihat pesanan Anda belum terkonfirmasi pembayarannya.

Batas waktu pembayaran: {expired_at}

Jika Anda sudah melakukan transfer manual, mohon upload bukti pembayaran melalui aplikasi agar kami dapat segera memprosesnya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/menunggu-konfirmasi', 'payment');
INSERT INTO public.reply_templates VALUES ('c5c27b1c-05d5-4ae6-83a0-7bb1c89c74f9', 'Pembayaran Gagal / Kedaluwarsa', 'Kami informasikan bahwa pembayaran untuk pesanan {order_number} telah gagal atau melewati batas waktu. ⚠️

Pesanan Anda telah kami batalkan secara otomatis. Anda dapat melakukan pemesanan ulang kapan saja.

Mohon pastikan saldo/limit mencukupi saat melakukan pembayaran berikutnya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/pembayaran-gagal', 'payment');
INSERT INTO public.reply_templates VALUES ('98ccc880-d214-4c26-bb39-c6fc1cc4bfe3', 'Cara Melakukan Pembayaran', 'Berikut metode pembayaran yang tersedia di UmrahMart:

💳 Transfer Bank (BCA, Mandiri, BNI, BRI)
📱 E-Wallet (GoPay, OVO, DANA, ShopeePay)
🏪 Gerai Retail (Alfamart, Indomaret)
💵 QRIS

Pilih metode yang paling nyaman untuk Anda saat checkout.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/cara-pembayaran', 'payment');
INSERT INTO public.reply_templates VALUES ('022b8cc0-8129-45d4-a9b5-0601ce941461', 'Pengajuan Refund Diterima', 'Pengajuan refund Anda telah kami terima dan sedang dalam proses review. ✅

Nomor tiket refund: {ticket_number}

Tim kami akan memverifikasi dalam 1–3 hari kerja. Anda akan mendapat notifikasi setelah proses selesai.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/pengajuan-diterima', 'refund');
INSERT INTO public.reply_templates VALUES ('60cf1bc4-50e9-49ca-a6fa-dc79ac495678', 'Refund Sedang Diproses', 'Refund Anda sedang dalam proses pencairan. 🔄

Jumlah refund : {refund_amount}
Tujuan        : {refund_destination}
Estimasi      : 3–7 hari kerja

Mohon bersabar dan pastikan rekening/e-wallet tujuan masih aktif.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/sedang-diproses', 'refund');
INSERT INTO public.reply_templates VALUES ('83ec97f6-b6fc-4e70-9491-ac41192fe241', 'Refund Berhasil', 'Dana refund sebesar {refund_amount} telah berhasil kami kirimkan ke {refund_destination}. 🎉

Jika dalam 1×24 jam dana belum masuk, mohon hubungi kami kembali dengan menyertakan nomor tiket {ticket_number}.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/berhasil', 'refund');
INSERT INTO public.reply_templates VALUES ('c0aced63-92c1-4857-8db0-91187f1822ad', 'Pesanan Telah Dikirim', 'Pesanan Anda sudah dalam perjalanan! 🚚

Nomor Pesanan : {order_number}
Kurir         : {courier_name}
Nomor Resi    : {tracking_number}

Lacak pengiriman di: {tracking_url}

Estimasi tiba: {estimated_arrival}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/pesanan-dikirim', 'shipping');
INSERT INTO public.reply_templates VALUES ('933272ef-be7b-44dd-9fb8-8eaa07c48592', 'Konfirmasi Penerimaan Paket', 'Halo, berdasarkan data kurir, paket Anda sudah dinyatakan terkirim pada {delivered_at}.

Apakah paket sudah Anda terima dengan kondisi baik? Mohon konfirmasi agar pesanan dapat diselesaikan. 📦✅', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/konfirmasi-penerimaan', 'shipping');
INSERT INTO public.reply_templates VALUES ('57dc278b-9cab-43bf-8836-901e1b4fe9cf', 'Paket Tertahan / Terlambat', 'Kami melihat terdapat kendala pada pengiriman paket Anda (no. resi: {tracking_number}). ⚠️

Tim kami sedang berkoordinasi dengan pihak kurir untuk menindaklanjuti hal ini. Kami akan segera menginformasikan perkembangannya. Mohon maaf atas ketidaknyamanannya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/paket-tertahan', 'shipping');
INSERT INTO public.reply_templates VALUES ('0c4d32dc-44fc-402b-9fa5-44a181768d95', 'Informasi Ketersediaan Stok', 'Terima kasih atas minat Anda pada produk {product_name}.

Saat ini stok produk tersebut {stock_status}.

Anda dapat mengaktifkan notifikasi "Ingatkan Saya" pada halaman produk agar kami bisa memberitahu Anda saat stok tersedia kembali.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/ketersediaan-stok', 'product');
INSERT INTO public.reply_templates VALUES ('f2db8965-96ee-4957-a034-9bafb1981291', 'Detail Spesifikasi Produk', 'Berikut informasi lengkap mengenai {product_name}:

📌 Bahan    : {material}
📐 Ukuran   : {size}
🎨 Warna    : {color}
🏷️ Berat    : {weight}
✅ Halal    : {halal_status}

Apakah ada pertanyaan lain mengenai produk ini?', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/spesifikasi-produk', 'product');
INSERT INTO public.reply_templates VALUES ('d1a8034b-7a77-47f1-b549-440ff1a13d66', 'Verifikasi Berhasil', 'Selamat! Verifikasi akun Anda telah berhasil. ✅

Akun Anda kini sudah terverifikasi dan dapat menikmati semua fitur UmrahMart tanpa batasan. Terima kasih atas kerja samanya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/verification/berhasil', 'verification');
INSERT INTO public.reply_templates VALUES ('31b16e1b-998a-4779-87a2-663ef09cffb4', 'Produk Tidak Lagi Tersedia', 'Mohon maaf, produk {product_name} saat ini sudah tidak tersedia di katalog kami. 😔

Kami merekomendasikan produk serupa yang mungkin sesuai dengan kebutuhan Anda:
👉 {alternative_product}

Silakan cek koleksi terbaru kami di aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/produk-tidak-tersedia', 'product');
INSERT INTO public.reply_templates VALUES ('9d2a3755-e301-4552-a55f-d1ca99ee9aab', 'Bantuan Reset Password', 'Untuk mereset password akun Anda, ikuti langkah berikut:

1️⃣ Buka halaman Login
2️⃣ Klik "Lupa Password?"
3️⃣ Masukkan email terdaftar Anda
4️⃣ Cek inbox email untuk link reset (cek juga folder Spam)
5️⃣ Klik link dan buat password baru

Link reset berlaku selama 60 menit. Jika belum menerima email, hubungi kami kembali.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/reset-password', 'account');
INSERT INTO public.reply_templates VALUES ('52d1c3ad-e619-47c3-930c-2497057c746c', 'Verifikasi Email', 'Untuk memverifikasi email Anda:

1️⃣ Cek inbox email {email} (termasuk folder Spam/Junk)
2️⃣ Buka email dari UmrahMart dengan subjek "Verifikasi Akun"
3️⃣ Klik tombol "Verifikasi Sekarang"

Jika email tidak ditemukan, klik "Kirim Ulang Verifikasi" di halaman akun Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/verifikasi-email', 'account');
INSERT INTO public.reply_templates VALUES ('849c3f8f-05e6-4ad7-b8a7-9b39914b99ed', 'Akun Diblokir / Dinonaktifkan', 'Kami melihat akun Anda saat ini tidak aktif. 🔒

Hal ini dapat terjadi karena:
• Aktivitas yang mencurigakan terdeteksi
• Pelanggaran syarat & ketentuan penggunaan
• Permintaan penonaktifan sebelumnya

Untuk mengajukan reaktivasi, mohon kirimkan data diri Anda ke email support@umrahmart.id.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/akun-diblokir', 'account');
INSERT INTO public.reply_templates VALUES ('c729c57e-34d3-475d-8d65-7d0ff314a51d', 'Keluhan Diterima dan Dicatat', 'Kami sangat menyesal mendengar pengalaman yang tidak menyenangkan ini. 🙏

Keluhan Anda telah kami catat dengan nomor tiket {ticket_number}. Tim kami akan menginvestigasi dan menghubungi Anda dalam 1×24 jam kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/keluhan-diterima', 'complaint');
INSERT INTO public.reply_templates VALUES ('08411013-5cb3-497d-87a5-14b18b27361d', 'Keluhan Sedang Diinvestigasi', 'Kami tengah menginvestigasi keluhan yang Anda sampaikan terkait {complaint_subject}.

Proses investigasi membutuhkan waktu maksimal 3 hari kerja. Kami mohon kesabaran Anda dan akan segera memberikan informasi lebih lanjut.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/sedang-investigasi', 'complaint');
INSERT INTO public.reply_templates VALUES ('2be83748-504e-4be7-8895-47ef175d1689', 'Keluhan Telah Diselesaikan', 'Keluhan Anda dengan nomor tiket {ticket_number} telah kami selesaikan. ✅

Solusi yang diberikan: {resolution}

Kami memohon maaf atas pengalaman yang kurang menyenangkan ini dan berterima kasih atas masukan Anda untuk perbaikan layanan kami.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/telah-diselesaikan', 'complaint');
INSERT INTO public.reply_templates VALUES ('29ef4e0e-fff3-4d1d-998f-b4b8c16956f6', 'Prosedur Pengembalian Barang', 'Berikut prosedur pengembalian barang (return) di UmrahMart:

1️⃣ Ajukan return melalui menu "Pesanan Saya" → "Ajukan Return"
2️⃣ Pilih produk dan alasan pengembalian
3️⃣ Unggah foto kondisi barang
4️⃣ Tunggu konfirmasi persetujuan (1–2 hari kerja)
5️⃣ Kirim barang ke alamat yang tertera

Syarat: barang dalam kondisi asli, belum dipakai, dan beserta kemasan lengkap. Maksimal 7 hari setelah diterima.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/prosedur-pengembalian', 'return');
INSERT INTO public.reply_templates VALUES ('8c893df1-0794-4a42-99d1-dbc94ca8c470', 'Pengajuan Return Disetujui', 'Kabar baik! Pengajuan return Anda untuk pesanan {order_number} telah disetujui. ✅

Silakan kirimkan barang ke:
📍 {return_address}
Atas nama: UmrahMart Returns

Setelah barang kami terima dan verifikasi, dana akan dikembalikan dalam 3–5 hari kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/disetujui', 'return');
INSERT INTO public.reply_templates VALUES ('a4203d8e-4615-4ef0-b3eb-1a90bc79924e', 'Pengajuan Return Ditolak', 'Mohon maaf, pengajuan return untuk pesanan {order_number} tidak dapat kami proses. ❌

Alasan: {rejection_reason}

Jika Anda merasa keberatan, Anda dapat mengajukan banding dengan menghubungi kami dan melampirkan bukti pendukung.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/ditolak', 'return');
INSERT INTO public.reply_templates VALUES ('81873554-d801-46b4-8c20-a68f0aa0b604', 'Informasi Promo Aktif', 'Halo! Berikut promo yang sedang berlangsung di UmrahMart 🎉:

🏷️ {promo_name}
💰 Diskon hingga {discount_value}
📅 Berlaku: {promo_start} – {promo_end}
📝 Syarat: {promo_terms}

Jangan sampai terlewat! Belanja sekarang sebelum promo berakhir.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/promo/info-promo-aktif', 'promo');
INSERT INTO public.reply_templates VALUES ('6dcdda2a-39dc-4762-b8f0-754619e8da71', 'Promo Tidak Berlaku untuk Pesanan Ini', 'Mohon maaf, promo yang Anda gunakan tidak berlaku untuk pesanan ini. ⚠️

Kemungkinan penyebab:
• Promo sudah berakhir
• Produk tidak termasuk dalam kategori promo
• Minimum pembelian tidak terpenuhi

Cek syarat & ketentuan promo di halaman "Promo" pada aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/promo/tidak-berlaku', 'promo');
INSERT INTO public.reply_templates VALUES ('716a7e0f-9c10-404c-8fe1-a019a38adc7d', 'Cara Menggunakan Voucher', 'Berikut cara menggunakan voucher di UmrahMart:

1️⃣ Tambahkan produk ke keranjang
2️⃣ Masuk ke halaman Checkout
3️⃣ Klik "Gunakan Voucher"
4️⃣ Masukkan kode voucher: {voucher_code}
5️⃣ Klik "Terapkan" — diskon akan langsung terpotong

Pastikan pesanan Anda memenuhi minimum pembelian voucher.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/voucher/cara-menggunakan', 'voucher');
INSERT INTO public.reply_templates VALUES ('56959391-4b07-4cb8-89ab-a61abbef4748', 'Voucher Tidak Valid atau Kedaluwarsa', 'Kode voucher yang Anda masukkan tidak dapat digunakan. ❌

Kemungkinan penyebab:
• Kode voucher salah atau sudah digunakan
• Voucher sudah kedaluwarsa ({expired_date})
• Akun Anda tidak memenuhi syarat voucher

Untuk bantuan lebih lanjut, silakan kirimkan screenshot kode voucher Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/voucher/tidak-valid', 'voucher');
INSERT INTO public.reply_templates VALUES ('7ca45115-0adc-4c23-a1a1-5924e6177e20', 'Masalah Teknis pada Aplikasi', 'Kami mohon maaf atas kendala teknis yang Anda alami. 🛠️

Beberapa langkah yang dapat dicoba:
1. Tutup dan buka kembali aplikasi
2. Pastikan koneksi internet stabil
3. Update aplikasi ke versi terbaru
4. Hapus cache aplikasi
5. Restart perangkat

Jika masalah berlanjut, mohon kirimkan screenshot error beserta tipe perangkat dan versi OS Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/masalah-teknis-aplikasi', 'technical');
INSERT INTO public.reply_templates VALUES ('c139c11f-0159-4e62-b765-58dcc55b58eb', 'Tidak Bisa Login ke Akun', 'Kami memahami betapa frustrasinya tidak bisa masuk ke akun. Mari kami bantu! 🔑

Silakan coba langkah berikut:
1. Pastikan email dan password sudah benar
2. Aktifkan Caps Lock/periksa huruf kapital
3. Gunakan fitur "Lupa Password" untuk reset
4. Coba login dari perangkat atau browser berbeda

Apakah Anda menerima pesan error tertentu? Mohon informasikan agar kami dapat membantu lebih lanjut.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/tidak-bisa-login', 'technical');
INSERT INTO public.reply_templates VALUES ('82f514a7-153d-4f59-9327-229778c82585', 'Fitur Sedang Dalam Perbaikan', 'Kami menginformasikan bahwa fitur {feature_name} saat ini sedang dalam perbaikan/maintenance. 🔧

Estimasi selesai: {maintenance_end}

Kami mohon maaf atas ketidaknyamanannya. Tim teknis kami sedang bekerja keras untuk memulihkan layanan secepatnya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/fitur-dalam-perbaikan', 'technical');
INSERT INTO public.reply_templates VALUES ('774c04c8-9d21-4d1a-8c86-3ab0ef2987a0', 'Verifikasi Identitas Diperlukan', 'Untuk keamanan akun Anda, kami perlu melakukan verifikasi identitas. 🔐

Dokumen yang diperlukan:
📄 KTP (foto depan, jelas dan tidak blur)
🤳 Selfie sambil memegang KTP

Kirimkan dokumen melalui email ke: verify@umrahmart.id
Subject: Verifikasi Identitas - {user_id}

Proses verifikasi membutuhkan 1–2 hari kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/verification/identitas-diperlukan', 'verification');
INSERT INTO public.reply_templates VALUES ('967434ba-1ba8-459e-9514-f8c5dfdc3ec1', 'Informasi Pendaftaran Vendor', 'Terima kasih atas minat Anda untuk bergabung sebagai vendor di UmrahMart! 🏪

Persyaratan pendaftaran:
✅ KTP pemilik usaha
✅ Foto toko/tempat usaha
✅ Logo bisnis
✅ Banner bisnis
✅ Buku rekening bank
✅ NPWP (opsional)

Proses review membutuhkan 2–3 hari kerja setelah semua dokumen lengkap.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/info-pendaftaran', 'vendor');
INSERT INTO public.reply_templates VALUES ('96cb369d-1535-4e9f-b794-2db6a703901f', 'Status Vendor Sedang Direview', 'Pendaftaran vendor Anda sedang dalam proses review oleh tim kami. ⏳

Kami akan menginformasikan hasilnya melalui email dan notifikasi aplikasi dalam 2–3 hari kerja.

Pastikan data dan dokumen yang Anda unggah sudah lengkap dan terbaca dengan jelas.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/sedang-direview', 'vendor');
INSERT INTO public.reply_templates VALUES ('bf62a18f-4d7d-48f0-99a3-fa1fae0cf4fc', 'Vendor Berhasil Diaktifkan', 'Selamat! Akun vendor Anda telah berhasil diaktifkan! 🎉

Anda sekarang dapat mulai:
🛍️ Menambahkan produk ke katalog
📊 Mengelola stok dan harga
📦 Menerima dan memproses pesanan

Untuk panduan lengkap, kunjungi halaman "Panduan Vendor" di dashboard Anda. Semoga sukses!', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/berhasil-diaktifkan', 'vendor');
INSERT INTO public.reply_templates VALUES ('df27202f-3ef0-4ac9-84ea-e14506b78e01', 'Notifikasi Stok Tersedia', 'Kabar gembira! 🎉 Produk {product_name} yang Anda tunggu-tunggu kini sudah tersedia kembali.

Stok tersisa: {remaining_stock} pcs

Segera pesan sebelum kehabisan! Klik di sini: {product_url}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/stock/stok-tersedia', 'stock');
INSERT INTO public.reply_templates VALUES ('30cf40bb-4b5a-42d1-a288-7b1b3ff5ceb5', 'Informasi Pre-Order', 'Produk {product_name} saat ini tersedia dalam mode Pre-Order. 📋

Detail Pre-Order:
📅 Batas pemesanan : {po_end_date}
🚚 Estimasi kirim  : {estimated_delivery}
💰 Harga PO        : {po_price}

Pesan sekarang untuk mendapatkan kepastian stok!', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/stock/info-pre-order', 'stock');
INSERT INTO public.reply_templates VALUES ('24e59108-be8b-4641-85f0-8f11128cf74c', 'Konfirmasi Pembatalan Pesanan', 'Kami telah menerima permintaan pembatalan untuk pesanan {order_number}. ℹ️

Untuk melanjutkan proses pembatalan, mohon konfirmasi alasan Anda:
a) Salah alamat pengiriman
b) Ingin ganti produk
c) Menemukan harga lebih murah
d) Berubah pikiran
e) Lainnya

Balas dengan huruf pilihan Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/konfirmasi-pembatalan', 'cancellation');
INSERT INTO public.reply_templates VALUES ('7ffccfea-bef8-45d5-b8cb-c46798a0da55', 'Pesanan Tidak Dapat Dibatalkan', 'Mohon maaf, pesanan {order_number} tidak dapat dibatalkan karena status pesanan sudah dalam tahap {order_status}. ⚠️

Jika barang sudah diterima dan terdapat masalah, Anda masih bisa mengajukan return dalam 7 hari setelah penerimaan.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/tidak-dapat-dibatalkan', 'cancellation');
INSERT INTO public.reply_templates VALUES ('f04503ae-fb8b-4e54-b17d-b2473014537f', 'Pembatalan Berhasil', 'Pesanan {order_number} berhasil dibatalkan. ✅

Jika ada pembayaran yang perlu dikembalikan:
💰 Jumlah refund : {refund_amount}
⏱️ Estimasi dana masuk : 3–7 hari kerja

Terima kasih atas pengertian Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/berhasil', 'cancellation');
INSERT INTO public.reply_templates VALUES ('45b72f3e-52d9-499e-8d31-faa17af981be', 'Jadwal Estimasi Pengiriman', 'Berikut estimasi pengiriman berdasarkan lokasi Anda:

🏙️ Jabodetabek    : 1–2 hari kerja
🗺️ Pulau Jawa     : 2–3 hari kerja
🏝️ Luar Jawa      : 3–5 hari kerja
🗾 Papua/Terpencil : 5–10 hari kerja

Estimasi dihitung sejak paket diserahkan ke kurir, tidak termasuk hari libur nasional.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/delivery/estimasi-pengiriman', 'delivery');
INSERT INTO public.reply_templates VALUES ('ca6a9164-b296-4daf-9131-6902d1a3aa6e', 'Pengiriman Tertunda', 'Kami mohon maaf atas keterlambatan pengiriman paket Anda (no. resi: {tracking_number}). 🙏

Kendala yang terjadi: {delay_reason}

Tim kami sedang berkoordinasi dengan pihak kurir {courier_name} untuk mempercepat proses pengiriman. Kami akan menginformasikan update terbaru sesegera mungkin.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/delivery/pengiriman-tertunda', 'delivery');
INSERT INTO public.reply_templates VALUES ('a217d188-4f32-4184-80a2-cadda1107c5d', 'Informasi Pengambilan di Toko (Pickup)', 'Pesanan Anda dapat diambil langsung di:

📍 {store_address}
🕐 Jam operasional: {store_hours}

Tunjukkan nomor pesanan {order_number} atau kode QR di aplikasi kepada petugas. Pesanan akan disimpan selama 3 hari sejak notifikasi siap diambil.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/pickup/info-pickup', 'pickup');
INSERT INTO public.reply_templates VALUES ('cb92a7d1-22cc-4fb6-afa0-3b7c0140eb93', 'Pesanan Siap Diambil', 'Pesanan Anda sudah siap untuk diambil! 🛍️

Nomor Pesanan : {order_number}
Lokasi        : {store_address}
Batas Ambil   : {pickup_deadline}

Jangan lupa bawa identitas diri saat pengambilan. Sampai jumpa!', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/pickup/siap-diambil', 'pickup');
INSERT INTO public.reply_templates VALUES ('a2940fd9-f18b-4024-8694-29be71b0cf23', 'Informasi Paket Langganan', 'Berikut informasi paket langganan UmrahMart Premium:

⭐ Paket Basic   : Rp 29.000/bulan — Gratis ongkir 2x/bulan
⭐ Paket Premium : Rp 59.000/bulan — Gratis ongkir unlimited + cashback 5%
⭐ Paket VIP     : Rp 99.000/bulan — Semua benefit + akses produk eksklusif

Daftar sekarang di menu "Langganan" pada aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/subscription/info-langganan', 'subscription');
INSERT INTO public.reply_templates VALUES ('fec12945-f26b-4e5b-892d-e53b916e5ad3', 'Masa Aktif Langganan Hampir Habis', 'Halo! Masa aktif langganan UmrahMart Premium Anda akan berakhir pada {expiry_date}. ⏰

Perpanjang sekarang dan nikmati terus benefit:
✅ Gratis ongkos kirim
✅ Cashback eksklusif member
✅ Early access produk baru

Klik di sini untuk perpanjang: {renewal_url}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/subscription/hampir-habis', 'subscription');
INSERT INTO public.reply_templates VALUES ('8ae2e394-0d0c-4fb7-8744-f13bbb7d448e', 'Permintaan Ulasan Produk', 'Assalamu''alaikum! Semoga pesanan Anda sudah diterima dengan baik. 😊

Kami akan sangat berterima kasih jika Anda meluangkan waktu untuk memberikan ulasan pada produk {product_name}.

Ulasan Anda sangat membantu calon pembeli lain dalam memilih produk yang tepat. Berikan ulasan di menu "Pesanan Saya" → "Beri Ulasan".', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/review/minta-ulasan', 'review');
INSERT INTO public.reply_templates VALUES ('011a94d2-8ce3-4594-87ba-63e9b46fe1eb', 'Terima Kasih Atas Ulasan Anda', 'Jazakallahu khairan atas ulasan yang telah Anda berikan! 🌟

Masukan Anda sangat berharga bagi kami untuk terus meningkatkan kualitas produk dan layanan. Semoga UmrahMart selalu bisa menjadi teman setia perjalanan ibadah Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/review/terima-kasih-ulasan', 'review');
INSERT INTO public.reply_templates VALUES ('23b40e79-e1f3-4c17-a3ee-64ec201cd2ae', 'Laporan Penipuan Diterima', 'Laporan penipuan Anda telah kami terima dan kami tangani dengan sangat serius. 🚨

Nomor laporan: {report_number}

Tim keamanan kami akan menginvestigasi dalam 1×24 jam. Jangan lakukan transfer atau berikan data sensitif kepada pihak yang mencurigakan.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/fraud/laporan-penipuan-diterima', 'fraud');
INSERT INTO public.reply_templates VALUES ('53916ccc-b29b-40cb-a795-dc5f7de740e5', 'Tindakan Keamanan Akun', 'Kami mendeteksi aktivitas tidak biasa pada akun Anda. 🔐

Sebagai tindakan keamanan, akun Anda telah kami sementara kunci.

Langkah yang perlu Anda lakukan:
1. Segera ubah password Anda
2. Aktifkan verifikasi 2 langkah
3. Hubungi kami melalui {support_contact} untuk membuka kunci akun

Jangan pernah bagikan OTP atau password kepada siapapun, termasuk tim kami.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/fraud/tindakan-keamanan-akun', 'fraud');
INSERT INTO public.reply_templates VALUES ('1a57d5dc-1751-43c9-bdf7-9ccfc1d9898c', 'Pertanyaan Tidak Dapat Kami Jawab Saat Ini', 'Terima kasih atas pertanyaan Anda. Pertanyaan ini membutuhkan penanganan dari tim khusus kami. 🔄

Kami akan meneruskan pertanyaan ini ke tim terkait dan akan menghubungi Anda kembali dalam 1×24 jam kerja.

Nomor tiket Anda: {ticket_number}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/other/eskalasi-ke-tim-terkait', 'other');
INSERT INTO public.reply_templates VALUES ('9a10b61b-1dc3-407a-b5a5-b80062a404b4', 'Cara Menghubungi Kami', 'Ada beberapa cara untuk menghubungi tim UmrahMart:

💬 Live Chat  : Aplikasi UmrahMart (menu "Bantuan")
📧 Email      : support@umrahmart.id
📞 Telepon    : 021-XXXX-XXXX (08.00–21.00 WIB)
📱 WhatsApp   : wa.me/628XXXXXXXXX
📘 Instagram  : @umrahmart.id

Kami siap membantu Anda! 🙏', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/other/cara-menghubungi-kami', 'other');


--
-- Data for Name: return_reasons; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.schema_migrations VALUES (52, false) ON CONFLICT DO NOTHING;


--
-- Data for Name: shipments; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.shipments VALUES ('0ed06a39-93f1-437a-91b6-e146e98cd932', 'd31dc0b1-0079-4f43-971a-d0f86fe64cfb', 'sicepat', 'REG', '', 'waiting_pickup', NULL, NULL, '1-2 day');
INSERT INTO public.shipments VALUES ('3627c0f5-11cf-466e-bc43-2d527edc629b', '4605d7ba-5a51-49d1-b793-29d82de1099d', 'jne', 'JTR', '', 'waiting_pickup', NULL, NULL, '3 day');
INSERT INTO public.shipments VALUES ('e1dedb0f-6ab3-4866-9e2e-0aa9adda7ef4', 'e0ee921c-9ed6-41e7-b1c4-ae3c55fc6280', 'sicepat', 'REG', '', 'waiting_pickup', NULL, NULL, '1-2 day');
INSERT INTO public.shipments VALUES ('18b0f3fd-56f4-4481-9720-05ff66b0b42b', '3c6b4adf-7643-4336-9088-a6ba755e30d4', 'jnt', 'EZ', 'JX7629871362', 'shipped', '2026-03-31 01:28:47.500271+00', NULL, '');
INSERT INTO public.shipments VALUES ('0b8604ea-3eb8-480a-bc15-d71583e0da48', '959a0bc5-6d71-4f34-a1da-8eb6ba454073', 'jnt', 'EZ', 'JX7629871362', 'delivered', '2026-03-31 01:30:04.952759+00', '2026-03-24 08:41:52+00', '');


--
-- Data for Name: ticket_subjects; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ticket_subjects VALUES (1, 'Barang Tidak Sampai', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (2, 'Barang Rusak', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (3, 'Pengembalian Dana', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (4, 'Barang Tidak Sesuai', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (5, 'Pengiriman Terlambat', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (6, 'Pembatalan Pesanan', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (7, 'Kesalahan Produk', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (8, 'Akun Bermasalah', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (9, 'Pembayaran Gagal', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (10, 'Voucher/Promo Tidak Berlaku', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (11, 'Pertanyaan Umum', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (12, 'Lainnya', true, '2026-03-31 00:34:09.934189+00') ON CONFLICT DO NOTHING;


--
-- Data for Name: tickets; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.tickets VALUES ('0b6082a1-173e-4719-986a-62b6541c9bce', 'TKT-20260331-0001', 'cf0add84-b149-40c7-ae68-fb9d02594cae', NULL, 'ORD-20260331-0242A9B9', '081200000006', 'Ahmad Subarkah', 'Barang Tidak Sampai', 'Saya sudah checkout di hari 31 Maret 2026, kok belum sampai ya?. Mohon dicek.', 'open', 'web', NULL, '2026-03-31 01:23:51.315238+00', '2026-03-31 01:23:51.315238+00', 'https://kemenhaj.s3.ap-southeast-1.amazonaws.com/ticket-attachments/2026/03/31/ba810cb9-1caf-42cb-b259-b013969b5a79', 'image/png', 1);
INSERT INTO public.tickets VALUES ('b3ed4e66-3e67-47b8-a93e-1146b416e064', 'TKT-20260331-0002', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', 'ORD-20260331-84E3D283', '081200000006', 'Ahmad Subarkah', 'Barang Tidak Sesuai', 'Barang sampai namun salah kirim merk', 'on_progress', 'web', NULL, '2026-03-31 01:31:44.254841+00', '2026-03-31 01:32:19.531951+00', 'https://kemenhaj.s3.ap-southeast-1.amazonaws.com/ticket-attachments/2026/03/31/8c321f23-4ba1-40a3-9250-08cc85db07de', 'image/png', 4);


--
-- Data for Name: ticket_messages; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ticket_messages VALUES ('e5b56e84-d944-4319-9ea1-876cb11184ad', 'b3ed4e66-3e67-47b8-a93e-1146b416e064', '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', 'Assalamu''alaikum, selamat datang di UmrahMart! 🕌

Kami hadir untuk memudahkan Anda mendapatkan perlengkapan haji dan umroh terbaik. Ada yang bisa kami bantu hari ini?', true, false, '2026-03-31 01:33:25.279677+00');
INSERT INTO public.ticket_messages VALUES ('bb0e54f2-fccb-471a-9dee-67a206f111fd', 'b3ed4e66-3e67-47b8-a93e-1146b416e064', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'Barangnya kok tidak sesuai merk ya?', false, false, '2026-03-31 01:39:16.647949+00');


--
-- Data for Name: ticket_attachments; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: ticket_status_logs; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.ticket_status_logs VALUES ('96914520-8ce0-4def-9095-6356e2da621a', 'b3ed4e66-3e67-47b8-a93e-1146b416e064', '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', 'open', 'on_progress', NULL, '2026-03-31 01:32:17.495983+00');
INSERT INTO public.ticket_status_logs VALUES ('5794cbd4-3ed1-4779-8419-059f979f63e1', 'b3ed4e66-3e67-47b8-a93e-1146b416e064', '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', 'open', 'on_progress', NULL, '2026-03-31 01:32:19.531841+00');


--
-- Data for Name: vendor_balances; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_balances VALUES ('297a0ab4-2168-4d54-9bf1-f526a476ac19', 100000.00, 0.00, 100000.00, 0.00, '2026-03-31 01:38:34.011887+00', 0.00);


--
-- Data for Name: vendor_bank_accounts; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vendor_banners; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_banners VALUES ('43c93ca5-e5fc-490a-ac5a-372779ba353c', '297a0ab4-2168-4d54-9bf1-f526a476ac19', 'Contoh banner', 'vendor-banners/297a0ab4-2168-4d54-9bf1-f526a476ac19/43c93ca5-e5fc-490a-ac5a-372779ba353c/1774921866425', '2026-03-31 01:51:06.425616+00', '2026-03-31 01:51:43.805196+00');


--
-- Data for Name: vendor_couriers; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_couriers VALUES ('2b58adaa-54c4-48ff-889f-d15f8683be8e', '297a0ab4-2168-4d54-9bf1-f526a476ac19', 6, true, '2026-03-31 01:16:53.888353+00');
INSERT INTO public.vendor_couriers VALUES ('2d3dbe85-4617-4d2c-9195-1c22fe94a606', '297a0ab4-2168-4d54-9bf1-f526a476ac19', 1, true, '2026-03-31 01:16:53.889477+00');
INSERT INTO public.vendor_couriers VALUES ('cd1134f1-f5b3-4cc7-9009-2a8fb175cafd', '297a0ab4-2168-4d54-9bf1-f526a476ac19', 2, true, '2026-03-31 01:16:53.890574+00');


--
-- Data for Name: vendor_documents; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_documents VALUES ('acae82b0-6bb1-41ba-97bb-3771a5141d4b', '297a0ab4-2168-4d54-9bf1-f526a476ac19', 'owner_document_id', 'vendor-onboardings/fed503f5-53f3-408e-82ad-03703894d97d/documents/owner_document_id/179e1ad3-cbd8-4b0b-9732-6773ad5b6257', 'application/pdf', 68609, NULL, '2aea50ec-ff64-4f82-b729-c2f4085a5622', NULL, NULL, '2026-03-31 00:38:29.394538+00', '2026-03-31 00:38:29.394538+00');


--
-- Data for Name: vendor_onboardings; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.vendor_onboardings VALUES ('fed503f5-53f3-408e-82ad-03703894d97d', '404lamfound@gmail.com', 'completed', NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, NULL, '2026-03-31 00:37:46.560083+00', '2026-03-31 00:38:29.394538+00', '2026-03-31 00:37:46.560083+00', '2026-03-31 00:38:29.394538+00', NULL, NULL, NULL, NULL, NULL);


--
-- Data for Name: vendor_vouchers; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vendor_voucher_products; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: vendor_withdrawals; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: wishlist_items; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Name: admin_contacts_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.admin_contacts_id_seq', 1, false);


--
-- Name: couriers_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.couriers_id_seq', 12, true);


--
-- Name: faqs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.faqs_id_seq', 1, false);


--
-- Name: return_reasons_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.return_reasons_id_seq', 1, false);


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


