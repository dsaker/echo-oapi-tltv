-- name: InsertTranslates :one
INSERT INTO translates (phrase_id, language_id, phrase, phrase_hint)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SelectTranslatesByTitleIdLangId :many
SELECT tr.* FROM titles t
                                        JOIN phrases p ON t.id = p.title_id
                                        JOIN translates tr ON p.id = tr.phrase_id AND tr.language_id = $1
WHERE tr.language_id = $1 and t.id = $2;

-- name: SelectExistsTranslates :one
SELECT EXISTS(
    SELECT 1 FROM titles t
                      JOIN phrases p ON t.id = p.title_id
                      JOIN translates tr ON p.id = tr.phrase_id AND tr.language_id = $1
    WHERE tr.language_id = $1 and t.id = $2 ) AS "exists";