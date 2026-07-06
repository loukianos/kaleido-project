-- name: CreateContract :one
INSERT INTO contracts (
    chain_id,
    address,
    deploy_tx_hash,
    base_uri,
    active
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetActiveContractByChainID :one
SELECT *
FROM contracts
WHERE chain_id = $1
  AND active = true;

-- name: GetContractByID :one
SELECT *
FROM contracts
WHERE id = $1;

-- name: ListContractsByChainID :many
SELECT *
FROM contracts
WHERE chain_id = $1
ORDER BY id;

-- name: DeactivateActiveContract :exec
UPDATE contracts
SET active = false,
    updated_at = now()
WHERE chain_id = $1
  AND active = true;

-- name: ActivateContract :one
UPDATE contracts
SET active = true,
    updated_at = now()
WHERE id = $1
RETURNING *;
