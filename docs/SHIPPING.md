# Shipping Integration Design (2026)

## Summary
This project now uses a simplified shipping integration for RajaOngkir. The previous `shipping_services` and `product_shipping_services` tables have been **removed**. All shipping logic is now based on vendor-selected couriers only.

## Key Changes
- **Removed Tables:**
  - `shipping_services` (UUID, name, code)
  - `product_shipping_services` (product_id, shipping_service_id)
- **Removed Code:**
  - All domain, repository, usecase, handler, and router logic for shipping services
  - All API endpoints for `/shipping-services`
- **Migration:**
  - See `migrations/000029_drop_shipping_services.up.sql` for schema changes

## New Shipping Flow
- Vendors select couriers via `vendor_couriers` table:
  - `vendor_id`, `courier_id` (int, references `couriers`)
- Checkout and shipping cost calculation uses only the couriers enabled by each vendor
- RajaOngkir API is called with a colon-separated list of courier codes (e.g. `jne:sicepat:jnt`)
- No more product-level shipping service validation

## Benefits
- **No data redundancy:** Only couriers are stored, no duplicate service info
- **Simpler maintenance:** No need to sync service codes with RajaOngkir
- **No mismatch risk:** All courier codes come directly from vendor selection

## Example RajaOngkir Request
```
curl --location 'https://rajaongkir.komerce.id/api/v1/calculate/district/domestic-cost' \
--header 'key: YOUR_API_KEY' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'origin=1391' \
--data-urlencode 'destination=1376' \
--data-urlencode 'weight=1000' \
--data-urlencode 'courier=jne:sicepat:jnt' \
--data-urlencode 'price=lowest'
```

## Migration Notes
- Run migration `000029_drop_shipping_services.up.sql` to remove old tables
- Remove all code references to shipping services
- Update seed SQL to remove shipping_services and product_shipping_services

## Table References
- **Couriers:** `couriers` (int id, code, name)
- **Vendor Couriers:** `vendor_couriers` (vendor_id, courier_id)

## API Impact
- `/api/v1/shipping-services` endpoint is **removed**
- Product creation no longer requires shipping service selection
- All shipping logic is now vendor/courier-based

---
_Last updated: March 6, 2026_
