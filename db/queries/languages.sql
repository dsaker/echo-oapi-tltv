-- name: SelectLanguagesById :one
SELECT * FROM languages WHERE id = $1;

-- name: ListLanguages :many
SELECT * FROM languages;