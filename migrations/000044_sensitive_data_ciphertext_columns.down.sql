ALTER TABLE vendor_bank_accounts
    ALTER COLUMN account_number TYPE VARCHAR(60);

ALTER TABLE vendor_onboardings
    ALTER COLUMN nik TYPE VARCHAR(32);
