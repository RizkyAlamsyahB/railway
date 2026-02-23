# XenPlatform API Reference — Rangkuman Endpoint untuk Marketplace Multi-Seller

> **Tujuan dokumen ini:** Referensi konteks untuk AI Agent dalam pengembangan fitur pembayaran marketplace dengan banyak seller menggunakan Xendit xenPlatform.
>
> **Base URL:** `https://api.xendit.co`  
> **Versi API:** xenPlatform 1.0

---

## Autentikasi

Semua endpoint xenPlatform menggunakan **HTTP Basic Authentication**.

- **Username:** API Key Anda (dari Xendit Dashboard → Settings → API Keys)
- **Password:** dikosongkan (string kosong)

Format header yang dihasilkan:
```
Authorization: Basic {BASE64(API_KEY + ":")}
```

Cara termudah menggunakan curl adalah dengan flag `-u`:
```bash
curl --request GET \
  --url https://api.xendit.co/v2/accounts \
  --header 'accept: application/json' \
  -u "{YOUR_API_KEY}:"
```

Atau menggunakan header eksplisit (seperti pada contoh di Xendit Docs):
```bash
--header 'authorization: Basic {BASE64(API_KEY + ":")}'
```

> **Permission API Key:** Setiap endpoint membutuhkan permission spesifik yang harus diaktifkan pada API Key di Xendit Dashboard → Settings → API Keys → Edit Permissions.

---

## Konsep Utama

xenPlatform memungkinkan platform marketplace untuk:
- Membuat dan mengelola **sub-account** untuk setiap seller
- Melakukan transaksi **atas nama** seller menggunakan header `for-user-id`
- **Split payment** otomatis ke beberapa akun saat pembayaran masuk
- **Transfer saldo** antar akun (master ↔ sub-account, atau antar sub-account)
- Melakukan **verifikasi KYC** seller melalui Account Holder API

### Tipe Sub-Account

Ada dua tipe sub-account di xenPlatform. Pilih berdasarkan model bisnis dan pengalaman yang ingin diberikan kepada seller.

> ⚠️ **Catatan penting:** `OWNED` sub-account saat ini **dinonaktifkan secara default untuk akun yang berbasis di Indonesia**. Hubungi support Xendit untuk mengaktifkannya.

| | `OWNED` | `MANAGED` |
|---|---|---|
| **Deskripsi** | Sub-account yang merepresentasikan bisnis yang **dimiliki atau dikelola oleh platform itu sendiri**. Tidak memerlukan verifikasi terpisah. | Sub-account yang merepresentasikan **merchant/seller pihak ketiga** yang transaksinya dilakukan atas nama mereka. |
| **Cocok untuk** | Platform yang bertransaksi atas nama merchant dan hanya perlu memisahkan saldo per partner | Platform yang memberikan pengalaman pembayaran lebih personal ke merchant, atau ketika customer bertransaksi langsung dengan merchant |
| **Verifikasi** | Langsung `LIVE` dan siap bertransaksi **tanpa verifikasi**. Verifikasi via API tetap bisa dilakukan jika diperlukan (misal untuk mengaktifkan channel tertentu atau jika platform adalah payment reseller). | Merchant harus menyelesaikan proses onboarding di dashboard Xendit sebelum bisa menerima pembayaran. |
| **Akses dashboard** | Seller **tidak bisa** melihat atau mengakses akun ini kecuali diundang secara eksplisit | Merchant **bisa melihat dan mengelola** dana mereka sendiri, sementara platform tetap punya kontrol penuh |
| **Nama yang tampil ke customer saat bayar** | **Nama Master Account** (nama platform) | **Nama Sub-Account** (nama merchant/seller) |
| **Pihak yang menerima invoice dari Xendit** | Master Account | Sub-Account (masing-masing merchant) |
| **Biaya transaksi dipotong dari** | Sub-Account | Sub-Account |
| **Opsi payout ke bank** | Disbursement saja | Disbursement **dan** Withdrawal (merchant bisa withdraw sendiri) |

