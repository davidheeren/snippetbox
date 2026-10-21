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

	srv := &http.Server{
		Addr:    args.addr,
		Handler: cfg.routes(),
	}

	logger.Info("Serving server", slog.String("url", "http://localhost"+args.addr))
	log.Fatal(srv.ListenAndServe())
}
