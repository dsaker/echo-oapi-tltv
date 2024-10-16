// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: phrases.sql

package db

import (
	"context"
)

const insertPhrases = `-- name: InsertPhrases :one
INSERT INTO phrases (title_id)
VALUES ($1)
RETURNING id, title_id
`

func (q *Queries) InsertPhrases(ctx context.Context, titleID int64) (Phrase, error) {
	row := q.db.QueryRowContext(ctx, insertPhrases, titleID)
	var i Phrase
	err := row.Scan(&i.ID, &i.TitleID)
	return i, err
}

const selectPhraseIdsByTitleId = `-- name: SelectPhraseIdsByTitleId :many
SELECT id FROM phrases
WHERE title_id = $1
ORDER BY id
`

func (q *Queries) SelectPhraseIdsByTitleId(ctx context.Context, titleID int64) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, selectPhraseIdsByTitleId, titleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectPhrasesFromTranslates = `-- name: SelectPhrasesFromTranslates :many
SELECT og.phrase_id, og.phrase, og.phrase_hint, new.phrase, new.phrase_hint
FROM
    (
        SELECT phrase_id, phrase, phrase_hint
        FROM translates t
        where t.language_id = $1
    ) og
    JOIN (
        SELECT phrase, phrase_hint, phrase_id
        FROM translates t
        WHERE t.language_id = $2
    ) new
ON og.phrase_id = new.phrase_id
WHERE og.phrase_id
IN (
    SELECT phrase_id from users_phrases
    WHERE user_id = $3 AND title_id = $4 AND language_id = $2
    ORDER BY  phrase_correct, phrase_id
    LIMIT $5
    )
`

type SelectPhrasesFromTranslatesParams struct {
	LanguageID   int16 `json:"language_id"`
	LanguageID_2 int16 `json:"language_id_2"`
	UserID       int64 `json:"user_id"`
	TitleID      int64 `json:"title_id"`
	Limit        int32 `json:"limit"`
}

type SelectPhrasesFromTranslatesRow struct {
	PhraseID     int64  `json:"phrase_id"`
	Phrase       string `json:"phrase"`
	PhraseHint   string `json:"phrase_hint"`
	Phrase_2     string `json:"phrase_2"`
	PhraseHint_2 string `json:"phrase_hint_2"`
}

