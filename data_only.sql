--
-- PostgreSQL database dump
--



-- Dumped from database version 17.8
-- Dumped by pg_dump version 17.9 (Ubuntu 17.9-1.pgdg24.04+1)
SET session_replication_role = replica;

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

SET SESSION AUTHORIZATION DEFAULT;

ALTER TABLE public.roles DISABLE TRIGGER ALL;

INSERT INTO public.roles VALUES (1, 'admin', 'Administrator') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (2, 'umkm', 'UMKM') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (3, 'customer', 'Customer') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (4, 'cs', 'Customer Service') ON CONFLICT DO NOTHING;
INSERT INTO public.roles VALUES (5, 'finance', 'Finance') ON CONFLICT DO NOTHING;


ALTER TABLE public.roles ENABLE TRIGGER ALL;

--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.users DISABLE TRIGGER ALL;

INSERT INTO public.users VALUES ('577c7f1a-8d9d-4c94-9758-920cc72e81d7', 'harundarat@gmail.com', 'System Administrator', NULL, '0811111', '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 1, 'active', '2026-03-31 00:39:54.52106+00', '2026-03-31 00:39:54.52106+00', '2026-03-31 00:39:54.52106+00', NULL);
INSERT INTO public.users VALUES ('8e3a1e10-8604-4dc4-87db-838cb5bea0bd', 'nobev17858@nexafilm.com', 'Ahmad Subarkah', '1990-01-02', NULL, '$2a$10$z6kONryRQp5jSS6LOXepc.O.A9SUvRR0XS.VaXSILyP6M777w7Mee', 2, 'active', '2026-04-04 12:05:33.823588+00', '2026-04-04 12:05:58.662114+00', '2026-04-04 12:05:58.662114+00', NULL);
INSERT INTO public.users VALUES ('3af6ae21-eb2a-45ff-b2cb-5010ef5e7e5c', 'rizkyalamsyah.dev@gmail.com', 'Rizky Alamsyah', '1990-01-02', NULL, '$2a$10$CKrCheEGOgWpATRZQB8T1.W97OtNxUutk252W5ukdyJSPmiMep5.m', 2, 'active', '2026-04-04 11:58:09.094933+00', '2026-04-04 12:01:26.504538+00', '2026-04-04 12:01:26.504538+00', NULL);
INSERT INTO public.users VALUES ('cf0add84-b149-40c7-ae68-fb9d02594cae', '404lamfound@gmail.com', 'Dev Customer', NULL, '081200000006', '$2a$10$mH97QEwzF5qUzTsS6g.5NenBhi0lKA/h1rlsc2nvcQt5t9DDEgbN6', 3, 'active', '2026-03-31 00:39:54.850767+00', '2026-03-31 00:39:54.850767+00', '2026-03-31 00:39:54.850767+00', NULL);


ALTER TABLE public.users ENABLE TRIGGER ALL;

--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.addresses DISABLE TRIGGER ALL;

INSERT INTO public.addresses VALUES ('0833971b-ab80-4049-ba1e-9ddb9d321bdb', '3af6ae21-eb2a-45ff-b2cb-5010ef5e7e5c', 'Warehouse Sedati', 'Admin warehouse A', '071530174921', '18', '583', '6001', '61253', 'Jl. Raya Sedati No. 99', true, 'JAWA TIMUR', ' SIDOARJO', 'SEDATI', '70995', 'SEDATI GEDE', NULL, NULL, NULL, '2026-04-04 12:17:19.977604+00', '2026-04-04 12:17:19.977604+00');
INSERT INTO public.addresses VALUES ('299885ab-890a-45ce-8cd4-2cf57757a826', '8e3a1e10-8604-4dc4-87db-838cb5bea0bd', 'Warehouse Palmerah Jakbar', 'Admin gudang palmerah', '098765432123', '10', '135', '1326', '11480', 'Jl. Palmera No. 46', true, 'DKI JAKARTA', ' JAKARTA BARAT', 'PALMERAH', '17502', 'PALMERAH', NULL, NULL, NULL, '2026-04-04 12:25:58.983544+00', '2026-04-04 12:25:58.983544+00');
INSERT INTO public.addresses VALUES ('df86c486-d12c-4522-87f6-895fc25d3088', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'Alamat Keputran', 'Zaki', '063819235476', '18', '577', '5899', '60265', 'Jl. Tegalsari kecamatan keputran no 123', true, 'JAWA TIMUR', ' SURABAYA', 'TEGALSARI', '69344', 'KEPUTRAN', NULL, NULL, NULL, '2026-04-04 12:37:44.356033+00', '2026-04-04 12:37:44.356033+00');


ALTER TABLE public.addresses ENABLE TRIGGER ALL;

--
-- Data for Name: admin_contacts; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.admin_contacts DISABLE TRIGGER ALL;



ALTER TABLE public.admin_contacts ENABLE TRIGGER ALL;

--
-- Data for Name: banners; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.banners DISABLE TRIGGER ALL;



ALTER TABLE public.banners ENABLE TRIGGER ALL;

--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.carts DISABLE TRIGGER ALL;

INSERT INTO public.carts VALUES ('997ab39f-3d75-4b32-ae5b-763fe2d908eb', 'cf0add84-b149-40c7-ae68-fb9d02594cae', 'active', '2026-04-04 12:35:53.707901+00', '2026-04-04 12:42:19.609555+00');


ALTER TABLE public.carts ENABLE TRIGGER ALL;

--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.categories DISABLE TRIGGER ALL;

INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000001', 'ca000001-0000-0000-0000-000000000001', 'Perlengkapan Haji & Umrah', 'perlengkapan-haji-umrah', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000002', 'ca000001-0000-0000-0000-000000000002', 'Oleh-oleh & Souvenir', 'oleh-oleh-souvenir', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000003', 'ca000001-0000-0000-0000-000000000003', 'Fashion Muslim', 'fashion-muslim', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000011', 'ca000001-0000-0000-0000-000000000001', 'Pakaian Ihram', 'pakaian-ihram', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000012', 'ca000001-0000-0000-0000-000000000001', 'Sajadah & Alat Shalat', 'sajadah-alat-shalat', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000013', 'ca000001-0000-0000-0000-000000000001', 'Tasbih & Hampers', 'tasbih-hampers', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000021', 'ca000001-0000-0000-0000-000000000002', 'Kurma', 'kurma', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000022', 'ca000001-0000-0000-0000-000000000002', 'Minyak & Herbal', 'minyak-herbal', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000023', 'ca000001-0000-0000-0000-000000000002', 'Air Zamzam', 'air-zamzam', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000031', 'ca000001-0000-0000-0000-000000000003', 'Mukena', 'mukena', true);
INSERT INTO public.categories VALUES ('ca000001-0000-0000-0000-000000000032', 'ca000001-0000-0000-0000-000000000003', 'Gamis & Jubah', 'gamis-jubah', true);


ALTER TABLE public.categories ENABLE TRIGGER ALL;

--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendors DISABLE TRIGGER ALL;

