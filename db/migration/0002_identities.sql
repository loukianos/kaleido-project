-- +goose Up
CREATE TABLE IF NOT EXISTS identities (
    id bigserial PRIMARY KEY,
    issuer text NOT NULL,
    subject text NOT NULL,
    role text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT identities_issuer_not_empty CHECK (issuer <> ''),
    CONSTRAINT identities_subject_not_empty CHECK (subject <> ''),
    CONSTRAINT identities_role_valid CHECK (role IN ('lender', 'servicer', 'admin')),
    CONSTRAINT identities_issuer_subject_unique UNIQUE (issuer, subject)
);

CREATE TABLE IF NOT EXISTS signing_keys (
    id bigserial PRIMARY KEY,
    identity_id bigint NOT NULL UNIQUE REFERENCES identities(id),
    address text NOT NULL UNIQUE,
    ciphertext bytea NOT NULL,
    encryptor text NOT NULL,
    key_version integer NOT NULL DEFAULT 1,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT signing_keys_address_not_empty CHECK (address <> ''),
    CONSTRAINT signing_keys_encryptor_not_empty CHECK (encryptor <> '')
);

ALTER TABLE loans
    ADD COLUMN IF NOT EXISTS lender_identity_id bigint REFERENCES identities(id);

-- +goose Down
ALTER TABLE loans DROP COLUMN IF EXISTS lender_identity_id;
DROP TABLE IF EXISTS signing_keys;
DROP TABLE IF EXISTS identities;
