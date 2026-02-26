# Petunjuk Teknis: Integrasi Xendit Payment Link untuk Marketplace Backend

> **Konteks Sistem:** Dokumen ini ditujukan untuk AI agent yang mengerjakan fitur pembayaran pada sistem backend marketplace. Sistem menggunakan **Xendit xenPlatform** sebagai payment gateway, dengan model **Master Account** (platform) yang memproses pembayaran atas nama **Sub-Account** (merchant/seller). Mekanisme pembayaran yang digunakan adalah **Payment Link** (Invoice API v2).

---

## 1. Konsep Dasar

### Payment Link adalah Invoice

Di Xendit, **Payment Link** dibangun di atas **Invoice API** (`POST /v2/invoices`). Ketika dokumentasi menyebut "Payment Link", artinya adalah invoice yang menghasilkan URL checkout yang di-host oleh Xendit. Customer tidak perlu membangun halaman checkout sendiri — cukup redirect customer ke `invoice_url` yang dikembalikan API.

### Arsitektur xenPlatform untuk Marketplace

```
┌─────────────────────────────────────────────┐
│           MASTER ACCOUNT (Platform)          │
│  - Pemegang API Key                         │
│  - Mengatur semua transaksi sub-account     │
└────────────────┬────────────────────────────┘
                 │ for-user-id header
        ┌────────┴─────────┐
        ▼                  ▼
┌──────────────┐   ┌──────────────┐
│ Sub-Account  │   │ Sub-Account  │
│  (Seller A)  │   │  (Seller B)  │
│  user_id: X  │   │  user_id: Y  │
└──────────────┘   └──────────────┘
```

**Poin penting:**
- Satu API Key Master Account digunakan untuk semua request.
- Untuk membuat Payment Link atas nama seller tertentu, sertakan header `for-user-id` berisi `user_id` sub-account seller tersebut.
- Dana masuk akan langsung tercatat di balance sub-account seller yang bersangkutan.

---

## 2. Konfigurasi Awal yang Diperlukan

Sebelum mengerjakan fitur, pastikan hal-hal berikut sudah siap:

| Kebutuhan | Keterangan |
|---|---|
| `XENDIT_API_KEY` | Secret key dari Master Account (Settings → API Key → beri permission `WRITE` untuk Money In) |
| `XENDIT_WEBHOOK_TOKEN` | Callback token untuk memverifikasi keaslian webhook |
| `WEBHOOK_URL` | URL endpoint backend yang akan menerima notifikasi dari Xendit |
| Sub-Account `user_id` | ID unik tiap seller, didapat saat pembuatan sub-account via xenPlatform API |

**Base URL API:**
```
https://api.xendit.co
```

**Autentikasi:** HTTP Basic Auth — gunakan API Key sebagai `username`, password dikosongkan.

```
Authorization: Basic base64(XENDIT_API_KEY:)
```

---

## 3. Alur Sistem Payment Link (End-to-End)

```
[Customer checkout] 
       │
       ▼
[Backend: POST /v2/invoices]  ← sertakan for-user-id: {seller_user_id}
       │
       ▼
[Xendit mengembalikan invoice_url]
       │
       ▼
[Backend redirect / kirim link ke customer]
       │
       ▼
[Customer membayar di halaman Xendit]
       │
       ├── Berhasil → redirect ke success_redirect_url
       │                + Xendit kirim webhook PAID ke backend
       │
       └── Gagal → customer bisa retry di halaman yang sama
                    (tidak perlu buat invoice baru)
```

---

## 4. API: Membuat Payment Link

### Endpoint

```
POST https://api.xendit.co/v2/invoices
```

### Headers

```http
Content-Type: application/json
Authorization: Basic {base64(XENDIT_API_KEY:)}
for-user-id: {seller_sub_account_user_id}
```

> ⚠️ Header `for-user-id` **wajib disertakan** untuk marketplace agar dana masuk ke balance seller yang benar, bukan ke balance master account.

### Request Body

```json
{
  "external_id": "ORDER-2024-001",
  "amount": 510000,
  "description": "Pembelian Air Conditioner dari Toko ABC",
  "invoice_duration": 86400,
  "customer": {
    "given_names": "Budi",
    "surname": "Santoso",
    "email": "budi@example.com",
    "mobile_number": "+6281234567890"
  },
  "success_redirect_url": "https://yourdomain.com/payment/success?order_id=ORDER-2024-001",
  "failure_redirect_url": "https://yourdomain.com/payment/failed?order_id=ORDER-2024-001",
  "currency": "IDR",
  "items": [
    {
      "name": "Air Conditioner 1 PK",
      "quantity": 1,
      "price": 510000,
      "category": "Electronic",
      "url": "https://yourdomain.com/products/ac-1pk"
    }
  ],
  "metadata": {
    "order_id": "ORDER-2024-001",
    "seller_id": "SELLER-XYZ"
  }
}
```

