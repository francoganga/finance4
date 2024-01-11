package main

import (
	"database/sql"
	"finance/internal/models"
	"finance/internal/parser"
	"finance/internal/utils"
	"fmt"
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
		fmt.Printf("matches_len=%v\n", len(matches))

		if err != nil {
			a.errorResponse(w, r, 500, "Error")
			return
		}

		for _, line := range matches {
			p := parser.FromInput(line)

			consu := p.ParseConsumo()

			if len(p.Errors()) > 0 {

				msg := ""

				for _, e := range p.Errors() {
					msg += e
				}

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
				Date:        pd.Format("2006-01-02"),
				Code:        sql.NullString{String: consu.Code, Valid: true},
				Description: consu.Description,
				Amount:      int64(consu.Amount),
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

	a.templates.Render("dashboard.html", w, pongo2.Context{
		"periods": periods,
	})
}

