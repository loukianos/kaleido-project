-- +goose Up
CREATE TABLE IF NOT EXISTS contracts (
    id bigserial PRIMARY KEY,
    chain_id bigint NOT NULL,
    address text NOT NULL,
    deploy_tx_hash text,
    base_uri text NOT NULL,
    active boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT contracts_address_not_empty CHECK (address <> ''),
    CONSTRAINT contracts_base_uri_not_empty CHECK (base_uri <> '')
);

CREATE UNIQUE INDEX IF NOT EXISTS contracts_chain_address_unique
    ON contracts (chain_id, lower(address));

CREATE UNIQUE INDEX IF NOT EXISTS contracts_one_active_per_chain
    ON contracts (chain_id)
    WHERE active;

CREATE TABLE IF NOT EXISTS loans (
    id bigserial PRIMARY KEY,
    token_id text,
    contract_id bigint REFERENCES contracts(id),
    borrower_ref text NOT NULL,
    lender_address text NOT NULL,
    principal_minor bigint NOT NULL,
    apr_bps integer NOT NULL,
    term_days bigint NOT NULL,
    interest_due_minor bigint NOT NULL,
    total_due_minor bigint NOT NULL,
    outstanding_minor bigint NOT NULL,
    status text NOT NULL,
    mint_operation_id bigint,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT loans_borrower_ref_not_empty CHECK (borrower_ref <> ''),
    CONSTRAINT loans_lender_address_not_empty CHECK (lender_address <> ''),
    CONSTRAINT loans_principal_positive CHECK (principal_minor > 0),
    CONSTRAINT loans_apr_bps_range CHECK (apr_bps >= 0 AND apr_bps <= 65535),
    CONSTRAINT loans_term_positive CHECK (term_days > 0),
    CONSTRAINT loans_amounts_non_negative CHECK (
        interest_due_minor >= 0
        AND total_due_minor >= principal_minor
        AND outstanding_minor >= 0
        AND outstanding_minor <= total_due_minor
    ),
    CONSTRAINT loans_status_valid CHECK (
        status IN ('originating', 'active', 'settling', 'repaid', 'defaulted')
    ),
    CONSTRAINT loans_token_requires_contract CHECK (token_id IS NULL OR contract_id IS NOT NULL)
);

CREATE UNIQUE INDEX IF NOT EXISTS loans_contract_token_unique
    ON loans (contract_id, token_id)
    WHERE token_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS loans_lender_status_idx
    ON loans (lower(lender_address), status);

CREATE TABLE IF NOT EXISTS repayments (
    id bigserial PRIMARY KEY,
    loan_id bigint NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    amount_minor bigint NOT NULL,
    external_ref text,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT repayments_amount_positive CHECK (amount_minor > 0)
);

CREATE UNIQUE INDEX IF NOT EXISTS repayments_external_ref_unique
    ON repayments (loan_id, external_ref)
    WHERE external_ref IS NOT NULL;

CREATE TABLE IF NOT EXISTS chain_operations (
    id bigserial PRIMARY KEY,
    kind text NOT NULL,
    status text NOT NULL,
    contract_id bigint REFERENCES contracts(id),
    loan_id bigint REFERENCES loans(id),
    tx_hash text,
    nonce bigint,
    error text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT chain_operations_kind_valid CHECK (
        kind IN ('deploy_contract', 'originate', 'transfer', 'settle', 'mark_defaulted')
    ),
    CONSTRAINT chain_operations_status_valid CHECK (
        status IN ('created', 'submitted', 'mined', 'applied', 'retryable')
    ),
    CONSTRAINT chain_operations_nonce_non_negative CHECK (nonce IS NULL OR nonce >= 0)
);

-- Entry point for finding stuck operations (status = 'retryable') for
-- manual or automated recovery.
CREATE INDEX IF NOT EXISTS chain_operations_status_idx
    ON chain_operations (status);

ALTER TABLE loans
    ADD CONSTRAINT loans_mint_operation_fk
    FOREIGN KEY (mint_operation_id) REFERENCES chain_operations(id)
    DEFERRABLE INITIALLY DEFERRED;

CREATE TABLE IF NOT EXISTS app_locks (
    name text PRIMARY KEY,
    holder text NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT app_locks_name_not_empty CHECK (name <> ''),
    CONSTRAINT app_locks_holder_not_empty CHECK (holder <> '')
);

-- +goose Down
DROP TABLE IF EXISTS app_locks;

ALTER TABLE IF EXISTS loans
    DROP CONSTRAINT IF EXISTS loans_mint_operation_fk;

DROP TABLE IF EXISTS chain_operations;
DROP TABLE IF EXISTS repayments;
DROP TABLE IF EXISTS loans;
DROP TABLE IF EXISTS contracts;