### Penjelasan Field Penting

| Field | Tipe | Wajib | Keterangan |
|---|---|---|---|
| `external_id` | string | ✅ | ID unik dari sistem kita (misalnya order ID). Digunakan untuk idempotency dan tracking. Tidak boleh duplikat untuk satu sub-account. |
| `amount` | number | ✅ | Total jumlah yang harus dibayar customer (IDR, tanpa desimal). |
| `description` | string | ✅ | Deskripsi transaksi yang ditampilkan di halaman checkout. |
| `invoice_duration` | number | ❌ | Durasi kedaluwarsa link dalam detik. Default: 86400 (24 jam). Maksimum: 31536000 (1 tahun). |
| `customer.email` | string | ❌ | Email customer untuk notifikasi dari Xendit. |
| `customer.mobile_number` | string | ❌ | Format E.164 (contoh: `+6281234567890`). |
| `success_redirect_url` | string | ❌ | URL redirect setelah pembayaran berhasil. |
| `failure_redirect_url` | string | ❌ | URL redirect jika pembayaran gagal/dibatalkan. |
| `currency` | string | ❌ | Default: `IDR`. |
| `items` | array | ❌ | Detail item yang dibeli. Ditampilkan di halaman checkout. |
| `metadata` | object | ❌ | Data tambahan bebas dari sistem kita. Tidak diproses Xendit, dikembalikan di webhook. Gunakan untuk menyimpan order ID, seller ID, dll. |

### Response Sukses (200)

```json
{
  "id": "6748105a77f16ebe0cc583a7",
  "external_id": "ORDER-2024-001",
  "user_id": "62440e322008e87fb29c1fd0",
  "status": "PENDING",
  "merchant_name": "Nama Toko Seller",
  "amount": 510000,
  "description": "Pembelian Air Conditioner dari Toko ABC",
  "expiry_date": "2024-11-29T06:40:26.353Z",
  "invoice_url": "https://checkout.xendit.co/web/6748105a77f16ebe0cc583a7",
  "available_banks": [...],
  "available_ewallets": [...],
  "available_qr_codes": [...],
  "currency": "IDR",
  "created": "2024-11-28T06:40:26.508Z",
  "metadata": {
    "order_id": "ORDER-2024-001",
    "seller_id": "SELLER-XYZ"
  }
}
```

**Field penting dari response yang harus disimpan ke database:**

| Field | Fungsi |
|---|---|
| `id` | Xendit Invoice ID — gunakan untuk query status atau expire invoice |
| `external_id` | Konfirmasi bahwa external_id kita sudah diterima |
| `invoice_url` | URL yang dikirimkan ke customer untuk melakukan pembayaran |
| `status` | Status awal selalu `PENDING` |
| `expiry_date` | Waktu kedaluwarsa link (ISO 8601) |

---

## 5. Siklus Status Payment Link

```
PENDING ──── bayar sebelum expired ──→ PAID
   │
   └─── expired / di-expire manual ──→ EXPIRED
```

| Status | Deskripsi |
|---|---|
| `PENDING` | Invoice baru dibuat, menunggu pembayaran customer. |
| `PAID` | Pembayaran berhasil. Webhook dikirim ke backend. |
| `EXPIRED` | Waktu habis (`invoice_duration` tercapai) atau di-expire manual. Tidak bisa diaktifkan kembali. Customer tidak bisa lagi membuka `invoice_url`. |

> ℹ️ **Tidak ada status `REFUNDED`** pada Payment Link. Refund harus dilacak secara terpisah melalui tab Transaction/Related Payment di dashboard Xendit.

---

## 6. Webhook: Menerima Notifikasi Pembayaran

### Konfigurasi Webhook

Webhook URL dikonfigurasi di Xendit Dashboard → Settings → Webhooks. Untuk xenPlatform, Master Account dapat set webhook untuk sub-account via API dengan `for-user-id` header.

**Tipe webhook yang relevan:** `invoice`

### Payload Webhook (Status PAID)

