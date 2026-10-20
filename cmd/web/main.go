package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handlerHome)
	mux.HandleFunc("GET /snippet/view/{id}", handlerSnippetView)
	mux.HandleFunc("GET /snippet/create", handlerSnippetCreate)
	mux.HandleFunc("POST /snippet/create", handlerSnippetCreatePost)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	port := "4000"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on: http://localhost:%s\n", port)
	log.Fatal(srv.ListenAndServe())
}
