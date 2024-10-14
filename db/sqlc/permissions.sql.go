// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: permissions.sql

package db

import (
	"context"
)

const deleteUserPermissionById = `-- name: DeleteUserPermissionById :exec
DELETE FROM users_permissions WHERE  user_id = $1 and permission_id = $2
`

type DeleteUserPermissionByIdParams struct {
	UserID       int64 `json:"user_id"`
	PermissionID int16 `json:"permission_id"`
}

func (q *Queries) DeleteUserPermissionById(ctx context.Context, arg DeleteUserPermissionByIdParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserPermissionById, arg.UserID, arg.PermissionID)
	return err
}

const insertUserPermission = `-- name: InsertUserPermission :one
INSERT INTO users_permissions (user_id, permission_id)
VALUES ($1, $2)
RETURNING user_id, permission_id
`

type InsertUserPermissionParams struct {
	UserID       int64 `json:"user_id"`
	PermissionID int16 `json:"permission_id"`
}

func (q *Queries) InsertUserPermission(ctx context.Context, arg InsertUserPermissionParams) (UsersPermission, error) {
	row := q.db.QueryRowContext(ctx, insertUserPermission, arg.UserID, arg.PermissionID)
	var i UsersPermission
	err := row.Scan(&i.UserID, &i.PermissionID)
	return i, err
}

const selectPermissionByCode = `-- name: SelectPermissionByCode :one
select id, code from permissions where code = $1
`

func (q *Queries) SelectPermissionByCode(ctx context.Context, code string) (Permission, error) {
	row := q.db.QueryRowContext(ctx, selectPermissionByCode, code)
	var i Permission
	err := row.Scan(&i.ID, &i.Code)
	return i, err
}

const selectUserPermissions = `-- name: SelectUserPermissions :many
SELECT p.code FROM users_permissions up
JOIN permissions p on p.id = up.permission_id
WHERE user_id = $1
`

func (q *Queries) SelectUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, selectUserPermissions, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		items = append(items, code)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
