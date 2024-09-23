// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
)

const deleteUserById = `-- name: DeleteUserById :exec
DELETE FROM users WHERE  id = $1
`

func (q *Queries) DeleteUserById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteUserById, id)
	return err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO users (name, email, hashed_password, title_id, flipped, og_language_id, new_language_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, title_id, name, email, hashed_password, flipped, og_language_id, new_language_id, created
`

type InsertUserParams struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	HashedPassword string `json:"hashed_password"`
	TitleID        int64  `json:"title_id"`
	Flipped        bool   `json:"flipped"`
	OgLanguageID   int64  `json:"og_language_id"`
	NewLanguageID  int64  `json:"new_language_id"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, insertUser,
		arg.Name,
		arg.Email,
		arg.HashedPassword,
		arg.TitleID,
		arg.Flipped,
		arg.OgLanguageID,
		arg.NewLanguageID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TitleID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.Flipped,
		&i.OgLanguageID,
		&i.NewLanguageID,
		&i.Created,
	)
	return i, err
}

const selectUserById = `-- name: SelectUserById :one
SELECT id, title_id, name, email, hashed_password, flipped, og_language_id, new_language_id, created FROM users WHERE  id = $1
`

func (q *Queries) SelectUserById(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TitleID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.Flipped,
		&i.OgLanguageID,
		&i.NewLanguageID,
		&i.Created,
	)
	return i, err
}

const selectUserByName = `-- name: SelectUserByName :one
SELECT id, title_id, name, email, hashed_password, flipped, og_language_id, new_language_id, created FROM users WHERE  name = $1
`

func (q *Queries) SelectUserByName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, selectUserByName, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TitleID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.Flipped,
		&i.OgLanguageID,
		&i.NewLanguageID,
		&i.Created,
	)
	return i, err
}

const updateUserById = `-- name: UpdateUserById :one
UPDATE users
SET title_id = $1, email = $2, flipped = $3, og_language_id = $4, new_language_id = $5, hashed_password = $6
WHERE id = $7
RETURNING id, title_id, name, email, hashed_password, flipped, og_language_id, new_language_id, created
`

type UpdateUserByIdParams struct {
	TitleID        int64  `json:"title_id"`
	Email          string `json:"email"`
	Flipped        bool   `json:"flipped"`
	OgLanguageID   int64  `json:"og_language_id"`
	NewLanguageID  int64  `json:"new_language_id"`
	HashedPassword string `json:"hashed_password"`
	ID             int64  `json:"id"`
}

func (q *Queries) UpdateUserById(ctx context.Context, arg UpdateUserByIdParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserById,
		arg.TitleID,
		arg.Email,
		arg.Flipped,
		arg.OgLanguageID,
		arg.NewLanguageID,
		arg.HashedPassword,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.TitleID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.Flipped,
		&i.OgLanguageID,
		&i.NewLanguageID,
		&i.Created,
	)
	return i, err
}
