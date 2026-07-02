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
