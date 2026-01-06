-- name: CreateUserDetails :one
INSERT INTO user_details (user_id, full_name, phone, address, date_of_birth, avatar_url, bio)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserDetailsByUserID :one
SELECT * FROM user_details WHERE user_id = $1 LIMIT 1;

-- name: UpdateUserDetails :exec
UPDATE user_details
SET full_name = COALESCE(NULLIF($2, ''), full_name),
    phone = COALESCE(NULLIF($3, ''), phone),
    address = COALESCE(NULLIF($4, ''), address),
    date_of_birth = COALESCE($5, date_of_birth),
    bio = COALESCE(NULLIF($6, ''), bio),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;

-- name: UpdateUserAvatar :exec
UPDATE user_details
SET avatar_url = $2, updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1;
