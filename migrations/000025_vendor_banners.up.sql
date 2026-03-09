-- Vendor banners for vendor storefront
CREATE TABLE vendor_banners (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id   UUID        NOT NULL REFERENCES vendors(id) ON DELETE CASCADE,
    title       VARCHAR(255) NOT NULL,
    image_url   TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_vendor_banners_vendor_id ON vendor_banners (vendor_id);
