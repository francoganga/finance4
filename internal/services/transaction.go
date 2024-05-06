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

func GenerateOverview(transactions []models.Transaction) Overview {

	salaryTransactions := lo.Filter(transactions, func(i models.Transaction, _ int) bool {
		return strings.Contains(i.Description, "haberes")
	})

	salary := lo.SumBy(salaryTransactions, func(i models.Transaction) int64 {
		return i.Amount
	})

	expenses := lo.Filter(transactions, func(i models.Transaction, _ int) bool {
		return i.Amount < 0
	})

	totalExpenses := lo.SumBy(expenses, func(i models.Transaction) int64 {
		return i.Amount
	})

	fmt.Printf("%d - %d\n", salary, totalExpenses)

	return Overview{
		Salary:          salary / 100,
		Expenses:        totalExpenses / 100,
		RemainingAmount: (salary + totalExpenses) / 100,
	}
}

func SearchLastMonthTransactions(search string, db *sql.DB) ([]models.Transaction, error) {

	rows, err := db.Query(`WITH lmt as (SELECT id, date, code, description, amount, balance FROM transactions WHERE strftime('%Y-%m', date) = (SELECT strftime('%Y-%m', date) FROM transactions order by date desc limit 1) ORDER BY id)

		SELECT * from lmt WHERE description like $1 OR amount = $2 OR code = $2 OR balance = $2 OR $2 = '';`, "%"+search+"%", search)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {

		var t models.Transaction

		err := rows.Scan(
			&t.ID,
			&t.Date,
			&t.Code,
			&t.Description,
			&t.Amount,
			&t.Balance,
		)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t)
	}

	return transactions, nil
}
