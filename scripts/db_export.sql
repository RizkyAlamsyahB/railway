--
-- PostgreSQL database dump
--

\restrict bIF4ooiSRpwcf5IUMzkjoAvUzVDvYTG2iw8yd3ryvByArb8IkxhtYzXuRbl5Bog

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
-- Name: enforce_product_image_limit(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.enforce_product_image_limit() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
DECLARE
    current_count INT;
BEGIN
    SELECT COUNT(*)
    INTO current_count
    FROM product_images pi
    WHERE pi.product_id = NEW.product_id
      AND (TG_OP <> 'UPDATE' OR pi.id <> OLD.id);

    IF current_count >= 10 THEN
        RAISE EXCEPTION 'maximum 10 images per product is allowed';
    END IF;

    RETURN NEW;
END;
$$;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.addresses (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    label character varying(40),
    recipient_name character varying(120),
    phone character varying(20),
    province_id character varying(20),
    city_id character varying(20),
    district_id character varying(20),
    postal_code character varying(10),
    address_line text NOT NULL,
    is_default boolean DEFAULT false NOT NULL,
    province_name character varying(100),
    city_name character varying(100),
    district_name character varying(100),
    subdistrict_id character varying(20),
    subdistrict_name character varying(100),
    notes text,
    latitude double precision,
    longitude double precision,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: admin_contacts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.admin_contacts (
    id integer NOT NULL,
    content text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: admin_contacts_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.admin_contacts_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: admin_contacts_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.admin_contacts_id_seq OWNED BY public.admin_contacts.id;


--
-- Name: banners; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.banners (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(255) NOT NULL,
    image_url text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: cart_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cart_items (
    id uuid NOT NULL,
    cart_id uuid NOT NULL,
    product_variant_id uuid NOT NULL,
    qty integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_cart_items_qty_positive CHECK ((qty > 0))
);


--
-- Name: carts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.carts (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    status character varying(16) DEFAULT 'active'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_carts_status CHECK (((status)::text = ANY ((ARRAY['active'::character varying, 'converted'::character varying, 'abandoned'::character varying])::text[])))
);


--
-- Name: categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.categories (
    id uuid NOT NULL,
    parent_id uuid,
    name character varying(80) NOT NULL,
    slug character varying(100) NOT NULL,
    is_active boolean DEFAULT true NOT NULL
);


--
-- Name: chat_conversations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_conversations (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    initiator_id uuid NOT NULL,
    participant_id uuid NOT NULL,
    last_message_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_chat_conv_not_self CHECK ((initiator_id <> participant_id))
);


--
-- Name: chat_messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.chat_messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    conversation_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    message text NOT NULL,
    attachment_url text,
    is_read boolean DEFAULT false NOT NULL,
    read_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: couriers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.couriers (
    id integer NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(100) NOT NULL,
    logo_url character varying(255),
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: couriers_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.couriers_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: couriers_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.couriers_id_seq OWNED BY public.couriers.id;


--
-- Name: email_verification_tokens; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.email_verification_tokens (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    email character varying(255) NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    consumed_at timestamp with time zone,
    invalidated_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: faqs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.faqs (
    id integer NOT NULL,
    category character varying(60) NOT NULL,
    question text NOT NULL,
    answer text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: faqs_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.faqs_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: faqs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.faqs_id_seq OWNED BY public.faqs.id;


--
-- Name: ledger_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_accounts (
    id uuid NOT NULL,
    code character varying(20) NOT NULL,
    name character varying(120) NOT NULL,
    account_type character varying(16) NOT NULL,
    normal_side character varying(1) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    CONSTRAINT ck_ledger_accounts_account_type CHECK (((account_type)::text = ANY ((ARRAY['asset'::character varying, 'liability'::character varying, 'equity'::character varying, 'revenue'::character varying, 'expense'::character varying])::text[]))),
    CONSTRAINT ck_ledger_accounts_normal_side CHECK (((normal_side)::text = ANY ((ARRAY['D'::character varying, 'C'::character varying])::text[])))
);


--
-- Name: ledger_journals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_journals (
    id uuid NOT NULL,
    journal_no character varying(40) NOT NULL,
    source_type character varying(24) NOT NULL,
    source_id uuid NOT NULL,
    event_time timestamp with time zone NOT NULL,
    description text,
    status character varying(16) DEFAULT 'posted'::character varying NOT NULL,
    created_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_ledger_journals_source_type CHECK (((source_type)::text = ANY ((ARRAY['order'::character varying, 'payment_invoice'::character varying, 'refund'::character varying, 'payout_batch'::character varying, 'manual'::character varying])::text[]))),
    CONSTRAINT ck_ledger_journals_status CHECK (((status)::text = ANY ((ARRAY['posted'::character varying, 'reversed'::character varying])::text[])))
);


--
-- Name: ledger_lines; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ledger_lines (
    id uuid NOT NULL,
    journal_id uuid NOT NULL,
    account_id uuid NOT NULL,
    debit numeric(18,2) DEFAULT 0 NOT NULL,
    credit numeric(18,2) DEFAULT 0 NOT NULL,
    currency character(3) DEFAULT 'IDR'::bpchar NOT NULL,
    reference character varying(80)
);


--
-- Name: notifications; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.notifications (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    type character varying(32) NOT NULL,
    title character varying(255) NOT NULL,
    message text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    read_at timestamp with time zone,
    CONSTRAINT ck_notifications_type CHECK (((type)::text = ANY ((ARRAY['chat'::character varying, 'ticket'::character varying, 'vendor_product'::character varying, 'order'::character varying, 'refund'::character varying, 'payout'::character varying, 'system'::character varying])::text[])))
);


--
-- Name: order_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.order_items (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    product_variant_id uuid NOT NULL,
    product_name_snapshot character varying(180) NOT NULL,
    sku_snapshot character varying(80) NOT NULL,
    qty integer NOT NULL,
    unit_price numeric(18,2) NOT NULL,
    line_total numeric(18,2) NOT NULL,
    CONSTRAINT ck_order_items_qty_positive CHECK ((qty > 0))
);


--
-- Name: order_status_history; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.order_status_history (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    old_status character varying(24),
    new_status character varying(24) NOT NULL,
    changed_by uuid,
    changed_at timestamp with time zone DEFAULT now() NOT NULL,
    notes text,
    CONSTRAINT ck_order_status_history_new_status CHECK (((new_status)::text = ANY ((ARRAY['pending_payment'::character varying, 'paid'::character varying, 'packed'::character varying, 'shipped'::character varying, 'completed'::character varying, 'canceled'::character varying, 'refunded'::character varying])::text[]))),
    CONSTRAINT ck_order_status_history_old_status CHECK (((old_status IS NULL) OR ((old_status)::text = ANY ((ARRAY['pending_payment'::character varying, 'paid'::character varying, 'packed'::character varying, 'shipped'::character varying, 'completed'::character varying, 'canceled'::character varying, 'refunded'::character varying])::text[]))))
);


--
-- Name: orders; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.orders (
    id uuid NOT NULL,
    order_no character varying(30) NOT NULL,
    user_id uuid NOT NULL,
    vendor_id uuid NOT NULL,
    shipping_address_snapshot jsonb NOT NULL,
    order_status character varying(24) DEFAULT 'pending_payment'::character varying NOT NULL,
    payment_status character varying(16) DEFAULT 'unpaid'::character varying NOT NULL,
    subtotal numeric(18,2) NOT NULL,
    shipping_fee numeric(18,2) DEFAULT 0 NOT NULL,
    platform_fee numeric(18,2) DEFAULT 0 NOT NULL,
    grand_total numeric(18,2) NOT NULL,
    placed_at timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_orders_order_status CHECK (((order_status)::text = ANY ((ARRAY['pending_payment'::character varying, 'paid'::character varying, 'packed'::character varying, 'shipped'::character varying, 'completed'::character varying, 'canceled'::character varying, 'refunded'::character varying])::text[]))),
    CONSTRAINT ck_orders_payment_status CHECK (((payment_status)::text = ANY ((ARRAY['unpaid'::character varying, 'paid_partial'::character varying, 'paid'::character varying, 'refunded'::character varying])::text[])))
);


--
-- Name: payment_events; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payment_events (
    id uuid NOT NULL,
    payment_invoice_id uuid NOT NULL,
    event_type character varying(50) NOT NULL,
    external_event_id character varying(120) NOT NULL,
    payload jsonb NOT NULL,
    received_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: payment_invoices; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payment_invoices (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    gateway character varying(20) DEFAULT 'xendit'::character varying NOT NULL,
    xendit_invoice_id character varying(100),
    external_invoice_id character varying(100) NOT NULL,
    invoice_url text,
    payment_method character varying(30),
    payment_channel character varying(30),
    amount numeric(18,2) NOT NULL,
    currency character(3) DEFAULT 'IDR'::bpchar NOT NULL,
    status character varying(16) DEFAULT 'pending'::character varying NOT NULL,
    expires_at timestamp with time zone,
    paid_at timestamp with time zone,
    raw_payload jsonb,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_payment_invoices_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'paid'::character varying, 'expired'::character varying, 'failed'::character varying])::text[])))
);


--
-- Name: payout_batches; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payout_batches (
    id uuid NOT NULL,
    vendor_id uuid NOT NULL,
    period_start date NOT NULL,
    period_end date NOT NULL,
    status character varying(16) DEFAULT 'draft'::character varying NOT NULL,
    total_gross numeric(18,2) NOT NULL,
    total_fee numeric(18,2) NOT NULL,
    total_net numeric(18,2) NOT NULL,
    paid_at timestamp with time zone,
    created_by uuid NOT NULL,
    xendit_payout_id character varying(128),
    channel_code character varying(32),
    description text,
    xendit_status character varying(32),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_payout_batches_status CHECK (((status)::text = ANY ((ARRAY['schedule'::character varying, 'completed'::character varying, 'failed'::character varying, 'on_hold'::character varying])::text[]))),
    CONSTRAINT ck_payout_batches_xendit_status CHECK (((xendit_status)::text = ANY ((ARRAY['ready'::character varying, 'schedule'::character varying, 'complete'::character varying, 'failed'::character varying, 'on hold'::character varying])::text[])))
);


--
-- Name: payout_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.payout_items (
    id uuid NOT NULL,
    payout_batch_id uuid NOT NULL,
    order_id uuid NOT NULL,
    gross_amount numeric(18,2) NOT NULL,
    platform_fee_amount numeric(18,2) NOT NULL,
    net_amount numeric(18,2) NOT NULL
);


--
-- Name: product_images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_images (
    id uuid NOT NULL,
    product_id uuid NOT NULL,
    image_url text NOT NULL,
    mime_type character varying(100),
    file_size_bytes integer,
    is_primary boolean DEFAULT false NOT NULL,
    sort_order integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_product_images_sort_order_non_negative CHECK ((sort_order >= 0))
);


--
-- Name: product_review_images; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_review_images (
    id uuid NOT NULL,
    review_id uuid NOT NULL,
    object_key text NOT NULL,
    mime_type character varying(100),
    file_size_bytes integer,
    sort_order integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_product_review_images_sort_non_negative CHECK ((sort_order >= 0))
);


--
-- Name: product_review_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_review_stats (
    product_id uuid NOT NULL,
    total_reviews bigint DEFAULT 0 NOT NULL,
    total_stars bigint DEFAULT 0 NOT NULL,
    average_rating numeric(4,2) DEFAULT 0 NOT NULL,
    star_0_count bigint DEFAULT 0 NOT NULL,
    star_1_count bigint DEFAULT 0 NOT NULL,
    star_2_count bigint DEFAULT 0 NOT NULL,
    star_3_count bigint DEFAULT 0 NOT NULL,
    star_4_count bigint DEFAULT 0 NOT NULL,
    star_5_count bigint DEFAULT 0 NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_product_review_stats_non_negative CHECK (((total_reviews >= 0) AND (total_stars >= 0) AND (star_0_count >= 0) AND (star_1_count >= 0) AND (star_2_count >= 0) AND (star_3_count >= 0) AND (star_4_count >= 0) AND (star_5_count >= 0)))
);


--
-- Name: product_reviews; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_reviews (
    id uuid NOT NULL,
    product_id uuid NOT NULL,
    order_id uuid NOT NULL,
    order_item_id uuid NOT NULL,
    user_id uuid NOT NULL,
    rating smallint NOT NULL,
    review_text text NOT NULL,
    status character varying(16) DEFAULT 'published'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_product_reviews_rating_range CHECK (((rating >= 1) AND (rating <= 5))),
    CONSTRAINT ck_product_reviews_status CHECK (((status)::text = ANY ((ARRAY['published'::character varying, 'hidden'::character varying])::text[])))
);


--
-- Name: product_variants; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.product_variants (
    id uuid NOT NULL,
    product_id uuid NOT NULL,
    sku character varying(80) NOT NULL,
    variant_name character varying(120) NOT NULL,
    price numeric(18,2) NOT NULL,
    currency character(3) DEFAULT 'IDR'::bpchar NOT NULL,
    stock_on_hand integer NOT NULL,
    weight_gram integer,
    is_default boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    CONSTRAINT ck_product_variants_price_positive CHECK ((price > (0)::numeric)),
    CONSTRAINT ck_product_variants_stock_non_negative CHECK ((stock_on_hand >= 0)),
    CONSTRAINT ck_product_variants_weight_non_negative CHECK (((weight_gram IS NULL) OR (weight_gram >= 0)))
);


--
-- Name: products; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.products (
    id uuid NOT NULL,
    vendor_id uuid NOT NULL,
    category_id uuid NOT NULL,
    name character varying(180) NOT NULL,
    slug character varying(220) NOT NULL,
    description text NOT NULL,
    status character varying(16) DEFAULT 'draft'::character varying NOT NULL,
    halal_ai_status character varying(16) DEFAULT 'pending'::character varying NOT NULL,
    halal_ai_notes text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_products_halal_ai_status CHECK (((halal_ai_status)::text = ANY ((ARRAY['pending'::character varying, 'passed'::character varying, 'failed'::character varying])::text[]))),
    CONSTRAINT ck_products_status CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'published'::character varying, 'blocked'::character varying, 'archived'::character varying])::text[])))
);


--
-- Name: refunds; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.refunds (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    payment_invoice_id uuid NOT NULL,
    amount numeric(18,2) NOT NULL,
    reason text,
    status character varying(16) DEFAULT 'requested'::character varying NOT NULL,
    requested_by uuid NOT NULL,
    processed_by uuid,
    processed_at timestamp with time zone,
    CONSTRAINT ck_refunds_amount_positive CHECK ((amount > (0)::numeric)),
    CONSTRAINT ck_refunds_status CHECK (((status)::text = ANY ((ARRAY['requested'::character varying, 'approved'::character varying, 'rejected'::character varying, 'processed'::character varying])::text[])))
);


