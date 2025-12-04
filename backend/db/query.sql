-- name: GetWord :one
SELECT sqlc.embed(word), sqlc.embed(word_type), sqlc.embed(word_translation)
FROM word
JOIN word_type on word.type_id = word_type.id
JOIN word_translation on word_translation.word_id = word.id
WHERE word.id = ? LIMIT 1;

-- name: GetWordByName :many
SELECT sqlc.embed(word), sqlc.embed(word_type), sqlc.embed(word_translation)
FROM word
JOIN word_type on word.type_id = word_type.id
JOIN word_translation on word_translation.word_id = word.id
WHERE word.name = ?;

-- name: FindWords :many
SELECT sqlc.embed(derivative_form), sqlc.embed(word), sqlc.embed(word_type), sqlc.embed(word_translation) 
FROM derivative_form
JOIN word on derivative_form.base_word_id = word.id
JOIN word_type on word.type_id = word_type.id
JOIN word_translation on word_translation.word_id = word.id
WHERE derivative_form.name = ?;