INSERT INTO public.vendors (id, owner_user_id, vendor_type, legal_name, display_name, description, status, approved_by, approved_at, status_reason, xendit_account_id, created_at, updated_at, total_sold) VALUES ('0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', '8e3a1e10-8604-4dc4-87db-838cb5bea0bd', 'souvenir_store', NULL, 'Toko Oleh Oleh Haji Test', NULL, 'active', '577c7f1a-8d9d-4c94-9758-920cc72e81d7', '2026-04-04 12:12:18.79223+00', NULL, '69d10022d4417937e489e76b', '2026-04-04 12:05:58.662114+00', '2026-04-04 12:12:18.79223+00', 0);
INSERT INTO public.vendors (id, owner_user_id, vendor_type, legal_name, display_name, description, status, approved_by, approved_at, status_reason, xendit_account_id, created_at, updated_at, total_sold) VALUES ('95e22a26-5c24-44cc-b207-1455cf8656ed', '3af6ae21-eb2a-45ff-b2cb-5010ef5e7e5c', 'souvenir_store', NULL, 'Toko Oleh Oleh Haji Test', NULL, 'active', '577c7f1a-8d9d-4c94-9758-920cc72e81d7', '2026-04-04 12:12:37.827469+00', NULL, '69d1003545d4336498d0da87', '2026-04-04 12:01:26.504538+00', '2026-04-04 12:12:37.827469+00', 0);


ALTER TABLE public.vendors ENABLE TRIGGER ALL;

--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.products DISABLE TRIGGER ALL;

