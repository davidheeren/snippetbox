package main

import "net/http"

func (cfg *config) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", cfg.handlerHome)
	mux.HandleFunc("GET /snippets/view/{id}", cfg.handlerSnippetView)
	mux.HandleFunc("GET /snippets/create", cfg.handlerSnippetCreate)
	mux.HandleFunc("POST /snippets/create", cfg.handlerSnippetCreatePost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return mux
}
