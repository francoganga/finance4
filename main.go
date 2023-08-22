package main

import (
	"database/sql"
	"finance/pkg/models"
	"finance/pkg/parser"
	"finance/pkg/utils"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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

	r.Post("/file", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseMultipartForm(32 << 20)

		if err != nil {
			w.Write([]byte("Error: " + err.Error()))
			return
		}

		files := r.MultipartForm.File["files"]

		for _, file := range files {
			matches, err := utils.GetMatchesFromFile(file)

			if err != nil {
				w.Write([]byte("Error: " + err.Error()))
				return
			}

			db, err := sql.Open("sqlite3", "file:app.db?cache=shared&mode=memory")

			fmt.Println("opened db")

			if err != nil {
				panic(err)
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

					w.Write([]byte("Error: " + msg))
					return
				}

				err := queries.CreateTransaction(r.Context(), models.CreateTransactionParams{
					Date:        consu.Date,
					Code:        sql.NullString{String: consu.Code, Valid: true},
					Description: consu.Description,
					Amount:      int64(consu.Amount),
				})

				if err != nil {
					w.Write([]byte("Error: " + err.Error()))
					return
				}

			}
		}

		w.Write([]byte("success"))
	})
	templ, err := template.New("base.gohtml").ParseFiles("templates/base.gohtml")

	if err != nil {
		panic(err)
	}

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {

		type TemplGlobals struct {
			Mode string
		}

		err := templ.Execute(w, TemplGlobals{
			Mode: appEnv,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	})
	http.ListenAndServe(":3000", r)
}

