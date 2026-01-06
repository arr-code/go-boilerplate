-- name: AssignRoleToUser :one
INSERT INTO user_roles (user_id, role_id, assigned_by)
VALUES ($1, $2, $3)
RETURNING *;

-- name: RemoveRoleFromUser :exec
DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2;

-- name: GetUserRoles :many
SELECT r.* FROM roles r
INNER JOIN user_roles ur ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: GetUserWithRoles :one
SELECT
    u.id,
    u.email,
    u.username,
    u.is_active,
    u.email_verified,
    u.last_login,
    u.created_at,
    u.updated_at,
    COALESCE(
        json_agg(
            json_build_object(
                'id', r.id,
                'role_name', r.role_name,
                'description', r.description,
                'permissions', r.permissions
            )
        ) FILTER (WHERE r.id IS NOT NULL),
        '[]'
    ) as roles
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE u.id = $1
GROUP BY u.id;

-- name: CheckUserHasRole :one
SELECT EXISTS(
    SELECT 1 FROM user_roles ur
    INNER JOIN roles r ON ur.role_id = r.id
    WHERE ur.user_id = $1 AND r.role_name = $2
) as has_role;

-- name: CheckUserHasAnyRole :one
SELECT EXISTS(
    SELECT 1 FROM user_roles ur
    INNER JOIN roles r ON ur.role_id = r.id
    WHERE ur.user_id = $1 AND r.role_name = ANY($2::text[])
) as has_any_role;

-- name: RemoveAllUserRoles :exec
DELETE FROM user_roles WHERE user_id = $1;
