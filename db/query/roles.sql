-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = $1 LIMIT 1;

-- name: GetRoleByName :one
SELECT * FROM roles WHERE role_name = $1 LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles ORDER BY role_name;

-- name: CreateRole :one
INSERT INTO roles (role_name, description, permissions)
VALUES ($1, $2, $3)
RETURNING *;
