-- name: SelectLanguagesById :one
select * from languages where id = $1;