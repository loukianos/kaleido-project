-- +goose Up
ALTER TABLE chain_operations
    ADD COLUMN IF NOT EXISTS attempts integer NOT NULL DEFAULT 0;

ALTER TABLE chain_operations
    DROP CONSTRAINT chain_operations_status_valid;
ALTER TABLE chain_operations
    ADD CONSTRAINT chain_operations_status_valid CHECK (
        status IN ('created', 'submitted', 'mined', 'applied', 'retryable', 'failed')
    );

ALTER TABLE loans
    DROP CONSTRAINT loans_status_valid;
ALTER TABLE loans
    ADD CONSTRAINT loans_status_valid CHECK (
        status IN ('originating', 'active', 'settling', 'repaid', 'defaulted', 'failed')
    );

ALTER TABLE loans
    ADD COLUMN IF NOT EXISTS external_ref text;

CREATE UNIQUE INDEX IF NOT EXISTS loans_external_ref_unique
    ON loans (external_ref)
    WHERE external_ref IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS loans_external_ref_unique;
ALTER TABLE loans DROP COLUMN IF EXISTS external_ref;
ALTER TABLE loans
    DROP CONSTRAINT loans_status_valid;
ALTER TABLE loans
    ADD CONSTRAINT loans_status_valid CHECK (
        status IN ('originating', 'active', 'settling', 'repaid', 'defaulted')
    );
ALTER TABLE chain_operations
    DROP CONSTRAINT chain_operations_status_valid;
ALTER TABLE chain_operations
    ADD CONSTRAINT chain_operations_status_valid CHECK (
        status IN ('created', 'submitted', 'mined', 'applied', 'retryable')
    );
ALTER TABLE chain_operations DROP COLUMN IF EXISTS attempts;
