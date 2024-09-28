// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: titles.sql

package db

import (
	"context"
)

const deleteTitleById = `-- name: DeleteTitleById :exec
DELETE FROM titles WHERE  id = $1
`

func (q *Queries) DeleteTitleById(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTitleById, id)
	return err
}

const insertTitle = `-- name: InsertTitle :one
INSERT INTO titles (title, num_subs, language_id, og_language_id)
VALUES ($1, $2, $3, $4)
RETURNING id, title, num_subs, language_id, og_language_id
`

type InsertTitleParams struct {
	Title        string `json:"title"`
	NumSubs      int32  `json:"num_subs"`
	LanguageID   int64  `json:"language_id"`
	OgLanguageID int64  `json:"og_language_id"`
}

func (q *Queries) InsertTitle(ctx context.Context, arg InsertTitleParams) (Title, error) {
	row := q.db.QueryRowContext(ctx, insertTitle,
		arg.Title,
		arg.NumSubs,
		arg.LanguageID,
		arg.OgLanguageID,
	)
	var i Title
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.NumSubs,
		&i.LanguageID,
		&i.OgLanguageID,
	)
	return i, err
}

const listTitles = `-- name: ListTitles :many
SELECT id, title, similarity(title, $1) AS similarity, num_subs, language_id, og_language_id
FROM titles
ORDER BY similarity
LIMIT $2
`

type ListTitlesParams struct {
	Similarity string `json:"similarity"`
	Limit      int32  `json:"limit"`
}

type ListTitlesRow struct {
	ID           int64   `json:"id"`
	Title        string  `json:"title"`
	Similarity   float32 `json:"similarity"`
	NumSubs      int32   `json:"num_subs"`
	LanguageID   int64   `json:"language_id"`
	OgLanguageID int64   `json:"og_language_id"`
}

func (q *Queries) ListTitles(ctx context.Context, arg ListTitlesParams) ([]ListTitlesRow, error) {
	rows, err := q.db.QueryContext(ctx, listTitles, arg.Similarity, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListTitlesRow{}
	for rows.Next() {
		var i ListTitlesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Similarity,
			&i.NumSubs,
			&i.LanguageID,
			&i.OgLanguageID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTitlesByLanguage = `-- name: ListTitlesByLanguage :many
SELECT title, similarity(title, $1) AS similarity, num_subs, language_id, og_language_id
FROM titles
WHERE language_id = $2
ORDER BY similarity
LIMIT $3
`

type ListTitlesByLanguageParams struct {
	Similarity string `json:"similarity"`
	LanguageID int64  `json:"language_id"`
	Limit      int32  `json:"limit"`
}

type ListTitlesByLanguageRow struct {
	Title        string  `json:"title"`
	Similarity   float32 `json:"similarity"`
	NumSubs      int32   `json:"num_subs"`
	LanguageID   int64   `json:"language_id"`
	OgLanguageID int64   `json:"og_language_id"`
}

func (q *Queries) ListTitlesByLanguage(ctx context.Context, arg ListTitlesByLanguageParams) ([]ListTitlesByLanguageRow, error) {
	rows, err := q.db.QueryContext(ctx, listTitlesByLanguage, arg.Similarity, arg.LanguageID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListTitlesByLanguageRow{}
	for rows.Next() {
		var i ListTitlesByLanguageRow
		if err := rows.Scan(
			&i.Title,
			&i.Similarity,
			&i.NumSubs,
			&i.LanguageID,
			&i.OgLanguageID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectTitleById = `-- name: SelectTitleById :one
SELECT id, title, num_subs, language_id, og_language_id FROM titles WHERE  id = $1
`

func (q *Queries) SelectTitleById(ctx context.Context, id int64) (Title, error) {
	row := q.db.QueryRowContext(ctx, selectTitleById, id)
	var i Title
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.NumSubs,
		&i.LanguageID,
		&i.OgLanguageID,
	)
	return i, err
}
