-- name: GetUser :one
SELECT *
FROM taxeer_user
WHERE telegram_user_id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO taxeer_user (telegram_user_id, chat_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateUserChatId :exec
UPDATE taxeer_user
SET chat_id = $2
WHERE id = $1;