package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (cfg *config) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", cfg.handlerHome)
	mux.HandleFunc("GET /snippets/view/{id}", cfg.handlerSnippetView)
	mux.HandleFunc("GET /snippets/create", cfg.handlerSnippetCreate)
	mux.HandleFunc("POST /snippets", cfg.handlerSnippetPost)
	mux.HandleFunc("DELETE /snippets/{id}", cfg.handlerSnippetDelete)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	standardMiddleware := alice.New(cfg.middlewareRecoverPanic, cfg.middlewareLogRequest, middlewareCommonHeaders)
	return standardMiddleware.Then(mux)
}