### Header Penting

- **`for-user-id`**: Sertakan di setiap API call standar Xendit untuk membuat transaksi atas nama sub-account tertentu. Isinya adalah `account.id` yang didapat dari Create Account.

---

## 1. Account Management API

### 1.1 Create Account

**`POST /v2/accounts`**  
**Permission:** Account Write

**Fungsi:** Membuat sub-account baru untuk seller di platform.

- Untuk tipe **`OWNED`**: Akun langsung berstatus `LIVE` setelah dibuat dan siap bertransaksi.
- Untuk tipe **`MANAGED`**: Akun dibuat dengan status `INVITED`. Xendit mengirim email undangan ke seller. Akun baru siap bertransaksi setelah seller menyelesaikan verifikasi (status berubah ke `LIVE`).

> ⚠️ `OWNED` dinonaktifkan secara default untuk platform berbasis **Indonesia**. Hubungi Xendit untuk mengaktifkannya.

**Rate limit:** Maksimal 5 write request per detik.

**Contoh Request:**
```bash
curl --request POST \
  --url https://api.xendit.co/v2/accounts \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}' \
  --header 'content-type: application/json' \
  --data '{
    "email": "seller@example.com",
    "type": "OWNED",
    "public_profile": {
      "business_name": "Nama Toko Seller"
    }
  }'
```

**Request Body:**
```json
{
  "email": "seller@example.com",
  "type": "OWNED",
  "public_profile": {
    "business_name": "Nama Toko Seller"
  }
}
```

**Response (200):**
```json
{
  "id": "5cafeb170a2b18519b1b8761",
  "type": "OWNED",
  "email": "seller@example.com",
  "public_profile": { "business_name": "Nama Toko Seller" },
  "status": "LIVE"
}
```

**Field penting di response:**
- `id` → Simpan ini sebagai `seller_account_id`. Digunakan di header `for-user-id` dan di semua API transfer/split.
- `status` → Harus `LIVE` sebelum bisa menerima pembayaran.

**Status akun yang mungkin:** `INVITED`, `REGISTERED`, `AWAITING_DOCS`, `PENDING_VERIFICATION`, `LIVE`, `SUSPENDED`

---

### 1.2 List Accounts

**`GET /v2/accounts`**  
**Permission:** Accounts Read

**Fungsi:** Mengambil daftar semua sub-account yang terhubung ke platform, dengan dukungan filter dan pagination.

**Contoh Request:**
```bash
curl --request GET \
  --url 'https://api.xendit.co/v2/accounts?limit=10&type=OWNED' \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}'
```

**Query Parameters (semua opsional):**
| Parameter | Tipe | Keterangan |
|-----------|------|-----------|
| `email` | array string | Filter berdasarkan email |
| `status` | array string | Filter berdasarkan status |
| `public_profile.business_name` | string | Filter berdasarkan nama bisnis |
| `type` | string | `MANAGED` atau `OWNED` |
| `created[gte]` / `created[lte]` | datetime | Filter berdasarkan tanggal buat |
| `limit` | number | Jumlah hasil (1–50, default 10) |
| `after_id` | string | Untuk pagination (gunakan dari field `links` di response) |

**Response (200):**
```json
{
  "data": [ /* array account objects */ ],
  "has_more": true,
  "links": [{ "href": "/v2/accounts?after_id=...", "rel": "next", "method": "GET" }]
}
```

---

### 1.3 Get Account

**`GET /v2/accounts/{id}`**  
**Permission:** Accounts Read

**Fungsi:** Mengambil detail satu sub-account berdasarkan ID-nya.

**Contoh Request:**
```bash
curl --request GET \
  --url https://api.xendit.co/v2/accounts/5cafeb170a2b18519b1b8761 \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}'
```

**Path Parameter:** `id` — Account ID (contoh: `5cafeb170a2b18519b1b8761`)

