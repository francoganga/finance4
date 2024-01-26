// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package models

import (
	"context"
	"database/sql"

	types "finance/internal/types"
)

const createTransaction = `-- name: CreateTransaction :exec
INSERT INTO transactions (
  date, code, description, amount, balance
) VALUES (
  ?, ?, ?, ?, ?
)
`

type CreateTransactionParams struct {
	Date        types.Date
	Code        sql.NullString
	Description string
	Amount      int64
	Balance     int64
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) error {
	_, err := q.db.ExecContext(ctx, createTransaction,
		arg.Date,
		arg.Code,
		arg.Description,
		arg.Amount,
		arg.Balance,
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
SELECT id, date, code, description, amount, balance FROM transactions
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
		&i.Balance,
	)
	return i, err
}

const lastMonthTransactions = `-- name: LastMonthTransactions :many
SELECT id, date, code, description, amount, balance FROM transactions WHERE strftime('%Y-%m', date) = (SELECT strftime('%Y-%m', date) FROM transactions order by date desc limit 1) ORDER BY date DESC
`

func (q *Queries) LastMonthTransactions(ctx context.Context) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, lastMonthTransactions)
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
			&i.Balance,
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

const listTransactions = `-- name: ListTransactions :many
SELECT id, date, code, description, amount, balance FROM transactions
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
			&i.Balance,
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
