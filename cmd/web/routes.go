package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (cfg *config) routes() http.Handler {
	mux := http.NewServeMux()

	dynamic := alice.New(cfg.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(cfg.handlerHome))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(cfg.handlerSnippetView))
	mux.Handle("GET /snippets/create", dynamic.ThenFunc(cfg.handlerSnippetCreate))
	mux.Handle("POST /snippets", dynamic.ThenFunc(cfg.handlerSnippetPost))
	mux.Handle("DELETE /snippets/{id}", dynamic.ThenFunc(cfg.handlerSnippetDelete))

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	standardMiddleware := alice.New(cfg.middlewareRecoverPanic, cfg.middlewareLogRequest, middlewareCommonHeaders)
	return standardMiddleware.Then(mux)
}
