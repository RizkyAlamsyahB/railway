# Vendor Registration Multi-Step

Dokumen ini menjelaskan alur registrasi vendor publik yang berlaku saat ini.

Flow vendor terdiri dari 1 fase:

- Fase 1: onboarding awal sampai akun vendor berhasil dibuat dengan status `draft`

Base path API:

```text
/api/v1
```

## Ringkasan Flow

| Step | Endpoint | Auth | Tujuan |
| --- | --- | --- | --- |
| 1 | `POST /vendors/register/request-otp` | Tidak perlu | Mengirim OTP ke email vendor |
| 2 | `POST /vendors/register/verify-otp` | Tidak perlu | Verifikasi OTP dan mendapatkan `onboarding_token` |
| 3 | `POST /vendors/register/password` | `Bearer onboarding_token` | Finalisasi registrasi publik, membuat akun vendor, dan mengembalikan auth token |

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

Catatan:

- Endpoint ini memulai flow registrasi dengan mengirim OTP ke email vendor.
- Frontend perlu menyimpan `email`, `expires_at`, dan `cooldown_until`.

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

Catatan:

- `onboarding_token` wajib dikirim pada step 3 melalui header `Authorization: Bearer <onboarding_token>`.
- Jika user verify OTP lagi untuk email yang sama sebelum registrasi selesai, backend akan me-reset onboarding lama ke status `otp_verified`.

## Step 3. Finalisasi Registrasi dan Set Password

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
  "message": "vendor registration completed",
  "data": {
    "access_token": "<jwt>",
    "vendor_id": "6c9f5d13-4c56-4f3d-9b38-2ddf7e0f8d5e",
    "image_url": null,
    "email": "vendor@example.com",
    "store_name": null,
    "vendor_type": null,
    "vendor_status": "draft"
  },
  "errors": null,
  "meta": null
}
```

Behavior:

- Hanya bisa dipanggil jika status onboarding masih `otp_verified`.
- Step ini langsung membuat akun user role `umkm`, profil vendor minimal, dan saldo vendor awal.
- `store_name` bernilai `null` sampai vendor mengisi nama toko pada flow lanjutan.
- `vendor_type` bernilai `null` sampai vendor memilih tipenya pada flow lanjutan.
- `vendor_status` awal adalah `draft`.
- Setelah step ini onboarding ditandai `completed` dan flow registrasi publik selesai.

## Catatan Implementasi

- Tidak ada lagi branch onboarding publik `individual` maupun `corporate`.
- Tidak ada lagi field `business_legal_type` pada flow registrasi publik.
- Proses dokumen dan legalitas vendor dilakukan di flow terpisah setelah vendor login.
