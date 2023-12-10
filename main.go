package main

import (
	"database/sql"
	"finance/pkg/models"
	"finance/pkg/parser"
	"finance/pkg/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/flosch/pongo2/v6"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	_ "modernc.org/sqlite"
)

type Handler func(http.ResponseWriter, *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, err)
		return
	}
}

func main() {

	appEnv := os.Getenv("APP_ENV")

	if appEnv == "" {
		appEnv = "dev"
	}

	fmt.Printf("appEnv=%v\n", appEnv)

	r := chi.NewRouter()
	fs := http.FileServer(http.Dir("./public"))

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Method("POST", "/file", Handler(func(w http.ResponseWriter, r *http.Request) error {
		err := r.ParseMultipartForm(32 << 20)

		if err != nil {
			return err
		}

		files := r.MultipartForm.File["files"]

		for _, file := range files {
			matches, err := utils.GetMatchesFromFile(file)
			fmt.Printf("matches_len=%v\n", len(matches))

			if err != nil {
				return err
			}

			db, err := sql.Open("sqlite", "file:app.db?cache=shared")

			fmt.Println("opened db")

			if err != nil {
				return err
			}

			queries := models.New(db)

			for _, line := range matches {
				p := parser.FromInput(line)

				consu := p.ParseConsumo()

				if len(p.Errors()) > 0 {

					msg := ""

					for _, e := range p.Errors() {
						msg += e
					}

					return err
				}

				err := queries.CreateTransaction(r.Context(), models.CreateTransactionParams{
					Date:        consu.Date,
					Code:        sql.NullString{String: consu.Code, Valid: true},
					Description: consu.Description,
					Amount:      int64(consu.Amount),
				})

				if err != nil {
					return err
				}
			}
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, map[string]string{"status": "ok"})

		return nil
	}))

	templ := pongo2.Must(pongo2.FromFile("templates/base.html"))

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		err := templ.ExecuteWriter(pongo2.Context{"mode": appEnv}, w)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})
	http.ListenAndServe(":3000", r)
}

