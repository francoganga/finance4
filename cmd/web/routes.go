package main

import (
	assets "finance/public"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.FS(assets.GetAssets()))

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		a.templates.Render("home.html", w, nil)
	})
	r.Get("/dashboard", a.Dashboard)
	r.Get("/transactions/{id}/edit", a.EditTransaction)
	r.Patch("/transactions/{id}/addLabel", a.AddLabel)
	r.Get("/month_overview/{period}", a.MonthOverview)
	r.Post("/file", a.HandleFile)

	api := chi.NewRouter()
	api.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
		return
	})

	api.Get("/lmt", a.Lmt)
	api.Get("/transactions", a.ApiTransactions)

	r.Mount("/api", api)

	return r
}