**Response (200):** Object account lengkap (sama seperti response Create Account).

**Error:**
- `404 DATA_NOT_FOUND` — ID tidak valid atau tidak ditemukan.

---

### 1.4 Update Account

**`PATCH /v2/accounts/{id}`**  
**Permission:** Account Write

**Fungsi:** Mengupdate informasi sub-account, seperti nama bisnis, email, atau menghubungkan Account Holder untuk verifikasi KYC.

**Contoh Request — Link Account Holder:**
```bash
curl --request PATCH \
  --url https://api.xendit.co/v2/accounts/5cafeb170a2b18519b1b8761 \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}' \
  --header 'content-type: application/json' \
  --data '{ "account_holder_id": "4376b7b0-1c44-46be-8640-828f79cdc8be" }'
```

**Contoh Use Case:**

1. **Update nama bisnis:**
```json
{ "public_profile": { "business_name": "Nama Baru" } }
```

2. **Link Account Holder (untuk KYC):**
```json
{ "account_holder_id": "4376b7b0-1c44-46be-8640-828f79cdc8be" }
```

**Response (200):** Object account lengkap, termasuk `account_holder_id` jika sudah di-link.

---

## 2. Split Rule API

### 2.1 Create Split Rule

**`POST /split_rules`**  
**Permission:** API key master (gunakan API key master, bukan sub-account)

**Fungsi:** Membuat aturan pembagian pembayaran (split) ke beberapa akun secara otomatis. `split_rule_id` yang dihasilkan kemudian disertakan di header `with-split-rule` saat membuat transaksi pembayaran standar Xendit.

**Cara kerja:** Saat pembayaran settle, sistem otomatis mendistribusikan dana ke akun-akun yang ditentukan dalam `routes`, berdasarkan flat amount atau persentase.

**Contoh Request:**
```bash
curl --request POST \
  --url https://api.xendit.co/split_rules \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}' \
  --header 'content-type: application/json' \
  --data '{
    "name": "Komisi Marketplace Standard",
    "description": "Platform fee 5% + biaya pengiriman Rp3.000",
    "routes": [
      {
        "flat_amount": 3000,
        "currency": "IDR",
        "destination_account_id": "ID_AKUN_TUJUAN_1",
        "reference_id": "shipping-fee"
      },
      {
        "percent_amount": 5.25,
        "currency": "IDR",
        "destination_account_id": "ID_AKUN_TUJUAN_2",
        "reference_id": "platform-commission"
      }
    ]
  }'
```

**Request Body:**
```json
{
  "name": "Komisi Marketplace Standard",
  "description": "Platform fee 5% + biaya pengiriman Rp3.000",
  "routes": [
    {
      "flat_amount": 3000,
      "currency": "IDR",
      "destination_account_id": "ID_AKUN_TUJUAN_1",
      "reference_id": "shipping-fee"
    },
    {
      "percent_amount": 5.25,
      "currency": "IDR",
      "destination_account_id": "ID_AKUN_TUJUAN_2",
      "reference_id": "platform-commission"
    }
  ]
}
```

**Field routes:**
| Field | Tipe | Keterangan |
|-------|------|-----------|
| `flat_amount` | number | Jumlah tetap (wajib jika `percent_amount` null) |
| `percent_amount` | number | Persentase 0–100, max 2 desimal (wajib jika `flat_amount` null) |
| `currency` | string | `IDR`, `PHP`, `VND`, `MYR`, `THB`, `SGD`, `USD` |
| `destination_account_id` | string | Business ID akun tujuan (master atau sub-account) |
| `reference_id` | string | Harus unik antar routes dalam Split Rule yang sama |

**Response (200):**
```json
{
  "id": "splitru_d9e069f2-4da7-4562-93b7-ded87023d749",
  "name": "Komisi Marketplace Standard",
  "routes": [ /* ... */ ]
}
```

