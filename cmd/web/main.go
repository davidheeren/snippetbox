package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
)

type arguments struct {
	addr string
}

type config struct {
	logger *slog.Logger
}

func main() {
	var args arguments
	flag.StringVar(&args.addr, "addr", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	cfg := &config{
		logger: logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", cfg.handlerHome)
	mux.HandleFunc("GET /snippet/view/{id}", cfg.handlerSnippetView)
	mux.HandleFunc("GET /snippet/create", cfg.handlerSnippetCreate)
	mux.HandleFunc("POST /snippet/create", cfg.handlerSnippetCreatePost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	srv := &http.Server{
		Addr:    args.addr,
		Handler: mux,
	}

	logger.Info("Serving server", slog.String("url", "http://localhost"+args.addr))
	log.Fatal(srv.ListenAndServe())
}
