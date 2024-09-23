-- name: InsertUserPermission :one
INSERT INTO users_permissions (user_id, permission_id)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteUserPermissionById :exec
DELETE FROM users_permissions WHERE  user_id = $1 and permission_id = $2;

-- name: SelectUserPermissions :many
SELECT p.code FROM users_permissions up
JOIN permissions p on p.id = up.permission_id
WHERE user_id = $1;

-- name: SelectPermissionByCode :one
select * from permissions where code = $1;