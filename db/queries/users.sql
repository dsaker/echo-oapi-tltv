-- name: InsertUser :one
INSERT INTO users (name, email, hashed_password, title_id, og_language_id, new_language_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users WHERE  id = $1;

-- name: SelectUserById :one
SELECT * FROM users WHERE  id = $1;

-- name: SelectUserByName :one
SELECT * FROM users WHERE  name = $1;

-- name: UpdateUserById :one
UPDATE users
SET title_id = $1, email = $2, og_language_id = $3, new_language_id = $4, hashed_password = $5
WHERE id = $6
RETURNING *;