**Field penting di response:**
- `id` (format `splitru_xxx`) → Gunakan di header `with-split-rule: {split_rule_id}` saat membuat Invoice atau Payment.

**Error umum:**
- `400 INVALID_FEE_AMOUNT` — Jumlah fee negatif atau format salah.
- `400 DUPLICATE_ERROR` — `reference_id` tidak unik antar routes.
- `404 DESTINATION_ACCOUNT_NOT_FOUND` — `destination_account_id` tidak valid atau tidak terhubung ke platform.

---

## 3. Account Holder API (KYC Seller)

Account Holder merepresentasikan entitas legal yang memegang akun Xendit. Berisi data bisnis, data individu (PIC), alamat, dan dokumen KYC.

**Kapan Account Holder diperlukan:**
- Selalu diperlukan untuk sub-account tipe **`MANAGED`** sebelum bisa menerima pembayaran.
- Untuk sub-account tipe **`OWNED`**, verifikasi via Account Holder diperlukan jika:
  1. Platform adalah **Payment Service Provider (reseller)** yang onboarding merchant pihak ketiga ke Xendit; atau
  2. Ingin mengaktifkan **payment channel yang memerlukan verifikasi tambahan** (contoh: kartu kredit, GCash PH).

**Alur kerja:**
1. **Create Account Holder** → dapatkan `account_holder_id`
2. **Update Account** → link `account_holder_id` ke sub-account seller
3. Poll **Get Account Holder** secara berkala untuk memantau status KYC
4. Jika diperlukan, **Update Account Holder** untuk mengirim ulang dokumen atau mengaktifkan payment channel baru

---

### 3.1 Create Account Holder

**`POST /account_holders`**  
**Permission:** Account Holder Write

**Fungsi:** Membuat object Account Holder yang berisi data legal dan dokumen KYC seller.

**Contoh Request:**
```bash
curl --request POST \
  --url https://api.xendit.co/account_holders \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}' \
  --header 'content-type: application/json' \
  --data '{
    "business_detail": {
      "type": "CORPORATION",
      "legal_name": "PT Contoh Seller",
      "trading_name": "Nama Toko",
      "description": "Deskripsi bisnis",
      "industry_category": "ELECTRONICS_AND_ACCESSORIES",
      "date_of_registration": "2023-01-01",
      "country_of_operation": "ID"
    },
    "individual_details": [{
      "type": "PIC",
      "role": "owner",
      "given_names": "John",
      "surname": "Doe",
      "phone_number": "+62811234567",
      "email": "john@example.com",
      "nationality": "ID",
      "date_of_birth": "1990-01-01",
      "gender": "MALE"
    }],
    "address": {
      "country": "ID",
      "city": "Jakarta",
      "province_state": "DKI Jakarta",
      "street_line1": "Jl. Sudirman No. 1",
      "postal_code": "10220"
    },
    "kyc_documents": [{
      "type": "NAMA_TIPE_DOKUMEN",
      "country": "ID",
      "file_id": "ID_FILE_DARI_UPLOAD_API"
    }],
    "website_url": "https://tokoanda.com",
    "phone_number": "+62811234567",
    "email": "seller@example.com"
  }'
```

