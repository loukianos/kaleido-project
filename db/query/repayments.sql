-- name: CreateRepayment :one
INSERT INTO repayments (
    loan_id,
    amount_minor,
    external_ref
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: ListRepaymentsByLoan :many
SELECT *
FROM repayments
WHERE loan_id = $1
ORDER BY created_at ASC, id ASC;
