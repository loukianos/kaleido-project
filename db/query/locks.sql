-- name: AcquireAppLock :one
INSERT INTO app_locks (
    name,
    holder,
    expires_at
) VALUES (
    $1, $2, $3
)
ON CONFLICT (name) DO UPDATE
SET holder = EXCLUDED.holder,
    expires_at = EXCLUDED.expires_at,
    updated_at = now()
WHERE app_locks.expires_at < now()
   OR app_locks.holder = EXCLUDED.holder
RETURNING *;

-- name: ReleaseAppLock :exec
DELETE FROM app_locks
WHERE name = $1
  AND holder = $2;
