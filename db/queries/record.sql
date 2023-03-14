-- name: GetAllRecordByUserId :many
SELECT *
FROM taxeer_record
WHERE taxeer_user_id = $1;

-- name: GetLastNRecordByUserId :many
SELECT *
FROM taxeer_record
WHERE taxeer_user_id = $1
ORDER BY date desc
LIMIT $2;

-- name: GetRecordByUserIdAndDateBetweenOrderedByDateDesc :many
SELECT *
FROM taxeer_record
WHERE taxeer_user_id = $1
AND date between $2 and $3
ORDER BY date desc;

-- name: CreateRecord :one
INSERT INTO taxeer_record (taxeer_user_id, "date", income_value, income_currency, rate)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateRecordIncomeValue :exec
UPDATE taxeer_record
SET income_value = $2
WHERE id = $1;

-- name: UpdateRecordIncomeCurrency :exec
UPDATE taxeer_record
SET income_currency = $2
WHERE id = $1;

-- name: UpdateRecordRate :exec
UPDATE taxeer_record
SET rate = $2
WHERE id = $1;