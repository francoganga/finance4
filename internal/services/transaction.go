package services

import (
	"database/sql"
	"finance/internal/models"
	"math"
	"strings"

	"github.com/samber/lo"
)

type Overview struct {
	Salary          int64
	Expenses        int64
	RemainingAmount int64
}

type Transaction struct {
	models.Transaction
	Label *string
}

type Metadata struct {
	CurrentPage  int
	PageSize     int
	FirstPage    int
	LastPage     int
	TotalRecords int
}

func (m Metadata) Pages() []int {
	pages := make([]int, m.LastPage)

	for i := 0; i < m.LastPage; i++ {
		pages[i] = i + 1
	}

	return pages
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func GenerateOverview(transactions []Transaction) Overview {

	salaryTransactions := lo.Filter(transactions, func(i Transaction, _ int) bool {
		return strings.Contains(i.Description, "haberes")
	})

	salary := lo.SumBy(salaryTransactions, func(i Transaction) int64 {
		return i.Amount
	})

	expenses := lo.Filter(transactions, func(i Transaction, _ int) bool {
		return i.Amount < 0
	})

	totalExpenses := lo.SumBy(expenses, func(i Transaction) int64 {
		return i.Amount
	})

	return Overview{
		Salary:          salary / 100,
		Expenses:        totalExpenses / 100,
		RemainingAmount: (salary + totalExpenses) / 100,
	}
}

type Filters struct {
	Page     int
	PageSize int
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f Filters) Limit() int {
	return f.PageSize
}

type SearchOpts struct {
	Search string
	Period string
	Filters
}

func SearchTransactions2(opts SearchOpts, db *sql.DB) ([]Transaction, Overview, Metadata, error) {

	rows, err := db.Query(`WITH lmt as (SELECT id, date, code, description, amount, balance, label_id FROM transactions WHERE (strftime('%Y-%m', date) = $3) OR ($3 = '') ORDER BY id)
				SELECT  t.*,
					l.name as label,
					SUM(case when amount > 0 then amount else 0 end) OVER(PARTITION BY strftime('%Y-%m', date)) / 100 as salary,
					SUM(case when amount < 0 then amount else 0 end) OVER(PARTITION BY strftime('%Y-%m', date)) / 100 as spends,
					SUM(amount) OVER(PARTITION BY strftime('%Y-%m', date)) / 100 as remainingMoney,
					-- metadata
					COUNT(*) OVER() as totalRows
				

				FROM lmt t
				LEFT JOIN label l on l.id = t.label_id
					WHERE description like $1 OR amount = $2 OR code = $2 OR balance = $2 OR $2 = ''
					LIMIT $4 OFFSET $5;`, "%"+opts.Search+"%", opts.Search, opts.Period, opts.Limit(), opts.Offset())

	if err != nil {
		return nil, Overview{}, Metadata{}, err
	}

	defer rows.Close()

	var transactions []Transaction
	var overview Overview

	totalRecords := 0

	for rows.Next() {

		var t Transaction

		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Code,
			&t.Description,
			&t.Amount,
			&t.Balance,
			&t.LabelID,
			&t.Label,
			&overview.Salary,
			&overview.Expenses,
			&overview.RemainingAmount,
			&totalRecords,
		)

		if err != nil {
			return nil, Overview{}, Metadata{}, err
		}

		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, Overview{}, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, opts.Page, opts.PageSize)

	return transactions, overview, metadata, nil
}

func SearchTransactions(opts SearchOpts, db *sql.DB) ([]Transaction, error) {

	rows, err := db.Query(`WITH lmt as (SELECT id, date, code, description, amount, balance, label_id FROM transactions WHERE (strftime('%Y-%m', date) = $3) OR ($3 = '')  ORDER BY id)

		SELECT t.*, l.name as label from lmt t LEFT JOIN label l on l.id = t.label_id WHERE description like $1 OR amount = $2 OR code = $2 OR balance = $2 OR $2 = '';`, "%"+opts.Search+"%", opts.Search, opts.Period)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {

		var t Transaction

		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Code,
			&t.Description,
			&t.Amount,
			&t.Balance,
			&t.LabelID,
			&t.Label,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil

}

func SearchLastMonthTransactions(search string, db *sql.DB) ([]Transaction, error) {

	rows, err := db.Query(`WITH lmt as (SELECT id, date, code, description, amount, balance, label_id FROM transactions WHERE strftime('%Y-%m', date) = (SELECT strftime('%Y-%m', date) FROM transactions order by date desc limit 1) ORDER BY id)

		SELECT t.*, l.name as label from lmt t LEFT JOIN label l on l.id = t.label_id WHERE description like $1 OR amount = $2 OR code = $2 OR balance = $2 OR $2 = '';`, "%"+search+"%", search)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {

		var t Transaction

		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Code,
			&t.Description,
			&t.Amount,
			&t.Balance,
			&t.LabelID,
			&t.Label,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}

func FindTransactionWithLabel(id int64, db *sql.DB) (Transaction, error) {

	var t Transaction

	err := db.QueryRow("SELECT t.*, l.name as label from transactions t LEFT JOIN label l on l.id = t.label_id where t.id = ?;", id).Scan(
		&t.ID,
		&t.Date,
		&t.Code,
		&t.Description,
		&t.Amount,
		&t.Balance,
		&t.LabelID,
		&t.Label,
	)

	if err != nil {
		return t, err
	}

	return t, nil
}
