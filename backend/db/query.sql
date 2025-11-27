-- name: GetWord :one
SELECT * FROM word
WHERE id = ? LIMIT 1;

-- name: FindWords :many
SELECT sqlc.embed(derivative_form), sqlc.embed(word) 
FROM derivative_form
JOIN word on derivative_form.base_word_id = word.id
WHERE derivative_form.name = ?;