--
-- Name: reply_templates; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.reply_templates (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying(150) NOT NULL,
    content text NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_by uuid NOT NULL,
    updated_by uuid,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    shortcut character varying(100) NOT NULL,
    category character varying(50) NOT NULL,
    CONSTRAINT ck_reply_templates_category CHECK (((category)::text = ANY ((ARRAY['general'::character varying, 'order'::character varying, 'payment'::character varying, 'refund'::character varying, 'shipping'::character varying, 'product'::character varying, 'account'::character varying, 'complaint'::character varying, 'return'::character varying, 'promo'::character varying, 'voucher'::character varying, 'technical'::character varying, 'verification'::character varying, 'vendor'::character varying, 'stock'::character varying, 'cancellation'::character varying, 'delivery'::character varying, 'pickup'::character varying, 'subscription'::character varying, 'review'::character varying, 'fraud'::character varying, 'other'::character varying])::text[])))
);


--
-- Name: return_reasons; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.return_reasons (
    id integer NOT NULL,
    reason character varying(255) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: return_reasons_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.return_reasons_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: return_reasons_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.return_reasons_id_seq OWNED BY public.return_reasons.id;


--
-- Name: roles; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.roles (
    id smallint NOT NULL,
    code character varying(30) NOT NULL,
    name character varying(60) NOT NULL
);


--
-- Name: roles_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

ALTER TABLE public.roles ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.roles_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version bigint NOT NULL,
    dirty boolean NOT NULL
);


--
-- Name: shipments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.shipments (
    id uuid NOT NULL,
    order_id uuid NOT NULL,
    courier_code character varying(30),
    service_type character varying(30),
    tracking_no character varying(80),
    shipment_status character varying(24) DEFAULT 'waiting_pickup'::character varying NOT NULL,
    shipped_at timestamp with time zone,
    delivered_at timestamp with time zone,
    komship_order_id integer,
    komship_order_no character varying(80),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_shipments_status CHECK (((shipment_status)::text = ANY ((ARRAY['waiting_pickup'::character varying, 'shipped'::character varying, 'delivered'::character varying, 'failed'::character varying, 'returned'::character varying])::text[])))
);


--
-- Name: ticket_attachments; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ticket_attachments (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ticket_id uuid NOT NULL,
    message_id uuid,
    file_url text NOT NULL,
    file_name character varying(255) NOT NULL,
    file_type character varying(50),
    file_size integer,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: ticket_messages; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ticket_messages (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ticket_id uuid NOT NULL,
    sender_id uuid NOT NULL,
    message text NOT NULL,
    is_from_cs boolean DEFAULT false NOT NULL,
    is_internal_note boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: ticket_status_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ticket_status_logs (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ticket_id uuid NOT NULL,
    changed_by uuid NOT NULL,
    old_status character varying(20),
    new_status character varying(20) NOT NULL,
    notes text,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: ticket_subjects; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.ticket_subjects (
    id integer NOT NULL,
    label character varying(100) NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: ticket_subjects_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.ticket_subjects_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: ticket_subjects_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.ticket_subjects_id_seq OWNED BY public.ticket_subjects.id;


--
-- Name: tickets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.tickets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    ticket_number character varying(20) NOT NULL,
    customer_id uuid NOT NULL,
    assigned_cs_id uuid,
    order_number character varying(50) NOT NULL,
    phone character varying(20) NOT NULL,
    reporter_name character varying(120) NOT NULL,
    subject character varying(255) NOT NULL,
    detail text NOT NULL,
    status character varying(20) DEFAULT 'open'::character varying NOT NULL,
    source character varying(10) DEFAULT 'web'::character varying NOT NULL,
    resolved_at timestamp with time zone,
    closed_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    attachment_url text,
    attachment_content_type character varying(50),
    subject_id integer NOT NULL,
    CONSTRAINT ck_tickets_source CHECK (((source)::text = ANY ((ARRAY['app'::character varying, 'web'::character varying])::text[]))),
    CONSTRAINT ck_tickets_status CHECK (((status)::text = ANY ((ARRAY['open'::character varying, 'on_progress'::character varying, 'resolved'::character varying, 'closed'::character varying])::text[])))
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id uuid NOT NULL,
    email character varying(255) NOT NULL,
    full_name character varying(120) NOT NULL,
    birth_date date,
    phone character varying(20),
    password_hash text NOT NULL,
    role_id smallint NOT NULL,
    status character varying(16) DEFAULT 'pending'::character varying NOT NULL,
    email_verified_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_users_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'active'::character varying, 'blocked'::character varying])::text[])))
);


--
-- Name: vendor_balances; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vendor_balances (
    vendor_id uuid NOT NULL,
    available_balance numeric(18,2) DEFAULT 0 NOT NULL,
    pending_balance numeric(18,2) DEFAULT 0 NOT NULL,
    total_earned numeric(18,2) DEFAULT 0 NOT NULL,
    total_withdrawn numeric(18,2) DEFAULT 0 NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_vendor_balances_available CHECK ((available_balance >= (0)::numeric)),
    CONSTRAINT ck_vendor_balances_pending CHECK ((pending_balance >= (0)::numeric))
);


--
-- Name: vendor_bank_accounts; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vendor_bank_accounts (
    id uuid NOT NULL,
    vendor_id uuid NOT NULL,
    bank_name character varying(120) NOT NULL,
    account_number character varying(60) NOT NULL,
    account_holder_name character varying(160) NOT NULL,
    verification_status character varying(16) DEFAULT 'pending'::character varying NOT NULL,
    rejection_reason text,
    verified_by uuid,
    verified_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_vendor_bank_accounts_verification_status CHECK (((verification_status)::text = ANY ((ARRAY['pending'::character varying, 'verified'::character varying, 'rejected'::character varying])::text[]))),
    CONSTRAINT ck_vendor_bank_accounts_verified_consistency CHECK (((((verification_status)::text = 'verified'::text) AND (verified_by IS NOT NULL) AND (verified_at IS NOT NULL)) OR (((verification_status)::text = ANY ((ARRAY['pending'::character varying, 'rejected'::character varying])::text[])) AND (verified_by IS NULL) AND (verified_at IS NULL))))
);


