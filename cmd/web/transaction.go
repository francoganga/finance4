package main

import (
	"database/sql"
	"finance/internal/models"
	"finance/internal/parser"
	"finance/internal/utils"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

func (a *application) HandleFile(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		a.errorResponse(w, r, 500, "Error")
	}

	files := r.MultipartForm.File["files"]

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

			err := a.queries.CreateTransaction(r.Context(), models.CreateTransactionParams{
				Date:        consu.Date,
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

