-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = ? LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM transactions
ORDER BY date;

-- name: CreateTransaction :exec
INSERT INTO transactions (
  date, code, description, amount
) VALUES (
  ?, ?, ?, ?
);

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = ?;
