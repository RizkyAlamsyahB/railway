# Vendor Registration Multi-Step

Dokumen ini menjelaskan alur registrasi vendor multi-step yang dipakai backend saat ini. Isi dokumen ini disusun dari implementasi endpoint backend.

Flow vendor terdiri dari 2 fase:

- Fase 1: onboarding awal sampai akun vendor berhasil dibuat dan vendor bisa login dengan status `draft`
- Fase 2: post-login completion untuk melengkapi rekening bank dan dokumen bisnis sampai status vendor berubah menjadi `submitted`

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
| 5a-individual | `POST /vendors/register/souvenir-store/individual/presign` | `Bearer onboarding_token` | Mendapat `upload_url` dan `object_key` untuk dokumen identitas owner |
| 5b-individual | `PUT <upload_url>` | Tidak perlu | Upload file identitas owner langsung ke storage |
| 5c-individual | `POST /vendors/register/souvenir-store/individual` | `Bearer onboarding_token` | Finalisasi registrasi awal vendor individual dan mendapatkan `vendor auth token` |
| 5a-corporate | `POST /vendors/register/souvenir-store/corporate/presign` | `Bearer onboarding_token` | Mendapat `upload_url` dan `object_key` untuk dokumen NIB |
| 5b-corporate | `PUT <upload_url>` | Tidak perlu | Upload file NIB langsung ke storage |
| 5c-corporate | `POST /vendors/register/souvenir-store/corporate` | `Bearer onboarding_token` | Finalisasi registrasi awal vendor corporate dan mendapatkan `vendor auth token` |
| 6a | `POST /vendors/bank-account` | `Bearer vendor_auth_token` | Menyimpan atau memperbarui rekening bank vendor |
| 6b | `POST /vendors/documents/presign` | `Bearer vendor_auth_token` | Mendapat `upload_url` dan `object_key` untuk dokumen completion |
| 6c | `PUT <upload_url>` | Tidak perlu | Upload file completion langsung ke storage |
| 6d | `POST /vendors/documents/confirm` | `Bearer vendor_auth_token` | Konfirmasi dokumen completion dan ubah status ke `submitted` bila lengkap |

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
- Setelah step ini, flow bercabang ke jalur `individual` atau `corporate`.

Catatan implementasi saat ini:

- Endpoint finalisasi onboarding ada di path `/vendors/register/souvenir-store/...`.
- Backend saat ini belum membedakan path finalisasi berdasarkan nilai `vendor_type` yang disimpan pada step ini.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau step tidak sesuai.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 5A. Jalur Individual: Presign Dokumen Identitas

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

- Hanya bisa dipanggil jika status onboarding sudah `store_info_completed`.
- Jalur ini dipakai untuk registrasi vendor dengan legal type `perorangan`.
- Frontend harus menentukan MIME file sebelum meminta presign URL.
- Frontend wajib simpan `upload_url` dan `object_key`.
- `object_key` ini harus dikirim lagi apa adanya pada step 5C individual.
- Endpoint ini belum menyelesaikan registrasi dan belum mengubah status onboarding.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau step tidak sesuai.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 5B. Jalur Individual: Upload File ke Presigned URL

Endpoint:

```http
PUT <upload_url>
Content-Type: <mime-type-file>
```

Behavior:

- File di-upload langsung ke storage menggunakan `upload_url` dari step 5A individual.
- Tidak memakai header `Authorization` backend.
- Gunakan `Content-Type` sesuai file asli.
- MIME type yang diterima backend saat submit individual:
  - `image/jpeg`
  - `application/pdf`
- Setelah upload sukses, frontend cukup tandai file sudah ter-upload dan simpan `object_key`.
- Response body dari storage bisa kosong. Jangan mengandalkan format JSON.

Catatan penting:

- Jangan ubah `object_key`.
- Jika file belum benar-benar ter-upload ke storage, step 5C individual akan gagal.

## Step 5C. Jalur Individual: Submit Legal Info dan Selesaikan Registrasi

Endpoint:

