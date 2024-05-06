package main

import (
	"database/sql"
	"encoding/json"
	"finance/internal/models"
	"finance/internal/parser"
	"finance/internal/services"
	"finance/internal/types"
	"finance/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/francoganga/ulari"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/samber/lo"
)

func (a *application) HandleFile(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		a.errorResponse(w, r, 500, "Error")
	}

	files, ok := r.MultipartForm.File["files"]

	if !ok {
		a.logger.Fatal("no files")
	}

	for _, file := range files {
		matches, err := utils.GetMatchesFromFile(file)

		if err != nil {
			a.errorResponse(w, r, 500, "Error")
			return
		}

		for _, line := range matches {
			p := parser.New(line)

			consu, err := p.Parse()
			if err != nil {
				a.errorResponse(w, r, 500, "Error")
				return
			}

			pd, err := time.Parse("02/01/06", consu.Date)

			if err != nil {
				a.logError(r, err)
				a.errorResponse(w, r, 500, "Error")
				return
			}

			err = a.queries.CreateTransaction(r.Context(), models.CreateTransactionParams{
				Date:        *types.NewDate(pd),
				Code:        sql.NullString{String: consu.Code, Valid: true},
				Description: consu.Description,
				Amount:      int64(consu.Amount),
				Balance:     int64(consu.Balance),
			})

			if err != nil {
				a.errorResponse(w, r, 500, "Error")
				return
			}
		}
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{"status": "ok"})
}

func (a *application) Dashboard(w http.ResponseWriter, r *http.Request) {

	search := r.URL.Query().Get("search")

	if hx := r.Header.Get("Hx-Request"); hx == "true" {
		lmt, err := services.SearchLastMonthTransactions(search, a.db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Printf("len(lmt)=%#v\n", len(lmt))

		err = a.templates.Render("transaction/_last_month_transactions.html", w, pongo2.Context{
			"transactions": lmt,
		})

		if err != nil {
			fmt.Println("error")
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	type period struct {
		Link string
		Text string
	}

	var periods []period

	rows, err := a.db.Query("select distinct strftime('%Y-%m', date) as m from transactions order by date desc")
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		var val string

		err := rows.Scan(&val)
		if err != nil {
			a.errorResponse(w, r, 500, err.Error())
			return
		}

		p, err := time.Parse("2006-01", val)
		if err != nil {
			a.errorResponse(w, r, 500, err.Error())
			return
		}

		period := period{
			Text: p.Format("2006 - Jan"),
			Link: p.Format("2006-01"),
		}

		periods = append(periods, period)

	}
	// latest transactions
	//select * from transactions where strftime('%Y-%m', date) = (select strftime('%Y-%m', date) from transactions order by date desc limit 1) order by date desc limit 1

	lts, err := services.SearchLastMonthTransactions(search, a.db)
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	salary := lo.MaxBy(lts, func(a models.Transaction, b models.Transaction) bool {
		return a.Amount > b.Amount
	})

	expenses := lo.Filter(lts, func(i models.Transaction, _ int) bool {
		return i.Amount < 0
	})

	totalExpenses := lo.SumBy(expenses, func(i models.Transaction) int64 {
		return i.Amount
	})

	remainingAmount := salary.Amount + totalExpenses

	a.templates.Render("dashboard.html", w, pongo2.Context{
		"periods":         periods,
		"transactions":    lts,
		"balance":         lts[len(lts)-1].Balance,
		"salary":          salary.Amount / 100,
		"expenses":        expenses,
		"totalExpenses":   totalExpenses / 100,
		"remainingAmount": remainingAmount / 100,
	})
}

func (a *application) NewTransaction(w http.ResponseWriter, r *http.Request) {

	form := ulari.NewFormData()
	now := time.Now().Format("2006-01-02")

	form.Add(ulari.DateField("date", now, "p-2", "border", "border-gray-300"))
	form.Add(ulari.StringField("name", "p-2", "border", "border-red-500", "focus:border-sky-400"))

	if err := a.templates.Render("transaction/new.html", w, pongo2.Context{
		"form": form,
	}); err != nil {
		a.errorResponse(w, r, 500, err.Error())
	}

}

func (a *application) Lmt(w http.ResponseWriter, r *http.Request) {
	lts, err := a.queries.LastMonthTransactions(r.Context())
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(lts)
}

func (a *application) MonthOverview(w http.ResponseWriter, r *http.Request) {

	period := chi.URLParam(r, "period")

	if period == "" {
		a.errorResponse(w, r, 500, "empty period")
	}

	query := "select * from transactions where strftime('%Y-%m', date) = ?"

	rows, err := a.db.Query(query, period)
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	defer rows.Close()

	var transactions []models.Transaction

	for rows.Next() {
		var transaction models.Transaction

		err := rows.Scan(
			&transaction.ID,
			&transaction.Date,
			&transaction.Code,
			&transaction.Description,
			&transaction.Amount,
			&transaction.Balance,
		)

		if err != nil {
			a.errorResponse(w, r, 500, err.Error())
			return
		}

		transactions = append(transactions, transaction)
	}

	overview := services.GenerateOverview(transactions)

	a.templates.Render("month_overview.html", w, pongo2.Context{
		"overview": overview,
	})
}
