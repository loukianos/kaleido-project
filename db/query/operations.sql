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
    signer_address = $4,
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
    attempts = attempts + 1,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetOperationFailed :one
UPDATE chain_operations
SET status = 'failed',
    error = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ListRetryableOperations :many
SELECT *
FROM chain_operations
WHERE status = 'retryable'
  AND attempts < @max_attempts::integer
ORDER BY updated_at
LIMIT @limit_count;

-- name: ListExhaustedOperations :many
SELECT *
FROM chain_operations
WHERE status = 'retryable'
  AND attempts >= @max_attempts::integer
ORDER BY updated_at
LIMIT @limit_count;

-- name: ListStaleSubmittedOperations :many
SELECT *
FROM chain_operations
WHERE status = 'submitted'
  AND updated_at < @stale_before::timestamptz
ORDER BY updated_at
LIMIT @limit_count;

-- name: AttachOperationContract :one
UPDATE chain_operations
SET contract_id = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;
