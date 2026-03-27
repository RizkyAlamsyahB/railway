UPDATE vendor_onboardings
SET status = 'password_set', updated_at = now()
WHERE status = 'store_info_completed';

ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_status;

ALTER TABLE vendor_onboardings
    ADD CONSTRAINT ck_vendor_onboardings_status CHECK (
        status IN ('otp_verified', 'password_set', 'completed')
    );

