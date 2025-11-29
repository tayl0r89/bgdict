-- name: GetWord :one
SELECT sqlc.embed(word), sqlc.embed(word_type)
FROM word
JOIN word_type on word.type_id = word_type.id
WHERE word.id = ? LIMIT 1;

-- name: FindWords :many
SELECT sqlc.embed(derivative_form), sqlc.embed(word), sqlc.embed(word_type) 
FROM derivative_form
JOIN word on derivative_form.base_word_id = word.id
JOIN word_type on word.type_id = word_type.id
WHERE derivative_form.name = ?;