-- name: GetOrCreateIdentity :one
INSERT INTO identities (
    issuer,
    subject,
    role
) VALUES (
    $1, $2, $3
)
ON CONFLICT (issuer, subject) DO UPDATE SET role = identities.role
RETURNING *;

-- name: GetIdentityByID :one
SELECT *
FROM identities
WHERE id = $1;

-- name: GetIdentityByIssuerSubject :one
SELECT *
FROM identities
WHERE issuer = $1
  AND subject = $2;

-- name: CreateSigningKey :one
INSERT INTO signing_keys (
    identity_id,
    address,
    ciphertext,
    encryptor,
    key_version
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetSigningKeyByIdentityID :one
SELECT *
FROM signing_keys
WHERE identity_id = $1;

-- name: GetSigningKeyByAddress :one
SELECT *
FROM signing_keys
WHERE lower(address) = lower($1);
