-- +goose Up
ALTER TABLE chain_operations
    ADD COLUMN IF NOT EXISTS signer_address text;

ALTER TABLE chain_operations
    DROP CONSTRAINT chain_operations_kind_valid;
ALTER TABLE chain_operations
    ADD CONSTRAINT chain_operations_kind_valid CHECK (
        kind IN ('deploy_contract', 'originate', 'transfer', 'settle', 'mark_defaulted', 'grant_role')
    );

-- +goose Down
ALTER TABLE chain_operations
    DROP CONSTRAINT chain_operations_kind_valid;
ALTER TABLE chain_operations
    ADD CONSTRAINT chain_operations_kind_valid CHECK (
        kind IN ('deploy_contract', 'originate', 'transfer', 'settle', 'mark_defaulted')
    );
ALTER TABLE chain_operations DROP COLUMN IF EXISTS signer_address;
