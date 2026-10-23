package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/davidheeren/snippetbox/internal/database"
	"github.com/go-playground/form/v4"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type arguments struct {
	addr string
}

type config struct {
	logger         *slog.Logger
	dbQueries      *database.Queries
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	var args arguments
	flag.StringVar(&args.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbURL, ok := os.LookupEnv("DB_URL")
	if !ok {
		log.Fatal("DB_URL environment variable must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	templateCache, err := newTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	cfg := &config{
		logger:         logger,
		dbQueries:      database.New(db),
		templateCache:  templateCache,
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
	}

	srv := &http.Server{
		Addr:     args.addr,
		Handler:  cfg.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Serving server", slog.String("url", "http://localhost"+args.addr))
	log.Fatal(srv.ListenAndServe())
}
