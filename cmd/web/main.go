package main

import (
	"database/sql"
	"finance/internal/models"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/francoganga/pongoe"
	_ "modernc.org/sqlite"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config    config
	logger    *log.Logger
	db        *sql.DB
	queries   *models.Queries
	templates *pongoe.Templates
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	tts := pongoe.LoadTemplates("templates")

	db, err := sql.Open("sqlite", "file:app.db?cache=shared")

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("opened db")

	queries := models.New(db)

	app := &application{
		config:    cfg,
		logger:    logger,
		db:        db,
		queries:   queries,
		templates: tts,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)

	err = srv.ListenAndServe()

	logger.Fatal(err)
}

