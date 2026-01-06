-- name: CreateUser :one
INSERT INTO users (email, username, password_hash, is_active, email_verified)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: GetUserByEmailOrUsername :one
SELECT * FROM users WHERE email = $1 OR username = $2 LIMIT 1;

-- name: UpdateUserLastLogin :exec
UPDATE users SET last_login = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeactivateUser :exec
UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ActivateUser :exec
UPDATE users SET is_active = true, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: VerifyUserEmail :exec
UPDATE users SET email_verified = true, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;
