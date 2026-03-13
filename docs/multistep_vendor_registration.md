# Vendor Registration Multi-Step

Dokumen ini menjelaskan alur registrasi vendor multi-step yang dipakai backend saat ini. Isi dokumen ini disusun dari implementasi endpoint backend dan referensi flow di `scripts/manual-tests/vendor-register-multistep.html`.

Base path API:

```text
/api/v1
```

## Ringkasan Flow

| Step | Endpoint | Auth | Tujuan |
| --- | --- | --- | --- |
| 1 | `POST /vendors/register/request-otp` | Tidak perlu | Mengirim OTP ke email vendor |
| 2 | `POST /vendors/register/verify-otp` | Tidak perlu | Verifikasi OTP dan mendapatkan `onboarding_token` |
| 3 | `POST /vendors/register/password` | `Bearer onboarding_token` | Menyimpan password akun vendor |
| 4 | `POST /vendors/register/store` | `Bearer onboarding_token` | Menyimpan data toko |
| 5a | `POST /vendors/register/legal-document/presign` | `Bearer onboarding_token` | Mendapat `upload_url` dan `object_key` untuk upload dokumen |
| 5b | `PUT <upload_url>` | Tidak perlu | Upload file langsung ke storage |
| 5c | `POST /vendors/register/legal-document/submit` | `Bearer onboarding_token` | Finalisasi registrasi dan mendapatkan `vendor auth token` |

## Format Response API

Semua endpoint API backend menggunakan envelope JSON berikut:

```json
{
  "success": true,
  "message": "vendor registration otp sent",
  "data": {},
  "errors": null,
  "meta": null
}
```

Catatan:

- Payload utama selalu ada di `data`.
- Saat gagal, frontend sebaiknya baca `message`.
- Untuk error validasi dari Gin binder, `message` biasanya `validation failed` dan detail validator ada di `errors`.
- Upload ke `presigned URL` bukan endpoint API backend, jadi response-nya tidak memakai envelope ini.

## Step 1. Request OTP

Endpoint:

```http
POST /api/v1/vendors/register/request-otp
Content-Type: application/json
```

Request body:

```json
{
  "email": "vendor@example.com"
}
```

Success response:

```json
{
  "success": true,
  "message": "vendor registration otp sent",
  "data": {
    "expires_at": "2026-03-12T10:00:00Z",
    "cooldown_until": "2026-03-12T09:31:00Z"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Mulai flow registrasi dengan mengirim OTP ke email.
- Frontend perlu simpan `email`, `expires_at`, dan `cooldown_until`.
- `cooldown_until` dipakai untuk disable tombol resend/request OTP.

Status yang perlu ditangani:

- `201` sukses.
- `400` payload tidak valid.
- `409` email sudah terdaftar, atau email tersebut sudah punya vendor.
- `429` request OTP terlalu cepat.

## Step 2. Verify OTP

Endpoint:

```http
POST /api/v1/vendors/register/verify-otp
Content-Type: application/json
```

Request body:

```json
{
  "email": "vendor@example.com",
  "code": "123456"
}
```

Success response:

```json
{
  "success": true,
  "message": "vendor registration otp verified",
  "data": {
    "onboarding_token": "<jwt>",
    "onboarding_expires_at": "2026-03-12T17:30:00Z",
    "status": "otp_verified"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- OTP diverifikasi dan backend menerbitkan `onboarding_token`.
- Token ini wajib dipakai pada step 3 sampai step 5c melalui header:

```http
Authorization: Bearer <onboarding_token>
```

- Frontend perlu simpan minimal:
  - `onboarding_token`
  - `onboarding_expires_at`
  - `status`
- Jika user verify OTP lagi untuk email yang sama sebelum registrasi selesai, backend akan me-reset draft onboarding lama dan memulai ulang dari status `otp_verified`.

Status yang perlu ditangani:

- `200` sukses.
- `400` OTP salah, OTP expired, atau payload tidak valid.
- `409` email sudah dipakai akun/vendor.
- `429` terlalu banyak percobaan OTP yang salah.

## Step 3. Save Password

Endpoint:

```http
POST /api/v1/vendors/register/password
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "password": "Password123!"
}
```

Success response:

```json
{
  "success": true,
  "message": "vendor registration password saved",
  "data": {
    "status": "password_set"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Hanya bisa dipanggil jika status onboarding masih `otp_verified`.
- Setelah sukses, frontend update state lokal ke `password_set`.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau step tidak sesuai.
- `401` onboarding token tidak valid / expired / format header salah.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 4. Save Store Info

Endpoint:

```http
POST /api/v1/vendors/register/store
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "store_name": "Toko Oleh Oleh Haji Test",
  "vendor_type": "souvenir_store"
}
```

`vendor_type` yang diizinkan:

- `souvenir_store`
- `ppiu`
- `hajj_dormitory`

Success response:

```json
{
  "success": true,
  "message": "vendor registration store saved",
  "data": {
    "status": "store_info_completed"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Hanya bisa dipanggil jika status onboarding sudah `password_set`.
- Setelah sukses, frontend update state lokal ke `store_info_completed`.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau step tidak sesuai.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 5a. Presign Legal Document

Endpoint:

```http
POST /api/v1/vendors/register/legal-document/presign
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "document_id_type": "ktp"
}
```

`document_id_type` yang diizinkan:

- `ktp`
- `passport`

Success response:

```json
{
  "success": true,
  "message": "vendor registration legal document presigned",
  "data": {
    "upload_url": "https://...",
    "object_key": "vendor-onboardings/<onboarding_id>/documents/owner_document_id/<file_id>",
    "expires_at": "2026-03-12T10:00:00Z",
    "document_id_type": "ktp",
    "doc_type": "owner_document_id"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Hanya bisa dipanggil jika status onboarding sudah `store_info_completed`.
- Frontend wajib simpan `upload_url` dan `object_key`.
- `object_key` ini harus dikirim lagi apa adanya pada step 5c.
- Endpoint ini belum menyelesaikan registrasi dan belum mengubah status onboarding.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau step tidak sesuai.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 5b. Upload File ke Presigned URL

Endpoint:

```http
PUT <upload_url>
Content-Type: <mime-type-file>
```

Behavior:

- File di-upload langsung ke storage menggunakan `upload_url` dari step 5a.
- Tidak memakai header `Authorization` backend.
- Gunakan `Content-Type` sesuai file asli.
- MIME type yang diterima backend saat submit:
  - `image/jpeg`
  - `image/png`
  - `image/webp`
  - `application/pdf`
- Setelah upload sukses, frontend cukup tandai file sudah ter-upload dan simpan `object_key`.
- Response body dari storage bisa kosong. Jangan mengandalkan format JSON.

Catatan penting:

- Jangan ubah `object_key`.
- Jika file belum benar-benar ter-upload ke storage, step 5c akan gagal.

## Step 5c. Submit Legal Info dan Selesaikan Registrasi

Endpoint:

```http
POST /api/v1/vendors/register/legal-document/submit
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "business_legal_type": "perorangan",
  "document_id_type": "ktp",
  "nik": "3173000000000001",
  "owner_name": "Ahmad Subarkah",
  "birth_date": "1990-01-02",
  "document_id_object_key": "vendor-onboardings/<onboarding_id>/documents/owner_document_id/<file_id>"
}
```

Value yang diizinkan:

- `business_legal_type`: `perorangan`, `korporasi`
- `document_id_type`: `ktp`, `passport`
- `birth_date`: format `YYYY-MM-DD`

Success response:

```json
{
  "success": true,
  "message": "vendor registration completed",
  "data": {
    "token": "<vendor_auth_jwt>",
    "vendor_id": "<uuid>",
    "image_url": "",
    "email": "vendor@example.com",
    "name": "Ahmad Subarkah",
    "vendor_type": "souvenir_store",
    "vendor_status": "draft",
    "display_name": "Toko Oleh Oleh Haji Test"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Ini step final dari flow multi-step.
- Backend akan:
  - menyimpan legal info ke draft onboarding,
  - membuat user baru dengan role `umkm`,
  - membuat vendor baru,
  - menyimpan dokumen `owner_document_id`,
  - menginisialisasi balance vendor,
  - menandai onboarding sebagai `completed`,
  - mengembalikan JWT auth vendor biasa.
- Setelah sukses, frontend harus berhenti memakai `onboarding_token` dan ganti ke `data.token`.

Catatan penting:

- `vendor_status` awal setelah registrasi adalah `draft`, bukan `active`.
- Artinya registrasi akun vendor sudah selesai, tetapi vendor belum fully approved untuk semua fitur seller.
- Jika frontend punya seller dashboard, sebaiknya arahkan user ke state "menunggu verifikasi / lengkapi onboarding lanjutan", bukan langsung ke flow vendor aktif penuh.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid, step tidak sesuai, format tanggal salah, `object_key` tidak cocok, file belum ada di storage, atau MIME type file tidak valid.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai atau email sudah dipakai.
- `500` unexpected server error.

## Endpoint Terkait Setelah Registrasi Selesai

Endpoint ini bukan bagian dari step registrasi awal, tapi biasanya dipakai sesudah step 5c:

### Login Vendor

```http
POST /api/v1/vendors/login
Content-Type: application/json
```

Request:

```json
{
  "email": "vendor@example.com",
  "password": "Password123!"
}
```

Behavior:

- Dipakai untuk login di sesi berikutnya.
- Vendor masih bisa login walaupun status-nya `draft`, `submitted`, atau `rejected`.
- Vendor tidak bisa login jika status `blocked`.

### Ambil Profil Vendor

```http
GET /api/v1/vendors/me
Authorization: Bearer <vendor_auth_jwt>
```

Behavior:

- Dipakai untuk hydrate session setelah frontend menyimpan `token` dari step 5c atau dari login biasa.

## State yang Sebaiknya Disimpan di Frontend

Minimal state lokal yang perlu dipersist:

- `email`
- `otp_expires_at`
- `cooldown_until`
- `onboarding_token`
- `onboarding_expires_at`
- `onboarding_status`
- `document_object_key`
- `document_uploaded`
- `vendor_auth_token`
- `vendor_id`
- `vendor_status`

Catatan:

- Saat ini tidak ada endpoint `GET` untuk membaca progress onboarding draft.
- Jika frontend kehilangan `onboarding_token`, cara paling aman adalah ulangi step verify OTP.
- Karena verify OTP me-reset draft onboarding untuk email yang sama, jangan lakukan auto-verify ulang tanpa konfirmasi user jika mereka sudah mengisi step lanjutan.

## Rekomendasi Implementasi Frontend

- Perlakukan flow ini sebagai wizard yang strict per-step.
- Simpan `onboarding_token` setelah step 2 dan kirim di semua request step 3-5c.
- Tampilkan countdown dari `cooldown_until` dan `onboarding_expires_at`.
- Pisahkan step `presign`, `upload`, dan `submit`; jangan gabungkan menjadi satu request.
- Simpan `object_key` dari response presign, bukan dari nama file lokal.
- Setelah step 5c sukses, hapus semua state onboarding lokal dan simpan `vendor_auth_token`.
- Jangan asumsikan vendor langsung `active` setelah registrasi selesai.