```http
POST /api/v1/vendors/register/souvenir-store/individual
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
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

- Ini step final dari flow onboarding awal untuk jalur individual.
- Backend akan otomatis menetapkan `business_legal_type = perorangan`.
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

## Step 5D. Jalur Corporate: Presign Dokumen NIB

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

- Hanya bisa dipanggil jika status onboarding sudah `store_info_completed`.
- Jalur ini dipakai untuk registrasi vendor dengan legal type `korporasi`.
- Frontend wajib simpan `upload_url` dan `object_key`.
- `object_key` ini harus dikirim lagi apa adanya pada step 5F corporate.
- Endpoint ini belum menyelesaikan registrasi dan belum mengubah status onboarding.

Status yang perlu ditangani:

- `200` sukses.
- `400` step tidak sesuai.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai.

## Step 5E. Jalur Corporate: Upload File NIB ke Presigned URL

Endpoint:

```http
PUT <upload_url>
Content-Type: application/pdf
```

Behavior:

- File di-upload langsung ke storage menggunakan `upload_url` dari step 5D corporate.
- Tidak memakai header `Authorization` backend.
- Untuk jalur corporate, file NIB harus berupa PDF.
- Setelah upload sukses, frontend cukup tandai file sudah ter-upload dan simpan `object_key`.
- Response body dari storage bisa kosong. Jangan mengandalkan format JSON.

Catatan penting:

- Jangan ubah `object_key`.
- Jika file belum benar-benar ter-upload ke storage, step 5F corporate akan gagal.
- Jika MIME type file bukan `application/pdf`, submit corporate akan gagal.

## Step 5F. Jalur Corporate: Submit Legal Info dan Selesaikan Registrasi

Endpoint:

```http
POST /api/v1/vendors/register/souvenir-store/corporate
Authorization: Bearer <onboarding_token>
Content-Type: application/json
```

Request body:

```json
{
  "nib": "1234567890123",
  "company_name": "PT Oleh Oleh Haji Nusantara",
  "established_date": "2020-01-02",
  "registered_address": "Jl. Contoh No. 1, Jakarta",
  "nib_document_object_key": "vendor-onboardings/<onboarding_id>/documents/business_nib/<file_id>"
}
```

Value yang perlu diperhatikan:

- `established_date`: format `YYYY-MM-DD`

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
    "vendor_status": "draft",
    "display_name": "Toko Oleh Oleh Haji Test"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Ini step final dari flow onboarding awal untuk jalur corporate.
- Backend akan otomatis menetapkan `business_legal_type = korporasi`.
- Backend akan:
  - menyimpan legal info corporate ke draft onboarding,
  - membuat user baru dengan role `umkm`,
  - membuat vendor baru,
  - menyimpan dokumen `business_nib`,
  - menginisialisasi balance vendor,
  - menandai onboarding sebagai `completed`,
  - mengembalikan JWT auth vendor biasa.
- Setelah sukses, frontend harus berhenti memakai `onboarding_token` dan ganti ke `data.token`.

Catatan penting:

- `vendor_status` awal setelah registrasi adalah `draft`, bukan `active`.
- Artinya registrasi akun vendor sudah selesai, tetapi vendor belum fully approved untuk semua fitur seller.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid, step tidak sesuai, format tanggal salah, `object_key` tidak cocok, file belum ada di storage, atau MIME type file tidak valid.
- `401` onboarding token tidak valid / expired.
- `404` onboarding draft tidak ditemukan.
- `409` onboarding sudah selesai atau email sudah dipakai.
- `500` unexpected server error.

## Fase 2. Post-Login Completion Sampai `submitted`

Setelah step 5C individual atau step 5F corporate sukses, vendor sudah punya akun, bisa login, dan `vendor_status` masih `draft`.

Vendor belum bisa lanjut ke alur seller aktif penuh sebelum fase completion ini selesai. Backend baru akan mengubah status menjadi `submitted` jika seluruh syarat berikut sudah terpenuhi:

- rekening bank vendor sudah tersimpan
- dokumen legal onboarding awal sudah ada:
  - `owner_document_id` untuk `perorangan`
  - `business_nib` untuk `korporasi`
- dokumen completion wajib sudah dikonfirmasi:
  - `store_photo`
  - `bank_account_proof`
  - `business_logo`
  - `business_banner`
- jika `business_legal_type = korporasi`, dokumen `business_npwp` juga wajib

Catatan:

- step completion hanya bisa dijalankan saat status vendor masih `draft`
- dokumen completion dibuat per dokumen lewat flow `presign -> upload -> confirm`
- backend akan mengevaluasi ulang kelengkapan completion saat simpan rekening bank maupun saat endpoint `confirm` dipanggil
- endpoint `confirm` bisa dipanggil berulang kali; status tetap `draft` sampai semua syarat lengkap

### Step 6A. Simpan Rekening Bank Vendor

Endpoint:

```http
POST /api/v1/vendors/bank-account
Authorization: Bearer <vendor_auth_jwt>
Content-Type: application/json
```

Request body:

```json
{
  "bank_name": "Bank Syariah Indonesia",
  "account_number": "1234567890",
  "account_holder_name": "Ahmad Subarkah"
}
```

Success response:

```json
{
  "success": true,
  "message": "vendor bank account saved successfully",
  "data": {
    "vendor_id": "<uuid>",
    "vendor_status": "submitted",
    "bank_account": {
      "id": "<uuid>",
      "vendor_id": "<uuid>",
      "bank_name": "Bank Syariah Indonesia",
      "account_number": "1234567890",
      "account_holder_name": "Ahmad Subarkah",
      "verification_status": "pending",
      "rejection_reason": null,
      "verified_by": null,
      "verified_at": null,
      "created_at": "2026-03-13T08:00:00Z",
      "updated_at": "2026-03-13T08:00:00Z"
    }
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Hanya bisa dipanggil saat vendor berstatus `draft`.
- Jika rekening bank sudah pernah ada, backend akan update data dan reset status verifikasi ke `pending`.
- Setelah rekening berhasil disimpan, backend langsung mengevaluasi ulang apakah seluruh syarat completion sudah lengkap.
- Jika seluruh syarat completion sudah lengkap pada saat save rekening, status vendor otomatis berubah menjadi `submitted`.
- Frontend cukup menyimpan `vendor_status`; detail rekening bisa dianggap source of truth dari response terakhir.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid atau status vendor bukan `draft`.
- `401` token vendor tidak valid / expired.
- `403` role bukan `umkm`.
- `404` vendor tidak ditemukan.

### Step 6B. Presign Dokumen Completion

Endpoint:

```http
POST /api/v1/vendors/documents/presign
Authorization: Bearer <vendor_auth_jwt>
Content-Type: application/json
```

Request body:

```json
{
  "doc_type": "store_photo"
}
```

`doc_type` yang diizinkan:

- `store_photo`
- `bank_account_proof`
- `business_logo`
- `business_banner`
- `business_npwp`

Success response:

```json
{
  "success": true,
  "message": "vendor document presigned successfully",
  "data": {
    "doc_type": "store_photo",
    "upload_url": "https://...",
    "object_key": "vendors/<vendor_id>/documents/store_photo/<file_id>",
    "expires_at": "2026-03-13T08:30:00Z"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Hanya bisa dipanggil saat vendor berstatus `draft`.
- dokumen legal onboarding awal (`owner_document_id` untuk `perorangan`, `business_nib` untuk `korporasi`) tidak boleh dipresign dari endpoint ini karena sudah dibuat pada onboarding awal.
- Frontend wajib menyimpan `object_key` untuk dikirim ke step 6D.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid, `doc_type` tidak didukung, atau status vendor bukan `draft`.
- `401` token vendor tidak valid / expired.
- `403` role bukan `umkm`.
- `404` vendor tidak ditemukan.

### Step 6C. Upload File Completion ke Presigned URL

Endpoint:

```http
PUT <upload_url>
Content-Type: <mime-type-file>
```

Behavior:

- File di-upload langsung ke storage menggunakan `upload_url` dari step 6B.
- Tidak memakai header `Authorization` backend.
- MIME type yang diterima backend saat confirm:
  - `image/jpeg`
  - `image/png`
  - `image/webp`
  - `application/pdf`
- Setelah upload sukses, frontend cukup tandai file completion sudah ter-upload dan simpan `object_key`.

### Step 6D. Confirm Dokumen Completion

Endpoint:

```http
POST /api/v1/vendors/documents/confirm
Authorization: Bearer <vendor_auth_jwt>
Content-Type: application/json
```

Request body:

```json
{
  "documents": [
    {
      "doc_type": "store_photo",
      "object_key": "vendors/<vendor_id>/documents/store_photo/<file_id>"
    }
  ]
}
```

Success response saat belum lengkap:

```json
{
  "success": true,
  "message": "documents confirmed successfully",
  "data": {
    "vendor_id": "<uuid>",
    "vendor_status": "draft",
    "documents_confirmed": 1
  },
  "errors": null,
  "meta": null
}
```

Success response saat semua syarat sudah lengkap:

```json
{
  "success": true,
  "message": "documents confirmed successfully",
  "data": {
    "vendor_id": "<uuid>",
    "vendor_status": "submitted",
    "documents_confirmed": 1
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Endpoint ini memverifikasi object di storage dan meng-update metadata dokumen.
- `doc_type` yang diterima hanya dokumen completion:
  - `store_photo`
  - `bank_account_proof`
  - `business_logo`
  - `business_banner`
  - `business_npwp`
- Jika semua syarat sudah lengkap, backend otomatis mengubah status vendor dari `draft` ke `submitted`.
- Jika belum lengkap, request tetap sukses tetapi `vendor_status` tetap `draft`.
- Endpoint ini bisa dipakai untuk memicu reevaluasi ulang walaupun dokumen yang dikirim sudah pernah dikonfirmasi sebelumnya.

Status yang perlu ditangani:

- `200` sukses.
- `400` payload tidak valid, `doc_type` tidak didukung, file belum ada di storage, MIME type tidak valid, ukuran file terlalu besar, atau status vendor bukan `draft`.
- `401` token vendor tidak valid / expired.
- `403` role bukan `umkm`.
- `404` vendor tidak ditemukan.
- `500` unexpected server error.

## Endpoint Terkait Setelah Registrasi Awal Selesai

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

- Dipakai untuk hydrate session setelah frontend menyimpan `token` dari step 5C individual, step 5F corporate, atau dari login biasa.
- Frontend juga bisa memakai endpoint ini untuk menentukan apakah vendor masih `draft`, sudah `submitted`, atau sudah `active`.

## State yang Sebaiknya Disimpan di Frontend

Minimal state lokal yang perlu dipersist:

- `email`
- `otp_expires_at`
- `cooldown_until`
- `onboarding_token`
- `onboarding_expires_at`
- `onboarding_status`
- `registration_path` (`individual` atau `corporate`) bila wizard perlu bisa di-resume
- `individual_document_object_key` untuk jalur individual
- `individual_document_uploaded` untuk jalur individual
- `corporate_nib_object_key` untuk jalur corporate
- `corporate_nib_uploaded` untuk jalur corporate
- `vendor_auth_token`
- `vendor_id`
- `vendor_status`
- `completion_document_object_key` per dokumen yang sedang diproses
- `completion_document_uploaded` per dokumen yang sedang diproses
- data rekening bank draft jika flow completion dipecah ke beberapa layar

Catatan:

- Saat ini tidak ada endpoint `GET` untuk membaca progress onboarding draft.
- Jika frontend kehilangan `onboarding_token`, cara paling aman adalah ulangi step verify OTP.
- Karena verify OTP me-reset draft onboarding untuk email yang sama, jangan lakukan auto-verify ulang tanpa konfirmasi user jika mereka sudah mengisi step lanjutan.

## Rekomendasi Implementasi Frontend

- Perlakukan flow ini sebagai wizard yang strict per-step.
- Simpan `onboarding_token` setelah step 2 dan kirim di semua request step 3 sampai finalisasi jalur individual/corporate.
- Tampilkan countdown dari `cooldown_until` dan `onboarding_expires_at`.
- Setelah step 4, arahkan user ke salah satu jalur:
  - jalur individual: `presign -> upload -> submit` via endpoint `/souvenir-store/individual`
  - jalur corporate: `presign -> upload -> submit` via endpoint `/souvenir-store/corporate`
- Simpan `object_key` dari response presign, bukan dari nama file lokal.
- Untuk jalur corporate, upload dokumen NIB sebagai PDF.
- Setelah finalisasi registrasi awal sukses, hentikan pemakaian `onboarding_token` dan ganti ke `vendor_auth_token`.
- Lanjutkan flow completion setelah login pertama: simpan rekening bank, lalu upload dan confirm dokumen completion satu per satu.
- Jangan asumsikan vendor langsung `active` setelah registrasi awal selesai; status awalnya tetap `draft`.
- Anggap `submitted` sebagai state "siap direview admin", bukan "vendor aktif penuh".
