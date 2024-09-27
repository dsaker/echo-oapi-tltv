-- name: InsertTitle :one
INSERT INTO titles (title, num_subs, language_id, og_language_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: SelectTitleById :one
SELECT * FROM titles WHERE  id = $1;

-- name: DeleteTitleById :exec
DELETE FROM titles WHERE  id = $1;

-- name: ListTitlesByLanguage :many
SELECT title, similarity(title, $1) AS similarity, num_subs, language_id, og_language_id
FROM titles
WHERE language_id = $2
ORDER BY similarity
LIMIT $3;

-- name: ListTitles :many
SELECT id, title, similarity(title, $1) AS similarity, num_subs, language_id, og_language_id
FROM titles
ORDER BY similarity
LIMIT $2;