```json
{
  "id": "6748105a77f16ebe0cc583a7",
  "external_id": "ORDER-2024-001",
  "user_id": "62440e322008e87fb29c1fd0",
  "status": "PAID",
  "amount": 510000,
  "paid_amount": 510000,
  "paid_at": "2024-11-28T06:41:29.776Z",
  "payment_method": "EWALLET",
  "payment_channel": "DANA",
  "payment_id": "ewc_99013716-9354-439e-9fea-e61a06e87ea3",
  "currency": "IDR",
  "description": "Pembelian Air Conditioner dari Toko ABC",
  "merchant_name": "Nama Toko Seller",
  "items": [...],
  "metadata": {
    "order_id": "ORDER-2024-001",
    "seller_id": "SELLER-XYZ"
  },
  "success_redirect_url": "https://yourdomain.com/payment/success",
  "failure_redirect_url": "https://yourdomain.com/payment/failed"
}
```

### Payload Webhook (Status EXPIRED)

```json
{
  "id": "6748105a77f16ebe0cc583a7",
  "external_id": "ORDER-2024-001",
  "status": "EXPIRED",
  ...
}
```

### Verifikasi Keaslian Webhook

**Wajib** — setiap request webhook dari Xendit menyertakan header `x-callback-token`. Bandingkan nilainya dengan `XENDIT_WEBHOOK_TOKEN` yang tersimpan di server. Jika tidak cocok, **abaikan request** (kemungkinan spoofed).

```javascript
// Contoh verifikasi (Node.js/Express)
app.post('/webhook/xendit/invoice', (req, res) => {
  const callbackToken = req.headers['x-callback-token'];
  
  if (callbackToken !== process.env.XENDIT_WEBHOOK_TOKEN) {
    return res.status(401).json({ message: 'Unauthorized' });
  }
  
  const payload = req.body;
  
  if (payload.status === 'PAID') {
    // Update status order di database
    // Misalnya: updateOrder(payload.external_id, 'PAID')
  }
  
  if (payload.status === 'EXPIRED') {
    // Update status order menjadi expired/cancelled
  }
  
  // Selalu balas 200 agar Xendit tidak retry
  return res.status(200).json({ received: true });
});
```

> ⚠️ **Selalu kembalikan HTTP 200** setelah memproses webhook. Jika Xendit tidak menerima 200, ia akan melakukan retry beberapa kali. Pastikan logika pemrosesan idempotent (aman dijalankan berkali-kali untuk ID yang sama).

---

## 7. API: Expire Payment Link Manual

Gunakan endpoint ini jika order dibatalkan oleh sistem (misalnya stok habis, atau customer membatalkan order) sebelum customer sempat bayar.

### Endpoint

```
POST https://api.xendit.co/v2/invoices/{invoice_id}/expire!
```

### Headers

```http
Authorization: Basic {base64(XENDIT_API_KEY:)}
for-user-id: {seller_sub_account_user_id}
```

### Contoh Request

```bash
curl -X POST \
  https://api.xendit.co/v2/invoices/6748105a77f16ebe0cc583a7/expire! \
  -u XENDIT_API_KEY: \
  -H "for-user-id: 62440e322008e87fb29c1fd0"
```

---

## 8. API: Cek Status Payment Link

Gunakan endpoint ini untuk polling status invoice (alternatif/komplemen dari webhook).

### Endpoint

```
GET https://api.xendit.co/v2/invoices/{invoice_id}
```

### Headers

```http
Authorization: Basic {base64(XENDIT_API_KEY:)}
for-user-id: {seller_sub_account_user_id}
```

> ℹ️ Sebaiknya utamakan webhook untuk update status. Gunakan polling hanya sebagai fallback (misalnya saat webhook delay atau gagal diterima).

---

## 9. Skema Database yang Direkomendasikan

Simpan data berikut pada tabel `payment_links` (atau gabungkan ke tabel `orders`):

```sql
CREATE TABLE payment_links (
  id              VARCHAR PRIMARY KEY,          -- internal ID sistem kita
  order_id        VARCHAR NOT NULL,             -- referensi ke tabel orders
  seller_id       VARCHAR NOT NULL,             -- referensi ke tabel sellers
  xendit_id       VARCHAR UNIQUE,               -- field "id" dari response Xendit
  external_id     VARCHAR UNIQUE NOT NULL,      -- field "external_id" yang kita kirim
  invoice_url     VARCHAR,                      -- URL checkout untuk customer
  amount          BIGINT NOT NULL,              -- dalam satuan terkecil (IDR: rupiah)
  currency        VARCHAR(3) DEFAULT 'IDR',
  status          VARCHAR DEFAULT 'PENDING',    -- PENDING | PAID | EXPIRED
  expiry_date     TIMESTAMP,
  paid_at         TIMESTAMP,                    -- dari webhook saat status PAID
  payment_method  VARCHAR,                      -- EWALLET, BANK_TRANSFER, dsb.
  payment_channel VARCHAR,                      -- DANA, OVO, BCA, dsb.
  created_at      TIMESTAMP DEFAULT NOW(),
  updated_at      TIMESTAMP DEFAULT NOW()
);
```

