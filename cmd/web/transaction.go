package main

import (
	"database/sql"
	"encoding/json"
	"finance/internal/models"
	"finance/internal/parser"
	"finance/internal/types"
	"finance/internal/utils"
	"net/http"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/render"
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

	var periods []time.Time

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

		periods = append(periods, p)

	}
	// latest transactions
	//select * from transactions where strftime('%Y-%m', date) = (select strftime('%Y-%m', date) from transactions order by date desc limit 1) order by date desc limit 1

	lts, err := a.queries.LastMonthTransactions(r.Context())
	if err != nil {
		a.errorResponse(w, r, 500, err.Error())
		return
	}

	a.templates.Render("dashboard.html", w, pongo2.Context{
		"periods": periods,
		"lts":     lts,
		"balance": lts[len(lts)-1].Balance,
	})
}

func (a *application) NewTransaction(w http.ResponseWriter, r *http.Request) {

	cd := time.Now().Format("2006-01-02")

	if err := a.templates.Render("transaction/new.html", w, pongo2.Context{
		"date": cd,
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

	json.NewEncoder(w).Encode(lts)
}