**Request Body (ringkasan field penting):**
```json
{
  "business_detail": {
    "type": "CORPORATION",           // CORPORATION, PARTNERSHIP, SOLE_PROPRIETORSHIP, INDIVIDUAL, dll.
    "legal_name": "PT Contoh Seller",
    "trading_name": "Nama Toko",
    "description": "Deskripsi bisnis",
    "industry_category": "ELECTRONICS_AND_ACCESSORIES",
    "date_of_registration": "2023-01-01",
    "country_of_operation": "ID"     // ID, PH, VN, MY, TH
  },
  "individual_details": [
    {
      "type": "PIC",                 // Person in Charge — minimal 1 wajib
      "role": "owner",
      "given_names": "John",
      "surname": "Doe",
      "phone_number": "+62811234567",
      "email": "john@example.com",
      "nationality": "ID",
      "date_of_birth": "1990-01-01",
      "gender": "MALE"
    }
  ],
  "address": {
    "country": "ID",
    "city": "Jakarta",
    "province_state": "DKI Jakarta",
    "street_line1": "Jl. Sudirman No. 1",
    "postal_code": "10220"
  },
  "kyc_documents": [
    {
      "type": "NAMA_TIPE_DOKUMEN",   // Lihat docs Xendit untuk daftar lengkap per negara
      "country": "ID",
      "file_id": "ID_FILE_DARI_UPLOAD_API"
    }
  ],
  "website_url": "https://tokoanda.com",
  "phone_number": "+62811234567",
  "email": "seller@example.com"
}
```

**Response (200):**
```json
{
  "id": "4376b7b0-1c44-46be-8640-828f79cdc8be",
  "kyc": { "status": "NOT_VERIFIED" },
  "created_at": "2023-03-30T11:41:57.881Z"
}
```

**Field penting:**
- `id` → Simpan sebagai `account_holder_id`, gunakan untuk link ke sub-account via Update Account.
- `kyc.status` → Dimulai dari `NOT_VERIFIED`. Pantau perubahannya dengan polling Get Account Holder.

---

### 3.2 Get Account Holder

**`GET /account_holders/{id}`**  
**Permission:** Account Holder Read

**Fungsi:** Mengambil detail Account Holder dan status KYC-nya saat ini. Gunakan untuk polling status verifikasi setelah submit dokumen.

**Contoh Request:**
```bash
curl --request GET \
  --url https://api.xendit.co/account_holders/4376b7b0-1c44-46be-8640-828f79cdc8be \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}'
```

**Path Parameter:** `id` — Account Holder ID

**Response (200):** Object Account Holder lengkap termasuk status `kyc` dan `capabilities`.

**KYC status yang mungkin:**
| Status | Keterangan | Tindakan |
|--------|-----------|----------|
| `NOT_VERIFIED` | Belum diverifikasi (status awal) | Tunggu proses review Xendit |
| `PASSED` | KYC berhasil, akun seller terverifikasi | Update status seller di DB → aktif |
| `FAILED` | KYC gagal | Notifikasi seller untuk menghubungi support |
| `RESUBMISSION_REQUIRED` | Dokumen tidak valid/kurang | Minta seller upload ulang via Update Account Holder |

---

### 3.3 Update Account Holder

**`PATCH /account_holders/{id}`**  
**Permission:** Account Holder Write

**Fungsi:** Mengupdate data/dokumen Account Holder, atau mengaktifkan payment channel baru untuk seller.

**Contoh Request — Aktifkan payment channel:**
```bash
curl --request PATCH \
  --url https://api.xendit.co/account_holders/4376b7b0-1c44-46be-8640-828f79cdc8be \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}' \
  --header 'content-type: application/json' \
  --data '{
    "capabilities": [
      { "type": "MONEY_IN", "channel_code": "ID_OVO" },
      { "type": "MONEY_IN", "channel_code": "ID_GOPAY" }
    ]
  }'
```

**Use Case 1 — Update dokumen KYC yang ditolak:**
```json
{
  "kyc_documents": [
    { "country": "ID", "type": "NAMA_DOKUMEN", "file_id": "FILE_ID_BARU" }
  ]
}
```

**Use Case 2 — Update URL website:**
```json
{ "website_url": "https://website-baru.com" }
```

**Use Case 3 — Aktifkan payment channel untuk seller:**
```json
{
  "capabilities": [
    { "type": "MONEY_IN", "channel_code": "ID_OVO" },
    { "type": "MONEY_IN", "channel_code": "ID_GOPAY" }
  ]
}
```

**Response (200):** Object Account Holder ter-update, termasuk array `capabilities` dengan status masing-masing channel.

