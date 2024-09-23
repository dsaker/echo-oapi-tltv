-- name: InsertUser :one
INSERT INTO users (name, email, hashed_password, title_id, flipped, og_language_id, new_language_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteUserById :exec
DELETE FROM users WHERE  id = $1;

-- name: SelectUserById :one
SELECT * FROM users WHERE  id = $1;

-- name: SelectUserByName :one
SELECT * FROM users WHERE  name = $1;

-- name: UpdateUserById :one
UPDATE users
SET title_id = $1, email = $2, flipped = $3, og_language_id = $4, new_language_id = $5, hashed_password = $6
WHERE id = $7
RETURNING *;