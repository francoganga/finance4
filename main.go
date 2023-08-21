package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"portal/pkg/utils"
	"strings"

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

			for _, line := range matches {
				fmt.Println(strings.Replace(line, "\n", "\\n", -1))
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

