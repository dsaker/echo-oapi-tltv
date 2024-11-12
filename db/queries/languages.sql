-- name: SelectLanguagesById :one
SELECT * FROM languages WHERE id = $1;

-- name: ListLanguagesSimilar :many
SELECT id, language, tag, similarity(language, $1) AS similarity
FROM languages
ORDER BY similarity desc;

-- name: ListLanguages :many
SELECT * FROM languages;