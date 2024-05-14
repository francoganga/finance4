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
	"strconv"
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
			"search":       search,
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

	salary := lo.MaxBy(lts, func(a services.Transaction, b services.Transaction) bool {
		return a.Amount > b.Amount
	})

	expenses := lo.Filter(lts, func(i services.Transaction, _ int) bool {
		return i.Amount < 0
	})

	totalExpenses := lo.SumBy(expenses, func(i services.Transaction) int64 {
		return i.Amount
	})

	remainingAmount := salary.Amount + totalExpenses

	labels, err := a.queries.ListLabels(r.Context())
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	a.templates.Render("dashboard.html", w, pongo2.Context{
		"labels":          labels,
		"periods":         periods,
		"transactions":    lts,
		"balance":         lts[len(lts)-1].Balance,
		"salary":          salary.Amount / 100,
		"expenses":        expenses,
		"totalExpenses":   totalExpenses / 100,
		"remainingAmount": remainingAmount / 100,
		"search":          search,
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

func (a *application) EditTransaction(w http.ResponseWriter, r *http.Request) {

	param := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	transaction, err := services.FindTransactionWithLabel(id, a.db)
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	labels, err := a.queries.ListLabels(r.Context())
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	a.templates.Render("transaction/_edit_transaction.html", w, pongo2.Context{
		"labels":      labels,
		"transaction": transaction,
	})
}

func (a *application) AddLabel(w http.ResponseWriter, r *http.Request) {

	param := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(param, 10, 64)

	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	transaction, err := services.FindTransactionWithLabel(id, a.db)

	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		default:
			a.errorResponse(w, r, 500, err.Error())
			return
		}
	}

	err = r.ParseForm()

	if len(r.Form["label_id"]) < 1 {
		a.errorResponse(w, r, 400, "Bad request")
		return
	}

	labelId, err := strconv.ParseInt(r.Form["label_id"][0], 10, 64)

	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	err = a.queries.AddLabelToTransaction(r.Context(), models.AddLabelToTransactionParams{
		LabelID: sql.NullInt64{Valid: true, Int64: labelId},
		ID:      id,
	})

	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	transaction, err = services.FindTransactionWithLabel(id, a.db)
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	a.templates.Render("transaction/_transaction.partial.html", w, pongo2.Context{
		"transaction": transaction,
	})

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

func (a *application) ApiTransactions(w http.ResponseWriter, r *http.Request) {

	opts := parseOpts(r.URL.Query())

	transactions, err := services.SearchTransactions(opts, a.db)

	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(transactions)
}

func (a *application) MonthOverview(w http.ResponseWriter, r *http.Request) {

	opts := parseOpts(r.URL.Query())

	transactions, err := services.SearchTransactions(opts, a.db)

	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	if hx := r.Header.Get("Hx-Request"); hx == "true" {

		err = a.templates.Render("transaction/_last_month_transactions.html", w, pongo2.Context{
			"transactions": transactions,
			"search":       opts.Search,
			"period":       opts.Period,
		})

		return
	}

	overview := services.GenerateOverview(transactions)

	labels, err := a.queries.ListLabels(r.Context())
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	a.templates.Render("month_overview.html", w, pongo2.Context{
		"overview":     overview,
		"transactions": transactions,
		"search":       opts.Search,
		"period":       opts.Period,
		"labels":       labels,
	})
}
