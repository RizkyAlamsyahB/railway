# Vendor Registration Multi-Step

Dokumen ini menjelaskan alur registrasi vendor multi-step yang dipakai backend saat ini. Isi dokumen ini disusun dari implementasi endpoint backend.

Flow vendor terdiri dari 1 fase:

- Fase 1: onboarding awal sampai akun vendor berhasil dibuat dengan status `submitted` dan siap untuk di-review oleh admin

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
| 4a-individual | `POST /vendors/register/souvenir-store/individual/presign` | `Bearer onboarding_token` | Mendapat `upload_url` dan `object_key` untuk dokumen identitas owner |
| 4b-individual | `PUT <upload_url>` | Tidak perlu | Upload file identitas owner langsung ke storage |
| 4c-individual | `POST /vendors/register/souvenir-store/individual` | `Bearer onboarding_token` | Finalisasi registrasi vendor individual — status langsung menjadi `submitted` |
| 4a-corporate | `POST /vendors/register/souvenir-store/corporate/presign` | `Bearer onboarding_token` | Mendapat `upload_url` dan `object_key` untuk dokumen NIB |
| 4b-corporate | `PUT <upload_url>` | Tidak perlu | Upload file NIB langsung ke storage |
| 4c-corporate | `POST /vendors/register/souvenir-store/corporate` | `Bearer onboarding_token` | Finalisasi registrasi vendor corporate — status langsung menjadi `submitted` |

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
- Token ini wajib dipakai pada step 3 sampai step 4c melalui header:

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
- Setelah step ini, flow bercabang ke jalur `individual` atau `corporate`.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau step tidak sesuai.
- `401` onboarding token tidak valid / expired / format header salah.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 4A. Jalur Individual: Presign Dokumen Identitas

Endpoint:

```http
POST /api/v1/vendors/register/souvenir-store/individual/presign
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "document_id_type": "ktp",
  "content_type": "image/jpeg"
}
```

`document_id_type` yang diizinkan:

- `ktp`
- `passport`

`content_type` yang diizinkan:

- `image/jpeg`
- `application/pdf`

Success response:

