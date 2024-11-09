-- name: SelectVoicesByLanguageId :many
SELECT * FROM voices WHERE language_id = $1;

-- name: ListVoices :many
SELECT * FROM voices;

-- name: SelectVoiceById :one
SELECT * FROM voices WHERE id = $1;