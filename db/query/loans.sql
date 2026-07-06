-- name: CreateLoan :one
INSERT INTO loans (
    borrower_ref,
    lender_address,
    lender_identity_id,
    principal_minor,
    apr_bps,
    term_days,
    interest_due_minor,
    total_due_minor,
    outstanding_minor,
    status,
    mint_operation_id,
    contract_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: SetLoanActive :one
UPDATE loans
SET token_id = $2,
    status = 'active',
    mint_operation_id = $3,
    contract_id = $4,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: GetLoanByID :one
SELECT *
FROM loans
WHERE id = $1;

-- name: GetLoanByIDForUpdate :one
SELECT *
FROM loans
WHERE id = $1
FOR UPDATE;

-- name: ListLoans :many
SELECT *
FROM loans
WHERE (@lender::text = '' OR lower(lender_address) = lower(@lender::text))
  AND (@status::text = '' OR status = @status::text)
ORDER BY id DESC
LIMIT @limit_count
OFFSET @offset_count;

-- name: UpdateLoanOutstandingAndStatus :one
UPDATE loans
SET outstanding_minor = $2,
    status = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateLoanLender :one
UPDATE loans
SET lender_address = $2,
    lender_identity_id = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: SetLoanStatus :one
UPDATE loans
SET status = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;
