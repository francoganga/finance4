-- name: GetTransaction :one
SELECT * FROM transactions
WHERE id = ? LIMIT 1;

-- name: GetTransaction2 :one
SELECT t.*, (SELECT name from label where id = t.label_id) as label FROM transactions t
WHERE t.id = ?;

-- name: ListTransactions :many
SELECT * FROM transactions
ORDER BY date;


-- name: LastMonthTransactions :many
SELECT * FROM transactions WHERE strftime('%Y-%m', date) = (SELECT strftime('%Y-%m', date) FROM transactions order by date desc limit 1) ORDER BY id;

-- name: CreateTransaction :exec
INSERT INTO transactions (
  date, code, description, amount, balance
) VALUES (
  ?, ?, ?, ?, ?
);

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = ?;

-- name: ListLabels :many
SELECT * FROM label;

-- name: FindLabel :one
SELECT * FROM label WHERE id = ?;