--
-- Name: vendor_banners; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vendor_banners (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    vendor_id uuid NOT NULL,
    title character varying(255) NOT NULL,
    image_url text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: vendor_couriers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vendor_couriers (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    vendor_id uuid NOT NULL,
    courier_id integer NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: vendor_documents; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vendor_documents (
    id uuid NOT NULL,
    vendor_id uuid NOT NULL,
    doc_type character varying(40) NOT NULL,
    file_url text NOT NULL,
    mime_type character varying(100),
    file_size_bytes integer,
    file_checksum character varying(128),
    uploaded_by uuid,
    verification_status character varying(16) DEFAULT 'pending'::character varying NOT NULL,
    rejection_reason text,
    verified_by uuid,
    verified_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_vendor_documents_doc_type CHECK (((doc_type)::text = ANY ((ARRAY['owner_ktp'::character varying, 'owner_passport'::character varying, 'business_npwp'::character varying, 'store_photo'::character varying, 'bank_account_proof'::character varying, 'business_logo'::character varying, 'business_banner'::character varying])::text[]))),
    CONSTRAINT ck_vendor_documents_verification_status CHECK (((verification_status)::text = ANY ((ARRAY['pending'::character varying, 'verified'::character varying, 'rejected'::character varying])::text[]))),
    CONSTRAINT ck_vendor_documents_verified_consistency CHECK (((((verification_status)::text = 'verified'::text) AND (verified_by IS NOT NULL) AND (verified_at IS NOT NULL)) OR (((verification_status)::text = ANY ((ARRAY['pending'::character varying, 'rejected'::character varying])::text[])) AND (verified_by IS NULL) AND (verified_at IS NULL))))
);


--
-- Name: vendor_withdrawals; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vendor_withdrawals (
    id uuid NOT NULL,
    vendor_id uuid NOT NULL,
    amount numeric(18,2) NOT NULL,
    channel_code character varying(32) NOT NULL,
    status character varying(16) DEFAULT 'pending'::character varying NOT NULL,
    xendit_payout_id character varying(128),
    xendit_status character varying(32),
    description text,
    failed_reason text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    fee_estimated numeric(18,2) DEFAULT 0 NOT NULL,
    fee_actual numeric(18,2),
    amount_net numeric(18,2),
    total_deducted numeric(18,2) DEFAULT 0 NOT NULL,
    fee_status character varying(16) DEFAULT 'pending'::character varying NOT NULL,
    xendit_transaction_id character varying(128),
    CONSTRAINT ck_vendor_withdrawals_amount CHECK ((amount >= (10000)::numeric)),
    CONSTRAINT ck_vendor_withdrawals_amount_net CHECK (((amount_net IS NULL) OR (amount_net >= (0)::numeric))),
    CONSTRAINT ck_vendor_withdrawals_fee_actual CHECK (((fee_actual IS NULL) OR (fee_actual >= (0)::numeric))),
    CONSTRAINT ck_vendor_withdrawals_fee_estimated CHECK ((fee_estimated >= (0)::numeric)),
    CONSTRAINT ck_vendor_withdrawals_fee_status CHECK (((fee_status)::text = ANY ((ARRAY['pending'::character varying, 'resolved'::character varying, 'failed'::character varying])::text[]))),
    CONSTRAINT ck_vendor_withdrawals_status CHECK (((status)::text = ANY ((ARRAY['pending'::character varying, 'processing'::character varying, 'completed'::character varying, 'failed'::character varying])::text[]))),
    CONSTRAINT ck_vendor_withdrawals_total_deducted CHECK ((total_deducted >= (0)::numeric))
);


--
-- Name: vendors; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.vendors (
    id uuid NOT NULL,
    owner_user_id uuid NOT NULL,
    vendor_type character varying(32) DEFAULT 'general_souvenir_store'::character varying NOT NULL,
    legal_name character varying(160),
    display_name character varying(120) NOT NULL,
    responsible_person_name character varying(120) NOT NULL,
    description text,
    status character varying(16) DEFAULT 'draft'::character varying NOT NULL,
    approved_by uuid,
    approved_at timestamp with time zone,
    status_reason text,
    xendit_account_id character varying(255),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    CONSTRAINT ck_vendors_status CHECK (((status)::text = ANY ((ARRAY['draft'::character varying, 'submitted'::character varying, 'active'::character varying, 'rejected'::character varying, 'blocked'::character varying])::text[]))),
    CONSTRAINT ck_vendors_vendor_type CHECK (((vendor_type)::text = ANY ((ARRAY['umrah_souvenir_store'::character varying, 'hajj_souvenir_store'::character varying, 'general_souvenir_store'::character varying])::text[])))
);


--
-- Name: wishlist_items; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.wishlist_items (
    id uuid NOT NULL,
    user_id uuid NOT NULL,
    product_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


--
-- Name: admin_contacts id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_contacts ALTER COLUMN id SET DEFAULT nextval('public.admin_contacts_id_seq'::regclass);


--
-- Name: couriers id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.couriers ALTER COLUMN id SET DEFAULT nextval('public.couriers_id_seq'::regclass);


--
-- Name: faqs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.faqs ALTER COLUMN id SET DEFAULT nextval('public.faqs_id_seq'::regclass);


--
-- Name: return_reasons id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.return_reasons ALTER COLUMN id SET DEFAULT nextval('public.return_reasons_id_seq'::regclass);


--
-- Name: ticket_subjects id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_subjects ALTER COLUMN id SET DEFAULT nextval('public.ticket_subjects_id_seq'::regclass);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.addresses (id, user_id, label, recipient_name, phone, province_id, city_id, district_id, postal_code, address_line, is_default, province_name, city_name, district_name, subdistrict_id, subdistrict_name, notes, latitude, longitude, created_at, updated_at) FROM stdin;
5ca38c87-0f51-49af-8cb7-f7d1fb3fdf92	aa000002-0000-0000-0000-000000000001	ALAMAT WAREHOUSE TAMBORA	ALAMAT WAREHOUSE TAMBORA	015381035182	10	135	1328	11220	ALAMAT WAREHOUSE TAMBORA	t	DKI JAKARTA	 JAKARTA BARAT	TAMBORA	17521	TAMBORA	\N	\N	\N	2026-03-06 02:09:24.455618+00	2026-03-06 02:09:24.455618+00
46a73c45-b207-46d7-8b62-daf2a06adc52	aa000002-0000-0000-0000-000000000002	ALAMAT WAREHOUSE GEDANGAN	ALAMAT WAREHOUSE GEDANGAN	015391530162	18	583	5997	61254	ALAMAT WAREHOUSE GEDANGAN	t	JAWA TIMUR	 SIDOARJO	GEDANGAN	70917	GEDANGAN	\N	\N	\N	2026-03-06 02:10:45.278021+00	2026-03-06 02:10:45.278021+00
ff233f85-ebbd-4f81-bf74-eba609a80355	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	WIYUNG	ALAMAT WIYUNG	0278461826193	18	577	5901	60228	WIYUNG	t	JAWA TIMUR	 SURABAYA	WIYUNG	69354	WIYUNG	\N	\N	\N	2026-03-06 02:04:44.356548+00	2026-03-06 02:04:44.356548+00
9886594d-1e2f-4000-986f-0251f609edb1	8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d	SEMARANG - BANYUMANIK	TOKO SEMARANG - BANYUMANIK	027491539102	12	560	5587	50264	SEMARANG - BANYUMANIK	f	JAWA TENGAH	 SEMARANG	BANYUMANIK	64907	BANYUMANIK	SEMARANG - BANYUMANIK	\N	\N	2026-03-09 03:50:59.567225+00	2026-03-09 03:51:05.401793+00
01656d0d-286f-4999-8a69-496b4d156255	8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d	ALAMAT WAREHOUSE WIYUNG	ALAMAT WAREHOUSE WIYUNG	0836193519301	18	577	5901	60228	ALAMAT WAREHOUSE WIYUNG	t	JAWA TIMUR	 SURABAYA	WIYUNG	69354	WIYUNG	\N	\N	\N	2026-03-06 02:08:01.743737+00	2026-03-09 03:51:05.408555+00
\.


--
-- Data for Name: admin_contacts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.admin_contacts (id, content, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: banners; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.banners (id, title, image_url, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: cart_items; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.cart_items (id, cart_id, product_variant_id, qty, created_at) FROM stdin;
b4892850-c04f-41db-936a-4b227d512de2	dc0a2e9b-2409-4597-bf42-49dbb790447a	c0000001-0000-0000-0000-000000000001	2	2026-03-06 02:14:44.149978+00
5f85c553-07bd-4bd0-9254-465415406ce2	fff51f12-4fa4-4e20-9616-5d781f47fd73	c0000001-0000-0000-0000-000000000001	1	2026-03-09 03:36:14.273154+00
ac88c77f-999b-4e65-bd1c-6f62d9a7f1ee	36ca2a82-3544-42b0-be84-8557a602d4a3	c0000001-0000-0000-0000-000000000007	1	2026-03-09 06:40:00.359628+00
\.


--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.carts (id, user_id, status, created_at, updated_at) FROM stdin;
dc0a2e9b-2409-4597-bf42-49dbb790447a	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	converted	2026-03-06 02:14:44.137578+00	2026-03-06 06:39:22.183961+00
fff51f12-4fa4-4e20-9616-5d781f47fd73	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	converted	2026-03-06 06:41:49.685037+00	2026-03-09 03:37:20.13297+00
36ca2a82-3544-42b0-be84-8557a602d4a3	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	converted	2026-03-09 03:51:43.211964+00	2026-03-09 06:46:53.850863+00
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.categories (id, parent_id, name, slug, is_active) FROM stdin;
ca000001-0000-0000-0000-000000000001	\N	Perlengkapan Haji & Umrah	perlengkapan-haji-umrah	t
ca000001-0000-0000-0000-000000000002	\N	Oleh-oleh & Souvenir	oleh-oleh-souvenir	t
ca000001-0000-0000-0000-000000000003	\N	Fashion Muslim	fashion-muslim	t
ca000001-0000-0000-0000-000000000011	ca000001-0000-0000-0000-000000000001	Pakaian Ihram	pakaian-ihram	t
ca000001-0000-0000-0000-000000000012	ca000001-0000-0000-0000-000000000001	Sajadah & Alat Shalat	sajadah-alat-shalat	t
ca000001-0000-0000-0000-000000000013	ca000001-0000-0000-0000-000000000001	Tasbih & Hampers	tasbih-hampers	t
ca000001-0000-0000-0000-000000000021	ca000001-0000-0000-0000-000000000002	Kurma	kurma	t
ca000001-0000-0000-0000-000000000022	ca000001-0000-0000-0000-000000000002	Minyak & Herbal	minyak-herbal	t
ca000001-0000-0000-0000-000000000023	ca000001-0000-0000-0000-000000000002	Air Zamzam	air-zamzam	t
ca000001-0000-0000-0000-000000000031	ca000001-0000-0000-0000-000000000003	Mukena	mukena	t
ca000001-0000-0000-0000-000000000032	ca000001-0000-0000-0000-000000000003	Gamis & Jubah	gamis-jubah	t
\.


--
-- Data for Name: chat_conversations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.chat_conversations (id, initiator_id, participant_id, last_message_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: chat_messages; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.chat_messages (id, conversation_id, sender_id, message, attachment_url, is_read, read_at, created_at) FROM stdin;
\.


--
-- Data for Name: couriers; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.couriers (id, code, name, logo_url, is_active, created_at) FROM stdin;
1	jne	JNE	\N	t	2026-03-06 01:58:19.265177+00
2	sicepat	SiCepat	\N	t	2026-03-06 01:58:19.265177+00
3	ide	IDExpress	\N	t	2026-03-06 01:58:19.265177+00
4	sap	SAP Express	\N	t	2026-03-06 01:58:19.265177+00
5	ninja	Ninja	\N	t	2026-03-06 01:58:19.265177+00
6	jnt	J&T Express	\N	t	2026-03-06 01:58:19.265177+00
7	tiki	TIKI	\N	t	2026-03-06 01:58:19.265177+00
8	wahana	Wahana Express	\N	t	2026-03-06 01:58:19.265177+00
9	pos	POS Indonesia	\N	t	2026-03-06 01:58:19.265177+00
10	sentral	Sentral Cargo	\N	t	2026-03-06 01:58:19.265177+00
11	lion	Lion Parcel	\N	t	2026-03-06 01:58:19.265177+00
12	rex	Royal Express Asia	\N	t	2026-03-06 01:58:19.265177+00
\.


--
-- Data for Name: email_verification_tokens; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.email_verification_tokens (id, user_id, email, token_hash, expires_at, consumed_at, invalidated_at, created_at) FROM stdin;
\.


--
-- Data for Name: faqs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.faqs (id, category, question, answer, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: ledger_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ledger_accounts (id, code, name, account_type, normal_side, is_active) FROM stdin;
d2b4a5eb-1555-4129-bfb5-03f0f440e081	1100	Payment Gateway Receivable	asset	D	t
980eb053-9b2c-403c-a952-7cd5d74d5e47	2100	Vendor Payable	liability	C	t
0a7c297a-a2c2-49a3-8ed7-10cc9332bd0b	4100	Platform Fee Revenue	revenue	C	t
0277561c-c81d-4c55-bab4-4f895e4a9b29	4200	Admin Fee Revenue	revenue	C	t
\.


--
-- Data for Name: ledger_journals; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ledger_journals (id, journal_no, source_type, source_id, event_time, description, status, created_by, created_at) FROM stdin;
\.


--
-- Data for Name: ledger_lines; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ledger_lines (id, journal_id, account_id, debit, credit, currency, reference) FROM stdin;
\.


--
-- Data for Name: notifications; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.notifications (id, user_id, type, title, message, created_at, read_at) FROM stdin;
\.


--
-- Data for Name: order_items; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.order_items (id, order_id, product_variant_id, product_name_snapshot, sku_snapshot, qty, unit_price, line_total) FROM stdin;
4feeae02-05a3-493a-b0a1-a66e3accfc6f	2413f18e-1462-4a6c-8696-d1d3f23494ea	c0000001-0000-0000-0000-000000000001	Kain Ihram Pria Premium	IHR-PRM-STD	2	75000.00	150000.00
865bb7a0-5a4f-4a7c-aa4e-5b01f9be15e7	d75fcc53-4e31-4245-a7d4-865e41ac27b1	c0000001-0000-0000-0000-000000000001	Kain Ihram Pria Premium	IHR-PRM-STD	2	75000.00	150000.00
9ae9a60e-cf33-4b50-9a29-51c4601f83c5	4354c494-e4dc-4b1b-9f2e-94dfd8e42ff8	c0000001-0000-0000-0000-000000000001	Kain Ihram Pria Premium	IHR-PRM-STD	2	75000.00	150000.00
eed6357c-d902-4bd5-9ae5-c5f14fc682ef	3df49783-a642-4a92-a95a-86de5ae501e4	c0000001-0000-0000-0000-000000000001	Kain Ihram Pria Premium	IHR-PRM-STD	1	75000.00	75000.00
791f089f-c016-44e0-8caf-dd82737fe46e	d8933c5f-1698-40ef-b1a8-40cfe5f45d3d	c0000001-0000-0000-0000-000000000007	Kurma Ajwa Madinah Asli	KRM-AJW-500	1	120000.00	120000.00
27841cdf-ae68-447d-afc5-ae1dd07d9f95	16698988-c415-4cf3-a9d2-664ec99407a3	c0000001-0000-0000-0000-000000000007	Kurma Ajwa Madinah Asli	KRM-AJW-500	1	120000.00	120000.00
\.


--
-- Data for Name: order_status_history; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.order_status_history (id, order_id, old_status, new_status, changed_by, changed_at, notes) FROM stdin;
dcde6f7d-a8f7-43b9-b1e8-289b62c9a28f	2413f18e-1462-4a6c-8696-d1d3f23494ea	\N	pending_payment	\N	2026-03-06 06:07:20.132364+00	\N
cf3e5a28-d550-46a0-9104-c2ce8e2316ce	2413f18e-1462-4a6c-8696-d1d3f23494ea	pending_payment	canceled	\N	2026-03-06 06:07:20.875562+00	Checkout rolled back: failed to create payment invoice: xendit API error (status 401, code INVALID_API_KEY): The API key provided is invalid. Please make sure to use the secret/public API key that you can obtain from the Xendit Dashboard. See https://developers.xendit.co/api-reference/#authentication for more details
a558bd1a-df71-412e-a39c-cd65e65b5fbc	d75fcc53-4e31-4245-a7d4-865e41ac27b1	\N	pending_payment	\N	2026-03-06 06:12:22.279086+00	\N
b0a2cd0d-8d5c-4d57-8a5d-2f0de7483f71	d75fcc53-4e31-4245-a7d4-865e41ac27b1	pending_payment	canceled	\N	2026-03-06 06:12:22.888173+00	Checkout rolled back: failed to create payment invoice: xendit API error (status 404, code XEN_PLATFORM_RELATIONSHIP_NOT_FOUND_ERROR): Cannot find sub-account data or business is not master account
54c0e353-149a-42c2-b53e-e3ee48d6dd6b	4354c494-e4dc-4b1b-9f2e-94dfd8e42ff8	\N	pending_payment	\N	2026-03-06 06:39:22.179691+00	\N
9f93fbc2-6359-4d79-b8ed-542d47d9f45d	3df49783-a642-4a92-a95a-86de5ae501e4	\N	pending_payment	\N	2026-03-09 03:37:20.129622+00	\N
3fb48346-03a3-41b6-b2ec-bf0c73f97a5e	d8933c5f-1698-40ef-b1a8-40cfe5f45d3d	\N	pending_payment	\N	2026-03-09 06:40:31.101452+00	\N
58e7d979-b428-485d-b450-bf4f1e52d30f	d8933c5f-1698-40ef-b1a8-40cfe5f45d3d	pending_payment	canceled	\N	2026-03-09 06:40:31.720197+00	Checkout rolled back: failed to create payment invoice: xendit API error (status 400, code API_VALIDATION_ERROR): Header 'for-user-id' is not in a valid XenPlatform sub-account ID format
4c3b3571-3da1-4b76-afab-1a9d8d4df96d	16698988-c415-4cf3-a9d2-664ec99407a3	\N	pending_payment	\N	2026-03-09 06:46:53.829696+00	\N
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.orders (id, order_no, user_id, vendor_id, shipping_address_snapshot, order_status, payment_status, subtotal, shipping_fee, platform_fee, grand_total, placed_at, created_at, updated_at) FROM stdin;
2413f18e-1462-4a6c-8696-d1d3f23494ea	ORD-20260306-FC389295	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	da000001-0000-0000-0000-000000000001	{"label": "WIYUNG", "phone": "0278461826193", "city_name": " SURABAYA", "address_id": "ff233f85-ebbd-4f81-bf74-eba609a80355", "is_default": true, "address_line": "WIYUNG", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "ALAMAT WIYUNG", "subdistrict_name": "WIYUNG"}	canceled	unpaid	150000.00	6000.00	6000.00	162000.00	2026-03-06 06:07:20.127604+00	2026-03-06 06:07:20.127604+00	2026-03-06 06:07:20.873926+00
d75fcc53-4e31-4245-a7d4-865e41ac27b1	ORD-20260306-C2D078ED	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	da000001-0000-0000-0000-000000000001	{"label": "WIYUNG", "phone": "0278461826193", "city_name": " SURABAYA", "address_id": "ff233f85-ebbd-4f81-bf74-eba609a80355", "is_default": true, "address_line": "WIYUNG", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "ALAMAT WIYUNG", "subdistrict_name": "WIYUNG"}	canceled	unpaid	150000.00	6000.00	6000.00	162000.00	2026-03-06 06:12:22.274566+00	2026-03-06 06:12:22.274566+00	2026-03-06 06:12:22.886718+00
4354c494-e4dc-4b1b-9f2e-94dfd8e42ff8	ORD-20260306-1A12C83F	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	da000001-0000-0000-0000-000000000001	{"label": "WIYUNG", "phone": "0278461826193", "city_name": " SURABAYA", "address_id": "ff233f85-ebbd-4f81-bf74-eba609a80355", "is_default": true, "address_line": "WIYUNG", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "ALAMAT WIYUNG", "subdistrict_name": "WIYUNG"}	paid	paid	150000.00	6000.00	6000.00	162000.00	2026-03-06 06:39:22.176503+00	2026-03-06 06:39:22.176503+00	2026-03-06 06:39:22.176503+00
3df49783-a642-4a92-a95a-86de5ae501e4	ORD-20260309-CA7FCEB2	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	da000001-0000-0000-0000-000000000001	{"label": "WIYUNG", "phone": "0278461826193", "city_name": " SURABAYA", "address_id": "ff233f85-ebbd-4f81-bf74-eba609a80355", "is_default": true, "postal_code": "60228", "address_line": "WIYUNG", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "ALAMAT WIYUNG", "subdistrict_name": "WIYUNG"}	pending_payment	unpaid	75000.00	6000.00	6000.00	87000.00	2026-03-09 03:37:20.114567+00	2026-03-09 03:37:20.114567+00	2026-03-09 03:37:20.114567+00
d8933c5f-1698-40ef-b1a8-40cfe5f45d3d	ORD-20260309-E71EA86E	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	da000001-0000-0000-0000-000000000002	{"label": "WIYUNG", "phone": "0278461826193", "city_name": " SURABAYA", "address_id": "ff233f85-ebbd-4f81-bf74-eba609a80355", "is_default": true, "postal_code": "60228", "address_line": "WIYUNG", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "ALAMAT WIYUNG", "subdistrict_name": "WIYUNG"}	canceled	unpaid	120000.00	18000.00	6000.00	144000.00	2026-03-09 06:40:30.992795+00	2026-03-09 06:40:30.992795+00	2026-03-09 06:40:31.718669+00
16698988-c415-4cf3-a9d2-664ec99407a3	ORD-20260309-8F2B8EC7	7d281d37-8f4a-48a8-9b33-e47c76b3d13c	da000001-0000-0000-0000-000000000002	{"label": "WIYUNG", "phone": "0278461826193", "city_name": " SURABAYA", "address_id": "ff233f85-ebbd-4f81-bf74-eba609a80355", "is_default": true, "postal_code": "60228", "address_line": "WIYUNG", "district_name": "WIYUNG", "province_name": "JAWA TIMUR", "recipient_name": "ALAMAT WIYUNG", "subdistrict_name": "WIYUNG"}	pending_payment	unpaid	120000.00	18000.00	6000.00	144000.00	2026-03-09 06:46:53.739214+00	2026-03-09 06:46:53.739214+00	2026-03-09 06:46:53.739214+00
\.


--
-- Data for Name: payment_events; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.payment_events (id, payment_invoice_id, event_type, external_event_id, payload, received_at) FROM stdin;
\.


--
-- Data for Name: payment_invoices; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.payment_invoices (id, order_id, gateway, xendit_invoice_id, external_invoice_id, invoice_url, payment_method, payment_channel, amount, currency, status, expires_at, paid_at, raw_payload, created_at, updated_at) FROM stdin;
791163be-195e-497b-8181-be246f79ad64	16698988-c415-4cf3-a9d2-664ec99407a3	xendit	noop-inv-9dcb2cc9	INV-ORD-20260309-8F2B8EC7-e9597e44	https://checkout-bypass.example.com/noop-inv-9dcb2cc9	\N	\N	144000.00	IDR	pending	2026-03-10 06:46:53+00	\N	\N	2026-03-09 06:46:53.739214+00	2026-03-09 06:46:53.739214+00
\.


--
-- Data for Name: payout_batches; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.payout_batches (id, vendor_id, period_start, period_end, status, total_gross, total_fee, total_net, paid_at, created_by, xendit_payout_id, channel_code, description, xendit_status, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: payout_items; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.payout_items (id, payout_batch_id, order_id, gross_amount, platform_fee_amount, net_amount) FROM stdin;
\.


--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.product_images (id, product_id, image_url, mime_type, file_size_bytes, is_primary, sort_order, created_at) FROM stdin;
a1b00001-0000-0000-0000-000000000001	b0000001-0000-0000-0000-000000000001	products/kain-ihram-pria-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000002	b0000001-0000-0000-0000-000000000001	products/kain-ihram-pria-2.jpg	image/jpeg	\N	f	1	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000003	b0000001-0000-0000-0000-000000000002	products/sajadah-travel-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000004	b0000001-0000-0000-0000-000000000003	products/tasbih-digital-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000005	b0000001-0000-0000-0000-000000000003	products/tasbih-digital-2.jpg	image/jpeg	\N	f	1	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000006	b0000001-0000-0000-0000-000000000004	products/kurma-ajwa-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000007	b0000001-0000-0000-0000-000000000004	products/kurma-ajwa-2.jpg	image/jpeg	\N	f	1	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000008	b0000001-0000-0000-0000-000000000005	products/minyak-zaitun-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000009	b0000001-0000-0000-0000-000000000006	products/air-zamzam-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000010	b0000001-0000-0000-0000-000000000007	products/mukena-katun-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000011	b0000001-0000-0000-0000-000000000007	products/mukena-katun-2.jpg	image/jpeg	\N	f	1	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000012	b0000001-0000-0000-0000-000000000008	products/gamis-pria-1.jpg	image/jpeg	\N	t	0	2026-03-06 01:58:19.265177+00
a1b00001-0000-0000-0000-000000000013	b0000001-0000-0000-0000-000000000008	products/gamis-pria-2.jpg	image/jpeg	\N	f	1	2026-03-06 01:58:19.265177+00
\.


--
-- Data for Name: product_review_images; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.product_review_images (id, review_id, object_key, mime_type, file_size_bytes, sort_order, created_at) FROM stdin;
\.


--
-- Data for Name: product_review_stats; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.product_review_stats (product_id, total_reviews, total_stars, average_rating, star_0_count, star_1_count, star_2_count, star_3_count, star_4_count, star_5_count, updated_at) FROM stdin;
\.


--
-- Data for Name: product_reviews; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.product_reviews (id, product_id, order_id, order_item_id, user_id, rating, review_text, status, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: product_variants; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.product_variants (id, product_id, sku, variant_name, price, currency, stock_on_hand, weight_gram, is_default, is_active) FROM stdin;
c0000001-0000-0000-0000-000000000002	b0000001-0000-0000-0000-000000000001	IHR-PRM-XL	Extra Large (130x240 cm)	95000.00	IDR	80	500	f	t
c0000001-0000-0000-0000-000000000003	b0000001-0000-0000-0000-000000000002	SJD-TRV-GRN	Hijau Tosca	45000.00	IDR	200	300	t	t
c0000001-0000-0000-0000-000000000004	b0000001-0000-0000-0000-000000000003	TSB-DGT-BLK	Hitam	25000.00	IDR	300	50	t	t
c0000001-0000-0000-0000-000000000005	b0000001-0000-0000-0000-000000000003	TSB-DGT-GLD	Gold	35000.00	IDR	200	55	f	t
c0000001-0000-0000-0000-000000000006	b0000001-0000-0000-0000-000000000004	KRM-AJW-250	250 gram	65000.00	IDR	120	300	f	t
c0000001-0000-0000-0000-000000000008	b0000001-0000-0000-0000-000000000004	KRM-AJW-1KG	1 Kilogram	220000.00	IDR	40	1050	f	t
c0000001-0000-0000-0000-000000000009	b0000001-0000-0000-0000-000000000005	MZT-EV-250	250 ml	85000.00	IDR	100	350	t	t
c0000001-0000-0000-0000-000000000010	b0000001-0000-0000-0000-000000000005	MZT-EV-500	500 ml	155000.00	IDR	60	650	f	t
c0000001-0000-0000-0000-000000000011	b0000001-0000-0000-0000-000000000006	ZMZ-5L	5 Liter	175000.00	IDR	30	5200	t	t
c0000001-0000-0000-0000-000000000012	b0000001-0000-0000-0000-000000000007	MKN-KTJ-WHT	Putih Polos	185000.00	IDR	60	450	t	t
c0000001-0000-0000-0000-000000000013	b0000001-0000-0000-0000-000000000007	MKN-KTJ-PNK	Dusty Pink	195000.00	IDR	45	460	f	t
c0000001-0000-0000-0000-000000000014	b0000001-0000-0000-0000-000000000008	GMS-ALH-M	M (Lingkar Dada 104cm)	250000.00	IDR	40	380	f	t
c0000001-0000-0000-0000-000000000015	b0000001-0000-0000-0000-000000000008	GMS-ALH-L	L (Lingkar Dada 110cm)	250000.00	IDR	55	400	t	t
c0000001-0000-0000-0000-000000000016	b0000001-0000-0000-0000-000000000008	GMS-ALH-XL	XL (Lingkar Dada 118cm)	265000.00	IDR	35	420	f	t
c0000001-0000-0000-0000-000000000001	b0000001-0000-0000-0000-000000000001	IHR-PRM-STD	Standard (115x220 cm)	75000.00	IDR	147	400	t	t
c0000001-0000-0000-0000-000000000007	b0000001-0000-0000-0000-000000000004	KRM-AJW-500	500 gram	120000.00	IDR	79	550	t	t
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.products (id, vendor_id, category_id, name, slug, description, status, halal_ai_status, halal_ai_notes, created_at, updated_at) FROM stdin;
b0000001-0000-0000-0000-000000000001	da000001-0000-0000-0000-000000000001	ca000001-0000-0000-0000-000000000011	Kain Ihram Pria Premium	kain-ihram-pria-premium	Kain ihram pria berbahan katun premium grade A, lembut dan nyaman untuk ibadah haji dan umrah. Jahitan rapi, tidak mudah kusut. Tersedia dalam ukuran Standard dan Extra Large.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
b0000001-0000-0000-0000-000000000002	da000001-0000-0000-0000-000000000001	ca000001-0000-0000-0000-000000000012	Sajadah Travel Lipat Premium	sajadah-travel-lipat-premium	Sajadah travel lipat praktis dengan tas penyimpanan waterproof. Bahan beludru soft-touch, ringan hanya 300 gram. Ideal untuk dibawa haji dan umrah.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
b0000001-0000-0000-0000-000000000003	da000001-0000-0000-0000-000000000001	ca000001-0000-0000-0000-000000000013	Tasbih Digital Premium 33 Biji	tasbih-digital-premium-33	Tasbih digital elektronik dengan counter otomatis. Baterai tahan hingga 6 bulan pemakaian normal. Material ABS anti-slip, nyaman digenggam.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
b0000001-0000-0000-0000-000000000004	da000001-0000-0000-0000-000000000002	ca000001-0000-0000-0000-000000000021	Kurma Ajwa Madinah Asli	kurma-ajwa-madinah-asli	Kurma Ajwa asli dari kebun Al-Madinah Al-Munawwarah. Tekstur lembut, rasa manis alami. Sering disebut sebagai kurma Nabi. Kemasan vacuum seal untuk menjaga kesegaran.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
b0000001-0000-0000-0000-000000000005	da000001-0000-0000-0000-000000000002	ca000001-0000-0000-0000-000000000022	Minyak Zaitun Extra Virgin Asli	minyak-zaitun-extra-virgin-asli	Minyak zaitun extra virgin 100% asli impor dari Timur Tengah. Cold-pressed, tanpa campuran. Cocok untuk kesehatan dan kecantikan.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
b0000001-0000-0000-0000-000000000006	da000001-0000-0000-0000-000000000002	ca000001-0000-0000-0000-000000000023	Air Zamzam Asli 5 Liter	air-zamzam-asli-5-liter	Air Zamzam asli dari sumur Zamzam, Mekkah. Dikemas dalam jeriken food-grade kedap udara. Sertifikat keaslian tersedia.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
b0000001-0000-0000-0000-000000000007	da000001-0000-0000-0000-000000000003	ca000001-0000-0000-0000-000000000031	Mukena Katun Jepang Premium	mukena-katun-jepang-premium	Mukena bahan katun Jepang super halus dan adem. Tidak menerawang, jatuh sempurna. Dilengkapi tas travel matching. Tersedia warna Putih dan Dusty Pink.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
b0000001-0000-0000-0000-000000000008	da000001-0000-0000-0000-000000000003	ca000001-0000-0000-0000-000000000032	Gamis Pria Al-Haramain	gamis-pria-al-haramain	Gamis pria model Al-Haramain lengan panjang. Bahan toyobo premium anti kusut, adem, dan tidak menerawang. Cocok untuk shalat dan acara formal.	published	passed	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
\.


--
-- Data for Name: refunds; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.refunds (id, order_id, payment_invoice_id, amount, reason, status, requested_by, processed_by, processed_at) FROM stdin;
\.


--
-- Data for Name: reply_templates; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.reply_templates (id, title, content, is_active, created_by, updated_by, created_at, updated_at, shortcut, category) FROM stdin;
b764b2a3-eda7-4fa5-8973-5af4ed80f15a	Selamat Datang di UmrahMart	Assalamu'alaikum, selamat datang di UmrahMart! 🕌\n\nKami hadir untuk memudahkan Anda mendapatkan perlengkapan haji dan umroh terbaik. Ada yang bisa kami bantu hari ini?	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/general/selamat-datang	general
ab5db7ec-d629-4c6d-bd0b-6898744975c9	Terima Kasih Telah Menghubungi Kami	Terima kasih telah menghubungi tim Customer Service UmrahMart. 🙏\n\nKami akan segera membantu Anda. Mohon tunggu sebentar ya.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/general/terima-kasih-menghubungi	general
89c571de-633b-4839-826b-8c07bc9a01eb	Apakah Masih Ada Yang Bisa Dibantu?	Apakah masih ada yang bisa kami bantu? 😊\n\nJika masalah Anda sudah teratasi, jangan lupa berikan ulasan untuk membantu kami meningkatkan layanan. Semoga ibadah haji/umroh Anda berjalan lancar. Aamiin.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/general/ada-yang-bisa-dibantu	general
a03e89a3-072e-46bf-938a-e980c5b8564a	Jam Operasional Customer Service	Tim Customer Service UmrahMart melayani Anda setiap hari:\n\n🕐 Senin – Jumat : 08.00 – 21.00 WIB\n🕐 Sabtu – Minggu : 09.00 – 18.00 WIB\n\nDi luar jam operasional, Anda tetap bisa meninggalkan pesan dan kami akan merespons saat jam kerja.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/general/jam-operasional	general
5a366885-ae98-4549-ba12-425ebe40bcb5	Konfirmasi Pesanan Diterima	Alhamdulillah, pesanan Anda telah kami terima! ✅\n\nNomor pesanan Anda: {order_number}\nStatus: Menunggu Pembayaran\n\nSilakan selesaikan pembayaran sebelum {expired_at} agar pesanan dapat segera kami proses. Jazakallahu khairan.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/order/konfirmasi-pesanan-diterima	order
99ba26a0-e0ed-4c49-a5ee-6ffbdc22450f	Pesanan Sedang Diproses	Kabar baik! Pesanan Anda dengan nomor {order_number} sedang kami proses. 📦\n\nEstimasi pesanan siap dikirim: 1–2 hari kerja.\nKami akan menginformasikan nomor resi pengiriman segera setelah paket diserahkan ke kurir.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/order/pesanan-diproses	order
67a8d847-500b-407b-a1b5-d7f3a517a781	Pesanan Berhasil Dibatalkan	Pesanan Anda dengan nomor {order_number} telah berhasil dibatalkan.\n\nJika Anda sudah melakukan pembayaran, proses refund akan kami lakukan dalam 3–7 hari kerja ke metode pembayaran asal.\n\nApabila ada pertanyaan lebih lanjut, jangan ragu untuk menghubungi kami kembali.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/order/pesanan-dibatalkan	order
fa2a2b13-6476-4c37-ac2e-17658935d8cd	Detail Pesanan	Berikut detail pesanan Anda:\n\n📋 Nomor Pesanan : {order_number}\n📅 Tanggal       : {order_date}\n💰 Total         : {total_amount}\n📍 Status        : {order_status}\n\nUntuk informasi lebih lengkap, silakan cek halaman "Pesanan Saya" di aplikasi.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/order/detail-pesanan	order
046c2619-fbf4-4084-986e-ee1234972eaf	Pembayaran Berhasil Diterima	Alhamdulillah, pembayaran Anda telah berhasil kami terima! 💚\n\nNomor pesanan : {order_number}\nJumlah        : {amount}\nMetode        : {payment_method}\n\nPesanan Anda akan segera kami proses. Terima kasih.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/payment/pembayaran-berhasil	payment
123b2d38-ff07-44da-a3a4-c0b064e801c8	Menunggu Konfirmasi Pembayaran	Halo, kami melihat pesanan Anda belum terkonfirmasi pembayarannya.\n\nBatas waktu pembayaran: {expired_at}\n\nJika Anda sudah melakukan transfer manual, mohon upload bukti pembayaran melalui aplikasi agar kami dapat segera memprosesnya.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/payment/menunggu-konfirmasi	payment
30295239-f60b-4108-a926-da75bea42823	Pembayaran Gagal / Kedaluwarsa	Kami informasikan bahwa pembayaran untuk pesanan {order_number} telah gagal atau melewati batas waktu. ⚠️\n\nPesanan Anda telah kami batalkan secara otomatis. Anda dapat melakukan pemesanan ulang kapan saja.\n\nMohon pastikan saldo/limit mencukupi saat melakukan pembayaran berikutnya.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/payment/pembayaran-gagal	payment
204c4201-f07c-49c8-b5cc-29390accb75a	Cara Melakukan Pembayaran	Berikut metode pembayaran yang tersedia di UmrahMart:\n\n💳 Transfer Bank (BCA, Mandiri, BNI, BRI)\n📱 E-Wallet (GoPay, OVO, DANA, ShopeePay)\n🏪 Gerai Retail (Alfamart, Indomaret)\n💵 QRIS\n\nPilih metode yang paling nyaman untuk Anda saat checkout.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/payment/cara-pembayaran	payment
1e4ecad1-7315-408f-8fee-a24b0f8a27c1	Pengajuan Refund Diterima	Pengajuan refund Anda telah kami terima dan sedang dalam proses review. ✅\n\nNomor tiket refund: {ticket_number}\n\nTim kami akan memverifikasi dalam 1–3 hari kerja. Anda akan mendapat notifikasi setelah proses selesai.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/refund/pengajuan-diterima	refund
e900736a-e42d-4ba3-8097-32d74c5434a0	Refund Sedang Diproses	Refund Anda sedang dalam proses pencairan. 🔄\n\nJumlah refund : {refund_amount}\nTujuan        : {refund_destination}\nEstimasi      : 3–7 hari kerja\n\nMohon bersabar dan pastikan rekening/e-wallet tujuan masih aktif.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/refund/sedang-diproses	refund
694a07b1-8b30-496e-9cb1-48cebc1bce5c	Refund Berhasil	Dana refund sebesar {refund_amount} telah berhasil kami kirimkan ke {refund_destination}. 🎉\n\nJika dalam 1×24 jam dana belum masuk, mohon hubungi kami kembali dengan menyertakan nomor tiket {ticket_number}.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/refund/berhasil	refund
8673e7ca-1f46-469f-ab64-8fdbafcdfff3	Pesanan Telah Dikirim	Pesanan Anda sudah dalam perjalanan! 🚚\n\nNomor Pesanan : {order_number}\nKurir         : {courier_name}\nNomor Resi    : {tracking_number}\n\nLacak pengiriman di: {tracking_url}\n\nEstimasi tiba: {estimated_arrival}	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/shipping/pesanan-dikirim	shipping
9539af5b-e0ce-4da3-be3d-f5536f6c75df	Konfirmasi Penerimaan Paket	Halo, berdasarkan data kurir, paket Anda sudah dinyatakan terkirim pada {delivered_at}.\n\nApakah paket sudah Anda terima dengan kondisi baik? Mohon konfirmasi agar pesanan dapat diselesaikan. 📦✅	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/shipping/konfirmasi-penerimaan	shipping
457ad0fa-4140-45b5-a017-99b6b5e39c6e	Paket Tertahan / Terlambat	Kami melihat terdapat kendala pada pengiriman paket Anda (no. resi: {tracking_number}). ⚠️\n\nTim kami sedang berkoordinasi dengan pihak kurir untuk menindaklanjuti hal ini. Kami akan segera menginformasikan perkembangannya. Mohon maaf atas ketidaknyamanannya.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/shipping/paket-tertahan	shipping
19a767b9-edf3-4456-9afe-1e4e546c4cc8	Informasi Ketersediaan Stok	Terima kasih atas minat Anda pada produk {product_name}.\n\nSaat ini stok produk tersebut {stock_status}.\n\nAnda dapat mengaktifkan notifikasi "Ingatkan Saya" pada halaman produk agar kami bisa memberitahu Anda saat stok tersedia kembali.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/product/ketersediaan-stok	product
ae88c368-57a6-4943-8f5b-c7c4e812ee2a	Detail Spesifikasi Produk	Berikut informasi lengkap mengenai {product_name}:\n\n📌 Bahan    : {material}\n📐 Ukuran   : {size}\n🎨 Warna    : {color}\n🏷️ Berat    : {weight}\n✅ Halal    : {halal_status}\n\nApakah ada pertanyaan lain mengenai produk ini?	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/product/spesifikasi-produk	product
4a80507a-3346-42cb-a3e1-06aec1c94efe	Produk Tidak Lagi Tersedia	Mohon maaf, produk {product_name} saat ini sudah tidak tersedia di katalog kami. 😔\n\nKami merekomendasikan produk serupa yang mungkin sesuai dengan kebutuhan Anda:\n👉 {alternative_product}\n\nSilakan cek koleksi terbaru kami di aplikasi.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/product/produk-tidak-tersedia	product
e90a750f-b9b3-4eaa-948b-4e329d3ecada	Bantuan Reset Password	Untuk mereset password akun Anda, ikuti langkah berikut:\n\n1️⃣ Buka halaman Login\n2️⃣ Klik "Lupa Password?"\n3️⃣ Masukkan email terdaftar Anda\n4️⃣ Cek inbox email untuk link reset (cek juga folder Spam)\n5️⃣ Klik link dan buat password baru\n\nLink reset berlaku selama 60 menit. Jika belum menerima email, hubungi kami kembali.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/account/reset-password	account
5827ded9-9083-4c3c-ac4d-80b3ad2ea022	Verifikasi Email	Untuk memverifikasi email Anda:\n\n1️⃣ Cek inbox email {email} (termasuk folder Spam/Junk)\n2️⃣ Buka email dari UmrahMart dengan subjek "Verifikasi Akun"\n3️⃣ Klik tombol "Verifikasi Sekarang"\n\nJika email tidak ditemukan, klik "Kirim Ulang Verifikasi" di halaman akun Anda.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/account/verifikasi-email	account
156e1f79-fee5-47c4-992a-2857748bfd79	Akun Diblokir / Dinonaktifkan	Kami melihat akun Anda saat ini tidak aktif. 🔒\n\nHal ini dapat terjadi karena:\n• Aktivitas yang mencurigakan terdeteksi\n• Pelanggaran syarat & ketentuan penggunaan\n• Permintaan penonaktifan sebelumnya\n\nUntuk mengajukan reaktivasi, mohon kirimkan data diri Anda ke email support@umrahmart.id.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/account/akun-diblokir	account
e68b0c9b-d6cd-462d-935f-a065e819f043	Keluhan Diterima dan Dicatat	Kami sangat menyesal mendengar pengalaman yang tidak menyenangkan ini. 🙏\n\nKeluhan Anda telah kami catat dengan nomor tiket {ticket_number}. Tim kami akan menginvestigasi dan menghubungi Anda dalam 1×24 jam kerja.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/complaint/keluhan-diterima	complaint
c81fd25f-dd50-46a7-ada5-1b40e8a342ee	Keluhan Sedang Diinvestigasi	Kami tengah menginvestigasi keluhan yang Anda sampaikan terkait {complaint_subject}.\n\nProses investigasi membutuhkan waktu maksimal 3 hari kerja. Kami mohon kesabaran Anda dan akan segera memberikan informasi lebih lanjut.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/complaint/sedang-investigasi	complaint
339ad257-6a07-4b6f-a28f-4a88570d9583	Keluhan Telah Diselesaikan	Keluhan Anda dengan nomor tiket {ticket_number} telah kami selesaikan. ✅\n\nSolusi yang diberikan: {resolution}\n\nKami memohon maaf atas pengalaman yang kurang menyenangkan ini dan berterima kasih atas masukan Anda untuk perbaikan layanan kami.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/complaint/telah-diselesaikan	complaint
758cfbe7-634b-4865-b579-1ad6c7ce20bb	Prosedur Pengembalian Barang	Berikut prosedur pengembalian barang (return) di UmrahMart:\n\n1️⃣ Ajukan return melalui menu "Pesanan Saya" → "Ajukan Return"\n2️⃣ Pilih produk dan alasan pengembalian\n3️⃣ Unggah foto kondisi barang\n4️⃣ Tunggu konfirmasi persetujuan (1–2 hari kerja)\n5️⃣ Kirim barang ke alamat yang tertera\n\nSyarat: barang dalam kondisi asli, belum dipakai, dan beserta kemasan lengkap. Maksimal 7 hari setelah diterima.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/return/prosedur-pengembalian	return
fee0744f-3470-4824-8a8a-d39854285d13	Pengajuan Return Disetujui	Kabar baik! Pengajuan return Anda untuk pesanan {order_number} telah disetujui. ✅\n\nSilakan kirimkan barang ke:\n📍 {return_address}\nAtas nama: UmrahMart Returns\n\nSetelah barang kami terima dan verifikasi, dana akan dikembalikan dalam 3–5 hari kerja.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/return/disetujui	return
efb48e03-7862-402a-b432-b4b890fc77ab	Pengajuan Return Ditolak	Mohon maaf, pengajuan return untuk pesanan {order_number} tidak dapat kami proses. ❌\n\nAlasan: {rejection_reason}\n\nJika Anda merasa keberatan, Anda dapat mengajukan banding dengan menghubungi kami dan melampirkan bukti pendukung.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/return/ditolak	return
2209a107-5f45-42d3-8b24-c38b9b609e94	Informasi Promo Aktif	Halo! Berikut promo yang sedang berlangsung di UmrahMart 🎉:\n\n🏷️ {promo_name}\n💰 Diskon hingga {discount_value}\n📅 Berlaku: {promo_start} – {promo_end}\n📝 Syarat: {promo_terms}\n\nJangan sampai terlewat! Belanja sekarang sebelum promo berakhir.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/promo/info-promo-aktif	promo
7470ed4e-cc43-4afb-8698-8d08336e0206	Promo Tidak Berlaku untuk Pesanan Ini	Mohon maaf, promo yang Anda gunakan tidak berlaku untuk pesanan ini. ⚠️\n\nKemungkinan penyebab:\n• Promo sudah berakhir\n• Produk tidak termasuk dalam kategori promo\n• Minimum pembelian tidak terpenuhi\n\nCek syarat & ketentuan promo di halaman "Promo" pada aplikasi.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/promo/tidak-berlaku	promo
85494871-711d-4d6d-bd81-5d1ef0dadb21	Cara Menggunakan Voucher	Berikut cara menggunakan voucher di UmrahMart:\n\n1️⃣ Tambahkan produk ke keranjang\n2️⃣ Masuk ke halaman Checkout\n3️⃣ Klik "Gunakan Voucher"\n4️⃣ Masukkan kode voucher: {voucher_code}\n5️⃣ Klik "Terapkan" — diskon akan langsung terpotong\n\nPastikan pesanan Anda memenuhi minimum pembelian voucher.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/voucher/cara-menggunakan	voucher
9624e831-9eac-4cd6-bfc1-07a5065fc4f5	Voucher Tidak Valid atau Kedaluwarsa	Kode voucher yang Anda masukkan tidak dapat digunakan. ❌\n\nKemungkinan penyebab:\n• Kode voucher salah atau sudah digunakan\n• Voucher sudah kedaluwarsa ({expired_date})\n• Akun Anda tidak memenuhi syarat voucher\n\nUntuk bantuan lebih lanjut, silakan kirimkan screenshot kode voucher Anda.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/voucher/tidak-valid	voucher
405df36d-ae2f-431f-bb3d-84b64222de15	Masalah Teknis pada Aplikasi	Kami mohon maaf atas kendala teknis yang Anda alami. 🛠️\n\nBeberapa langkah yang dapat dicoba:\n1. Tutup dan buka kembali aplikasi\n2. Pastikan koneksi internet stabil\n3. Update aplikasi ke versi terbaru\n4. Hapus cache aplikasi\n5. Restart perangkat\n\nJika masalah berlanjut, mohon kirimkan screenshot error beserta tipe perangkat dan versi OS Anda.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/technical/masalah-teknis-aplikasi	technical
219f350b-a329-4706-b1b5-2f99ff18962e	Tidak Bisa Login ke Akun	Kami memahami betapa frustrasinya tidak bisa masuk ke akun. Mari kami bantu! 🔑\n\nSilakan coba langkah berikut:\n1. Pastikan email dan password sudah benar\n2. Aktifkan Caps Lock/periksa huruf kapital\n3. Gunakan fitur "Lupa Password" untuk reset\n4. Coba login dari perangkat atau browser berbeda\n\nApakah Anda menerima pesan error tertentu? Mohon informasikan agar kami dapat membantu lebih lanjut.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/technical/tidak-bisa-login	technical
16535fd9-4831-4e45-9346-53e59d61990e	Fitur Sedang Dalam Perbaikan	Kami menginformasikan bahwa fitur {feature_name} saat ini sedang dalam perbaikan/maintenance. 🔧\n\nEstimasi selesai: {maintenance_end}\n\nKami mohon maaf atas ketidaknyamanannya. Tim teknis kami sedang bekerja keras untuk memulihkan layanan secepatnya.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/technical/fitur-dalam-perbaikan	technical
98729e07-8f8e-4782-beb5-b79db3e943f8	Verifikasi Identitas Diperlukan	Untuk keamanan akun Anda, kami perlu melakukan verifikasi identitas. 🔐\n\nDokumen yang diperlukan:\n📄 KTP (foto depan, jelas dan tidak blur)\n🤳 Selfie sambil memegang KTP\n\nKirimkan dokumen melalui email ke: verify@umrahmart.id\nSubject: Verifikasi Identitas - {user_id}\n\nProses verifikasi membutuhkan 1–2 hari kerja.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/verification/identitas-diperlukan	verification
3ed0230d-d315-4ed7-9d87-269803e4bfcc	Verifikasi Berhasil	Selamat! Verifikasi akun Anda telah berhasil. ✅\n\nAkun Anda kini sudah terverifikasi dan dapat menikmati semua fitur UmrahMart tanpa batasan. Terima kasih atas kerja samanya.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/verification/berhasil	verification
1b3e1789-f74a-4e30-a3e8-a8a761bba507	Informasi Pendaftaran Vendor	Terima kasih atas minat Anda untuk bergabung sebagai vendor di UmrahMart! 🏪\n\nPersyaratan pendaftaran:\n✅ KTP pemilik usaha\n✅ Foto toko/tempat usaha\n✅ Logo bisnis\n✅ Banner bisnis\n✅ Buku rekening bank\n✅ NPWP (opsional)\n\nProses review membutuhkan 2–3 hari kerja setelah semua dokumen lengkap.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/vendor/info-pendaftaran	vendor
593511f3-93a3-440f-933f-1f987d5df330	Status Vendor Sedang Direview	Pendaftaran vendor Anda sedang dalam proses review oleh tim kami. ⏳\n\nKami akan menginformasikan hasilnya melalui email dan notifikasi aplikasi dalam 2–3 hari kerja.\n\nPastikan data dan dokumen yang Anda unggah sudah lengkap dan terbaca dengan jelas.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/vendor/sedang-direview	vendor
1109f196-4b72-4548-b068-6f7f67ffd62b	Vendor Berhasil Diaktifkan	Selamat! Akun vendor Anda telah berhasil diaktifkan! 🎉\n\nAnda sekarang dapat mulai:\n🛍️ Menambahkan produk ke katalog\n📊 Mengelola stok dan harga\n📦 Menerima dan memproses pesanan\n\nUntuk panduan lengkap, kunjungi halaman "Panduan Vendor" di dashboard Anda. Semoga sukses!	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/vendor/berhasil-diaktifkan	vendor
ea325369-0913-474d-a5e4-4f17ca5a125e	Notifikasi Stok Tersedia	Kabar gembira! 🎉 Produk {product_name} yang Anda tunggu-tunggu kini sudah tersedia kembali.\n\nStok tersisa: {remaining_stock} pcs\n\nSegera pesan sebelum kehabisan! Klik di sini: {product_url}	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/stock/stok-tersedia	stock
676d2aa6-973c-4409-a4d9-1a3ce91a56e8	Informasi Pre-Order	Produk {product_name} saat ini tersedia dalam mode Pre-Order. 📋\n\nDetail Pre-Order:\n📅 Batas pemesanan : {po_end_date}\n🚚 Estimasi kirim  : {estimated_delivery}\n💰 Harga PO        : {po_price}\n\nPesan sekarang untuk mendapatkan kepastian stok!	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/stock/info-pre-order	stock
aaa4463e-a5ad-45e4-a1c7-0fd58c0d8f91	Konfirmasi Pembatalan Pesanan	Kami telah menerima permintaan pembatalan untuk pesanan {order_number}. ℹ️\n\nUntuk melanjutkan proses pembatalan, mohon konfirmasi alasan Anda:\na) Salah alamat pengiriman\nb) Ingin ganti produk\nc) Menemukan harga lebih murah\nd) Berubah pikiran\ne) Lainnya\n\nBalas dengan huruf pilihan Anda.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/cancellation/konfirmasi-pembatalan	cancellation
0da87072-6206-424e-96bc-ba5c186e10aa	Pesanan Tidak Dapat Dibatalkan	Mohon maaf, pesanan {order_number} tidak dapat dibatalkan karena status pesanan sudah dalam tahap {order_status}. ⚠️\n\nJika barang sudah diterima dan terdapat masalah, Anda masih bisa mengajukan return dalam 7 hari setelah penerimaan.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/cancellation/tidak-dapat-dibatalkan	cancellation
da5f8b53-3f0c-4e0e-9ad3-dd744be612e2	Pembatalan Berhasil	Pesanan {order_number} berhasil dibatalkan. ✅\n\nJika ada pembayaran yang perlu dikembalikan:\n💰 Jumlah refund : {refund_amount}\n⏱️ Estimasi dana masuk : 3–7 hari kerja\n\nTerima kasih atas pengertian Anda.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/cancellation/berhasil	cancellation
f267dbd9-8120-4b47-936f-dbd92ec618c2	Jadwal Estimasi Pengiriman	Berikut estimasi pengiriman berdasarkan lokasi Anda:\n\n🏙️ Jabodetabek    : 1–2 hari kerja\n🗺️ Pulau Jawa     : 2–3 hari kerja\n🏝️ Luar Jawa      : 3–5 hari kerja\n🗾 Papua/Terpencil : 5–10 hari kerja\n\nEstimasi dihitung sejak paket diserahkan ke kurir, tidak termasuk hari libur nasional.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/delivery/estimasi-pengiriman	delivery
287b0854-7fe7-4364-8942-42b56d8834ed	Pengiriman Tertunda	Kami mohon maaf atas keterlambatan pengiriman paket Anda (no. resi: {tracking_number}). 🙏\n\nKendala yang terjadi: {delay_reason}\n\nTim kami sedang berkoordinasi dengan pihak kurir {courier_name} untuk mempercepat proses pengiriman. Kami akan menginformasikan update terbaru sesegera mungkin.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/delivery/pengiriman-tertunda	delivery
6b5e9418-79f8-436c-837d-83fe58a6ecb1	Informasi Pengambilan di Toko (Pickup)	Pesanan Anda dapat diambil langsung di:\n\n📍 {store_address}\n🕐 Jam operasional: {store_hours}\n\nTunjukkan nomor pesanan {order_number} atau kode QR di aplikasi kepada petugas. Pesanan akan disimpan selama 3 hari sejak notifikasi siap diambil.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/pickup/info-pickup	pickup
bf72ddda-2258-4c8d-99e8-73198a54d6cf	Pesanan Siap Diambil	Pesanan Anda sudah siap untuk diambil! 🛍️\n\nNomor Pesanan : {order_number}\nLokasi        : {store_address}\nBatas Ambil   : {pickup_deadline}\n\nJangan lupa bawa identitas diri saat pengambilan. Sampai jumpa!	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/pickup/siap-diambil	pickup
81d65fdc-ec2a-4ec9-8838-97163f56c016	Informasi Paket Langganan	Berikut informasi paket langganan UmrahMart Premium:\n\n⭐ Paket Basic   : Rp 29.000/bulan — Gratis ongkir 2x/bulan\n⭐ Paket Premium : Rp 59.000/bulan — Gratis ongkir unlimited + cashback 5%\n⭐ Paket VIP     : Rp 99.000/bulan — Semua benefit + akses produk eksklusif\n\nDaftar sekarang di menu "Langganan" pada aplikasi.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/subscription/info-langganan	subscription
0033c4fc-1d54-41d2-8778-1e4f71e10a90	Masa Aktif Langganan Hampir Habis	Halo! Masa aktif langganan UmrahMart Premium Anda akan berakhir pada {expiry_date}. ⏰\n\nPerpanjang sekarang dan nikmati terus benefit:\n✅ Gratis ongkos kirim\n✅ Cashback eksklusif member\n✅ Early access produk baru\n\nKlik di sini untuk perpanjang: {renewal_url}	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/subscription/hampir-habis	subscription
15081062-8c4a-41f7-b4cf-9d2b28eb5264	Permintaan Ulasan Produk	Assalamu'alaikum! Semoga pesanan Anda sudah diterima dengan baik. 😊\n\nKami akan sangat berterima kasih jika Anda meluangkan waktu untuk memberikan ulasan pada produk {product_name}.\n\nUlasan Anda sangat membantu calon pembeli lain dalam memilih produk yang tepat. Berikan ulasan di menu "Pesanan Saya" → "Beri Ulasan".	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/review/minta-ulasan	review
e3d2e475-5166-4313-8b52-5ff2de44c14c	Terima Kasih Atas Ulasan Anda	Jazakallahu khairan atas ulasan yang telah Anda berikan! 🌟\n\nMasukan Anda sangat berharga bagi kami untuk terus meningkatkan kualitas produk dan layanan. Semoga UmrahMart selalu bisa menjadi teman setia perjalanan ibadah Anda.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/review/terima-kasih-ulasan	review
4eea82bf-4b06-4152-943d-e233d7526699	Laporan Penipuan Diterima	Laporan penipuan Anda telah kami terima dan kami tangani dengan sangat serius. 🚨\n\nNomor laporan: {report_number}\n\nTim keamanan kami akan menginvestigasi dalam 1×24 jam. Jangan lakukan transfer atau berikan data sensitif kepada pihak yang mencurigakan.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/fraud/laporan-penipuan-diterima	fraud
65d11440-cb07-45bf-bffb-5098c3186242	Tindakan Keamanan Akun	Kami mendeteksi aktivitas tidak biasa pada akun Anda. 🔐\n\nSebagai tindakan keamanan, akun Anda telah kami sementara kunci.\n\nLangkah yang perlu Anda lakukan:\n1. Segera ubah password Anda\n2. Aktifkan verifikasi 2 langkah\n3. Hubungi kami melalui {support_contact} untuk membuka kunci akun\n\nJangan pernah bagikan OTP atau password kepada siapapun, termasuk tim kami.	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/fraud/tindakan-keamanan-akun	fraud
d114e9ae-ca54-44db-b918-c64b8b8a4f8b	Pertanyaan Tidak Dapat Kami Jawab Saat Ini	Terima kasih atas pertanyaan Anda. Pertanyaan ini membutuhkan penanganan dari tim khusus kami. 🔄\n\nKami akan meneruskan pertanyaan ini ke tim terkait dan akan menghubungi Anda kembali dalam 1×24 jam kerja.\n\nNomor tiket Anda: {ticket_number}	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/other/eskalasi-ke-tim-terkait	other
3dd5915f-349b-400d-865f-4e71c1556127	Cara Menghubungi Kami	Ada beberapa cara untuk menghubungi tim UmrahMart:\n\n💬 Live Chat  : Aplikasi UmrahMart (menu "Bantuan")\n📧 Email      : support@umrahmart.id\n📞 Telepon    : 021-XXXX-XXXX (08.00–21.00 WIB)\n📱 WhatsApp   : wa.me/628XXXXXXXXX\n📘 Instagram  : @umrahmart.id\n\nKami siap membantu Anda! 🙏	t	cc000001-0000-0000-0000-000000000001	\N	2026-03-06 06:19:54.454902+00	2026-03-06 06:19:54.454902+00	/other/cara-menghubungi-kami	other
\.


--
-- Data for Name: return_reasons; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.return_reasons (id, reason, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: roles; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.roles (id, code, name) FROM stdin;
1	admin	Administrator
2	umkm	UMKM
3	customer	Customer
4	cs	Customer Service
5	finance	Finance
\.


--
-- Data for Name: schema_migrations; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.schema_migrations (version, dirty) FROM stdin;
33	f
\.


--
-- Data for Name: shipments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.shipments (id, order_id, courier_code, service_type, tracking_no, shipment_status, shipped_at, delivered_at, komship_order_id, komship_order_no, created_at, updated_at) FROM stdin;
9712d2df-3af8-4ad6-94ef-fadf2b32cde6	2413f18e-1462-4a6c-8696-d1d3f23494ea	jnt	EZ		waiting_pickup	\N	\N	\N	\N	2026-03-06 06:07:20.127604+00	2026-03-06 06:07:20.127604+00
f182918e-8262-436a-a213-f9e1a51814f5	d75fcc53-4e31-4245-a7d4-865e41ac27b1	jnt	EZ		waiting_pickup	\N	\N	\N	\N	2026-03-06 06:12:22.274566+00	2026-03-06 06:12:22.274566+00
7e22f896-8f65-4e38-9ce5-e40e452b0671	4354c494-e4dc-4b1b-9f2e-94dfd8e42ff8	jnt	EZ		waiting_pickup	\N	\N	\N	\N	2026-03-06 06:39:22.176503+00	2026-03-06 06:39:22.176503+00
9a6ac9ba-8e08-43b7-975a-55900b8bcd78	3df49783-a642-4a92-a95a-86de5ae501e4	jnt	EZ		waiting_pickup	\N	\N	\N	\N	2026-03-09 03:37:20.114567+00	2026-03-09 03:37:20.114567+00
\.


--
-- Data for Name: ticket_attachments; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ticket_attachments (id, ticket_id, message_id, file_url, file_name, file_type, file_size, created_at) FROM stdin;
\.


--
-- Data for Name: ticket_messages; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ticket_messages (id, ticket_id, sender_id, message, is_from_cs, is_internal_note, created_at) FROM stdin;
\.


--
-- Data for Name: ticket_status_logs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ticket_status_logs (id, ticket_id, changed_by, old_status, new_status, notes, created_at) FROM stdin;
\.


--
-- Data for Name: ticket_subjects; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.ticket_subjects (id, label, is_active, created_at) FROM stdin;
1	Barang Tidak Sampai	t	2026-03-05 09:05:37.479599+00
2	Barang Rusak	t	2026-03-05 09:05:37.479599+00
3	Pengembalian Dana	t	2026-03-05 09:05:37.479599+00
4	Barang Tidak Sesuai	t	2026-03-05 09:05:37.479599+00
5	Pengiriman Terlambat	t	2026-03-05 09:05:37.479599+00
6	Pembatalan Pesanan	t	2026-03-05 09:05:37.479599+00
7	Kesalahan Produk	t	2026-03-05 09:05:37.479599+00
8	Akun Bermasalah	t	2026-03-05 09:05:37.479599+00
9	Pembayaran Gagal	t	2026-03-05 09:05:37.479599+00
10	Voucher/Promo Tidak Berlaku	t	2026-03-05 09:05:37.479599+00
11	Pertanyaan Umum	t	2026-03-05 09:05:37.479599+00
12	Lainnya	t	2026-03-05 09:05:37.479599+00
\.


--
-- Data for Name: tickets; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.tickets (id, ticket_number, customer_id, assigned_cs_id, order_number, phone, reporter_name, subject, detail, status, source, resolved_at, closed_at, created_at, updated_at, attachment_url, attachment_content_type, subject_id) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.users (id, email, full_name, birth_date, phone, password_hash, role_id, status, email_verified_at, created_at, updated_at) FROM stdin;
aa000002-0000-0000-0000-000000000001	umkm2@dev.local	Toko Oleh-Oleh Tanah Suci	\N	081299000002	$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y	2	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
aa000002-0000-0000-0000-000000000002	umkm3@dev.local	Perlengkapan Ibadah Barokah	\N	081299000003	$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y	2	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
aa000002-0000-0000-0000-000000000003	customer2@dev.local	Siti Aisyah	\N	081299000004	$2a$10$xB1d32uugV6/A4Qd7UXA2.zWk67QihWS1V7zfJvrWFA2joqHVGG9y	3	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
7907d953-ba02-40a4-b403-07b68d8471e8	admin@example.com	Admin	\N	08120000000	$2a$10$dc41C09en634mTZdCnRZiejmB3Z7alOhfOh3r/xsPCzpD4ObMQitu	1	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d	umkm@dev.local	Dev UMKM	\N	081200000001	$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq	2	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
7d281d37-8f4a-48a8-9b33-e47c76b3d13c	customer@dev.local	Dev Customer	\N	081200000002	$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq	3	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
cc000001-0000-0000-0000-000000000001	cs@dev.local	Dev Customer Service 1	\N	081200000003	$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq	4	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
cc000001-0000-0000-0000-000000000002	cs2@dev.local	Dev Customer Service 2	\N	081200000005	$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq	4	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
cc000001-0000-0000-0000-000000000003	cs3@dev.local	Dev Customer Service 3	\N	081200000006	$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq	4	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
ff000001-0000-0000-0000-000000000001	finance@dev.local	Dev Finance	\N	081200000004	$2a$10$1Sv.ravjLnQNaGpTB0e9PO0N34Pnqnr0a5nwsDq7yZmARc7cI.WWq	5	active	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
\.


--
-- Data for Name: vendor_balances; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vendor_balances (vendor_id, available_balance, pending_balance, total_earned, total_withdrawn, updated_at) FROM stdin;
da000001-0000-0000-0000-000000000001	0.00	0.00	0.00	0.00	2026-03-06 01:58:19.265177+00
da000001-0000-0000-0000-000000000002	0.00	0.00	0.00	0.00	2026-03-06 01:58:19.265177+00
da000001-0000-0000-0000-000000000003	0.00	0.00	0.00	0.00	2026-03-06 01:58:19.265177+00
\.


--
-- Data for Name: vendor_bank_accounts; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vendor_bank_accounts (id, vendor_id, bank_name, account_number, account_holder_name, verification_status, rejection_reason, verified_by, verified_at, created_at, updated_at) FROM stdin;
db000001-0000-0000-0000-000000000001	da000001-0000-0000-0000-000000000001	Bank Syariah Indonesia	7788001122	PT Ihram Jaya Nusantara	pending	\N	\N	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
db000001-0000-0000-0000-000000000002	da000001-0000-0000-0000-000000000002	Bank Muamalat	1122334455	CV Oleh-Oleh Tanah Suci	pending	\N	\N	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
db000001-0000-0000-0000-000000000003	da000001-0000-0000-0000-000000000003	Bank BCA Syariah	5566001122	UD Barokah Store	pending	\N	\N	\N	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
\.


--
-- Data for Name: vendor_banners; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vendor_banners (id, vendor_id, title, image_url, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: vendor_couriers; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vendor_couriers (id, vendor_id, courier_id, is_active, created_at) FROM stdin;
f711b8ae-791c-436a-b12e-eca1de41326a	da000001-0000-0000-0000-000000000003	1	t	2026-03-06 02:11:02.534896+00
7c50a798-4335-48cb-8887-611ef6102f55	da000001-0000-0000-0000-000000000002	2	t	2026-03-06 02:12:54.909634+00
a9591f36-edd5-4b26-a0b7-d2e6b4545849	da000001-0000-0000-0000-000000000001	3	t	2026-03-09 03:51:26.108854+00
2dce4d62-1e1f-4850-9890-b76b001e3aa6	da000001-0000-0000-0000-000000000001	6	t	2026-03-09 03:51:26.109688+00
09e2cccd-76c4-4767-bfaf-77a73f130eea	da000001-0000-0000-0000-000000000001	1	t	2026-03-09 03:51:26.11001+00
e0407467-ad4f-4fac-a624-2ee910dd2a71	da000001-0000-0000-0000-000000000001	11	t	2026-03-09 03:51:26.110199+00
96bc39e4-bf67-4bb0-a863-2b31f41cf725	da000001-0000-0000-0000-000000000001	5	t	2026-03-09 03:51:26.110365+00
e5275084-a9ba-424d-a92b-666fb8e3c93a	da000001-0000-0000-0000-000000000001	9	t	2026-03-09 03:51:26.110545+00
fe387630-be24-44a8-9a99-643408d09f53	da000001-0000-0000-0000-000000000001	12	t	2026-03-09 03:51:26.110722+00
b0ee122a-efd0-4cbd-b5f6-c8b66e4ccd86	da000001-0000-0000-0000-000000000001	4	t	2026-03-09 03:51:26.110841+00
fff65272-25b2-46ed-b9b5-096b5a5d2549	da000001-0000-0000-0000-000000000001	10	t	2026-03-09 03:51:26.111+00
b2d73ada-2271-45c0-b829-fd63e007b5bf	da000001-0000-0000-0000-000000000001	2	t	2026-03-09 03:51:26.111162+00
9a604f86-ffb3-4c3f-a9d7-d61ef748b633	da000001-0000-0000-0000-000000000001	7	t	2026-03-09 03:51:26.111301+00
46b94a9e-d543-4a9a-a5dd-ed852d6d439b	da000001-0000-0000-0000-000000000001	8	t	2026-03-09 03:51:26.111418+00
\.


--
-- Data for Name: vendor_documents; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vendor_documents (id, vendor_id, doc_type, file_url, mime_type, file_size_bytes, file_checksum, uploaded_by, verification_status, rejection_reason, verified_by, verified_at, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: vendor_withdrawals; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vendor_withdrawals (id, vendor_id, amount, channel_code, status, xendit_payout_id, xendit_status, description, failed_reason, created_at, updated_at, fee_estimated, fee_actual, amount_net, total_deducted, fee_status, xendit_transaction_id) FROM stdin;
\.


--
-- Data for Name: vendors; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.vendors (id, owner_user_id, vendor_type, legal_name, display_name, responsible_person_name, description, status, approved_by, approved_at, status_reason, xendit_account_id, created_at, updated_at) FROM stdin;
da000001-0000-0000-0000-000000000002	aa000002-0000-0000-0000-000000000001	umrah_souvenir_store	CV Oleh-Oleh Tanah Suci	TOKO WAREHOUSE JAKARTA BARAT - TAMBORA	Toko Oleh-Oleh Tanah Suci	Pusat oleh-oleh haji dan umrah langsung dari Tanah Suci. Kurma premium, minyak zaitun asli, dan souvenir eksklusif Mekkah-Madinah.	active	7907d953-ba02-40a4-b403-07b68d8471e8	2026-03-06 01:58:19.265177+00	\N	fake-xendit-for-dev	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
da000001-0000-0000-0000-000000000003	aa000002-0000-0000-0000-000000000002	general_souvenir_store	UD Barokah Store	TOKO WAREHOUSE SIDOARJO - GEDANGAN	Perlengkapan Ibadah Barokah	Menyediakan berbagai perlengkapan ibadah berkualitas dengan harga terjangkau. Mulai dari mukena, tasbih, Al-Quran, hingga hampers haji umrah.	active	7907d953-ba02-40a4-b403-07b68d8471e8	2026-03-06 01:58:19.265177+00	\N	fake-xendit-for-dev	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
da000001-0000-0000-0000-000000000001	8ab0c9c7-4dac-4087-9e47-cd4c7317ef6d	hajj_souvenir_store	PT Ihram Jaya Nusantara	TOKO WAREHOUSE SURABAYA - WIYUNG	Dev UMKM	Toko perlengkapan haji & umrah terlengkap di Jakarta. Menyediakan kain ihram, sajadah, dan perlengkapan ibadah lainnya dengan kualitas premium.	active	\N	\N	\N	fake-xendit-for-dev	2026-03-06 01:58:19.265177+00	2026-03-06 01:58:19.265177+00
\.


--
-- Data for Name: wishlist_items; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.wishlist_items (id, user_id, product_id, created_at) FROM stdin;
\.


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
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: admin_contacts admin_contacts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_contacts
    ADD CONSTRAINT admin_contacts_pkey PRIMARY KEY (id);


--
-- Name: banners banners_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.banners
    ADD CONSTRAINT banners_pkey PRIMARY KEY (id);


--
-- Name: cart_items cart_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT cart_items_pkey PRIMARY KEY (id);


--
-- Name: carts carts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: categories categories_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_slug_key UNIQUE (slug);


--
-- Name: chat_conversations chat_conversations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_conversations
    ADD CONSTRAINT chat_conversations_pkey PRIMARY KEY (id);


--
-- Name: chat_messages chat_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT chat_messages_pkey PRIMARY KEY (id);


--
-- Name: couriers couriers_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.couriers
    ADD CONSTRAINT couriers_code_key UNIQUE (code);


--
-- Name: couriers couriers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.couriers
    ADD CONSTRAINT couriers_pkey PRIMARY KEY (id);


--
-- Name: email_verification_tokens email_verification_tokens_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT email_verification_tokens_pkey PRIMARY KEY (id);


--
-- Name: faqs faqs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.faqs
    ADD CONSTRAINT faqs_pkey PRIMARY KEY (id);


--
-- Name: ledger_accounts ledger_accounts_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_accounts
    ADD CONSTRAINT ledger_accounts_code_key UNIQUE (code);


--
-- Name: ledger_accounts ledger_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_accounts
    ADD CONSTRAINT ledger_accounts_pkey PRIMARY KEY (id);


--
-- Name: ledger_journals ledger_journals_journal_no_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_journals
    ADD CONSTRAINT ledger_journals_journal_no_key UNIQUE (journal_no);


--
-- Name: ledger_journals ledger_journals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_journals
    ADD CONSTRAINT ledger_journals_pkey PRIMARY KEY (id);


--
-- Name: ledger_lines ledger_lines_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_lines
    ADD CONSTRAINT ledger_lines_pkey PRIMARY KEY (id);


--
-- Name: notifications notifications_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT notifications_pkey PRIMARY KEY (id);


--
-- Name: order_items order_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT order_items_pkey PRIMARY KEY (id);


--
-- Name: order_status_history order_status_history_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_status_history
    ADD CONSTRAINT order_status_history_pkey PRIMARY KEY (id);


--
-- Name: orders orders_order_no_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_order_no_key UNIQUE (order_no);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: payment_events payment_events_external_event_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_events
    ADD CONSTRAINT payment_events_external_event_id_key UNIQUE (external_event_id);


--
-- Name: payment_events payment_events_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_events
    ADD CONSTRAINT payment_events_pkey PRIMARY KEY (id);


--
-- Name: payment_invoices payment_invoices_external_invoice_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_invoices
    ADD CONSTRAINT payment_invoices_external_invoice_id_key UNIQUE (external_invoice_id);


--
-- Name: payment_invoices payment_invoices_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_invoices
    ADD CONSTRAINT payment_invoices_pkey PRIMARY KEY (id);


--
-- Name: payout_batches payout_batches_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payout_batches
    ADD CONSTRAINT payout_batches_pkey PRIMARY KEY (id);


--
-- Name: payout_items payout_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payout_items
    ADD CONSTRAINT payout_items_pkey PRIMARY KEY (id);


--
-- Name: product_images product_images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_images
    ADD CONSTRAINT product_images_pkey PRIMARY KEY (id);


--
-- Name: product_review_images product_review_images_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_review_images
    ADD CONSTRAINT product_review_images_pkey PRIMARY KEY (id);


--
-- Name: product_review_stats product_review_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_review_stats
    ADD CONSTRAINT product_review_stats_pkey PRIMARY KEY (product_id);


--
-- Name: product_reviews product_reviews_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT product_reviews_pkey PRIMARY KEY (id);


--
-- Name: product_variants product_variants_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_variants
    ADD CONSTRAINT product_variants_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: products products_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_slug_key UNIQUE (slug);


--
-- Name: refunds refunds_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refunds
    ADD CONSTRAINT refunds_pkey PRIMARY KEY (id);


--
-- Name: reply_templates reply_templates_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reply_templates
    ADD CONSTRAINT reply_templates_pkey PRIMARY KEY (id);


--
-- Name: return_reasons return_reasons_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.return_reasons
    ADD CONSTRAINT return_reasons_pkey PRIMARY KEY (id);


--
-- Name: roles roles_code_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_code_key UNIQUE (code);


--
-- Name: roles roles_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.roles
    ADD CONSTRAINT roles_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: shipments shipments_order_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT shipments_order_id_key UNIQUE (order_id);


--
-- Name: shipments shipments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT shipments_pkey PRIMARY KEY (id);


--
-- Name: ticket_attachments ticket_attachments_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_attachments
    ADD CONSTRAINT ticket_attachments_pkey PRIMARY KEY (id);


--
-- Name: ticket_messages ticket_messages_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_messages
    ADD CONSTRAINT ticket_messages_pkey PRIMARY KEY (id);


--
-- Name: ticket_status_logs ticket_status_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_status_logs
    ADD CONSTRAINT ticket_status_logs_pkey PRIMARY KEY (id);


--
-- Name: ticket_subjects ticket_subjects_label_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_subjects
    ADD CONSTRAINT ticket_subjects_label_key UNIQUE (label);


--
-- Name: ticket_subjects ticket_subjects_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_subjects
    ADD CONSTRAINT ticket_subjects_pkey PRIMARY KEY (id);


--
-- Name: tickets tickets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_pkey PRIMARY KEY (id);


--
-- Name: tickets tickets_ticket_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_ticket_number_key UNIQUE (ticket_number);


--
-- Name: chat_conversations uq_chat_conv_pair; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_conversations
    ADD CONSTRAINT uq_chat_conv_pair UNIQUE (initiator_id, participant_id);


--
-- Name: product_reviews uq_product_reviews_user_order_item; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT uq_product_reviews_user_order_item UNIQUE (user_id, order_item_id);


--
-- Name: product_variants uq_product_variants_product_sku; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_variants
    ADD CONSTRAINT uq_product_variants_product_sku UNIQUE (product_id, sku);


--
-- Name: reply_templates uq_reply_templates_shortcut; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reply_templates
    ADD CONSTRAINT uq_reply_templates_shortcut UNIQUE (shortcut);


--
-- Name: vendors uq_vendors_owner_user; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT uq_vendors_owner_user UNIQUE (owner_user_id);


--
-- Name: users users_phone_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_phone_key UNIQUE (phone);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: vendor_balances vendor_balances_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_balances
    ADD CONSTRAINT vendor_balances_pkey PRIMARY KEY (vendor_id);


--
-- Name: vendor_bank_accounts vendor_bank_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_bank_accounts
    ADD CONSTRAINT vendor_bank_accounts_pkey PRIMARY KEY (id);


--
-- Name: vendor_bank_accounts vendor_bank_accounts_vendor_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_bank_accounts
    ADD CONSTRAINT vendor_bank_accounts_vendor_id_key UNIQUE (vendor_id);


--
-- Name: vendor_banners vendor_banners_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_banners
    ADD CONSTRAINT vendor_banners_pkey PRIMARY KEY (id);


--
-- Name: vendor_couriers vendor_couriers_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_couriers
    ADD CONSTRAINT vendor_couriers_pkey PRIMARY KEY (id);


--
-- Name: vendor_couriers vendor_couriers_vendor_id_courier_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_couriers
    ADD CONSTRAINT vendor_couriers_vendor_id_courier_id_key UNIQUE (vendor_id, courier_id);


--
-- Name: vendor_documents vendor_documents_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_documents
    ADD CONSTRAINT vendor_documents_pkey PRIMARY KEY (id);


--
-- Name: vendor_withdrawals vendor_withdrawals_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_withdrawals
    ADD CONSTRAINT vendor_withdrawals_pkey PRIMARY KEY (id);


--
-- Name: vendors vendors_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT vendors_pkey PRIMARY KEY (id);


--
-- Name: wishlist_items wishlist_items_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wishlist_items
    ADD CONSTRAINT wishlist_items_pkey PRIMARY KEY (id);


--
-- Name: idx_addresses_user_default; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_addresses_user_default ON public.addresses USING btree (user_id, is_default) WHERE (is_default = true);


--
-- Name: idx_carts_user_status_updated_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_carts_user_status_updated_at ON public.carts USING btree (user_id, status, updated_at DESC);


--
-- Name: idx_chat_conv_initiator_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_conv_initiator_id ON public.chat_conversations USING btree (initiator_id);


--
-- Name: idx_chat_conv_last_msg; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_conv_last_msg ON public.chat_conversations USING btree (last_message_at DESC);


--
-- Name: idx_chat_conv_participant_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_conv_participant_id ON public.chat_conversations USING btree (participant_id);


--
-- Name: idx_chat_messages_conv_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_conv_id ON public.chat_messages USING btree (conversation_id);


--
-- Name: idx_chat_messages_created; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_created ON public.chat_messages USING btree (conversation_id, created_at);


--
-- Name: idx_chat_messages_unread; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_chat_messages_unread ON public.chat_messages USING btree (conversation_id, is_read) WHERE (NOT is_read);


--
-- Name: idx_email_verification_tokens_expires_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_email_verification_tokens_expires_at ON public.email_verification_tokens USING btree (expires_at);


--
-- Name: idx_email_verification_tokens_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_email_verification_tokens_user_created_at ON public.email_verification_tokens USING btree (user_id, created_at DESC);


--
-- Name: idx_faqs_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_faqs_category ON public.faqs USING btree (category);


--
-- Name: idx_ledger_journals_source; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ledger_journals_source ON public.ledger_journals USING btree (source_type, source_id);


--
-- Name: idx_ledger_lines_account; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ledger_lines_account ON public.ledger_lines USING btree (account_id);


--
-- Name: idx_notifications_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_user_created_at ON public.notifications USING btree (user_id, created_at DESC);


--
-- Name: idx_notifications_user_read_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_notifications_user_read_at ON public.notifications USING btree (user_id, read_at);


--
-- Name: idx_orders_status_payment_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_orders_status_payment_status ON public.orders USING btree (order_status, payment_status);


--
-- Name: idx_orders_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_orders_user_created_at ON public.orders USING btree (user_id, created_at);


--
-- Name: idx_orders_vendor_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_orders_vendor_created_at ON public.orders USING btree (vendor_id, created_at);


--
-- Name: idx_payment_invoices_order_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payment_invoices_order_status ON public.payment_invoices USING btree (order_id, status);


--
-- Name: idx_payout_batches_vendor_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_payout_batches_vendor_status ON public.payout_batches USING btree (vendor_id, status);


--
-- Name: idx_product_images_product_sort; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_images_product_sort ON public.product_images USING btree (product_id, sort_order);


--
-- Name: idx_product_review_images_review_sort; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_review_images_review_sort ON public.product_review_images USING btree (review_id, sort_order);


--
-- Name: idx_product_reviews_product_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_reviews_product_created_at ON public.product_reviews USING btree (product_id, created_at DESC);


--
-- Name: idx_product_reviews_product_rating; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_reviews_product_rating ON public.product_reviews USING btree (product_id, rating);


--
-- Name: idx_product_reviews_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_reviews_user_created_at ON public.product_reviews USING btree (user_id, created_at DESC);


--
-- Name: idx_product_variants_product_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_product_variants_product_active ON public.product_variants USING btree (product_id, is_active);


--
-- Name: idx_products_vendor_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_products_vendor_status ON public.products USING btree (vendor_id, status);


--
-- Name: idx_reply_templates_active; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reply_templates_active ON public.reply_templates USING btree (is_active) WHERE (is_active = true);


--
-- Name: idx_reply_templates_category; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reply_templates_category ON public.reply_templates USING btree (category);


--
-- Name: idx_reply_templates_shortcut; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_reply_templates_shortcut ON public.reply_templates USING btree (shortcut);


--
-- Name: idx_ticket_attachments_message_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ticket_attachments_message_id ON public.ticket_attachments USING btree (message_id);


--
-- Name: idx_ticket_attachments_ticket_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ticket_attachments_ticket_id ON public.ticket_attachments USING btree (ticket_id);


--
-- Name: idx_ticket_messages_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ticket_messages_created_at ON public.ticket_messages USING btree (ticket_id, created_at);


--
-- Name: idx_ticket_messages_ticket_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ticket_messages_ticket_id ON public.ticket_messages USING btree (ticket_id);


--
-- Name: idx_ticket_status_logs_ticket_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_ticket_status_logs_ticket_id ON public.ticket_status_logs USING btree (ticket_id);


--
-- Name: idx_tickets_assigned_cs; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tickets_assigned_cs ON public.tickets USING btree (assigned_cs_id);


--
-- Name: idx_tickets_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tickets_created_at ON public.tickets USING btree (created_at DESC);


--
-- Name: idx_tickets_customer_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tickets_customer_id ON public.tickets USING btree (customer_id);


--
-- Name: idx_tickets_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tickets_status ON public.tickets USING btree (status);


--
-- Name: idx_tickets_subject_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_tickets_subject_id ON public.tickets USING btree (subject_id);


--
-- Name: idx_vendor_bank_accounts_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_bank_accounts_status ON public.vendor_bank_accounts USING btree (verification_status);


--
-- Name: idx_vendor_banners_vendor_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_banners_vendor_id ON public.vendor_banners USING btree (vendor_id);


--
-- Name: idx_vendor_couriers_vendor; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_couriers_vendor ON public.vendor_couriers USING btree (vendor_id);


--
-- Name: idx_vendor_documents_vendor_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_documents_vendor_status ON public.vendor_documents USING btree (vendor_id, verification_status);


--
-- Name: idx_vendor_withdrawals_fee_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_withdrawals_fee_status ON public.vendor_withdrawals USING btree (fee_status);


--
-- Name: idx_vendor_withdrawals_status; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_withdrawals_status ON public.vendor_withdrawals USING btree (status);


--
-- Name: idx_vendor_withdrawals_vendor_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_withdrawals_vendor_id ON public.vendor_withdrawals USING btree (vendor_id);


--
-- Name: idx_vendor_withdrawals_xendit_payout_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_vendor_withdrawals_xendit_payout_id ON public.vendor_withdrawals USING btree (xendit_payout_id);


--
-- Name: idx_wishlist_items_product; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_wishlist_items_product ON public.wishlist_items USING btree (product_id);


--
-- Name: idx_wishlist_items_user_created_at; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_wishlist_items_user_created_at ON public.wishlist_items USING btree (user_id, created_at DESC, id DESC);


--
-- Name: uq_cart_items_cart_variant; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_cart_items_cart_variant ON public.cart_items USING btree (cart_id, product_variant_id);


--
-- Name: uq_carts_user_active; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_carts_user_active ON public.carts USING btree (user_id) WHERE ((status)::text = 'active'::text);


--
-- Name: uq_email_verification_tokens_token_hash; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_email_verification_tokens_token_hash ON public.email_verification_tokens USING btree (token_hash);


--
-- Name: uq_email_verification_tokens_user_active; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_email_verification_tokens_user_active ON public.email_verification_tokens USING btree (user_id) WHERE ((consumed_at IS NULL) AND (invalidated_at IS NULL));


--
-- Name: uq_payment_invoices_xendit_id; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_payment_invoices_xendit_id ON public.payment_invoices USING btree (xendit_invoice_id) WHERE (xendit_invoice_id IS NOT NULL);


--
-- Name: uq_payout_items_batch_order; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_payout_items_batch_order ON public.payout_items USING btree (payout_batch_id, order_id);


--
-- Name: uq_product_images_primary_per_product; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_product_images_primary_per_product ON public.product_images USING btree (product_id) WHERE is_primary;


--
-- Name: uq_product_variants_default_per_product; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_product_variants_default_per_product ON public.product_variants USING btree (product_id) WHERE is_default;


--
-- Name: uq_users_email_lower; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_users_email_lower ON public.users USING btree (lower((email)::text));


--
-- Name: uq_vendor_documents_vendor_doc_type; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_vendor_documents_vendor_doc_type ON public.vendor_documents USING btree (vendor_id, doc_type);


--
-- Name: uq_wishlist_items_user_product; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX uq_wishlist_items_user_product ON public.wishlist_items USING btree (user_id, product_id);


--
-- Name: product_images trg_product_images_max_10; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_product_images_max_10 BEFORE INSERT OR UPDATE OF product_id ON public.product_images FOR EACH ROW EXECUTE FUNCTION public.enforce_product_image_limit();


--
-- Name: addresses fk_addresses_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT fk_addresses_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: cart_items fk_cart_items_cart; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES public.carts(id);


--
-- Name: cart_items fk_cart_items_product_variant; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cart_items
    ADD CONSTRAINT fk_cart_items_product_variant FOREIGN KEY (product_variant_id) REFERENCES public.product_variants(id);


--
-- Name: carts fk_carts_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT fk_carts_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: categories fk_categories_parent; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT fk_categories_parent FOREIGN KEY (parent_id) REFERENCES public.categories(id);


--
-- Name: chat_conversations fk_chat_conv_initiator; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_conversations
    ADD CONSTRAINT fk_chat_conv_initiator FOREIGN KEY (initiator_id) REFERENCES public.users(id);


--
-- Name: chat_conversations fk_chat_conv_participant; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_conversations
    ADD CONSTRAINT fk_chat_conv_participant FOREIGN KEY (participant_id) REFERENCES public.users(id);


--
-- Name: chat_messages fk_cmsg_conv; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT fk_cmsg_conv FOREIGN KEY (conversation_id) REFERENCES public.chat_conversations(id) ON DELETE CASCADE;


--
-- Name: chat_messages fk_cmsg_sender; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.chat_messages
    ADD CONSTRAINT fk_cmsg_sender FOREIGN KEY (sender_id) REFERENCES public.users(id);


--
-- Name: email_verification_tokens fk_email_verification_tokens_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.email_verification_tokens
    ADD CONSTRAINT fk_email_verification_tokens_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: ledger_journals fk_ledger_journals_created_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_journals
    ADD CONSTRAINT fk_ledger_journals_created_by FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: ledger_lines fk_ledger_lines_account; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_lines
    ADD CONSTRAINT fk_ledger_lines_account FOREIGN KEY (account_id) REFERENCES public.ledger_accounts(id);


--
-- Name: ledger_lines fk_ledger_lines_journal; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ledger_lines
    ADD CONSTRAINT fk_ledger_lines_journal FOREIGN KEY (journal_id) REFERENCES public.ledger_journals(id);


--
-- Name: notifications fk_notifications_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.notifications
    ADD CONSTRAINT fk_notifications_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE;


--
-- Name: order_items fk_order_items_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT fk_order_items_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: order_items fk_order_items_product_variant; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_items
    ADD CONSTRAINT fk_order_items_product_variant FOREIGN KEY (product_variant_id) REFERENCES public.product_variants(id);


--
-- Name: order_status_history fk_order_status_history_changed_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_status_history
    ADD CONSTRAINT fk_order_status_history_changed_by FOREIGN KEY (changed_by) REFERENCES public.users(id);


--
-- Name: order_status_history fk_order_status_history_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.order_status_history
    ADD CONSTRAINT fk_order_status_history_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: orders fk_orders_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: orders fk_orders_vendor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT fk_orders_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: payment_events fk_payment_events_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_events
    ADD CONSTRAINT fk_payment_events_invoice FOREIGN KEY (payment_invoice_id) REFERENCES public.payment_invoices(id);


--
-- Name: payment_invoices fk_payment_invoices_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payment_invoices
    ADD CONSTRAINT fk_payment_invoices_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: payout_batches fk_payout_batches_created_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payout_batches
    ADD CONSTRAINT fk_payout_batches_created_by FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: payout_batches fk_payout_batches_vendor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payout_batches
    ADD CONSTRAINT fk_payout_batches_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: payout_items fk_payout_items_batch; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payout_items
    ADD CONSTRAINT fk_payout_items_batch FOREIGN KEY (payout_batch_id) REFERENCES public.payout_batches(id);


--
-- Name: payout_items fk_payout_items_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.payout_items
    ADD CONSTRAINT fk_payout_items_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: product_images fk_product_images_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_images
    ADD CONSTRAINT fk_product_images_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: product_review_images fk_product_review_images_review; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_review_images
    ADD CONSTRAINT fk_product_review_images_review FOREIGN KEY (review_id) REFERENCES public.product_reviews(id) ON DELETE CASCADE;


--
-- Name: product_review_stats fk_product_review_stats_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_review_stats
    ADD CONSTRAINT fk_product_review_stats_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: product_reviews fk_product_reviews_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT fk_product_reviews_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: product_reviews fk_product_reviews_order_item; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT fk_product_reviews_order_item FOREIGN KEY (order_item_id) REFERENCES public.order_items(id);


--
-- Name: product_reviews fk_product_reviews_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT fk_product_reviews_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: product_reviews fk_product_reviews_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_reviews
    ADD CONSTRAINT fk_product_reviews_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: product_variants fk_product_variants_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.product_variants
    ADD CONSTRAINT fk_product_variants_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: products fk_products_category; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: products fk_products_vendor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT fk_products_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: refunds fk_refunds_invoice; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refunds
    ADD CONSTRAINT fk_refunds_invoice FOREIGN KEY (payment_invoice_id) REFERENCES public.payment_invoices(id);


--
-- Name: refunds fk_refunds_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refunds
    ADD CONSTRAINT fk_refunds_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: refunds fk_refunds_processed_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refunds
    ADD CONSTRAINT fk_refunds_processed_by FOREIGN KEY (processed_by) REFERENCES public.users(id);


--
-- Name: refunds fk_refunds_requested_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.refunds
    ADD CONSTRAINT fk_refunds_requested_by FOREIGN KEY (requested_by) REFERENCES public.users(id);


--
-- Name: reply_templates fk_rtpl_creator; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reply_templates
    ADD CONSTRAINT fk_rtpl_creator FOREIGN KEY (created_by) REFERENCES public.users(id);


--
-- Name: reply_templates fk_rtpl_updater; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.reply_templates
    ADD CONSTRAINT fk_rtpl_updater FOREIGN KEY (updated_by) REFERENCES public.users(id);


--
-- Name: shipments fk_shipments_order; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.shipments
    ADD CONSTRAINT fk_shipments_order FOREIGN KEY (order_id) REFERENCES public.orders(id);


--
-- Name: ticket_attachments fk_tattach_message; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_attachments
    ADD CONSTRAINT fk_tattach_message FOREIGN KEY (message_id) REFERENCES public.ticket_messages(id) ON DELETE SET NULL;


--
-- Name: ticket_attachments fk_tattach_ticket; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_attachments
    ADD CONSTRAINT fk_tattach_ticket FOREIGN KEY (ticket_id) REFERENCES public.tickets(id) ON DELETE CASCADE;


--
-- Name: tickets fk_tickets_cs; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT fk_tickets_cs FOREIGN KEY (assigned_cs_id) REFERENCES public.users(id);


--
-- Name: tickets fk_tickets_customer; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT fk_tickets_customer FOREIGN KEY (customer_id) REFERENCES public.users(id);


--
-- Name: ticket_messages fk_tmsg_sender; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_messages
    ADD CONSTRAINT fk_tmsg_sender FOREIGN KEY (sender_id) REFERENCES public.users(id);


--
-- Name: ticket_messages fk_tmsg_ticket; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_messages
    ADD CONSTRAINT fk_tmsg_ticket FOREIGN KEY (ticket_id) REFERENCES public.tickets(id) ON DELETE CASCADE;


--
-- Name: ticket_status_logs fk_tslog_ticket; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_status_logs
    ADD CONSTRAINT fk_tslog_ticket FOREIGN KEY (ticket_id) REFERENCES public.tickets(id) ON DELETE CASCADE;


--
-- Name: ticket_status_logs fk_tslog_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.ticket_status_logs
    ADD CONSTRAINT fk_tslog_user FOREIGN KEY (changed_by) REFERENCES public.users(id);


--
-- Name: users fk_users_role; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT fk_users_role FOREIGN KEY (role_id) REFERENCES public.roles(id);


--
-- Name: vendor_balances fk_vendor_balances_vendor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_balances
    ADD CONSTRAINT fk_vendor_balances_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: vendor_bank_accounts fk_vendor_bank_accounts_vendor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_bank_accounts
    ADD CONSTRAINT fk_vendor_bank_accounts_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: vendor_bank_accounts fk_vendor_bank_accounts_verified_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_bank_accounts
    ADD CONSTRAINT fk_vendor_bank_accounts_verified_by FOREIGN KEY (verified_by) REFERENCES public.users(id);


--
-- Name: vendor_documents fk_vendor_documents_uploaded_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_documents
    ADD CONSTRAINT fk_vendor_documents_uploaded_by FOREIGN KEY (uploaded_by) REFERENCES public.users(id);


--
-- Name: vendor_documents fk_vendor_documents_vendor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_documents
    ADD CONSTRAINT fk_vendor_documents_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: vendor_documents fk_vendor_documents_verified_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_documents
    ADD CONSTRAINT fk_vendor_documents_verified_by FOREIGN KEY (verified_by) REFERENCES public.users(id);


--
-- Name: vendor_withdrawals fk_vendor_withdrawals_vendor; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_withdrawals
    ADD CONSTRAINT fk_vendor_withdrawals_vendor FOREIGN KEY (vendor_id) REFERENCES public.vendors(id);


--
-- Name: vendors fk_vendors_approved_by; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT fk_vendors_approved_by FOREIGN KEY (approved_by) REFERENCES public.users(id);


--
-- Name: vendors fk_vendors_owner; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendors
    ADD CONSTRAINT fk_vendors_owner FOREIGN KEY (owner_user_id) REFERENCES public.users(id);


--
-- Name: wishlist_items fk_wishlist_items_product; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wishlist_items
    ADD CONSTRAINT fk_wishlist_items_product FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: wishlist_items fk_wishlist_items_user; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.wishlist_items
    ADD CONSTRAINT fk_wishlist_items_user FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: tickets tickets_subject_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.tickets
    ADD CONSTRAINT tickets_subject_id_fkey FOREIGN KEY (subject_id) REFERENCES public.ticket_subjects(id);


--
-- Name: vendor_banners vendor_banners_vendor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_banners
    ADD CONSTRAINT vendor_banners_vendor_id_fkey FOREIGN KEY (vendor_id) REFERENCES public.vendors(id) ON DELETE CASCADE;


--
-- Name: vendor_couriers vendor_couriers_courier_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_couriers
    ADD CONSTRAINT vendor_couriers_courier_id_fkey FOREIGN KEY (courier_id) REFERENCES public.couriers(id) ON DELETE CASCADE;


--
-- Name: vendor_couriers vendor_couriers_vendor_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.vendor_couriers
    ADD CONSTRAINT vendor_couriers_vendor_id_fkey FOREIGN KEY (vendor_id) REFERENCES public.vendors(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

\unrestrict bIF4ooiSRpwcf5IUMzkjoAvUzVDvYTG2iw8yd3ryvByArb8IkxhtYzXuRbl5Bog

