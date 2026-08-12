ALTER TABLE ftn_identities
    ADD COLUMN IF NOT EXISTS password_salt BYTEA;

ALTER TABLE ftn_identities
    ADD COLUMN IF NOT EXISTS password_version SMALLINT NOT NULL DEFAULT 1;

CREATE INDEX IF NOT EXISTS idx_ftn_identities_status
    ON ftn_identities(status);

-- Production provisioning must populate password_salt and password_hash together.
-- Existing rows with NULL password_salt must not be enabled for password login.
