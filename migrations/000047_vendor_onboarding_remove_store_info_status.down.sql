ALTER TABLE vendor_onboardings
    DROP CONSTRAINT IF EXISTS ck_vendor_onboardings_status;

ALTER TABLE vendor_onboardings
    ADD CONSTRAINT ck_vendor_onboardings_status CHECK (
        status IN ('otp_verified', 'password_set', 'store_info_completed', 'completed')
    );

