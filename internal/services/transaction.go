package services

import (
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