**Status capabilities:** `VERIFICATION_IN_PROGRESS`, `LIVE`, `RESUBMISSION_REQUIRED`, `DECLINED`

**Error penting:**
- `403 KYC_VERIFICATION_IN_PROGRESS` — Tidak bisa update saat KYC sedang diproses.
- `403 CHANNEL_HAS_BEEN_ACTIVATED` — Channel sudah aktif, tidak perlu aktivasi ulang.
- `400 INSUFFICIENT_ACCOUNT_HOLDER_DATA` — Data tidak cukup untuk aktivasi channel yang diminta.

---

## 4. Transfers API

### 4.1 Create Transfer

**`POST /transfers`**  
**Permission:** API key master platform (sub-account tidak bisa buat transfer)

**Fungsi:** Mentransfer saldo antar akun dalam ekosistem Xendit: master → sub-account, sub-account → master, atau antar sub-account.

**Contoh Request:**
```bash
curl --request POST \
  --url https://api.xendit.co/transfers \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}' \
  --header 'content-type: application/json' \
  --data '{
    "reference": "transfer-unik-ref-001",
    "amount": 500000,
    "source_user_id": "ID_AKUN_SUMBER",
    "destination_user_id": "ID_AKUN_TUJUAN"
  }'
```

**Request Body:**
```json
{
  "reference": "transfer-unik-ref-001",
  "amount": 500000,
  "source_user_id": "ID_AKUN_SUMBER",
  "destination_user_id": "ID_AKUN_TUJUAN"
}
```

**Field:**
| Field | Keterangan |
|-------|-----------|
| `reference` | Unik per transfer. Jika reference sama dengan payload berbeda → error. Jika sama dengan payload sama → idempotent retry. |
| `amount` | Harus > 0. IDR tidak boleh desimal. PHP max 2 desimal. |
| `source_user_id` | Business ID akun sumber (bisa master atau sub-account) |
| `destination_user_id` | Business ID akun tujuan |

**Response (200):**
```json
{
  "transfer_id": "bd1cc56b-...",
  "reference": "transfer-unik-ref-001",
  "source_user_id": "...",
  "destination_user_id": "...",
  "amount": 500000,
  "status": "SUCCESSFUL"
}
```

**Status transfer:** `SUCCESSFUL`, `PENDING`, `FAILED`

**Error penting:**
- `400 INSUFFICIENT_BALANCE` — Saldo akun sumber tidak cukup.
- `400 INVALID_SOURCE_OR_DESTINATION_ERROR` — ID akun tidak valid atau tidak terhubung ke platform.
- `403 XEN_PLATFORM_SUB_ACCOUNT_NOT_LIVE` — Akun sumber atau tujuan belum berstatus LIVE.
- `403 DUPLICATE_REFERENCE` — Reference sudah digunakan dengan payload berbeda.
- `425 TRANSFER_IN_PROGRESS` — Transfer sedang diproses; gunakan Get Transfer by Reference untuk cek status.

---

### 4.2 Get Transfer by Reference

**`GET /transfers?reference={reference}`**  
**Permission:** API key master platform

**Fungsi:** Mengambil detail transfer berdasarkan reference. Berguna untuk cek status transfer yang masih `PENDING` atau setelah mendapat error `425 TRANSFER_IN_PROGRESS`.

**Contoh Request:**
```bash
curl --request GET \
  --url 'https://api.xendit.co/transfers?reference=transfer-unik-ref-001' \
  --header 'accept: application/json' \
  --header 'authorization: Basic {BASE64(API_KEY + ":")}'
```

**Query Parameter:** `reference` — Reference string yang digunakan saat Create Transfer.

**Response (200):** Object transfer dengan status terkini.

---

## Ringkasan Alur Integrasi Marketplace Multi-Seller

