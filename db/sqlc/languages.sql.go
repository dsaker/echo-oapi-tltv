// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: languages.sql

package db

import (
	"context"
)

const listLanguages = `-- name: ListLanguages :many
SELECT id, language, tag FROM languages
`

func (q *Queries) ListLanguages(ctx context.Context) ([]Language, error) {
	rows, err := q.db.QueryContext(ctx, listLanguages)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Language{}
	for rows.Next() {
		var i Language
		if err := rows.Scan(&i.ID, &i.Language, &i.Tag); err != nil {
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

const listLanguagesSimilar = `-- name: ListLanguagesSimilar :many
SELECT id, language, tag, similarity(language, $1) AS similarity
FROM languages
ORDER BY similarity desc
`

type ListLanguagesSimilarRow struct {
	ID         int16   `json:"id"`
	Language   string  `json:"language"`
	Tag        string  `json:"tag"`
	Similarity float32 `json:"similarity"`
}

func (q *Queries) ListLanguagesSimilar(ctx context.Context, similarity string) ([]ListLanguagesSimilarRow, error) {
	rows, err := q.db.QueryContext(ctx, listLanguagesSimilar, similarity)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListLanguagesSimilarRow{}
	for rows.Next() {
		var i ListLanguagesSimilarRow
		if err := rows.Scan(
			&i.ID,
			&i.Language,
			&i.Tag,
			&i.Similarity,
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

const selectLanguagesById = `-- name: SelectLanguagesById :one
SELECT id, language, tag FROM languages WHERE id = $1
`

func (q *Queries) SelectLanguagesById(ctx context.Context, id int16) (Language, error) {
	row := q.db.QueryRowContext(ctx, selectLanguagesById, id)
	var i Language
	err := row.Scan(&i.ID, &i.Language, &i.Tag)
	return i, err
}
