BEGIN;

-- Auto‑received: order masih shipped + delivered_at > 48 jam
UPDATE shipments s
SET delivered_at = NOW() - interval '49 hours',
    shipped_at   = NOW() - interval '60 hours'
FROM orders o
WHERE o.order_no = 'SEED-SHP-0001'
  AND s.order_id = o.id;

UPDATE orders
SET order_status   = 'shipped',
    payment_status = 'paid',
    updated_at     = NOW() - interval '2 hours'
WHERE order_no = 'SEED-SHP-0001';

-- Settlement: order sudah received + updated_at > 24 jam
UPDATE orders
SET order_status   = 'received',
    payment_status = 'paid',
    updated_at     = NOW() - interval '25 hours'
WHERE order_no = 'SEED-REC-0001';

-- Pastikan escrow untuk order received ada (gunakan subtotal order)
-- Jalankan sekali agar tidak double-assign.
UPDATE vendor_balances vb
SET escrow_balance = vb.escrow_balance + o.subtotal,
    total_earned   = vb.total_earned + o.subtotal,
    updated_at     = NOW()
FROM orders o
WHERE o.order_no = 'SEED-REC-0001'
  AND vb.vendor_id = o.vendor_id;

COMMIT;
