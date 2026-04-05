# Checkout Flow — Dokumentasi Lengkap

## Daftar Isi

1. [Gambaran Umum](#gambaran-umum)
2. [Alur Checkout (Step-by-Step)](#alur-checkout-step-by-step)
3. [API Endpoints](#api-endpoints)
4. [Ongkos Kirim (Shipping)](#ongkos-kirim-shipping)
5. [Biaya Platform](#biaya-platform)
6. [Pembayaran (Xendit)](#pembayaran-xendit)
7. [Status Order & Transisi](#status-order--transisi)
8. [Request & Response DTO](#request--response-dto)
9. [Skema Database](#skema-database)
10. [Mekanisme Kompensasi (Rollback)](#mekanisme-kompensasi-rollback)

---

## Gambaran Umum

Checkout pada UmrahMart mengubah **item cart yang dipilih** menjadi satu atau lebih **order** — masing-masing satu order per vendor. Setiap order memiliki invoice pembayaran sendiri di Xendit (sub-account vendor). Ongkir dihitung secara real-time melalui **RajaOngkir Komerce API** berdasarkan alamat pengiriman pelanggan dan alamat gudang vendor.

```
┌─────────┐     ┌──────────────────┐     ┌──────────────┐     ┌───────────────┐
│  Cart    │────▶│ Checkout Preview │────▶│ Place Order   │────▶│ Xendit Invoice│
│ (items)  │     │ (pilih ongkir)   │     │ (create order)│     │ (bayar)       │
└─────────┘     └──────────────────┘     └──────────────┘     └───────────────┘
```

---

## Alur Checkout (Step-by-Step)

### Step 1 — Isi Keranjang

Customer menambahkan product variant ke keranjang via `POST /api/v1/users/cart/items`.

- Cart dibuat otomatis (lazy-create) — satu cart aktif per user (enforced via UNIQUE constraint di DB).
- Jika variant yang sama sudah ada di cart, quantity akan di-increment (upsert).
- Item baru otomatis memiliki `is_selected = true`.
- Validasi saat menambah item:
  - Variant harus aktif & exist
  - Product harus berstatus `published`
  - Stok mencukupi (`stock_on_hand >= qty`)

### Step 1.5 — Pilih Item yang Akan Di-checkout

Customer dapat mengubah status pilihan item via `PATCH /api/v1/users/cart/items/:itemId/selection`.

- Item dengan `is_selected = true` akan ikut pada checkout preview dan checkout.
- Item dengan `is_selected = false` tetap berada di cart, tetapi diabaikan oleh preview dan checkout.

### Step 2 — Preview Checkout

Customer memanggil `POST /api/v1/users/checkout/preview` dengan `address_id`.

1. **Validasi alamat**: Alamat harus milik user yang login dan memiliki `district_id`.
2. **Ambil item cart terpilih** (`is_selected = true`) dari cart aktif. Gagal jika tidak ada item terpilih.
3. **Enrich setiap item terpilih**: Ambil data variant + product terkini. Gagal jika ada item yang sudah tidak tersedia.
4. **Kelompokkan item terpilih per vendor**.
5. **Untuk setiap vendor**:
   - Ambil alamat gudang vendor (alamat default pemilik vendor — harus punya `district_id`).
   - Ambil daftar kurir yang dipilih vendor dari tabel `vendor_couriers`.
   - Hitung subtotal dan total berat (default **500 gram** per item jika berat tidak diset).
   - Panggil **RajaOngkir** `CalculateDomesticCost` dengan origin district, destination district, total berat, dan kode kurir yang dipisah titik dua.
6. **Return response**: Item terpilih dikelompokkan per vendor, lengkap dengan `shipping_options` (pilihan ongkir) dan `platform_fee` tetap **Rp 6.000**.

### Step 3 — Place Order (Checkout)

Customer memanggil `POST /api/v1/users/checkout` dengan `address_id`, `shipping_choices`, dan optional `notes`.

1. **Validasi ulang alamat** dan parse pilihan pengiriman per vendor.
2. **Ambil item cart terpilih**, enrich, dan validasi stok (`qty <= stock_on_hand`).
3. **Ambil data user** untuk customer data di Xendit.
4. **Untuk setiap vendor group dari item terpilih**:
   - Validasi vendor punya `XenditAccountID`.
   - Resolve alamat gudang + kurir vendor.
   - Build order items, hitung subtotal + total berat.
   - **Hitung ulang ongkir di server** via RajaOngkir (biaya dari client **TIDAK** dipercaya).
   - Match `courier_code` + `service` yang dipilih customer terhadap response RajaOngkir.
   - Hitung `grand_total = subtotal + shipping_fee + platform_fee (Rp 6.000)`.
   - Dalam satu **DB transaction**:
     - Buat record `orders` + `order_items` + `order_status_history`
     - **Kurangi `stock_on_hand` secara atomik** (gagal jika stok tidak cukup)
   - Panggil Xendit `CreateInvoice` pada sub-account vendor, masa berlaku **24 jam**.
   - Buat record `payment_invoices` (status = `pending`).
5. **Hapus item terpilih yang berhasil diproses** dari cart aktif.
6. **Biarkan item yang tidak dipilih tetap berada di cart aktif**.
7. **Return response**: Array order dengan `invoice_url`, `order_id`, `order_no`, rincian biaya, dan `expires_at`.

### Step 4 — Pembayaran

Pembayaran diproses melalui **webhook Xendit** ke `POST /api/v1/webhooks/xendit/invoice`.

- **PAID**: Invoice → `paid`, Order → `paid`, jurnal ledger (double-entry) dicatat, saldo vendor di-credit sebesar `subtotal`.
- **EXPIRED**: Invoice → `expired`, Order → `canceled`, stok di-restore secara atomik.
- Idempotency: Event di-deduplikasi berdasarkan `external_event_id` (`{xendit_id}:{status}`).

### Step 5 — Penyelesaian Order

Customer bisa menandai order sebagai selesai via `POST /api/v1/users/orders/:orderId/complete` dari status `paid`, `packed`, atau `shipped` → `completed`.

---

## API Endpoints

### Cart (Customer Only)

| Method | Path | Deskripsi |
|--------|------|-----------|
| `GET` | `/api/v1/users/cart` | Ambil cart aktif beserta items |
| `POST` | `/api/v1/users/cart/items` | Tambah product variant ke cart |
| `PATCH` | `/api/v1/users/cart/items/:itemId` | Update quantity item |
| `PATCH` | `/api/v1/users/cart/items/:itemId/selection` | Ubah status selected item untuk checkout |
| `DELETE` | `/api/v1/users/cart/items/:itemId` | Hapus satu item |
| `DELETE` | `/api/v1/users/cart` | Kosongkan seluruh cart |

### Checkout (Customer Only)

| Method | Path | Deskripsi |
|--------|------|-----------|
| `POST` | `/api/v1/users/checkout/preview` | Preview checkout untuk item cart terpilih — items per vendor + opsi ongkir |
| `POST` | `/api/v1/users/checkout` | Place order untuk item cart terpilih, buat invoice Xendit, return payment URL |

### Order (Customer Only)

| Method | Path | Deskripsi |
|--------|------|-----------|
| `GET` | `/api/v1/users/orders` | Daftar order customer (paginated, filter by status) |
| `POST` | `/api/v1/users/orders/:orderId/complete` | Tandai order selesai |
| `GET` | `/api/v1/order-statuses` | Daftar semua status order yang valid (public) |

### Webhook

| Method | Path | Deskripsi |
|--------|------|-----------|
| `POST` | `/api/v1/webhooks/xendit/invoice` | Webhook Xendit (verifikasi via `x-callback-token`) |

### Lokasi Pengiriman (Customer & Vendor)

| Method | Path | Deskripsi |
|--------|------|-----------|
| `GET` | `/api/v1/locations/provinces` | Daftar provinsi (alias resmi untuk form alamat umum/vendor) |
| `GET` | `/api/v1/locations/cities?province_id=` | Daftar kota per provinsi |
| `GET` | `/api/v1/locations/districts?city_id=` | Daftar kecamatan per kota |
| `GET` | `/api/v1/locations/subdistricts?district_id=` | Daftar kelurahan per kecamatan |
| `GET` | `/api/v1/shipping/provinces` | Alias backward-compatible untuk provinces |
| `GET` | `/api/v1/shipping/cities?province_id=` | Alias backward-compatible untuk cities |
| `GET` | `/api/v1/shipping/districts?city_id=` | Alias backward-compatible untuk districts |
| `GET` | `/api/v1/shipping/subdistricts?district_id=` | Alias backward-compatible untuk subdistricts |
| `GET` | `/api/v1/shipping/couriers` | Daftar semua kurir aktif |

### Manajemen Kurir Vendor (Vendor Only)

| Method | Path | Deskripsi |
|--------|------|-----------|
| `GET` | `/api/v1/vendors/couriers` | Ambil daftar kurir vendor |
| `PUT` | `/api/v1/vendors/couriers` | Ganti seluruh daftar kurir vendor |
| `DELETE` | `/api/v1/vendors/couriers/:courierId` | Hapus satu kurir dari vendor |

---

## Ongkos Kirim (Shipping)

### Provider

Ongkir dihitung menggunakan **RajaOngkir Komerce API** (`https://rajaongkir.komerce.id`).

### Cara Kerja Perhitungan Ongkir

```
┌──────────────────┐    ┌──────────────┐    ┌────────────────┐
│ Alamat Gudang    │    │ Alamat       │    │ Berat Total    │
│ Vendor           │    │ Customer     │    │ Barang         │
│ (district_id)    │    │ (district_id)│    │ (gram)         │
└────────┬─────────┘    └──────┬───────┘    └───────┬────────┘
         │                     │                    │
         └─────────┬───────────┘                    │
                   ▼                                │
         ┌──────────────────┐                       │
         │ Kurir yang       │                       │
         │ dipilih vendor   │◀──────────────────────┘
         │ (vendor_couriers)│
         └────────┬─────────┘
                  ▼
         ┌──────────────────┐
         │  RajaOngkir API  │
         │  (domestic-cost) │
         └────────┬─────────┘
                  ▼
         ┌──────────────────┐
         │ Daftar Opsi      │
         │ Ongkir           │
         │ (kurir, service,  │
         │  biaya, ETD)     │
         └──────────────────┘
```

### Parameter Perhitungan

| Parameter | Sumber | Keterangan |
|-----------|--------|------------|
| **Origin** | `district_id` dari alamat default pemilik vendor | Alamat gudang/asal pengiriman |
| **Destination** | `district_id` dari alamat pengiriman customer | Alamat tujuan pengiriman |
| **Weight** | Jumlah `weight_gram` × `qty` per item | Default **500 gram** per item jika berat tidak diset |
| **Couriers** | Kode kurir dari `vendor_couriers` | Digabung dengan separator `:` (contoh: `jne:sicepat:jnt`) |

### Contoh Request RajaOngkir

```bash
curl --location 'https://rajaongkir.komerce.id/api/v1/calculate/district/domestic-cost' \
  --header 'key: YOUR_API_KEY' \
  --header 'Content-Type: application/x-www-form-urlencoded' \
  --data-urlencode 'origin=1391' \
  --data-urlencode 'destination=1376' \
  --data-urlencode 'weight=1000' \
  --data-urlencode 'courier=jne:sicepat:jnt' \
  --data-urlencode 'price=lowest'
```

### Keamanan Ongkir

> **Penting**: Ongkir yang ditampilkan di preview bersifat **informatif**. Saat checkout (place order), ongkir **selalu dihitung ulang di server**. Client hanya mengirim `courier_code` + `service`, **bukan harga ongkir**. Ini mencegah manipulasi biaya pengiriman.

### Kurir yang Didukung

Terdapat 12 kurir domestik yang tersedia:

| Kode | Nama |
|------|------|
| `jne` | JNE |
| `sicepat` | SiCepat |
| `ide` | IDExpress |
| `sap` | SAP Express |
| `ninja` | Ninja |
| `jnt` | J&T Express |
| `tiki` | TIKI |
| `wahana` | Wahana Express |
| `pos` | POS Indonesia |
| `sentral` | Sentral Cargo |
| `lion` | Lion Parcel |
| `rex` | Royal Express Asia |

Vendor memilih kurir mana saja yang mereka gunakan melalui endpoint manajemen kurir. Hanya kurir yang dipilih vendor yang akan muncul di opsi ongkir saat checkout.

### Alur Vendor Memilih Kurir

1. Admin/Vendor melihat daftar kurir tersedia: `GET /api/v1/shipping/couriers`
2. Vendor mengatur kurir yang dipakai: `PUT /api/v1/vendors/couriers`
3. Saat checkout, hanya kurir yang dipilih vendor yang dikirim ke RajaOngkir

---

## Biaya Platform

Setiap order dikenakan **platform fee tetap sebesar Rp 6.000**, terdiri dari:

| Komponen | Jumlah |
|----------|--------|
| Admin Fee | Rp 5.000 |
| App Fee | Rp 1.000 |
| **Total Platform Fee** | **Rp 6.000** |

### Rumus Grand Total

```
grand_total = subtotal + shipping_fee + platform_fee
```

Contoh:
- Subtotal barang: Rp 200.000
- Ongkir (JNE REG): Rp 15.000
- Platform fee: Rp 6.000
- **Grand total: Rp 221.000**

---

## Pembayaran (Xendit)

### Arsitektur Multi-Vendor

Setiap order menghasilkan **invoice Xendit yang terpisah** pada **sub-account vendor** (menggunakan XenPlatform `for-user-id` header). Artinya pembayaran langsung masuk ke akun Xendit masing-masing vendor.

### Detail Invoice

| Parameter | Nilai |
|-----------|-------|
| Gateway | Xendit Invoice API |
| Masa Berlaku | 24 jam (86.400 detik) |
| Format External ID | `INV-{order_no}-{8 karakter random}` |
| Redirect Success | `{frontend_url}/payment/success?order_id={id}` |
| Redirect Failure | `{frontend_url}/payment/failed?order_id={id}` |

### Pencatatan Keuangan (Ledger)

Saat pembayaran berhasil (PAID), sistem mencatat jurnal double-entry:

| Akun | Debit | Credit |
|------|-------|--------|
| `1100` Payment Gateway Receivable | paid_amount | — |
| `2100` Vendor Payable | — | subtotal |
| `4100` Platform Fee Revenue | — | Rp 1.000 |
| `4200` Admin Fee Revenue | — | Rp 5.000 |

Saldo vendor (`vendor.balance`) di-credit sebesar **subtotal** order untuk mekanisme withdrawal mandiri.

### Verifikasi Webhook

Header `x-callback-token` harus sesuai dengan `webhookVerificationToken` yang dikonfigurasi, untuk memastikan request benar-benar berasal dari Xendit.

---

## Status Order & Transisi

### Daftar Status

| Status | Deskripsi |
|--------|-----------|
| `pending_payment` | Status awal setelah order dibuat, menunggu pembayaran |
| `paid` | Pembayaran dikonfirmasi via webhook Xendit |
| `packed` | Vendor sudah mengemas pesanan |
| `shipped` | Pengiriman dikirim (nomor resi sudah ada) |
| `completed` | Customer konfirmasi terima barang |
| `canceled` | Pembayaran expired atau rollback checkout |
| `refunded` | Pembayaran telah di-refund |

### Diagram Transisi Status

```
                      ┌────────────────────────────────────────────┐
                      │                                            │
                      ▼                                            │
┌─────────────────┐  PAID   ┌──────┐  pack  ┌────────┐  ship  ┌─────────┐  confirm ┌───────────┐
│ pending_payment │────────▶│ paid │───────▶│ packed │───────▶│ shipped │────────▶│ completed │
└────────┬────────┘         └──┬───┘        └────┬───┘        └────┬────┘         └───────────┘
         │                     │                 │                 │
         │  EXPIRED            │    complete      │    complete     │    complete
         ▼                     └─────────────────┴─────────────────┘
┌──────────┐                                     │
│ canceled │◀────────────────────────────────────┘ (customer bisa complete
└──────────┘                                        dari paid/packed/shipped)
```

### Status Pengiriman (Shipment)

| Status | Deskripsi |
|--------|-----------|
| `waiting_pickup` | Menunggu pick-up oleh kurir |
| `shipped` | Sudah dikirim |
| `delivered` | Sudah sampai tujuan |
| `failed` | Pengiriman gagal |
| `returned` | Barang dikembalikan |

---

## Request & Response DTO

### Checkout Preview

**Request** — `POST /api/v1/users/checkout/preview`

```json
{
  "address_id": "uuid-alamat-customer"
}
```

**Response**

```json
{
  "success": true,
  "message": "Checkout preview",
  "data": {
    "address": {
      "id": "uuid",
      "label": "Rumah",
      "recipient_name": "Ahmad",
      "phone": "08123456789",
      "full_address": "Jl. Mawar No. 10",
      "district": "Kebayoran Baru",
      "city": "Jakarta Selatan",
      "province": "DKI Jakarta"
    },
    "vendors": [
      {
        "vendor_id": "uuid-vendor",
        "vendor_name": "Toko Perlengkapan Haji",
        "items": [
          {
            "product_variant_id": "uuid-variant",
            "product_name": "Kain Ihram Premium",
            "variant_name": "Putih - L",
            "image_url": "https://s3.../presigned-url",
            "price": 100000,
            "qty": 2,
            "subtotal": 200000,
            "weight_gram": 500
          }
        ],
        "subtotal": 200000,
        "total_weight_gram": 1000,
        "shipping_options": [
          {
            "name": "JNE",
            "code": "jne",
            "service": "REG",
            "description": "Layanan Reguler",
            "cost": 15000,
            "etd": "2-3"
          },
          {
            "name": "SiCepat",
            "code": "sicepat",
            "service": "BEST",
            "description": "Belanja Express Sampai Tepat Waktu",
            "cost": 13000,
            "etd": "1-2"
          }
        ]
      }
    ],
    "platform_fee": 6000
  }
}
```

Catatan:
- Preview hanya memproses item dengan `is_selected = true`.
- Jika tidak ada item cart yang dipilih, endpoint akan gagal dengan error `no selected cart items`.

### Place Order (Checkout)

**Request** — `POST /api/v1/users/checkout`

```json
{
  "address_id": "uuid-alamat-customer",
  "shipping_choices": [
    {
      "vendor_id": "uuid-vendor",
      "courier_code": "jne",
      "service": "REG"
    }
  ],
  "notes": "Tolong packing rapi ya"
}
```

Catatan:
- Checkout hanya memproses item dengan `is_selected = true`.
- Setelah checkout sukses, item yang berhasil diproses dihapus dari cart aktif.
- Item yang tidak dipilih tetap berada di cart aktif.
- Field `notes` saat ini diterima oleh API tetapi belum diproses lebih lanjut oleh backend.

**Response**

```json
{
  "success": true,
  "message": "Checkout successful",
  "data": {
    "orders": [
      {
        "order_id": "uuid-order",
        "order_no": "ORD-20260311-ABCD1234",
        "vendor_id": "uuid-vendor",
        "vendor_name": "Toko Perlengkapan Haji",
        "subtotal": 200000,
        "shipping_fee": 15000,
        "platform_fee": 6000,
        "grand_total": 221000,
        "invoice_url": "https://checkout.xendit.co/web/...",
        "expires_at": "2026-03-12T12:00:00Z"
      }
    ]
  }
}
```

### Xendit Webhook Payload

**Request** — `POST /api/v1/webhooks/xendit/invoice`

```json
{
  "id": "xendit-invoice-id-123",
  "external_id": "INV-ORD-20260311-ABCD1234-XY12ZW34",
  "status": "PAID",
  "amount": 221000,
  "paid_amount": 221000,
  "paid_at": "2026-03-11T13:00:00Z",
  "payment_method": "BANK_TRANSFER",
  "payment_channel": "BCA",
  "currency": "IDR"
}
```

---

## Skema Database

### Tabel `orders`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID (PK) | |
| `order_no` | VARCHAR(30) UNIQUE | Format: `ORD-YYYYMMDD-XXXXXXXX` |
| `user_id` | UUID FK → users | Customer pemesan |
| `vendor_id` | UUID FK → vendors | Satu order per vendor per checkout |
| `shipping_address_snapshot` | JSONB | Salinan alamat saat order dibuat |
| `order_status` | VARCHAR(24) | Default `pending_payment` |
| `payment_status` | VARCHAR(16) | Default `unpaid` |
| `subtotal` | NUMERIC(18,2) | Total harga barang |
| `shipping_fee` | NUMERIC(18,2) | Ongkos kirim |
| `platform_fee` | NUMERIC(18,2) | Biaya platform (Rp 6.000) |
| `grand_total` | NUMERIC(18,2) | subtotal + shipping + platform |
| `placed_at` | TIMESTAMPTZ | Waktu order dibuat |
| `created_at` / `updated_at` | TIMESTAMPTZ | Timestamp standar |

### Tabel `order_items`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID (PK) | |
| `order_id` | UUID FK → orders | |
| `product_variant_id` | UUID FK → product_variants | |
| `product_name_snapshot` | VARCHAR(180) | Nama produk saat order |
| `sku_snapshot` | VARCHAR(80) | SKU saat order |
| `qty` | INT | Jumlah |
| `unit_price` | NUMERIC(18,2) | Harga per unit |
| `line_total` | NUMERIC(18,2) | qty × unit_price |

### Tabel `order_status_history`

Mencatat setiap perubahan status: `old_status`, `new_status`, `changed_by`, `changed_at`, `notes`.

### Tabel `shipments`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID (PK) | |
| `order_id` | UUID UNIQUE FK → orders | 1:1 dengan order |
| `courier_code` | VARCHAR(30) | Kode kurir |
| `service_type` | VARCHAR(30) | Tipe layanan |
| `tracking_no` | VARCHAR(80) | Nomor resi |
| `shipment_status` | VARCHAR(24) | Status pengiriman |
| `shipped_at` / `delivered_at` | TIMESTAMPTZ | Waktu kirim/terima |

### Tabel `payment_invoices`

| Kolom | Tipe | Keterangan |
|-------|------|------------|
| `id` | UUID (PK) | |
| `order_id` | UUID FK → orders | |
| `gateway` | VARCHAR(20) | Selalu `xendit` |
| `xendit_invoice_id` | VARCHAR(100) | ID invoice dari Xendit |
| `external_invoice_id` | VARCHAR(100) UNIQUE | Format: `INV-{order_no}-{8chars}` |
| `invoice_url` | TEXT | URL halaman checkout Xendit |
| `status` | VARCHAR(16) | `pending` / `paid` / `expired` / `failed` |
| `amount` | NUMERIC(18,2) | = grand_total |
| `expires_at` / `paid_at` | TIMESTAMPTZ | |
| `payment_method` / `payment_channel` | VARCHAR(30) | Diisi saat PAID |
| `raw_payload` | JSONB | Raw webhook payload |

### Tabel `payment_events`

Tabel idempotency: `external_event_id` UNIQUE, menyimpan setiap event webhook Xendit.

### Tabel `carts` / `cart_items`

- `carts`: Satu cart aktif per user.
- `cart_items`: UNIQUE `cart_id` + `product_variant_id`, menyimpan `qty` dan `is_selected`.
- `is_selected`: Menentukan apakah item ikut dalam checkout preview dan checkout.
- Partial checkout tidak mengubah status cart; hanya item terpilih yang sukses diproses akan dihapus dari cart aktif.

### Tabel `couriers` / `vendor_couriers`

- `couriers`: Tabel master kurir (serial PK, `code`, `name`, `logo_url`, `is_active`).
- `vendor_couriers`: Tabel join (`vendor_id`, `courier_id`) UNIQUE, `is_active`.

---

## Mekanisme Kompensasi (Rollback)

Jika terjadi kegagalan di tengah proses checkout (misal: checkout multi-vendor, vendor pertama berhasil tapi vendor kedua gagal), sistem akan melakukan **kompensasi otomatis** untuk order yang sudah terbuat:

1. **Invoice Xendit** yang sudah dibuat akan di-expire-kan.
2. **Payment invoice** di DB ditandai status `failed`.
3. **Order** diubah ke status `canceled`.
4. **Stok** yang sudah dikurangi akan **di-restore** kembali.

Ini memastikan konsistensi data meskipun terjadi partial failure.

---

_Terakhir diperbarui: 11 Maret 2026_