```json
{
  "success": true,
  "message": "vendor registration individual document presigned",
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

- Hanya bisa dipanggil jika status onboarding sudah `password_set`.
- Jalur ini dipakai untuk registrasi vendor dengan legal type `perorangan`.
- Frontend harus menentukan MIME file sebelum meminta presign URL.
- Frontend wajib simpan `upload_url` dan `object_key`.
- `object_key` ini harus dikirim lagi apa adanya pada step 4C individual.
- Endpoint ini belum menyelesaikan registrasi dan belum mengubah status onboarding.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau step tidak sesuai.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 4B. Jalur Individual: Upload File ke Presigned URL

Endpoint:

```http
PUT <upload_url>
Content-Type: <mime-type-file>
```

Behavior:

- File di-upload langsung ke storage menggunakan `upload_url` dari step 4A individual.
- Tidak memakai header `Authorization` backend.
- Gunakan `Content-Type` sesuai file asli.
- MIME type yang diterima backend saat submit individual:
  - `image/jpeg`
  - `application/pdf`
- Setelah upload sukses, frontend cukup tandai file sudah ter-upload dan simpan `object_key`.
- Response body dari storage bisa kosong. Jangan mengandalkan format JSON.

Catatan penting:

- Jangan ubah `object_key`.
- Jika file belum benar-benar ter-upload ke storage, step 4C individual akan gagal.

## Step 4C. Jalur Individual: Submit Info Toko, Legal Info, dan Selesaikan Registrasi

Endpoint:

```http
POST /api/v1/vendors/register/souvenir-store/individual
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "store_name": "Toko Oleh Oleh Haji Test",
  "document_id_type": "ktp",
  "nik": "3173000000000001",
  "owner_name": "Ahmad Subarkah",
  "birth_date": "1990-01-02",
  "document_id_object_key": "vendor-onboardings/<onboarding_id>/documents/owner_document_id/<file_id>"
}
```

Value yang diizinkan:

- `document_id_type`: `ktp`, `passport`
- `birth_date`: format `YYYY-MM-DD`
- `store_name`: wajib, maksimal 120 karakter

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
    "vendor_status": "submitted",
    "display_name": "Toko Oleh Oleh Haji Test"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Ini step final dari flow onboarding awal untuk jalur individual.
- Step ini sekaligus menyimpan info toko (`store_name`) dan data legal owner.
- Backend akan otomatis menetapkan `business_legal_type = perorangan` dan `vendor_type = souvenir_store`.
- Backend akan:
  - menyimpan info toko dan legal info ke draft onboarding,
  - membuat user baru dengan role `umkm`,
  - membuat vendor baru,
  - menyimpan dokumen `owner_document_id`,
  - menginisialisasi balance vendor,
  - menandai onboarding sebagai `completed`,
  - mengembalikan JWT auth vendor biasa.
- Setelah sukses, frontend harus berhenti memakai `onboarding_token` dan ganti ke `data.token`.

Catatan penting:

- `vendor_status` setelah registrasi langsung menjadi `submitted` dan siap untuk di-review oleh admin.
- Vendor belum bisa menjalankan fitur seller aktif penuh sampai admin mengubah status menjadi `active`.
- Arahkan user ke halaman "menunggu verifikasi admin" setelah registrasi selesai.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid, step tidak sesuai, format tanggal salah, `object_key` tidak cocok, file belum ada di storage, atau MIME type file tidak valid.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai atau email sudah dipakai.
- `500` unexpected server error.

## Step 4D. Jalur Corporate: Presign Dokumen NIB

Endpoint:

```http
POST /api/v1/vendors/register/souvenir-store/corporate/presign
Authorization: Bearer <onboarding_token>
```

Request body:

- Tidak ada request body.

Success response:

```json
{
  "success": true,
  "message": "vendor registration corporate document presigned",
  "data": {
    "upload_url": "https://...",
    "object_key": "vendor-onboardings/<onboarding_id>/documents/business_nib/<file_id>",
    "expires_at": "2026-03-12T10:00:00Z",
    "document_id_type": "",
    "doc_type": "business_nib"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Hanya bisa dipanggil jika status onboarding sudah `password_set`.
- Jalur ini dipakai untuk registrasi vendor dengan legal type `korporasi`.
- Frontend wajib simpan `upload_url` dan `object_key`.
- `object_key` ini harus dikirim lagi apa adanya pada step 4F corporate.
- Endpoint ini belum menyelesaikan registrasi dan belum mengubah status onboarding.

Status yang perlu ditangani:

- `200` sukses.
- `400` step tidak sesuai.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 4E. Jalur Corporate: Upload File NIB ke Presigned URL

Endpoint:

```http
PUT <upload_url>
Content-Type: application/pdf
```

Behavior:

- File di-upload langsung ke storage menggunakan `upload_url` dari step 4D corporate.
- Tidak memakai header `Authorization` backend.
- Untuk jalur corporate, file NIB harus berupa PDF.
- Setelah upload sukses, frontend cukup tandai file sudah ter-upload dan simpan `object_key`.
- Response body dari storage bisa kosong. Jangan mengandalkan format JSON.

Catatan penting:

- Jangan ubah `object_key`.
- Jika file belum benar-benar ter-upload ke storage, step 4F corporate akan gagal.
- Jika MIME type file bukan `application/pdf`, submit corporate akan gagal.

## Step 4F. Jalur Corporate: Submit Info Toko, Legal Info, dan Selesaikan Registrasi

Endpoint:

```http
POST /api/v1/vendors/register/souvenir-store/corporate
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "store_name": "Toko Oleh Oleh Haji Test",
  "nib": "1234567890123",
  "company_name": "PT Oleh Oleh Haji Nusantara",
  "established_date": "2020-01-02",
  "registered_address": "Jl. Contoh No. 1, Jakarta",
  "nib_document_object_key": "vendor-onboardings/<onboarding_id>/documents/business_nib/<file_id>"
}
```

Value yang perlu diperhatikan:

- `established_date`: format `YYYY-MM-DD`
- `store_name`: wajib, maksimal 120 karakter

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
    "name": "PT Oleh Oleh Haji Nusantara",
    "vendor_type": "souvenir_store",
    "vendor_status": "submitted",
    "display_name": "Toko Oleh Oleh Haji Test"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Ini step final dari flow onboarding awal untuk jalur corporate.
- Step ini sekaligus menyimpan info toko (`store_name`) dan data legal corporate.
- Backend akan otomatis menetapkan `business_legal_type = korporasi` dan `vendor_type = souvenir_store`.
- Backend akan:
  - menyimpan info toko dan legal info corporate ke draft onboarding,
  - membuat user baru dengan role `umkm`,
  - membuat vendor baru,
  - menyimpan dokumen `business_nib`,
  - menginisialisasi balance vendor,
  - menandai onboarding sebagai `completed`,
  - mengembalikan JWT auth vendor biasa.
- Setelah sukses, frontend harus berhenti memakai `onboarding_token` dan ganti ke `data.token`.

Catatan penting:

- `vendor_status` setelah registrasi langsung menjadi `submitted` dan siap untuk di-review oleh admin.
- Vendor belum bisa menjalankan fitur seller aktif penuh sampai admin mengubah status menjadi `active`.
- Arahkan user ke halaman "menunggu verifikasi admin" setelah registrasi selesai.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid, step tidak sesuai, format tanggal salah, `object_key` tidak cocok, file belum ada di storage, atau MIME type file tidak valid.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai atau email sudah dipakai.
- `500` unexpected server error.

## Rekomendasi Implementasi Frontend

- Perlakukan flow ini sebagai wizard yang strict per-step.
- Simpan `onboarding_token` setelah step 2 dan kirim di semua request step 3 sampai finalisasi jalur individual/corporate.
- Tampilkan countdown dari `cooldown_until` dan `onboarding_expires_at`.
- Setelah step 3, arahkan user ke salah satu jalur:
  - jalur individual: `presign -> upload -> submit` via endpoint `/souvenir-store/individual`
  - jalur corporate: `presign -> upload -> submit` via endpoint `/souvenir-store/corporate`
- Simpan `object_key` dari response presign, bukan dari nama file lokal.
- Untuk jalur corporate, upload dokumen NIB sebagai PDF.
- Setelah finalisasi registrasi sukses, `vendor_status` langsung menjadi `submitted`.
- Hentikan pemakaian `onboarding_token` dan arahkan user ke halaman "menunggu verifikasi admin".
- Status `submitted` berarti "siap direview admin" — vendor belum bisa menjalankan fitur seller aktif penuh sampai admin mengubah status menjadi `active`.