INSERT INTO public.products VALUES ('9b2b9956-3d69-4c05-98bf-c87b03fec86d', '95e22a26-5c24-44cc-b207-1455cf8656ed', 'ca000001-0000-0000-0000-000000000021', 'Kurma Ajwa', 'kurma-ajwa', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 'published', 'pending', NULL, '2026-04-04 12:18:30.516593+00', '2026-04-04 12:18:30.516593+00');
INSERT INTO public.products VALUES ('f0d3da47-3a75-43b1-8a5d-8e18d00f4a7d', '95e22a26-5c24-44cc-b207-1455cf8656ed', 'ca000001-0000-0000-0000-000000000003', 'Baju Musllim Lengan Panjang', 'baju-musllim-lengan-panjang', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 'published', 'pending', NULL, '2026-04-04 12:22:42.65505+00', '2026-04-04 12:22:42.65505+00');
INSERT INTO public.products VALUES ('cd83d2ab-b6cc-4b37-989c-b732beb8a897', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 'ca000001-0000-0000-0000-000000000023', 'Air Zamzam 5 Liter Asli 100% 1 DUS (5 Liter) Zam Zam Kemasan Galon Kecil Oleh Oleh Haji dan Umroh', 'air-zamzam-5-liter-asli-100-1-dus-5-liter-zam-zam-kemasan-galon-kecil-oleh-oleh-haji-dan-umroh', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 'published', 'pending', NULL, '2026-04-04 12:31:46.354883+00', '2026-04-04 12:31:46.354883+00');
INSERT INTO public.products VALUES ('9f5bfde2-09b4-44c3-b080-a17099d9a797', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 'ca000001-0000-0000-0000-000000000023', 'Minyak Zaitun Mustika Ratu 175ml - Melembabkan, Melembutkan, Menyegarkan', 'minyak-zaitun-mustika-ratu-175ml-melembabkan-melembutkan-menyegarkan', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.', 'published', 'pending', NULL, '2026-04-04 12:33:39.236631+00', '2026-04-04 12:33:39.236631+00');


ALTER TABLE public.products ENABLE TRIGGER ALL;

--
-- Data for Name: product_variants; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.product_variants DISABLE TRIGGER ALL;

INSERT INTO public.product_variants VALUES ('6680cd95-04b4-4459-a798-68f99cb9f03f', '9b2b9956-3d69-4c05-98bf-c87b03fec86d', 'TOKOO-KURMA000-00', 'Default', 75000.00, 'IDR', 20, 500, true, true);
INSERT INTO public.product_variants VALUES ('c369fc1f-f812-4f65-8967-c65fab818812', 'f0d3da47-3a75-43b1-8a5d-8e18d00f4a7d', 'TOKOO-BAJUM001-01', 'Baju Muslim LPJ Sage', 170000.00, 'IDR', 25, 250, false, true);
INSERT INTO public.product_variants VALUES ('a402ab6b-f02d-4eaa-b10b-3d7e8893d333', 'cd83d2ab-b6cc-4b37-989c-b732beb8a897', 'TOKOO-AIRZA000-00', 'Default', 298072.00, 'IDR', 15, 6000, true, true);
INSERT INTO public.product_variants VALUES ('a8d80d68-2834-46ab-a90d-b8eff9eb1ec0', '9f5bfde2-09b4-44c3-b080-a17099d9a797', 'TOKOO-MINYA001-00', 'Default', 35000.00, 'IDR', 48, 175, true, true);
INSERT INTO public.product_variants VALUES ('94d6f287-58e4-4698-bc3e-f30e2e9439c8', 'f0d3da47-3a75-43b1-8a5d-8e18d00f4a7d', 'TOKOO-BAJUM001-00', 'Default', 170000.00, 'IDR', 17, 250, true, true);


ALTER TABLE public.product_variants ENABLE TRIGGER ALL;

--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.cart_items DISABLE TRIGGER ALL;



ALTER TABLE public.cart_items ENABLE TRIGGER ALL;

--
-- Data for Name: chat_conversations; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.chat_conversations DISABLE TRIGGER ALL;



ALTER TABLE public.chat_conversations ENABLE TRIGGER ALL;

--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.chat_messages DISABLE TRIGGER ALL;



ALTER TABLE public.chat_messages ENABLE TRIGGER ALL;

--
-- Data for Name: couriers; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.couriers DISABLE TRIGGER ALL;

INSERT INTO public.couriers VALUES (1, 'jne', 'JNE', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (2, 'sicepat', 'SiCepat', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (3, 'ide', 'IDExpress', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (4, 'sap', 'SAP Express', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (5, 'ninja', 'Ninja', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (6, 'jnt', 'J&T Express', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (7, 'tiki', 'TIKI', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (8, 'wahana', 'Wahana Express', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (9, 'pos', 'POS Indonesia', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (10, 'sentral', 'Sentral Cargo', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (11, 'lion', 'Lion Parcel', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;
INSERT INTO public.couriers VALUES (12, 'rex', 'Royal Express Asia', NULL, true, '2026-04-04 08:43:39.566922+00') ON CONFLICT DO NOTHING;


ALTER TABLE public.couriers ENABLE TRIGGER ALL;

--
-- Data for Name: email_verification_tokens; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.email_verification_tokens DISABLE TRIGGER ALL;



ALTER TABLE public.email_verification_tokens ENABLE TRIGGER ALL;

--
-- Data for Name: faqs; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.faqs DISABLE TRIGGER ALL;



ALTER TABLE public.faqs ENABLE TRIGGER ALL;

--
-- Data for Name: ledger_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.ledger_accounts DISABLE TRIGGER ALL;

INSERT INTO public.ledger_accounts VALUES ('2bbd864c-d440-46fa-8ddd-94116057b385', '1100', 'Payment Gateway Receivable', 'asset', 'D', true) ON CONFLICT DO NOTHING;
INSERT INTO public.ledger_accounts VALUES ('e47fd62c-529e-4786-b216-ac4e9f79fac9', '2100', 'Vendor Payable', 'liability', 'C', true) ON CONFLICT DO NOTHING;
INSERT INTO public.ledger_accounts VALUES ('6017dd31-2620-4665-bb23-f31c999db0bc', '4100', 'Platform Fee Revenue', 'revenue', 'C', true) ON CONFLICT DO NOTHING;
INSERT INTO public.ledger_accounts VALUES ('3ccfa4e4-3c56-4433-adc6-e1fefb5f0584', '4200', 'Admin Fee Revenue', 'revenue', 'C', true) ON CONFLICT DO NOTHING;


ALTER TABLE public.ledger_accounts ENABLE TRIGGER ALL;

--
-- Data for Name: ledger_journals; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.ledger_journals DISABLE TRIGGER ALL;

INSERT INTO public.ledger_journals VALUES ('a7d5c118-ba7e-456f-bf42-7b37bbaaa4af', 'JRN-PAY-ORD-20260404-6E8B747E-827845ac', 'payment_invoice', 'c746ecd5-cb8d-45b3-a6f2-9a8ecffbd1de', '2026-04-04 12:44:13.481736+00', 'Payment received for order ORD-20260404-6E8B747E', 'posted', NULL, '2026-04-04 12:44:13.481736+00');
INSERT INTO public.ledger_journals VALUES ('dacde11b-84c2-4206-aba1-970f25f35aaa', 'JRN-PAY-ORD-20260404-C542C29A-4d30578f', 'payment_invoice', '6444f73a-49aa-4dd1-b31c-cb352545dbf6', '2026-04-04 12:45:01.071215+00', 'Payment received for order ORD-20260404-C542C29A', 'posted', NULL, '2026-04-04 12:45:01.071215+00');
INSERT INTO public.ledger_journals VALUES ('120634fc-a664-45b0-955f-88c00c68d74d', 'JRN-PAY-ORD-20260404-50B5ABC2-3cec1f1b', 'payment_invoice', '0fa1f7a7-c377-4c35-b77a-1b1d38b9c28c', '2026-04-04 12:45:35.524032+00', 'Payment received for order ORD-20260404-50B5ABC2', 'posted', NULL, '2026-04-04 12:45:35.524032+00');


ALTER TABLE public.ledger_journals ENABLE TRIGGER ALL;

--
-- Data for Name: ledger_lines; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.ledger_lines DISABLE TRIGGER ALL;

INSERT INTO public.ledger_lines VALUES ('86e9ed7f-dcc4-4ece-87c5-d99c6da78437', 'a7d5c118-ba7e-456f-bf42-7b37bbaaa4af', '2bbd864c-d440-46fa-8ddd-94116057b385', 181000.00, 0.00, 'IDR', 'ORD-20260404-6E8B747E');
INSERT INTO public.ledger_lines VALUES ('1f751baf-a0d1-434f-bb30-12e33a9b37ec', 'a7d5c118-ba7e-456f-bf42-7b37bbaaa4af', 'e47fd62c-529e-4786-b216-ac4e9f79fac9', 0.00, 170000.00, 'IDR', 'ORD-20260404-6E8B747E');
INSERT INTO public.ledger_lines VALUES ('51ac56a1-15c2-49df-937c-6fddecb37b68', 'a7d5c118-ba7e-456f-bf42-7b37bbaaa4af', '6017dd31-2620-4665-bb23-f31c999db0bc', 0.00, 1000.00, 'IDR', 'ORD-20260404-6E8B747E');
INSERT INTO public.ledger_lines VALUES ('2fb1bfee-b91d-4054-8b1c-f8a801c0dac1', 'a7d5c118-ba7e-456f-bf42-7b37bbaaa4af', '3ccfa4e4-3c56-4433-adc6-e1fefb5f0584', 0.00, 5000.00, 'IDR', 'ORD-20260404-6E8B747E');
INSERT INTO public.ledger_lines VALUES ('4fcc1ca1-cff7-40c8-8026-9956b1a32958', 'dacde11b-84c2-4206-aba1-970f25f35aaa', '2bbd864c-d440-46fa-8ddd-94116057b385', 357000.00, 0.00, 'IDR', 'ORD-20260404-C542C29A');
INSERT INTO public.ledger_lines VALUES ('ee333a70-2393-4200-946b-149ef997c7db', 'dacde11b-84c2-4206-aba1-970f25f35aaa', 'e47fd62c-529e-4786-b216-ac4e9f79fac9', 0.00, 340000.00, 'IDR', 'ORD-20260404-C542C29A');
INSERT INTO public.ledger_lines VALUES ('a8bdef3e-5b06-44f6-b8f5-050af7ecb1b1', 'dacde11b-84c2-4206-aba1-970f25f35aaa', '6017dd31-2620-4665-bb23-f31c999db0bc', 0.00, 1000.00, 'IDR', 'ORD-20260404-C542C29A');
INSERT INTO public.ledger_lines VALUES ('28b35f98-1139-445b-87ec-25a27127105e', 'dacde11b-84c2-4206-aba1-970f25f35aaa', '3ccfa4e4-3c56-4433-adc6-e1fefb5f0584', 0.00, 5000.00, 'IDR', 'ORD-20260404-C542C29A');
INSERT INTO public.ledger_lines VALUES ('0c0b8767-8fa8-41ca-b776-6edc6ca65f02', '120634fc-a664-45b0-955f-88c00c68d74d', '2bbd864c-d440-46fa-8ddd-94116057b385', 59000.00, 0.00, 'IDR', 'ORD-20260404-50B5ABC2');
INSERT INTO public.ledger_lines VALUES ('64cbb518-e100-4b2f-8d30-954a4bb8c77f', '120634fc-a664-45b0-955f-88c00c68d74d', 'e47fd62c-529e-4786-b216-ac4e9f79fac9', 0.00, 35000.00, 'IDR', 'ORD-20260404-50B5ABC2');
INSERT INTO public.ledger_lines VALUES ('aef8880f-705b-42f0-88ec-ba775acf744e', '120634fc-a664-45b0-955f-88c00c68d74d', '6017dd31-2620-4665-bb23-f31c999db0bc', 0.00, 1000.00, 'IDR', 'ORD-20260404-50B5ABC2');
INSERT INTO public.ledger_lines VALUES ('07a8fe11-566f-4d39-9138-3ad80aa3283f', '120634fc-a664-45b0-955f-88c00c68d74d', '3ccfa4e4-3c56-4433-adc6-e1fefb5f0584', 0.00, 5000.00, 'IDR', 'ORD-20260404-50B5ABC2');


ALTER TABLE public.ledger_lines ENABLE TRIGGER ALL;

--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.notifications DISABLE TRIGGER ALL;



ALTER TABLE public.notifications ENABLE TRIGGER ALL;

--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.orders DISABLE TRIGGER ALL;

INSERT INTO public.orders VALUES ('c1400457-d9a5-4c1f-9684-d0fb9dbb520f', 'ORD-20260404-A4DED2FC', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', '{"label": "Alamat Keputran", "phone": "063819235476", "city_name": " SURABAYA", "address_id": "df86c486-d12c-4522-87f6-895fc25d3088", "is_default": true, "postal_code": "60265", "address_line": "Jl. Tegalsari kecamatan keputran no 123", "district_name": "TEGALSARI", "province_name": "JAWA TIMUR", "recipient_name": "Zaki", "subdistrict_name": "KEPUTRAN"}', 'pending_payment', 'unpaid', 35000.00, 20000.00, 6000.00, 61000.00, '2026-04-04 12:39:32.72551+00', '2026-04-04 12:39:32.72551+00', '2026-04-04 12:39:32.72551+00');
INSERT INTO public.orders VALUES ('2ba24abb-a8ea-4b5f-9f7c-e4ad699c9a82', 'ORD-20260404-C542C29A', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '95e22a26-5c24-44cc-b207-1455cf8656ed', '{"label": "Alamat Keputran", "phone": "063819235476", "city_name": " SURABAYA", "address_id": "df86c486-d12c-4522-87f6-895fc25d3088", "is_default": true, "postal_code": "60265", "address_line": "Jl. Tegalsari kecamatan keputran no 123", "district_name": "TEGALSARI", "province_name": "JAWA TIMUR", "recipient_name": "Zaki", "subdistrict_name": "KEPUTRAN"}', 'paid', 'paid', 340000.00, 11000.00, 6000.00, 357000.00, '2026-04-04 12:40:53.294759+00', '2026-04-04 12:40:53.294759+00', '2026-04-04 12:45:01.067627+00');
INSERT INTO public.orders VALUES ('4ac5dd9b-6039-4d95-8279-cd8437e21a2e', 'ORD-20260404-50B5ABC2', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', '{"label": "Alamat Keputran", "phone": "063819235476", "city_name": " SURABAYA", "address_id": "df86c486-d12c-4522-87f6-895fc25d3088", "is_default": true, "postal_code": "60265", "address_line": "Jl. Tegalsari kecamatan keputran no 123", "district_name": "TEGALSARI", "province_name": "JAWA TIMUR", "recipient_name": "Zaki", "subdistrict_name": "KEPUTRAN"}', 'paid', 'paid', 35000.00, 18000.00, 6000.00, 59000.00, '2026-04-04 12:38:06.650417+00', '2026-04-04 12:38:06.650417+00', '2026-04-04 12:45:35.520117+00');
INSERT INTO public.orders VALUES ('4fe7042c-0b4c-4444-8104-ecbd6fd4bc4a', 'ORD-20260404-6E8B747E', 'cf0add84-b149-40c7-ae68-fb9d02594cae', '95e22a26-5c24-44cc-b207-1455cf8656ed', '{"label": "Alamat Keputran", "phone": "063819235476", "city_name": " SURABAYA", "address_id": "df86c486-d12c-4522-87f6-895fc25d3088", "is_default": true, "postal_code": "60265", "address_line": "Jl. Tegalsari kecamatan keputran no 123", "district_name": "TEGALSARI", "province_name": "JAWA TIMUR", "recipient_name": "Zaki", "subdistrict_name": "KEPUTRAN"}', 'processing', 'paid', 170000.00, 5000.00, 6000.00, 181000.00, '2026-04-04 12:42:19.119166+00', '2026-04-04 12:42:19.119166+00', '2026-04-04 12:45:56.40852+00');


ALTER TABLE public.orders ENABLE TRIGGER ALL;

--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.order_items DISABLE TRIGGER ALL;

INSERT INTO public.order_items VALUES ('92614d10-ebd6-477d-ac7f-6571f1b11dd8', '4ac5dd9b-6039-4d95-8279-cd8437e21a2e', 'a8d80d68-2834-46ab-a90d-b8eff9eb1ec0', 'Minyak Zaitun Mustika Ratu 175ml - Melembabkan, Melembutkan, Menyegarkan', 'TOKOO-MINYA001-00', 1, 35000.00, 35000.00);
INSERT INTO public.order_items VALUES ('4263f2ad-2113-4f50-8e65-76e424a0f82a', 'c1400457-d9a5-4c1f-9684-d0fb9dbb520f', 'a8d80d68-2834-46ab-a90d-b8eff9eb1ec0', 'Minyak Zaitun Mustika Ratu 175ml - Melembabkan, Melembutkan, Menyegarkan', 'TOKOO-MINYA001-00', 1, 35000.00, 35000.00);
INSERT INTO public.order_items VALUES ('cfbbc17f-5752-4cb5-b9dc-61e02a6d0f85', '2ba24abb-a8ea-4b5f-9f7c-e4ad699c9a82', '94d6f287-58e4-4698-bc3e-f30e2e9439c8', 'Baju Musllim Lengan Panjang', 'TOKOO-BAJUM001-00', 2, 170000.00, 340000.00);
INSERT INTO public.order_items VALUES ('03bad991-7974-4907-8845-c78867fc6754', '4fe7042c-0b4c-4444-8104-ecbd6fd4bc4a', '94d6f287-58e4-4698-bc3e-f30e2e9439c8', 'Baju Musllim Lengan Panjang', 'TOKOO-BAJUM001-00', 1, 170000.00, 170000.00);


ALTER TABLE public.order_items ENABLE TRIGGER ALL;

--
-- Data for Name: order_status_history; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.order_status_history DISABLE TRIGGER ALL;

INSERT INTO public.order_status_history VALUES ('64febef3-c800-4835-a7b0-47ff9fc13687', '4ac5dd9b-6039-4d95-8279-cd8437e21a2e', NULL, 'pending_payment', NULL, '2026-04-04 12:38:06.805658+00', NULL);
INSERT INTO public.order_status_history VALUES ('d3169ec3-51fb-481c-b116-15c9ca911dcd', 'c1400457-d9a5-4c1f-9684-d0fb9dbb520f', NULL, 'pending_payment', NULL, '2026-04-04 12:39:32.874468+00', NULL);
INSERT INTO public.order_status_history VALUES ('51105594-e396-45d2-af85-3b4517e0b065', '2ba24abb-a8ea-4b5f-9f7c-e4ad699c9a82', NULL, 'pending_payment', NULL, '2026-04-04 12:40:53.34183+00', NULL);
INSERT INTO public.order_status_history VALUES ('f3216cec-a289-4e38-b1f8-d0fe99f8dec8', '4fe7042c-0b4c-4444-8104-ecbd6fd4bc4a', NULL, 'pending_payment', NULL, '2026-04-04 12:42:19.171374+00', NULL);
INSERT INTO public.order_status_history VALUES ('73d2939e-1d0d-4856-a594-f2cffbe37757', '4fe7042c-0b4c-4444-8104-ecbd6fd4bc4a', 'pending_payment', 'paid', NULL, '2026-04-04 12:44:13.477605+00', 'Payment received via BANK_TRANSFER/BNI');
INSERT INTO public.order_status_history VALUES ('83fb596f-e47a-45ec-80df-e5fb8d97d153', '2ba24abb-a8ea-4b5f-9f7c-e4ad699c9a82', 'pending_payment', 'paid', NULL, '2026-04-04 12:45:01.068136+00', 'Payment received via BANK_TRANSFER/BNI');
INSERT INTO public.order_status_history VALUES ('2e343093-1d75-4e89-b243-d2bbbadb5794', '4ac5dd9b-6039-4d95-8279-cd8437e21a2e', 'pending_payment', 'paid', NULL, '2026-04-04 12:45:35.520637+00', 'Payment received via BANK_TRANSFER/BNI');
INSERT INTO public.order_status_history VALUES ('01eb3104-b702-4ec9-b0c9-1222ff736957', '4fe7042c-0b4c-4444-8104-ecbd6fd4bc4a', 'paid', 'processing', '3af6ae21-eb2a-45ff-b2cb-5010ef5e7e5c', '2026-04-04 12:45:56.409124+00', 'Order accepted by vendor');


ALTER TABLE public.order_status_history ENABLE TRIGGER ALL;

--
-- Data for Name: otp_codes; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.otp_codes DISABLE TRIGGER ALL;

INSERT INTO public.otp_codes VALUES ('c228a11a-9894-4e39-a9c0-d47cf2435309', NULL, 'rizkyalamsyah703@gmail.com', 'email_verification', 'email', 'lYeFzbqY2ZZ85Xj7Nf1STCVoC3KnculLOBK66nYK6_U', '2026-04-04 08:55:12.13209+00', 0, '2026-04-04 08:50:40.980839+00', NULL, '2026-04-04 08:50:12.13209+00');
INSERT INTO public.otp_codes VALUES ('85a42964-37a7-410d-8691-43c47d206c54', NULL, 'rizkyalamsyah703@gmail.com', 'email_verification', 'email', 'P3mnUiEZPwQhytcnWeeCirrTz-W-ll94S5aJJZMOKXw', '2026-04-04 09:21:08.455788+00', 0, '2026-04-04 09:16:32.202321+00', NULL, '2026-04-04 09:16:08.455788+00');
INSERT INTO public.otp_codes VALUES ('3a584de3-b456-43bd-9b9b-416e0360c6e3', NULL, 'rizkyalamsyah.dev@gmail.com', 'email_verification', 'email', 'sTa6R3hxqTjigY4K2XUmxSVq8iOv3AQb4G18sjGsjUk', '2026-04-04 09:27:28.99549+00', 0, '2026-04-04 09:22:48.297728+00', NULL, '2026-04-04 09:22:28.99549+00');
INSERT INTO public.otp_codes VALUES ('aaecb669-7496-4751-8226-6c63e2a120dc', NULL, 'rizkyalamsyah.dev@gmail.com', 'email_verification', 'email', 'GxNBulHqVnfRC5w98JcAOSTSOuFPJ2haEgS38lgH3aM', '2026-04-04 09:36:44.001388+00', 1, '2026-04-04 09:32:51.709821+00', NULL, '2026-04-04 09:31:44.001388+00');
INSERT INTO public.otp_codes VALUES ('cb56d1a1-c268-41b9-9e59-d90af2650fa9', NULL, 'rizkyalamsyah703@gmail.com', 'email_verification', 'email', '0o9x0db7IO4c9A7Sc4DlCw0n-U6U4Vm04KmAVzwW5aE', '2026-04-04 09:50:28.687365+00', 0, '2026-04-04 09:45:54.772724+00', NULL, '2026-04-04 09:45:28.687365+00');
INSERT INTO public.otp_codes VALUES ('1573bf45-0de7-45ab-afde-ce4d97d9d8f1', NULL, 'rizkyalamsyah703@gmail.com', 'email_verification', 'email', 'ERUatrnU2R3bF48fc7erl_abhKk6WPo3Yxy7DqhXRJ4', '2026-04-04 10:21:25.717627+00', 0, '2026-04-04 10:16:47.75169+00', NULL, '2026-04-04 10:16:25.717627+00');
INSERT INTO public.otp_codes VALUES ('0a98a00a-cbb2-4d24-9ba6-5adc273e8f53', NULL, 'rizkyalamsyah703@gmail.com', 'email_verification', 'email', 'ygN0UPvFR6PFcZFDhD3lkbARUlMe8AAMaZwo78Gd2LY', '2026-04-04 10:40:52.037751+00', 0, '2026-04-04 10:36:14.346698+00', NULL, '2026-04-04 10:35:52.037751+00');
INSERT INTO public.otp_codes VALUES ('f4d9e3f7-30f1-4e28-82d6-913b0a3de5a5', NULL, 'r.alamsyah.8e@gmail.com', 'email_verification', 'email', 'QgygA6yMg4n0GsoD8E5CFrZUlStmWTb76b24yigRZug', '2026-04-04 10:45:13.005058+00', 1, '2026-04-04 10:40:37.589599+00', NULL, '2026-04-04 10:40:13.005058+00');
INSERT INTO public.otp_codes VALUES ('b29a4635-9ef5-43c9-a916-9d15da6b31ee', NULL, 'r.alamsyah.8e@gmail.com', 'email_verification', 'email', '6IuU9OeLrzhQXyreZeolQwNjOZ0LdFEfcz0OZLiczps', '2026-04-04 10:47:04.572612+00', 0, '2026-04-04 10:42:25.152119+00', NULL, '2026-04-04 10:42:04.572612+00');
INSERT INTO public.otp_codes VALUES ('b33298c3-3207-43ba-b4cd-3cbc9aaca822', NULL, 'r.alamsyah.8e@gmail.com', 'email_verification', 'email', 'XN9EvMy3uUeiT-cXEX2XE1DRq_tHkzYFVfylgSJXAMA', '2026-04-04 10:50:56.015415+00', 0, '2026-04-04 10:46:13.461056+00', NULL, '2026-04-04 10:45:56.015415+00');
INSERT INTO public.otp_codes VALUES ('5ee534a3-f6de-4e0c-8ca2-496e92269c69', NULL, 'r.alamsyah.8e@gmail.com', 'email_verification', 'email', 'xbPK0sSFEYuQxA9lTdcCVBT5CkZDs1GLJyG4jVOg3PM', '2026-04-04 10:55:18.747366+00', 0, '2026-04-04 10:50:35.547753+00', NULL, '2026-04-04 10:50:18.747366+00');
INSERT INTO public.otp_codes VALUES ('e657b266-337d-48cd-bc8f-8c273b29d1ce', NULL, 'r.alamsyah.8e@gmail.com', 'email_verification', 'email', 'evf-KS6a1QOAwR6K23Rmq9xsydb6H0RMlAJD6TqvXBI', '2026-04-04 10:59:21.391233+00', 0, '2026-04-04 10:54:39.557142+00', NULL, '2026-04-04 10:54:21.391233+00');
INSERT INTO public.otp_codes VALUES ('fb25e1ec-ead4-465d-9386-5bef0c8857dc', NULL, 'r.alamsyah.8e@gmail.com', 'email_verification', 'email', 'EAijpv3hTWtPLJCax490OfifywwqOr_zDa7VyDxTp4Q', '2026-04-04 11:02:36.290171+00', 0, '2026-04-04 10:57:58.462527+00', NULL, '2026-04-04 10:57:36.290171+00');
INSERT INTO public.otp_codes VALUES ('f84a347a-dc6d-4681-88a2-af9fa03aa9aa', NULL, 'rizkyalamsyah.dev@gmail.com', 'email_verification', 'email', 'ibKvJGSki6WAsbW9WyPxgpcejGTr7vZome3B2Z9Xj7Q', '2026-04-04 12:02:54.640239+00', 0, '2026-04-04 11:58:09.092689+00', NULL, '2026-04-04 11:57:54.640239+00');
INSERT INTO public.otp_codes VALUES ('b7b11344-ef25-43ae-841a-65a33e1e7076', NULL, '404lamfound@gmail.com', 'email_verification', 'email', 'xLg2Gv0QK6taSeSmfq5qKW9DLBpt_SsrxVN7phwZG6o', '2026-04-04 12:09:44.512268+00', 0, NULL, NULL, '2026-04-04 12:04:44.512268+00');
INSERT INTO public.otp_codes VALUES ('891a1f29-5417-4306-89be-04433d3d1957', NULL, 'nobev17858@nexafilm.com', 'email_verification', 'email', 'uI6pnDFSNxWn7AEYfarTU7Zbk9hKYjua5YfGoWhugMc', '2026-04-04 12:10:13.669755+00', 0, '2026-04-04 12:05:33.821398+00', NULL, '2026-04-04 12:05:13.669755+00');


ALTER TABLE public.otp_codes ENABLE TRIGGER ALL;

--
-- Data for Name: payment_invoices; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.payment_invoices DISABLE TRIGGER ALL;

INSERT INTO public.payment_invoices VALUES ('0be69b7d-eafa-48f6-8d6d-c36379d80dec', 'c1400457-d9a5-4c1f-9684-d0fb9dbb520f', 'xendit', '69d106843263d1649f21810c', 'INV-ORD-20260404-A4DED2FC-e12c2f1f', 'https://checkout-staging.xendit.co/web/69d106843263d1649f21810c', NULL, NULL, 61000.00, 'IDR', 'pending', '2026-04-05 12:39:32.939+00', NULL, NULL, '2026-04-04 12:39:32.72551+00', '2026-04-04 12:39:32.72551+00');
INSERT INTO public.payment_invoices VALUES ('c746ecd5-cb8d-45b3-a6f2-9a8ecffbd1de', '4fe7042c-0b4c-4444-8104-ecbd6fd4bc4a', 'xendit', '69d1072b1c934dec38d3e15f', 'INV-ORD-20260404-6E8B747E-28f562dd', 'https://checkout-staging.xendit.co/web/69d1072b1c934dec38d3e15f', 'BANK_TRANSFER', 'BNI', 181000.00, 'IDR', 'paid', '2026-04-05 12:42:19.239+00', '2026-04-04 12:42:24+00', '{"webhook_payload": {"id": "69d1072b1c934dec38d3e15f", "amount": 181000, "status": "PAID", "paid_at": "2026-04-04T12:42:24.000Z", "user_id": "69d1003545d4336498d0da87", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260404-6E8B747E-28f562dd", "paid_amount": 181000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}}', '2026-04-04 12:42:19.119166+00', '2026-04-04 12:44:13.474165+00');
INSERT INTO public.payment_invoices VALUES ('6444f73a-49aa-4dd1-b31c-cb352545dbf6', '2ba24abb-a8ea-4b5f-9f7c-e4ad699c9a82', 'xendit', '69d106d53263d1649f21818b', 'INV-ORD-20260404-C542C29A-b2d04aa9', 'https://checkout-staging.xendit.co/web/69d106d53263d1649f21818b', 'BANK_TRANSFER', 'BNI', 357000.00, 'IDR', 'paid', '2026-04-05 12:40:55.879+00', '2026-04-04 12:41:05+00', '{"webhook_payload": {"id": "69d106d53263d1649f21818b", "amount": 357000, "status": "PAID", "paid_at": "2026-04-04T12:41:05.000Z", "user_id": "69d1003545d4336498d0da87", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260404-C542C29A-b2d04aa9", "paid_amount": 357000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}}', '2026-04-04 12:40:53.294759+00', '2026-04-04 12:45:01.065264+00');
INSERT INTO public.payment_invoices VALUES ('0fa1f7a7-c377-4c35-b77a-1b1d38b9c28c', '4ac5dd9b-6039-4d95-8279-cd8437e21a2e', 'xendit', '69d1062e1c934dec38d3dfcf', 'INV-ORD-20260404-50B5ABC2-98b64360', 'https://checkout-staging.xendit.co/web/69d1062e1c934dec38d3dfcf', 'BANK_TRANSFER', 'BNI', 59000.00, 'IDR', 'paid', '2026-04-05 12:38:09.395+00', '2026-04-04 12:38:25+00', '{"webhook_payload": {"id": "69d1062e1c934dec38d3dfcf", "amount": 59000, "status": "PAID", "paid_at": "2026-04-04T12:38:25.000Z", "user_id": "69d10022d4417937e489e76b", "currency": "IDR", "metadata": null, "external_id": "INV-ORD-20260404-50B5ABC2-98b64360", "paid_amount": 59000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}}', '2026-04-04 12:38:06.650417+00', '2026-04-04 12:45:35.51779+00');


ALTER TABLE public.payment_invoices ENABLE TRIGGER ALL;

--
-- Data for Name: payment_events; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.payment_events DISABLE TRIGGER ALL;

INSERT INTO public.payment_events VALUES ('cd80621b-7ccf-4e16-9406-a6fe74dc038f', 'c746ecd5-cb8d-45b3-a6f2-9a8ecffbd1de', 'invoice.paid', '69d1072b1c934dec38d3e15f:PAID', '{"id": "69d1072b1c934dec38d3e15f", "amount": 181000, "status": "PAID", "paid_at": "2026-04-04T12:42:24.000Z", "external_id": "INV-ORD-20260404-6E8B747E-28f562dd", "paid_amount": 181000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}', '2026-04-04 12:44:13.471222+00');
INSERT INTO public.payment_events VALUES ('821cb106-d1cd-46aa-ba46-e00d4ba3e8da', '6444f73a-49aa-4dd1-b31c-cb352545dbf6', 'invoice.paid', '69d106d53263d1649f21818b:PAID', '{"id": "69d106d53263d1649f21818b", "amount": 357000, "status": "PAID", "paid_at": "2026-04-04T12:41:05.000Z", "external_id": "INV-ORD-20260404-C542C29A-b2d04aa9", "paid_amount": 357000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}', '2026-04-04 12:45:01.062719+00');
INSERT INTO public.payment_events VALUES ('5aa3a071-6699-42bf-a24a-0c2a15920728', '0fa1f7a7-c377-4c35-b77a-1b1d38b9c28c', 'invoice.paid', '69d1062e1c934dec38d3dfcf:PAID', '{"id": "69d1062e1c934dec38d3dfcf", "amount": 59000, "status": "PAID", "paid_at": "2026-04-04T12:38:25.000Z", "external_id": "INV-ORD-20260404-50B5ABC2-98b64360", "paid_amount": 59000, "payment_method": "BANK_TRANSFER", "payment_channel": "BNI"}', '2026-04-04 12:45:35.515224+00');


ALTER TABLE public.payment_events ENABLE TRIGGER ALL;

--
-- Data for Name: payout_batches; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.payout_batches DISABLE TRIGGER ALL;



ALTER TABLE public.payout_batches ENABLE TRIGGER ALL;

--
-- Data for Name: payout_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.payout_items DISABLE TRIGGER ALL;



ALTER TABLE public.payout_items ENABLE TRIGGER ALL;

--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.product_images DISABLE TRIGGER ALL;

INSERT INTO public.product_images VALUES ('fc5f4ec5-5b40-4ae9-b262-a55026a4de77', '9b2b9956-3d69-4c05-98bf-c87b03fec86d', 'products/9b2b9956-3d69-4c05-98bf-c87b03fec86d/images/fc5f4ec5-5b40-4ae9-b262-a55026a4de77/kurma-ajwa.jpg', 'image/jpeg', 11930, true, 0, '2026-04-04 12:18:30.516593+00');
INSERT INTO public.product_images VALUES ('248541aa-a53c-462c-a42b-13f3e94573ab', 'f0d3da47-3a75-43b1-8a5d-8e18d00f4a7d', 'products/f0d3da47-3a75-43b1-8a5d-8e18d00f4a7d/images/248541aa-a53c-462c-a42b-13f3e94573ab/baju-muslim-lpj-sage.jpg', 'image/jpeg', 89740, true, 0, '2026-04-04 12:22:42.65505+00');
INSERT INTO public.product_images VALUES ('564ee406-1e52-4999-b947-330a936af94b', 'cd83d2ab-b6cc-4b37-989c-b732beb8a897', 'products/cd83d2ab-b6cc-4b37-989c-b732beb8a897/images/564ee406-1e52-4999-b947-330a936af94b/airzamzam', 'image/jpeg', 578109, true, 0, '2026-04-04 12:31:46.354883+00');
INSERT INTO public.product_images VALUES ('fae6e64e-ead5-436a-8185-658ccf7f5b38', '9f5bfde2-09b4-44c3-b080-a17099d9a797', 'products/9f5bfde2-09b4-44c3-b080-a17099d9a797/images/fae6e64e-ead5-436a-8185-658ccf7f5b38/airzamzam', 'image/webp', 60818, true, 0, '2026-04-04 12:33:39.236631+00');


ALTER TABLE public.product_images ENABLE TRIGGER ALL;

--
-- Data for Name: product_promotions; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.product_promotions DISABLE TRIGGER ALL;



ALTER TABLE public.product_promotions ENABLE TRIGGER ALL;

--
-- Data for Name: product_reviews; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.product_reviews DISABLE TRIGGER ALL;



ALTER TABLE public.product_reviews ENABLE TRIGGER ALL;

--
-- Data for Name: product_review_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.product_review_images DISABLE TRIGGER ALL;



ALTER TABLE public.product_review_images ENABLE TRIGGER ALL;

--
-- Data for Name: product_review_stats; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.product_review_stats DISABLE TRIGGER ALL;



ALTER TABLE public.product_review_stats ENABLE TRIGGER ALL;

--
-- Data for Name: return_reasons; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.return_reasons DISABLE TRIGGER ALL;



ALTER TABLE public.return_reasons ENABLE TRIGGER ALL;

--
-- Data for Name: user_bank_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.user_bank_accounts DISABLE TRIGGER ALL;



ALTER TABLE public.user_bank_accounts ENABLE TRIGGER ALL;

--
-- Data for Name: refunds; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.refunds DISABLE TRIGGER ALL;



ALTER TABLE public.refunds ENABLE TRIGGER ALL;

--
-- Data for Name: refund_evidences; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.refund_evidences DISABLE TRIGGER ALL;



ALTER TABLE public.refund_evidences ENABLE TRIGGER ALL;

--
-- Data for Name: reply_templates; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.reply_templates DISABLE TRIGGER ALL;



ALTER TABLE public.reply_templates ENABLE TRIGGER ALL;

--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.schema_migrations DISABLE TRIGGER ALL;

INSERT INTO public.schema_migrations VALUES (60, false) ON CONFLICT DO NOTHING;


ALTER TABLE public.schema_migrations ENABLE TRIGGER ALL;

--
-- Data for Name: shipments; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.shipments DISABLE TRIGGER ALL;

INSERT INTO public.shipments VALUES ('d769f77e-f471-4abf-a2a4-c866fa40eee8', '4ac5dd9b-6039-4d95-8279-cd8437e21a2e', 'sicepat', 'REG', '', 'waiting_pickup', NULL, NULL, '2-3 day');
INSERT INTO public.shipments VALUES ('6105f76a-a518-4763-aa78-d077f3ed0815', 'c1400457-d9a5-4c1f-9684-d0fb9dbb520f', 'jne', 'REG', '', 'waiting_pickup', NULL, NULL, '1 day');
INSERT INTO public.shipments VALUES ('cd96dc03-b831-47c8-a112-d6da4f530e18', '2ba24abb-a8ea-4b5f-9f7c-e4ad699c9a82', 'sicepat', 'BEST', '', 'waiting_pickup', NULL, NULL, '1 day');
INSERT INTO public.shipments VALUES ('9adab346-3fca-4a7e-bb81-bdc9b95dc1a6', '4fe7042c-0b4c-4444-8104-ecbd6fd4bc4a', 'wahana', 'Ekonomis', '', 'waiting_pickup', NULL, NULL, '2 day');


ALTER TABLE public.shipments ENABLE TRIGGER ALL;

--
-- Data for Name: ticket_subjects; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.ticket_subjects DISABLE TRIGGER ALL;

INSERT INTO public.ticket_subjects VALUES (1, 'Barang Tidak Sampai', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (2, 'Barang Rusak', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (3, 'Pengembalian Dana', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (4, 'Barang Tidak Sesuai', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (5, 'Pengiriman Terlambat', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (6, 'Pembatalan Pesanan', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (7, 'Kesalahan Produk', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (8, 'Akun Bermasalah', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (9, 'Pembayaran Gagal', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (10, 'Voucher/Promo Tidak Berlaku', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (11, 'Pertanyaan Umum', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;
INSERT INTO public.ticket_subjects VALUES (12, 'Lainnya', true, '2026-04-04 08:43:39.364438+00') ON CONFLICT DO NOTHING;


ALTER TABLE public.ticket_subjects ENABLE TRIGGER ALL;

--
-- Data for Name: tickets; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.tickets DISABLE TRIGGER ALL;



ALTER TABLE public.tickets ENABLE TRIGGER ALL;

--
-- Data for Name: ticket_messages; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.ticket_messages DISABLE TRIGGER ALL;



ALTER TABLE public.ticket_messages ENABLE TRIGGER ALL;

--
-- Data for Name: ticket_attachments; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.ticket_attachments DISABLE TRIGGER ALL;



ALTER TABLE public.ticket_attachments ENABLE TRIGGER ALL;

--
-- Data for Name: ticket_status_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.ticket_status_logs DISABLE TRIGGER ALL;



ALTER TABLE public.ticket_status_logs ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_balances; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_balances DISABLE TRIGGER ALL;

INSERT INTO public.vendor_balances VALUES ('95e22a26-5c24-44cc-b207-1455cf8656ed', 0.00, 0.00, 0.00, 0.00, '2026-04-04 12:01:26.504538+00', 0.00);
INSERT INTO public.vendor_balances VALUES ('0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 0.00, 0.00, 0.00, 0.00, '2026-04-04 12:05:58.662114+00', 0.00);


ALTER TABLE public.vendor_balances ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_bank_accounts; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_bank_accounts DISABLE TRIGGER ALL;



ALTER TABLE public.vendor_bank_accounts ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_banners; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_banners DISABLE TRIGGER ALL;



ALTER TABLE public.vendor_banners ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_couriers; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_couriers DISABLE TRIGGER ALL;

INSERT INTO public.vendor_couriers VALUES ('cdf8b04b-3aaa-4ed5-a1e8-97a65a38f8ae', '95e22a26-5c24-44cc-b207-1455cf8656ed', 6, true, '2026-04-04 12:17:36.790144+00');
INSERT INTO public.vendor_couriers VALUES ('2b9863b7-967d-490f-9146-326b5cb793c0', '95e22a26-5c24-44cc-b207-1455cf8656ed', 1, true, '2026-04-04 12:17:36.791111+00');
INSERT INTO public.vendor_couriers VALUES ('5a6970ce-32b3-4bb6-be21-365480999783', '95e22a26-5c24-44cc-b207-1455cf8656ed', 2, true, '2026-04-04 12:17:36.791735+00');
INSERT INTO public.vendor_couriers VALUES ('2eba7e35-0b6e-4fed-8b24-05ff340b0d44', '95e22a26-5c24-44cc-b207-1455cf8656ed', 8, true, '2026-04-04 12:17:36.792193+00');
INSERT INTO public.vendor_couriers VALUES ('05a469b3-1ada-4348-808f-221386d561f9', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 6, true, '2026-04-04 12:29:11.177617+00');
INSERT INTO public.vendor_couriers VALUES ('918b68e0-2612-414b-b735-77f06291a381', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 1, true, '2026-04-04 12:29:11.178222+00');
INSERT INTO public.vendor_couriers VALUES ('633b3771-8cce-4fb5-9dc3-ff11f0902c97', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 5, true, '2026-04-04 12:29:11.178567+00');
INSERT INTO public.vendor_couriers VALUES ('eb739033-306c-41db-a930-e8420ef6e4a2', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 9, true, '2026-04-04 12:29:11.178948+00');
INSERT INTO public.vendor_couriers VALUES ('976eeeb6-c2a8-42b7-b1f0-696ec5bb63c7', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 2, true, '2026-04-04 12:29:11.179378+00');


ALTER TABLE public.vendor_couriers ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_documents; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_documents DISABLE TRIGGER ALL;

INSERT INTO public.vendor_documents VALUES ('90ffcb9d-c480-4f45-90a0-2c598048799e', '95e22a26-5c24-44cc-b207-1455cf8656ed', 'owner_document_id', 'vendor-onboardings/f72462ae-7926-492e-b041-96a1dd0aaeec/documents/owner_document_id/389c7aab-7df3-4597-86e9-c7080ed6228d', 'application/pdf', 2212, NULL, '3af6ae21-eb2a-45ff-b2cb-5010ef5e7e5c', NULL, NULL, '2026-04-04 12:01:26.504538+00', '2026-04-04 12:01:26.504538+00');
INSERT INTO public.vendor_documents VALUES ('4e86aae1-ced0-435f-9b36-8615ad4458cf', '0a12a4f5-2fa5-4bfb-8e7a-ac457a21d4ee', 'owner_document_id', 'vendor-onboardings/39cc684e-2f2e-4421-aada-b07c86c65fe8/documents/owner_document_id/ccb2d7d0-57b8-4228-af4b-be14b2ab8147', 'application/pdf', 2212, NULL, '8e3a1e10-8604-4dc4-87db-838cb5bea0bd', NULL, NULL, '2026-04-04 12:05:58.662114+00', '2026-04-04 12:05:58.662114+00');


ALTER TABLE public.vendor_documents ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_onboardings; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_onboardings DISABLE TRIGGER ALL;

INSERT INTO public.vendor_onboardings (id, email, status, password_hash, otp_verified_at, completed_at, created_at, updated_at) VALUES ('5d3b0319-4987-4cf7-9cf7-bb6831b3249b', 'rizkyalamsyah703@gmail.com', 'otp_verified', '$2a$10$.WXd0AO/nEzehmSg3zYtDe0wmQsfOs58k.mef7s.D1J8yRXV8m1Ve', '2026-04-04 10:36:14.348744+00', NULL, '2026-04-04 08:50:40.983267+00', '2026-04-04 10:36:16.745772+00');
INSERT INTO public.vendor_onboardings (id, email, status, password_hash, otp_verified_at, completed_at, created_at, updated_at) VALUES ('07cd6182-c145-44f3-a624-5bc50d6b6a5e', 'r.alamsyah.8e@gmail.com', 'otp_verified', '$2a$10$rBAtdeTH3xUZTKXsqGrLvOs.qwtQ/BAEo.B8BT9tZKaFqNIxMlrsy', '2026-04-04 10:57:58.464594+00', NULL, '2026-04-04 10:40:37.591553+00', '2026-04-04 10:57:59.838004+00');
INSERT INTO public.vendor_onboardings (id, email, status, password_hash, otp_verified_at, completed_at, created_at, updated_at) VALUES ('f72462ae-7926-492e-b041-96a1dd0aaeec', 'rizkyalamsyah.dev@gmail.com', 'completed', NULL, '2026-04-04 11:58:09.094933+00', '2026-04-04 12:01:26.504538+00', '2026-04-04 09:22:48.299603+00', '2026-04-04 12:01:26.504538+00');
INSERT INTO public.vendor_onboardings (id, email, status, password_hash, otp_verified_at, completed_at, created_at, updated_at) VALUES ('39cc684e-2f2e-4421-aada-b07c86c65fe8', 'nobev17858@nexafilm.com', 'completed', NULL, '2026-04-04 12:05:33.823588+00', '2026-04-04 12:05:58.662114+00', '2026-04-04 12:05:33.823588+00', '2026-04-04 12:05:58.662114+00');


ALTER TABLE public.vendor_onboardings ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_vouchers; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_vouchers DISABLE TRIGGER ALL;



ALTER TABLE public.vendor_vouchers ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_voucher_products; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_voucher_products DISABLE TRIGGER ALL;



ALTER TABLE public.vendor_voucher_products ENABLE TRIGGER ALL;

--
-- Data for Name: vendor_withdrawals; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.vendor_withdrawals DISABLE TRIGGER ALL;



ALTER TABLE public.vendor_withdrawals ENABLE TRIGGER ALL;

--
-- Data for Name: wishlist_items; Type: TABLE DATA; Schema: public; Owner: postgres
--

ALTER TABLE public.wishlist_items DISABLE TRIGGER ALL;



ALTER TABLE public.wishlist_items ENABLE TRIGGER ALL;

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
SET session_replication_role = DEFAULT;


