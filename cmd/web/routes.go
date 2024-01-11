package main

import (
	"finance/templates"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *application) routes() *chi.Mux {
	r := chi.NewRouter()

	fs := http.FileServer(http.Dir("./public"))

	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	r.Get("/", templ.Handler(templates.Base()).ServeHTTP)

	r.Post("/file", a.HandleFile)

	return r
}

