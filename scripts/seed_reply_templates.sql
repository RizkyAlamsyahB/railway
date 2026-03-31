-- Seed reply_templates with created_by = 789da1a5-b9e4-4fa0-aaca-6ed900e31b10
-- Run: psql -d haji_umroh_store -f scripts/seed_reply_templates.sql

INSERT INTO public.reply_templates VALUES ('c0e89a64-456f-4df8-bf77-2b8acd4b1b25', 'Salam Pembuka', 'Assalamualaikum, terima kasih telah menghubungi kami. Ada yang bisa kami bantu?', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-13 07:24:21.226459+00', '2026-03-13 07:24:21.226459+00', '/salam', 'general');
INSERT INTO public.reply_templates VALUES ('db924680-cc08-4bf0-9a4c-d54605ff89c1', 'Selamat Datang di UmrahMart', 'Assalamu''alaikum, selamat datang di UmrahMart! 🕌

Kami hadir untuk memudahkan Anda mendapatkan perlengkapan haji dan umroh terbaik. Ada yang bisa kami bantu hari ini?', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/selamat-datang', 'general');
INSERT INTO public.reply_templates VALUES ('b7ec135c-0073-43ec-a7ee-65fecb5668b8', 'Terima Kasih Telah Menghubungi Kami', 'Terima kasih telah menghubungi tim Customer Service UmrahMart. 🙏

Kami akan segera membantu Anda. Mohon tunggu sebentar ya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/terima-kasih-menghubungi', 'general');
INSERT INTO public.reply_templates VALUES ('9950affe-68af-4a87-bddd-57cb2b2d4522', 'Apakah Masih Ada Yang Bisa Dibantu?', 'Apakah masih ada yang bisa kami bantu? 😊

Jika masalah Anda sudah teratasi, jangan lupa berikan ulasan untuk membantu kami meningkatkan layanan. Semoga ibadah haji/umroh Anda berjalan lancar. Aamiin.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/ada-yang-bisa-dibantu', 'general');
INSERT INTO public.reply_templates VALUES ('2210ed2d-6243-47da-af10-cd789f2416b0', 'Jam Operasional Customer Service', 'Tim Customer Service UmrahMart melayani Anda setiap hari:

🕐 Senin – Jumat : 08.00 – 21.00 WIB
🕐 Sabtu – Minggu : 09.00 – 18.00 WIB

Di luar jam operasional, Anda tetap bisa meninggalkan pesan dan kami akan merespons saat jam kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/general/jam-operasional', 'general');
INSERT INTO public.reply_templates VALUES ('fd06b343-9d5a-49bf-ad6a-427fa4b238b0', 'Konfirmasi Pesanan Diterima', 'Alhamdulillah, pesanan Anda telah kami terima! ✅

Nomor pesanan Anda: {order_number}
Status: Menunggu Pembayaran

Silakan selesaikan pembayaran sebelum {expired_at} agar pesanan dapat segera kami proses. Jazakallahu khairan.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/konfirmasi-pesanan-diterima', 'order');
INSERT INTO public.reply_templates VALUES ('2c22e30d-7e7a-484f-ab6e-bd3b6f90dd29', 'Pesanan Sedang Diproses', 'Kabar baik! Pesanan Anda dengan nomor {order_number} sedang kami proses. 📦

Estimasi pesanan siap dikirim: 1–2 hari kerja.
Kami akan menginformasikan nomor resi pengiriman segera setelah paket diserahkan ke kurir.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/pesanan-diproses', 'order');
INSERT INTO public.reply_templates VALUES ('0dcf4408-01fb-49de-b782-03c949175719', 'Pesanan Berhasil Dibatalkan', 'Pesanan Anda dengan nomor {order_number} telah berhasil dibatalkan.

Jika Anda sudah melakukan pembayaran, proses refund akan kami lakukan dalam 3–7 hari kerja ke metode pembayaran asal.

Apabila ada pertanyaan lebih lanjut, jangan ragu untuk menghubungi kami kembali.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/pesanan-dibatalkan', 'order');
INSERT INTO public.reply_templates VALUES ('09e4155c-ccaf-4466-99e4-6f09c12682c8', 'Detail Pesanan', 'Berikut detail pesanan Anda:

📋 Nomor Pesanan : {order_number}
📅 Tanggal       : {order_date}
💰 Total         : {total_amount}
📍 Status        : {order_status}

Untuk informasi lebih lengkap, silakan cek halaman "Pesanan Saya" di aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/order/detail-pesanan', 'order');
INSERT INTO public.reply_templates VALUES ('039e2bf9-5a35-47fa-b5c6-df6550e1c4bd', 'Pembayaran Berhasil Diterima', 'Alhamdulillah, pembayaran Anda telah berhasil kami terima! 💚

Nomor pesanan : {order_number}
Jumlah        : {amount}
Metode        : {payment_method}

Pesanan Anda akan segera kami proses. Terima kasih.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/pembayaran-berhasil', 'payment');
INSERT INTO public.reply_templates VALUES ('279010e0-341c-42b1-8818-bf5e5b0d2bea', 'Menunggu Konfirmasi Pembayaran', 'Halo, kami melihat pesanan Anda belum terkonfirmasi pembayarannya.

Batas waktu pembayaran: {expired_at}

Jika Anda sudah melakukan transfer manual, mohon upload bukti pembayaran melalui aplikasi agar kami dapat segera memprosesnya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/menunggu-konfirmasi', 'payment');
INSERT INTO public.reply_templates VALUES ('c5c27b1c-05d5-4ae6-83a0-7bb1c89c74f9', 'Pembayaran Gagal / Kedaluwarsa', 'Kami informasikan bahwa pembayaran untuk pesanan {order_number} telah gagal atau melewati batas waktu. ⚠️

Pesanan Anda telah kami batalkan secara otomatis. Anda dapat melakukan pemesanan ulang kapan saja.

Mohon pastikan saldo/limit mencukupi saat melakukan pembayaran berikutnya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/pembayaran-gagal', 'payment');
INSERT INTO public.reply_templates VALUES ('98ccc880-d214-4c26-bb39-c6fc1cc4bfe3', 'Cara Melakukan Pembayaran', 'Berikut metode pembayaran yang tersedia di UmrahMart:

💳 Transfer Bank (BCA, Mandiri, BNI, BRI)
📱 E-Wallet (GoPay, OVO, DANA, ShopeePay)
🏪 Gerai Retail (Alfamart, Indomaret)
💵 QRIS

Pilih metode yang paling nyaman untuk Anda saat checkout.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/payment/cara-pembayaran', 'payment');
INSERT INTO public.reply_templates VALUES ('022b8cc0-8129-45d4-a9b5-0601ce941461', 'Pengajuan Refund Diterima', 'Pengajuan refund Anda telah kami terima dan sedang dalam proses review. ✅

Nomor tiket refund: {ticket_number}

Tim kami akan memverifikasi dalam 1–3 hari kerja. Anda akan mendapat notifikasi setelah proses selesai.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/pengajuan-diterima', 'refund');
INSERT INTO public.reply_templates VALUES ('60cf1bc4-50e9-49ca-a6fa-dc79ac495678', 'Refund Sedang Diproses', 'Refund Anda sedang dalam proses pencairan. 🔄

Jumlah refund : {refund_amount}
Tujuan        : {refund_destination}
Estimasi      : 3–7 hari kerja

Mohon bersabar dan pastikan rekening/e-wallet tujuan masih aktif.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/sedang-diproses', 'refund');
INSERT INTO public.reply_templates VALUES ('83ec97f6-b6fc-4e70-9491-ac41192fe241', 'Refund Berhasil', 'Dana refund sebesar {refund_amount} telah berhasil kami kirimkan ke {refund_destination}. 🎉

Jika dalam 1×24 jam dana belum masuk, mohon hubungi kami kembali dengan menyertakan nomor tiket {ticket_number}.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/refund/berhasil', 'refund');
INSERT INTO public.reply_templates VALUES ('c0aced63-92c1-4857-8db0-91187f1822ad', 'Pesanan Telah Dikirim', 'Pesanan Anda sudah dalam perjalanan! 🚚

Nomor Pesanan : {order_number}
Kurir         : {courier_name}
Nomor Resi    : {tracking_number}

Lacak pengiriman di: {tracking_url}

Estimasi tiba: {estimated_arrival}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/pesanan-dikirim', 'shipping');
INSERT INTO public.reply_templates VALUES ('933272ef-be7b-44dd-9fb8-8eaa07c48592', 'Konfirmasi Penerimaan Paket', 'Halo, berdasarkan data kurir, paket Anda sudah dinyatakan terkirim pada {delivered_at}.

Apakah paket sudah Anda terima dengan kondisi baik? Mohon konfirmasi agar pesanan dapat diselesaikan. 📦✅', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/konfirmasi-penerimaan', 'shipping');
INSERT INTO public.reply_templates VALUES ('57dc278b-9cab-43bf-8836-901e1b4fe9cf', 'Paket Tertahan / Terlambat', 'Kami melihat terdapat kendala pada pengiriman paket Anda (no. resi: {tracking_number}). ⚠️

Tim kami sedang berkoordinasi dengan pihak kurir untuk menindaklanjuti hal ini. Kami akan segera menginformasikan perkembangannya. Mohon maaf atas ketidaknyamanannya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/shipping/paket-tertahan', 'shipping');
INSERT INTO public.reply_templates VALUES ('0c4d32dc-44fc-402b-9fa5-44a181768d95', 'Informasi Ketersediaan Stok', 'Terima kasih atas minat Anda pada produk {product_name}.

Saat ini stok produk tersebut {stock_status}.

Anda dapat mengaktifkan notifikasi "Ingatkan Saya" pada halaman produk agar kami bisa memberitahu Anda saat stok tersedia kembali.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/ketersediaan-stok', 'product');
INSERT INTO public.reply_templates VALUES ('f2db8965-96ee-4957-a034-9bafb1981291', 'Detail Spesifikasi Produk', 'Berikut informasi lengkap mengenai {product_name}:

📌 Bahan    : {material}
📐 Ukuran   : {size}
🎨 Warna    : {color}
🏷️ Berat    : {weight}
✅ Halal    : {halal_status}

Apakah ada pertanyaan lain mengenai produk ini?', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/spesifikasi-produk', 'product');
INSERT INTO public.reply_templates VALUES ('d1a8034b-7a77-47f1-b549-440ff1a13d66', 'Verifikasi Berhasil', 'Selamat! Verifikasi akun Anda telah berhasil. ✅

Akun Anda kini sudah terverifikasi dan dapat menikmati semua fitur UmrahMart tanpa batasan. Terima kasih atas kerja samanya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/verification/berhasil', 'verification');
INSERT INTO public.reply_templates VALUES ('31b16e1b-998a-4779-87a2-663ef09cffb4', 'Produk Tidak Lagi Tersedia', 'Mohon maaf, produk {product_name} saat ini sudah tidak tersedia di katalog kami. 😔

Kami merekomendasikan produk serupa yang mungkin sesuai dengan kebutuhan Anda:
👉 {alternative_product}

Silakan cek koleksi terbaru kami di aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/product/produk-tidak-tersedia', 'product');
INSERT INTO public.reply_templates VALUES ('9d2a3755-e301-4552-a55f-d1ca99ee9aab', 'Bantuan Reset Password', 'Untuk mereset password akun Anda, ikuti langkah berikut:

1️⃣ Buka halaman Login
2️⃣ Klik "Lupa Password?"
3️⃣ Masukkan email terdaftar Anda
4️⃣ Cek inbox email untuk link reset (cek juga folder Spam)
5️⃣ Klik link dan buat password baru

Link reset berlaku selama 60 menit. Jika belum menerima email, hubungi kami kembali.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/reset-password', 'account');
INSERT INTO public.reply_templates VALUES ('52d1c3ad-e619-47c3-930c-2497057c746c', 'Verifikasi Email', 'Untuk memverifikasi email Anda:

1️⃣ Cek inbox email {email} (termasuk folder Spam/Junk)
2️⃣ Buka email dari UmrahMart dengan subjek "Verifikasi Akun"
3️⃣ Klik tombol "Verifikasi Sekarang"

Jika email tidak ditemukan, klik "Kirim Ulang Verifikasi" di halaman akun Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/verifikasi-email', 'account');
INSERT INTO public.reply_templates VALUES ('849c3f8f-05e6-4ad7-b8a7-9b39914b99ed', 'Akun Diblokir / Dinonaktifkan', 'Kami melihat akun Anda saat ini tidak aktif. 🔒

Hal ini dapat terjadi karena:
• Aktivitas yang mencurigakan terdeteksi
• Pelanggaran syarat & ketentuan penggunaan
• Permintaan penonaktifan sebelumnya

Untuk mengajukan reaktivasi, mohon kirimkan data diri Anda ke email support@umrahmart.id.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/account/akun-diblokir', 'account');
INSERT INTO public.reply_templates VALUES ('c729c57e-34d3-475d-8d65-7d0ff314a51d', 'Keluhan Diterima dan Dicatat', 'Kami sangat menyesal mendengar pengalaman yang tidak menyenangkan ini. 🙏

Keluhan Anda telah kami catat dengan nomor tiket {ticket_number}. Tim kami akan menginvestigasi dan menghubungi Anda dalam 1×24 jam kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/keluhan-diterima', 'complaint');
INSERT INTO public.reply_templates VALUES ('08411013-5cb3-497d-87a5-14b18b27361d', 'Keluhan Sedang Diinvestigasi', 'Kami tengah menginvestigasi keluhan yang Anda sampaikan terkait {complaint_subject}.

Proses investigasi membutuhkan waktu maksimal 3 hari kerja. Kami mohon kesabaran Anda dan akan segera memberikan informasi lebih lanjut.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/sedang-investigasi', 'complaint');
INSERT INTO public.reply_templates VALUES ('2be83748-504e-4be7-8895-47ef175d1689', 'Keluhan Telah Diselesaikan', 'Keluhan Anda dengan nomor tiket {ticket_number} telah kami selesaikan. ✅

Solusi yang diberikan: {resolution}

Kami memohon maaf atas pengalaman yang kurang menyenangkan ini dan berterima kasih atas masukan Anda untuk perbaikan layanan kami.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/complaint/telah-diselesaikan', 'complaint');
INSERT INTO public.reply_templates VALUES ('29ef4e0e-fff3-4d1d-998f-b4b8c16956f6', 'Prosedur Pengembalian Barang', 'Berikut prosedur pengembalian barang (return) di UmrahMart:

1️⃣ Ajukan return melalui menu "Pesanan Saya" → "Ajukan Return"
2️⃣ Pilih produk dan alasan pengembalian
3️⃣ Unggah foto kondisi barang
4️⃣ Tunggu konfirmasi persetujuan (1–2 hari kerja)
5️⃣ Kirim barang ke alamat yang tertera

Syarat: barang dalam kondisi asli, belum dipakai, dan beserta kemasan lengkap. Maksimal 7 hari setelah diterima.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/prosedur-pengembalian', 'return');
INSERT INTO public.reply_templates VALUES ('8c893df1-0794-4a42-99d1-dbc94ca8c470', 'Pengajuan Return Disetujui', 'Kabar baik! Pengajuan return Anda untuk pesanan {order_number} telah disetujui. ✅

Silakan kirimkan barang ke:
📍 {return_address}
Atas nama: UmrahMart Returns

Setelah barang kami terima dan verifikasi, dana akan dikembalikan dalam 3–5 hari kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/disetujui', 'return');
INSERT INTO public.reply_templates VALUES ('a4203d8e-4615-4ef0-b3eb-1a90bc79924e', 'Pengajuan Return Ditolak', 'Mohon maaf, pengajuan return untuk pesanan {order_number} tidak dapat kami proses. ❌

Alasan: {rejection_reason}

Jika Anda merasa keberatan, Anda dapat mengajukan banding dengan menghubungi kami dan melampirkan bukti pendukung.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/return/ditolak', 'return');
INSERT INTO public.reply_templates VALUES ('81873554-d801-46b4-8c20-a68f0aa0b604', 'Informasi Promo Aktif', 'Halo! Berikut promo yang sedang berlangsung di UmrahMart 🎉:

🏷️ {promo_name}
💰 Diskon hingga {discount_value}
📅 Berlaku: {promo_start} – {promo_end}
📝 Syarat: {promo_terms}

Jangan sampai terlewat! Belanja sekarang sebelum promo berakhir.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/promo/info-promo-aktif', 'promo');
INSERT INTO public.reply_templates VALUES ('6dcdda2a-39dc-4762-b8f0-754619e8da71', 'Promo Tidak Berlaku untuk Pesanan Ini', 'Mohon maaf, promo yang Anda gunakan tidak berlaku untuk pesanan ini. ⚠️

Kemungkinan penyebab:
• Promo sudah berakhir
• Produk tidak termasuk dalam kategori promo
• Minimum pembelian tidak terpenuhi

Cek syarat & ketentuan promo di halaman "Promo" pada aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/promo/tidak-berlaku', 'promo');
INSERT INTO public.reply_templates VALUES ('716a7e0f-9c10-404c-8fe1-a019a38adc7d', 'Cara Menggunakan Voucher', 'Berikut cara menggunakan voucher di UmrahMart:

1️⃣ Tambahkan produk ke keranjang
2️⃣ Masuk ke halaman Checkout
3️⃣ Klik "Gunakan Voucher"
4️⃣ Masukkan kode voucher: {voucher_code}
5️⃣ Klik "Terapkan" — diskon akan langsung terpotong

Pastikan pesanan Anda memenuhi minimum pembelian voucher.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/voucher/cara-menggunakan', 'voucher');
INSERT INTO public.reply_templates VALUES ('56959391-4b07-4cb8-89ab-a61abbef4748', 'Voucher Tidak Valid atau Kedaluwarsa', 'Kode voucher yang Anda masukkan tidak dapat digunakan. ❌

Kemungkinan penyebab:
• Kode voucher salah atau sudah digunakan
• Voucher sudah kedaluwarsa ({expired_date})
• Akun Anda tidak memenuhi syarat voucher

Untuk bantuan lebih lanjut, silakan kirimkan screenshot kode voucher Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/voucher/tidak-valid', 'voucher');
INSERT INTO public.reply_templates VALUES ('7ca45115-0adc-4c23-a1a1-5924e6177e20', 'Masalah Teknis pada Aplikasi', 'Kami mohon maaf atas kendala teknis yang Anda alami. 🛠️

Beberapa langkah yang dapat dicoba:
1. Tutup dan buka kembali aplikasi
2. Pastikan koneksi internet stabil
3. Update aplikasi ke versi terbaru
4. Hapus cache aplikasi
5. Restart perangkat

Jika masalah berlanjut, mohon kirimkan screenshot error beserta tipe perangkat dan versi OS Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/masalah-teknis-aplikasi', 'technical');
INSERT INTO public.reply_templates VALUES ('c139c11f-0159-4e62-b765-58dcc55b58eb', 'Tidak Bisa Login ke Akun', 'Kami memahami betapa frustrasinya tidak bisa masuk ke akun. Mari kami bantu! 🔑

Silakan coba langkah berikut:
1. Pastikan email dan password sudah benar
2. Aktifkan Caps Lock/periksa huruf kapital
3. Gunakan fitur "Lupa Password" untuk reset
4. Coba login dari perangkat atau browser berbeda

Apakah Anda menerima pesan error tertentu? Mohon informasikan agar kami dapat membantu lebih lanjut.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/tidak-bisa-login', 'technical');
INSERT INTO public.reply_templates VALUES ('82f514a7-153d-4f59-9327-229778c82585', 'Fitur Sedang Dalam Perbaikan', 'Kami menginformasikan bahwa fitur {feature_name} saat ini sedang dalam perbaikan/maintenance. 🔧

Estimasi selesai: {maintenance_end}

Kami mohon maaf atas ketidaknyamanannya. Tim teknis kami sedang bekerja keras untuk memulihkan layanan secepatnya.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/technical/fitur-dalam-perbaikan', 'technical');
INSERT INTO public.reply_templates VALUES ('774c04c8-9d21-4d1a-8c86-3ab0ef2987a0', 'Verifikasi Identitas Diperlukan', 'Untuk keamanan akun Anda, kami perlu melakukan verifikasi identitas. 🔐

Dokumen yang diperlukan:
📄 KTP (foto depan, jelas dan tidak blur)
🤳 Selfie sambil memegang KTP

Kirimkan dokumen melalui email ke: verify@umrahmart.id
Subject: Verifikasi Identitas - {user_id}

Proses verifikasi membutuhkan 1–2 hari kerja.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/verification/identitas-diperlukan', 'verification');
INSERT INTO public.reply_templates VALUES ('967434ba-1ba8-459e-9514-f8c5dfdc3ec1', 'Informasi Pendaftaran Vendor', 'Terima kasih atas minat Anda untuk bergabung sebagai vendor di UmrahMart! 🏪

Persyaratan pendaftaran:
✅ KTP pemilik usaha
✅ Foto toko/tempat usaha
✅ Logo bisnis
✅ Banner bisnis
✅ Buku rekening bank
✅ NPWP (opsional)

Proses review membutuhkan 2–3 hari kerja setelah semua dokumen lengkap.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/info-pendaftaran', 'vendor');
INSERT INTO public.reply_templates VALUES ('96cb369d-1535-4e9f-b794-2db6a703901f', 'Status Vendor Sedang Direview', 'Pendaftaran vendor Anda sedang dalam proses review oleh tim kami. ⏳

Kami akan menginformasikan hasilnya melalui email dan notifikasi aplikasi dalam 2–3 hari kerja.

Pastikan data dan dokumen yang Anda unggah sudah lengkap dan terbaca dengan jelas.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/sedang-direview', 'vendor');
INSERT INTO public.reply_templates VALUES ('bf62a18f-4d7d-48f0-99a3-fa1fae0cf4fc', 'Vendor Berhasil Diaktifkan', 'Selamat! Akun vendor Anda telah berhasil diaktifkan! 🎉

Anda sekarang dapat mulai:
🛍️ Menambahkan produk ke katalog
📊 Mengelola stok dan harga
📦 Menerima dan memproses pesanan

Untuk panduan lengkap, kunjungi halaman "Panduan Vendor" di dashboard Anda. Semoga sukses!', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/vendor/berhasil-diaktifkan', 'vendor');
INSERT INTO public.reply_templates VALUES ('df27202f-3ef0-4ac9-84ea-e14506b78e01', 'Notifikasi Stok Tersedia', 'Kabar gembira! 🎉 Produk {product_name} yang Anda tunggu-tunggu kini sudah tersedia kembali.

Stok tersisa: {remaining_stock} pcs

Segera pesan sebelum kehabisan! Klik di sini: {product_url}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/stock/stok-tersedia', 'stock');
INSERT INTO public.reply_templates VALUES ('30cf40bb-4b5a-42d1-a288-7b1b3ff5ceb5', 'Informasi Pre-Order', 'Produk {product_name} saat ini tersedia dalam mode Pre-Order. 📋

Detail Pre-Order:
📅 Batas pemesanan : {po_end_date}
🚚 Estimasi kirim  : {estimated_delivery}
💰 Harga PO        : {po_price}

Pesan sekarang untuk mendapatkan kepastian stok!', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/stock/info-pre-order', 'stock');
INSERT INTO public.reply_templates VALUES ('24e59108-be8b-4641-85f0-8f11128cf74c', 'Konfirmasi Pembatalan Pesanan', 'Kami telah menerima permintaan pembatalan untuk pesanan {order_number}. ℹ️

Untuk melanjutkan proses pembatalan, mohon konfirmasi alasan Anda:
a) Salah alamat pengiriman
b) Ingin ganti produk
c) Menemukan harga lebih murah
d) Berubah pikiran
e) Lainnya

Balas dengan huruf pilihan Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/konfirmasi-pembatalan', 'cancellation');
INSERT INTO public.reply_templates VALUES ('7ffccfea-bef8-45d5-b8cb-c46798a0da55', 'Pesanan Tidak Dapat Dibatalkan', 'Mohon maaf, pesanan {order_number} tidak dapat dibatalkan karena status pesanan sudah dalam tahap {order_status}. ⚠️

Jika barang sudah diterima dan terdapat masalah, Anda masih bisa mengajukan return dalam 7 hari setelah penerimaan.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/tidak-dapat-dibatalkan', 'cancellation');
INSERT INTO public.reply_templates VALUES ('f04503ae-fb8b-4e54-b17d-b2473014537f', 'Pembatalan Berhasil', 'Pesanan {order_number} berhasil dibatalkan. ✅

Jika ada pembayaran yang perlu dikembalikan:
💰 Jumlah refund : {refund_amount}
⏱️ Estimasi dana masuk : 3–7 hari kerja

Terima kasih atas pengertian Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/cancellation/berhasil', 'cancellation');
INSERT INTO public.reply_templates VALUES ('45b72f3e-52d9-499e-8d31-faa17af981be', 'Jadwal Estimasi Pengiriman', 'Berikut estimasi pengiriman berdasarkan lokasi Anda:

🏙️ Jabodetabek    : 1–2 hari kerja
🗺️ Pulau Jawa     : 2–3 hari kerja
🏝️ Luar Jawa      : 3–5 hari kerja
🗾 Papua/Terpencil : 5–10 hari kerja

Estimasi dihitung sejak paket diserahkan ke kurir, tidak termasuk hari libur nasional.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/delivery/estimasi-pengiriman', 'delivery');
INSERT INTO public.reply_templates VALUES ('ca6a9164-b296-4daf-9131-6902d1a3aa6e', 'Pengiriman Tertunda', 'Kami mohon maaf atas keterlambatan pengiriman paket Anda (no. resi: {tracking_number}). 🙏

Kendala yang terjadi: {delay_reason}

Tim kami sedang berkoordinasi dengan pihak kurir {courier_name} untuk mempercepat proses pengiriman. Kami akan menginformasikan update terbaru sesegera mungkin.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/delivery/pengiriman-tertunda', 'delivery');
INSERT INTO public.reply_templates VALUES ('a217d188-4f32-4184-80a2-cadda1107c5d', 'Informasi Pengambilan di Toko (Pickup)', 'Pesanan Anda dapat diambil langsung di:

📍 {store_address}
🕐 Jam operasional: {store_hours}

Tunjukkan nomor pesanan {order_number} atau kode QR di aplikasi kepada petugas. Pesanan akan disimpan selama 3 hari sejak notifikasi siap diambil.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/pickup/info-pickup', 'pickup');
INSERT INTO public.reply_templates VALUES ('cb92a7d1-22cc-4fb6-afa0-3b7c0140eb93', 'Pesanan Siap Diambil', 'Pesanan Anda sudah siap untuk diambil! 🛍️

Nomor Pesanan : {order_number}
Lokasi        : {store_address}
Batas Ambil   : {pickup_deadline}

Jangan lupa bawa identitas diri saat pengambilan. Sampai jumpa!', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/pickup/siap-diambil', 'pickup');
INSERT INTO public.reply_templates VALUES ('a2940fd9-f18b-4024-8694-29be71b0cf23', 'Informasi Paket Langganan', 'Berikut informasi paket langganan UmrahMart Premium:

⭐ Paket Basic   : Rp 29.000/bulan — Gratis ongkir 2x/bulan
⭐ Paket Premium : Rp 59.000/bulan — Gratis ongkir unlimited + cashback 5%
⭐ Paket VIP     : Rp 99.000/bulan — Semua benefit + akses produk eksklusif

Daftar sekarang di menu "Langganan" pada aplikasi.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/subscription/info-langganan', 'subscription');
INSERT INTO public.reply_templates VALUES ('fec12945-f26b-4e5b-892d-e53b916e5ad3', 'Masa Aktif Langganan Hampir Habis', 'Halo! Masa aktif langganan UmrahMart Premium Anda akan berakhir pada {expiry_date}. ⏰

Perpanjang sekarang dan nikmati terus benefit:
✅ Gratis ongkos kirim
✅ Cashback eksklusif member
✅ Early access produk baru

Klik di sini untuk perpanjang: {renewal_url}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/subscription/hampir-habis', 'subscription');
INSERT INTO public.reply_templates VALUES ('8ae2e394-0d0c-4fb7-8744-f13bbb7d448e', 'Permintaan Ulasan Produk', 'Assalamu''alaikum! Semoga pesanan Anda sudah diterima dengan baik. 😊

Kami akan sangat berterima kasih jika Anda meluangkan waktu untuk memberikan ulasan pada produk {product_name}.

Ulasan Anda sangat membantu calon pembeli lain dalam memilih produk yang tepat. Berikan ulasan di menu "Pesanan Saya" → "Beri Ulasan".', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/review/minta-ulasan', 'review');
INSERT INTO public.reply_templates VALUES ('011a94d2-8ce3-4594-87ba-63e9b46fe1eb', 'Terima Kasih Atas Ulasan Anda', 'Jazakallahu khairan atas ulasan yang telah Anda berikan! 🌟

Masukan Anda sangat berharga bagi kami untuk terus meningkatkan kualitas produk dan layanan. Semoga UmrahMart selalu bisa menjadi teman setia perjalanan ibadah Anda.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/review/terima-kasih-ulasan', 'review');
INSERT INTO public.reply_templates VALUES ('23b40e79-e1f3-4c17-a3ee-64ec201cd2ae', 'Laporan Penipuan Diterima', 'Laporan penipuan Anda telah kami terima dan kami tangani dengan sangat serius. 🚨

Nomor laporan: {report_number}

Tim keamanan kami akan menginvestigasi dalam 1×24 jam. Jangan lakukan transfer atau berikan data sensitif kepada pihak yang mencurigakan.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/fraud/laporan-penipuan-diterima', 'fraud');
INSERT INTO public.reply_templates VALUES ('53916ccc-b29b-40cb-a795-dc5f7de740e5', 'Tindakan Keamanan Akun', 'Kami mendeteksi aktivitas tidak biasa pada akun Anda. 🔐

Sebagai tindakan keamanan, akun Anda telah kami sementara kunci.

Langkah yang perlu Anda lakukan:
1. Segera ubah password Anda
2. Aktifkan verifikasi 2 langkah
3. Hubungi kami melalui {support_contact} untuk membuka kunci akun

Jangan pernah bagikan OTP atau password kepada siapapun, termasuk tim kami.', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/fraud/tindakan-keamanan-akun', 'fraud');
INSERT INTO public.reply_templates VALUES ('1a57d5dc-1751-43c9-bdf7-9ccfc1d9898c', 'Pertanyaan Tidak Dapat Kami Jawab Saat Ini', 'Terima kasih atas pertanyaan Anda. Pertanyaan ini membutuhkan penanganan dari tim khusus kami. 🔄

Kami akan meneruskan pertanyaan ini ke tim terkait dan akan menghubungi Anda kembali dalam 1×24 jam kerja.

Nomor tiket Anda: {ticket_number}', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/other/eskalasi-ke-tim-terkait', 'other');
INSERT INTO public.reply_templates VALUES ('9a10b61b-1dc3-407a-b5a5-b80062a404b4', 'Cara Menghubungi Kami', 'Ada beberapa cara untuk menghubungi tim UmrahMart:

💬 Live Chat  : Aplikasi UmrahMart (menu "Bantuan")
📧 Email      : support@umrahmart.id
📞 Telepon    : 021-XXXX-XXXX (08.00–21.00 WIB)
📱 WhatsApp   : wa.me/628XXXXXXXXX
📘 Instagram  : @umrahmart.id

Kami siap membantu Anda! 🙏', true, '789da1a5-b9e4-4fa0-aaca-6ed900e31b10', NULL, '2026-03-16 05:32:32.157552+00', '2026-03-16 05:32:32.157552+00', '/other/cara-menghubungi-kami', 'other');
