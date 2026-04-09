-- Import seller accounts from:
--   Data Haji Umrah Store - Sheet1.csv
--
-- How to run from repository root:
--   psql "$DATABASE_URL" -f scripts/import_umkm_sellers_from_csv.sql
--
-- What this script does:
-- 1. Reads the CSV with psql \copy.
-- 2. Keeps the first row per email and skips invalid email rows.
-- 3. Creates or updates users so they can log in:
--    - role = umkm
--    - status = active
--    - email_verified_at is set
--    - password_hash uses bcrypt via pgcrypto
-- 4. Creates a draft vendor profile for each imported seller if it does not exist yet.
--
-- Notes:
-- - The current CSV yields 208 unique valid sellers after deduplication.
-- - Duplicate or suspicious phone numbers are set to NULL to avoid the users.phone UNIQUE constraint.
-- - Re-running this script will reset imported users' password hash to the password value from the CSV.

BEGIN;

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TEMP TABLE tmp_umkm_import_raw (
    line_no BIGSERIAL PRIMARY KEY,
    source_no TEXT,
    nama_umkm TEXT,
    jenis_produk TEXT,
    keterangan TEXT,
    harga TEXT,
    gambar_produk TEXT,
    link_gambar_produk TEXT,
    nomor_kontak TEXT,
    nama_pemilik TEXT,
    alamat_usaha TEXT,
    provinsi TEXT,
    unused_col TEXT,
    email TEXT,
    raw_password TEXT,
    pic TEXT
) ON COMMIT DROP;

\copy tmp_umkm_import_raw (
    source_no,
    nama_umkm,
    jenis_produk,
    keterangan,
    harga,
    gambar_produk,
    link_gambar_produk,
    nomor_kontak,
    nama_pemilik,
    alamat_usaha,
    provinsi,
    unused_col,
    email,
    raw_password,
    pic
) FROM 'Data Haji Umrah Store - Sheet1.csv' WITH (FORMAT csv, HEADER true, ENCODING 'UTF8');

CREATE TEMP TABLE tmp_umkm_import_clean ON COMMIT DROP AS
WITH normalized AS (
    SELECT
        r.line_no,
        lower(btrim(r.email)) AS email,
        NULLIF(btrim(r.nama_umkm), '') AS nama_umkm,
        NULLIF(btrim(r.nama_pemilik), '') AS nama_pemilik,
        NULLIF(btrim(r.jenis_produk), '') AS jenis_produk,
        NULLIF(btrim(r.keterangan), '') AS keterangan,
        NULLIF(btrim(r.alamat_usaha), '') AS alamat_usaha,
        NULLIF(btrim(r.provinsi), '') AS provinsi,
        NULLIF(btrim(r.raw_password), '') AS raw_password,
        CASE
            WHEN COALESCE(r.nomor_kontak, '') ~ '[eE][+-]?[0-9]+' THEN NULL
            ELSE NULLIF(regexp_replace(COALESCE(r.nomor_kontak, ''), '\D', '', 'g'), '')
        END AS phone_digits
    FROM tmp_umkm_import_raw r
),
filtered AS (
    SELECT
        n.*,
        COALESCE(n.nama_umkm, n.nama_pemilik) AS user_full_name,
        COALESCE(n.nama_umkm, n.nama_pemilik) AS store_name,
        COALESCE(n.nama_pemilik, n.nama_umkm) AS responsible_person_name
    FROM normalized n
    WHERE n.email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'
      AND COALESCE(n.nama_umkm, n.nama_pemilik) IS NOT NULL
      AND n.raw_password IS NOT NULL
),
deduped AS (
    SELECT DISTINCT ON (f.email)
        f.*
    FROM filtered f
    ORDER BY f.email, f.line_no
),
phone_ranked AS (
    SELECT
        d.*,
        CASE
            WHEN d.phone_digits IS NULL THEN NULL
            WHEN length(d.phone_digits) < 8 THEN NULL
            WHEN count(*) OVER (PARTITION BY d.phone_digits) > 1 THEN NULL
            ELSE left(d.phone_digits, 20)
        END AS phone_for_user
    FROM deduped d
)
SELECT
    p.line_no,
    p.email,
    left(p.user_full_name, 120) AS user_full_name,
    left(p.store_name, 120) AS store_name,
    left(p.responsible_person_name, 120) AS responsible_person_name,
    left(COALESCE(p.responsible_person_name, p.store_name), 160) AS legal_name,
    p.phone_for_user,
    p.raw_password,
    p.jenis_produk,
    p.keterangan,
    p.alamat_usaha,
    p.provinsi,
    NULLIF(
        concat_ws(
            E'\n',
            CASE WHEN p.jenis_produk IS NOT NULL THEN 'Produk: ' || p.jenis_produk END,
            CASE WHEN p.keterangan IS NOT NULL THEN 'Catatan: ' || p.keterangan END,
            CASE WHEN p.alamat_usaha IS NOT NULL THEN 'Alamat: ' || p.alamat_usaha END,
            CASE WHEN p.provinsi IS NOT NULL THEN 'Provinsi: ' || p.provinsi END
        ),
        ''
    ) AS vendor_description
