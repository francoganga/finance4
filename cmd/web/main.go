package main

import (
	"database/sql"
	financeLogger "finance/internal/logger"
	"finance/internal/models"
	"finance/migrations"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/francoganga/pongoe"
	"github.com/pressly/goose/v3"
	sqldblogger "github.com/simukti/sqldb-logger"
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

	serveCmd := flag.NewFlagSet("serve", flag.ExitOnError)
	serveCmd.IntVar(&cfg.port, "port", 4000, "API server port")
	serveCmd.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		serveCmd.Parse(os.Args[2:])
		serve(cfg)
	case "migrate":
		migrate()
	default:
		log.Fatalf("unknown command %q", os.Args[1])
	}

}

func serve(cfg config) {

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	tts := pongoe.LoadTemplates("templates")

	db, err := sql.Open("sqlite", "file:app.db?cache=shared")

	if err != nil {
		log.Fatal(err)
	}

	lh := financeLogger.New(os.Stdout, nil)

	loggerAdapter := financeLogger.NewSQLLogger(slog.New(lh))

	db = sqldblogger.OpenDriver("file:app.db?cache=shared", db.Driver(), loggerAdapter)

	err = db.Ping()

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

func migrate() {
	goose.SetBaseFS(migrations.GetMigrations())

	if err := goose.SetDialect("sqlite"); err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite", "file:app.db?cache=shared")

	if err != nil {
		panic(err)
	}

	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}

	fmt.Println("migrated db successfully")
}

func usage() {
	fmt.Fprint(os.Stderr, "usage: <exe> [-port port] [-env env] [-db-dsn dsn]\n")
}

