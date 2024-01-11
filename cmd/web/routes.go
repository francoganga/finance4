package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./public"))

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		a.templates.Render("home.html", w, nil)
	})
	r.Get("/dashboard", a.Dashboard)

	r.Post("/file", a.HandleFile)

	return r
}

