-- name: CreateSession :one
INSERT INTO sessions (user_id, refresh_token, expires_at, ip_address, user_agent)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, refresh_token;

-- name: UpdateSessionByToken :one
UPDATE sessions
SET refresh_token = $2,
    updated_at = now(),
    expires_at = $3
WHERE refresh_token = $1
  AND revoked_at IS NULL
  AND expires_at > now()
RETURNING id, refresh_token;

-- name: RevokeSessionByToken :exec
UPDATE sessions
SET revoked_at = now(),
    updated_at = now()
WHERE refresh_token = $1
  AND revoked_at IS NULL;

-- name: FindValidSessionByToken :one
SELECT * FROM sessions
WHERE refresh_token = $1
  AND revoked_at IS NULL
  AND expires_at > now();