---

## 10. Pola Implementasi yang Direkomendasikan

### Saat Customer Checkout

```
1. Validasi order (stok, harga, dsb.)
2. Generate external_id unik (misalnya: "ORD-{order_id}-{timestamp}")
3. Ambil seller_user_id dari database
4. Panggil POST /v2/invoices dengan for-user-id: {seller_user_id}
5. Simpan xendit_id, external_id, invoice_url ke database
6. Set status order = "WAITING_PAYMENT"
7. Return invoice_url ke frontend → redirect customer
```

### Saat Menerima Webhook

```
1. Verifikasi x-callback-token
2. Parse payload, ambil external_id dan status
3. Cek apakah external_id sudah pernah diproses (idempotency check)
4. Jika status = PAID:
   - Update order status = "PAID"
   - Update payment_link status = "PAID"
   - Simpan paid_at, payment_method, payment_channel
   - Trigger proses selanjutnya (notifikasi seller, update stok, dsb.)
5. Jika status = EXPIRED:
   - Update order status = "EXPIRED" atau "CANCELLED"
6. Return HTTP 200
```

### Saat Order Dibatalkan Sistem

```
1. Ambil xendit_id dari database berdasarkan order_id
2. Panggil POST /v2/invoices/{xendit_id}/expire! dengan for-user-id yang sesuai
3. Update status order = "CANCELLED"
4. (Webhook EXPIRED akan dikirim Xendit secara otomatis)
```

---

## 11. Error Handling

| HTTP Code | Error Code Xendit | Penanganan |
|---|---|---|
| `400` | `API_VALIDATION_ERROR` | Periksa field yang tidak valid di `message` response |
| `400` | `DUPLICATE_ERROR` | `external_id` sudah digunakan. Gunakan ID yang berbeda. |
| `401` | `AUTHENTICATION_ERROR` | API Key salah atau tidak memiliki permission |
| `403` | `FORBIDDEN_ERROR` | `for-user-id` tidak valid atau sub-account tidak aktif |
| `404` | `DATA_NOT_FOUND` | Invoice ID tidak ditemukan |
| `409` | `INVOICE_ALREADY_PAID` | Invoice sudah dibayar, tidak bisa di-expire |
| `500` | `SERVER_ERROR` | Error di sisi Xendit. Coba lagi dengan exponential backoff. |

---

## 12. Catatan Penting & Gotchas

- **`external_id` harus unik per sub-account**, bukan per master account. Artinya dua seller berbeda boleh memiliki `external_id` yang sama, tapi satu seller tidak boleh.
- **Invoice yang `EXPIRED` tidak bisa diaktifkan kembali.** Jika customer perlu bayar lagi, buat invoice baru dengan `external_id` yang berbeda.
- **Tidak ada status `FAILED` pada Payment Link.** Jika customer gagal bayar, mereka bisa mencoba ulang di halaman yang sama selama invoice belum `EXPIRED`.
- **Idempotency webhook** — implementasikan pengecekan apakah event sudah pernah diproses sebelum menjalankan business logic, karena Xendit bisa mengirim ulang webhook jika tidak menerima 200.
- **Gunakan `metadata`** untuk menyimpan data kontekstual (order ID internal, seller ID, dsb.) — data ini akan dikembalikan di payload webhook sehingga memudahkan mapping ke data sistem.
- **Mode staging/test:** Gunakan URL `https://checkout-staging.xendit.co` (dikembalikan otomatis jika menggunakan API Key test environment).

---

## 13. Referensi

| Sumber | URL |
|---|---|
| Dokumentasi Payment Link Overview | https://docs.xendit.co/v1/docs/payment-links-api-overview |
| API Reference Invoice | https://developers.xendit.co/api-reference/#create-invoice |
| xenPlatform Overview | https://docs.xendit.co/docs/xenplatform-overview |
| Webhook Settings | https://docs.xendit.co/apidocs/set-webhook-url |
| Accept Payments for Sub-Accounts | https://docs.xendit.co/xenplatform/accept-payments/ |
