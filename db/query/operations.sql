-- name: CreateChainOperation :one
INSERT INTO chain_operations (
    kind,
    status,
    contract_id,
    loan_id
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetChainOperationByID :one
SELECT *
FROM chain_operations
WHERE id = $1;

-- name: SetOperationSubmitted :one
UPDATE chain_operations
SET status = 'submitted',
    tx_hash = $2,
    nonce = $3,
    error = NULL,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetOperationMined :one
UPDATE chain_operations
SET status = 'mined',
    error = NULL,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetOperationApplied :one
UPDATE chain_operations
SET status = 'applied',
    error = NULL,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetOperationRetryable :one
UPDATE chain_operations
SET status = 'retryable',
    error = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: AttachOperationContract :one
UPDATE chain_operations
SET contract_id = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;
