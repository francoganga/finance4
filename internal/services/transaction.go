package services

import (
	"database/sql"
	"finance/internal/models"
	"fmt"
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

	fmt.Printf("%d - %d\n", salary, totalExpenses)

	return Overview{
		Salary:          salary / 100,
		Expenses:        totalExpenses / 100,
		RemainingAmount: (salary + totalExpenses) / 100,
	}
}

func SearchTransactions(search string, month string, db *sql.DB) ([]Transaction, error) {

	rows, err := db.Query(`WITH lmt as (SELECT id, date, code, description, amount, balance, label_id FROM transactions WHERE (strftime('%Y-%m', date) = $3) OR ($3 = '')  ORDER BY id)

		SELECT t.*, l.name as label from lmt t LEFT JOIN label l on l.id = t.label_id WHERE description like $1 OR amount = $2 OR code = $2 OR balance = $2 OR $2 = '';`, "%"+search+"%", search, month)

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
