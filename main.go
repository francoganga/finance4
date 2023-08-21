package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
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

	r.Get("/hw", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, map[string]string{
			"message": "Hello, World!",
		})
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

