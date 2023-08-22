// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package models

import (
	"context"
	"database/sql"
)

const createTransaction = `-- name: CreateTransaction :exec
INSERT INTO transactions (
  date, code, description, amount
) VALUES (
  ?, ?, ?, ?
)
`

type CreateTransactionParams struct {
	Date        string
	Code        sql.NullString
	Description string
	Amount      int64
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) error {
	_, err := q.db.ExecContext(ctx, createTransaction,
		arg.Date,
		arg.Code,
		arg.Description,
		arg.Amount,
	)
	return err
}

const deleteTransaction = `-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = ?
`

func (q *Queries) DeleteTransaction(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteTransaction, id)
	return err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, date, code, description, amount FROM transactions
WHERE id = ? LIMIT 1
`

func (q *Queries) GetTransaction(ctx context.Context, id int64) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, id)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.Date,
		&i.Code,
		&i.Description,
		&i.Amount,
	)
	return i, err
}

const listTransactions = `-- name: ListTransactions :many
SELECT id, date, code, description, amount FROM transactions
ORDER BY date
`

func (q *Queries) ListTransactions(ctx context.Context) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, listTransactions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Transaction
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.Date,
			&i.Code,
			&i.Description,
			&i.Amount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}