FROM phone_ranked p;

CREATE TEMP TABLE tmp_umkm_import_final ON COMMIT DROP AS
SELECT
    c.*,
    CASE
        WHEN c.phone_for_user IS NULL THEN NULL
        WHEN EXISTS (
            SELECT 1
            FROM users u
            WHERE u.phone = c.phone_for_user
              AND lower(u.email) <> c.email
        ) THEN NULL
        ELSE c.phone_for_user
    END AS phone_final
FROM tmp_umkm_import_clean c;

WITH role_umkm AS (
    SELECT id
    FROM roles
    WHERE code = 'umkm'
)
INSERT INTO users (
    id,
    email,
    full_name,
    phone,
    password_hash,
    role_id,
    status,
    email_verified_at,
    created_at,
    updated_at
)
SELECT
    gen_random_uuid(),
    s.email,
    s.user_full_name,
    s.phone_final,
    crypt(s.raw_password, gen_salt('bf')),
    r.id,
    'active',
    now(),
    now(),
    now()
FROM tmp_umkm_import_final s
CROSS JOIN role_umkm r
ON CONFLICT ((lower(email))) DO UPDATE
SET
    full_name = EXCLUDED.full_name,
    phone = EXCLUDED.phone,
    password_hash = EXCLUDED.password_hash,
    role_id = EXCLUDED.role_id,
    status = 'active',
    email_verified_at = COALESCE(users.email_verified_at, now()),
    updated_at = now();

INSERT INTO vendors (
    id,
    owner_user_id,
    vendor_type,
    legal_name,
    display_name,
    responsible_person_name,
    description,
    status,
    created_at,
    updated_at
)
SELECT
    gen_random_uuid(),
    u.id,
    'general_souvenir_store',
    s.legal_name,
    s.store_name,
    s.responsible_person_name,
    s.vendor_description,
    'draft',
    now(),
    now()
FROM tmp_umkm_import_final s
JOIN users u
    ON lower(u.email) = s.email
WHERE NOT EXISTS (
    SELECT 1
    FROM vendors v
    WHERE v.owner_user_id = u.id
);

SELECT
    COUNT(*) AS prepared_rows
FROM tmp_umkm_import_final;

SELECT
    COUNT(*) AS users_ready_to_login
FROM users u
JOIN roles r
    ON r.id = u.role_id
JOIN tmp_umkm_import_final s
    ON lower(u.email) = s.email
WHERE r.code = 'umkm'
  AND u.status = 'active'
  AND u.email_verified_at IS NOT NULL;

SELECT
    COUNT(*) AS vendor_profiles_available
FROM vendors v
JOIN users u
    ON u.id = v.owner_user_id
JOIN tmp_umkm_import_final s
    ON lower(u.email) = s.email;

COMMIT;
