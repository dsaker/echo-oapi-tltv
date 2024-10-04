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
INSERT INTO titles (title, num_subs, og_language_id)
VALUES ($1, $2, $3)
RETURNING id, title, num_subs, og_language_id
`

type InsertTitleParams struct {
	Title        string `json:"title"`
	NumSubs      int16  `json:"num_subs"`
	OgLanguageID int16  `json:"og_language_id"`
}

func (q *Queries) InsertTitle(ctx context.Context, arg InsertTitleParams) (Title, error) {
	row := q.db.QueryRowContext(ctx, insertTitle, arg.Title, arg.NumSubs, arg.OgLanguageID)
	var i Title
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.NumSubs,
		&i.OgLanguageID,
	)
	return i, err
}

const listTitles = `-- name: ListTitles :many
SELECT id, title, similarity(title, $1) AS similarity, num_subs, og_language_id
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
	NumSubs      int16   `json:"num_subs"`
	OgLanguageID int16   `json:"og_language_id"`
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

const listTitlesByOgLanguage = `-- name: ListTitlesByOgLanguage :many
SELECT title, similarity(title, $1) AS similarity, num_subs, og_language_id
FROM titles
WHERE og_language_id = $2
ORDER BY similarity
LIMIT $3
`

type ListTitlesByOgLanguageParams struct {
	Similarity   string `json:"similarity"`
	OgLanguageID int16  `json:"og_language_id"`
	Limit        int32  `json:"limit"`
}

type ListTitlesByOgLanguageRow struct {
	Title        string  `json:"title"`
	Similarity   float32 `json:"similarity"`
	NumSubs      int16   `json:"num_subs"`
	OgLanguageID int16   `json:"og_language_id"`
}

func (q *Queries) ListTitlesByOgLanguage(ctx context.Context, arg ListTitlesByOgLanguageParams) ([]ListTitlesByOgLanguageRow, error) {
	rows, err := q.db.QueryContext(ctx, listTitlesByOgLanguage, arg.Similarity, arg.OgLanguageID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListTitlesByOgLanguageRow{}
	for rows.Next() {
		var i ListTitlesByOgLanguageRow
		if err := rows.Scan(
			&i.Title,
			&i.Similarity,
			&i.NumSubs,
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
SELECT id, title, num_subs, og_language_id FROM titles WHERE  id = $1
`

func (q *Queries) SelectTitleById(ctx context.Context, id int64) (Title, error) {
	row := q.db.QueryRowContext(ctx, selectTitleById, id)
	var i Title
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.NumSubs,
		&i.OgLanguageID,
	)
	return i, err
}
