-- name: InsertTitle :one
INSERT INTO titles (title, num_subs, language_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: SelectTitleById :one
SELECT * FROM titles WHERE  id = $1;

-- name: DeleteTitleById :exec
DELETE FROM titles WHERE  id = $1;

-- name: ListTitlesByLanguage :many
SELECT count(*) OVER(), id, title, similarity(title, $1) AS similarity, num_subs
FROM titles
WHERE language_id = $2
LIMIT $3;

-- name: ListTitles :many
SELECT count(*) OVER(), id, title, similarity(title, $1) AS similarity, num_subs
FROM titles
LIMIT $2;