### Onboarding Seller Baru — Tipe OWNED
Cocok jika platform ingin kontrol penuh, seller tidak perlu akses dashboard Xendit, dan nama platform yang tampil ke customer.
```
1. POST /v2/accounts               → Buat sub-account (type: "OWNED") — langsung LIVE
2. Simpan account.id sebagai seller_account_id
3. [Opsional] Jika perlu aktivasi channel tertentu (kartu kredit, dll.):
   a. POST /account_holders        → Buat Account Holder dengan data KYC seller
   b. PATCH /v2/accounts/{id}      → Link account_holder_id ke sub-account
   c. Poll GET /account_holders/{id} sampai kyc.status = PASSED
   d. PATCH /account_holders/{id}  → Aktifkan payment channels
   e. Poll GET /account_holders/{id} sampai capabilities[].status = LIVE
```

### Onboarding Seller Baru — Tipe MANAGED
Cocok jika seller adalah merchant pihak ketiga yang ingin nama mereka tampil ke customer, bisa akses dashboard sendiri, dan bisa withdrawal mandiri.
```
1. POST /v2/accounts               → Buat sub-account (type: "MANAGED") — status: INVITED
   → Xendit otomatis kirim email undangan ke seller
2. Seller menyelesaikan registrasi di dashboard Xendit
3. Poll GET /v2/accounts/{id} sampai status = LIVE (menandakan seller sudah terverifikasi)
4. [Opsional] POST /account_holders  → Buat Account Holder via API jika ingin
   proses KYC melalui platform (tidak hanya dari dashboard seller)
5. [Opsional] PATCH /account_holders/{id} → Aktifkan payment channels tambahan
```

### Setup Split Payment untuk Seller
```
1. POST /split_rules               → Buat Split Rule (definisikan komisi platform & seller)
2. Simpan split_rule_id dari response
3. Saat buat Invoice/Payment → sertakan header: with-split-rule: {split_rule_id}
```

### Transfer Dana Manual
```
1. POST /transfers                 → Transfer saldo antar akun
2. Jika dapat 425 → GET /transfers?reference={ref} untuk cek status terkini
3. Periksa status SUCCESSFUL / PENDING / FAILED
```

### Manajemen Akun
```
GET  /v2/accounts            → List semua seller (support filter & pagination)
GET  /v2/accounts/{id}       → Detail seller tertentu
PATCH /v2/accounts/{id}      → Update info seller atau link Account Holder
GET  /account_holders/{id}   → Cek status KYC dan capabilities seller
```

---

## Error Codes Umum

| Error Code | Keterangan |
|-----------|-----------|
| `API_VALIDATION_ERROR` | Request body tidak valid |
| `DATA_NOT_FOUND` | ID tidak ditemukan |
| `DESTINATION_ACCOUNT_NOT_FOUND` | `destination_account_id` tidak valid |
| `INSUFFICIENT_BALANCE` | Saldo tidak cukup |
| `INVALID_SOURCE_OR_DESTINATION_ERROR` | Source/destination ID tidak valid |
| `DUPLICATE_REFERENCE` | Reference sudah digunakan dengan payload berbeda |
| `XEN_PLATFORM_SUB_ACCOUNT_NOT_LIVE` | Akun belum berstatus LIVE |
| `INVALID_FEE_AMOUNT` | Jumlah fee tidak valid |
| `FEATURE_NOT_ACTIVATED` | Fitur belum diaktifkan untuk akun ini |
| `KYC_VERIFICATION_IN_PROGRESS` | Tidak bisa update saat KYC sedang berlangsung |
| `CHANNEL_HAS_BEEN_ACTIVATED` | Channel sudah aktif |
| `TRANSFER_IN_PROGRESS` | Transfer sedang diproses, cek status dengan GET |
| `MISMATCH_PAYLOAD_FOR_REFERENCE` | Reference digunakan ulang dengan payload berbeda |
| `INSUFFICIENT_ACCOUNT_HOLDER_DATA` | Data Account Holder tidak cukup untuk aktivasi channel |