func (q *Queries) SelectPhrasesFromTranslates(ctx context.Context, arg SelectPhrasesFromTranslatesParams) ([]SelectPhrasesFromTranslatesRow, error) {
	rows, err := q.db.QueryContext(ctx, selectPhrasesFromTranslates,
		arg.LanguageID,
		arg.LanguageID_2,
		arg.UserID,
		arg.TitleID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SelectPhrasesFromTranslatesRow{}
	for rows.Next() {
		var i SelectPhrasesFromTranslatesRow
		if err := rows.Scan(
			&i.PhraseID,
			&i.Phrase,
			&i.PhraseHint,
			&i.Phrase_2,
			&i.PhraseHint_2,
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

const selectPhrasesFromTranslatesWithCorrect = `-- name: SelectPhrasesFromTranslatesWithCorrect :many
SELECT og.phrase_id, og.phrase, og.phrase_hint, new.phrase, new.phrase_hint, up.phrase_correct
FROM
    (
        SELECT phrase_id, phrase, phrase_hint
        FROM translates t
        where t.language_id = $1
    ) og
    JOIN (
        SELECT phrase, phrase_hint, phrase_id
        FROM translates t
        WHERE t.language_id = $2
    ) new
ON og.phrase_id = new.phrase_id
    JOIN (
        SELECT phrase_id, phrase_correct
        FROM users_phrases u
        WHERE u.user_id = $3 AND u.title_id = $4 AND u.language_id = $2
    ) up
ON new.phrase_id = up.phrase_id
WHERE og.phrase_id
IN
    (
      SELECT phrase_id from users_phrases
      WHERE user_id = $3 AND title_id = $4 AND language_id = $2
      ORDER BY  phrase_correct, phrase_id
      LIMIT $5
    )
`

type SelectPhrasesFromTranslatesWithCorrectParams struct {
	LanguageID   int16 `json:"language_id"`
	LanguageID_2 int16 `json:"language_id_2"`
	UserID       int64 `json:"user_id"`
	TitleID      int64 `json:"title_id"`
	Limit        int32 `json:"limit"`
}

type SelectPhrasesFromTranslatesWithCorrectRow struct {
	PhraseID      int64  `json:"phrase_id"`
	Phrase        string `json:"phrase"`
	PhraseHint    string `json:"phrase_hint"`
	Phrase_2      string `json:"phrase_2"`
	PhraseHint_2  string `json:"phrase_hint_2"`
	PhraseCorrect int16  `json:"phrase_correct"`
}

func (q *Queries) SelectPhrasesFromTranslatesWithCorrect(ctx context.Context, arg SelectPhrasesFromTranslatesWithCorrectParams) ([]SelectPhrasesFromTranslatesWithCorrectRow, error) {
	rows, err := q.db.QueryContext(ctx, selectPhrasesFromTranslatesWithCorrect,
		arg.LanguageID,
		arg.LanguageID_2,
		arg.UserID,
		arg.TitleID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []SelectPhrasesFromTranslatesWithCorrectRow{}
	for rows.Next() {
		var i SelectPhrasesFromTranslatesWithCorrectRow
		if err := rows.Scan(
			&i.PhraseID,
			&i.Phrase,
			&i.PhraseHint,
			&i.Phrase_2,
			&i.PhraseHint_2,
			&i.PhraseCorrect,
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

const selectUsersPhrasesByCorrect = `-- name: SelectUsersPhrasesByCorrect :many
SELECT phrase_id from users_phrases
WHERE user_id = $1 and title_id = $2 and language_id = $3
ORDER BY  phrase_correct, phrase_id
limit $4
`

type SelectUsersPhrasesByCorrectParams struct {
	UserID     int64 `json:"user_id"`
	TitleID    int64 `json:"title_id"`
	LanguageID int16 `json:"language_id"`
	Limit      int32 `json:"limit"`
}

func (q *Queries) SelectUsersPhrasesByCorrect(ctx context.Context, arg SelectUsersPhrasesByCorrectParams) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, selectUsersPhrasesByCorrect,
		arg.UserID,
		arg.TitleID,
		arg.LanguageID,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var phrase_id int64
		if err := rows.Scan(&phrase_id); err != nil {
			return nil, err
		}
		items = append(items, phrase_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectUsersPhrasesByIds = `-- name: SelectUsersPhrasesByIds :one
SELECT user_id, title_id, phrase_id, language_id, phrase_correct from users_phrases
WHERE user_id = $1 and language_id = $2 and phrase_id = $3
`

type SelectUsersPhrasesByIdsParams struct {
	UserID     int64 `json:"user_id"`
	LanguageID int16 `json:"language_id"`
	PhraseID   int64 `json:"phrase_id"`
}

func (q *Queries) SelectUsersPhrasesByIds(ctx context.Context, arg SelectUsersPhrasesByIdsParams) (UsersPhrase, error) {
	row := q.db.QueryRowContext(ctx, selectUsersPhrasesByIds, arg.UserID, arg.LanguageID, arg.PhraseID)
	var i UsersPhrase
	err := row.Scan(
		&i.UserID,
		&i.TitleID,
		&i.PhraseID,
		&i.LanguageID,
		&i.PhraseCorrect,
	)
	return i, err
}

const updateUsersPhrasesByThreeIds = `-- name: UpdateUsersPhrasesByThreeIds :one
UPDATE users_phrases
SET user_id = $1, title_id = $2, phrase_id = $3, language_id = $4, phrase_correct = $5
WHERE user_id = $1 AND phrase_id = $3 AND language_id = $4
RETURNING user_id, title_id, phrase_id, language_id, phrase_correct
`

type UpdateUsersPhrasesByThreeIdsParams struct {
	UserID        int64 `json:"user_id"`
	TitleID       int64 `json:"title_id"`
	PhraseID      int64 `json:"phrase_id"`
	LanguageID    int16 `json:"language_id"`
	PhraseCorrect int16 `json:"phrase_correct"`
}

func (q *Queries) UpdateUsersPhrasesByThreeIds(ctx context.Context, arg UpdateUsersPhrasesByThreeIdsParams) (UsersPhrase, error) {
	row := q.db.QueryRowContext(ctx, updateUsersPhrasesByThreeIds,
		arg.UserID,
		arg.TitleID,
		arg.PhraseID,
		arg.LanguageID,
		arg.PhraseCorrect,
	)
	var i UsersPhrase
	err := row.Scan(
		&i.UserID,
		&i.TitleID,
		&i.PhraseID,
		&i.LanguageID,
		&i.PhraseCorrect,
	)
	return i, err
}
